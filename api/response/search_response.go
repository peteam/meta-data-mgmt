package response

import (
	"gopkg.in/couchbase/gocb.v1"
)

type SearchResponse struct {
	Header             *Header             `json:"header"`
	SearchResponseData *SearchResponseData `json:"data,omitempty"`
}
type SearchResponseData struct {
	SearchResults []string                  `json:"ids"`
	SearchFacets  []*gocb.SearchResultFacet `json:"facets"`
}
