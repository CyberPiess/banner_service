package banner

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CyberPiess/banner_sevice/internal/infrastructure/postgres/banner"
)

type bannerStorage interface {
	Get(ctx context.Context, bannerRequest banner.BannerRequest) (banner.BannerResponse, error)
}

type BannerService struct {
	store bannerStorage
}

func NewBannerService(storage bannerStorage) *BannerService {
	return &BannerService{store: storage}
}

func (b *BannerService) SearchBanner(ctx context.Context, bannerFilter Filter, user User) (BannerEntity, error) {
	//TODO: добавить проверки

	bannerRequest := banner.BannerRequest{
		TagId:           bannerFilter.TagId,
		FeatureId:       bannerFilter.FeatureId,
		UseLastRevision: bannerFilter.UseLastRevision,
	}

	banner, err := b.store.Get(ctx, bannerRequest)
	if err != nil {
		return BannerEntity{}, err
	}

	bannerEntity, err := b.createBannerEntity(banner)
	if err != nil {
		return BannerEntity{}, err
	}

	return bannerEntity, nil
}

func (b *BannerService) createBannerEntity(banner banner.BannerResponse) (BannerEntity, error) {

	var bannerEntity BannerEntity
	err := json.Unmarshal([]byte(banner.Content), &bannerEntity.Content)
	if err != nil {
		return bannerEntity, fmt.Errorf("500")
	}

	return bannerEntity, nil
}
