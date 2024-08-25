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
)

func (gc *GVClient) HandleMatrixMessage(ctx context.Context, msg *bridgev2.MatrixMessage) (message *bridgev2.MatrixMessageResponse, err error) {
	req := &gvproto.ReqSendSMS{
		ThreadID: string(msg.Portal.ID),
		TransactionID: &gvproto.ReqSendSMS_WrappedTxnID{
			ID: libgv.GenerateTransactionID(),
		},
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
		default:
			return nil, fmt.Errorf("%w %s", bridgev2.ErrUnsupportedMediaType, msg.Content.Info.MimeType)
		}
		data, err := gc.Main.Bridge.Bot.DownloadMedia(ctx, msg.Content.URL, msg.Content.File)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", bridgev2.ErrMediaDownloadFailed, err)
		}
		fileName := msg.Content.Body
		caption := ""
		if msg.Content.FileName != "" && msg.Content.FileName != msg.Content.Body {
			fileName = msg.Content.FileName
			caption = msg.Content.Body
		}
		mediaURL, err := gc.Client.UploadPhoto(ctx, fileName, msg.Content.Info.MimeType, data)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", bridgev2.ErrMediaReuploadFailed, err)
		}
		req.Media = &gvproto.ReqSendSMS_Media{
			Type: mediaType,
			URL:  mediaURL,
		}
		req.Text = caption
	default:
		return nil, fmt.Errorf("%w %s", bridgev2.ErrUnsupportedMessageType, msg.Content.MsgType)
	}
	resp, err := gc.Client.SendMessage(ctx, req)
	if err != nil {
		return nil, err
	}
	return &bridgev2.MatrixMessageResponse{
		DB: &database.Message{
			ID:        networkid.MessageID(resp.ThreadItemID),
			Timestamp: time.UnixMilli(resp.TimestampMS),
		},
	}, nil
}

func (gc *GVClient) HandleMatrixReadReceipt(ctx context.Context, msg *bridgev2.MatrixReadReceipt) error {
	resp, err := gc.Client.UpdateThreadAttributes(ctx, &gvproto.ReqUpdateAttributes{
		Attributes: &gvproto.ThreadAttributes{
			ThreadID: string(msg.Portal.ID),
			Read:     true,
		},
		OtherAttributes: &gvproto.ThreadAttributes{
			Read: true,
		},
		UnknownInt: 1,
	})
	zerolog.Ctx(ctx).Trace().Any("resp", resp).Msg("Update attributes response")
	return err
}
