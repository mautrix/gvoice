package connector

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/rs/zerolog"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"

	"go.mau.fi/mautrix-gvoice/pkg/libgv"
	"go.mau.fi/mautrix-gvoice/pkg/libgv/gvproto"
)

var (
	_ bridgev2.IdentifierResolvingNetworkAPI = (*GVClient)(nil)
	_ bridgev2.GroupCreatingNetworkAPI       = (*GVClient)(nil)
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
	// Google Voice uses the number as the contact name when the number isn't saved.
	contact := &gvproto.Contact{PhoneNumber: phone, Name: phone}
	if !createChat {
		return &bridgev2.ResolveIdentifierResponse{
			Ghost:    ghost,
			UserID:   userID,
			UserInfo: gc.wrapUserInfo(ctx, contact),
		}, nil
	}
	threadID := "t." + phone
	var thread *gvproto.Thread
	resp, err := gc.Client.GetThread(ctx, threadID, 1, "")
	newThread := libgv.IsNotFoundError(err)
	if newThread {
		// There's no create call: the first message creates the thread (see isNewThread).
		thread = &gvproto.Thread{
			ID:           threadID,
			IsText:       true,
			PhoneNumbers: []string{phone},
			Contacts:     []*gvproto.Contact{contact},
		}
	} else if err != nil {
		return nil, wrapRespError(fmt.Errorf("failed to get thread: %w", err))
	} else {
		thread = resp.Thread
	}
	chatInfo := gc.wrapChatInfo(ctx, thread)
	chatInfo.CanBackfill = !newThread
	return &bridgev2.ResolveIdentifierResponse{
		Ghost:    ghost,
		UserID:   userID,
		UserInfo: chatInfo.Members.MemberMap[userID].UserInfo,
		Chat: &bridgev2.CreateChatResponse{
			PortalKey:  gc.makePortalKey(thread.ID),
			PortalInfo: chatInfo,
		},
	}, nil
}

const newGroupPrefix = "new."

func (gc *GVClient) CreateGroup(ctx context.Context, params *bridgev2.GroupCreateParams) (*bridgev2.CreateChatResponse, error) {
	phones := make([]string, len(params.Participants))
	for i, userID := range params.Participants {
		prefix, phone, _ := strings.Cut(string(userID), ".")
		if prefix != gc.UserLogin.Metadata.(*UserLoginMetadata).Prefix {
			return nil, bridgev2.RespError(mautrix.MInvalidParam.WithMessage("user ID %s is not on this account", userID))
		}
		phones[i] = phone
	}
	slices.Sort(phones)
	phones = slices.Compact(phones)
	if len(phones) < 2 {
		return nil, bridgev2.RespError(mautrix.MInvalidParam.WithMessage("group texts need at least two other people"))
	}
	contacts := make([]*gvproto.Contact, len(phones))
	for i, phone := range phones {
		contacts[i] = &gvproto.Contact{PhoneNumber: phone, Name: phone}
	}
	thread := &gvproto.Thread{
		ID:           newGroupPrefix + strings.Join(phones, ","),
		IsText:       true,
		PhoneNumbers: phones,
		Contacts:     contacts,
	}
	portalKey := gc.makePortalKey(thread.ID)
	chatInfo := gc.wrapChatInfo(ctx, thread)
	chatInfo.CanBackfill = false
	var portal *bridgev2.Portal
	if params.RoomID != "" {
		var err error
		portal, err = gc.Main.Bridge.GetPortalByKey(ctx, portalKey)
		if err != nil {
			return nil, fmt.Errorf("failed to get portal by key: %w", err)
		}
		err = portal.UpdateMatrixRoomID(ctx, params.RoomID, bridgev2.UpdateMatrixRoomIDParams{
			OverwriteOldPortal: true,
			TombstoneOldRoom:   true,
			DeleteOldRoom:      true,
			ChatInfo:           chatInfo,
			ChatInfoSource:     gc.UserLogin,
		})
		if err != nil {
			return nil, err
		}
	}
	return &bridgev2.CreateChatResponse{
		PortalKey:  portalKey,
		PortalInfo: chatInfo,
		Portal:     portal,
	}, nil
}
