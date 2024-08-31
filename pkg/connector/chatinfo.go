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
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/rs/zerolog"
	"go.mau.fi/util/ptr"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/event"

	"go.mau.fi/mautrix-gvoice/pkg/libgv"
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

func (gc *GVClient) wrapChatInfo(ctx context.Context, info *gvproto.Thread) *bridgev2.ChatInfo {
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
			UserInfo:    gc.wrapUserInfo(ctx, member),
		}
	}
	wrapped.ExtraUpdates = func(ctx context.Context, portal *bridgev2.Portal) bool {
		slices.Sort(info.PhoneNumbers)
		meta := portal.Metadata.(*PortalMetadata)
		if !slices.Equal(info.PhoneNumbers, meta.Participants) {
			meta.Participants = info.PhoneNumbers
			return true
		}
		return false
	}
	return &wrapped
}

const MaxAvatarSize = 5 * 1024 * 1024

func (gc *GVClient) wrapAvatar(url string) *bridgev2.Avatar {
	return &bridgev2.Avatar{
		ID: networkid.AvatarID(url),
		Get: func(ctx context.Context) ([]byte, error) {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				return nil, err
			}
			req.Header.Set("User-Agent", libgv.UserAgent)
			req.Header.Set("Referer", libgv.Origin+"/")
			req.Header.Set("Accept", "image/webp,image/png,image/*;q=0.8,*/*;q=0.5")
			resp, err := gc.Client.HTTP.Do(req)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()
			if resp.ContentLength > MaxAvatarSize {
				return nil, fmt.Errorf("image too large")
			}
			return io.ReadAll(io.LimitReader(resp.Body, MaxAvatarSize))
		},
		Remove: url == "",
	}
}

func (gc *GVClient) wrapUserInfo(ctx context.Context, info *gvproto.Contact) *bridgev2.UserInfo {
	contactInfo := gc.fetchContactInfo(ctx, info.PhoneNumber)
	if contactInfo == nil {
		contactInfo = &ProcessedContact{}
	}
	identifiers := make([]string, len(contactInfo.Phones)+len(contactInfo.Emails))
	for i, phone := range contactInfo.Phones {
		identifiers[i] = fmt.Sprintf("tel:%s", phone)
	}
	for i, email := range contactInfo.Emails {
		identifiers[i+len(contactInfo.Phones)] = fmt.Sprintf("mailto:%s", email)
	}
	mainPhoneURI := fmt.Sprintf("tel:%s", info.PhoneNumber)
	if !slices.Contains(identifiers, mainPhoneURI) {
		identifiers = append(identifiers, mainPhoneURI)
	}
	return &bridgev2.UserInfo{
		Identifiers: identifiers,
		Name:        ptr.Ptr(gc.Main.Config.FormatDisplayname(info, contactInfo)),
		Avatar:      gc.wrapAvatar(contactInfo.AvatarURL),
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

type ProcessedContact struct {
	AvatarURL string
	FirstName string
	Name      string
	Emails    []string
	Phones    []string
}

func processContact(person *gvproto.Person) *ProcessedContact {
	if person == nil {
		return nil
	}
	var out ProcessedContact
	for _, method := range person.GetContactMethods() {
		displayInfo := method.GetDisplayInfo()
		if displayInfo.GetPhoto().GetType() != gvproto.ContactDisplayInfo_Photo_MONOGRAM && (out.AvatarURL == "" || displayInfo.GetPrimary()) {
			out.AvatarURL = displayInfo.GetPhoto().GetURL()
		}
		if displayInfo.GetName().GetValue() != "" && (out.Name == "" || displayInfo.GetPrimary()) {
			out.Name = displayInfo.GetName().GetValue()
		}
		if displayInfo.GetName().GetGivenName() != "" && (out.FirstName == "" || displayInfo.GetPrimary()) {
			out.FirstName = displayInfo.GetName().GetGivenName()
		}
		if method.GetPhone().GetCanonicalValue() != "" {
			out.Phones = append(out.Phones, method.GetPhone().GetCanonicalValue())
		}
		if method.GetEmail().GetEmail() != "" {
			out.Emails = append(out.Emails, method.GetEmail().GetEmail())
		}
	}
	slices.Sort(out.Phones)
	slices.Sort(out.Emails)
	out.Phones = slices.Compact(out.Phones)
	out.Emails = slices.Compact(out.Emails)
	return &out
}

func (gc *GVClient) fetchContactInfo(ctx context.Context, e164 string) *ProcessedContact {
	gc.contactCacheLock.Lock()
	defer gc.contactCacheLock.Unlock()
	if person, ok := gc.contactCache[e164]; ok {
		return person
	}
	zerolog.Ctx(ctx).Debug().Str("e164", e164).Msg("Fetching contact info")
	resp, err := gc.Client.LookupContact(ctx, e164)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Str("e164", e164).Msg("Failed to fetch contact info")
		return nil
	}
	match, ok := resp[e164]
	if !ok {
		zerolog.Ctx(ctx).Warn().Str("e164", e164).Msg("Contact not found in response")
	} else if match.GetFailureType() != gvproto.RespLookupContacts_Match_NO_FAILURE {
		if match.GetFailureType() == gvproto.RespLookupContacts_Match_PERMANENT {
			gc.contactCache[e164] = nil
		}
		zerolog.Ctx(ctx).Debug().
			Str("e164", e164).
			Stringer("failure", match.GetFailureType()).
			Msg("Contact not found")
	} else {
		person := processContact(match.GetAutocompletion().GetPerson())
		gc.contactCache[e164] = person
		return person
	}
	return nil
}
