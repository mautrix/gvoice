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
	"context"
	"fmt"
	"io"
	"maps"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"go.mau.fi/mautrix-gvoice/pkg/libgv/gvproto"
)

type Client struct {
	HTTP     *http.Client
	AuthUser string

	EventHandler func(context.Context, any)

	cookies     map[string]string
	cookiesLock sync.RWMutex
}

func NewClient(cookies map[string]string) *Client {
	if cookies == nil {
		cookies = make(map[string]string)
	}
	return &Client{
		HTTP:     &http.Client{Timeout: 120 * time.Second},
		cookies:  cookies,
		AuthUser: "0",
	}
}

func (c *Client) GetCookies() map[string]string {
	c.cookiesLock.RLock()
	defer c.cookiesLock.RUnlock()
	if c.cookies == nil {
		return make(map[string]string)
	}
	return maps.Clone(c.cookies)
}

func (c *Client) dispatchEvent(ctx context.Context, evt any) {
	if handler := c.EventHandler; handler != nil {
		handler(ctx, evt)
	}
}

func (c *Client) updateCookies(ctx context.Context, resp *http.Response) {
	respCookies := resp.Cookies()
	if len(respCookies) == 0 {
		return
	}
	c.cookiesLock.Lock()
	defer c.cookiesLock.Unlock()
	if c.cookies == nil {
		c.cookies = make(map[string]string)
	}
	cookiesChanged := false
	for _, cookie := range respCookies {
		if cookie.Expires.Before(time.Now()) || cookie.MaxAge < 0 {
			delete(c.cookies, cookie.Name)
			cookiesChanged = true
			continue
		}
		if c.cookies[cookie.Name] != cookie.Value {
			c.cookies[cookie.Name] = cookie.Value
			cookiesChanged = cookiesChanged || (cookie.Name != "__Secure-1PSIDCC" && cookie.Name != "__Secure-3PSIDCC" && cookie.Name != "SIDCC")
		}
	}
	if cookiesChanged {
		c.dispatchEvent(ctx, &CookieChanged{Cookies: maps.Clone(c.cookies)})
	}
}

func (c *Client) GetAccount(ctx context.Context) (*gvproto.RespGetAccount, error) {
	return ReadProtoResponse[*gvproto.RespGetAccount](
		c.MakeRequest(ctx, http.MethodPost, EndpointGetAccount, nil, nil, &gvproto.ReqGetAccount{
			UnknownInt2: 1,
		}),
	)
}

func GenerateTransactionID() int64 {
	return rand.Int64N(100000000000000)
}

func (c *Client) SendMessage(ctx context.Context, req *gvproto.ReqSendSMS) (*gvproto.RespSendSMS, error) {
	if req.TrackingData == nil {
		req.TrackingData = &gvproto.ReqSendSMS_TrackingData{Data: "!"}
	}
	if req.TransactionID == nil {
		req.TransactionID = &gvproto.ReqSendSMS_WrappedTxnID{ID: GenerateTransactionID()}
	}
	return ReadProtoResponse[*gvproto.RespSendSMS](
		c.MakeRequest(ctx, http.MethodPost, EndpointSendSMS, nil, nil, req),
	)
}

func (c *Client) UpdateThreadAttributes(ctx context.Context, req *gvproto.ReqUpdateAttributes) (*gvproto.RespUpdateAttributes, error) {
	return ReadProtoResponse[*gvproto.RespUpdateAttributes](
		c.MakeRequest(ctx, http.MethodPost, EndpointUpdateAttributes, nil, nil, req),
	)
}

func (c *Client) GetThread(ctx context.Context, threadID string, messageCount int, paginationToken string) (*gvproto.RespGetThread, error) {
	return ReadProtoResponse[*gvproto.RespGetThread](
		c.MakeRequest(ctx, http.MethodPost, EndpointGetThread, nil, nil, &gvproto.ReqGetThread{
			ThreadID:          threadID,
			MaybeMessageCount: int32(messageCount),
			PaginationToken:   paginationToken,
			UnknownWrapper: &gvproto.UnknownWrapper{
				UnknownInt2: 1,
				UnknownInt3: 1,
			},
		}),
	)
}

