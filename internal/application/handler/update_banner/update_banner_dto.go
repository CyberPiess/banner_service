package updatebanner

import "github.com/CyberPiess/banner_sevice/internal/domain/banner"

type UpdateBannerDTO struct {
	ID int
}

type ErrorBody struct {
	Error string `json:"error"`
}

func createEntityFromDTO(dataFromPath UpdateBannerDTO) banner.BannerEntity {
	return banner.BannerEntity{ID: dataFromPath.ID}
}
