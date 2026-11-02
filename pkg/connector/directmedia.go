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

func (gv *GVConnector) Download(ctx context.Context, mediaID networkid.MediaID, params map[string]string) (mediaproxy.GetMediaResponse, error) {
	loginLen, n := binary.Uvarint(mediaID)
	if n <= 0 || loginLen == 0 || loginLen >= uint64(len(mediaID)-n) {
		return nil, fmt.Errorf("invalid direct media ID")
	}
	loginID := networkid.UserLoginID(mediaID[n : n+int(loginLen)])
	attachmentID := string(mediaID[n+int(loginLen):])
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
