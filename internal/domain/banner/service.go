package banner

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/CyberPiess/banner_sevice/internal/infrastructure/postgres/banner"
)

type bannerStorage interface {
	Get(bannerRequest banner.BannerRequest) (banner.BannerResponse, error)
	IfTokenValid(token string) (bool, error)
	IfBannerExists(featureId int, tagId int) (bool, error)
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

	return bannerEntity, validToken, nil
}

func (b *BannerService) createBannerEntity(banner banner.BannerResponse) (BannerEntity, error) {

	var bannerEntity BannerEntity
	err := json.Unmarshal([]byte(banner.Content), &bannerEntity.Content)
	if err != nil {
		return bannerEntity, fmt.Errorf("500")
	}

	return bannerEntity, nil
}

func (b *BannerService) verifyData(bannerFilter Filter) (banner.BannerRequest, error) {
	var bannerRequest banner.BannerRequest
	var err error

	bannerRequest.TagId, err = strconv.Atoi(bannerFilter.TagId)
	if err != nil {
		return banner.BannerRequest{}, err
	}

	bannerRequest.FeatureId, err = strconv.Atoi(bannerFilter.FeatureId)
	if err != nil {
		return banner.BannerRequest{}, err
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
