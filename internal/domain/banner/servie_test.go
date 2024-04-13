package banner_service

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	mock "github.com/CyberPiess/banner_service/internal/domain/banner/mocks"
	"github.com/CyberPiess/banner_service/internal/infrastructure/logging"
	banner_storage "github.com/CyberPiess/banner_service/internal/infrastructure/postgres/banner"
	redis "github.com/CyberPiess/banner_service/internal/infrastructure/redis/cache"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type argsForUserGetBanner struct {
	banberFilter GetFilter
	user         User
}

func TestSearchBanner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger, err := logging.LoggerCreate(logging.Config{LogLevel: "info",
		LogFile: "banner_storage_test.log"})
	if err != nil {
		log.Fatal("error init logger")
	}

	mockBannerStorage := mock.NewMockbannerStorage(ctrl)
	mockBannerCache := mock.NewMockredisCache(ctrl)

	bannerService := NewBannerService(mockBannerStorage, mockBannerCache, logger)

	bannerSQL := []banner_storage.BannerEntitySql{{
		Content: `{"some_string":"somestring"}`,
	}}

	isActive := false

	mockBannerStorage.EXPECT().IfTokenValid(gomock.Any()).Return(true, nil).Times(7)
	mockBannerCache.EXPECT().IfCacheExists(gomock.Any()).Return(int64(0), nil).Times(6)
	mockBannerStorage.EXPECT().IfBannerExists(gomock.Any(), gomock.Any()).Return(true, nil).Times(3)
	mockBannerStorage.EXPECT().Get(gomock.Any()).Return(bannerSQL, nil).Times(2)
	mockBannerCache.EXPECT().AddToCache(gomock.Any(), gomock.Any()).Return(nil)
	mockBannerCache.EXPECT().AddToCache(gomock.Any(), gomock.Any()).Return(fmt.Errorf("some error"))
	mockBannerStorage.EXPECT().Get(gomock.Any()).Return([]banner_storage.BannerEntitySql{}, fmt.Errorf("some error"))
	mockBannerStorage.EXPECT().IfBannerExists(gomock.Any(), gomock.Any()).Return(false, nil).Times(2)
	mockBannerCache.EXPECT().AddToCache(gomock.Any(), gomock.Any()).Return(nil)
	mockBannerCache.EXPECT().AddToCache(gomock.Any(), gomock.Any()).Return(fmt.Errorf("redis_error"))
	mockBannerStorage.EXPECT().IfBannerExists(gomock.Any(), gomock.Any()).Return(false, fmt.Errorf("some error"))
	mockBannerCache.EXPECT().IfCacheExists(gomock.Any()).Return(int64(1), nil)
	mockBannerCache.EXPECT().GetFromCache(gomock.Any()).Return(redis.RedisEntity{Content: `{"some_string":"somestring"}`}, nil)
	mockBannerStorage.EXPECT().IfTokenValid(gomock.Any()).Return(false, nil)
	mockBannerStorage.EXPECT().IfTokenValid(gomock.Any()).Return(false, fmt.Errorf("some error"))

	tests := []struct {
		name       string
		args       argsForUserGetBanner
		wantError  error
		wantBanner BannerEntity
		wantToken  bool
	}{
		{
			name: "Correct Data",
			args: argsForUserGetBanner{
				user: User{Token: "some_token"},
				banberFilter: GetFilter{TagId: 1,
					FeatureId: 1},
			},
			wantError:  nil,
			wantBanner: BannerEntity{Content: map[string]interface{}{"some_string": "somestring"}, IsActive: &isActive},
			wantToken:  true,
		},
		{
			name: "Error from Redis",
			args: argsForUserGetBanner{
				user: User{Token: "some_token"},
				banberFilter: GetFilter{TagId: 1,
					FeatureId: 1},
			},
			wantError:  nil,
			wantBanner: BannerEntity{Content: map[string]interface{}{"some_string": "somestring"}, IsActive: &isActive},
			wantToken:  true,
		},
		{
			name: "Internal error from db",
			args: argsForUserGetBanner{
				user: User{Token: "some_token"},
				banberFilter: GetFilter{TagId: 1,
					FeatureId: 1},
			},
			wantError:  fmt.Errorf("some error"),
			wantBanner: BannerEntity{},
			wantToken:  true,
		},
		{
			name: "Banner doesnot exist",
			args: argsForUserGetBanner{
				user: User{Token: "some_token"},
				banberFilter: GetFilter{TagId: 1,
					FeatureId: 1},
			},
			wantError:  nil,
			wantBanner: BannerEntity{},
			wantToken:  true,
		},
		{
			name: "Error adding to cahce",
			args: argsForUserGetBanner{
				user: User{Token: "some_token"},
				banberFilter: GetFilter{TagId: 1,
					FeatureId: 1},
			},
			wantError:  fmt.Errorf("redis_error"),
			wantBanner: BannerEntity{},
			wantToken:  true,
		},
		{
			name: "Inner error while getting banner doesnot exist",
			args: argsForUserGetBanner{
				user: User{Token: "some_token"},
				banberFilter: GetFilter{TagId: 1,
					FeatureId: 1},
			},
			wantError:  fmt.Errorf("some error"),
			wantBanner: BannerEntity{},
			wantToken:  true,
		},
		{
			name: "Banner in cache",
			args: argsForUserGetBanner{
				user: User{Token: "some_token"},
				banberFilter: GetFilter{TagId: 1,
					FeatureId: 1},
			},
			wantError:  nil,
			wantBanner: BannerEntity{Content: map[string]interface{}{"some_string": "somestring"}},
			wantToken:  true,
		},
		{
			name: "Invalid token",
			args: argsForUserGetBanner{
				user: User{Token: "some_token"},
				banberFilter: GetFilter{TagId: 1,
					FeatureId: 1},
			},
			wantError:  nil,
			wantBanner: BannerEntity{},
			wantToken:  false,
		},
		{
			name: "Inner error while grting token",
			args: argsForUserGetBanner{
				user: User{Token: "some_token"},
				banberFilter: GetFilter{TagId: 1,
					FeatureId: 1},
			},
			wantError:  fmt.Errorf("some error"),
			wantBanner: BannerEntity{},
			wantToken:  false,
		},
	}

	for _, tt := range tests {
		bannerReturned, returnedToken, err := bannerService.SearchBanner(tt.args.banberFilter, tt.args.user)
		assert.Equal(t, tt.wantError, err)
		assert.Equal(t, tt.wantBanner, bannerReturned)
		assert.Equal(t, tt.wantToken, returnedToken)
	}
}

