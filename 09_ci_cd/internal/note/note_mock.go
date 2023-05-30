// Code generated by MockGen. DO NOT EDIT.
// Source: note.go

// Package note is a generated GoMock package.
package note

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockNoteRepo is a mock of NoteRepo interface.
type MockNoteRepo struct {
	ctrl     *gomock.Controller
	recorder *MockNoteRepoMockRecorder
}

// MockNoteRepoMockRecorder is the mock recorder for MockNoteRepo.
type MockNoteRepoMockRecorder struct {
	mock *MockNoteRepo
}

// NewMockNoteRepo creates a new mock instance.
func NewMockNoteRepo(ctrl *gomock.Controller) *MockNoteRepo {
	mock := &MockNoteRepo{ctrl: ctrl}
	mock.recorder = &MockNoteRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNoteRepo) EXPECT() *MockNoteRepoMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockNoteRepo) Add(note *Note) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", note)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MockNoteRepoMockRecorder) Add(note interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockNoteRepo)(nil).Add), note)
}

// Delete mocks base method.
func (m *MockNoteRepo) Delete(id uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockNoteRepoMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockNoteRepo)(nil).Delete), id)
}

// GetAll mocks base method.
func (m *MockNoteRepo) GetAll() ([]*Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]*Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockNoteRepoMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockNoteRepo)(nil).GetAll))
}

// GetByID mocks base method.
func (m *MockNoteRepo) GetByID(id uint64) (*Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", id)
	ret0, _ := ret[0].(*Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockNoteRepoMockRecorder) GetByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockNoteRepo)(nil).GetByID), id)
}

// Update mocks base method.
func (m *MockNoteRepo) Update(newNote *Note) (*Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", newNote)
	ret0, _ := ret[0].(*Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockNoteRepoMockRecorder) Update(newNote interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockNoteRepo)(nil).Update), newNote)
}
