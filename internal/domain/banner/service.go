//go:generate mockgen -source=service.go -destination=mocks/mock.go
package banner

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/CyberPiess/banner_sevice/internal/infrastructure/postgres/banner"
	redis "github.com/CyberPiess/banner_sevice/internal/infrastructure/redis/cache"
)

type bannerStorage interface {
	Get(bannerParams banner.GetUserBannerCriteria) ([]banner.BannerEntitySql, error)
	IfTokenValid(token string) (bool, error)
	IfBannerExists(featureId int, tagId int) (bool, error)
	IfAdminTokenValid(token string) (bool, error)
	SearchBannerByID(bannerID int) (bool, error)
	GetAllBanners(bannerParams banner.GetBannersListCriteria) ([]banner.BannerEntitySql, error)
	PostBanner(postBannerParams banner.BannerPutPostCriteria) (int, error)
	PutBanner(putBannerParams banner.BannerPutPostCriteria) error
	DeleteBanner(deleteBannerParams banner.BannerPutPostCriteria) error
}

type redisCache interface {
	AddToCache(key string, redisDTO redis.RedisEntity) error
	GetFromCache(key string) (redis.RedisEntity, error)
	DeleteFromCache(key string) error
}

type BannerService struct {
	store bannerStorage
	redis redisCache
}

func NewBannerService(storage bannerStorage, redis redisCache) *BannerService {
	return &BannerService{store: storage,
		redis: redis}
}

func (b *BannerService) SearchBanner(bannerFilter GetFilter, user User) (BannerEntity, bool, error) {
	validToken, err := b.store.IfTokenValid(user.Token)
	if err != nil {
		return BannerEntity{}, validToken, err
	}
	if !validToken {
		return BannerEntity{}, validToken, nil
	}
	cacheKey := fmt.Sprintf("tag_id=%d&feature_id=%d", bannerFilter.TagId, bannerFilter.FeatureId)

	if !bannerFilter.UseLastRevision {
		foundInCahce, err := b.redis.GetFromCache(cacheKey)
		if err != nil {
			log.Printf("error redis: %s", err.Error())
		}
		if foundInCahce.Content != "" {
			var bannerEntity BannerEntity
			err := json.Unmarshal([]byte(foundInCahce.Content), &bannerEntity.Content)
			if err != nil {
				return BannerEntity{}, validToken, err
			}
			return bannerEntity, validToken, err
		}
	}

	bannerExists, err := b.store.IfBannerExists(bannerFilter.TagId, bannerFilter.FeatureId)
	if err != nil {
		return BannerEntity{}, validToken, err
	}
	if !bannerExists {
		return BannerEntity{}, validToken, nil
	}

	bannerParams := banner.GetUserBannerCriteria{
		TagId:           bannerFilter.TagId,
		FeatureId:       bannerFilter.FeatureId,
		UseLastRevision: bannerFilter.UseLastRevision,
	}

	banner, err := b.store.Get(bannerParams)
	if err != nil {
		return BannerEntity{}, validToken, err
	}

	bannerEntity, err := b.createBannerEntity(banner)
	if err != nil {
		return BannerEntity{}, validToken, err
	}

	bannerContent, err := json.Marshal(bannerEntity[0].Content)
	if err != nil {
		log.Printf("error redis: %s", err.Error())
	}
	err = b.redis.AddToCache(cacheKey, redis.RedisEntity{Content: string(bannerContent)})
	if err != nil {
		log.Printf("error redis: %s", err.Error())
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

	bannerRequest := banner.GetBannersListCriteria{
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
		return 0, accessPermited, err
	}
	postBanner := banner.BannerPutPostCriteria{
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
		return bannerExists, accessPermited, nil
	}
	putBanner := banner.BannerPutPostCriteria{
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

	deleteBanner := banner.BannerPutPostCriteria{
		ID: newDeleteBanner.ID,
	}

	err = b.store.DeleteBanner(deleteBanner)

	return bannerExists, accessPermited, err
}

func (b *BannerService) createBannerEntity(bannerList []banner.BannerEntitySql) ([]BannerEntity, error) {

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
