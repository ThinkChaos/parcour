// Code generated by MockGen. DO NOT EDIT.
// Source: jobgroup.go

// Package jobgroup is a generated GoMock package.
package jobgroup

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockJobGroup is a mock of JobGroup interface.
type MockJobGroup struct {
	ctrl     *gomock.Controller
	recorder *MockJobGroupMockRecorder
}

// MockJobGroupMockRecorder is the mock recorder for MockJobGroup.
type MockJobGroupMockRecorder struct {
	mock *MockJobGroup
}

// NewMockJobGroup creates a new mock instance.
func NewMockJobGroup(ctrl *gomock.Controller) *MockJobGroup {
	mock := &MockJobGroup{ctrl: ctrl}
	mock.recorder = &MockJobGroupMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJobGroup) EXPECT() *MockJobGroupMockRecorder {
	return m.recorder
}

// Cancel mocks base method.
func (m *MockJobGroup) Cancel() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Cancel")
}

// Cancel indicates an expected call of Cancel.
func (mr *MockJobGroupMockRecorder) Cancel() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cancel", reflect.TypeOf((*MockJobGroup)(nil).Cancel))
}

// Close mocks base method.
func (m *MockJobGroup) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockJobGroupMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockJobGroup)(nil).Close))
}

// Ctx mocks base method.
func (m *MockJobGroup) Ctx() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ctx")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Ctx indicates an expected call of Ctx.
func (mr *MockJobGroupMockRecorder) Ctx() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ctx", reflect.TypeOf((*MockJobGroup)(nil).Ctx))
}

// Go mocks base method.
func (m *MockJobGroup) Go(job Job) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Go", job)
}

// Go indicates an expected call of Go.
func (mr *MockJobGroupMockRecorder) Go(job interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Go", reflect.TypeOf((*MockJobGroup)(nil).Go), job)
}

// Wait mocks base method.
func (m *MockJobGroup) Wait() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Wait")
	ret0, _ := ret[0].(error)
	return ret0
}

// Wait indicates an expected call of Wait.
func (mr *MockJobGroupMockRecorder) Wait() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Wait", reflect.TypeOf((*MockJobGroup)(nil).Wait))
}

// WaitCtx mocks base method.
func (m *MockJobGroup) WaitCtx(ctx context.Context) (error, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitCtx", ctx)
	ret0, _ := ret[0].(error)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// WaitCtx indicates an expected call of WaitCtx.
func (mr *MockJobGroupMockRecorder) WaitCtx(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitCtx", reflect.TypeOf((*MockJobGroup)(nil).WaitCtx), ctx)
}

// private mocks base method.
func (m *MockJobGroup) private() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "private")
}

// private indicates an expected call of private.
func (mr *MockJobGroupMockRecorder) private() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "private", reflect.TypeOf((*MockJobGroup)(nil).private))
}

// MockjobGroup is a mock of jobGroup interface.
type MockjobGroup struct {
	ctrl     *gomock.Controller
	recorder *MockjobGroupMockRecorder
}

// MockjobGroupMockRecorder is the mock recorder for MockjobGroup.
type MockjobGroupMockRecorder struct {
	mock *MockjobGroup
}

// NewMockjobGroup creates a new mock instance.
func NewMockjobGroup(ctrl *gomock.Controller) *MockjobGroup {
	mock := &MockjobGroup{ctrl: ctrl}
	mock.recorder = &MockjobGroupMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockjobGroup) EXPECT() *MockjobGroupMockRecorder {
	return m.recorder
}

// Cancel mocks base method.
func (m *MockjobGroup) Cancel() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Cancel")
}

// Cancel indicates an expected call of Cancel.
func (mr *MockjobGroupMockRecorder) Cancel() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cancel", reflect.TypeOf((*MockjobGroup)(nil).Cancel))
}

// Close mocks base method.
func (m *MockjobGroup) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockjobGroupMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockjobGroup)(nil).Close))
}

// Ctx mocks base method.
func (m *MockjobGroup) Ctx() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ctx")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Ctx indicates an expected call of Ctx.
func (mr *MockjobGroupMockRecorder) Ctx() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ctx", reflect.TypeOf((*MockjobGroup)(nil).Ctx))
}

// Go mocks base method.
func (m *MockjobGroup) Go(job Job) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Go", job)
}

// Go indicates an expected call of Go.
func (mr *MockjobGroupMockRecorder) Go(job interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Go", reflect.TypeOf((*MockjobGroup)(nil).Go), job)
}

// Wait mocks base method.
func (m *MockjobGroup) Wait() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Wait")
	ret0, _ := ret[0].(error)
	return ret0
}

// Wait indicates an expected call of Wait.
func (mr *MockjobGroupMockRecorder) Wait() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Wait", reflect.TypeOf((*MockjobGroup)(nil).Wait))
}

// WaitCtx mocks base method.
func (m *MockjobGroup) WaitCtx(ctx context.Context) (error, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitCtx", ctx)
	ret0, _ := ret[0].(error)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// WaitCtx indicates an expected call of WaitCtx.
func (mr *MockjobGroupMockRecorder) WaitCtx(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitCtx", reflect.TypeOf((*MockjobGroup)(nil).WaitCtx), ctx)
}

// init mocks base method.
func (m *MockjobGroup) init(arg0 context.Context, arg1 JobGroup) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "init", arg0, arg1)
}

// init indicates an expected call of init.
func (mr *MockjobGroupMockRecorder) init(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "init", reflect.TypeOf((*MockjobGroup)(nil).init), arg0, arg1)
}

// launch mocks base method.
func (m *MockjobGroup) launch(arg0 *boundJob) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "launch", arg0)
}

// launch indicates an expected call of launch.
func (mr *MockjobGroupMockRecorder) launch(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "launch", reflect.TypeOf((*MockjobGroup)(nil).launch), arg0)
}

// private mocks base method.
func (m *MockjobGroup) private() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "private")
}

// private indicates an expected call of private.
func (mr *MockjobGroupMockRecorder) private() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "private", reflect.TypeOf((*MockjobGroup)(nil).private))
}

// saveErr mocks base method.
func (m *MockjobGroup) saveErr(arg0 error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "saveErr", arg0)
}

// saveErr indicates an expected call of saveErr.
func (mr *MockjobGroupMockRecorder) saveErr(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "saveErr", reflect.TypeOf((*MockjobGroup)(nil).saveErr), arg0)
}

// savePanic mocks base method.
func (m *MockjobGroup) savePanic(arg0 any) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "savePanic", arg0)
}

// savePanic indicates an expected call of savePanic.
func (mr *MockjobGroupMockRecorder) savePanic(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "savePanic", reflect.TypeOf((*MockjobGroup)(nil).savePanic), arg0)
}
