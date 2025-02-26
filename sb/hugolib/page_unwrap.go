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
	"fmt"

	"github.com/strawberry-tools/strawberry/common/types"
	"github.com/strawberry-tools/strawberry/resources/page"
)

// Wraps a Page.
type pageWrapper interface {
	page() page.Page
}

// unwrapPage is used in equality checks and similar.
func unwrapPage(in any) (page.Page, error) {
	switch v := in.(type) {
	case *pageState:
		return v, nil
	case pageWrapper:
		return v.page(), nil
	case types.Unwrapper:
		return unwrapPage(v.Unwrapv())
	case page.Page:
		return v, nil
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("unwrapPage: %T not supported", in)
	}
}

func mustUnwrapPage(in any) page.Page {
	p, err := unwrapPage(in)
	if err != nil {
		panic(err)
	}

	return p
}
