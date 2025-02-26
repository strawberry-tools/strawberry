// Copyright 2024 The Strawberry Tools team. All rights reserved.
// Copyright 2024 The Hugo Authors. All rights reserved.
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
	"runtime"

	"github.com/strawberry-tools/strawberry/common/hugo"

	"github.com/bep/simplecobra"
	"github.com/spf13/cobra"
)

func newEnvCommand() simplecobra.Commander {
	return &simpleCommand{
		name:  "env",
		short: "Print Strawberry version and environment info",
		long:  "Print Strawberry version and environment info. This is useful in Strawberry bug reports",
		run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
			r.Printf("%s\n", hugo.BuildVersionString())
			r.Printf("GOOS=%q\n", runtime.GOOS)
			r.Printf("GOARCH=%q\n", runtime.GOARCH)
			r.Printf("GOVERSION=%q\n", runtime.Version())

			if r.isVerbose() {
				deps := hugo.GetDependencyList()
				for _, dep := range deps {
					r.Printf("%s\n", dep)
				}
			} else {
				// These are also included in the GetDependencyList above;
				// always print these as these are most likely the most useful to know about.
				deps := hugo.GetDependencyListNonGo()
				for _, dep := range deps {
					r.Printf("%s\n", dep)
				}
			}
			return nil
		},
		withc: func(cmd *cobra.Command, r *rootCommand) {
			cmd.ValidArgsFunction = cobra.NoFileCompletions
		},
	}
}
