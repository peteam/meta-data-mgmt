package service

import (
	"fmt"

	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/entity"
	"gopkg.in/couchbase/gocb.v1"
	"gopkg.in/couchbase/gocb.v1/cbft"
)

func (s *Service) SearchResourceByFields(pageSize int, pageNumber int, searchEntity *entity.SearchByFieldsBodyEntity) (gocb.SearchResults, error) {
	searchQuery := BuildSearchByFieldsQuery(pageSize, pageNumber, searchEntity)
	searchResult, err := s.repo.SearchResource(searchQuery)
	return searchResult, err
}

// BuildSearchByFieldsQuery build the query for repo doing search
func BuildSearchByFieldsQuery(pageSize int, pageNo int,
	searchEntity *entity.SearchByFieldsBodyEntity) *gocb.SearchQuery {
	fmt.Printf("\nPage size:%d \nPage number: %d\n", pageSize, pageNo)

	matchSearchFieldGroup := make([]*entity.ResourceSearchField, 0)
	notMatchSearchFieldGroup := make([]*entity.ResourceSearchField, 0)
	searchFieldArr := searchEntity.Input

	for _, searchField := range searchFieldArr {
		if searchField.ShoudMatch == false {
			notMatchSearchFieldGroup = append(notMatchSearchFieldGroup, searchField)
		} else {
			matchSearchFieldGroup = append(matchSearchFieldGroup, searchField)
		}
	}

	boolQuery := cbft.NewBooleanQuery()

	conjunctionQuery := cbft.NewConjunctionQuery(cbft.NewMatchQuery(searchEntity.URN).Field("urn"))
	isMatchAll := (searchEntity.Match == entity.SearchFieldAll)
	if isMatchAll == true {
		if len(matchSearchFieldGroup) != 0 {

			for _, matchSearchField := range matchSearchFieldGroup {
				conjunctionQuery.And(cbft.NewMatchQuery(matchSearchField.Term).Field(matchSearchField.Field))
			}

		}
		if len(notMatchSearchFieldGroup) != 0 {
			disjunctionQuery := cbft.NewDisjunctionQuery(cbft.NewMatchQuery(notMatchSearchFieldGroup[0].Term).Field(notMatchSearchFieldGroup[0].Field))

			for notMatchIdx, notMatchSearchField := range notMatchSearchFieldGroup {
				if notMatchIdx > 0 {
					disjunctionQuery.Or(cbft.NewMatchQuery(notMatchSearchField.Term).Field(notMatchSearchField.Field))
				}
			}
			boolQuery.MustNot(disjunctionQuery)
		}
	} else {
		if len(searchFieldArr) != 0 {
			disjunctionQuery := cbft.NewDisjunctionQuery(cbft.NewMatchQuery(searchFieldArr[0].Term).Field(searchFieldArr[0].Field))

			for matchAllIdx, matchAllField := range searchFieldArr {

				if matchAllIdx > 0 {
					disjunctionQuery.Or(cbft.NewMatchQuery(matchAllField.Term).Field(matchAllField.Field))
				}

			}

			boolQuery.Should(disjunctionQuery).ShouldMin(1)
		}
	}

	boolQuery.Must(conjunctionQuery)

	query := gocb.NewSearchQuery("all-field",
		boolQuery).Limit(pageSize).Skip(pageSize * (pageNo - 1))
	if searchEntity.Facet != nil && searchEntity.Facet.Field != "" {
		facetLabel := "Grouping"
		if searchEntity.Facet.Label != "" {
			facetLabel = searchEntity.Facet.Label
		}
		facetLimit := 5
		if searchEntity.Facet.Limit != 0 {
			facetLimit = searchEntity.Facet.Limit
		}
		query.AddFacet(facetLabel, cbft.NewTermFacet(searchEntity.Facet.Field, facetLimit))
	}
	return query
}
