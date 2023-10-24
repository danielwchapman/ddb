// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	ddb "github.com/danielwchapman/ddb"
	gomock "github.com/golang/mock/gomock"
)

// MockClientInterface is a mock of ClientInterface interface.
type MockClientInterface struct {
	ctrl     *gomock.Controller
	recorder *MockClientInterfaceMockRecorder
}

// MockClientInterfaceMockRecorder is the mock recorder for MockClientInterface.
type MockClientInterfaceMockRecorder struct {
	mock *MockClientInterface
}

// NewMockClientInterface creates a new mock instance.
func NewMockClientInterface(ctrl *gomock.Controller) *MockClientInterface {
	mock := &MockClientInterface{ctrl: ctrl}
	mock.recorder = &MockClientInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClientInterface) EXPECT() *MockClientInterfaceMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockClientInterface) Delete(ctx context.Context, pk, sk string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, pk, sk)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockClientInterfaceMockRecorder) Delete(ctx, pk, sk interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockClientInterface)(nil).Delete), ctx, pk, sk)
}

// Get mocks base method.
func (m *MockClientInterface) Get(ctx context.Context, pk, sk string, out any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, pk, sk, out)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockClientInterfaceMockRecorder) Get(ctx, pk, sk, out interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockClientInterface)(nil).Get), ctx, pk, sk, out)
}

// Put mocks base method.
func (m *MockClientInterface) Put(ctx context.Context, condition *string, row any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", ctx, condition, row)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put.
func (mr *MockClientInterfaceMockRecorder) Put(ctx, condition, row interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockClientInterface)(nil).Put), ctx, condition, row)
}

// TransactPuts mocks base method.
func (m *MockClientInterface) TransactPuts(ctx context.Context, token string, rows ...ddb.PutRow) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, token}
	for _, a := range rows {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "TransactPuts", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// TransactPuts indicates an expected call of TransactPuts.
func (mr *MockClientInterfaceMockRecorder) TransactPuts(ctx, token interface{}, rows ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, token}, rows...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransactPuts", reflect.TypeOf((*MockClientInterface)(nil).TransactPuts), varargs...)
}

// Update mocks base method.
func (m *MockClientInterface) Update(ctx context.Context, pk, sk string, opts ...ddb.Option) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, pk, sk}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Update", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockClientInterfaceMockRecorder) Update(ctx, pk, sk interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, pk, sk}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockClientInterface)(nil).Update), varargs...)
}

// MockDeleter is a mock of Deleter interface.
type MockDeleter struct {
	ctrl     *gomock.Controller
	recorder *MockDeleterMockRecorder
}

// MockDeleterMockRecorder is the mock recorder for MockDeleter.
type MockDeleterMockRecorder struct {
	mock *MockDeleter
}

// NewMockDeleter creates a new mock instance.
func NewMockDeleter(ctrl *gomock.Controller) *MockDeleter {
	mock := &MockDeleter{ctrl: ctrl}
	mock.recorder = &MockDeleterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeleter) EXPECT() *MockDeleterMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockDeleter) Delete(ctx context.Context, pk, sk string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, pk, sk)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockDeleterMockRecorder) Delete(ctx, pk, sk interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockDeleter)(nil).Delete), ctx, pk, sk)
}

// MockGetter is a mock of Getter interface.
type MockGetter struct {
	ctrl     *gomock.Controller
	recorder *MockGetterMockRecorder
}

// MockGetterMockRecorder is the mock recorder for MockGetter.
type MockGetterMockRecorder struct {
	mock *MockGetter
}

// NewMockGetter creates a new mock instance.
func NewMockGetter(ctrl *gomock.Controller) *MockGetter {
	mock := &MockGetter{ctrl: ctrl}
	mock.recorder = &MockGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGetter) EXPECT() *MockGetterMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockGetter) Get(ctx context.Context, pk, sk string, out any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, pk, sk, out)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockGetterMockRecorder) Get(ctx, pk, sk, out interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockGetter)(nil).Get), ctx, pk, sk, out)
}

// MockPutter is a mock of Putter interface.
type MockPutter struct {
	ctrl     *gomock.Controller
	recorder *MockPutterMockRecorder
}

// MockPutterMockRecorder is the mock recorder for MockPutter.
type MockPutterMockRecorder struct {
	mock *MockPutter
}

// NewMockPutter creates a new mock instance.
func NewMockPutter(ctrl *gomock.Controller) *MockPutter {
	mock := &MockPutter{ctrl: ctrl}
	mock.recorder = &MockPutterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPutter) EXPECT() *MockPutterMockRecorder {
	return m.recorder
}

// Put mocks base method.
func (m *MockPutter) Put(ctx context.Context, condition *string, row any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", ctx, condition, row)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put.
func (mr *MockPutterMockRecorder) Put(ctx, condition, row interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockPutter)(nil).Put), ctx, condition, row)
}

// MockTransactPutter is a mock of TransactPutter interface.
type MockTransactPutter struct {
	ctrl     *gomock.Controller
	recorder *MockTransactPutterMockRecorder
}

// MockTransactPutterMockRecorder is the mock recorder for MockTransactPutter.
type MockTransactPutterMockRecorder struct {
	mock *MockTransactPutter
}

// NewMockTransactPutter creates a new mock instance.
func NewMockTransactPutter(ctrl *gomock.Controller) *MockTransactPutter {
	mock := &MockTransactPutter{ctrl: ctrl}
	mock.recorder = &MockTransactPutterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactPutter) EXPECT() *MockTransactPutterMockRecorder {
	return m.recorder
}

// TransactPuts mocks base method.
func (m *MockTransactPutter) TransactPuts(ctx context.Context, token string, rows ...ddb.PutRow) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, token}
	for _, a := range rows {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "TransactPuts", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// TransactPuts indicates an expected call of TransactPuts.
func (mr *MockTransactPutterMockRecorder) TransactPuts(ctx, token interface{}, rows ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, token}, rows...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransactPuts", reflect.TypeOf((*MockTransactPutter)(nil).TransactPuts), varargs...)
}

// MockUpdater is a mock of Updater interface.
type MockUpdater struct {
	ctrl     *gomock.Controller
	recorder *MockUpdaterMockRecorder
}

// MockUpdaterMockRecorder is the mock recorder for MockUpdater.
type MockUpdaterMockRecorder struct {
	mock *MockUpdater
}

// NewMockUpdater creates a new mock instance.
func NewMockUpdater(ctrl *gomock.Controller) *MockUpdater {
	mock := &MockUpdater{ctrl: ctrl}
	mock.recorder = &MockUpdaterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUpdater) EXPECT() *MockUpdaterMockRecorder {
	return m.recorder
}

// Update mocks base method.
func (m *MockUpdater) Update(ctx context.Context, pk, sk string, opts ...ddb.Option) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, pk, sk}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Update", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockUpdaterMockRecorder) Update(ctx, pk, sk interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, pk, sk}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUpdater)(nil).Update), varargs...)
}
