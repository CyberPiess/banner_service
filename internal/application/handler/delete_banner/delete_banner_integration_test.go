package delete_banner_test

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	deleteBanner "github.com/CyberPiess/banner_service/internal/application/handler/delete_banner"
	banner_service "github.com/CyberPiess/banner_service/internal/domain/banner"
	"github.com/CyberPiess/banner_service/internal/infrastructure/logging"
	"github.com/CyberPiess/banner_service/internal/infrastructure/postgres"
	storage "github.com/CyberPiess/banner_service/internal/infrastructure/postgres/banner"
	rd "github.com/CyberPiess/banner_service/internal/infrastructure/redis"
	bannerCache "github.com/CyberPiess/banner_service/internal/infrastructure/redis/cache"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDbInstance *sql.DB
var testRedisInstance *redis.Client

func TestMain(m *testing.M) {
	testDB := postgres.SetupTestDatabase()
	testRedis := rd.SetupTestRedis()
	testDbInstance = testDB.DbInstance
	testRedisInstance = testRedis.RedisInstance
	defer testDB.TearDown()
	os.Exit(m.Run())
}

type args struct {
	w     *httptest.ResponseRecorder
	r     *http.Request
	token string
	id    string
}

func TestIntegrationDeleteBanner(t *testing.T) {

	logger, err := logging.LoggerCreate(logging.Config{LogLevel: "info",
		LogFile: "delete_banner_integration_test.log"})
	if err != nil {
		log.Fatal("error init logger")
	}

	bannerStore := storage.NewBannerRepository(testDbInstance, logger)
	bannerCache := bannerCache.NewBannerCache(testRedisInstance, logger)
	bannerService := banner_service.NewBannerService(bannerStore, bannerCache, logger)

	deleteBannerHandler := deleteBanner.NewDeleteBannerHandler(bannerService, logger)

	tests := []struct {
		name string
		args args
		want int
		body string
	}{
		{
			name: "Correct request",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodDelete,
					"http://localhost:8080/banner/{id}", nil),
				token: "admin_token",
				id:    "1"},
			want: 204,
			body: "",
		},
		{
			name: "Not Found",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodDelete,
					"http://localhost:8080/banner/{id}", nil),
				token: "admin_token",
				id:    "1"},
			want: 404,
			body: "",
		},
		{
			name: "User Token",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodDelete,
					"http://localhost:8080/banner/{id}", nil),
				token: "user_token",
				id:    "1"},
			want: 403,
			body: "",
		},
	}

	for _, tt := range tests {
		w := tt.args.w
		req := tt.args.r
		req.Header.Set("token", tt.args.token)
		req = mux.SetURLVars(req, map[string]string{"id": tt.args.id})
		deleteBannerHandler.DeleteBanner(w, req)
		data, err := io.ReadAll(w.Body)

		require.NoError(t, err)
		assert.Equal(t, tt.want, w.Code)
		assert.Equal(t, tt.body, string(data))
	}
}
