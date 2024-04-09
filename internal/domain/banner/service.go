package banner

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/CyberPiess/banner_sevice/internal/infrastructure/postgres/banner"
)

type bannerStorage interface {
	Get(bannerRequest banner.BannerRequest) ([]banner.BannerResponse, error)
	IfTokenValid(token string) (bool, error)
	IfBannerExists(featureId int, tagId int) (bool, error)
	IfAdminTokenValid(token string) (bool, error)
	SearchBannerByID(bannerID int) (bool, error)
	GetAllBanners(bannerParams banner.BannerRequest) ([]banner.BannerResponse, error)
	PostBanner(postBannerParams banner.BannerPutPostRequest) (int, error)
	PutBanner(putBannerParams banner.BannerPutPostRequest) error
	DeleteBanner(deleteBannerParams banner.BannerPutPostRequest) error
}
type BannerService struct {
	store bannerStorage
}

func NewBannerService(storage bannerStorage) *BannerService {
	return &BannerService{store: storage}
}

func (b *BannerService) SearchBanner(bannerFilter GetFilter, user User) (BannerEntity, bool, error) {
	validToken, err := b.store.IfTokenValid(user.Token)
	if err != nil {
		return BannerEntity{}, validToken, err
	}
	if !validToken {
		return BannerEntity{}, validToken, nil
	}

	bannerExists, err := b.store.IfBannerExists(bannerFilter.TagId, bannerFilter.FeatureId)
	if err != nil {
		return BannerEntity{}, validToken, err
	}
	if !bannerExists {
		return BannerEntity{}, validToken, nil
	}

	bannerRequest := banner.BannerRequest{
		TagId:           bannerFilter.TagId,
		FeatureId:       bannerFilter.FeatureId,
		UseLastRevision: bannerFilter.UseLastRevision,
	}

	banner, err := b.store.Get(bannerRequest)
	if err != nil {
		return BannerEntity{}, validToken, err
	}

	bannerEntity, err := b.createBannerEntity(banner)
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

	bannerRequest := banner.BannerRequest{
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
	postBanner := banner.BannerPutPostRequest{
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
	putBanner := banner.BannerPutPostRequest{
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

	deleteBanner := banner.BannerPutPostRequest{
		ID: newDeleteBanner.ID,
	}

	err = b.store.DeleteBanner(deleteBanner)

	return bannerExists, accessPermited, err
}

func (b *BannerService) createBannerEntity(bannerList []banner.BannerResponse) ([]BannerEntity, error) {

	var bannerEntityList []BannerEntity
	for _, b := range bannerList {
		bannerEntity := BannerEntity{
			ID:        b.ID,
			TagId:     b.TagId,
			FeatureId: b.FeatureId,
			IsActive:  &b.IsActive,
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
		}
		err := json.Unmarshal([]byte(b.Content), &bannerEntity.Content)

		if err != nil {
			return bannerEntityList, fmt.Errorf("500")
		}
		bannerEntityList = append(bannerEntityList, bannerEntity)

	}

	return bannerEntityList, nil
}
