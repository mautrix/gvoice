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

const UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36"
const CHUserAgent = `"Chromium";v="128", "Not;A=Brand";v="24", "Google Chrome";v="128"`
const CHPlatform = `"Linux"`
const ClientVersion = "665865172"
const JavaScriptUserAgent = "google-api-javascript-client/1.1.0"
const WaaXUserAgent = "grpc-web-javascript/0.1"

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
	APIKey        = "AIzaSyDTYc1N4xiODyrQYK0Kl6g_y279LjYkrBg"
	UploadOPI     = "111538494"
	WaaAPIKey     = "AIzaSyBGb5fGAyC-pRcRU6MUHb__b_vKha71HRE"
	WaaRequestKey = "/JR8jsAkqotcKsEKhXic"
)

const (
	Origin = "https://voice.google.com"

	APIDomain      = "clients6.google.com"
	RealtimeDomain = "signaler-pa." + APIDomain
	ContactsDomain = "peoplestack-pa." + APIDomain
	UploadDomain   = "docs.google.com"
	WaaDomain      = "waa-pa." + APIDomain

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

	ContactsBaseURL              = "https://" + ContactsDomain + "/$rpc/peoplestack.PeopleStackAutocompleteService"
	EndpointAutocompleteContacts = ContactsBaseURL + "/Autocomplete"
	EndpointLookupContacts       = ContactsBaseURL + "/Lookup"

	WaaBaseURL        = "https://" + WaaDomain + "/$rpc/google.internal.waa.v1.Waa"
	EndpointCreateWaa = WaaBaseURL + "/Create"
	EndpointPingWaa   = WaaBaseURL + "/Ping"

	PushAPIBaseURL   = "https://www.googleapis.com"
	PushRegistration = PushAPIBaseURL + "/voice/v1/voiceclient/api2notifications/registerdestination?alt=proto"
)
