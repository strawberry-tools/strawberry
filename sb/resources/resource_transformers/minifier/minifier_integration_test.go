// Copyright 2021 The Hugo Authors. All rights reserved.
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

package minifier_test

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/strawberry-tools/strawberry/hugolib"
)

// Issue 8954
func TestTransformMinify(t *testing.T) {
	c := qt.New(t)

	files := `
-- assets/js/test.js --
new Date(2002, 04, 11)
-- config.toml --
-- layouts/index.html --
{{ $js := resources.Get "js/test.js" | minify }}
<script>
{{ $js.Content }}
</script>
`

	b, err := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           c,
			TxtarString: files,
		},
	).BuildE()

	b.Assert(err, qt.IsNotNil)
	b.Assert(err, qt.ErrorMatches, "(?s).*legacy octal numbers.*line 1.*")
}
