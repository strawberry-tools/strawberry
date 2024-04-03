// Copyright 2024 The Strawberry Tools team All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"context"
	"errors"

	"github.com/bep/simplecobra"
	"github.com/spf13/cobra"

	"github.com/strawberry-tools/strawberry/common/hugo"
)

var vType string

func newVersionCmd() simplecobra.Commander {

	return &simpleCommand{
		use:   "version",
		short: "Print the version number of Strawberry",
		long: `This will print the Strawberry version information. There
are flags available for scripting.`,
		run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {

			var theType hugo.VersionType

			if vType == "regular" {
				theType = hugo.VersionRegular
			} else if vType == "short" {
				theType = hugo.VersionShort
			} else if vType == "detailed" {
				theType = hugo.VersionDetailed
			} else {
				return errors.New("Invalid value for --type.")
			}

			r.Println(hugo.PrintStrawberryVersion(theType))

			return nil
		},
		withc: func(cmd *cobra.Command, r *rootCommand) {
			cmd.Flags().StringVarP(&vType, "type", "", "regular", "level of information to display: short, regular, or detailed")
		},
	}
}
