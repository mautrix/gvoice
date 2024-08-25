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
	"bytes"
	"context"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"go.mau.fi/util/exmime"
	"maunium.net/go/mautrix/bridge/status"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/bridgev2/simplevent"
	"maunium.net/go/mautrix/event"

	"go.mau.fi/mautrix-gvoice/pkg/libgv"
	"go.mau.fi/mautrix-gvoice/pkg/libgv/gvproto"
)

const (
	BackgroundRefreshInterval = 15 * time.Minute
	MinRefreshInterval        = 10 * time.Second
	MinRefreshBurstCount      = 5
	RefreshDelay              = 500 * time.Millisecond
)

var _ bridgev2.BackfillingNetworkAPI = (*GVClient)(nil)

func isNewMessages(evt *libgv.RealtimeEvent) bool {
	w := evt.GetDataWrapper()
	if len(w) == 0 || len(w[0].GetData()) == 0 {
		return false
	}
	d := w[0].GetData()[0].GetEvent().GetSub2().GetData()
	if len(d) == 0 {
		return false
	}
	return bytes.HasPrefix(d[0].GetUnknownBytes(), []byte("["))
}

func (gc *GVClient) handleRealtimeEvent(ctx context.Context, rawEvt any) {
	switch evt := rawEvt.(type) {
	case *libgv.CookieChanged:
		gc.UserLogin.Metadata.(*UserLoginMetadata).Cookies = evt.Cookies
		err := gc.UserLogin.Save(ctx)
		if err != nil {
			zerolog.Ctx(ctx).Err(err).Msg("Failed to save cookies")
		}
	case *libgv.RealtimeConnected:
		gc.UserLogin.BridgeState.Send(status.BridgeState{StateEvent: status.StateConnected})
		go func() {
			gc.wakeupMessageFetcher <- struct{}{}
			zerolog.Ctx(ctx).Debug().Msg("Woke up message fetcher from connected event")
		}()
	case *libgv.RealtimeEvent:
		if isNewMessages(evt) {
			select {
			case gc.wakeupMessageFetcher <- struct{}{}:
				zerolog.Ctx(ctx).Debug().Msg("Woke up message fetcher from realtime event")
			default:
			}
		}
	}
}

func (gc *GVClient) fetchNewMessagesLoop(ctx context.Context) {
	ctxDone := ctx.Done()
	ticker := time.NewTicker(BackgroundRefreshInterval)
	defer ticker.Stop()
	for {
		select {
		case <-gc.wakeupMessageFetcher:
		case <-ticker.C:
		case <-ctxDone:
			zerolog.Ctx(ctx).Debug().Msg("Stopping message fetcher loop")
			return
		}
		time.Sleep(RefreshDelay)
		err := gc.fetchEventsLimiter.Wait(ctx)
		if err != nil {
			zerolog.Ctx(ctx).Debug().Err(err).Msg("Failed to wait for ratelimiter")
			return
		}
		go gc.fetchNewMessages(ctx)
	}
}

func (gc *GVClient) fetchNewMessages(ctx context.Context) {
	gc.fetchEventsLock.Lock()
	defer gc.fetchEventsLock.Unlock()
	zerolog.Ctx(ctx).Debug().Msg("Fetching new messages")
	resp, err := gc.Client.ListThreads(ctx, "")
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to list threads")
		return
	}
	for _, thread := range resp.Threads {
		if len(thread.Messages) == 0 {
			continue
		}
		lastMessageTS := time.UnixMilli(thread.Messages[0].Timestamp)
		portalKey := gc.makePortalKey(thread.ID)
		prevMsg, ok := gc.lastEvents[thread.ID]
		gc.lastEvents[thread.ID] = lastMessageTS
		if !ok {
			gc.Main.Bridge.QueueRemoteEvent(gc.UserLogin, &simplevent.ChatResync{
				EventMeta: simplevent.EventMeta{
					Type:         bridgev2.RemoteEventChatResync,
					PortalKey:    portalKey,
					CreatePortal: true,
				},
				ChatInfo:            gc.wrapChatInfo(thread),
				LatestMessageTS:     lastMessageTS,
				BundledBackfillData: thread,
			})
			continue
		}
		for _, msg := range thread.Messages {
			ts, txnID, sender := getMessageMeta(msg)
			if !ts.After(prevMsg) {
				break
			}
			gc.Main.Bridge.QueueRemoteEvent(gc.UserLogin, &simplevent.Message[*gvproto.Message]{
				EventMeta: simplevent.EventMeta{
					Type: bridgev2.RemoteEventMessage,
					LogContext: func(c zerolog.Context) zerolog.Context {
						return c.
							Int64("timestamp", msg.Timestamp).
							Int64("txn_id", msg.TransactionID).
							Str("top_level_sender", msg.GetContact().GetPhoneNumber()).
							Str("mms_sender", msg.GetMMS().GetSenderPhoneNumber())
					},
					PortalKey:   portalKey,
					Sender:      sender,
					Timestamp:   ts,
					StreamOrder: msg.Timestamp,
				},
				ConvertMessageFunc: gc.convertMessage,
				Data:               msg,
				ID:                 networkid.MessageID(msg.ID),
				TransactionID:      txnID,
			})
		}
	}
}

