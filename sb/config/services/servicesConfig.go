// Copyright 2024 The Strawberry Tools team. All rights reserved.
// Copyright 2019 The Hugo Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package services

import (
	"github.com/mitchellh/mapstructure"
	"github.com/strawberry-tools/strawberry/config"
)

const (
	servicesConfigKey = "services"

	disqusShortnameKey = "disqusshortname"
	googleAnalyticsKey = "googleanalytics"
	rssLimitKey        = "rssLimit"
	jsonFeedLimitKey   = "jsonFeedLimit"
	jsonFeedFullKey    = "jsonFeedFull"
)

// Config is a privacy configuration for all the relevant services in Hugo.
type Config struct {
	Disqus          Disqus
	GoogleAnalytics GoogleAnalytics
	Instagram       Instagram
	Twitter         Twitter
	RSS             RSS
	JSONFeed        JSONFeed
}

// Disqus holds the functional configuration settings related to the Disqus template.
type Disqus struct {
	// A Shortname is the unique identifier assigned to a Disqus site.
	Shortname string
}

// GoogleAnalytics holds the functional configuration settings related to the Google Analytics template.
type GoogleAnalytics struct {
	// The GA tracking ID.
	ID string
}

// Instagram holds the functional configuration settings related to the Instagram shortcodes.
type Instagram struct {
	// The Simple variant of the Instagram is decorated with Bootstrap 4 card classes.
	// This means that if you use Bootstrap 4 or want to provide your own CSS, you want
	// to disable the inline CSS provided by Hugo.
	DisableInlineCSS bool

	// App or Client Access Token.
	// If you are using a Client Access Token, remember that you must combine it with your App ID
	// using a pipe symbol (<APPID>|<CLIENTTOKEN>) otherwise the request will fail.
	AccessToken string
}

// Twitter holds the functional configuration settings related to the Twitter shortcodes.
type Twitter struct {
	// The Simple variant of Twitter is decorated with a basic set of inline styles.
	// This means that if you want to provide your own CSS, you want
	// to disable the inline CSS provided by Hugo.
	DisableInlineCSS bool
}

// RSS holds the functional configuration settings related to the RSS feeds.
type RSS struct {
	// Limit the number of pages.
	Limit int
}

// JSONFeed holds the config settings for a JSON Feed: https://www.jsonfeed.org/
type JSONFeed struct {
	Limit int  // limit number of pages included in feed
	Full  bool // Whether or not to use the full content or just the summary
}

// DecodeConfig creates a services Config from a given Hugo configuration.
func DecodeConfig(cfg config.Provider) (c Config, err error) {
	m := cfg.GetStringMap(servicesConfigKey)

	err = mapstructure.WeakDecode(m, &c)

	// Keep backwards compatibility.
	if c.GoogleAnalytics.ID == "" {
		// Try the global config
		c.GoogleAnalytics.ID = cfg.GetString(googleAnalyticsKey)
	}
	if c.Disqus.Shortname == "" {
		c.Disqus.Shortname = cfg.GetString(disqusShortnameKey)
	}

	if c.RSS.Limit == 0 {
		c.RSS.Limit = cfg.GetInt(rssLimitKey)
	}

	if c.JSONFeed.Limit == 0 {
		c.JSONFeed.Limit = cfg.GetInt(jsonFeedLimitKey)
	}

	if c.JSONFeed.Full == false {
		c.JSONFeed.Full = cfg.GetBool(jsonFeedFullKey)
	}

	return
}