func (c *Client) ListThreads(ctx context.Context, folder gvproto.ThreadFolder, versionToken string) (*gvproto.RespListThreads, error) {
	req := &gvproto.ReqListThreads{
		Folder:       folder,
		UnknownInt2:  20,
		UnknownInt3:  15,
		VersionToken: versionToken,
		UnknownWrapper: &gvproto.UnknownWrapper{
			UnknownInt2: 1,
			UnknownInt3: 1,
		},
	}
	if req.VersionToken != "" {
		req.UnknownInt2 = 10
	}
	return ReadProtoResponse[*gvproto.RespListThreads](
		c.MakeRequest(ctx, http.MethodPost, EndpointListThreads, nil, nil, req),
	)
}

func (c *Client) DeleteThread(ctx context.Context, threadID string) (*gvproto.RespDeleteThread, error) {
	return ReadProtoResponse[*gvproto.RespDeleteThread](
		c.MakeRequest(ctx, http.MethodPost, EndpointDeleteThread, nil, nil, &gvproto.ReqDeleteThread{
			ThreadID: threadID,
		}),
	)
}

func (c *Client) DownloadAttachment(ctx context.Context, mediaID string) (data []byte, mime string, err error) {
	var resp *http.Response
	resp, err = c.MakeRequest(ctx, http.MethodGet, fmt.Sprintf(EndpointDownloadTemplate, c.AuthUser, mediaID), url.Values{
		"s": {strconv.Itoa(int(gvproto.Attachment_Metadata_ORIGINAL))},
	}, nil, nil)
	if err != nil {
		return
	}
	mime = resp.Header.Get("Content-Type")
	data, err = io.ReadAll(resp.Body)
	return
}

func (c *Client) AutocompleteContacts(ctx context.Context, query string) (*gvproto.RespAutocompleteContacts, error) {
	var maxResults int32 = 15
	if query == "" {
		maxResults = 500
	}
	return ReadProtoResponse[*gvproto.RespAutocompleteContacts](
		c.MakeRequest(ctx, http.MethodPost, EndpointAutocompleteContacts, nil, nil, &gvproto.ReqAutocompleteContacts{
			UnknownInt1:  243,
			Query:        query,
			UnknownInts3: []int32{1, 2},
			MaxResults:   maxResults,
		}),
	)
}

func (c *Client) LookupContact(ctx context.Context, phones ...string) (map[string]*gvproto.RespLookupContacts_Match, error) {
	req := &gvproto.ReqLookupContacts{
		UnknownInt1:  243,
		UnknownInts2: []int32{1, 2},
		Targets:      make([]*gvproto.ContactID, len(phones)),
	}
	for i, phone := range phones {
		req.Targets[i] = &gvproto.ContactID{Phone: phone}
	}
	resp, err := ReadProtoResponse[*gvproto.RespLookupContacts](
		c.MakeRequest(ctx, http.MethodPost, EndpointLookupContacts, nil, nil, req),
	)
	if err != nil {
		return nil, err
	}
	contacts := make(map[string]*gvproto.RespLookupContacts_Match, len(resp.Matches))
	for _, match := range resp.Matches {
		contacts[match.ID.Phone] = match
	}
	return contacts, nil
}

func (c *Client) CreateWaa(ctx context.Context) (*gvproto.CreatedWaa, error) {
	resp, err := ReadProtoResponse[*gvproto.RespCreateWaa](
		c.MakeRequest(ctx, http.MethodPost, EndpointCreateWaa, nil, nil, &gvproto.ReqCreateWaa{
			RequestKey: WaaRequestKey,
		}),
	)
	return resp.GetWaa(), err
}

func (c *Client) PingWaa(ctx context.Context, sig string, val int64) error {
	_, err := ReadProtoResponse[*gvproto.RespPingWaa](
		c.MakeRequest(ctx, http.MethodPost, EndpointPingWaa, nil, nil, &gvproto.ReqPingWaa{
			RequestKey: WaaRequestKey,
			Payload:    sig,
			I1:         72,
			I2:         val,
		}),
	)
	return err
}
