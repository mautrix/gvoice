// mautrix-gvoice - A Matrix-Google Voice puppeting bridge.
// Copyright (C) 2024 Tulir Asokan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package connector

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/status"

	"go.mau.fi/mautrix-gvoice/pkg/libgv"
)

const (
	GVNotLoggedIn    status.BridgeStateErrorCode = "gv-not-logged-in"
	GVBadCredentials status.BridgeStateErrorCode = "gv-bad-credentials"
	GVConnectError   status.BridgeStateErrorCode = "gv-connect-error"
	GVDisconnected   status.BridgeStateErrorCode = "gv-transient-disconnect"
	GVRealtimeError  status.BridgeStateErrorCode = "gv-realtime-error"
	GVTooManyRetries status.BridgeStateErrorCode = "gv-too-many-retries"
	GVRateLimited    status.BridgeStateErrorCode = "gv-rate-limited"
)

func init() {
	status.BridgeStateHumanErrors.Update(status.BridgeStateErrorMap{
		GVNotLoggedIn:    "You're not logged into Google Voice",
		GVBadCredentials: "Your Google Voice session expired, please log in again",
		GVConnectError:   "Failed to connect to Google Voice, retrying",
		GVDisconnected:   "Disconnected from Google Voice, reconnecting",
		GVRealtimeError:  "Error in the Google Voice realtime connection",
		GVTooManyRetries: "Repeatedly failed to reconnect to Google Voice, retrying",
		GVRateLimited:    "Google Voice is rate limiting connections, retrying",
	})
}

var RetryTransientErrorTimeout = 1 * time.Minute
var RetryTransientErrorLongTimeout = 10 * time.Minute

const (
	MaxTransientRetries     = 5
	HealthyRealtimeDuration = 5 * time.Minute
)

type GVClient struct {
	Main      *GVConnector
	UserLogin *bridgev2.UserLogin
	Client    *libgv.Client
	LoggedIn  atomic.Bool

	lastEvents           map[string]time.Time
	wakeupMessageFetcher chan struct{}
	fetchEventsLimiter   *rate.Limiter
	fetchEventsLock      sync.Mutex
	threadSyncStarted    atomic.Bool

	contactCache     map[string]*ProcessedContact
	contactCacheLock sync.Mutex

	stopRealtime     atomic.Pointer[context.CancelFunc]
	requestSignature atomic.Pointer[requestSignatureFunc]
	stopWait         sync.WaitGroup
}

func (gv *GVConnector) LoadUserLogin(ctx context.Context, login *bridgev2.UserLogin) error {
	cookies := login.Metadata.(*UserLoginMetadata).Cookies
	lgvClient := libgv.NewClient(cookies)
	gvClient := &GVClient{
		Main:      gv,
		Client:    lgvClient,
		UserLogin: login,

		lastEvents:           make(map[string]time.Time),
		wakeupMessageFetcher: make(chan struct{}),
		fetchEventsLimiter:   rate.NewLimiter(rate.Every(MinRefreshInterval), MinRefreshBurstCount),
		contactCache:         make(map[string]*ProcessedContact),
	}
	gvClient.LoggedIn.Store(len(cookies) > 0)
	lgvClient.EventHandler = gvClient.handleRealtimeEvent
	login.Client = gvClient
	return nil
}

var _ bridgev2.NetworkAPI = (*GVClient)(nil)

func (gc *GVClient) Connect(ctx context.Context) {
	gc.UserLogin.BridgeState.Send(status.BridgeState{StateEvent: status.StateConnecting})
	_, _ = gc.Main.Bridge.GetGhostByID(ctx, "")
	_, err := gc.Client.GetAccount(ctx)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to get account on connect")
		if ctx.Err() == nil && gc.handleConnectError(err) {
			retryCtx, cancel := context.WithCancel(gc.Main.Bridge.BackgroundCtx)
			gc.stopRealtime.Store(&cancel)
			go gc.retryConnect(retryCtx, cancel)
		}
		return
	}
	go gc.connectRealtime()
}

// handleConnectError sends the bridge state for a failed account check and returns whether to retry.
func (gc *GVClient) handleConnectError(err error) bool {
	if libgv.IsAuthError(err) {
		gc.sendBadCredentials(err)
		return false
	}
	code := GVConnectError
	if isRateLimitError(err) {
		code = GVRateLimited
	}
	gc.UserLogin.BridgeState.Send(status.BridgeState{
		StateEvent: status.StateTransientDisconnect,
		Error:      code,
		Info:       map[string]any{"go_error": err.Error()},
	})
	return true
}

