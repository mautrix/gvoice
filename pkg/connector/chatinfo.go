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
	"fmt"
	"strings"

	"go.mau.fi/util/ptr"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/event"

	"go.mau.fi/mautrix-gvoice/pkg/libgv/gvproto"
)

func (gc *GVClient) IsThisUser(ctx context.Context, userID networkid.UserID) bool {
	return false
}

func (gc *GVClient) GetChatInfo(ctx context.Context, portal *bridgev2.Portal) (*bridgev2.ChatInfo, error) {
	return nil, fmt.Errorf("not implemented")
}

func (gc *GVClient) GetUserInfo(ctx context.Context, ghost *bridgev2.Ghost) (*bridgev2.UserInfo, error) {
	return nil, nil
}

func (gc *GVClient) makePortalKey(threadID string) networkid.PortalKey {
	return networkid.PortalKey{
		ID:       networkid.PortalID(threadID),
		Receiver: gc.UserLogin.ID,
	}
}

func (gc *GVClient) makeUserID(e164 string) networkid.UserID {
	return networkid.UserID(fmt.Sprintf("%s.%s", gc.UserLogin.Metadata.(*UserLoginMetadata).Prefix, e164))
}

func (gc *GVClient) wrapChatInfo(info *gvproto.Thread) *bridgev2.ChatInfo {
	var wrapped bridgev2.ChatInfo
	wrapped.Members = &bridgev2.ChatMemberList{
		IsFull:           true,
		TotalMemberCount: len(info.PhoneNumbers) + 1,
		MemberMap:        make(map[networkid.UserID]bridgev2.ChatMember, len(info.PhoneNumbers)+1),
		PowerLevels: &bridgev2.PowerLevelOverrides{
			Events: map[event.Type]int{
				event.EventReaction: 99,
			},
		},
	}
	wrapped.Members.MemberMap[""] = bridgev2.ChatMember{
		EventSender: bridgev2.EventSender{IsFromMe: true},
	}
	wrapped.CanBackfill = true
	if len(info.PhoneNumbers) == 1 {
		wrapped.Members.OtherUserID = gc.makeUserID(info.PhoneNumbers[0])
		wrapped.Type = ptr.Ptr(database.RoomTypeDM)
	} else {
		wrapped.Type = ptr.Ptr(database.RoomTypeDefault)
	}
	for _, member := range info.Contacts {
		if strings.HasPrefix(member.PhoneNumber, "Group Message.") {
			continue
		}
		userID := gc.makeUserID(member.PhoneNumber)
		wrapped.Members.MemberMap[userID] = bridgev2.ChatMember{
			EventSender: bridgev2.EventSender{Sender: userID},
			UserInfo:    gc.wrapUserInfo(member),
		}
	}
	return &wrapped
}

func (gc *GVClient) wrapUserInfo(info *gvproto.Contact) *bridgev2.UserInfo {
	return &bridgev2.UserInfo{
		Identifiers: []string{fmt.Sprintf("tel:%s", info.PhoneNumber)},
		Name:        ptr.Ptr(gc.Main.Config.FormatDisplayname(info.PhoneNumber, info.Name)),
		IsBot:       ptr.Ptr(false),
		ExtraUpdates: func(ctx context.Context, ghost *bridgev2.Ghost) bool {
			meta := ghost.Metadata.(*GhostMetadata)
			if meta.Phone != info.PhoneNumber {
				meta.Phone = info.PhoneNumber
				return true
			}
			return false
		},
	}
}
