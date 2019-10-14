package repository

import (
	entity "cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/entity"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockDbRepository is a mock of DbRepository interface
type MockDbRepository struct {
	ctrl     *gomock.Controller
	recorder *MockDbRepositoryMockRecorder
}

// MockDbRepositoryMockRecorder is the mock recorder for MockDbRepository
type MockDbRepositoryMockRecorder struct {
	mock *MockDbRepository
}

// NewMockDbRepository creates a new mock instance
func NewMockDbRepository(ctrl *gomock.Controller) *MockDbRepository {
	mock := &MockDbRepository{ctrl: ctrl}
	mock.recorder = &MockDbRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDbRepository) EXPECT() *MockDbRepositoryMockRecorder {
	return m.recorder
}

// Insert mocks base method
func (m *MockDbRepository) Insert(item *entity.Foo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", item)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert
func (mr *MockDbRepositoryMockRecorder) Insert(item interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockDbRepository)(nil).Insert), item)
}

// Upsert mocks base method
func (m *MockDbRepository) Upsert(item *entity.Foo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upsert", item)
	ret0, _ := ret[0].(error)
	return ret0
}

// Upsert indicates an expected call of Upsert
func (mr *MockDbRepositoryMockRecorder) Upsert(item interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upsert", reflect.TypeOf((*MockDbRepository)(nil).Upsert), item)
}

// Remove mocks base method
func (m *MockDbRepository) Remove(userId, itemId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", userId, itemId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove
func (mr *MockDbRepositoryMockRecorder) Remove(userId, itemId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockDbRepository)(nil).Remove), userId, itemId)
}

// List mocks base method
func (m *MockDbRepository) List() ([]*entity.ResourceType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List")
	ret0, _ := ret[0].([]*entity.ResourceType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockDbRepositoryMockRecorder) List() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockDbRepository)(nil).List))
}

// Count mocks base method
func (m *MockDbRepository) Count(userId string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", userId)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count
func (mr *MockDbRepositoryMockRecorder) Count(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockDbRepository)(nil).Count), userId)
}

// Retrieve mocks base method
func (m *MockDbRepository) Retrieve(resourceId, contentType string) (*entity.Content, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Retrieve", resourceId, contentType)
	ret0, _ := ret[0].(*entity.Content)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Retrieve indicates an expected call of Retrieve
func (mr *MockDbRepositoryMockRecorder) Retrieve(resourceId, contentType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Retrieve", reflect.TypeOf((*MockDbRepository)(nil).Retrieve), resourceId, contentType)
}

// RetrievewithOptional mocks base method
func (m *MockDbRepository) RetrievewithOptional(resourceId, contentType, providerName, catalogType, entityStatus, pageNumber, pageSize string) (*entity.Content, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrievewithOptional", resourceId, contentType, providerName, catalogType, entityStatus, pageNumber, pageSize)
	ret0, _ := ret[0].(*entity.Content)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrievewithOptional indicates an expected call of RetrievewithOptional
func (mr *MockDbRepositoryMockRecorder) RetrievewithOptional(resourceId, contentType, providerName, catalogType, entityStatus, pageNumber, pageSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrievewithOptional", reflect.TypeOf((*MockDbRepository)(nil).RetrievewithOptional), resourceId, contentType, providerName, catalogType, entityStatus, pageNumber, pageSize)
}

// MultiRetrievewithOptional mocks base method
func (m *MockDbRepository) MultiRetrievewithOptional(resourceId, contentType, providerName, catalogType, entityStatus, pageNumber, pageSize string) (*entity.MultiContent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MultiRetrievewithOptional", resourceId, contentType, providerName, catalogType, entityStatus, pageNumber, pageSize)
	ret0, _ := ret[0].(*entity.MultiContent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MultiRetrievewithOptional indicates an expected call of MultiRetrievewithOptional
func (mr *MockDbRepositoryMockRecorder) MultiRetrievewithOptional(resourceId, contentType, providerName, catalogType, entityStatus, pageNumber, pageSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MultiRetrievewithOptional", reflect.TypeOf((*MockDbRepository)(nil).MultiRetrievewithOptional), resourceId, contentType, providerName, catalogType, entityStatus, pageNumber, pageSize)
}

// PaginationCount mocks base method
func (m *MockDbRepository) PaginationCount(resourceId, contentType, providerName, catalogType, entityStatus, pageNumber, pageSize string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PaginationCount", resourceId, contentType, providerName, catalogType, entityStatus, pageNumber, pageSize)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PaginationCount indicates an expected call of PaginationCount
func (mr *MockDbRepositoryMockRecorder) PaginationCount(resourceId, contentType, providerName, catalogType, entityStatus, pageNumber, pageSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PaginationCount", reflect.TypeOf((*MockDbRepository)(nil).PaginationCount), resourceId, contentType, providerName, catalogType, entityStatus, pageNumber, pageSize)
}

// Readyz mocks base method
func (m *MockDbRepository) Readyz(item *entity.Health) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Readyz", item)
	ret0, _ := ret[0].(error)
	return ret0
}

// Readyz indicates an expected call of Readyz
func (mr *MockDbRepositoryMockRecorder) Readyz(item interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Readyz", reflect.TypeOf((*MockDbRepository)(nil).Readyz), item)
}

// Healthz mocks base method
func (m *MockDbRepository) Healthz() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Healthz")
	ret0, _ := ret[0].(error)
	return ret0
}

// Healthz indicates an expected call of Healthz
func (mr *MockDbRepositoryMockRecorder) Healthz() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Healthz", reflect.TypeOf((*MockDbRepository)(nil).Healthz))
}
