package postbanner

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	postBanner "github.com/CyberPiess/banner_sevice/internal/application/handler/post_banner/mocks"
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

	postBannerHandler := NewPostBannerHandler(mockBannerService)

	mockBannerService.EXPECT().PostBanner(gomock.Any(), gomock.Any()).Return(1, true, nil)
	mockBannerService.EXPECT().PostBanner(gomock.Any(), gomock.Any()).Return(0, false, fmt.Errorf("unauthorized user"))
	mockBannerService.EXPECT().PostBanner(gomock.Any(), gomock.Any()).Return(0, false, nil)
	mockBannerService.EXPECT().PostBanner(gomock.Any(), gomock.Any()).Return(0, false, fmt.Errorf("some error"))

	objectToSerializeToJSON := "{ \"tag_ids\": [1,2,3], \"feature_id\": 1, \"content\": {\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}, \"is_active\": false}"

	validBody := strings.NewReader(objectToSerializeToJSON)
	validBody2 := strings.NewReader(objectToSerializeToJSON)
	validBody3 := strings.NewReader(objectToSerializeToJSON)
	validBody4 := strings.NewReader(objectToSerializeToJSON)

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
			name: "Unauthorized user",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"http://localhost:8080/banner", validBody3),
				token: "some_token"},
			want: 403,
		},
		{
			name: "Unauthorized user",
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
