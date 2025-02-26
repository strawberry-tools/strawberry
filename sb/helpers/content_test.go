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

package helpers_test

import (
	"bytes"
	"html/template"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/strawberry-tools/strawberry/config"
	"github.com/strawberry-tools/strawberry/helpers"
)

func TestTrimShortHTML(t *testing.T) {
	tests := []struct {
		markup string
		input  []byte
		output []byte
	}{
		{"markdown", []byte(""), []byte("")},
		{"markdown", []byte("Plain text"), []byte("Plain text")},
		{"markdown", []byte("<p>Simple paragraph</p>"), []byte("Simple paragraph")},
		{"markdown", []byte("\n  \n \t  <p> \t Whitespace\nHTML  \n\t </p>\n\t"), []byte("Whitespace\nHTML")},
		{"markdown", []byte("<p>Multiple</p><p>paragraphs</p>"), []byte("<p>Multiple</p><p>paragraphs</p>")},
		{"markdown", []byte("<p>Nested<p>paragraphs</p></p>"), []byte("<p>Nested<p>paragraphs</p></p>")},
		{"markdown", []byte("<p>Hello</p>\n<ul>\n<li>list1</li>\n<li>list2</li>\n</ul>"), []byte("<p>Hello</p>\n<ul>\n<li>list1</li>\n<li>list2</li>\n</ul>")},
		// Issue 11698
		{"markdown", []byte("<h2 id=`a`>b</h2>\n\n<p>c</p>"), []byte("<h2 id=`a`>b</h2>\n\n<p>c</p>")},
		// Issue 12369
		{"markdown", []byte("<div class=\"paragraph\">\n<p>foo</p>\n</div>"), []byte("<div class=\"paragraph\">\n<p>foo</p>\n</div>")},
		{"asciidoc", []byte("<div class=\"paragraph\">\n<p>foo</p>\n</div>"), []byte("foo")},
	}

	c := newTestContentSpec(nil)
	for i, test := range tests {
		output := c.TrimShortHTML(test.input, test.markup)
		if !bytes.Equal(test.output, output) {
			t.Errorf("Test %d failed. Expected %q got %q", i, test.output, output)
		}
	}
}

func BenchmarkTrimShortHTML(b *testing.B) {
	c := newTestContentSpec(nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.TrimShortHTML([]byte("<p>Simple paragraph</p>"), "markdown")
	}
}

func TestBytesToHTML(t *testing.T) {
	c := qt.New(t)
	c.Assert(helpers.BytesToHTML([]byte("dobedobedo")), qt.Equals, template.HTML("dobedobedo"))
}

var benchmarkTruncateString = strings.Repeat("This is a sentence about nothing.", 20)

func BenchmarkTestTruncateWordsToWholeSentence(b *testing.B) {
	c := newTestContentSpec(nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.TruncateWordsToWholeSentence(benchmarkTruncateString)
	}
}

func TestTruncateWordsToWholeSentence(t *testing.T) {
	type test struct {
		input, expected string
		max             int
		truncated       bool
	}
	data := []test{
		{"a b c", "a b c", 12, false},
		{"a b c", "a b c", 3, false},
		{"a", "a", 1, false},
		{"This is a sentence.", "This is a sentence.", 5, false},
		{"This is also a sentence!", "This is also a sentence!", 1, false},
		{"To be. Or not to be. That's the question.", "To be.", 1, true},
		{" \nThis is not a sentence\nAnd this is another", "This is not a sentence", 4, true},
		{"", "", 10, false},
		{"This... is a more difficult test?", "This... is a more difficult test?", 1, false},
	}
	for i, d := range data {
		cfg := config.New()
		cfg.Set("summaryLength", d.max)
		c := newTestContentSpec(cfg)
		output, truncated := c.TruncateWordsToWholeSentence(d.input)
		if d.expected != output {
			t.Errorf("Test %d failed. Expected %q got %q", i, d.expected, output)
		}

		if d.truncated != truncated {
			t.Errorf("Test %d failed. Expected truncated=%t got %t", i, d.truncated, truncated)
		}
	}
}

