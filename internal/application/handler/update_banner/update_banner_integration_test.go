package update_banner_test

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	updateBanner "github.com/CyberPiess/banner_service/internal/application/handler/update_banner"
	bannerService "github.com/CyberPiess/banner_service/internal/domain/banner"
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

func TestIntegrationUpdateBanner(t *testing.T) {

	logger, err := logging.LoggerCreate(logging.Config{LogLevel: "info",
		LogFile: "update_banner_integration_test.log"})
	if err != nil {
		log.Fatal("error init logger")
	}

	bannerStore := storage.NewBannerRepository(testDbInstance, logger)
	bannerCache := bannerCache.NewBannerCache(testRedisInstance, logger)
	bannerService := bannerService.NewBannerService(bannerStore, bannerCache, logger)
	updateBannerHandler := updateBanner.NewPutBannerHandler(bannerService, logger)

	validRequestBody := "{ \"tag_ids\": [81, 93, 54], \"feature_id\": 55, \"content\": {\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}, \"is_active\": false}"
	validBody := strings.NewReader(validRequestBody)
	validBody2 := strings.NewReader(validRequestBody)

	tests := []struct {
		name string
		args args
		want int
		body string
	}{
		{
			name: "Access denied",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodDelete,
					"http://localhost:8080/banner/{id}", validBody),
				token: "admin_token",
				id:    "1",
			},
			want: 200,
			body: ``,
		},
		{
			name: "Access denied",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodDelete,
					"http://localhost:8080/banner/{id}", validBody2),
				token: "user_token",
				id:    "1",
			},
			want: 403,
			body: ``,
		},
	}

	for _, tt := range tests {
		w := tt.args.w
		req := tt.args.r
		req.Header.Set("token", tt.args.token)
		req = mux.SetURLVars(req, map[string]string{"id": tt.args.id})
		updateBannerHandler.PutBanner(w, req)
		data, err := io.ReadAll(w.Body)

		require.NoError(t, err)
		assert.Equal(t, tt.want, w.Code)
		assert.Equal(t, tt.body, string(data))
	}
}
