package connector

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"

	"go.mau.fi/mautrix-gvoice/pkg/libgv"
	"go.mau.fi/mautrix-gvoice/pkg/libgv/gvproto"
)

type mediaTransport func(*http.Request) (*http.Response, error)

func (f mediaTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type mediaUpload struct {
	bridgev2.MatrixAPI
	data []byte
	mime string
}

func (m *mediaUpload) UploadMedia(_ context.Context, _ id.RoomID, data []byte, _ string, mimeType string) (id.ContentURIString, *event.EncryptedFileInfo, error) {
	m.data, m.mime = data, mimeType
	return "mxc://test/upload", nil, nil
}

func TestConvertMediaReupload(t *testing.T) {
	for _, tc := range []struct {
		name       string
		direct     bool
		attachment string
		header     string
	}{
		{"unknown type with header", true, "", "image/png"},
		{"unknown type without header", true, "", ""},
		{"direct media disabled", false, "image/png", "application/octet-stream"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data := []byte("\x89PNG\r\n\x1a\nfixture")
			requests := 0
			client := libgv.NewClient(nil)
			client.HTTP.Transport = mediaTransport(func(req *http.Request) (*http.Response, error) {
				requests++
				if req.Method != http.MethodGet {
					t.Fatalf("unexpected method: %s", req.Method)
				}
				return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {tc.header}}, Body: io.NopCloser(bytes.NewReader(data)), Request: req}, nil
			})
			gc := &GVClient{Main: &GVConnector{DirectMedia: tc.direct}, Client: client}
			portal := &bridgev2.Portal{Portal: &database.Portal{MXID: "!room:test"}}
			upload := &mediaUpload{}
			var content event.MessageEventContent
			gc.convertMedia(context.Background(), portal, upload, &gvproto.Attachment{ID: "attachment", MimeType: tc.attachment}, &content)
			if requests != 1 || !bytes.Equal(upload.data, data) || upload.mime != "image/png" {
				t.Fatalf("unexpected reupload: requests=%d mime=%q data=%x", requests, upload.mime, upload.data)
			}
			if content.URL != "mxc://test/upload" || content.MsgType != event.MsgImage || content.Info.MimeType != "image/png" || content.Info.Size != len(data) {
				t.Fatalf("unexpected content: %+v", content)
			}
		})
	}
}
