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

package libgv

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"go.mau.fi/util/pblite"

	"go.mau.fi/mautrix-gvoice/pkg/libgv/gvproto"
	"go.mau.fi/mautrix-gvoice/pkg/libgv/utf16chunk"
)

const reqChooseServer = `[[null,null,null,[7,5],null,[null,[null,1],[[["3"]]]]]]`

func (c *Client) SubscribeRealtimeChannel(ctx context.Context) (string, *gvproto.WebChannelSession, error) {
	if c == nil {
		return "", nil, fmt.Errorf("google Voice client is not initialized - ensure you call NewClient() with valid cookies before subscribing to channels")
	}
	chooseServerResp, err := ReadProtoResponse[*gvproto.RespChooseServer](
		c.MakeRequest(ctx, http.MethodPost, EndpointRealtimeChooseServer, nil, http.Header{
			"Content-Type": {ContentTypePBLite},
		}, reqChooseServer),
	)
	if err != nil {
		return "", nil, fmt.Errorf("failed to choose server: %w", err)
	}
	query := url.Values{
		"VER":        {"8"},
		"gsessionid": {chooseServerResp.GSessionID},
		"RID":        {strconv.Itoa(rand.IntN(100000))}, // TODO is this correct?
		"CVER":       {"22"},
		"t":          {"1"},
		//"zx":         {"???"},
	}
	extraHeaders := http.Header{
		"X-WebChannel-Content-Type": {ContentTypePBLite},
	}
	body := url.Values{
		"count": {"7"},
		"ofs":   {"0"},
		// TODO the first null could be current timestamp as microseconds?
		"req0___data__": {`[[["1",[null,null,null,[7,5],null,[null,[null,1],[[["2"]]]],null,1,2],null,3]]]`},
		"req1___data__": {`[[["2",[null,null,null,[7,5],null,[null,[null,1],[[["3"]]]],null,1,2],null,3]]]`},
		"req2___data__": {`[[["3",[null,null,null,[7,5],null,[null,[null,1],[[["3"]]]],null,1,2],null,3]]]`},
		"req3___data__": {`[[["4",[null,null,null,[7,5],null,[null,[null,1],[[["1"]]]],null,1,2],null,3]]]`},
		"req4___data__": {`[[["5",[null,null,null,[7,5],null,[null,[null,1],[[["1"]]]],null,1,2],null,3]]]`},
		"req5___data__": {`[[["6",[null,null,null,[7,5],null,[null,[null,1],[[["1"]]]],null,1,2],null,3]]]`},
		"req6___data__": {`[[["9",[null,null,null,[7,5],null,[null,[null,1],[[["1"]]]],null,1,2],null,3]]]`},
	}
	createChannelResp, err := ReadProtoResponse[*gvproto.RespCreateChannel](
		c.MakeRequest(ctx, http.MethodPost, EndpointRealtimeChannel, query, extraHeaders, body),
	)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create channel: %w", err)
	}
	if createChannelResp.GetData() == nil {
		return "", nil, fmt.Errorf("create channel response has no data - this may indicate authentication issues or Google Voice API changes")
	}
	if createChannelResp.GetData().GetSession() == nil {
		return "", nil, fmt.Errorf("create channel response has no session - verify your Google Voice account is properly configured")
	}
	return chooseServerResp.GSessionID, createChannelResp.GetData().GetSession(), nil
}

var ErrTooManyUnknownSID = errors.New("too many repeated unknown SID errors")

var noopSuffix = []byte(`,["noop"]]`)

const ForceResubscribeInterval = 1 * time.Hour

