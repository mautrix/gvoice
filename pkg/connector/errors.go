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
	"errors"
	"fmt"
	"net/http"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/event"

	"go.mau.fi/mautrix-gvoice/pkg/libgv"
)

func isRateLimitError(err error) bool {
	var respErr *libgv.ResponseError
	return errors.As(err, &respErr) && respErr.Resp != nil && respErr.Resp.StatusCode == http.StatusTooManyRequests
}

// wrapSendError turns an error from a Google Voice request made in a Matrix event handler into a message status.
func wrapSendError(err error) error {
	if err == nil || errors.As(err, &bridgev2.MessageStatus{}) {
		return err
	}
	var respErr *libgv.ResponseError
	if !errors.As(err, &respErr) || respErr.Resp == nil {
		return bridgev2.WrapErrorInStatus(err).
			WithMessage("failed to reach Google Voice").
			WithErrorReason(event.MessageStatusGenericError).
			WithSendNotice(true)
	}
	switch code := respErr.Resp.StatusCode; {
	case code == http.StatusTooManyRequests:
		return bridgev2.WrapErrorInStatus(err).
			WithMessage("Google Voice is rate limiting you, try again in a few minutes").
			WithErrorReason(event.MessageStatusNetworkError).
			WithIsCertain(true).
			WithSendNotice(true)
	case respErr.IsAuthError():
		return bridgev2.WrapErrorInStatus(err).
			WithMessage("Google Voice rejected the request, you may need to log in again").
			WithErrorReason(event.MessageStatusNetworkError).
			WithIsCertain(true).
			WithSendNotice(true)
	case code == http.StatusNotFound:
		return bridgev2.WrapErrorInStatus(err).
			WithMessage("conversation not found on Google Voice").
			WithErrorReason(event.MessageStatusNetworkError).
			WithStatus(event.MessageStatusFail).
			WithIsCertain(true).
			WithSendNotice(true)
	case code >= 400 && code < 500:
		return bridgev2.WrapErrorInStatus(err).
			WithMessage("Google Voice rejected the request").
			WithErrorReason(event.MessageStatusNetworkError).
			WithIsCertain(true).
			WithSendNotice(true)
	default:
		return bridgev2.WrapErrorInStatus(err).
			WithMessage("Google Voice server error").
			WithErrorReason(event.MessageStatusNetworkError).
			WithSendNotice(true)
	}
}

// wrapRespError turns an error from a Google Voice request made during login or provisioning into an HTTP error,
// keeping the original error in the chain. Clients retry 429 and 502-504 automatically, so those aren't used.
func wrapRespError(err error) error {
	var mapped bridgev2.RespError
	var respErr *libgv.ResponseError
	switch {
	case libgv.IsAuthError(err):
		mapped = bridgev2.RespError(mautrix.MForbidden.WithMessage("Google rejected the cookies. Log in again with fresh cookies from a logged-in voice.google.com tab."))
	case isRateLimitError(err):
		mapped = bridgev2.RespError(mautrix.MLimitExceeded.WithMessage("Google Voice is rate limiting requests. Please wait a few minutes and try again.").WithStatus(http.StatusBadRequest))
	case errors.As(err, &respErr):
		mapped = bridgev2.RespError(mautrix.MUnknown.WithMessage("Google Voice returned an error"))
	default:
		mapped = bridgev2.RespError(mautrix.MUnknown.WithMessage("Failed to reach Google Voice"))
	}
	return fmt.Errorf("%w: %w", mapped, err)
}
