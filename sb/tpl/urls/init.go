// Copyright 2017 The Hugo Authors. All rights reserved.
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

package urls

import (
	"context"

	"github.com/strawberry-tools/strawberry/deps"
	"github.com/strawberry-tools/strawberry/tpl/internal"
)

const name = "urls"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New(d)

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(cctx context.Context, args ...any) (any, error) { return ctx, nil },
		}

		ns.AddMethodMapping(ctx.AbsURL,
			[]string{"absURL"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.AbsLangURL,
			[]string{"absLangURL"},
			[][2]string{},
		)
		ns.AddMethodMapping(ctx.Ref,
			[]string{"ref"},
			[][2]string{},
		)
		ns.AddMethodMapping(ctx.RelURL,
			[]string{"relURL"},
			[][2]string{},
		)
		ns.AddMethodMapping(ctx.RelLangURL,
			[]string{"relLangURL"},
			[][2]string{},
		)
		ns.AddMethodMapping(ctx.RelRef,
			[]string{"relref"},
			[][2]string{},
		)
		ns.AddMethodMapping(ctx.URLize,
			[]string{"urlize"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Anchorize,
			[]string{"anchorize"},
			[][2]string{
				{`{{ "This is a title" | anchorize }}`, `this-is-a-title`},
			},
		)

		ns.AddMethodMapping(ctx.JoinPath,
			nil,
			[][2]string{
				{`{{ urls.JoinPath "https://example.org" "foo" }}`, `https://example.org/foo`},
				{`{{ urls.JoinPath (slice "a" "b") }}`, `a/b`},
			},
		)

		return ns
	}

	internal.AddTemplateFuncsNamespace(f)
}