func (c *Client) RunRealtimeChannel(ctx context.Context) error {
	gSessionID, channel, err := c.SubscribeRealtimeChannel(ctx)
	if err != nil {
		return err
	}
	if channel == nil {
		return fmt.Errorf("channel is nil after subscription - failed to establish realtime connection, check network and authentication")
	}
	lastResubscribe := time.Now()
	log := zerolog.Ctx(ctx)
	var ackID uint64
	failedRequests := 0
	for {
		if time.Since(lastResubscribe) > ForceResubscribeInterval {
			log.Debug().Msg("Forcing re-subscribe as channel has been alive for too long")
			gSessionID, channel, err = c.SubscribeRealtimeChannel(ctx)
			if err != nil {
				return fmt.Errorf("failed to re-subscribe after timeout: %w", err)
			}
			if channel == nil {
				return fmt.Errorf("channel is nil after re-subscription - failed to re-establish realtime connection, verify authentication")
			}
			lastResubscribe = time.Now()
			ackID = 0
		}
		query := url.Values{
			"VER":        {"8"},
			"gsessionid": {gSessionID},
			"RID":        {"rpc"},
			"SID":        {channel.SessionID},
			"AID":        {strconv.FormatUint(ackID, 10)},
			"CI":         {"0"},
			"TYPE":       {"xmlhttp"},
			"t":          {"1"},
			//"zx":         {"???"},
		}
		log.Trace().
			Uint64("ack_id", ackID).
			Str("session_id", channel.SessionID).
			Str("gsessionid", gSessionID).
			Msg("Making new realtime long polling request")
		resp, err := c.MakeRequest(ctx, http.MethodGet, EndpointRealtimeChannel, query, nil, nil)
		if err != nil {
			var httpErr *ResponseError
			if errors.As(err, &httpErr) && httpErr.Body != nil && bytes.Contains(httpErr.Body, []byte("Unknown SID")) {
				failedRequests++
				if failedRequests > 10 {
					return ErrTooManyUnknownSID
				}
				sleep := time.Duration(failedRequests-1) * 2 * time.Second
				log.Debug().Stringer("sleep_duration", sleep).Msg("Unknown SID error, re-subscribing")
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(sleep):
				}
				gSessionID, channel, err = c.SubscribeRealtimeChannel(ctx)
				if err != nil {
					return fmt.Errorf("failed to re-subscribe after unknown SID: %w", err)
				}
				if channel == nil {
					return fmt.Errorf("channel is nil after re-subscription following unknown SID - connection lost, try restarting the bridge")
				}
				lastResubscribe = time.Now()
				ackID = 0
				continue
			}
			return fmt.Errorf("failed to send channel request: %w", err)
		}
		if failedRequests != 0 || ackID == 0 {
			log.Debug().
				Int("failure_count", failedRequests).
				Uint64("ack_id", ackID).
				Msg("Realtime long polling connected, sending event")
			c.dispatchEvent(ctx, &RealtimeConnected{})
		}
		reader := utf16chunk.NewReader(resp.Body)
	ReadLoop:
		for {
			var chunk []byte
			chunk, err = reader.ReadChunk()
			if err != nil {
				_ = resp.Body.Close()
				if errors.Is(err, io.EOF) {
					break
				} else if errors.Is(err, context.Canceled) {
					return ctx.Err()
				}
				log.Err(err).Msg("Failed to read chunk")
				break
			}
			var entries []json.RawMessage
			err = json.Unmarshal(chunk, &entries)
			if err != nil {
				_ = resp.Body.Close()
				log.Err(err).Msg("Failed to parse chunk into list")
				break
			}
			for i, entry := range entries {
				if bytes.HasSuffix(entry, noopSuffix) {
					var parsed gvproto.WebChannelNoopEvent
					err = pblite.Unmarshal(entry, &parsed)
					if err != nil {
						log.Debug().RawJSON("failed_entry", entry).Msg("Failed to parse noop channel entry")
						return fmt.Errorf("failed to parse entry #%d: %w", i+1, err)
					}
					ackID = parsed.ArrayID
				} else {
					var parsed gvproto.WebChannelEvent
					err = pblite.Unmarshal(entry, &parsed)
					if err != nil {
						log.Debug().RawJSON("failed_entry", entry).Msg("Failed to parse channel entry")
						return fmt.Errorf("failed to parse entry #%d: %w", i+1, err)
					}
					if len(parsed.GetDataWrapper()) == 1 && 
						parsed.GetDataWrapper()[0].GetAltData() != nil && 
						parsed.GetDataWrapper()[0].GetAltData().GetReconnect() {
						_ = resp.Body.Close()
						log.Debug().Msg("Got event that probably means we need to reconnect")
						gSessionID, channel, err = c.SubscribeRealtimeChannel(ctx)
						if err != nil {
							return fmt.Errorf("failed to re-subscribe after timeout: %w", err)
						}
						if channel == nil {
							return fmt.Errorf("channel is nil after re-subscription following reconnect event - persistent connection issues, check network stability")
						}
						lastResubscribe = time.Now()
						ackID = 0
						break ReadLoop
					}
					ackID = parsed.ArrayID
					c.dispatchEvent(ctx, &RealtimeEvent{WebChannelEvent: &parsed})
				}
				failedRequests = 0
			}
		}
	}
}
