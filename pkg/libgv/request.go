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
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"go.mau.fi/util/pblite"
	"google.golang.org/protobuf/proto"

	"go.mau.fi/mautrix-gvoice/pkg/libgv/utf16chunk"
)

func SAPISIDHash(origin, sapisid string) string {
	// Copied from libgm - TODO: move to shared library
	ts := time.Now().Unix()
	hash := sha1.Sum([]byte(fmt.Sprintf("%d %s %s", ts, sapisid, origin)))
	return fmt.Sprintf("SAPISIDHASH %d_%x", ts, hash[:])
}

func (c *Client) prepareHeaders(req *http.Request) {
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("X-Goog-AuthUser", c.AuthUser)
	if req.URL.Host == APIDomain {
		req.Header.Set("X-Client-Version", ClientVersion)
		req.Header.Set("X-ClientDetails", ClientDetails)
		req.Header.Set("X-JavaScript-User-Agent", JavaScriptUserAgent)
		req.Header.Set("X-Requested-With", "XMLHttpRequest")
		req.Header.Set("X-Goog-Encode-Response-If-Executable", "base64")
	}
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	if req.URL.Host == RealtimeDomain {
		req.Header.Set("Sec-Fetch-Site", "same-site")
	} else {
		req.Header.Set("Sec-Fetch-Site", "same-origin")
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	if req.URL.Host == UploadDomain {
		req.Header.Set("Origin", "https://"+UploadDomain)
		req.Header.Set("Referer", "https://"+UploadDomain+"/")
	} else {
		req.Header.Set("Origin", Origin)
		req.Header.Set("Referer", Origin+"/")
	}
	c.cookiesLock.RLock()
	defer c.cookiesLock.RUnlock()
	for name, value := range c.cookies {
		req.AddCookie(&http.Cookie{
			Name:  name,
			Value: value,
		})
		if name == "SAPISID" {
			req.Header.Set("Authorization", SAPISIDHash(Origin, value))
		}
	}
}

var (
	errMissingContentType = errors.New("missing Content-Type header")
	errInvalidBodyType    = errors.New("invalid body type")
	errGetWithBody        = errors.New("GET requests can't have a body")
)

type ResponseError struct {
	Req  *http.Request
	Resp *http.Response
	Body []byte
}

func (re *ResponseError) Error() string {
	return fmt.Sprintf("unexpected status code %d", re.Resp.StatusCode)
}

func (c *Client) MakeRequest(ctx context.Context, method, baseAddr string, query url.Values, headers http.Header, body any) (*http.Response, error) {
	parsedAddr, err := url.Parse(baseAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	if strings.HasSuffix(parsedAddr.Host, APIDomain) {
		if query == nil {
			query = make(url.Values)
		}
		query.Set("key", APIKey)
		if parsedAddr.Host == APIDomain {
			query.Set("alt", "proto")
		}
	}
	if query != nil {
		parsedAddr.RawQuery = query.Encode()
	}
	if headers == nil {
		headers = make(http.Header)
	}

	var realBody io.Reader
	switch typedBody := body.(type) {
	case proto.Message:
		var bodyBytes []byte
		if headers.Get("Content-Type") == ContentTypeProtobuf {
			bodyBytes, err = proto.Marshal(typedBody)
		} else {
			bodyBytes, err = pblite.Marshal(typedBody)
			headers.Set("Content-Type", ContentTypePBLite)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		realBody = bytes.NewReader(bodyBytes)
	case url.Values:
		realBody = strings.NewReader(typedBody.Encode())
		headers.Set("Content-Type", ContentTypeFormData)
	case io.Reader:
		realBody = typedBody
	case string:
		realBody = strings.NewReader(typedBody)
	case []byte:
		realBody = bytes.NewReader(typedBody)
	case nil:
		// no request body
	default:
		return nil, fmt.Errorf("%w %T", errInvalidBodyType, body)
	}
	if realBody != nil {
		if headers.Get("Content-Type") == "" {
			return nil, errMissingContentType
		}
		if method == http.MethodGet {
			return nil, errGetWithBody
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, parsedAddr.String(), realBody)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare request: %w", err)
	}
	req.Header = headers
	c.prepareHeaders(req)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	} else if resp.StatusCode >= 400 {
		data, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		logEvt := zerolog.Ctx(ctx).Warn().Int("status_code", resp.StatusCode)
		if zerolog.Ctx(ctx).GetLevel() == zerolog.TraceLevel {
			if json.Valid(data) {
				logEvt.RawJSON("response_json", data)
			} else if len(data) < 4096 {
				logEvt.Bytes("response_data", data)
			} else {
				logEvt.Str("response_data", "response too long to log")
			}
		}
		logEvt.Msg("Request failed")
		return nil, &ResponseError{
			Req:  req,
			Resp: resp,
			Body: data,
		}
	}
	if strings.HasSuffix(parsedAddr.Host, APIDomain) {
		c.updateCookies(ctx, resp)
	}
	return resp, nil
}

func ReadProtoResponse[T proto.Message](resp *http.Response, err error) (T, error) {
	var out T
	if err != nil {
		return out, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	plainMime, _, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	safetyMime, _, _ := mime.ParseMediaType(resp.Header.Get("X-Goog-Safety-Content-Type"))
	realMime := plainMime
	if safetyMime != "" {
		realMime = safetyMime
	}
	var data []byte
	if realMime == ContentTypePlainText {
		data, err = utf16chunk.NewReader(resp.Body).ReadChunk()
		if err != nil {
			return out, fmt.Errorf("failed to read response chunk: %w", err)
		}
	} else {
		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return out, fmt.Errorf("failed to read response: %w", err)
		}
	}

	out = out.ProtoReflect().New().Interface().(T)
	switch realMime {
	case ContentTypeProtobuf:
		if plainMime == ContentTypePlainText {
			var n int
			n, err = base64.StdEncoding.Decode(data, data)
			if err != nil {
				return out, fmt.Errorf("failed to decode base64: %w", err)
			}
			data = data[:n]
		}
		err = proto.Unmarshal(data, out)
	case ContentTypePBLite, ContentTypePlainText:
		err = pblite.Unmarshal(data, out)
	default:
		return out, fmt.Errorf("unknown response content type %q", realMime)
	}
	if err != nil {
		return out, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return out, nil
}

func ReadJSONResponse[T any](resp *http.Response, err error) (T, error) {
	var out T
	if err != nil {
		return out, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	err = json.NewDecoder(resp.Body).Decode(&out)
	if err != nil {
		return out, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return out, nil
}