type argsForGetAllBaners struct {
	bannerFilter GetAllFilter
	user         User
}

func TestSearchAllBanners(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBannerStorage := mock.NewMockbannerStorage(ctrl)
	mockBannerCache := mock.NewMockredisCache(ctrl)

	logger, err := logging.LoggerCreate(logging.Config{LogLevel: "info",
		LogFile: "banner_storage_test.log"})
	if err != nil {
		log.Fatal("error init logger")
	}

	bannerService := NewBannerService(mockBannerStorage, mockBannerCache, logger)
	var timeToRet sql.NullTime
	isActive := true

	bannerSQL := []banner_storage.BannerEntitySql{{
		ID:        1,
		Content:   `{"some_string":"somestring"}`,
		TagId:     []int{1, 2, 3},
		IsActive:  true,
		CreatedAt: timeToRet.Time,
		UpdatedAt: timeToRet.Time,
	}}

	mockBannerStorage.EXPECT().IfAdminTokenValid(gomock.Any()).Return(true, nil).Times(2)
	mockBannerStorage.EXPECT().GetAllBanners(gomock.Any()).Return(bannerSQL, nil)
	mockBannerStorage.EXPECT().GetAllBanners(gomock.Any()).Return([]banner_storage.BannerEntitySql{}, fmt.Errorf("some error"))
	mockBannerStorage.EXPECT().IfAdminTokenValid(gomock.Any()).Return(false, nil)
	mockBannerStorage.EXPECT().IfAdminTokenValid(gomock.Any()).Return(false, fmt.Errorf("some error"))

	tests := []struct {
		name       string
		args       argsForGetAllBaners
		wantError  error
		wantBanner []BannerEntity
		wantToken  bool
	}{
		{
			name: "Correct Data",
			args: argsForGetAllBaners{
				user:         User{Token: "some_token"},
				bannerFilter: GetAllFilter{TagId: 1},
			},
			wantError: nil,
			wantBanner: []BannerEntity{
				{
					ID:        1,
					Content:   map[string]interface{}{"some_string": "somestring"},
					TagId:     []int{1, 2, 3},
					IsActive:  &isActive,
					CreatedAt: timeToRet.Time,
					UpdatedAt: timeToRet.Time,
				},
			},
			wantToken: true,
		},
		{
			name: "Internal server error while getting banners",
			args: argsForGetAllBaners{
				user:         User{Token: "some_token"},
				bannerFilter: GetAllFilter{TagId: 1},
			},
			wantError:  fmt.Errorf("some error"),
			wantBanner: []BannerEntity{},
			wantToken:  true,
		},
		{
			name: "Access permited",
			args: argsForGetAllBaners{
				user:         User{Token: "some_token"},
				bannerFilter: GetAllFilter{TagId: 1},
			},
			wantError:  nil,
			wantBanner: []BannerEntity{},
			wantToken:  false,
		},
		{
			name: "Internal error while verifying token",
			args: argsForGetAllBaners{
				user:         User{Token: "some_token"},
				bannerFilter: GetAllFilter{TagId: 1},
			},
			wantError:  fmt.Errorf("some error"),
			wantBanner: []BannerEntity{},
			wantToken:  false,
		},
	}

	for _, tt := range tests {
		bannerReturned, returnedToken, err := bannerService.SearchAllBanners(tt.args.bannerFilter, tt.args.user)
		assert.Equal(t, tt.wantError, err)
		assert.Equal(t, tt.wantBanner, bannerReturned)
		assert.Equal(t, tt.wantToken, returnedToken)
	}
}

