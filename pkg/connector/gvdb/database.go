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
	"strconv"

	"github.com/rs/zerolog"
	"go.mau.fi/util/dbutil"
	"maunium.net/go/mautrix/bridgev2/networkid"
)

type GVDB struct {
	*dbutil.Database
}

var table dbutil.UpgradeTable

//go:embed *.sql
var upgrades embed.FS

func init() {
	table.RegisterFS(upgrades)
}

func New(db *dbutil.Database, log zerolog.Logger) *GVDB {
	db = db.Child("gvoice_version", table, dbutil.ZeroLogger(log))
	return &GVDB{
		Database: db,
	}
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
