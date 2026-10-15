package connector

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"

	"go.mau.fi/mautrix-gvoice/pkg/libgv/gvproto"
)

var (
	_ bridgev2.IdentifierResolvingNetworkAPI = (*GVClient)(nil)
	_ bridgev2.IdentifierValidatingNetwork   = (*GVConnector)(nil)
)

func isOnlyDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func (gv *GVConnector) ValidateUserID(id networkid.UserID) bool {
	prefix, phone, ok := strings.Cut(string(id), ".")
	return ok && isOnlyDigits(prefix) && isOnlyDigits(strings.TrimPrefix(phone, "+"))
}

var ErrNotInternationalNumber = errors.New("phone number must be in international format")

func (gc *GVClient) normalizePhoneNumber(identifier string) (string, error) {
	phone, err := bridgev2.CleanNonInternationalPhoneNumber(identifier)
	if err != nil {
		return "", err
	} else if strings.HasPrefix(phone, "+") {
		return bridgev2.CleanPhoneNumber(phone)
	}
	// Numbers saved without a country code are assumed to be local when the Google Voice number is North American.
	if strings.HasPrefix(gc.UserLogin.RemoteProfile.Phone, "+1") {
		if len(phone) == 10 {
			return "+1" + phone, nil
		} else if len(phone) == 11 && phone[0] == '1' {
			return "+" + phone, nil
		}
	}
	return "", ErrNotInternationalNumber
}

func (gc *GVClient) ResolveIdentifier(ctx context.Context, identifier string, createChat bool) (*bridgev2.ResolveIdentifierResponse, error) {
	var phone string
	if gc.Main.ValidateUserID(networkid.UserID(identifier)) {
		var prefix string
		prefix, phone, _ = strings.Cut(identifier, ".")
		if prefix != gc.UserLogin.Metadata.(*UserLoginMetadata).Prefix {
			return nil, fmt.Errorf("%w: prefix mismatch", bridgev2.ErrResolveIdentifierTryNext)
		}
	} else {
		var err error
		phone, err = gc.normalizePhoneNumber(identifier)
		if err != nil {
			zerolog.Ctx(ctx).Debug().Str("input_identifier", identifier).Msg("Invalid phone number passed to ResolveIdentifier")
			return nil, bridgev2.WrapRespErrManual(err, mautrix.MInvalidParam.ErrCode, http.StatusBadRequest)
		}
	}
	userID := gc.makeUserID(phone)
	ghost, err := gc.Main.Bridge.GetGhostByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ghost: %w", err)
	}
	var userInfo *bridgev2.UserInfo
	// Don't replace the name Google Voice gave an existing ghost with the bare number.
	if ghost.Name == "" {
		userInfo = gc.wrapUserInfo(ctx, &gvproto.Contact{PhoneNumber: phone, Name: phone})
	}
	resp := &bridgev2.ResolveIdentifierResponse{
		Ghost:    ghost,
		UserID:   userID,
		UserInfo: userInfo,
	}
	if createChat {
		// Google Voice doesn't have a separate create step: 1:1 SMS threads are always
		// identified by the other party's number, and the thread is created by the first message.
		threadID := "t." + phone
		chatInfo := gc.wrapChatInfo(ctx, &gvproto.Thread{ID: threadID, IsText: true, PhoneNumbers: []string{phone}})
		chatInfo.Members.MemberMap[userID] = bridgev2.ChatMember{
			EventSender: bridgev2.EventSender{Sender: userID},
			UserInfo:    userInfo,
		}
		resp.Chat = &bridgev2.CreateChatResponse{
			PortalKey:  gc.makePortalKey(threadID),
			PortalInfo: chatInfo,
		}
	}
	return resp, nil
}
