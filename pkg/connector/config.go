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

package connector

import (
	_ "embed"
	"strings"
	"text/template"

	up "go.mau.fi/util/configupgrade"
	"gopkg.in/yaml.v3"

	"go.mau.fi/mautrix-gvoice/pkg/libgv/gvproto"
)

//go:embed example-config.yaml
var ExampleConfig string

type Config struct {
	DisplaynameTemplate string `yaml:"displayname_template"`

	displaynameTemplate *template.Template `yaml:"-"`
}

type umConfig Config

func (c *Config) UnmarshalYAML(node *yaml.Node) error {
	err := node.Decode((*umConfig)(c))
	if err != nil {
		return err
	}

	c.displaynameTemplate, err = template.New("displayname").Parse(c.DisplaynameTemplate)
	if err != nil {
		return err
	}
	return nil
}

type DisplaynameTemplateArgs struct {
	PhoneNumber string
	Name        string
	Contact     *ProcessedContact
}

func (gv *GVConnector) GetConfig() (example string, data any, upgrader up.Upgrader) {
	return ExampleConfig, &gv.Config, up.SimpleUpgrader(upgradeConfig)
}

func (c Config) FormatDisplayname(info *gvproto.Contact, contact *ProcessedContact) string {
	var buf strings.Builder
	_ = c.displaynameTemplate.Execute(&buf, DisplaynameTemplateArgs{
		PhoneNumber: info.PhoneNumber,
		Name:        info.Name,
		Contact:     contact,
	})
	return buf.String()
}

func upgradeConfig(helper up.Helper) {
	helper.Copy(up.Str, "displayname_template")
}
