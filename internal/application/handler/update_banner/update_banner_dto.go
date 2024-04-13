package update_banner

import bannerService "github.com/CyberPiess/banner_service/internal/domain/banner"

type UpdateBannerDTO struct {
	ID int
}

type ErrorBody struct {
	Error string `json:"error"`
}

func createEntityFromDTO(dataFromPath UpdateBannerDTO) bannerService.BannerEntity {
	return bannerService.BannerEntity{ID: dataFromPath.ID}
}
