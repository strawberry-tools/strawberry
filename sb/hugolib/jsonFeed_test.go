// Copyright 2024 The Strawberry Tools team. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package hugolib

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/strawberry-tools/strawberry/deps"

	"github.com/gopherlibs/jsonfeed/jsonfeed"
)

func TestJSONFeedOutput(t *testing.T) {

	t.Parallel()

	jsonFeedLimit := len(weightedSources) - 1

	cfg, fs := newTestCfg()
	cfg.Set("baseURL", "http://auth/bub/")
	cfg.Set("title", "JSON Feed Test")
	cfg.Set("jsonFeedLimit", jsonFeedLimit)
	cfg.Set("jsonFeedFull", false)
	th, configs := newTestHelperFromProvider(cfg, fs, t)

	jsonFeedURI := "index.feed.json"

	for _, src := range weightedSources {
		writeSource(t, fs, filepath.Join("content", "sect", src[0]), src[1])
	}

	buildSingleSite(t, deps.DepsCfg{Fs: fs, Configs: configs}, BuildCfg{})

	// Home feed
	th.assertFileContent(filepath.Join("public", jsonFeedURI), "https://jsonfeed.org/version/1.1", "JSON Feed Test", "http://auth/bub/")
	// Section feed
	th.assertFileContent(filepath.Join("public", "sect", jsonFeedURI), "https://jsonfeed.org/version/1.1", "Sects on JSON Feed Test", "http://auth/bub/")
	// Taxonomy feed
	th.assertFileContent(filepath.Join("public", "categories", "hugo", jsonFeedURI), "https://jsonfeed.org/version/1.1", "Hugo on JSON Feed Test", "http://auth/bub/")

	// JSON Feed Item Limit
	content := readWorkingDir(t, fs, filepath.Join("public", jsonFeedURI))
	c := strings.Count(content, "content_html")
	if c != jsonFeedLimit {
		t.Errorf("incorrect JSON Feed item count: expected %d, got %d", jsonFeedLimit, c)
	}

	// Encoded summary
	th.assertFileContent(filepath.Join("public", jsonFeedURI), "https://jsonfeed.org/version/1.1", "description", `A \u003cem\u003ecustom\u003c/em\u003e summary`)

	// Validate JSON Feed with external library
	homeFeed := readWorkingDir(th, th.Fs, filepath.Join("public", jsonFeedURI))
	jFeed, err := jsonfeed.Parse(strings.NewReader(homeFeed))
	if err != nil {
		t.Error("home JSON Feed cannot be parsed by JSONfeed library")
	}

	if errs := jFeed.Validate(); len(errs) > 0 {
		t.Error("the home JSON Feed is not valid according to schema")
	}
}
