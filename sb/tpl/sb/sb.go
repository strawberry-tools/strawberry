// Copyright 2024 The Strawberry Tools team. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package sb

import (
	"encoding/base64"
	"html/template"

	qrcode "github.com/skip2/go-qrcode"
	"github.com/strawberry-tools/strawberry/deps"
)

func New(d *deps.Deps) *Namespace {
	return &Namespace{deps: d}
}

type Namespace struct {
	deps *deps.Deps
}

func (ns *Namespace) QRCoder(url any) (template.HTML, error) {

	var png []byte

	png, err := qrcode.Encode(url.(string), qrcode.Medium, 256)
	if err != nil {
		return "", err
	}

	encodedImage := base64.StdEncoding.EncodeToString(png)

	fullImage := "data:image/png;base64," + encodedImage

	return template.HTML(fullImage), nil
}
