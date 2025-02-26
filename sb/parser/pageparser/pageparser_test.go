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

package pageparser

import (
	"bytes"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/strawberry-tools/strawberry/parser/metadecoders"
)

func BenchmarkParse(b *testing.B) {
	start := `


---
title: "Front Matters"
description: "It really does"
---

This is some summary. This is some summary. This is some summary. This is some summary.

 <!--more-->


`
	input := []byte(start + strings.Repeat(strings.Repeat("this is text", 30)+"{{< myshortcode >}}This is some inner content.{{< /myshortcode >}}", 10))
	cfg := Config{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := parseBytes(input, cfg, lexIntroSection); err != nil {
			b.Fatal(err)
		}
	}
}

func TestFormatFromFrontMatterType(t *testing.T) {
	c := qt.New(t)
	for _, test := range []struct {
		typ    ItemType
		expect metadecoders.Format
	}{
		{TypeFrontMatterJSON, metadecoders.JSON},
		{TypeFrontMatterTOML, metadecoders.TOML},
		{TypeFrontMatterYAML, metadecoders.YAML},
		{TypeFrontMatterORG, metadecoders.ORG},
		{TypeIgnore, ""},
	} {
		c.Assert(FormatFromFrontMatterType(test.typ), qt.Equals, test.expect)
	}
}

func TestIsProbablyItemsSource(t *testing.T) {
	c := qt.New(t)

	input := ` {{< foo >}} `
	items, err := collectStringMain(input)
	c.Assert(err, qt.IsNil)

	c.Assert(IsProbablySourceOfItems([]byte(input), items), qt.IsTrue)
	c.Assert(IsProbablySourceOfItems(bytes.Repeat([]byte(" "), len(input)), items), qt.IsFalse)
	c.Assert(IsProbablySourceOfItems([]byte(`{{< foo >}}  `), items), qt.IsFalse)
	c.Assert(IsProbablySourceOfItems([]byte(``), items), qt.IsFalse)
}

func TestHasShortcode(t *testing.T) {
	c := qt.New(t)

	c.Assert(HasShortcode("{{< foo >}}"), qt.IsTrue)
	c.Assert(HasShortcode("aSDasd  SDasd aSD\n\nasdfadf{{% foo %}}\nasdf"), qt.IsTrue)
	c.Assert(HasShortcode("{{</* foo */>}}"), qt.IsFalse)
	c.Assert(HasShortcode("{{%/* foo */%}}"), qt.IsFalse)
}

func BenchmarkHasShortcode(b *testing.B) {
	withShortcode := strings.Repeat("this is text", 30) + "{{< myshortcode >}}This is some inner content.{{< /myshortcode >}}" + strings.Repeat("this is text", 30)
	withoutShortcode := strings.Repeat("this is text", 30) + "This is some inner content." + strings.Repeat("this is text", 30)
	b.Run("Match", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			HasShortcode(withShortcode)
		}
	})

	b.Run("NoMatch", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			HasShortcode(withoutShortcode)
		}
	})
}

func TestSummaryDividerStartingFromMain(t *testing.T) {
	c := qt.New(t)

	input := `aaa <!--more--> bbb`
	items, err := collectStringMain(input)
	c.Assert(err, qt.IsNil)

	c.Assert(items, qt.HasLen, 4)
	c.Assert(items[1].Type, qt.Equals, TypeLeadSummaryDivider)
}
