package connector

import (
	"context"
	"encoding/json"
	"fmt"

	"go.mau.fi/util/random"
	"maunium.net/go/mautrix/bridgev2"
)

var (
	_ bridgev2.PushableNetworkAPI          = (*GVClient)(nil)
	_ bridgev2.BackgroundSyncingNetworkAPI = (*GVClient)(nil)
)

var pushConfig = &bridgev2.PushConfig{
	FCM: &bridgev2.FCMPushConfig{SenderID: "301778431048"},
}

func (gc *GVClient) GetPushConfigs() *bridgev2.PushConfig {
	return pushConfig
}

func (gc *GVClient) RegisterPushNotifications(ctx context.Context, pushType bridgev2.PushType, token string) error {
	if !gc.IsLoggedIn() {
		return bridgev2.ErrNotLoggedIn
	}
	if pushType != bridgev2.PushTypeFCM {
		return fmt.Errorf("unsupported push type %s", pushType)
	}
	if token == "" {
		return fmt.Errorf("empty push token")
	}
	meta := gc.UserLogin.Metadata.(*UserLoginMetadata)
	if meta.PushDeviceID == "" {
		meta.PushDeviceID = fmt.Sprintf("%X", random.Bytes(32))
		if err := gc.UserLogin.Save(ctx); err != nil {
			meta.PushDeviceID = ""
			return fmt.Errorf("failed to save push device ID: %w", err)
		}
	}
	return gc.Client.RegisterPushNotifications(ctx, meta.PushDeviceID, token)
}

func (gc *GVClient) ConnectBackground(ctx context.Context, params *bridgev2.ConnectBackgroundParams) error {
	if !gc.IsLoggedIn() {
		return bridgev2.ErrNotLoggedIn
	}
	var notification struct {
		Data struct {
			ThreadID string `json:"thread_id"`
		} `json:"data"`
	}
	if params != nil {
		_ = json.Unmarshal(params.RawData, &notification)
	}
	gc.loadInitialContacts(ctx)
	return gc.fetchNewMessages(ctx, notification.Data.ThreadID)
}
