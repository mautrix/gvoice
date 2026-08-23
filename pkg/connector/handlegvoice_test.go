package connector

import (
	"testing"

	"github.com/stretchr/testify/require"
	"maunium.net/go/mautrix/bridgev2/networkid"

	"go.mau.fi/mautrix-gvoice/pkg/connector/gvdb"
	"go.mau.fi/mautrix-gvoice/pkg/libgv/gvproto"
)

func TestCanonicalThreadIDUsesTextThreadForCalls(t *testing.T) {
	textThreadID := canonicalThreadID(&gvproto.Thread{
		ID:           "text-thread",
		IsText:       true,
		PhoneNumbers: []string{"+15551234567"},
	}, map[string]string{"+15551234567": "text-thread"})
	require.Equal(t, "text-thread", textThreadID)

	callThreadID := canonicalThreadID(&gvproto.Thread{
		ID:           "call-thread",
		PhoneNumbers: []string{"+15551234567"},
		Folders:      []gvproto.ThreadFolder{gvproto.ThreadFolder_ALL_CALL_THREADS},
	}, map[string]string{"+15551234567": "text-thread"})
	require.Equal(t, "text-thread", callThreadID)
}

func TestCanonicalThreadIDUsesTextThreadForVoicemail(t *testing.T) {
	voicemailThreadID := canonicalThreadID(&gvproto.Thread{
		ID:           "voicemail-thread",
		PhoneNumbers: []string{"+15551234567"},
		Folders:      []gvproto.ThreadFolder{gvproto.ThreadFolder_ALL_VOICEMAIL_THREADS},
	}, map[string]string{"+15551234567": "text-thread"})
	require.Equal(t, "text-thread", voicemailThreadID)
}

func TestCanonicalThreadIDLeavesStandaloneCallsAlone(t *testing.T) {
	callThreadID := canonicalThreadID(&gvproto.Thread{
		ID:           "call-thread",
		PhoneNumbers: []string{"+15551234567"},
		Folders:      []gvproto.ThreadFolder{gvproto.ThreadFolder_ALL_CALL_THREADS},
	}, nil)
	require.Equal(t, "call-thread", callThreadID)
}

func TestCanonicalThreadIDLeavesGroupThreadsAlone(t *testing.T) {
	// Multiple participants -> no unambiguous participant key -> never merged.
	callThreadID := canonicalThreadID(&gvproto.Thread{
		ID:           "group-call-thread",
		PhoneNumbers: []string{"+15551234567", "+15557654321"},
		Folders:      []gvproto.ThreadFolder{gvproto.ThreadFolder_ALL_CALL_THREADS},
	}, map[string]string{"+15551234567": "text-thread"})
	require.Equal(t, "group-call-thread", callThreadID)
}

func TestThreadParticipantKeyFallsBackToContacts(t *testing.T) {
	thread := &gvproto.Thread{
		Contacts: []*gvproto.Contact{{
			PhoneNumber: "+15551234567",
		}},
	}

	require.Equal(t, "+15551234567", threadParticipantKey(thread))
}

func TestThreadParticipantKeyIgnoresGroupPseudoContacts(t *testing.T) {
	thread := &gvproto.Thread{
		Contacts: []*gvproto.Contact{
			{PhoneNumber: "Group Message.abc"},
			{PhoneNumber: "+15551234567"},
		},
	}

	require.Equal(t, "+15551234567", threadParticipantKey(thread))
}

func TestBuildTextThreadIndexPicksMostRecentThreadPerParticipant(t *testing.T) {
	index := buildTextThreadIndex([]*gvproto.Thread{
		{
			ID:           "older-text-thread",
			IsText:       true,
			PhoneNumbers: []string{"+15551234567"},
			Messages:     []*gvproto.Message{{Timestamp: 100}},
		},
		{
			ID:           "newer-text-thread",
			IsText:       true,
			PhoneNumbers: []string{"+15551234567"},
			Messages:     []*gvproto.Message{{Timestamp: 200}},
		},
		{
			// Not a text thread, must be ignored when building the index.
			ID:           "call-thread",
			PhoneNumbers: []string{"+15551234567"},
			Folders:      []gvproto.ThreadFolder{gvproto.ThreadFolder_ALL_CALL_THREADS},
		},
	})

	require.Equal(t, "newer-text-thread", index["+15551234567"])
}

func TestMakeMessageIDNamespacesMergedThreadMessages(t *testing.T) {
	msg := &gvproto.Message{ID: "msg-1"}

	require.Equal(t, networkid.MessageID("msg-1"), makeMessageID("text-thread", "text-thread", msg))
	require.Equal(t, networkid.MessageID("call-thread\nmsg-1"), makeMessageID("call-thread", "text-thread", msg))
}

func TestThreadCursorRoundTrip(t *testing.T) {
	cursor := encodeThreadCursor("call-thread", "next-page")
	threadID, token := decodeThreadCursor("text-thread", cursor)
	require.Equal(t, "call-thread", threadID)
	require.Equal(t, "next-page", token)

	threadID, token = decodeThreadCursor("text-thread", "legacy-token")
	require.Equal(t, "text-thread", threadID)
	require.Equal(t, "legacy-token", token)
}

func TestCanonicalThreadIDMergesIntoStaleDBTextPortal(t *testing.T) {
	// The GV list API only returns the most recent threads, so a contact with
	// an inactive text thread can be absent from the fetch response. The
	// bridge still knows their text portal from its own database, and call
	// threads for that contact must merge into it instead of creating a new
	// portal per call.
	callThread := &gvproto.Thread{
		ID:           "call-thread",
		PhoneNumbers: []string{"+15551234567"},
		Folders:      []gvproto.ThreadFolder{gvproto.ThreadFolder_ALL_CALL_THREADS},
	}
	fetched := buildTextThreadIndex([]*gvproto.Thread{callThread})
	merged := mergeTextThreadIndex(fetched, []gvdb.TextPortal{{
		ID:           "t.+15551234567",
		Participants: []string{"+15551234567"},
	}})
	require.Equal(t, "t.+15551234567", canonicalThreadID(callThread, merged))
}

func TestMergeTextThreadIndexPrefersFetchedOverDB(t *testing.T) {
	// A text thread seen in the fetch response is more current than the
	// database copy and must win for the same participant.
	fetched := map[string]string{"+11111111111": "t.+11111111111"}
	merged := mergeTextThreadIndex(fetched, []gvdb.TextPortal{
		{ID: "t.+11111111111-old", Participants: []string{"+11111111111"}},
		{ID: "t.+15551234567", Participants: []string{"+15551234567"}},
		{ID: "t.+12223334444", Participants: []string{"+12223334444", "+13334445555"}},
		{ID: "t.+19998887777", Participants: nil},
	})
	require.Equal(t, "t.+11111111111", merged["+11111111111"])
	require.Equal(t, "t.+15551234567", merged["+15551234567"])
	require.Equal(t, "t.+12223334444", merged["+12223334444"])
	require.Equal(t, "t.+12223334444", merged["+13334445555"])
	_, ok := merged["+19998887777"]
	require.False(t, ok)
}
