// Copyright 2024 The Strawberry Tools team. All rights reserved.
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

package hugolib

import (
	"testing"
)

func TestEmbeddedShortcodes(t *testing.T) {

	t.Run("with theme", func(t *testing.T) {
		t.Parallel()

		files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["taxonomy", "term", "RSS", "sitemap", "robotsTXT", "page", "section"]
ignoreErrors = ["error-missing-instagram-accesstoken"]
[params]
foo = "bar"
-- content/_index.md --
---
title: "Home"
---

## Figure

{{< figure src="image.png" >}}

## Gist

{{< gist spf13 7896402 >}}

## Highlight

{{< highlight go >}}
package main
{{< /highlight >}}

## Instagram

{{< instagram BWNjjyYFxVx >}}

## Tweet

{{< tweet user="1626985695280603138" id="877500564405444608" >}}

## Vimeo

{{< vimeo 20097015 >}}

## YouTube

{{< youtube PArFPgHrNZM >}}

## Param

Foo: {{< param foo >}}

## Mastodon

{{< mastodon url="https://mastodon.social/@popey/101544533764122938" >}}

## QR code

{{< qrcode url="https://google.com" >}}

-- layouts/index.html --
Content: {{ .Content }}|
`
		b := Test(t, files)

		b.AssertFileContent("public/index.html", `
<figure>
https://gist.github.com/spf13/7896402.js
<span style="color:#a6e22e">main</span></span>
https://t.co/X94FmYDEZJ
https://www.youtube.com/embed/PArFPgHrNZM
Foo: bar
https://mastodon.social/@popey/101544533764122938
data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAQAAAAEAAQMAAABmvDolAAAABlBMVEX///8AAABVwtN+AAABZ0lEQVR42uyYwc3sIAyEHXHgSAkphdZSGqVQAkcOiHkaQ1bJ+1dbAGYOK23ynax4xrZsbW1t/RJUzVXxQJGQz/HkWgvo/HHNVV99EQKZT6I54GCZmrBOYKlOlipZBVzV70XEOCC+igQAeVFg9IXwgyihsA5fG2d14OOT9IdQfhjp2sCUApit8U2rA/1Ad0DzUIfQvIgp4loLONC1MTwq/4asjRGTGAP43mHE4vAH5Pc0aAQAumsOVXwNQOH7iGQPYHAOn9T/J/KZ4tNIlwDYGU0cE6CGEgrTID1j0Qhw9IM+SXdgLN6FehiIDaDLwb7gvMg5ifPiu1BGAP1gxpbEzgnAnzFpDWBOQToWMxblTDFFa8DzSDKPA/8ZqQ3gPpIwF4sEaKEk2gPurXmcSEbjpPd2sA7wWYIYB+/boDFgjMV6JnnHohFgHknmlsQdiT4p1oDbJ4E6xqR5JbnWAra2tizqXwAAAP//EnuWL6vnRrMAAAAASUVORK5CYII=



`)
	})
}
