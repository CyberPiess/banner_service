package banner

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/CyberPiess/banner_sevice/internal/infrastructure/postgres/banner"
)

type bannerStorage interface {
	Get(bannerRequest banner.BannerRequest) ([]banner.BannerResponse, error)
	IfTokenValid(token string) (bool, error)
	IfBannerExists(featureId int, tagId int) (bool, error)
	IfAdminTokenValid(token string) (bool, error)
	GetAllBanners(bannerParams banner.BannerRequest) ([]banner.BannerResponse, error)
	PostBanner(postBannerParams banner.BannerPostRequest) (int64, error)
}
type BannerService struct {
	store bannerStorage
}

func NewBannerService(storage bannerStorage) *BannerService {
	return &BannerService{store: storage}
}

func (b *BannerService) SearchBanner(bannerFilter Filter, user User) (BannerEntity, bool, error) {
	if !b.ifTokenSupplied(user) {
		return BannerEntity{}, false, fmt.Errorf("unauthorized user")
	}

	validToken, err := b.store.IfTokenValid(user.Token)
	if err != nil {
		return BannerEntity{}, validToken, err
	}
	if !validToken {
		return BannerEntity{}, validToken, nil
	}

	bannerRequest, err := b.verifyData(bannerFilter)
	if err != nil {
		return BannerEntity{}, validToken, fmt.Errorf("wrong data supplied")
	} else if bannerRequest.TagId == 0 || bannerRequest.FeatureId == 0 {
		return BannerEntity{}, validToken, fmt.Errorf("wrong data supplied")
	}

	bannerExists, err := b.store.IfBannerExists(bannerRequest.TagId, bannerRequest.FeatureId)
	if err != nil {
		return BannerEntity{}, validToken, err
	}
	if !bannerExists {
		return BannerEntity{}, validToken, nil
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

func (b *BannerService) SearchAllBanners(bannerFilter Filter, user User) ([]BannerEntity, bool, error) {
	if !b.ifTokenSupplied(user) {
		return []BannerEntity{}, false, fmt.Errorf("unauthorized user")
	}

	validToken, err := b.store.IfAdminTokenValid(user.Token)
	if err != nil {
		return []BannerEntity{}, validToken, err
	}
	if !validToken {
		return []BannerEntity{}, validToken, nil
	}

	bannerRequest, err := b.verifyData(bannerFilter)
	if err != nil {
		return []BannerEntity{}, validToken, fmt.Errorf("wrong data supplied")
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

func (b *BannerService) PostBanner(newPostBanner BannerEntity, user User) (int64, bool, error) {
	if !b.ifTokenSupplied(user) {
		return 0, false, fmt.Errorf("unauthorized user")
	}
	validToken, err := b.store.IfAdminTokenValid(user.Token)
	if err != nil {
		return 0, validToken, err
	}
	if !validToken {
		return 0, validToken, nil
	}

	someString, err := json.Marshal(newPostBanner.Content)
	if err != nil {
		return 0, validToken, err
	}
	postBanner := banner.BannerPostRequest{
		TagIds:    newPostBanner.TagId,
		FeatureId: newPostBanner.FeatureId,
		Content:   string(someString[:]),
		IsActive:  newPostBanner.IsActive,
		CreatedAt: time.Now(),
		UpdatedAt: newPostBanner.UpdatedAt,
	}

	createdID, err := b.store.PostBanner(postBanner)
	if err != nil {
		return 0, validToken, err
	}
	return createdID, validToken, nil
}

func (b *BannerService) createBannerEntity(bannerList []banner.BannerResponse) ([]BannerEntity, error) {

	var bannerEntityList []BannerEntity
	for _, b := range bannerList {
		bannerEntity := BannerEntity{
			ID:        b.ID,
			TagId:     b.TagId,
			FeatureId: b.FeatureId,
			IsActive:  b.IsActive,
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

func (b *BannerService) verifyData(bannerFilter Filter) (banner.BannerRequest, error) {
	var bannerRequest banner.BannerRequest
	var err error

	if bannerFilter.TagId != "" {
		bannerRequest.TagId, err = strconv.Atoi(bannerFilter.TagId)
		if err != nil {
			return banner.BannerRequest{}, err
		}
	}

	if bannerFilter.FeatureId != "" {
		bannerRequest.FeatureId, err = strconv.Atoi(bannerFilter.FeatureId)
		if err != nil {
			return banner.BannerRequest{}, err
		}
	}

	if bannerFilter.Limit != "" {
		bannerRequest.Limit, err = strconv.Atoi(bannerFilter.Limit)
		if err != nil {
			return banner.BannerRequest{}, err
		}
	}

	if bannerFilter.Offset != "" {
		bannerRequest.Offset, err = strconv.Atoi(bannerFilter.Offset)
		if err != nil {
			return banner.BannerRequest{}, err
		}
	}

	if bannerFilter.UseLastRevision == "" {
		bannerRequest.UseLastRevision = false
		return bannerRequest, nil
	}

	bannerRequest.UseLastRevision, err = strconv.ParseBool(bannerFilter.UseLastRevision)
	if err != nil {
		return banner.BannerRequest{}, err
	}
	return bannerRequest, nil
}

func (b *BannerService) ifTokenSupplied(user User) bool {
	return user.Token != ""
}
