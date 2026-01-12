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
)

var RetryTransientErrorTimeout = 1 * time.Minute

const MaxTransientRetries = 5

type GVClient struct {
	Main      *GVConnector
	UserLogin *bridgev2.UserLogin
	Client    *libgv.Client
	LoggedIn  bool

	lastEvents           map[string]time.Time
	wakeupMessageFetcher chan struct{}
	fetchEventsLimiter   *rate.Limiter
	fetchEventsLock      sync.Mutex

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
		LoggedIn:  len(cookies) > 0,

		lastEvents:           make(map[string]time.Time),
		wakeupMessageFetcher: make(chan struct{}),
		fetchEventsLimiter:   rate.NewLimiter(rate.Every(MinRefreshInterval), MinRefreshBurstCount),
		contactCache:         make(map[string]*ProcessedContact),
	}
	lgvClient.EventHandler = gvClient.handleRealtimeEvent
	login.Client = gvClient
	return nil
}

var _ bridgev2.NetworkAPI = (*GVClient)(nil)

func (gc *GVClient) Connect(ctx context.Context) {
	_, _ = gc.Main.Bridge.GetGhostByID(ctx, "")
	_, err := gc.Client.GetAccount(ctx)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to get account on connect")
		if ctx.Err() == nil {
			if libgv.IsAuthError(err) {
				gc.UserLogin.BridgeState.Send(status.BridgeState{
					StateEvent: status.StateBadCredentials,
					Error:      GVBadCredentials,
					Info:       map[string]any{"go_error": err.Error()},
				})
			} else {
				gc.UserLogin.BridgeState.Send(status.BridgeState{
					StateEvent: status.StateUnknownError,
					Error:      GVConnectError,
					Info:       map[string]any{"go_error": err.Error()},
				})
			}
		}
		return
	}
	go gc.connectRealtime()
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
		err := gc.Client.RunRealtimeChannel(ctx)
		if errors.Is(err, context.Canceled) || ctx.Err() != nil {
			log.Debug().Err(err).Msg("Realtime channel disconnected with context canceled")
		} else if err == nil {
			log.Warn().Msg("Realtime channel disconnected without error")
			transientRetries = 0
		} else if libgv.IsAuthError(err) {
			log.Err(err).Msg("Realtime channel disconnected with auth error")
			gc.UserLogin.BridgeState.Send(status.BridgeState{
				StateEvent: status.StateBadCredentials,
				Error:      GVBadCredentials,
				Info:       map[string]any{"go_error": err.Error()},
			})
			return
		} else if errors.Is(err, libgv.ErrTooManyUnknownSID) {
			log.Err(err).Msg("Realtime channel failed with too many unknown SID errors")
			gc.UserLogin.BridgeState.Send(status.BridgeState{
				StateEvent: status.StateUnknownError,
				Error:      GVTooManyRetries,
				Info:       map[string]any{"go_error": err.Error()},
			})
			return
		} else {
			transientRetries++
			if transientRetries > MaxTransientRetries {
				log.Err(err).Int("retries", transientRetries).Msg("Realtime channel exceeded max transient retries")
				gc.UserLogin.BridgeState.Send(status.BridgeState{
					StateEvent: status.StateUnknownError,
					Error:      GVTooManyRetries,
					Info:       map[string]any{"go_error": err.Error()},
				})
				return
			}
			log.Err(err).Int("retry", transientRetries).Msg("Realtime channel disconnected with transient error")
			gc.UserLogin.BridgeState.Send(status.BridgeState{
				StateEvent: status.StateTransientDisconnect,
				Error:      GVDisconnected,
				Info:       map[string]any{"go_error": err.Error()},
			})
			select {
			case <-ctx.Done():
			case <-time.After(RetryTransientErrorTimeout):
				log.Info().Int("retry", transientRetries).Msg("Retrying connection after transient error")
				continue
			}
		}
		return
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
	return gc.LoggedIn
}

func (gc *GVClient) LogoutRemote(ctx context.Context) {
	// TODO is logging out possible?
}
