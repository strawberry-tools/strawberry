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

// This file is autogenerated.

package page

import (
	"encoding/json"
	"github.com/strawberry-tools/strawberry/config"
	"time"
)

func MarshalPageToJSON(p Page) ([]byte, error) {
	date := p.Date()
	lastmod := p.Lastmod()
	publishDate := p.PublishDate()
	expiryDate := p.ExpiryDate()
	aliases := p.Aliases()
	bundleType := p.BundleType()
	description := p.Description()
	draft := p.Draft()
	isHome := p.IsHome()
	keywords := p.Keywords()
	kind := p.Kind()
	layout := p.Layout()
	linkTitle := p.LinkTitle()
	isNode := p.IsNode()
	isPage := p.IsPage()
	path := p.Path()
	slug := p.Slug()
	lang := p.Lang()
	isSection := p.IsSection()
	section := p.Section()
	sitemap := p.Sitemap()
	typ := p.Type()
	weight := p.Weight()

	s := struct {
		Date        time.Time
		Lastmod     time.Time
		PublishDate time.Time
		ExpiryDate  time.Time
		Aliases     []string
		BundleType  string
		Description string
		Draft       bool
		IsHome      bool
		Keywords    []string
		Kind        string
		Layout      string
		LinkTitle   string
		IsNode      bool
		IsPage      bool
		Path        string
		Slug        string
		Lang        string
		IsSection   bool
		Section     string
		Sitemap     config.SitemapConfig
		Type        string
		Weight      int
	}{
		Date:        date,
		Lastmod:     lastmod,
		PublishDate: publishDate,
		ExpiryDate:  expiryDate,
		Aliases:     aliases,
		BundleType:  bundleType,
		Description: description,
		Draft:       draft,
		IsHome:      isHome,
		Keywords:    keywords,
		Kind:        kind,
		Layout:      layout,
		LinkTitle:   linkTitle,
		IsNode:      isNode,
		IsPage:      isPage,
		Path:        path,
		Slug:        slug,
		Lang:        lang,
		IsSection:   isSection,
		Section:     section,
		Sitemap:     sitemap,
		Type:        typ,
		Weight:      weight,
	}

	return json.Marshal(&s)
}
