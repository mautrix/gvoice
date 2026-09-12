package libgv

import (
	"context"
	"net/http"

	"go.mau.fi/mautrix-gvoice/pkg/libgv/gvproto"
)

func (c *Client) RegisterPushNotifications(ctx context.Context, deviceID, token string) error {
	_, err := ReadProtoResponse[*gvproto.RespRegisterDestination](
		c.MakeRequest(ctx, http.MethodPost, APIBaseURL+"/api2notifications/registerdestination", nil, http.Header{
			"Content-Type": {ContentTypeProtobuf},
			"X-Gv-Rpc-Id":  {"1639"},
		}, &gvproto.ReqRegisterDestination{
			Destination: &gvproto.ReqRegisterDestination_Destination{
				DeviceID: deviceID,
				Fcm:      &gvproto.ReqRegisterDestination_Destination_FCM{Token: token},
			},
		}),
	)
	return err
}
