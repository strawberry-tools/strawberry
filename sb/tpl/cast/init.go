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

package cast

import (
	"context"

	"github.com/strawberry-tools/strawberry/deps"
	"github.com/strawberry-tools/strawberry/tpl/internal"
)

const name = "cast"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New()

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(cctx context.Context, args ...any) (any, error) { return ctx, nil },
		}

		ns.AddMethodMapping(ctx.ToInt,
			[]string{"int"},
			[][2]string{
				{`{{ "1234" | int | printf "%T" }}`, `int`},
			},
		)

		ns.AddMethodMapping(ctx.ToString,
			[]string{"string"},
			[][2]string{
				{`{{ 1234 | string | printf "%T" }}`, `string`},
			},
		)

		ns.AddMethodMapping(ctx.ToFloat,
			[]string{"float"},
			[][2]string{
				{`{{ "1234" | float | printf "%T" }}`, `float64`},
			},
		)

		return ns
	}

	internal.AddTemplateFuncsNamespace(f)
}
