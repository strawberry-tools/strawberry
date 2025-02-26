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

package page

import (
	"fmt"
	"sort"
	"strings"

	"github.com/strawberry-tools/strawberry/compare"
	"github.com/strawberry-tools/strawberry/langs"
)

// The TaxonomyList is a list of all taxonomies and their values
// e.g. List['tags'] => TagTaxonomy (from above)
type TaxonomyList map[string]Taxonomy

func (tl TaxonomyList) String() string {
	return fmt.Sprintf("TaxonomyList(%d)", len(tl))
}

// A Taxonomy is a map of keywords to a list of pages.
// For example
//
//	TagTaxonomy['technology'] = WeightedPages
//	TagTaxonomy['go']  =  WeightedPages
type Taxonomy map[string]WeightedPages

// OrderedTaxonomy is another representation of an Taxonomy using an array rather than a map.
// Important because you can't order a map.
type OrderedTaxonomy []OrderedTaxonomyEntry

// getOneOPage returns one page in the taxonomy,
// nil if there is none.
func (t OrderedTaxonomy) getOneOPage() Page {
	if len(t) == 0 {
		return nil
	}
	return t[0].Pages()[0]
}

// OrderedTaxonomyEntry is similar to an element of a Taxonomy, but with the key embedded (as name)
// e.g:  {Name: Technology, WeightedPages: TaxonomyPages}
type OrderedTaxonomyEntry struct {
	Name string
	WeightedPages
}

// Get the weighted pages for the given key.
func (i Taxonomy) Get(key string) WeightedPages {
	return i[strings.ToLower(key)]
}

// Count the weighted pages for the given key.
func (i Taxonomy) Count(key string) int { return len(i[strings.ToLower(key)]) }

// TaxonomyArray returns an ordered taxonomy with a non defined order.
func (i Taxonomy) TaxonomyArray() OrderedTaxonomy {
	ies := make([]OrderedTaxonomyEntry, len(i))
	count := 0
	for k, v := range i {
		ies[count] = OrderedTaxonomyEntry{Name: k, WeightedPages: v}
		count++
	}
	return ies
}

// Alphabetical returns an ordered taxonomy sorted by key name.
func (i Taxonomy) Alphabetical() OrderedTaxonomy {
	ia := i.TaxonomyArray()
	p := ia.getOneOPage()
	if p == nil {
		return ia
	}
	currentSite := p.Site().Current()
	coll := langs.GetCollator1(currentSite.Language())
	coll.Lock()
	defer coll.Unlock()
	name := func(i1, i2 *OrderedTaxonomyEntry) bool {
		return coll.CompareStrings(i1.Name, i2.Name) < 0
	}
	oiBy(name).Sort(ia)
	return ia
}

// ByCount returns an ordered taxonomy sorted by # of pages per key.
// If taxonomies have the same # of pages, sort them alphabetical
func (i Taxonomy) ByCount() OrderedTaxonomy {
	count := func(i1, i2 *OrderedTaxonomyEntry) bool {
		li1 := len(i1.WeightedPages)
		li2 := len(i2.WeightedPages)

		if li1 == li2 {
			return compare.LessStrings(i1.Name, i2.Name)
		}
		return li1 > li2
	}

	ia := i.TaxonomyArray()
	oiBy(count).Sort(ia)
	return ia
}

// Page returns the taxonomy page or nil if the taxonomy has no terms.
func (i Taxonomy) Page() Page {
	for _, v := range i {
		return v.Page().Parent()
	}
	return nil
}

// Pages returns the Pages for this taxonomy.
func (ie OrderedTaxonomyEntry) Pages() Pages {
	return ie.WeightedPages.Pages()
}

// Count returns the count the pages in this taxonomy.
func (ie OrderedTaxonomyEntry) Count() int {
	return len(ie.WeightedPages)
}

// Term returns the name given to this taxonomy.
func (ie OrderedTaxonomyEntry) Term() string {
	return ie.Name
}

// Reverse reverses the order of the entries in this taxonomy.
func (t OrderedTaxonomy) Reverse() OrderedTaxonomy {
	for i, j := 0, len(t)-1; i < j; i, j = i+1, j-1 {
		t[i], t[j] = t[j], t[i]
	}

	return t
}

// A type to implement the sort interface for TaxonomyEntries.
type orderedTaxonomySorter struct {
	taxonomy OrderedTaxonomy
	by       oiBy
}

// Closure used in the Sort.Less method.
type oiBy func(i1, i2 *OrderedTaxonomyEntry) bool

func (by oiBy) Sort(taxonomy OrderedTaxonomy) {
	ps := &orderedTaxonomySorter{
		taxonomy: taxonomy,
		by:       by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Stable(ps)
}

// Len is part of sort.Interface.
func (s *orderedTaxonomySorter) Len() int {
	return len(s.taxonomy)
}

// Swap is part of sort.Interface.
func (s *orderedTaxonomySorter) Swap(i, j int) {
	s.taxonomy[i], s.taxonomy[j] = s.taxonomy[j], s.taxonomy[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *orderedTaxonomySorter) Less(i, j int) bool {
	return s.by(&s.taxonomy[i], &s.taxonomy[j])
}
