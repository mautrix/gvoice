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
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/event"

	"go.mau.fi/mautrix-gvoice/pkg/libgv"
	"go.mau.fi/mautrix-gvoice/pkg/libgv/gvproto"
)

var (
	_ bridgev2.ReadReceiptHandlingNetworkAPI = (*GVClient)(nil)
	_ bridgev2.DeleteChatHandlingNetworkAPI  = (*GVClient)(nil)
	_ bridgev2.UserBlockingNetworkAPI        = (*GVClient)(nil)
)

func (gc *GVClient) HandleMatrixMessage(ctx context.Context, msg *bridgev2.MatrixMessage) (message *bridgev2.MatrixMessageResponse, err error) {
	req := &gvproto.ReqSendSMS{
		ThreadID: string(msg.Portal.ID),
		TransactionID: &gvproto.ReqSendSMS_WrappedTxnID{
			ID: libgv.GenerateTransactionID(),
		},
	}
	recipients := msg.Portal.Metadata.(*PortalMetadata).Participants
	if newThread, err := gc.isNewThread(ctx, msg.Portal); err != nil {
		return nil, fmt.Errorf("failed to check if thread exists: %w", err)
	} else if newThread {
		// Like the web app, address the first message of a new thread to the recipients, which creates the thread.
		req.ThreadID = ""
		req.Recipients = recipients
	}
	var electronStatus string
	if rs := gc.requestSignature.Load(); rs != nil {
		td, err := (*rs)(ctx, req.ThreadID, recipients, req.TransactionID.ID)
		if err != nil {
			zerolog.Ctx(ctx).Err(err).Msg("Failed to generate signature")
			electronStatus = "failed"
		} else if td != "" {
			req.TrackingData = &gvproto.ReqSendSMS_TrackingData{Data: td}
			electronStatus = "ok"
		} else {
			electronStatus = "failed/empty"
		}
	} else {
		zerolog.Ctx(ctx).Debug().Msg("Electron not running, sending message without signature data")
		electronStatus = "unavailable"
	}
	switch msg.Content.MsgType {
	case event.MsgText, event.MsgNotice:
		req.Text = msg.Content.Body
	case event.MsgEmote:
		req.Text = "/me " + msg.Content.Body
	case event.MsgImage:
		var mediaType gvproto.ReqSendSMS_Media_Type
		switch msg.Content.Info.MimeType {
		case "image/png":
			mediaType = gvproto.ReqSendSMS_Media_PNG
		case "image/jpeg":
			mediaType = gvproto.ReqSendSMS_Media_JPEG
		case "image/bmp":
			mediaType = gvproto.ReqSendSMS_Media_BMP
		case "image/tiff":
			mediaType = gvproto.ReqSendSMS_Media_TIFF
		case "image/gif":
			mediaType = gvproto.ReqSendSMS_Media_GIF
		default:
			return nil, fmt.Errorf("%w %s", bridgev2.ErrUnsupportedMediaType, msg.Content.Info.MimeType)
		}
		data, err := gc.Main.Bridge.Bot.DownloadMedia(ctx, msg.Content.URL, msg.Content.File)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", bridgev2.ErrMediaDownloadFailed, err)
		}
		caption := ""
		if msg.Content.FileName != "" && msg.Content.FileName != msg.Content.Body {
			caption = msg.Content.Body
		}
		req.Media = &gvproto.ReqSendSMS_Media{
			Type: mediaType,
			Data: data,
		}
		req.Text = caption
	default:
		return nil, fmt.Errorf("%w %s", bridgev2.ErrUnsupportedMessageType, msg.Content.MsgType)
	}
	resp, err := gc.Client.SendMessage(ctx, req)
	if err != nil {
		if isRateLimitError(err) {
			err = fmt.Errorf("%w (electron status: %s)", err, electronStatus)
		}
		return nil, wrapSendError(err)
	}
	if strings.HasPrefix(string(msg.Portal.ID), newGroupPrefix) && resp.ThreadID != "" {
		result, _, err := gc.Main.Bridge.ReIDPortal(ctx, msg.Portal.PortalKey, gc.makePortalKey(resp.ThreadID))
		if err != nil {
			zerolog.Ctx(ctx).Err(err).Str("thread_id", resp.ThreadID).Msg("Failed to re-ID new group portal")
		} else {
			zerolog.Ctx(ctx).Debug().Str("thread_id", resp.ThreadID).Stringer("result", result).Msg("Re-IDed new group portal")
		}
	}
	return &bridgev2.MatrixMessageResponse{
		DB: &database.Message{
			ID:        networkid.MessageID(resp.ThreadItemID),
			Timestamp: time.UnixMilli(resp.TimestampMS),
		},
	}, nil
}

// isNewThread returns true for chats started from Matrix that don't have a Google Voice thread yet.
func (gc *GVClient) isNewThread(ctx context.Context, portal *bridgev2.Portal) (bool, error) {
	if strings.HasPrefix(string(portal.ID), newGroupPrefix) {
		return true, nil
	}
	if first, err := gc.Main.Bridge.DB.Message.GetFirstPortalMessage(ctx, portal.PortalKey); err != nil || first != nil {
		return false, err
	}
	_, err := gc.Client.GetThread(ctx, string(portal.ID), 1, "")
	if libgv.IsNotFoundError(err) {
		return true, nil
	}
	return false, wrapSendError(err)
}

