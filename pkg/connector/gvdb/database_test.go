package gvdb

import (
	"context"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"go.mau.fi/util/dbutil"

	"maunium.net/go/mautrix/bridgev2/networkid"
)

func TestGetTextPortals(t *testing.T) {
	base, err := dbutil.NewWithDialect(":memory:", "sqlite3")
	require.NoError(t, err)
	defer base.Close()
	db := New(base, zerolog.Nop())

	// The portal table belongs to bridgev2's schema; the gvoice section only
	// reads from it, so a minimal stand-in table suffices here.
	_, err = base.RawDB.Exec(`
		CREATE TABLE portal (
			id        TEXT NOT NULL,
			metadata  jsonb NOT NULL,
			receiver  TEXT NOT NULL
		)
	`)
	require.NoError(t, err)
	_, err = base.RawDB.Exec(`
		INSERT INTO portal (id, metadata, receiver) VALUES
			('t.+15551234567', '{"participants":["+15551234567"]}', 'login-1'),
			('c.PCIGVAWKCLJLEITCHINGUGSSWDIEAQBAQAUQAEKAB', '{"participants":["+15551234567"]}', 'login-1'),
			('t.+12223334444', '{"participants":["+12223334444","+13334445555"]}', 'login-1'),
			('t.+19998887777', '{"participants":null}', 'login-1'),
			('t.+17778889999', '{"participants":["+17778889999"]}', 'login-2'),
			('g.Group Message.abc', '{"participants":["+11111111111"]}', 'login-1')
	`)
	require.NoError(t, err)

	portals, err := db.GetTextPortals(context.Background(), "login-1")
	require.NoError(t, err)
	require.Len(t, portals, 3)
	byID := make(map[networkid.PortalID][]string, len(portals))
	for _, portal := range portals {
		byID[portal.ID] = portal.Participants
	}
	require.Equal(t, []string{"+15551234567"}, byID["t.+15551234567"])
	require.Equal(t, []string{"+12223334444", "+13334445555"}, byID["t.+12223334444"])
	require.Empty(t, byID["t.+19998887777"])
	require.NotContains(t, byID, "c.PCIGVAWKCLJLEITCHINGUGSSWDIEAQBAQAUQAEKAB")
	require.NotContains(t, byID, "g.Group Message.abc")
	require.NotContains(t, byID, "t.+17778889999")
}