func (gc *GVClient) FetchMessages(ctx context.Context, params bridgev2.FetchMessagesParams) (*bridgev2.FetchMessagesResponse, error) {
	if params.Count <= 0 {
		return nil, fmt.Errorf("count must be positive")
	}
	var thread *gvproto.Thread
	var messagesToConvert []*gvproto.Message
	if params.BundledData != nil {
		thread = params.BundledData.(*gvproto.Thread)
		messagesToConvert = thread.Messages
	}
	if params.Forward {
		if thread == nil {
			resp, err := gc.Client.GetThread(ctx, string(params.Portal.ID), 100, "")
			if err != nil {
				return nil, fmt.Errorf("failed to fetch latest messages: %w", err)
			}
			thread = resp.Thread
			messagesToConvert = thread.Messages
		}
		didCutOff := false
		if params.AnchorMessage != nil {
			for i, msg := range messagesToConvert {
				if networkid.MessageID(msg.ID) == params.AnchorMessage.ID || !time.UnixMilli(msg.Timestamp).After(params.AnchorMessage.Timestamp) {
					messagesToConvert = messagesToConvert[:i]
					didCutOff = true
				}
			}
		}
		for len(messagesToConvert) < params.Count && !didCutOff && thread.PaginationToken != "" {
			resp, err := gc.Client.GetThread(ctx, string(params.Portal.ID), 100, thread.PaginationToken)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch messages: %w", err)
			}
			thread = resp.Thread
			messagesToConvert = append(messagesToConvert, thread.Messages...)
			if params.AnchorMessage != nil {
				for i, msg := range messagesToConvert {
					if networkid.MessageID(msg.ID) == params.AnchorMessage.ID || !time.UnixMilli(msg.Timestamp).After(params.AnchorMessage.Timestamp) {
						messagesToConvert = messagesToConvert[:i]
						didCutOff = true
					}
				}
			}
		}
	} else {
		paginationToken := string(params.Cursor)
		if paginationToken == "" {
			if params.AnchorMessage == nil {
				return nil, fmt.Errorf("can't backward backfill without either cursor or anchor message")
			}
			paginationToken = strconv.FormatInt(params.AnchorMessage.Timestamp.UnixMilli(), 10)
		}
		for len(messagesToConvert) < params.Count && paginationToken != "" {
			resp, err := gc.Client.GetThread(ctx, string(params.Portal.ID), 100, paginationToken)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch messages: %w", err)
			}
			thread = resp.Thread
			messagesToConvert = append(messagesToConvert, thread.Messages...)
			paginationToken = thread.PaginationToken
		}
		if thread == nil {
			return nil, fmt.Errorf("unexpected state: no thread")
		}
	}
	convertedMessages := make([]*bridgev2.BackfillMessage, len(messagesToConvert))
	for i, msg := range messagesToConvert {
		ts, txnID, sender := getMessageMeta(msg)
		converted, _ := gc.convertMessage(ctx, params.Portal, gc.Main.Bridge.Bot, msg)
		convertedMessages[i] = &bridgev2.BackfillMessage{
			ConvertedMessage: converted,
			ID:               networkid.MessageID(msg.ID),
			Sender:           sender,
			TxnID:            txnID,
			Timestamp:        ts,
			StreamOrder:      msg.Timestamp,
		}
	}
	slices.Reverse(convertedMessages)
	return &bridgev2.FetchMessagesResponse{
		Messages: convertedMessages,
		Cursor:   networkid.PaginationCursor(thread.PaginationToken),
		HasMore:  thread.PaginationToken != "",
		MarkRead: thread.Read,
	}, nil
}

