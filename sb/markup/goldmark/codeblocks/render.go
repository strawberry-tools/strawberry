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

package codeblocks

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/strawberry-tools/strawberry/common/herrors"
	htext "github.com/strawberry-tools/strawberry/common/text"
	"github.com/strawberry-tools/strawberry/markup/converter/hooks"
	"github.com/strawberry-tools/strawberry/markup/goldmark/internal/render"
	"github.com/strawberry-tools/strawberry/markup/highlight/chromalexers"
	"github.com/strawberry-tools/strawberry/markup/internal/attributes"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type (
	codeBlocksExtension struct{}
	htmlRenderer        struct{}
)

func New() goldmark.Extender {
	return &codeBlocksExtension{}
}

func (e *codeBlocksExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(&Transformer{}, 100),
		),
	)
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(newHTMLRenderer(), 100),
	))
}

func newHTMLRenderer() renderer.NodeRenderer {
	r := &htmlRenderer{}
	return r
}

func (r *htmlRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindCodeBlock, r.renderCodeBlock)
}

func (r *htmlRenderer) renderCodeBlock(w util.BufWriter, src []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	ctx := w.(*render.Context)

	if entering {
		return ast.WalkContinue, nil
	}

	n := node.(*codeBlock)
	lang := getLang(n.b, src)
	renderer := ctx.RenderContext().GetRenderer(hooks.CodeBlockRendererType, lang)
	if renderer == nil {
		return ast.WalkStop, fmt.Errorf("no code renderer found for %q", lang)
	}

	ordinal := n.ordinal

	var buff bytes.Buffer

	l := n.b.Lines().Len()
	for i := 0; i < l; i++ {
		line := n.b.Lines().At(i)
		buff.Write(line.Value(src))
	}

	s := htext.Chomp(buff.String())

	var info []byte
	if n.b.Info != nil {
		info = n.b.Info.Segment.Value(src)
	}

	attrtp := attributes.AttributesOwnerCodeBlockCustom
	if isd, ok := renderer.(hooks.IsDefaultCodeBlockRendererProvider); (ok && isd.IsDefaultCodeBlockRenderer()) || chromalexers.Get(lang) != nil {
		// We say that this is a Chroma code block if it's the default code block renderer
		// or if the language is supported by Chroma.
		attrtp = attributes.AttributesOwnerCodeBlockChroma
	}

	// IsDefaultCodeBlockRendererProvider
	attrs, attrStr, err := getAttributes(n.b, info)
	if err != nil {
		return ast.WalkStop, &herrors.TextSegmentError{Err: err, Segment: attrStr}
	}
	cbctx := &codeBlockContext{
		page:             ctx.DocumentContext().Document,
		pageInner:        r.getPageInner(ctx),
		lang:             lang,
		code:             s,
		ordinal:          ordinal,
		AttributesHolder: attributes.New(attrs, attrtp),
	}

	cbctx.createPos = func() htext.Position {
		if resolver, ok := renderer.(hooks.ElementPositionResolver); ok {
			return resolver.ResolvePosition(cbctx)
		}
		return htext.Position{
			Filename:     ctx.DocumentContext().Filename,
			LineNumber:   1,
			ColumnNumber: 1,
		}
	}

	cr := renderer.(hooks.CodeBlockRenderer)

	err = cr.RenderCodeblock(
		ctx.RenderContext().Ctx,
		w,
		cbctx,
	)
	if err != nil {
		return ast.WalkContinue, herrors.NewFileErrorFromPos(err, cbctx.createPos())
	}

	return ast.WalkContinue, nil
}

func (r *htmlRenderer) getPageInner(rctx *render.Context) any {
	pid := rctx.PeekPid()
	if pid > 0 {
		if lookup := rctx.DocumentContext().DocumentLookup; lookup != nil {
			if v := rctx.DocumentContext().DocumentLookup(pid); v != nil {
				return v
			}
		}
	}
	return rctx.DocumentContext().Document
}

type codeBlockContext struct {
	page      any
	pageInner any
	lang      string
	code      string
	ordinal   int

	// This is only used in error situations and is expensive to create,
	// to delay creation until needed.
	pos       htext.Position
	posInit   sync.Once
	createPos func() htext.Position

	*attributes.AttributesHolder
}

func (c *codeBlockContext) Page() any {
	return c.page
}

func (c *codeBlockContext) PageInner() any {
	return c.pageInner
}

func (c *codeBlockContext) Type() string {
	return c.lang
}

func (c *codeBlockContext) Inner() string {
	return c.code
}

func (c *codeBlockContext) Ordinal() int {
	return c.ordinal
}

func (c *codeBlockContext) Position() htext.Position {
	c.posInit.Do(func() {
		c.pos = c.createPos()
	})
	return c.pos
}

func getLang(node *ast.FencedCodeBlock, src []byte) string {
	langWithAttributes := string(node.Language(src))
	lang, _, _ := strings.Cut(langWithAttributes, "{")
	return lang
}

func getAttributes(node *ast.FencedCodeBlock, infostr []byte) ([]ast.Attribute, string, error) {
	if node.Attributes() != nil {
		return node.Attributes(), "", nil
	}
	if infostr != nil {
		attrStartIdx := -1
		attrEndIdx := -1

		for idx, char := range infostr {
			if attrEndIdx == -1 && char == '{' {
				attrStartIdx = idx
			}
			if attrStartIdx != -1 && char == '}' {
				attrEndIdx = idx
				break
			}
		}

		if attrStartIdx != -1 && attrEndIdx != -1 {
			n := ast.NewTextBlock() // dummy node for storing attributes
			attrStr := infostr[attrStartIdx : attrEndIdx+1]
			if attrs, hasAttr := parser.ParseAttributes(text.NewReader(attrStr)); hasAttr {
				for _, attr := range attrs {
					n.SetAttribute(attr.Name, attr.Value)
				}
				return n.Attributes(), "", nil
			} else {
				return nil, string(attrStr), errors.New("failed to parse Markdown attributes; you may need to quote the values")
			}
		}
	}
	return nil, "", nil
}
