// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/krobus00/storage-service/internal/model (interfaces: ObjectRepository)

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

// MockObjectRepository is a mock of ObjectRepository interface.
type MockObjectRepository struct {
	ctrl     *gomock.Controller
	recorder *MockObjectRepositoryMockRecorder
}

// MockObjectRepositoryMockRecorder is the mock recorder for MockObjectRepository.
type MockObjectRepositoryMockRecorder struct {
	mock *MockObjectRepository
}

// NewMockObjectRepository creates a new mock instance.
func NewMockObjectRepository(ctrl *gomock.Controller) *MockObjectRepository {
	mock := &MockObjectRepository{ctrl: ctrl}
	mock.recorder = &MockObjectRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockObjectRepository) EXPECT() *MockObjectRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockObjectRepository) Create(arg0 context.Context, arg1 *model.ObjectPayload) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockObjectRepositoryMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockObjectRepository)(nil).Create), arg0, arg1)
}

// FindByID mocks base method.
func (m *MockObjectRepository) FindByID(arg0 context.Context, arg1 string) (*model.Object, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", arg0, arg1)
	ret0, _ := ret[0].(*model.Object)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockObjectRepositoryMockRecorder) FindByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockObjectRepository)(nil).FindByID), arg0, arg1)
}

// GeneratePresignedURL mocks base method.
func (m *MockObjectRepository) GeneratePresignedURL(arg0 context.Context, arg1 *model.Object) (*model.GetPresignedURLResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GeneratePresignedURL", arg0, arg1)
	ret0, _ := ret[0].(*model.GetPresignedURLResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GeneratePresignedURL indicates an expected call of GeneratePresignedURL.
func (mr *MockObjectRepositoryMockRecorder) GeneratePresignedURL(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GeneratePresignedURL", reflect.TypeOf((*MockObjectRepository)(nil).GeneratePresignedURL), arg0, arg1)
}

// InjectDB mocks base method.
func (m *MockObjectRepository) InjectDB(arg0 *gorm.DB) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InjectDB", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// InjectDB indicates an expected call of InjectDB.
func (mr *MockObjectRepositoryMockRecorder) InjectDB(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InjectDB", reflect.TypeOf((*MockObjectRepository)(nil).InjectDB), arg0)
}

// InjectRedisClient mocks base method.
func (m *MockObjectRepository) InjectRedisClient(arg0 *redis.Client) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InjectRedisClient", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// InjectRedisClient indicates an expected call of InjectRedisClient.
func (mr *MockObjectRepositoryMockRecorder) InjectRedisClient(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InjectRedisClient", reflect.TypeOf((*MockObjectRepository)(nil).InjectRedisClient), arg0)
}

// InjectS3Client mocks base method.
func (m *MockObjectRepository) InjectS3Client(arg0 model.S3Client) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InjectS3Client", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// InjectS3Client indicates an expected call of InjectS3Client.
func (mr *MockObjectRepositoryMockRecorder) InjectS3Client(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InjectS3Client", reflect.TypeOf((*MockObjectRepository)(nil).InjectS3Client), arg0)
}