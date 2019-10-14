package repository

import (
	reflect "reflect"

	entity "cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/entity"
	gomock "github.com/golang/mock/gomock"
	"gopkg.in/couchbase/gocb.v1"
)

// InsertResource mocks base method
func (m *MockDbRepository) InsertResource(item *entity.Resource) (*entity.Resource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertResource", item)
	ret0, _ := ret[0].(*entity.Resource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert
func (mr *MockDbRepositoryMockRecorder) InsertResource(item interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertResource", reflect.TypeOf((*MockDbRepository)(nil).InsertResource), item)
}

// SearchResource mocks base method
func (m *MockDbRepository) SearchResource(query *gocb.SearchQuery) (gocb.SearchResults, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchResource", query)
	ret0, _ := ret[0].(gocb.SearchResults)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Insert
func (mr *MockDbRepositoryMockRecorder) SearchResource(item interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchResource", reflect.TypeOf((*MockDbRepository)(nil).InsertResource), item)
}

// DeleteResource mocks base method
func (m *MockDbRepository) DeleteResource(cbKey string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteResource", cbKey)

	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Insert
func (mr *MockDbRepositoryMockRecorder) DeleteResource(item interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteResource", reflect.TypeOf((*MockDbRepository)(nil).DeleteResource), item)
}
