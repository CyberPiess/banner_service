package create_banner

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	postBanner "github.com/CyberPiess/banner_service/internal/application/handler/create_banner/mocks"
	"github.com/CyberPiess/banner_service/internal/infrastructure/logging"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type args struct {
	w     *httptest.ResponseRecorder
	r     *http.Request
	token string
}

func TestPostBanner(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBannerService := postBanner.NewMockpostBannerService(ctrl)

	logger, err := logging.LoggerCreate(logging.Config{LogLevel: "info",
		LogFile: "create_banner_test.log"})
	if err != nil {
		log.Fatal("error init logger")
	}

	postBannerHandler := NewPostBannerHandler(mockBannerService, logger)

	mockBannerService.EXPECT().PostBanner(gomock.Any(), gomock.Any()).Return(1, true, nil)
	mockBannerService.EXPECT().PostBanner(gomock.Any(), gomock.Any()).Return(0, false, nil)
	mockBannerService.EXPECT().PostBanner(gomock.Any(), gomock.Any()).Return(0, false, fmt.Errorf("some error"))

	validRequestBody := "{ \"tag_ids\": [1,2,3], \"feature_id\": 1, \"content\": {\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}, \"is_active\": false}"
	absentTagId := "{ \"feature_id\": 1, \"content\": {\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}, \"is_active\": false}"
	absentFeatureId := "{\"tag_ids\": [1,2,3], \"content\": {\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}, \"is_active\": false}"
	absentContent := "{\"tag_ids\": [1,2,3], \"feature_id\": 1, \"is_active\": false}"
	absentIsActive := "{\"tag_ids\": [1,2,3], \"content\": {\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}}"

	validBody := strings.NewReader(validRequestBody)
	validBody2 := strings.NewReader(validRequestBody)
	validBody3 := strings.NewReader(validRequestBody)
	validBody4 := strings.NewReader(validRequestBody)
	noTagID := strings.NewReader(absentTagId)
	noFeatureID := strings.NewReader(absentFeatureId)
	noContent := strings.NewReader(absentContent)
	noIsActive := strings.NewReader(absentIsActive)

	test := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Correct data",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"http://localhost:8080/banner", validBody),
				token: "some_token"},
			want: 201,
		},
		{
			name: "Empty body",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"http://localhost:8080/banner", nil),
				token: "some_token"},
			want: 400,
		},
		{
			name: "Empty tag_ids",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"http://localhost:8080/banner", noTagID),
				token: "some_token"},
			want: 400,
		},
		{
			name: "Empty feature_id",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"http://localhost:8080/banner", noFeatureID),
				token: "some_token"},
			want: 400,
		},
		{
			name: "Empty feature_id",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"http://localhost:8080/banner", noContent),
				token: "some_token"},
			want: 400,
		},
		{
			name: "Empty feature_id",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"http://localhost:8080/banner", noIsActive),
				token: "some_token"},
			want: 400,
		},
		{
			name: "Unauthorized user",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"http://localhost:8080/banner", validBody2),
				token: ""},
			want: 401,
		},
		{
			name: "Access denied",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"http://localhost:8080/banner", validBody3),
				token: "some_token"},
			want: 403,
		},
		{
			name: "Inner error",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"http://localhost:8080/banner", validBody4),
				token: "some_token"},
			want: 500,
		},
	}

	for _, tt := range test {
		req := tt.args.r
		w := tt.args.w

		req.Header.Set("token", tt.args.token)
		postBannerHandler.PostBanner(w, req)
		assert.Equal(t, tt.want, w.Code)
	}

}
