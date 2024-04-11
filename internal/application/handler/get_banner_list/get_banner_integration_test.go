package getbannerlist_test

import (
	"database/sql"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	adminBanner "github.com/CyberPiess/banner_sevice/internal/application/handler/get_banner_list"
	"github.com/CyberPiess/banner_sevice/internal/domain/banner"
	"github.com/CyberPiess/banner_sevice/internal/infrastructure/postgres"
	storage "github.com/CyberPiess/banner_sevice/internal/infrastructure/postgres/banner"
	rd "github.com/CyberPiess/banner_sevice/internal/infrastructure/redis"
	bannerCache "github.com/CyberPiess/banner_sevice/internal/infrastructure/redis/cache"
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
	os.Exit(m.Run())
}

type args struct {
	w     *httptest.ResponseRecorder
	r     *http.Request
	token string
}

func TestIntegrationGetAllBanners(t *testing.T) {

	bannerStore := storage.NewBannerRepository(testDbInstance)
	bannerCache := bannerCache.NewBannerCache(testRedisInstance)
	bannerService := banner.NewBannerService(bannerStore, bannerCache)

	getAllBannersHandler := adminBanner.NewGetAllBannersHandler(bannerService)

	responseTagIdFeatureId := `[{"banner_id":1,"tag_ids":[1,2,3],"feature_id":1,"content":{"text":"some_text","title":"some_title","url":"some_url"},"is_active":true,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}]`
	responseTagId := `[{"banner_id":1,"tag_ids":[1,2,3],"feature_id":1,"content":{"text":"some_text","title":"some_title","url":"some_url"},"is_active":true,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"},{"banner_id":3,"tag_ids":[1,2,3],"feature_id":2,"content":{"text":"some_text","title":"some_title","url":"some_url"},"is_active":false,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}]`
	responseFeatureId := `[{"banner_id":1,"tag_ids":[1,2,3],"feature_id":1,"content":{"text":"some_text","title":"some_title","url":"some_url"},"is_active":true,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"},{"banner_id":2,"tag_ids":[4,5,6],"feature_id":1,"content":{"text":"some_text","title":"some_title","url":"some_url"},"is_active":false,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}]`
	responseAllBanners := `[{"banner_id":1,"tag_ids":[1,2,3],"feature_id":1,"content":{"text":"some_text","title":"some_title","url":"some_url"},"is_active":true,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"},{"banner_id":2,"tag_ids":[4,5,6],"feature_id":1,"content":{"text":"some_text","title":"some_title","url":"some_url"},"is_active":false,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"},{"banner_id":3,"tag_ids":[1,2,3],"feature_id":2,"content":{"text":"some_text","title":"some_title","url":"some_url"},"is_active":false,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}]`
	responseLimitOne := `[{"banner_id":1,"tag_ids":[1,2,3],"feature_id":1,"content":{"text":"some_text","title":"some_title","url":"some_url"},"is_active":true,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}]`
	responseOffset := `[{"banner_id":3,"tag_ids":[1,2,3],"feature_id":2,"content":{"text":"some_text","title":"some_title","url":"some_url"},"is_active":false,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}]`
	responseLimitAndOffset := `[{"banner_id":2,"tag_ids":[4,5,6],"feature_id":1,"content":{"text":"some_text","title":"some_title","url":"some_url"},"is_active":false,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}]`

	tests := []struct {
		name string
		args args
		want int
		body string
	}{
		{
			name: "Correct request with tag_id and feature_id",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodDelete,
					"http://localhost:8080/banner?tag_id=1&feature_id=1", nil),
				token: "admin_token",
			},
			want: 200,
			body: responseTagIdFeatureId,
		},
		{
			name: "Correct request with tag_id",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodDelete,
					"http://localhost:8080/banner?tag_id=1", nil),
				token: "admin_token",
			},
			want: 200,
			body: responseTagId,
		},
		{
			name: "Correct request with feature_id",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodDelete,
					"http://localhost:8080/banner?feature_id=1", nil),
				token: "admin_token",
			},
			want: 200,
			body: responseFeatureId,
		},
		{
			name: "Correct request for all banners",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodDelete,
					"http://localhost:8080/banner", nil),
				token: "admin_token",
			},
			want: 200,
			body: responseAllBanners,
		},
		{
			name: "Correct request for all banners with limit",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodDelete,
					"http://localhost:8080/banner?limit=1", nil),
				token: "admin_token",
			},
			want: 200,
			body: responseLimitOne,
		},
		{
			name: "Correct request for all banners with offset",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodDelete,
					"http://localhost:8080/banner?offset=2", nil),
				token: "admin_token",
			},
			want: 200,
			body: responseOffset,
		},
		{
			name: "Correct request for all banners limit and offset",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodDelete,
					"http://localhost:8080/banner?limit=1&offset=1", nil),
				token: "admin_token",
			},
			want: 200,
			body: responseLimitAndOffset,
		},
		{
			name: "Correct request returning no banners",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodDelete,
					"http://localhost:8080/banner?tag_id=10", nil),
				token: "admin_token",
			},
			want: 200,
			body: "[]",
		},
		{
			name: "Access permited",
			args: args{w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodDelete,
					"http://localhost:8080/banner?tag_id=10", nil),
				token: "user_token",
			},
			want: 403,
			body: "",
		},
	}

	for _, tt := range tests {
		w := tt.args.w
		req := tt.args.r
		req.Header.Set("token", tt.args.token)
		getAllBannersHandler.GetAllBanners(w, req)
		data, err := io.ReadAll(w.Body)

		require.NoError(t, err)
		assert.Equal(t, tt.want, w.Code)
		assert.Equal(t, tt.body, string(data))
	}
}
