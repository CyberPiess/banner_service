package adminbanner

import (
	"context"

	"github.com/CyberPiess/banner_sevice/internal/domain/banner"
)

type bannerService interface {
	SearchBanner(ctx context.Context, bannerFilter banner.Filter, user banner.User) (banner.BannerEntity, bool, error)
}
