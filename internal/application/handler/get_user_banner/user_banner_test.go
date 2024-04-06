package userbanner

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	appBanner "github.com/CyberPiess/banner_sevice/internal/application/handler/get_user_banner/mocks"
	"github.com/CyberPiess/banner_sevice/internal/domain/banner"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type args struct {
	w     *httptest.ResponseRecorder
	r     *http.Request
	token string
}

func TestGetUserBanner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBannerService := appBanner.NewMockbannerService(ctrl)

	bannerHandler := NewBannerHandler(mockBannerService)

	mockBannerService.EXPECT().SearchBanner(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(banner.BannerEntity{Content: map[string]interface{}{"url": "some_url", "text": "some_text", "title": "some_title"}},
			true, nil).Times(2)

	mockBannerService.EXPECT().SearchBanner(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(banner.BannerEntity{}, true, nil)

	mockBannerService.EXPECT().SearchBanner(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(banner.BannerEntity{Content: map[string]interface{}{"url": "some_url", "text": "some_text", "title": "some_title"}},
			false,
			nil)

	mockBannerService.EXPECT().SearchBanner(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(banner.BannerEntity{Content: map[string]interface{}{"url": "some_url", "text": "some_text", "title": "some_title"}},
			true,
			fmt.Errorf("some error from db"))

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Correct request without use_last_revision",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"http://localhost:8080/user_banner?tag_id=1&feature_id=123", nil),
				token: "some_token"},
			want: 200,
		},
		{
			name: "Correct request with use_last_revision",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"http://localhost:8080/user_banner?tag_id=1&feature_id=123&use_last_rebision=true", nil),
				token: "some_token"},
			want: 200,
		},
		{
			name: "Empty token",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"http://localhost:8080/user_banner?tag_id=1&feature_id=123&use_last_rebision=true", nil),
				token: ""},
			want: 401,
		},
		{
			name: "No tag_id",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"http://localhost:8080/user_banner?feature_id=123", nil),
				token: ""},
			want: 400,
		},
		{
			name: "No feature_id",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"http://localhost:8080/user_banner?tag_id=1", nil),
				token: ""},
			want: 400,
		},
		{
			name: "No query params",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"http://localhost:8080/user_banner", nil),
				token: ""},
			want: 400,
		},
		{
			name: "Wrong format",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"http://localhost:8080/user_banne?tag_id=a&feature_id=123r", nil),
				token: ""},
			want: 400,
		},
		{
			name: "Method POST",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"http://localhost:8080/users", nil),
				token: ""},
			want: 400,
		},
		{
			name: "If returned empty content",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"http://localhost:8080/user_banner?tag_id=1&feature_id=123", nil),
				token: "some_token"},
			want: 404,
		},
		{
			name: "If userToken wasnot found in token tables",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"http://localhost:8080/user_banner?tag_id=1&feature_id=123", nil),
				token: "some_token"},
			want: 403,
		},
		{
			name: "Inner error",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"http://localhost:8080/user_banner?tag_id=1&feature_id=123", nil),
				token: "some_token"},
			want: 500,
		},
	}

	for _, tt := range tests {
		req := tt.args.r
		w := tt.args.w

		req.Header.Set("token", tt.args.token)

		bannerHandler.GetUserBanner(w, req)
		assert.Equal(t, tt.want, w.Code)
	}

}
