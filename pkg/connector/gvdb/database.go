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

package gvdb

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/rs/zerolog"
	"go.mau.fi/util/dbutil"
	"maunium.net/go/mautrix/bridgev2/networkid"
)

type GVDB struct {
	*dbutil.Database
}

var table = dbutil.BuildUpgradeTable().WithFS(upgrades).Finish()

//go:embed *.sql
var upgrades embed.FS

func New(db *dbutil.Database, log zerolog.Logger) *GVDB {
	db = db.Child("gvoice_version", table, dbutil.ZeroLogger(log))
	return &GVDB{
		Database: db,
	}
}

// TextPortal is a 1:1 text conversation portal known to the bridge, used to
// seed the call/voicemail merge index with text threads that the Google Voice
// list API no longer returns (it only lists the most recent threads).
type TextPortal struct {
	ID           networkid.PortalID
	Participants []string
}

// GetTextPortals returns the bridge's own text portals for a login, with the
// participants stored in their metadata. Call threads for a participant whose
// text thread has fallen off the fetched thread list can still be merged into
// the portal using this data.
func (db *GVDB) GetTextPortals(ctx context.Context, loginID networkid.UserLoginID) ([]TextPortal, error) {
	rows, err := db.Query(ctx, `
		SELECT id, metadata
		FROM portal
		WHERE id LIKE 't.%' AND receiver = $1
	`, loginID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var portals []TextPortal
	for rows.Next() {
		var (
			id       networkid.PortalID
			metadata string
		)
		if err = rows.Scan(&id, &metadata); err != nil {
			return nil, err
		}
		var meta struct {
			Participants []string `json:"participants"`
		}
		if err = json.Unmarshal([]byte(metadata), &meta); err != nil {
			return nil, fmt.Errorf("failed to parse portal metadata: %w", err)
		}
		portals = append(portals, TextPortal{ID: id, Participants: meta.Participants})
	}
	return portals, rows.Err()
}

func (db *GVDB) GetLoginPrefix(ctx context.Context, loginID networkid.UserLoginID) (string, error) {
	var rowID int64
	err := db.QueryRow(ctx, `
		INSERT INTO gvoice_login_prefix (login_id)
		VALUES ($1)
		ON CONFLICT (login_id) DO UPDATE SET login_id=gvoice_login_prefix.login_id
		RETURNING prefix
	`, loginID).Scan(&rowID)
	return strconv.FormatInt(rowID, 10), err
}
