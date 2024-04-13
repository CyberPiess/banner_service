package delete_banner

import banner_service "github.com/CyberPiess/banner_service/internal/domain/banner"

type BannerDeleteDTO struct {
	ID int
}

type ErrorBody struct {
	Error string `json:"error"`
}

func createEntityFromDTO(dataFromPath BannerDeleteDTO) banner_service.BannerEntity {
	return banner_service.BannerEntity{
		ID: dataFromPath.ID,
	}
}
