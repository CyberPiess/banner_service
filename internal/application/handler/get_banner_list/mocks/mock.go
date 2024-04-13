// Code generated by MockGen. DO NOT EDIT.
// Source: get_banner.go

// Package mock_get_banner_list is a generated GoMock package.
package mock_get_banner_list

import (
	reflect "reflect"

	banner_service "github.com/CyberPiess/banner_service/internal/domain/banner"
	gomock "github.com/golang/mock/gomock"
	logrus "github.com/sirupsen/logrus"
)

// MockgetAllBannerService is a mock of getAllBannerService interface.
type MockgetAllBannerService struct {
	ctrl     *gomock.Controller
	recorder *MockgetAllBannerServiceMockRecorder
}

// MockgetAllBannerServiceMockRecorder is the mock recorder for MockgetAllBannerService.
type MockgetAllBannerServiceMockRecorder struct {
	mock *MockgetAllBannerService
}

// NewMockgetAllBannerService creates a new mock instance.
func NewMockgetAllBannerService(ctrl *gomock.Controller) *MockgetAllBannerService {
	mock := &MockgetAllBannerService{ctrl: ctrl}
	mock.recorder = &MockgetAllBannerServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockgetAllBannerService) EXPECT() *MockgetAllBannerServiceMockRecorder {
	return m.recorder
}

// DeleteBanner mocks base method.
func (m *MockgetAllBannerService) DeleteBanner(newDeleteBanner banner_service.BannerEntity, user banner_service.User) (bool, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBanner", newDeleteBanner, user)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DeleteBanner indicates an expected call of DeleteBanner.
func (mr *MockgetAllBannerServiceMockRecorder) DeleteBanner(newDeleteBanner, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBanner", reflect.TypeOf((*MockgetAllBannerService)(nil).DeleteBanner), newDeleteBanner, user)
}

// PostBanner mocks base method.
func (m *MockgetAllBannerService) PostBanner(newPostBanner banner_service.BannerEntity, user banner_service.User) (int, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostBanner", newPostBanner, user)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// PostBanner indicates an expected call of PostBanner.
func (mr *MockgetAllBannerServiceMockRecorder) PostBanner(newPostBanner, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostBanner", reflect.TypeOf((*MockgetAllBannerService)(nil).PostBanner), newPostBanner, user)
}

// PutBanner mocks base method.
func (m *MockgetAllBannerService) PutBanner(newPutBanner banner_service.BannerEntity, user banner_service.User) (bool, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutBanner", newPutBanner, user)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// PutBanner indicates an expected call of PutBanner.
func (mr *MockgetAllBannerServiceMockRecorder) PutBanner(newPutBanner, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutBanner", reflect.TypeOf((*MockgetAllBannerService)(nil).PutBanner), newPutBanner, user)
}

// SearchAllBanners mocks base method.
func (m *MockgetAllBannerService) SearchAllBanners(bannerFilter banner_service.GetAllFilter, user banner_service.User) ([]banner_service.BannerEntity, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchAllBanners", bannerFilter, user)
	ret0, _ := ret[0].([]banner_service.BannerEntity)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SearchAllBanners indicates an expected call of SearchAllBanners.
func (mr *MockgetAllBannerServiceMockRecorder) SearchAllBanners(bannerFilter, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchAllBanners", reflect.TypeOf((*MockgetAllBannerService)(nil).SearchAllBanners), bannerFilter, user)
}

// SearchBanner mocks base method.
func (m *MockgetAllBannerService) SearchBanner(bannerFilter banner_service.GetFilter, user banner_service.User) (banner_service.BannerEntity, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchBanner", bannerFilter, user)
	ret0, _ := ret[0].(banner_service.BannerEntity)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SearchBanner indicates an expected call of SearchBanner.
func (mr *MockgetAllBannerServiceMockRecorder) SearchBanner(bannerFilter, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchBanner", reflect.TypeOf((*MockgetAllBannerService)(nil).SearchBanner), bannerFilter, user)
}

// Mocklogger is a mock of logger interface.
type Mocklogger struct {
	ctrl     *gomock.Controller
	recorder *MockloggerMockRecorder
}

// MockloggerMockRecorder is the mock recorder for Mocklogger.
type MockloggerMockRecorder struct {
	mock *Mocklogger
}

// NewMocklogger creates a new mock instance.
func NewMocklogger(ctrl *gomock.Controller) *Mocklogger {
	mock := &Mocklogger{ctrl: ctrl}
	mock.recorder = &MockloggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mocklogger) EXPECT() *MockloggerMockRecorder {
	return m.recorder
}

// WithFields mocks base method.
func (m *Mocklogger) WithFields(fields logrus.Fields) *logrus.Entry {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithFields", fields)
	ret0, _ := ret[0].(*logrus.Entry)
	return ret0
}

// WithFields indicates an expected call of WithFields.
func (mr *MockloggerMockRecorder) WithFields(fields interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithFields", reflect.TypeOf((*Mocklogger)(nil).WithFields), fields)
}
