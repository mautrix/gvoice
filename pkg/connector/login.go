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
	"slices"

	"github.com/rs/zerolog"
	"maunium.net/go/mautrix/bridge/status"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/networkid"

	"go.mau.fi/mautrix-gvoice/pkg/libgv"
)

const (
	LoginFlowIDCookies = "cookies"

	LoginStepIDCookies  = "fi.mau.gvoice.cookies"
	LoginStepIDComplete = "fi.mau.gvoice.complete"
)

var loginFlows = []bridgev2.LoginFlow{{
	Name:        "Cookies",
	Description: "Log in with your Google account by providing cookies from voice.google.com",
	ID:          LoginFlowIDCookies,
}}

func (gv *GVConnector) GetLoginFlows() []bridgev2.LoginFlow {
	return loginFlows
}

func (gv *GVConnector) CreateLogin(ctx context.Context, user *bridgev2.User, flowID string) (bridgev2.LoginProcess, error) {
	if flowID != LoginFlowIDCookies {
		return nil, fmt.Errorf("unknown login flow ID %q", flowID)
	}
	return &GVLogin{User: user}, nil
}

type GVLogin struct {
	User *bridgev2.User
}

var _ bridgev2.LoginProcessCookies = (*GVLogin)(nil)

var cookieLoginStep = &bridgev2.LoginStep{
	Type:         bridgev2.LoginStepTypeCookies,
	StepID:       LoginStepIDCookies,
	Instructions: "Enter a JSON object with your cookies, or a cURL command copied from browser devtools.",
	CookiesParams: &bridgev2.LoginCookiesParams{
		URL:    "https://voice.google.com/signup",
		Fields: []bridgev2.LoginCookieField{},
	},
}

func init() {
	cookies := map[string]string{
		"OSID":             "voice.google.com",
		"COMPASS":          "clients6.google.com",
		"SID":              ".google.com",
		"HSID":             ".google.com",
		"SSID":             ".google.com",
		"APISID":           ".google.com",
		"SAPISID":          ".google.com",
		"__Secure-1PSIDTS": ".google.com",
	}
	requiredCookies := []string{"SID", "HSID", "SSID", "APISID", "SAPISID"}
	cookieLoginStep.CookiesParams.Fields = make([]bridgev2.LoginCookieField, len(cookies))
	i := 0
	for cookie, domain := range cookies {
		cookieLoginStep.CookiesParams.Fields[i] = bridgev2.LoginCookieField{
			ID:       cookie,
			Required: slices.Contains(requiredCookies, cookie),
			Sources: []bridgev2.LoginCookieFieldSource{{
				Type:         bridgev2.LoginCookieTypeCookie,
				Name:         cookie,
				CookieDomain: domain,
			}},
		}
		i++
	}
}

func (gl *GVLogin) Start(ctx context.Context) (*bridgev2.LoginStep, error) {
	return cookieLoginStep, nil
}

func (gl *GVLogin) Cancel() {}

func (gl *GVLogin) SubmitCookies(ctx context.Context, cookies map[string]string) (*bridgev2.LoginStep, error) {
	cli := libgv.NewClient(cookies)
	acc, err := cli.GetAccount(ctx)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("Failed to get account")
		return nil, err
	}
	// TODO is google account email available somehow?
	ul, err := gl.User.NewLogin(ctx, &database.UserLogin{
		ID:         networkid.UserLoginID(fmt.Sprintf("%s|%s", gl.User.MXID, acc.Account.PrimaryDestinationID)),
		RemoteName: acc.Account.PrimaryDestinationID,
		RemoteProfile: status.RemoteProfile{
			Phone: acc.Account.PrimaryDestinationID,
		},
		Metadata: &UserLoginMetadata{
			Cookies: cli.GetCookies(),
		},
	}, &bridgev2.NewLoginParams{
		DeleteOnConflict: false,
	})
	if err != nil {
		return nil, err
	}
	go ul.Client.(*GVClient).connectRealtime()
	return &bridgev2.LoginStep{
		Type:         bridgev2.LoginStepTypeComplete,
		StepID:       LoginStepIDComplete,
		Instructions: fmt.Sprintf("Successfully logged in as %s", acc.Account.PrimaryDestinationID),
		CompleteParams: &bridgev2.LoginCompleteParams{
			UserLoginID: ul.ID,
			UserLogin:   ul,
		},
	}, nil
}
