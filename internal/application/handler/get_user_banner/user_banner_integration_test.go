package user_banner_test

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	userBanner "github.com/CyberPiess/banner_service/internal/application/handler/get_user_banner"
	bannerService "github.com/CyberPiess/banner_service/internal/domain/banner"
	"github.com/CyberPiess/banner_service/internal/infrastructure/logging"
	"github.com/CyberPiess/banner_service/internal/infrastructure/postgres"
	storage "github.com/CyberPiess/banner_service/internal/infrastructure/postgres/banner"
	rd "github.com/CyberPiess/banner_service/internal/infrastructure/redis"
	bannerCachePkg "github.com/CyberPiess/banner_service/internal/infrastructure/redis/cache"
	"github.com/go-redis/redis"

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
	defer testRedis.TearDown()
	os.Exit(m.Run())
}

type args struct {
	w     *httptest.ResponseRecorder
	r     *http.Request
	token string
}

func TestIntegrationGetUserBanner(t *testing.T) {

	logger, err := logging.LoggerCreate(logging.Config{LogLevel: "info",
		LogFile: "get_user_banner_integration_test.log"})
	if err != nil {
		log.Fatal("error init logger")
	}

	bannerStore := storage.NewBannerRepository(testDbInstance, logger)
	bannerCache := bannerCachePkg.NewBannerCache(testRedisInstance, logger)
	testRedisDTO := bannerCachePkg.RedisEntity{Content: `{"title": "sometitle"}`}
	bannerCache.AddToCache("tag_id=5&feature_id=1", testRedisDTO)
	bannerService := bannerService.NewBannerService(bannerStore, bannerCache, logger)

	bannerHandler := userBanner.NewBannerHandler(bannerService, logger)

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
					http.MethodGet,
					"http://localhost:8080/user_banner?tag_id=3&feature_id=1", nil),
				token: "user_token"},
			want: 200,
			body: `{"text":"some_text","title":"some_title","url":"some_url"}`,
		},
		{
			name: "Correct request to cache",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"http://localhost:8080/user_banner?tag_id=5&feature_id=1", nil),
				token: "user_token"},
			want: 200,
			body: `{"title":"sometitle"}`,
		},
		{
			name: "No banner",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"http://localhost:8080/user_banner?tag_id=6&feature_id=1", nil),
				token: "user_token"},
			want: 404,
			body: ``,
		},
	}

	for _, tt := range tests {
		w := tt.args.w
		req := tt.args.r
		req.Header.Set("token", tt.args.token)

		bannerHandler.GetUserBanner(w, req)
		data, err := io.ReadAll(w.Body)
		require.NoError(t, err)
		assert.Equal(t, tt.want, w.Code)
		assert.Equal(t, tt.body, string(data))

	}
}
