package adminbanner

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	appBanner "github.com/CyberPiess/banner_sevice/internal/application/handler/get_banner/mocks"
	"github.com/CyberPiess/banner_sevice/internal/domain/banner"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type args struct {
	w     *httptest.ResponseRecorder
	r     *http.Request
	token string
}

func TestGetAllBanners(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBannerService := appBanner.NewMockbannerService(ctrl)

	bannerHandler := NewBannerHandler(mockBannerService)

	testBannerEntity := banner.BannerEntity{ID: 1,
		Content:   map[string]interface{}{"url": "some_url", "text": "some_text", "title": "some_title"},
		TagId:     []int{1, 2, 3},
		FeatureId: 1,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now()}

	mockBannerService.EXPECT().SearchAllBanners(gomock.Any(), gomock.Any()).
		Return([]banner.BannerEntity{testBannerEntity}, true, nil)
	mockBannerService.EXPECT().SearchAllBanners(gomock.Any(), gomock.Any()).
		Return([]banner.BannerEntity{}, false, fmt.Errorf("unauthorized user"))
	mockBannerService.EXPECT().SearchAllBanners(gomock.Any(), gomock.Any()).
		Return([]banner.BannerEntity{}, false, nil)
	mockBannerService.EXPECT().SearchAllBanners(gomock.Any(), gomock.Any()).
		Return([]banner.BannerEntity{}, true, fmt.Errorf("some inner error"))

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Correct request with no params",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"http://localhost:8080/user_banner", nil),
				token: "some_token"},
			want: 200,
		},
		{
			name: "Request without token",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"http://localhost:8080/user_banner", nil),
				token: ""},
			want: 401,
		},
		{
			name: "Permission_denied",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"http://localhost:8080/user_banner", nil),
				token: "not_admin_token"},
			want: 403,
		},
		{
			name: "Internal error",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"http://localhost:8080/user_banner", nil),
				token: "admin"},
			want: 500,
		},
	}

	for _, tt := range tests {
		req := tt.args.r
		w := tt.args.w

		req.Header.Set("token", tt.args.token)

		bannerHandler.GetAllBanners(w, req)
		assert.Equal(t, tt.want, w.Code)
	}
}
