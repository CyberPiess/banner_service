package createbanner

import "github.com/CyberPiess/banner_sevice/internal/domain/banner"

type CreateDTO struct {
	Content   map[string]interface{} `json:"content"`
	TagId     []int                  `json:"tag_ids"`
	FeatureId int                    `json:"feature_id"`
	IsActive  *bool                  `json:"is_active"`
}

type ErrorBody struct {
	Error string `json:"error"`
}

func createEntityFromDTO(dataFromBody CreateDTO) banner.BannerEntity {
	return banner.BannerEntity{
		Content:   dataFromBody.Content,
		FeatureId: dataFromBody.FeatureId,
		TagId:     dataFromBody.TagId,
		IsActive:  dataFromBody.IsActive,
	}
}
