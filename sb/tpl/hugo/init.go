// Copyright 2018 The Hugo Authors. All rights reserved.
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

// Package hugo provides template functions for accessing the Site Hugo object.
package hugo

import (
	"context"

	"github.com/strawberry-tools/strawberry/deps"
	"github.com/strawberry-tools/strawberry/tpl/internal"
)

const name = "hugo"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		if d.Site == nil {
			panic("no site in deps")
		}
		h := d.Site.Hugo()

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(cctx context.Context, args ...any) (any, error) { return h, nil },
		}

		// We just add the Hugo struct as the namespace here. No method mappings.

		return ns
	}

	internal.AddTemplateFuncsNamespace(f)
}