func TestTruncateWordsByRune(t *testing.T) {
	type test struct {
		input, expected string
		max             int
		truncated       bool
	}
	data := []test{
		{"", "", 1, false},
		{"a b c", "a b c", 12, false},
		{"a b c", "a b c", 3, false},
		{"a", "a", 1, false},
		{"Hello 中国", "", 0, true},
		{"这是中文，全中文。", "这是中文，", 5, true},
		{"Hello 中国", "Hello 中", 2, true},
		{"Hello 中国", "Hello 中国", 3, false},
		{"Hello中国 Good 好的", "Hello中国 Good 好", 9, true},
		{"This is a sentence.", "This is", 2, true},
		{"This is also a sentence!", "This", 1, true},
		{"To be. Or not to be. That's the question.", "To be. Or not", 4, true},
		{" \nThis is    not a sentence\n ", "This is not", 3, true},
	}
	for i, d := range data {
		cfg := config.New()
		cfg.Set("summaryLength", d.max)
		c := newTestContentSpec(cfg)
		output, truncated := c.TruncateWordsByRune(strings.Fields(d.input))
		if d.expected != output {
			t.Errorf("Test %d failed. Expected %q got %q", i, d.expected, output)
		}

		if d.truncated != truncated {
			t.Errorf("Test %d failed. Expected truncated=%t got %t", i, d.truncated, truncated)
		}
	}
}

func TestExtractTOCNormalContent(t *testing.T) {
	content := []byte("<nav>\n<ul>\nTOC<li><a href=\"#")

	actualTocLessContent, actualToc := helpers.ExtractTOC(content)
	expectedTocLess := []byte("TOC<li><a href=\"#")
	expectedToc := []byte("<nav id=\"TableOfContents\">\n<ul>\n")

	if !bytes.Equal(actualTocLessContent, expectedTocLess) {
		t.Errorf("Actual tocless (%s) did not equal expected (%s) tocless content", actualTocLessContent, expectedTocLess)
	}

	if !bytes.Equal(actualToc, expectedToc) {
		t.Errorf("Actual toc (%s) did not equal expected (%s) toc content", actualToc, expectedToc)
	}
}

func TestExtractTOCGreaterThanSeventy(t *testing.T) {
	content := []byte("<nav>\n<ul>\nTOC This is a very long content which will definitely be greater than seventy, I promise you that.<li><a href=\"#")

	actualTocLessContent, actualToc := helpers.ExtractTOC(content)
	// Because the start of Toc is greater than 70+startpoint of <li> content and empty TOC will be returned
	expectedToc := []byte("")

	if !bytes.Equal(actualTocLessContent, content) {
		t.Errorf("Actual tocless (%s) did not equal expected (%s) tocless content", actualTocLessContent, content)
	}

	if !bytes.Equal(actualToc, expectedToc) {
		t.Errorf("Actual toc (%s) did not equal expected (%s) toc content", actualToc, expectedToc)
	}
}

func TestExtractNoTOC(t *testing.T) {
	content := []byte("TOC")

	actualTocLessContent, actualToc := helpers.ExtractTOC(content)
	expectedToc := []byte("")

	if !bytes.Equal(actualTocLessContent, content) {
		t.Errorf("Actual tocless (%s) did not equal expected (%s) tocless content", actualTocLessContent, content)
	}

	if !bytes.Equal(actualToc, expectedToc) {
		t.Errorf("Actual toc (%s) did not equal expected (%s) toc content", actualToc, expectedToc)
	}
}

var totalWordsBenchmarkString = strings.Repeat("Hugo Rocks ", 200)

func TestTotalWords(t *testing.T) {
	for i, this := range []struct {
		s     string
		words int
	}{
		{"Two, Words!", 2},
		{"Word", 1},
		{"", 0},
		{"One, Two,      Three", 3},
		{totalWordsBenchmarkString, 400},
	} {
		actualWordCount := helpers.TotalWords(this.s)

		if actualWordCount != this.words {
			t.Errorf("[%d] Actual word count (%d) for test string (%s) did not match %d", i, actualWordCount, this.s, this.words)
		}
	}
}

func BenchmarkTotalWords(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wordCount := helpers.TotalWords(totalWordsBenchmarkString)
		if wordCount != 400 {
			b.Fatal("Wordcount error")
		}
	}
}
