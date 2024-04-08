// Code generated by MockGen. DO NOT EDIT.
// Source: post_banner.go

// Package mock_postbanner is a generated GoMock package.
package mock_postbanner

import (
	reflect "reflect"

	banner "github.com/CyberPiess/banner_sevice/internal/domain/banner"
	gomock "github.com/golang/mock/gomock"
)

// MockpostBannerService is a mock of postBannerService interface.
type MockpostBannerService struct {
	ctrl     *gomock.Controller
	recorder *MockpostBannerServiceMockRecorder
}

// MockpostBannerServiceMockRecorder is the mock recorder for MockpostBannerService.
type MockpostBannerServiceMockRecorder struct {
	mock *MockpostBannerService
}

// NewMockpostBannerService creates a new mock instance.
func NewMockpostBannerService(ctrl *gomock.Controller) *MockpostBannerService {
	mock := &MockpostBannerService{ctrl: ctrl}
	mock.recorder = &MockpostBannerServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockpostBannerService) EXPECT() *MockpostBannerServiceMockRecorder {
	return m.recorder
}

// PostBanner mocks base method.
func (m *MockpostBannerService) PostBanner(newBanner banner.BannerEntity, user banner.User) (int64, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostBanner", newBanner, user)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// PostBanner indicates an expected call of PostBanner.
func (mr *MockpostBannerServiceMockRecorder) PostBanner(newBanner, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostBanner", reflect.TypeOf((*MockpostBannerService)(nil).PostBanner), newBanner, user)
}

// SearchAllBanners mocks base method.
func (m *MockpostBannerService) SearchAllBanners(bannerFilter banner.Filter, user banner.User) ([]banner.BannerEntity, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchAllBanners", bannerFilter, user)
	ret0, _ := ret[0].([]banner.BannerEntity)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SearchAllBanners indicates an expected call of SearchAllBanners.
func (mr *MockpostBannerServiceMockRecorder) SearchAllBanners(bannerFilter, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchAllBanners", reflect.TypeOf((*MockpostBannerService)(nil).SearchAllBanners), bannerFilter, user)
}

// SearchBanner mocks base method.
func (m *MockpostBannerService) SearchBanner(bannerFilter banner.Filter, user banner.User) (banner.BannerEntity, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchBanner", bannerFilter, user)
	ret0, _ := ret[0].(banner.BannerEntity)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SearchBanner indicates an expected call of SearchBanner.
func (mr *MockpostBannerServiceMockRecorder) SearchBanner(bannerFilter, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchBanner", reflect.TypeOf((*MockpostBannerService)(nil).SearchBanner), bannerFilter, user)
}