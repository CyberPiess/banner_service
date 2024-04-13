//go:generate mockgen -source=service.go -destination=mocks/mock.go
package banner_service

import (
	"encoding/json"
	"fmt"
	"time"

	storage "github.com/CyberPiess/banner_service/internal/infrastructure/postgres/banner"
	redis "github.com/CyberPiess/banner_service/internal/infrastructure/redis/cache"
	"github.com/sirupsen/logrus"
)

type bannerStorage interface {
	Get(bannerParams storage.GetUserBannerCriteria) ([]storage.BannerEntitySql, error)
	IfTokenValid(token string) (bool, error)
	IfBannerExists(featureId int, tagId int) (bool, error)
	IfAdminTokenValid(token string) (bool, error)
	SearchBannerByID(bannerID int) (bool, error)
	GetAllBanners(bannerParams storage.GetBannersListCriteria) ([]storage.BannerEntitySql, error)
	PostBanner(postBannerParams storage.BannerPutPostCriteria) (int, error)
	PutBanner(putBannerParams storage.BannerPutPostCriteria) error
	DeleteBanner(deleteBannerParams storage.BannerPutPostCriteria) error
}

type logger interface {
	WithFields(fields logrus.Fields) *logrus.Entry
}

type redisCache interface {
	AddToCache(key string, redisDTO redis.RedisEntity) error
	IfCacheExists(key string) (int64, error)
	GetFromCache(key string) (redis.RedisEntity, error)
}

type BannerService struct {
	store  bannerStorage
	redis  redisCache
	logger logger
}

func NewBannerService(storage bannerStorage, redis redisCache, logger logger) *BannerService {
	return &BannerService{store: storage,
		redis:  redis,
		logger: logger}
}

func (b *BannerService) SearchBanner(bannerFilter GetFilter, user User) (BannerEntity, bool, error) {
	validToken, err := b.store.IfTokenValid(user.Token)
	if err != nil {
		return BannerEntity{}, validToken, err
	}
	if !validToken {
		return BannerEntity{}, validToken, nil
	}
	cacheKey := b.createCacheKey(bannerFilter.TagId, bannerFilter.FeatureId)

	if !bannerFilter.UseLastRevision {

		exists, err := b.redis.IfCacheExists(cacheKey)
		if err != nil {
			return BannerEntity{}, validToken, err
		}

		if exists == 1 {
			var bannerEntity BannerEntity
			foundInCahce, err := b.redis.GetFromCache(cacheKey)
			if err != nil {
				return BannerEntity{}, validToken, err
			}
			if foundInCahce.Content != "" {
				err := json.Unmarshal([]byte(foundInCahce.Content), &bannerEntity.Content)
				if err != nil {
					return BannerEntity{}, validToken, err
				}
			}
			return bannerEntity, validToken, err
		}

	}

	bannerExists, err := b.store.IfBannerExists(bannerFilter.TagId, bannerFilter.FeatureId)
	if err != nil {
		return BannerEntity{}, validToken, err
	}
	if !bannerExists {
		err = b.redis.AddToCache(cacheKey, redis.RedisEntity{Content: ""})
		if err != nil {
			return BannerEntity{}, validToken, err
		}
		return BannerEntity{}, validToken, nil
	}

	bannerParams := storage.GetUserBannerCriteria{
		TagId:     bannerFilter.TagId,
		FeatureId: bannerFilter.FeatureId,
	}

	banner, err := b.store.Get(bannerParams)
	if err != nil {
		return BannerEntity{}, validToken, err
	}

	bannerEntity, err := b.createBannerEntity(banner)
	if err != nil {
		b.logger.WithFields(logrus.Fields{
			"package":  "banner_service",
			"function": "SearchBanner",
			"error":    err,
		}).Error("Error unmarshalling content in createBannerEntity")
		return BannerEntity{}, validToken, err
	}

	bannerContent, err := json.Marshal(bannerEntity[0].Content)
	if err != nil {
		b.logger.WithFields(logrus.Fields{
			"package":  "banner_service",
			"function": "SearchBanner",
			"error":    err,
		}).Error("Error marshalling content")
		return BannerEntity{}, validToken, err
	}
	err = b.redis.AddToCache(cacheKey, redis.RedisEntity{Content: string(bannerContent)})
	if err != nil {
		return BannerEntity{}, validToken, err
	}

	return bannerEntity[0], validToken, nil
}

