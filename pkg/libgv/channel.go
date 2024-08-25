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
	return chooseServerResp.GSessionID, createChannelResp.GetData().GetSession(), nil
}

var noopSuffix = []byte(`,["noop"]]`)

func (c *Client) RunRealtimeChannel(ctx context.Context) error {
	gSessionID, channel, err := c.SubscribeRealtimeChannel(ctx)
	if err != nil {
		return err
	}
	log := zerolog.Ctx(ctx)
	var ackID uint64
	failedRequests := 0
	initialConnect := true
	for {
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
			if errors.As(err, &httpErr) && bytes.Contains(httpErr.Body, []byte("Unknown SID")) {
				failedRequests++
				if failedRequests > 10 {
					return errors.New("too many repeated unknown SID errors")
				}
				sleep := time.Duration(failedRequests-1) * 2 * time.Second
				log.Debug().Stringer("sleep_duration", sleep).Msg("Unknown SID error, re-subscribing")
				time.Sleep(sleep)
				gSessionID, channel, err = c.SubscribeRealtimeChannel(ctx)
				if err != nil {
					return fmt.Errorf("failed to re-subscribe after unknown SID: %w", err)
				}
				ackID = 0
				continue
			}
			return fmt.Errorf("failed to send channel request: %w", err)
		}
		if failedRequests != 0 || initialConnect {
			log.Debug().
				Int("failure_count", failedRequests).
				Bool("is_initial", initialConnect).
				Msg("Realtime long polling connected, sending event")
			c.dispatchEvent(ctx, &RealtimeConnected{})
		}
		reader := utf16chunk.NewReader(resp.Body)
		for {
			var chunk []byte
			chunk, err = reader.ReadChunk()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return fmt.Errorf("failed to read chunk: %w", err)
			}
			var entries []json.RawMessage
			err = json.Unmarshal(chunk, &entries)
			if err != nil {
				return fmt.Errorf("failed to parse chunk into list: %w", err)
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
					ackID = parsed.ArrayID
					c.dispatchEvent(ctx, &RealtimeEvent{WebChannelEvent: &parsed})
				}
				failedRequests = 0
				initialConnect = false
			}
		}
	}
}
