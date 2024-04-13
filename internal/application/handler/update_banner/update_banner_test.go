package update_banner

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	putBanner "github.com/CyberPiess/banner_service/internal/application/handler/update_banner/mocks"
	"github.com/CyberPiess/banner_service/internal/infrastructure/logging"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type args struct {
	w     *httptest.ResponseRecorder
	r     *http.Request
	token string
}

func TestPutBanner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger, err := logging.LoggerCreate(logging.Config{LogLevel: "info",
		LogFile: "update_banner_test.log"})
	if err != nil {
		log.Fatal("error init logger")
	}

	mockBannerService := putBanner.NewMockputBannerService(ctrl)
	putBannerHandler := NewPutBannerHandler(mockBannerService, logger)

	mockBannerService.EXPECT().PutBanner(gomock.Any(), gomock.Any()).Return(true, true, nil)
	mockBannerService.EXPECT().PutBanner(gomock.Any(), gomock.Any()).Return(false, false, nil)
	mockBannerService.EXPECT().PutBanner(gomock.Any(), gomock.Any()).Return(false, true, nil)
	mockBannerService.EXPECT().PutBanner(gomock.Any(), gomock.Any()).Return(true, true, fmt.Errorf("some inner error"))

	objectToSerializeToJSON := "{ \"tag_ids\": [1,2,3], \"feature_id\": 1, \"content\": {\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}, \"is_active\": false}"

	validBody := strings.NewReader(objectToSerializeToJSON)
	validBody2 := strings.NewReader(objectToSerializeToJSON)
	validBody3 := strings.NewReader(objectToSerializeToJSON)
	validBody4 := strings.NewReader(objectToSerializeToJSON)
	validBody5 := strings.NewReader(objectToSerializeToJSON)

	test := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Coorect data",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPut,
					"http://localhost:8080/banner/{id}", validBody),
				token: "some_token"},
			want: 200,
		},
		{
			name: "Empty body",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"http://localhost:8080/banner/{id}", nil),
				token: "some_token"},
			want: 400,
		},
		{
			name: "Unauthorized user",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"http://localhost:8080/banner/{id}", validBody2),
				token: ""},
			want: 401,
		},
		{
			name: "Access denied",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"http://localhost:8080/banner/{id}", validBody3),
				token: "some_token"},
			want: 403,
		},
		{
			name: "Banner not found",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"http://localhost:8080/banner/{id}", validBody4),
				token: "some_token"},
			want: 404,
		},
		{
			name: "Inner error",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"http://localhost:8080/banner/{id}", validBody5),
				token: "some_token"},
			want: 500,
		},
	}

	for _, tt := range test {
		w := tt.args.w
		req := tt.args.r
		req = mux.SetURLVars(req, map[string]string{"id": "1"})

		req.Header.Set("token", tt.args.token)
		putBannerHandler.PutBanner(w, req)
		assert.Equal(t, tt.want, w.Code)
	}
}
