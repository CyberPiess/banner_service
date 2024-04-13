package delete_banner

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	delBanner "github.com/CyberPiess/banner_service/internal/application/handler/delete_banner/mocks"
	"github.com/CyberPiess/banner_service/internal/infrastructure/logging"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type args struct {
	w     *httptest.ResponseRecorder
	r     *http.Request
	token string
	id    string
}

func TestDeleteBanner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBannerService := delBanner.NewMockdeleteBannerService(ctrl)

	logger, err := logging.LoggerCreate(logging.Config{LogLevel: "info",
		LogFile: "delete_banner_test.log"})
	if err != nil {
		log.Fatal("error init logger")
	}

	delBannerHandler := NewDeleteBannerHandler(mockBannerService, logger)

	mockBannerService.EXPECT().DeleteBanner(gomock.Any(), gomock.Any()).Return(true, true, nil)
	mockBannerService.EXPECT().DeleteBanner(gomock.Any(), gomock.Any()).Return(false, false, nil)
	mockBannerService.EXPECT().DeleteBanner(gomock.Any(), gomock.Any()).Return(false, true, nil)
	mockBannerService.EXPECT().DeleteBanner(gomock.Any(), gomock.Any()).Return(true, true, fmt.Errorf("some internal error"))

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
					http.MethodPut,
					"http://localhost:8080/banner/{id}", nil),
				token: "some_token",
				id:    "1",
			},
			want: 204,
		},
		{
			name: "Incorrect id",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPut,
					"http://localhost:8080/banner/{id}", nil),
				token: "some_token",
				id:    "a",
			},
			want: 400,
		},
		{
			name: "Empty token",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPut,
					"http://localhost:8080/banner/{id}", nil),
				token: "",
				id:    "1",
			},
			want: 401,
		},
		{
			name: "Access denied",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPut,
					"http://localhost:8080/banner/{id}", nil),
				token: "some_not_admin_token",
				id:    "1",
			},
			want: 403,
		},
		{
			name: "Banner not found",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPut,
					"http://localhost:8080/banner/{id}", nil),
				token: "some_token",
				id:    "1",
			},
			want: 404,
		},
		{
			name: "Internal error",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPut,
					"http://localhost:8080/banner/{id}", nil),
				token: "some_token",
				id:    "1",
			},
			want: 500,
		},
	}

	for _, tt := range test {
		w := tt.args.w
		req := tt.args.r

		req = mux.SetURLVars(req, map[string]string{"id": tt.args.id})
		req.Header.Set("token", tt.args.token)
		delBannerHandler.DeleteBanner(w, req)
		assert.Equal(t, tt.want, w.Code)
	}
}
