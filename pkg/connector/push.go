// mautrix-telegram - A Matrix-Telegram puppeting bridge.
// Copyright (C) 2025 Beeper
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

	"maunium.net/go/mautrix/bridgev2"
)

var _ bridgev2.PushableNetworkAPI = (*GVClient)(nil)

func (g *GVClient) RegisterPushNotifications(ctx context.Context, pushType bridgev2.PushType, token string) error {
	return nil
}

func (g *GVClient) GetPushConfigs() *bridgev2.PushConfig {
	return &bridgev2.PushConfig{
		FCM: &bridgev2.FCMPushConfig{
			SenderID: "301778431048",
		},
	}
}
