// Copyright 2019 The Hugo Authors. All rights reserved.
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

package minifier

import (
	"context"
	"testing"

	"github.com/strawberry-tools/strawberry/config/testconfig"
	"github.com/strawberry-tools/strawberry/resources/resource"

	qt "github.com/frankban/quicktest"
	"github.com/strawberry-tools/strawberry/resources/resource_transformers/htesting"
)

func TestTransform(t *testing.T) {
	c := qt.New(t)

	d := testconfig.GetTestDeps(nil, nil)
	t.Cleanup(func() { c.Assert(d.Close(), qt.IsNil) })

	client, _ := New(d.ResourceSpec)
	r, err := htesting.NewResourceTransformerForSpec(d.ResourceSpec, "hugo.html", "<h1>   Hugo Rocks!   </h1>")
	c.Assert(err, qt.IsNil)

	transformed, err := client.Minify(r)
	c.Assert(err, qt.IsNil)

	c.Assert(transformed.RelPermalink(), qt.Equals, "/hugo.min.html")
	content, err := transformed.(resource.ContentProvider).Content(context.Background())
	c.Assert(err, qt.IsNil)
	c.Assert(content, qt.Equals, "<h1>Hugo Rocks!</h1>")
}
