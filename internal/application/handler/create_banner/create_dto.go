package create_banner

import bannerService "github.com/CyberPiess/banner_service/internal/domain/banner"

type CreateDTO struct {
	Content   map[string]interface{} `json:"content"`
	TagId     []int                  `json:"tag_ids"`
	FeatureId int                    `json:"feature_id"`
	IsActive  *bool                  `json:"is_active"`
}

type ErrorBody struct {
	Error string `json:"error"`
}

func createEntityFromDTO(dataFromBody CreateDTO) bannerService.BannerEntity {
	return bannerService.BannerEntity{
		Content:   dataFromBody.Content,
		FeatureId: dataFromBody.FeatureId,
		TagId:     dataFromBody.TagId,
		IsActive:  dataFromBody.IsActive,
	}
}
