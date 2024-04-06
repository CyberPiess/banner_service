package banner

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CyberPiess/banner_sevice/internal/infrastructure/postgres/banner"
)

type bannerStorage interface {
	Get(ctx context.Context, bannerRequest banner.BannerRequest) (banner.BannerResponse, error)
	IfTokenValid(token string) (bool, error)
}
type BannerService struct {
	store bannerStorage
}

func NewBannerService(storage bannerStorage) *BannerService {
	return &BannerService{store: storage}
}

func (b *BannerService) SearchBanner(ctx context.Context, bannerFilter Filter, user User) (BannerEntity, bool, error) {
	validToken, err := b.store.IfTokenValid(user.Token)
	if err != nil {
		return BannerEntity{}, validToken, err
	}

	if !validToken {
		return BannerEntity{}, validToken, nil
	}

	bannerRequest := banner.BannerRequest{
		TagId:           bannerFilter.TagId,
		FeatureId:       bannerFilter.FeatureId,
		UseLastRevision: bannerFilter.UseLastRevision,
	}

	banner, err := b.store.Get(ctx, bannerRequest)
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
