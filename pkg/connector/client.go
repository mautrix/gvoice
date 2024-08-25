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

	"golang.org/x/time/rate"
	"maunium.net/go/mautrix/bridge/status"
	"maunium.net/go/mautrix/bridgev2"

	"go.mau.fi/mautrix-gvoice/pkg/libgv"
)

type GVClient struct {
	Main      *GVConnector
	UserLogin *bridgev2.UserLogin
	Client    *libgv.Client
	LoggedIn  bool

	lastEvents           map[string]time.Time
	wakeupMessageFetcher chan struct{}
	fetchEventsLimiter   *rate.Limiter
	fetchEventsLock      sync.Mutex

	stopRealtime atomic.Pointer[context.CancelFunc]
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
	}
	lgvClient.EventHandler = gvClient.handleRealtimeEvent
	login.Client = gvClient
	return nil
}

var _ bridgev2.NetworkAPI = (*GVClient)(nil)

func (gc *GVClient) Connect(ctx context.Context) error {
	_, _ = gc.Main.Bridge.GetGhostByID(ctx, "")
	_, err := gc.Client.GetAccount(ctx)
	if err != nil {
		// TODO split out bad credentials
		gc.UserLogin.BridgeState.Send(status.BridgeState{
			StateEvent: status.StateUnknownError,
			Error:      "gv-connect-error",
			Info:       map[string]any{"go_error": err.Error()},
		})
		return nil
	}
	go gc.connectRealtime()
	return nil
}

func (gc *GVClient) connectRealtime() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go gc.fetchNewMessagesLoop(gc.UserLogin.Log.With().Str("component", "fetch messages loop").Logger().WithContext(ctx))
	log := gc.UserLogin.Log.With().Str("component", "realtime channel").Logger()
	gc.stopRealtime.Store(&cancel)
	ctx = log.WithContext(ctx)
	err := gc.Client.RunRealtimeChannel(ctx)
	if errors.Is(err, context.Canceled) {
		log.Debug().Msg("Realtime channel disconnected with context canceled")
	} else if err == nil {
		log.Warn().Msg("Realtime channel disconnected without error")
	} else {
		log.Err(err).Msg("Realtime channel disconnected with unknown error")
		gc.UserLogin.BridgeState.Send(status.BridgeState{
			StateEvent: status.StateUnknownError,
			Error:      "gv-realtime-unknown-error",
			Info:       map[string]any{"go_error": err.Error()},
		})
	}
}

func (gc *GVClient) Disconnect() {
	if stop := gc.stopRealtime.Swap(nil); stop != nil {
		(*stop)()
	}
}

func (gc *GVClient) IsLoggedIn() bool {
	return gc.LoggedIn
}

func (gc *GVClient) LogoutRemote(ctx context.Context) {
	// TODO is logging out possible?
}
