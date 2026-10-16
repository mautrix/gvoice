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

	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/event"
)

var generalCaps = &bridgev2.NetworkGeneralCapabilities{
	DisappearingMessages: false,
	AggressiveUpdateInfo: false,
	Provisioning: bridgev2.ProvisioningCapabilities{
		ResolveIdentifier: bridgev2.ResolveIdentifierCapabilities{
			CreateDM: true,
			AnyPhone: true,
		},
	},
}

func (gv *GVConnector) GetCapabilities() *bridgev2.NetworkGeneralCapabilities {
	return generalCaps
}

func (gv *GVConnector) GetBridgeInfoVersion() (info, caps int) {
	return 1, 4
}

var roomCaps = &event.RoomFeatures{
	ID:            "fi.mau.gvoice.capabilities.2026_10_16",
	MaxTextLength: 160,
	File: map[event.MessageType]*event.FileFeatures{
		event.MsgImage: {
			MaxSize:          2 * 1024 * 1024,
			Caption:          event.CapLevelFullySupported,
			MaxCaptionLength: 160,
			MimeTypes: map[string]event.CapabilitySupportLevel{
				"image/png":  event.CapLevelFullySupported,
				"image/jpeg": event.CapLevelFullySupported,
				"image/bmp":  event.CapLevelFullySupported,
				"image/tiff": event.CapLevelFullySupported,
				"image/gif":  event.CapLevelFullySupported,
			},
		},
		event.CapMsgGIF: {
			MaxSize:          2 * 1024 * 1024,
			Caption:          event.CapLevelFullySupported,
			MaxCaptionLength: 160,
			MimeTypes: map[string]event.CapabilitySupportLevel{
				"image/gif": event.CapLevelFullySupported,
			},
		},
	},
	DeleteChat: true,
}

var dmCaps *event.RoomFeatures

func init() {
	dmCaps = roomCaps.Clone()
	dmCaps.ID += "+dm"
	dmCaps.BlockUser = true
	dmCaps.ReportSpam = true
}

func (gc *GVClient) GetCapabilities(ctx context.Context, portal *bridgev2.Portal) *event.RoomFeatures {
	if portal.RoomType == database.RoomTypeDM {
		return dmCaps
	}
	return roomCaps
}