type argsForPutPostDelete struct {
	banner BannerEntity
	user   User
}

func TestPostBanner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBannerStorage := mock.NewMockbannerStorage(ctrl)
	mockBannerCache := mock.NewMockredisCache(ctrl)

	logger, err := logging.LoggerCreate(logging.Config{LogLevel: "info",
		LogFile: "banner_storage_test.log"})
	if err != nil {
		log.Fatal("error init logger")
	}

	bannerService := NewBannerService(mockBannerStorage, mockBannerCache, logger)

	mockBannerStorage.EXPECT().IfAdminTokenValid(gomock.Any()).Return(true, nil).Times(2)
	mockBannerStorage.EXPECT().PostBanner(gomock.Any()).Return(1, nil)
	mockBannerStorage.EXPECT().PostBanner(gomock.Any()).Return(0, fmt.Errorf("some error"))
	mockBannerStorage.EXPECT().IfAdminTokenValid(gomock.Any()).Return(false, nil)
	mockBannerStorage.EXPECT().IfAdminTokenValid(gomock.Any()).Return(false, fmt.Errorf("some error"))

	isActive := true

	tests := []struct {
		name      string
		args      argsForPutPostDelete
		wantError error
		wantId    int
		wantToken bool
	}{
		{
			name: "Correct Data",
			args: argsForPutPostDelete{
				user: User{Token: "some_token"},
				banner: BannerEntity{
					Content:   map[string]interface{}{"some_string": "somestring"},
					TagId:     []int{1, 2, 3},
					FeatureId: 1,
					IsActive:  &isActive,
					CreatedAt: time.Now(),
				},
			},
			wantError: nil,
			wantId:    1,
			wantToken: true,
		},
		{
			name: "Internal error while creating banner",
			args: argsForPutPostDelete{
				user: User{Token: "some_token"},
				banner: BannerEntity{
					Content:   map[string]interface{}{"some_string": "somestring"},
					TagId:     []int{1, 2, 3},
					FeatureId: 1,
					IsActive:  &isActive,
					CreatedAt: time.Now(),
				},
			},
			wantError: fmt.Errorf("some error"),
			wantId:    0,
			wantToken: true,
		},
		{
			name: "Access denied",
			args: argsForPutPostDelete{
				user: User{Token: "some_token"},
				banner: BannerEntity{
					Content:   map[string]interface{}{"some_string": "somestring"},
					TagId:     []int{1, 2, 3},
					FeatureId: 1,
					IsActive:  &isActive,
					CreatedAt: time.Now(),
				},
			},
			wantError: nil,
			wantId:    0,
			wantToken: false,
		},
		{
			name: "Inner error while verifying token",
			args: argsForPutPostDelete{
				user: User{Token: "some_token"},
				banner: BannerEntity{
					Content:   map[string]interface{}{"some_string": "somestring"},
					TagId:     []int{1, 2, 3},
					FeatureId: 1,
					IsActive:  &isActive,
					CreatedAt: time.Now(),
				},
			},
			wantError: fmt.Errorf("some error"),
			wantId:    0,
			wantToken: false,
		},
	}

	for _, tt := range tests {
		bannerIdReturned, returnedToken, err := bannerService.PostBanner(tt.args.banner, tt.args.user)
		assert.Equal(t, tt.wantError, err)
		assert.Equal(t, tt.wantId, bannerIdReturned)
		assert.Equal(t, tt.wantToken, returnedToken)
	}
}

