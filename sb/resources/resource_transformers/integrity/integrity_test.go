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

package integrity

import (
	"context"
	"testing"

	"github.com/strawberry-tools/strawberry/config/testconfig"
	"github.com/strawberry-tools/strawberry/resources/resource"

	qt "github.com/frankban/quicktest"
	"github.com/strawberry-tools/strawberry/resources/resource_transformers/htesting"
)

func TestHashFromAlgo(t *testing.T) {
	for _, algo := range []struct {
		name string
		bits int
	}{
		{"md5", 128},
		{"sha256", 256},
		{"sha384", 384},
		{"sha512", 512},
		{"shaman", -1},
	} {
		t.Run(algo.name, func(t *testing.T) {
			c := qt.New(t)
			h, err := newHash(algo.name)
			if algo.bits > 0 {
				c.Assert(err, qt.IsNil)
				c.Assert(h.Size(), qt.Equals, algo.bits/8)
			} else {
				c.Assert(err, qt.Not(qt.IsNil))
				c.Assert(err.Error(), qt.Contains, "use either md5, sha256, sha384 or sha512")
			}
		})
	}
}

func TestTransform(t *testing.T) {
	c := qt.New(t)

	d := testconfig.GetTestDeps(nil, nil)
	t.Cleanup(func() { c.Assert(d.Close(), qt.IsNil) })

	client := New(d.ResourceSpec)

	r, err := htesting.NewResourceTransformerForSpec(d.ResourceSpec, "hugo.txt", "Hugo Rocks!")
	c.Assert(err, qt.IsNil)

	transformed, err := client.Fingerprint(r, "")

	c.Assert(err, qt.IsNil)
	c.Assert(transformed.RelPermalink(), qt.Equals, "/hugo.a5ad1c6961214a55de53c1ce6e60d27b6b761f54851fa65e33066460dfa6a0db.txt")
	c.Assert(transformed.Data(), qt.DeepEquals, map[string]any{"Integrity": "sha256-pa0caWEhSlXeU8HObmDSe2t2H1SFH6ZeMwZkYN+moNs="})
	content, err := transformed.(resource.ContentProvider).Content(context.Background())
	c.Assert(err, qt.IsNil)
	c.Assert(content, qt.Equals, "Hugo Rocks!")
}
