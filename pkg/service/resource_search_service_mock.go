package service

import (
	entity "cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/entity"
	"gopkg.in/couchbase/gocb.v1"
)

// ListResource mocks base method
func (m *MockMetaDataMgmtService) SearchResourceByFields(pageSize int, pageNumber int, searchEntity *entity.SearchByFieldsBodyEntity) (gocb.SearchResults, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchResourceByFields", pageSize, pageNumber, searchEntity)
	ret0, _ := ret[0].(gocb.SearchResults)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}
