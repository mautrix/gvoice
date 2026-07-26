// mautrix-gvoice - A Matrix-Google Voice puppeting bridge.
//
// Starting new chats. Upstream doesn't implement this, so conversations could only be replied
// to, never begun — you had to open the Google Voice app first (Beeper documents the same
// limitation). Google Voice SMS thread ids are deterministic: "t." followed by the peer's E.164
// number, e.g. "t.+15551234567" — confirmed from live portal ids. So resolving an identifier is
// pure string work: no round trip, and the first send goes through the ordinary code path with a
// thread id the server already agrees with.

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
)

var _ bridgev2.IdentifierResolvingNetworkAPI = (*GVClient)(nil)

// threadIDForNumber builds the SMS thread id Google Voice uses for a 1:1 conversation.
func threadIDForNumber(e164 string) string {
	return "t." + e164
}

// normalizeGVNumber turns user input into E.164. Google Voice is US-only, so a bare 10-digit
// number is assumed to be +1 and an 11-digit number starting with 1 just gains the plus. Anything
// already written with a leading + is kept as-is so international peers still work.
func normalizeGVNumber(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	hadPlus := strings.HasPrefix(trimmed, "+")
	var digits strings.Builder
	for _, r := range trimmed {
		if r >= '0' && r <= '9' {
			digits.WriteRune(r)
		}
	}
	d := digits.String()
	switch {
	case d == "":
		return "", fmt.Errorf("%q doesn't contain a phone number", input)
	case hadPlus:
		return "+" + d, nil
	case len(d) == 10:
		return "+1" + d, nil
	case len(d) == 11 && d[0] == '1':
		return "+" + d, nil
	case len(d) < 7:
		return "", fmt.Errorf("%q is too short to be a phone number", input)
	default:
		// Longer bare number without a plus: assume it's already country-coded.
		return "+" + d, nil
	}
}

// newDMChatInfo describes a 1:1 portal we're creating locally. GetChatInfo is "not implemented"
// on this bridge, so the central module cannot fetch this for us — if we don't supply it here,
// portal creation fails. Participants must be populated too: HandleMatrixMessage reads it to
// build the WAA signature's recipient list.
func (gc *GVClient) newDMChatInfo(e164 string, userID networkid.UserID) *bridgev2.ChatInfo {
	info := &bridgev2.ChatInfo{
		Type: ptr.Ptr(database.RoomTypeDM),
		Members: &bridgev2.ChatMemberList{
			IsFull:           true,
			TotalMemberCount: 2,
			OtherUserID:      userID,
			MemberMap: map[networkid.UserID]bridgev2.ChatMember{
				"":     {EventSender: bridgev2.EventSender{IsFromMe: true}},
				userID: {EventSender: bridgev2.EventSender{Sender: userID}},
			},
			PowerLevels: &bridgev2.PowerLevelOverrides{
				Events: map[event.Type]int{
					event.EventReaction: 99,
				},
			},
		},
	}
	info.ExtraUpdates = func(ctx context.Context, portal *bridgev2.Portal) bool {
		meta := portal.Metadata.(*PortalMetadata)
		if len(meta.Participants) == 0 {
			meta.Participants = []string{e164}
			return true
		}
		return false
	}
	return info
}

// ResolveIdentifier maps a phone number to its ghost and 1:1 portal, creating the portal when
// asked. Nothing is sent here — the thread materialises server-side on the first message, which
// is why the id has to be derivable rather than returned by an API call.
func (gc *GVClient) ResolveIdentifier(ctx context.Context, identifier string, createChat bool) (*bridgev2.ResolveIdentifierResponse, error) {
	e164, err := normalizeGVNumber(identifier)
	if err != nil {
		return nil, err
	}
	userID := gc.makeUserID(e164)
	resp := &bridgev2.ResolveIdentifierResponse{
		UserID: userID,
	}
	if createChat {
		resp.Chat = &bridgev2.CreateChatResponse{
			PortalKey:  gc.makePortalKey(threadIDForNumber(e164)),
			PortalInfo: gc.newDMChatInfo(e164, userID),
		}
	}
	return resp, nil
}