func getMessageMeta(msg *gvproto.Message) (ts time.Time, txnID networkid.TransactionID, sender bridgev2.EventSender) {
	ts = time.UnixMilli(msg.Timestamp)
	if msg.Type == gvproto.Message_SMS_OUT {
		sender.IsFromMe = true
	} else if senderNum := msg.GetMMS().GetSenderPhoneNumber(); senderNum != "" {
		sender.Sender = networkid.UserID(senderNum)
	} else {
		sender.Sender = networkid.UserID(msg.GetContact().GetPhoneNumber())
	}
	if msg.TransactionID != 0 {
		txnID = networkid.TransactionID(strconv.FormatInt(msg.TransactionID, 10))
	}
	return
}

func (gc *GVClient) convertMessage(ctx context.Context, portal *bridgev2.Portal, intent bridgev2.MatrixAPI, msg *gvproto.Message) (*bridgev2.ConvertedMessage, error) {
	var content event.MessageEventContent
	output := &bridgev2.ConvertedMessage{
		Parts: []*bridgev2.ConvertedMessagePart{{
			Type:    event.EventMessage,
			Content: &content,
		}},
	}
	content.MsgType = event.MsgText
	if msg.MMS != nil {
		if msg.MMS.Subject != "" {
			content.Body = fmt.Sprintf("**%s**\n%s", msg.MMS.Subject, msg.MMS.Text)
			content.FormattedBody = fmt.Sprintf("<strong>%s</strong><br>%s", msg.MMS.Subject, msg.MMS.Text)
			content.Format = event.FormatHTML
		} else {
			content.Body = msg.MMS.Text
		}
		for i, att := range msg.MMS.Attachments {
			if i == 0 {
				gc.convertMedia(ctx, portal, intent, att, &content)
			} else {
				var anotherPart event.MessageEventContent
				gc.convertMedia(ctx, portal, intent, att, &anotherPart)
				output.Parts = append(output.Parts, &bridgev2.ConvertedMessagePart{
					ID:      networkid.PartID(strconv.Itoa(i)),
					Type:    event.EventMessage,
					Content: &anotherPart,
				})
			}
		}
	} else {
		content.Body = msg.Text
	}
	return output, nil
}

func (gc *GVClient) convertMedia(ctx context.Context, portal *bridgev2.Portal, intent bridgev2.MatrixAPI, att *gvproto.Attachment, into *event.MessageEventContent) {
	if att.Status == gvproto.Attachment_NOT_SUPPORTED {
		addMediaFailure(into, "File type not supported by Google Voice")
		return
	}
	data, mimeType, err := gc.Client.DownloadAttachment(ctx, att.ID)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to download attachment")
		addMediaFailure(into, "Failed to download attachment")
		return
	}
	if att.MimeType != "" {
		mimeType = att.MimeType
	}
	if mimeType == "" {
		mimeType = http.DetectContentType(data)
	}
	origMeta := &gvproto.Attachment_Metadata{}
	for _, meta := range att.Metadata {
		if meta.Size == gvproto.Attachment_Metadata_ORIGINAL {
			origMeta = meta
			break
		} else if origMeta == nil || meta.Width > origMeta.Width {
			origMeta = meta
		}
	}
	prefix := strings.Split(mimeType, "/")[0]
	var msgtype event.MessageType
	switch prefix {
	case "image":
		msgtype = event.MsgImage
	case "audio":
		msgtype = event.MsgAudio
	case "video":
		msgtype = event.MsgVideo
	default:
		prefix = "file"
		msgtype = event.MsgFile
	}
	fileName := prefix + exmime.ExtensionFromMimetype(mimeType)
	into.URL, into.File, err = intent.UploadMedia(ctx, portal.MXID, data, fileName, mimeType)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to reupload attachment")
		addMediaFailure(into, "Failed to reupload attachment")
		return
	}
	into.FileName = fileName
	into.MsgType = msgtype
	into.Info = &event.FileInfo{
		MimeType: mimeType,
		Width:    int(origMeta.Width),
		Height:   int(origMeta.Height),
		Size:     len(data),
	}
}

func addMediaFailure(into *event.MessageEventContent, message string) {
	if into.Body != "" {
		into.Body = fmt.Sprintf("%s\n\n%s", message, into.Body)
		if into.FormattedBody != "" {
			into.FormattedBody = fmt.Sprintf("<p>%s</p><p>%s</p>", message, into.FormattedBody)
		}
	} else {
		into.Body = message
		into.MsgType = event.MsgNotice
	}
}
