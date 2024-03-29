// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/krobus00/storage-service/internal/model (interfaces: ObjectWhitelistTypeRepository)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	redis "github.com/go-redis/redis/v8"
	gomock "github.com/golang/mock/gomock"
	model "github.com/krobus00/storage-service/internal/model"
	gorm "gorm.io/gorm"
)

// MockObjectWhitelistTypeRepository is a mock of ObjectWhitelistTypeRepository interface.
type MockObjectWhitelistTypeRepository struct {
	ctrl     *gomock.Controller
	recorder *MockObjectWhitelistTypeRepositoryMockRecorder
}

// MockObjectWhitelistTypeRepositoryMockRecorder is the mock recorder for MockObjectWhitelistTypeRepository.
type MockObjectWhitelistTypeRepositoryMockRecorder struct {
	mock *MockObjectWhitelistTypeRepository
}

// NewMockObjectWhitelistTypeRepository creates a new mock instance.
func NewMockObjectWhitelistTypeRepository(ctrl *gomock.Controller) *MockObjectWhitelistTypeRepository {
	mock := &MockObjectWhitelistTypeRepository{ctrl: ctrl}
	mock.recorder = &MockObjectWhitelistTypeRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockObjectWhitelistTypeRepository) EXPECT() *MockObjectWhitelistTypeRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockObjectWhitelistTypeRepository) Create(arg0 context.Context, arg1 *model.ObjectWhitelistType) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockObjectWhitelistTypeRepositoryMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockObjectWhitelistTypeRepository)(nil).Create), arg0, arg1)
}

// DeleteByTypeIDAndExt mocks base method.
func (m *MockObjectWhitelistTypeRepository) DeleteByTypeIDAndExt(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByTypeIDAndExt", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByTypeIDAndExt indicates an expected call of DeleteByTypeIDAndExt.
func (mr *MockObjectWhitelistTypeRepositoryMockRecorder) DeleteByTypeIDAndExt(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByTypeIDAndExt", reflect.TypeOf((*MockObjectWhitelistTypeRepository)(nil).DeleteByTypeIDAndExt), arg0, arg1, arg2)
}

// FindByTypeIDAndExt mocks base method.
func (m *MockObjectWhitelistTypeRepository) FindByTypeIDAndExt(arg0 context.Context, arg1, arg2 string) (*model.ObjectWhitelistType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByTypeIDAndExt", arg0, arg1, arg2)
	ret0, _ := ret[0].(*model.ObjectWhitelistType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByTypeIDAndExt indicates an expected call of FindByTypeIDAndExt.
func (mr *MockObjectWhitelistTypeRepositoryMockRecorder) FindByTypeIDAndExt(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByTypeIDAndExt", reflect.TypeOf((*MockObjectWhitelistTypeRepository)(nil).FindByTypeIDAndExt), arg0, arg1, arg2)
}

// InjectDB mocks base method.
func (m *MockObjectWhitelistTypeRepository) InjectDB(arg0 *gorm.DB) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InjectDB", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// InjectDB indicates an expected call of InjectDB.
func (mr *MockObjectWhitelistTypeRepositoryMockRecorder) InjectDB(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InjectDB", reflect.TypeOf((*MockObjectWhitelistTypeRepository)(nil).InjectDB), arg0)
}

// InjectRedisClient mocks base method.
func (m *MockObjectWhitelistTypeRepository) InjectRedisClient(arg0 *redis.Client) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InjectRedisClient", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// InjectRedisClient indicates an expected call of InjectRedisClient.
func (mr *MockObjectWhitelistTypeRepositoryMockRecorder) InjectRedisClient(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InjectRedisClient", reflect.TypeOf((*MockObjectWhitelistTypeRepository)(nil).InjectRedisClient), arg0)
}
