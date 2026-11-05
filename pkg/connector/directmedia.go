package connector

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"net/http"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/mediaproxy"

	"go.mau.fi/mautrix-gvoice/pkg/libgv"
)

var _ bridgev2.DirectMediableNetwork = (*GVConnector)(nil)

func (gv *GVConnector) SetUseDirectMedia() {
	gv.DirectMedia = true
}

func makeDirectMediaID(loginID networkid.UserLoginID, attachmentID string) networkid.MediaID {
	mediaID := binary.AppendUvarint([]byte{1}, uint64(len(loginID)))
	mediaID = append(mediaID, loginID...)
	return networkid.MediaID(append(mediaID, attachmentID...))
}

func parseDirectMediaID(mediaID networkid.MediaID) (networkid.UserLoginID, string, error) {
	if len(mediaID) == 0 || mediaID[0] != 1 {
		return "", "", fmt.Errorf("unsupported direct media ID version")
	}
	mediaID = mediaID[1:]
	loginLen, n := binary.Uvarint(mediaID)
	if n <= 0 || loginLen == 0 || loginLen >= uint64(len(mediaID)-n) {
		return "", "", fmt.Errorf("invalid direct media ID")
	}
	return networkid.UserLoginID(mediaID[n : n+int(loginLen)]), string(mediaID[n+int(loginLen):]), nil
}

func (gv *GVConnector) Download(ctx context.Context, mediaID networkid.MediaID, params map[string]string) (mediaproxy.GetMediaResponse, error) {
	loginID, attachmentID, err := parseDirectMediaID(mediaID)
	if err != nil {
		return nil, err
	}
	login, err := gv.Bridge.GetExistingUserLoginByID(ctx, loginID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user login: %w", err)
	} else if login == nil || !login.Client.IsLoggedIn() {
		return nil, bridgev2.ErrNotLoggedIn
	}
	data, mimeType, err := login.Client.(*GVClient).Client.DownloadAttachment(ctx, attachmentID)
	if err != nil {
		var respErr *libgv.ResponseError
		if errors.As(err, &respErr) {
			switch respErr.Resp.StatusCode {
			case http.StatusNotFound:
				return nil, mautrix.MNotFound.WithMessage("Media not found on Google Voice")
			case http.StatusUnauthorized, http.StatusForbidden:
				return nil, mautrix.MForbidden.WithMessage("Google Voice denied access to the media")
			}
		}
		return nil, mautrix.MUnknown.WithMessage("Failed to download media from Google Voice").WithStatus(http.StatusBadGateway)
	}
	resp := mediaproxy.GetMediaResponseRawData(data)
	if mimeType != "" {
		resp.ContentType = mimeType
	}
	return resp, nil
}
