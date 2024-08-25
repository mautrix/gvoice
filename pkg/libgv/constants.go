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
	"net/url"
)

const UserAgent = "Mozilla/5.0 (X11; Linux x86_64; rv:130.0) Gecko/20100101 Firefox/130.0"
const ClientVersion = "665865172"
const JavaScriptUserAgent = "google-api-javascript-client/1.1.0"

var ClientDetails = url.Values{
	"appVersion": {"5.0 (X11)"},
	"platform":   {"Linux x86_64"},
	"userAgent":  {UserAgent},
}.Encode()

const (
	ContentTypeProtobuf  = "application/x-protobuf"
	ContentTypePBLite    = "application/json+protobuf"
	ContentTypeFormData  = "application/x-www-form-urlencoded"
	ContentTypePlainText = "text/plain"
)

const (
	APIKey    = "AIzaSyDTYc1N4xiODyrQYK0Kl6g_y279LjYkrBg"
	UploadOPI = "111538494"
)

const (
	Origin = "https://voice.google.com"

	APIDomain      = "clients6.google.com"
	RealtimeDomain = "signaler-pa." + APIDomain
	UploadDomain   = "docs.google.com"

	APIBaseURL                    = "https://" + APIDomain + "/voice/v1/voiceclient"
	EndpointGetAccount            = APIBaseURL + "/account/get"
	EndpointGetThread             = APIBaseURL + "/api2thread/get"
	EndpointListThreads           = APIBaseURL + "/api2thread/list"
	EndpointGetSIPRegisterInfo    = APIBaseURL + "/sipregisterinfo/get"
	EndpointGetThreadingInfo      = APIBaseURL + "/threadinginfo/get"
	EndpointSendSMS               = APIBaseURL + "/api2thread/sendsms"
	EndpointUpdateAttributes      = APIBaseURL + "/thread/updateattributes"
	EndpointBatchUpdateAttributes = APIBaseURL + "/thread/batchupdateattributes"
	EndpointMarkAllRead           = APIBaseURL + "/thread/markallread"
	EndpointDeleteThread          = APIBaseURL + "/thread/delete"

	RealtimeBaseURL              = "https://" + RealtimeDomain
	EndpointRealtimeChannel      = RealtimeBaseURL + "/punctual/multi-watch/channel"
	EndpointRealtimeChooseServer = RealtimeBaseURL + "/punctual/v1/chooseServer"

	EndpointUpload           = "https://" + UploadDomain + "/upload/photos/resumable"
	EndpointDownloadTemplate = "https://voice.google.com/u/%s/a/i/%s"
)