func TestPutBanner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBannerStorage := mock.NewMockbannerStorage(ctrl)
	mockBannerCache := mock.NewMockredisCache(ctrl)

	isActive := true

	logger, err := logging.LoggerCreate(logging.Config{LogLevel: "info",
		LogFile: "banner_storage_test.log"})
	if err != nil {
		log.Fatal("error init logger")
	}

	bannerService := NewBannerService(mockBannerStorage, mockBannerCache, logger)

	mockBannerStorage.EXPECT().IfAdminTokenValid(gomock.Any()).Return(true, nil).Times(4)
	mockBannerStorage.EXPECT().SearchBannerByID(gomock.Any()).Return(true, nil).Times(2)
	mockBannerStorage.EXPECT().PutBanner(gomock.Any()).Return(nil)
	mockBannerStorage.EXPECT().PutBanner(gomock.Any()).Return(fmt.Errorf("some error"))
	mockBannerStorage.EXPECT().SearchBannerByID(gomock.Any()).Return(false, nil)
	mockBannerStorage.EXPECT().SearchBannerByID(gomock.Any()).Return(false, fmt.Errorf("some error"))
	mockBannerStorage.EXPECT().IfAdminTokenValid(gomock.Any()).Return(false, nil)
	mockBannerStorage.EXPECT().IfAdminTokenValid(gomock.Any()).Return(false, fmt.Errorf("some error"))

	tests := []struct {
		name            string
		args            argsForPutPostDelete
		wantError       error
		wantBannerExist bool
		wantToken       bool
	}{
		{
			name: "Correct Data",
			args: argsForPutPostDelete{
				user: User{Token: "some_token"},
				banner: BannerEntity{
					Content:   map[string]interface{}{"some_string": "somestring"},
					TagId:     []int{1, 2, 3},
					FeatureId: 1,
					IsActive:  &isActive,
					CreatedAt: time.Now(),
				},
			},
			wantError:       nil,
			wantBannerExist: true,
			wantToken:       true,
		},
		{
			name: "Internal error while updating",
			args: argsForPutPostDelete{
				user: User{Token: "some_token"},
				banner: BannerEntity{
					Content:   map[string]interface{}{"some_string": "somestring"},
					TagId:     []int{1, 2, 3},
					FeatureId: 1,
					IsActive:  &isActive,
					CreatedAt: time.Now(),
				},
			},
			wantError:       fmt.Errorf("some error"),
			wantBannerExist: true,
			wantToken:       true,
		},
		{
			name: "Banner not found",
			args: argsForPutPostDelete{
				user: User{Token: "some_token"},
				banner: BannerEntity{
					Content:   map[string]interface{}{"some_string": "somestring"},
					TagId:     []int{1, 2, 3},
					FeatureId: 1,
					IsActive:  &isActive,
					CreatedAt: time.Now(),
				},
			},
			wantError:       nil,
			wantBannerExist: false,
			wantToken:       true,
		},
		{
			name: "Inner error while banner search",
			args: argsForPutPostDelete{
				user: User{Token: "some_token"},
				banner: BannerEntity{
					Content:   map[string]interface{}{"some_string": "somestring"},
					TagId:     []int{1, 2, 3},
					FeatureId: 1,
					IsActive:  &isActive,
					CreatedAt: time.Now(),
				},
			},
			wantError:       fmt.Errorf("some error"),
			wantBannerExist: false,
			wantToken:       true,
		},
		{
			name: "Access permited",
			args: argsForPutPostDelete{
				user: User{Token: "some_token"},
				banner: BannerEntity{
					Content:   map[string]interface{}{"some_string": "somestring"},
					TagId:     []int{1, 2, 3},
					FeatureId: 1,
					IsActive:  &isActive,
					CreatedAt: time.Now(),
				},
			},
			wantError:       nil,
			wantBannerExist: false,
			wantToken:       false,
		},
		{
			name: "Inner error while validating token",
			args: argsForPutPostDelete{
				user: User{Token: "some_token"},
				banner: BannerEntity{
					Content:   map[string]interface{}{"some_string": "somestring"},
					TagId:     []int{1, 2, 3},
					FeatureId: 1,
					IsActive:  &isActive,
					CreatedAt: time.Now(),
				},
			},
			wantError:       fmt.Errorf("some error"),
			wantBannerExist: false,
			wantToken:       false,
		},
	}

	for _, tt := range tests {
		updateResult, resultToken, err := bannerService.PutBanner(tt.args.banner, tt.args.user)
		assert.Equal(t, tt.wantError, err)
		assert.Equal(t, tt.wantBannerExist, updateResult)
		assert.Equal(t, tt.wantToken, resultToken)
	}
}