func (gc *GVClient) getMergedThreads(ctx context.Context, portal *bridgev2.Portal, start, end time.Time) (threadIDs []string, hasOwnMessages bool, err error) {
	messages, err := gc.Main.Bridge.DB.Message.GetMessagesBetweenTimeQuery(ctx, portal.PortalKey, start, end)
	if err != nil {
		return nil, false, err
	}
	for _, msg := range messages {
		threadID := msg.Metadata.(*MessageMetadata).ThreadID
		if threadID == "" || threadID == string(portal.ID) {
			hasOwnMessages = true
		} else if !slices.Contains(threadIDs, threadID) {
			threadIDs = append(threadIDs, threadID)
		}
	}
	return threadIDs, hasOwnMessages, nil
}

func (gc *GVClient) HandleMatrixReadReceipt(ctx context.Context, msg *bridgev2.MatrixReadReceipt) error {
	start := msg.ReadUpTo.Add(-7 * 24 * time.Hour)
	if msg.LastRead.After(start) {
		start = msg.LastRead
	}
	threadIDs, hasOwnMessages, err := gc.getMergedThreads(ctx, msg.Portal, start, msg.ReadUpTo)
	if err != nil {
		return fmt.Errorf("failed to get merged threads: %w", err)
	}
	if hasOwnMessages {
		threadIDs = append(threadIDs, string(msg.Portal.ID))
	}
	var errs []error
	for _, threadID := range threadIDs {
		resp, err := gc.Client.UpdateThreadAttributes(ctx, &gvproto.ReqUpdateAttributes{
			Attributes: &gvproto.ThreadAttributes{
				ThreadID: threadID,
				Read:     true,
			},
			OtherAttributes: &gvproto.ThreadAttributes{
				Read: true,
			},
			UnknownInt: 1,
		})
		zerolog.Ctx(ctx).Trace().Str("thread_id", threadID).Any("resp", resp).Msg("Update attributes response")
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to mark %s as read: %w", threadID, err))
		}
	}
	return errors.Join(errs...)
}

func (gc *GVClient) HandleMatrixDeleteChat(ctx context.Context, chat *bridgev2.MatrixDeleteChat) error {
	threadIDs, hasOwnMessages, err := gc.getMergedThreads(ctx, chat.Portal, time.Unix(0, 0), time.Now().Add(24*time.Hour))
	if err != nil {
		return fmt.Errorf("failed to get merged threads: %w", err)
	}
	for _, threadID := range threadIDs {
		_, err = gc.Client.DeleteThread(ctx, threadID)
		if err != nil {
			zerolog.Ctx(ctx).Err(err).Str("thread_id", threadID).Msg("Failed to delete merged thread")
		}
	}
	if hasOwnMessages || len(threadIDs) == 0 {
		_, err = gc.Client.DeleteThread(ctx, string(chat.Portal.ID))
		return wrapSendError(err)
	}
	return nil
}

func (gc *GVClient) HandleMatrixBlockUser(ctx context.Context, msg *bridgev2.MatrixBlockUser) error {
	if msg.Content.ReportSpam && !msg.Content.Block {
		return bridgev2.WrapErrorInStatus(errors.New("spam reporting requires blocking the user")).
			WithMessage("Google Voice can only report spam when blocking the user").
			WithErrorReason(event.MessageStatusUnsupported).
			WithStatus(event.MessageStatusFail).
			WithIsCertain(true).
			WithSendNotice(true)
	}
	err := gc.setThreadBlocked(ctx, string(msg.Portal.ID), msg.Content.Block, msg.Content.ReportSpam)
	if !libgv.IsNotFoundError(err) {
		return wrapSendError(err)
	}
	// Blocks apply to the number, so a merged call thread works when there's no SMS thread.
	threadIDs, _, dbErr := gc.getMergedThreads(ctx, msg.Portal, time.Unix(0, 0), time.Now().Add(24*time.Hour))
	if dbErr != nil {
		return fmt.Errorf("failed to get merged threads: %w", dbErr)
	} else if len(threadIDs) == 0 {
		return bridgev2.WrapErrorInStatus(fmt.Errorf("conversation doesn't exist on Google Voice yet: %w", err)).
			WithMessage("this conversation doesn't exist on Google Voice yet").
			WithErrorReason(event.MessageStatusUnsupported).
			WithStatus(event.MessageStatusFail).
			WithIsCertain(true).
			WithSendNotice(true)
	}
	return wrapSendError(gc.setThreadBlocked(ctx, threadIDs[0], msg.Content.Block, msg.Content.ReportSpam))
}

func (gc *GVClient) setThreadBlocked(ctx context.Context, threadID string, blocked, spam bool) error {
	_, err := gc.Client.UpdateThreadAttributes(ctx, &gvproto.ReqUpdateAttributes{
		Attributes: &gvproto.ThreadAttributes{
			ThreadID:  threadID,
			IsBlocked: blocked,
			IsSpammer: spam,
		},
		OtherAttributes: &gvproto.ThreadAttributes{
			IsBlocked: true,
			IsSpammer: spam,
		},
		UnknownInt: 1,
	})
	return err
}
