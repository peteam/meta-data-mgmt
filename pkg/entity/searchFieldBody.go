package entity

// SearchFieldAny match by any field
const SearchFieldAny = "ANY"

// SearchFieldAll match all fields
const SearchFieldAll = "ALL"

// SearchByFieldsBodyEntity contains all necessary information for doing full text search
type SearchByFieldsBodyEntity struct {
	URN   string                 `json:"urn"`
	Match string                 `json:"match"`
	Facet *SearchFacet           `json:"facet"`
	Input []*ResourceSearchField `json:"input"`
}

// ResourceSearchField contains search term for couchbase FTS
type ResourceSearchField struct {
	Term       string `json:"term"`
	Field      string `json:"field"`
	ShoudMatch bool   `json:"should_match,omitempty"`
}

// SearchFacet - grouping search result by field (facet field) this group limit by facet_limit and label is facet_label
// Example: grouping research resource by "category" (facet field) - not more than 5 group returned (facet limit)
// grouping has label out put is: Categories (facet label)
type SearchFacet struct {
	Field string `json:"field"`
	Label string `json:"label"`
	Limit int    `json:"limit"`
}
