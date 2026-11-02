package connector

import (
	"context"
	"encoding/json"
	"fmt"

	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/mediaproxy"
)

var _ bridgev2.DirectMediableNetwork = (*GVConnector)(nil)

type directMediaID struct {
	LoginID      networkid.UserLoginID `json:"login_id"`
	AttachmentID string                `json:"attachment_id"`
}

func (gv *GVConnector) SetUseDirectMedia() {
	gv.DirectMedia = true
}

func (gv *GVConnector) Download(ctx context.Context, mediaID networkid.MediaID, params map[string]string) (mediaproxy.GetMediaResponse, error) {
	var info directMediaID
	if err := json.Unmarshal([]byte(mediaID), &info); err != nil || info.LoginID == "" || info.AttachmentID == "" {
		return nil, fmt.Errorf("invalid direct media ID")
	}
	login, err := gv.Bridge.GetExistingUserLoginByID(ctx, info.LoginID)
	if err != nil {
		return nil, err
	} else if login == nil || !login.Client.IsLoggedIn() {
		return nil, bridgev2.ErrNotLoggedIn
	}
	data, _, err := login.Client.(*GVClient).Client.DownloadAttachment(ctx, info.AttachmentID)
	if err != nil {
		return nil, err
	}
	return mediaproxy.GetMediaResponseRawData(data), nil
}