func TestDeleteBanner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBannerStorage := mock.NewMockbannerStorage(ctrl)
	mockBannerCache := mock.NewMockredisCache(ctrl)

	logger, err := logging.LoggerCreate(logging.Config{LogLevel: "info",
		LogFile: "banner_storage_test.log"})
	if err != nil {
		log.Fatal("error init logger")
	}

	bannerService := NewBannerService(mockBannerStorage, mockBannerCache, logger)

	mockBannerStorage.EXPECT().IfAdminTokenValid(gomock.Any()).Return(true, nil).Times(4)
	mockBannerStorage.EXPECT().SearchBannerByID(gomock.Any()).Return(true, nil).Times(2)
	mockBannerStorage.EXPECT().DeleteBanner(gomock.Any()).Return(nil)
	mockBannerStorage.EXPECT().DeleteBanner(gomock.Any()).Return(fmt.Errorf("some error"))
	mockBannerStorage.EXPECT().SearchBannerByID(gomock.Any()).Return(false, nil)
	mockBannerStorage.EXPECT().SearchBannerByID(gomock.Any()).Return(false, fmt.Errorf("some error"))
	mockBannerStorage.EXPECT().IfAdminTokenValid(gomock.Any()).Return(false, nil)
	mockBannerStorage.EXPECT().IfAdminTokenValid(gomock.Any()).Return(false, fmt.Errorf("some error"))

	tests := []struct {
		name            string
		args            argsForPutPostDelete
		wantError       error
		wantBannerExist bool
		wantToken       bool
	}{
		{
			name: "Correct Data",
			args: argsForPutPostDelete{
				user: User{Token: "some_token"},
				banner: BannerEntity{
					ID: 1,
				},
			},
			wantError:       nil,
			wantBannerExist: true,
			wantToken:       true,
		},
		{
			name: "Internal error while deleting",
			args: argsForPutPostDelete{
				user: User{Token: "some_token"},
				banner: BannerEntity{
					ID: 1,
				},
			},
			wantError:       fmt.Errorf("some error"),
			wantBannerExist: true,
			wantToken:       true,
		},
		{
			name: "Banner not found",
			args: argsForPutPostDelete{
				user: User{Token: "some_token"},
				banner: BannerEntity{
					ID: 1,
				},
			},
			wantError:       nil,
			wantBannerExist: false,
			wantToken:       true,
		},
		{
			name: "Inner error while banner search",
			args: argsForPutPostDelete{
				user: User{Token: "some_token"},
				banner: BannerEntity{
					ID: 1,
				},
			},
			wantError:       fmt.Errorf("some error"),
			wantBannerExist: false,
			wantToken:       true,
		},
		{
			name: "Access permited",
			args: argsForPutPostDelete{
				user: User{Token: "some_token"},
				banner: BannerEntity{
					ID: 1,
				},
			},
			wantError:       nil,
			wantBannerExist: false,
			wantToken:       false,
		},
		{
			name: "Inner error while validating token",
			args: argsForPutPostDelete{
				user: User{Token: "some_token"},
				banner: BannerEntity{
					ID: 1,
				},
			},
			wantError:       fmt.Errorf("some error"),
			wantBannerExist: false,
			wantToken:       false,
		},
	}

	for _, tt := range tests {
		updateResult, resultToken, err := bannerService.DeleteBanner(tt.args.banner, tt.args.user)
		assert.Equal(t, tt.wantError, err)
		assert.Equal(t, tt.wantBannerExist, updateResult)
		assert.Equal(t, tt.wantToken, resultToken)
	}
}
