package userbanner_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	userBanner "github.com/CyberPiess/banner_sevice/internal/application/handler/get_user_banner"
	"github.com/CyberPiess/banner_sevice/internal/domain/banner"
	"github.com/CyberPiess/banner_sevice/internal/infrastructure/postgres"
	storage "github.com/CyberPiess/banner_sevice/internal/infrastructure/postgres/banner"

	"github.com/stretchr/testify/assert"
)

var testDbInstance *sql.DB

func TestMain(m *testing.M) {
	testDB := postgres.SetupTestDatabase()
	testDbInstance = testDB.DbInstance
	defer testDB.TearDown()
	os.Exit(m.Run())
}

type args struct {
	w     *httptest.ResponseRecorder
	r     *http.Request
	token string
}

func TestIntegrationGetUserBanner(t *testing.T) {
	bannerStore := storage.NewBannerRepository(testDbInstance)
	bannerService := banner.NewBannerService(bannerStore)

	bannerHandler := userBanner.NewBannerHandler(bannerService)

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
					"http://localhost:8080/user_banner?tag_id=1&feature_id=1", nil),
				token: "user_token"},
			want: 200,
		},
	}

	for _, tt := range tests {
		w := tt.args.w
		req := tt.args.r
		req.Header.Set("token", tt.args.token)

		bannerHandler.GetUserBanner(w, req)
		assert.Equal(t, tt.want, w.Code)
	}
}
