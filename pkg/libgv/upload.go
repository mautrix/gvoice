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
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type ReqStartUploadFieldInlined struct {
	Name        string `json:"name"`
	Content     string `json:"content"`
	ContentType string `json:"contentType"`
}

type ReqStartUploadFieldExternal struct {
	Name     string   `json:"name"`
	FileName string   `json:"filename"`
	Put      struct{} `json:"put"`
	Size     int      `json:"size"`
}

type ReqStartUploadField struct {
	Inlined  *ReqStartUploadFieldInlined  `json:"inlined,omitempty"`
	External *ReqStartUploadFieldExternal `json:"external,omitempty"`
}

type ReqStartUpload struct {
	ProtocolVersion      string `json:"protocolVersion"` // "0.8"
	CreateSessionRequest struct {
		Fields []ReqStartUploadField `json:"fields"`
	} `json:"createSessionRequest"`
}

type RespFinishUpload struct {
	SessionStatus struct {
		State                  string `json:"state"`
		ExternalFieldTransfers []struct {
			Name             string `json:"name"`
			Status           string `json:"status"`
			BytesTransferred int    `json:"bytesTransferred"`
			BytesTotal       int    `json:"bytesTotal"`
			PutInfo          struct {
				URL string `json:"url"`
			} `json:"putInfo"`
			ContentType string `json:"content_type"`
		} `json:"externalFieldTransfers"`
		AdditionalInfo struct {
			UploaderServiceGoogleRupioAdditionalInfo struct {
				CompletionInfo struct {
					Status               string `json:"status"`
					CustomerSpecificInfo struct {
						Kind           string  `json:"kind"`
						AlbumID        string  `json:"albumid"`
						PhotoID        string  `json:"photoid"`
						PhotoMediaKey  string  `json:"photoMediaKey"`
						AlbumMediaKey  string  `json:"albumMediaKey"`
						Width          int     `json:"width"`
						Height         int     `json:"height"`
						URL            string  `json:"url"`
						Title          string  `json:"title"`
						Description    string  `json:"description"`
						Username       string  `json:"username"`
						PhotoPageURL   string  `json:"photoPageUrl"`
						AlbumPageURL   string  `json:"albumPageUrl"`
						Rotation       int     `json:"rotation"`
						MimeType       string  `json:"mimeType"`
						Timestamp      float64 `json:"timestamp"`
						AutoDownsized  bool    `json:"autoDownsized"`
						AutoEnhance    bool    `json:"autoEnhance"`
						Position       float64 `json:"position"`
						ImageVersion   int     `json:"imageVersion"`
						IsPhotoDeduped bool    `json:"isPhotoDeduped"`
					} `json:"customerSpecificInfo"`
				} `json:"completionInfo"`
			} `json:"uploader_service.GoogleRupioAdditionalInfo"`
		} `json:"additionalInfo"`
		UploadID string `json:"upload_id"`
	} `json:"sessionStatus"`
}

func (c *Client) UploadPhoto(ctx context.Context, fileName, mimeType string, photo []byte) (string, error) {
	startReq := &ReqStartUpload{ProtocolVersion: "0.8"}
	startReq.CreateSessionRequest.Fields = []ReqStartUploadField{{
		External: &ReqStartUploadFieldExternal{
			Name:     "file",
			FileName: fileName,
			Size:     len(photo),
		},
	}, {
		Inlined: &ReqStartUploadFieldInlined{
			Name:        "title",
			Content:     fileName,
			ContentType: "text/plain",
		},
	}, {
		Inlined: &ReqStartUploadFieldInlined{
			Name:        "addtime",
			Content:     strconv.FormatInt(time.Now().UnixMilli(), 10),
			ContentType: "text/plain",
		},
	}, {
		Inlined: &ReqStartUploadFieldInlined{
			Name:        "onepick_version",
			Content:     "v2",
			ContentType: "text/plain",
		},
	}, {
		Inlined: &ReqStartUploadFieldInlined{
			Name:        "onepick_host_id",
			Content:     "36", // TODO what is this?
			ContentType: "text/plain",
		},
	}, {
		Inlined: &ReqStartUploadFieldInlined{
			Name:        "album_mode",
			Content:     "temporary",
			ContentType: "text/plain",
		},
	}, {
		Inlined: &ReqStartUploadFieldInlined{
			Name:        "silo_id",
			Content:     "26", // TODO what is this?
			ContentType: "text/plain",
		},
	}}
	startReqBytes, err := json.Marshal(startReq)
	if err != nil {
		return "", err
	}
	resp, err := c.MakeRequest(ctx, http.MethodPost, EndpointUpload, url.Values{
		"authuser": {c.AuthUser},
		"opi":      {UploadOPI},
	}, http.Header{
		"Content-Type":                        {ContentTypeFormData},
		"X-Goog-Upload-Command":               {"start"},
		"X-Goog-Upload-Header-Content-Length": {strconv.Itoa(len(photo))},
		"X-Goog-Upload-Header-Content-Type":   {mimeType},
		"X-Goog-Upload-Protocol":              {"resumable"},
	}, startReqBytes)
	if err != nil {
		return "", err
	}
	_ = resp.Body.Close()
	uploadURL := resp.Header.Get("X-Goog-Upload-URL")
	if uploadURL == "" {
		return "", fmt.Errorf("missing X-Goog-Upload-URL header")
	}
	respData, err := ReadJSONResponse[*RespFinishUpload](
		c.MakeRequest(ctx, http.MethodPost, uploadURL, nil, http.Header{
			"Content-Type":          {mimeType},
			"X-Goog-Upload-Command": {"upload, finalize"},
			"X-Goog-Upload-Offset":  {"0"},
		}, bytes.NewReader(photo)),
	)
	if err != nil {
		return "", err
	}
	uploadedURL := respData.SessionStatus.AdditionalInfo.UploaderServiceGoogleRupioAdditionalInfo.CompletionInfo.CustomerSpecificInfo.URL
	if uploadedURL == "" {
		return "", fmt.Errorf("missing uploaded URL")
	}
	return uploadedURL, nil
}
