package tplimpl_test

import (
	"path/filepath"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/strawberry-tools/strawberry/hugolib"
	"github.com/strawberry-tools/strawberry/tpl"
)

func TestPrintUnusedTemplates(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
baseURL = 'http://example.com/'
printUnusedTemplates=true
-- content/p1.md --
---
title: "P1"
---
{{< usedshortcode >}}
-- layouts/baseof.html --
{{ block "main" . }}{{ end }}
-- layouts/baseof.json --
{{ block "main" . }}{{ end }}
-- layouts/index.html --
{{ define "main" }}FOO{{ end }}
-- layouts/_default/single.json --
-- layouts/_default/single.html --
{{ define "main" }}MAIN{{ end }}
-- layouts/post/single.html --
{{ define "main" }}MAIN{{ end }}
-- layouts/partials/usedpartial.html --
-- layouts/partials/unusedpartial.html --
-- layouts/shortcodes/usedshortcode.html --
{{ partial "usedpartial.html" }}
-- layouts/shortcodes/unusedshortcode.html --

	`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
		},
	)
	b.Build()

	unused := b.H.Tmpl().(tpl.UnusedTemplatesProvider).UnusedTemplates()

	var names []string
	for _, tmpl := range unused {
		names = append(names, tmpl.Name())
	}

	b.Assert(names, qt.DeepEquals, []string{"_default/single.json", "baseof.json", "partials/unusedpartial.html", "post/single.html", "shortcodes/unusedshortcode.html"})
	b.Assert(unused[0].Filename(), qt.Equals, filepath.Join(b.Cfg.WorkingDir, "layouts/_default/single.json"))
}

// Verify that the new keywords in Go 1.18 is available.
func TestGo18Constructs(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
baseURL = 'http://example.com/'
disableKinds = ["section", "home", "rss", "taxonomy",  "term", "rss"]
-- content/p1.md --
---
title: "P1"
---
-- layouts/partials/counter.html --
{{ if .Scratch.Get "counter" }}{{ .Scratch.Add "counter" 1 }}{{ else }}{{ .Scratch.Set "counter" 1 }}{{ end }}{{ return true }}
-- layouts/_default/single.html --
continue:{{ range seq 5 }}{{ if eq . 2 }}{{continue}}{{ end }}{{ . }}{{ end }}:END:
break:{{ range seq 5 }}{{ if eq . 2 }}{{break}}{{ end }}{{ . }}{{ end }}:END:
continue2:{{ range seq 5 }}{{ if eq . 2 }}{{ continue }}{{ end }}{{ . }}{{ end }}:END:
break2:{{ range seq 5 }}{{ if eq . 2 }}{{ break }}{{ end }}{{ . }}{{ end }}:END:

counter1: {{ partial "counter.html" . }}/{{ .Scratch.Get "counter" }}
and1: {{ if (and false (partial "counter.html" .)) }}true{{ else }}false{{ end }}
or1: {{ if (or true (partial "counter.html" .)) }}true{{ else }}false{{ end }}
and2: {{ if (and true (partial "counter.html" .)) }}true{{ else }}false{{ end }}
or2: {{ if (or false (partial "counter.html" .)) }}true{{ else }}false{{ end }}


counter2: {{ .Scratch.Get "counter" }}


	`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
			NeedsOsFS:   true,
		},
	)
	b.Build()

	b.AssertFileContent("public/p1/index.html", `
continue:1345:END:
break:1:END:
continue2:1345:END:
break2:1:END:
counter1: true/1
and1: false
or1: true
and2: true
or2: true
counter2: 3
`)
}

