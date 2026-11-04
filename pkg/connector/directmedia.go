package connector

import (
	"context"
	"encoding/binary"
	"fmt"

	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/mediaproxy"
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
		return nil, err
	} else if login == nil || !login.Client.IsLoggedIn() {
		return nil, bridgev2.ErrNotLoggedIn
	}
	data, _, err := login.Client.(*GVClient).Client.DownloadAttachment(ctx, attachmentID)
	if err != nil {
		return nil, err
	}
	return mediaproxy.GetMediaResponseRawData(data), nil
}