func (b *BannerService) SearchAllBanners(bannerFilter GetAllFilter, user User) ([]BannerEntity, bool, error) {
	validToken, err := b.store.IfAdminTokenValid(user.Token)
	if err != nil {
		return []BannerEntity{}, validToken, err
	}
	if !validToken {
		return []BannerEntity{}, validToken, nil
	}

	bannerRequest := storage.GetBannersListCriteria{
		TagId:     bannerFilter.TagId,
		FeatureId: bannerFilter.FeatureId,
		Limit:     bannerFilter.Limit,
		Offset:    bannerFilter.Offset,
	}

	banner, err := b.store.GetAllBanners(bannerRequest)
	if err != nil {
		return []BannerEntity{}, validToken, err
	}

	resultBanners, err := b.createBannerEntity(banner)
	if err != nil {
		return []BannerEntity{}, validToken, err
	}

	return resultBanners, validToken, nil
}

func (b *BannerService) PostBanner(newPostBanner BannerEntity, user User) (int, bool, error) {
	accessPermited, err := b.store.IfAdminTokenValid(user.Token)
	if err != nil {
		return 0, accessPermited, err
	}
	if !accessPermited {
		return 0, accessPermited, nil
	}

	someString, err := json.Marshal(newPostBanner.Content)
	if err != nil {
		b.logger.WithFields(logrus.Fields{
			"package":  "banner_service",
			"function": "PostBanner",
			"error":    err,
		}).Error("Error marshalling content")
		return 0, accessPermited, err
	}
	postBanner := storage.BannerPutPostCriteria{
		TagIds:    newPostBanner.TagId,
		FeatureId: newPostBanner.FeatureId,
		Content:   string(someString[:]),
		IsActive:  *newPostBanner.IsActive,
		CreatedAt: time.Now(),
		UpdatedAt: newPostBanner.UpdatedAt,
	}

	createdID, err := b.store.PostBanner(postBanner)
	if err != nil {
		return 0, accessPermited, err
	}
	return createdID, accessPermited, nil
}

func (b *BannerService) PutBanner(newPutBanner BannerEntity, user User) (bool, bool, error) {
	accessPermited, err := b.store.IfAdminTokenValid(user.Token)
	if err != nil {
		return false, accessPermited, err
	}
	if !accessPermited {
		return false, accessPermited, nil
	}

	bannerExists, err := b.store.SearchBannerByID(newPutBanner.ID)
	if err != nil {
		return false, accessPermited, err
	}
	if !bannerExists {
		return bannerExists, accessPermited, nil
	}

	someString, err := json.Marshal(newPutBanner.Content)
	if err != nil {
		b.logger.WithFields(logrus.Fields{
			"package":  "banner_service",
			"function": "PutBanner",
			"error":    err,
		}).Error("Error marshalling content")
		return bannerExists, accessPermited, nil
	}
	putBanner := storage.BannerPutPostCriteria{
		TagIds:    newPutBanner.TagId,
		FeatureId: newPutBanner.FeatureId,
		Content:   string(someString[:]),
		ID:        newPutBanner.ID,
	}

	if newPutBanner.IsActive != nil {
		putBanner.IfFlagActiveIsSet = *newPutBanner.IsActive
	}

	err = b.store.PutBanner(putBanner)
	return bannerExists, accessPermited, err
}

func (b *BannerService) DeleteBanner(newDeleteBanner BannerEntity, user User) (bool, bool, error) {

	accessPermited, err := b.store.IfAdminTokenValid(user.Token)
	if err != nil {
		return false, accessPermited, err
	}
	if !accessPermited {
		return false, accessPermited, nil
	}

	bannerExists, err := b.store.SearchBannerByID(newDeleteBanner.ID)
	if err != nil {
		return false, accessPermited, err
	}
	if !bannerExists {
		return bannerExists, accessPermited, nil
	}

	deleteBanner := storage.BannerPutPostCriteria{
		ID: newDeleteBanner.ID,
	}

	err = b.store.DeleteBanner(deleteBanner)

	return bannerExists, accessPermited, err
}

func (b *BannerService) createBannerEntity(bannerList []storage.BannerEntitySql) ([]BannerEntity, error) {

	var bannerEntityList []BannerEntity
	for _, b := range bannerList {
		var bannerEntity BannerEntity
		isActive := b.IsActive
		bannerEntity.ID = b.ID
		bannerEntity.TagId = b.TagId
		bannerEntity.FeatureId = b.FeatureId
		bannerEntity.IsActive = &isActive
		bannerEntity.CreatedAt = b.CreatedAt
		bannerEntity.UpdatedAt = b.UpdatedAt

		err := json.Unmarshal([]byte(b.Content), &bannerEntity.Content)

		if err != nil {
			return bannerEntityList, fmt.Errorf("500")
		}
		bannerEntityList = append(bannerEntityList, bannerEntity)

	}

	return bannerEntityList, nil
}

func (b *BannerService) createCacheKey(tagID int, featureID int) string {
	return fmt.Sprintf("tag_id=%d&feature_id=%d", tagID, featureID)
}