// Issue 10495
func TestCommentsBeforeBlockDefinition(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
baseURL = 'http://example.com/'
-- content/s1/p1.md --
---
title: "S1P1"
---
-- content/s2/p1.md --
---
title: "S2P1"
---
-- content/s3/p1.md --
---
title: "S3P1"
---
-- layouts/_default/baseof.html --
{{ block "main" . }}{{ end }}
-- layouts/s1/single.html --
{{/* foo */}}
{{ define "main" }}{{ .Title }}{{ end }}
-- layouts/s2/single.html --
{{- /* foo */}}
{{ define "main" }}{{ .Title }}{{ end }}
-- layouts/s3/single.html --
{{- /* foo */ -}}
{{ define "main" }}{{ .Title }}{{ end }}
	`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	)
	b.Build()

	b.AssertFileContent("public/s1/p1/index.html", `S1P1`)
	b.AssertFileContent("public/s2/p1/index.html", `S2P1`)
	b.AssertFileContent("public/s3/p1/index.html", `S3P1`)
}

func TestGoTemplateBugs(t *testing.T) {
	t.Run("Issue 11112", func(t *testing.T) {
		t.Parallel()

		files := `
-- config.toml --
-- layouts/index.html --
{{ $m := dict "key" "value" }}
{{ $k := "" }}
{{ $v := "" }}
{{ range $k, $v = $m }}
{{ $k }} = {{ $v }}
{{ end }}
	`

		b := hugolib.NewIntegrationTestBuilder(
			hugolib.IntegrationTestConfig{
				T:           t,
				TxtarString: files,
			},
		)
		b.Build()

		b.AssertFileContent("public/index.html", `key = value`)
	})
}

func TestSecurityAllowActionJSTmpl(t *testing.T) {
	filesTemplate := `
-- config.toml --
SECURITYCONFIG
-- layouts/index.html --
<script>
var a = §§{{.Title }}§§;
</script>
	`

	files := strings.ReplaceAll(filesTemplate, "SECURITYCONFIG", "")

	b, err := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	).BuildE()

	// This used to fail, but not in >= Hugo 0.121.0.
	b.Assert(err, qt.IsNil)
}

func TestGoogleAnalyticsTemplate(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','section','rss','sitemap','taxonomy','term']
[privacy.googleAnalytics]
disable = false
respectDoNotTrack = true
[services.googleAnalytics]
id = 'G-0123456789'
-- layouts/index.html --
{{ template "_internal/google_analytics.html" . }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html",
		`<script async src="https://www.googletagmanager.com/gtag/js?id=G-0123456789"></script>`,
		`var dnt = (navigator.doNotTrack || window.doNotTrack || navigator.msDoNotTrack);`,
	)
}

func TestDisqusTemplate(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','section','rss','sitemap','taxonomy','term']
[services.disqus]
shortname = 'foo'
[privacy.disqus]
disable = false
-- layouts/index.html --
{{ template "_internal/disqus.html" . }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html",
		`s.src = '//' + "foo" + '.disqus.com/embed.js';`,
	)
}

func TestSitemap(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home','rss','section','taxonomy','term']
[sitemap]
disable = true
-- content/p1.md --
---
title: p1
sitemap:
  p1_disable: foo
---
-- content/p2.md --
---
title: p2

---
-- layouts/_default/single.html --
{{ .Title }}
`

	// Test A: Exclude all pages via site config.
	b := hugolib.Test(t, files)
	b.AssertFileContentExact("public/sitemap.xml",
		"<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\"?>\n<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\"\n  xmlns:xhtml=\"http://www.w3.org/1999/xhtml\">\n  \n</urlset>\n",
	)

	// Test B: Include all pages via site config.
	files_b := strings.ReplaceAll(files, "disable = true", "disable = false")
	b = hugolib.Test(t, files_b)
	b.AssertFileContentExact("public/sitemap.xml",
		"<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\"?>\n<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\"\n  xmlns:xhtml=\"http://www.w3.org/1999/xhtml\">\n  <url>\n    <loc>/p1/</loc>\n  </url><url>\n    <loc>/p2/</loc>\n  </url>\n</urlset>\n",
	)

	// Test C: Exclude all pages via site config, but include p1 via front matter.
	files_c := strings.ReplaceAll(files, "p1_disable: foo", "disable: false")
	b = hugolib.Test(t, files_c)
	b.AssertFileContentExact("public/sitemap.xml",
		"<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\"?>\n<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\"\n  xmlns:xhtml=\"http://www.w3.org/1999/xhtml\">\n  <url>\n    <loc>/p1/</loc>\n  </url>\n</urlset>\n",
	)

	// Test D:  Include all pages via site config, but exclude p1 via front matter.
	files_d := strings.ReplaceAll(files_b, "p1_disable: foo", "disable: true")
	b = hugolib.Test(t, files_d)
	b.AssertFileContentExact("public/sitemap.xml",
		"<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\"?>\n<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\"\n  xmlns:xhtml=\"http://www.w3.org/1999/xhtml\">\n  <url>\n    <loc>/p2/</loc>\n  </url>\n</urlset>\n",
	)
}
