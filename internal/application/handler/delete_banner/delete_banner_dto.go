package deletebanner

import "github.com/CyberPiess/banner_sevice/internal/domain/banner"

type BannerDeleteDTO struct {
	ID int
}

type ErrorBody struct {
	Error string `json:"error"`
}

func createEntityFromDTO(dataFromPath BannerDeleteDTO) banner.BannerEntity {
	return banner.BannerEntity{
		ID: dataFromPath.ID,
	}
}
