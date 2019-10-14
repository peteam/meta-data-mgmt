package service

import (
	reflect "reflect"

	entity "cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockMetaDataMgmtService is a mock of MetaDataMgmtService interface
type MockMetaDataMgmtService struct {
	ctrl     *gomock.Controller
	recorder *MockMetaDataMgmtServiceMockRecorder
}

// MockMetaDataMgmtServiceMockRecorder is the mock recorder for MockMetaDataMgmtService
type MockMetaDataMgmtServiceMockRecorder struct {
	mock *MockMetaDataMgmtService
}

// NewMockMetaDataMgmtService creates a new mock instance
func NewMockMetaDataMgmtService(ctrl *gomock.Controller) *MockMetaDataMgmtService {
	mock := &MockMetaDataMgmtService{ctrl: ctrl}
	mock.recorder = &MockMetaDataMgmtServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMetaDataMgmtService) EXPECT() *MockMetaDataMgmtServiceMockRecorder {
	return m.recorder
}

// GetContent mocks base method
func (m *MockMetaDataMgmtService) GetContent(resourceId, contentType string) (*entity.Content, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContent", resourceId, contentType)
	ret0, _ := ret[0].(*entity.Content)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContent indicates an expected call of GetContent
func (mr *MockMetaDataMgmtServiceMockRecorder) GetContent(resourceId, contentType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContent", reflect.TypeOf((*MockMetaDataMgmtService)(nil).GetContent), resourceId, contentType)
}

// GetContentOptionalService mocks base method
func (m *MockMetaDataMgmtService) GetContentOptionalService(resourceId, contentType, providerName, catalogType, entityStatus, pageNumber, pageSize string) (*entity.Content, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContentOptionalService", resourceId, contentType, providerName, catalogType, entityStatus, pageNumber, pageSize)
	ret0, _ := ret[0].(*entity.Content)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContentOptionalService indicates an expected call of GetContentOptionalService
func (mr *MockMetaDataMgmtServiceMockRecorder) GetContentOptionalService(resourceId, contentType, providerName, catalogType, entityStatus, pageNumber, pageSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContentOptionalService", reflect.TypeOf((*MockMetaDataMgmtService)(nil).GetContentOptionalService), resourceId, contentType, providerName, catalogType, entityStatus, pageNumber, pageSize)
}

// GetMultiContentOptionalService mocks base method
func (m *MockMetaDataMgmtService) GetMultiContentOptionalService(resourceId, contentType, providerName, catalogType, entityStatus, pageNumber, pageSize string) (*entity.MultiContent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMultiContentOptionalService", resourceId, contentType, providerName, catalogType, entityStatus, pageNumber, pageSize)
	ret0, _ := ret[0].(*entity.MultiContent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMultiContentOptionalService indicates an expected call of GetMultiContentOptionalService
func (mr *MockMetaDataMgmtServiceMockRecorder) GetMultiContentOptionalService(resourceId, contentType, providerName, catalogType, entityStatus, pageNumber, pageSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMultiContentOptionalService", reflect.TypeOf((*MockMetaDataMgmtService)(nil).GetMultiContentOptionalService), resourceId, contentType, providerName, catalogType, entityStatus, pageNumber, pageSize)
}

// TotalContentCount mocks base method
func (m *MockMetaDataMgmtService) TotalContentCount(resourceId, contentType, providerName, catalogType, entityStatus, pageNumber, pageSize string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TotalContentCount", resourceId, contentType, providerName, catalogType, entityStatus, pageNumber, pageSize)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TotalContentCount indicates an expected call of TotalContentCount
func (mr *MockMetaDataMgmtServiceMockRecorder) TotalContentCount(resourceId, contentType, providerName, catalogType, entityStatus, pageNumber, pageSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TotalContentCount", reflect.TypeOf((*MockMetaDataMgmtService)(nil).TotalContentCount), resourceId, contentType, providerName, catalogType, entityStatus, pageNumber, pageSize)
}

// CountResource mocks base method
func (m *MockMetaDataMgmtService) CountResource(resourceType string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountResource", resourceType)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountResource indicates an expected call of CountResource
func (mr *MockMetaDataMgmtServiceMockRecorder) CountResource(resourceType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountResource", reflect.TypeOf((*MockMetaDataMgmtService)(nil).CountResource), resourceType)
}

// ListResource mocks base method
func (m *MockMetaDataMgmtService) ListResource() ([]*entity.ResourceType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListResource")
	ret0, _ := ret[0].([]*entity.ResourceType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListResource mocks base method
func (m *MockMetaDataMgmtService) AddResource(resource *entity.Resource) (*entity.Resource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddResource", resource)
	ret0, _ := ret[0].(*entity.Resource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListResource indicates an expected call of ListResource
func (mr *MockMetaDataMgmtServiceMockRecorder) ListResource() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListResource", reflect.TypeOf((*MockMetaDataMgmtService)(nil).ListResource))
}

// ListResource indicates an expected call of ListResource
func (mr *MockMetaDataMgmtServiceMockRecorder) AddResource() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddResource", reflect.TypeOf((*MockMetaDataMgmtService)(nil).AddResource))
}

// ListResource mocks base method
func (m *MockMetaDataMgmtService) DeleteResource(cbKey string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteResource", cbKey)

	ret0, _ := ret[0].(error)
	return ret0
}

// ListResource indicates an expected call of ListResource
func (mr *MockMetaDataMgmtServiceMockRecorder) DeleteResource() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteResource", reflect.TypeOf((*MockMetaDataMgmtService)(nil).DeleteResource))
}

// Readyz mocks base method
func (m *MockMetaDataMgmtService) Readyz() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Readyz")
	ret0, _ := ret[0].(error)
	return ret0
}

// Readyz indicates an expected call of Readyz
func (mr *MockMetaDataMgmtServiceMockRecorder) Readyz() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Readyz", reflect.TypeOf((*MockMetaDataMgmtService)(nil).Readyz))
}

// Healthz mocks base method
func (m *MockMetaDataMgmtService) Healthz() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Healthz")
	ret0, _ := ret[0].(error)
	return ret0
}

// Healthz indicates an expected call of Healthz
func (mr *MockMetaDataMgmtServiceMockRecorder) Healthz() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Healthz", reflect.TypeOf((*MockMetaDataMgmtService)(nil).Healthz))
}