func (gc *GVClient) retryConnect(ctx context.Context, cancel context.CancelFunc) {
	defer cancel()
	log := gc.UserLogin.Log.With().Str("action", "retry connect").Logger()
	ctx = log.WithContext(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(RetryTransientErrorTimeout):
		}
		_, err := gc.Client.GetAccount(ctx)
		if ctx.Err() != nil {
			return
		} else if err == nil {
			gc.connectRealtime()
			return
		}
		log.Err(err).Msg("Failed to get account on connect retry")
		if !gc.handleConnectError(err) {
			return
		}
	}
}

func (gc *GVClient) checkSessionAfterAuthError(ctx context.Context) {
	if !gc.LoggedIn.Load() {
		return
	}
	callCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err := gc.Client.GetAccount(callCtx)
	// A relogin while GetAccount was running replaces the login's client, and the new one must keep its state.
	if libgv.IsAuthError(err) && ctx.Err() == nil && gc.UserLogin.Client == gc {
		gc.sendBadCredentials(err)
	}
}

func (gc *GVClient) sendBadCredentials(err error) {
	gc.LoggedIn.Store(false)
	gc.UserLogin.BridgeState.Send(status.BridgeState{
		StateEvent: status.StateBadCredentials,
		Error:      GVBadCredentials,
		UserAction: status.UserActionRelogin,
		Info:       map[string]any{"go_error": err.Error()},
	})
	if stop := gc.stopRealtime.Load(); stop != nil {
		(*stop)()
	}
}

func (gc *GVClient) connectRealtime() {
	ctx, cancel := context.WithCancel(gc.Main.Bridge.BackgroundCtx)
	gc.stopRealtime.Store(&cancel)
	defer cancel()

	gc.loadInitialContacts(gc.UserLogin.Log.With().Str("action", "load initial contacts").Logger().WithContext(ctx))
	go gc.fetchNewMessagesLoop(gc.UserLogin.Log.With().Str("component", "fetch messages loop").Logger().WithContext(ctx))
	go gc.runElectron(ctx)

	log := gc.UserLogin.Log.With().Str("component", "realtime channel").Logger()
	ctx = log.WithContext(ctx)
	transientRetries := 0
	for {
		startedAt := time.Now()
		err := gc.Client.RunRealtimeChannel(ctx)
		if time.Since(startedAt) > HealthyRealtimeDuration {
			transientRetries = 0
		}
		retryIn := RetryTransientErrorTimeout
		if errors.Is(err, context.Canceled) || ctx.Err() != nil {
			log.Debug().Err(err).Msg("Realtime channel disconnected with context canceled")
			return
		} else if err == nil {
			log.Warn().Msg("Realtime channel disconnected without error")
			transientRetries = 0
			gc.UserLogin.BridgeState.Send(status.BridgeState{
				StateEvent: status.StateTransientDisconnect,
				Error:      GVDisconnected,
			})
		} else if libgv.IsAuthError(err) {
			log.Err(err).Msg("Realtime channel disconnected with auth error")
			gc.sendBadCredentials(err)
			return
		} else {
			transientRetries++
			code := GVDisconnected
			if transientRetries > MaxTransientRetries {
				code = GVTooManyRetries
				retryIn = RetryTransientErrorLongTimeout
			}
			log.Err(err).
				Int("retry", transientRetries).
				Stringer("retry_in", retryIn).
				Msg("Realtime channel disconnected with transient error")
			gc.UserLogin.BridgeState.Send(status.BridgeState{
				StateEvent: status.StateTransientDisconnect,
				Error:      code,
				Info:       map[string]any{"go_error": err.Error()},
			})
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(retryIn):
			log.Info().Int("retry", transientRetries).Msg("Retrying connection after transient error")
		}
	}
}

func (gc *GVClient) loadInitialContacts(ctx context.Context) {
	resp, err := gc.Client.AutocompleteContacts(ctx, "")
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to load initial contacts")
		return
	}
	gc.contactCacheLock.Lock()
	defer gc.contactCacheLock.Unlock()
	for _, contact := range resp.Results {
		for _, method := range contact.GetPerson().GetContactMethods() {
			e164 := method.GetPhone().GetCanonicalValue()
			if e164 != "" {
				_, alreadyExists := gc.contactCache[e164]
				if !alreadyExists || method.GetDisplayInfo().GetPrimary() {
					gc.contactCache[e164] = processContact(contact.GetPerson())
				}
			}
		}
	}
}

func (gc *GVClient) Disconnect() {
	if stop := gc.stopRealtime.Swap(nil); stop != nil {
		(*stop)()
	}
	gc.stopWait.Wait()
}

func (gc *GVClient) IsLoggedIn() bool {
	return gc.LoggedIn.Load()
}

func (gc *GVClient) LogoutRemote(ctx context.Context) {
	// TODO is logging out possible?
	gc.Disconnect()
}
