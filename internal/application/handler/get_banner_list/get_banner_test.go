package get_banner_list

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	getBanner "github.com/CyberPiess/banner_service/internal/application/handler/get_banner_list/mocks"
	bannerService "github.com/CyberPiess/banner_service/internal/domain/banner"
	"github.com/CyberPiess/banner_service/internal/infrastructure/logging"
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

	logger, err := logging.LoggerCreate(logging.Config{LogLevel: "info",
		LogFile: "get_banner_list_test.log"})
	if err != nil {
		log.Fatal("error init logger")
	}
	mockBannerService := getBanner.NewMockgetAllBannerService(ctrl)

	bannerHandler := NewGetAllBannersHandler(mockBannerService, logger)
	isActive := true

	testBannerEntity := bannerService.BannerEntity{ID: 1,
		Content:   map[string]interface{}{"url": "some_url", "text": "some_text", "title": "some_title"},
		TagId:     []int{1, 2, 3},
		FeatureId: 1,
		IsActive:  &isActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now()}

	mockBannerService.EXPECT().SearchAllBanners(gomock.Any(), gomock.Any()).
		Return([]bannerService.BannerEntity{testBannerEntity}, true, nil)
	mockBannerService.EXPECT().SearchAllBanners(gomock.Any(), gomock.Any()).
		Return([]bannerService.BannerEntity{}, false, nil)
	mockBannerService.EXPECT().SearchAllBanners(gomock.Any(), gomock.Any()).
		Return([]bannerService.BannerEntity{}, true, fmt.Errorf("some inner error"))

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
					"http://localhost:8080/banner", nil),
				token: ""},
			want: 401,
		},
		{
			name: "Permission_denied",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"http://localhost:8080/banner", nil),
				token: "not_admin_token"},
			want: 403,
		},
		{
			name: "Internal error",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"http://localhost:8080/banner", nil),
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
