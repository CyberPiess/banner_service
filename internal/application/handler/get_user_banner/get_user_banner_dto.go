package userbanner

import (
	"encoding/json"

	"github.com/CyberPiess/banner_sevice/internal/domain/banner"
)

type GetUserBannerDTO struct {
	TagId           int  `schema:"tag_id,required"`
	FeatureId       int  `schema:"feature_id,required"`
	UseLastRevision bool `schema:"use_last_revision,default:false"`
}

type ErrorBody struct {
	Error string `json:"error"`
}

func createFilterFromDTO(dataFromSchema GetUserBannerDTO) banner.GetFilter {
	return banner.GetFilter{TagId: dataFromSchema.TagId,
		FeatureId:       dataFromSchema.FeatureId,
		UseLastRevision: dataFromSchema.UseLastRevision}
}

func createFromEntity(entity banner.BannerEntity) ([]byte, error) {
	jsonContent, err := json.Marshal(entity.Content)
	if err != nil {
		return nil, err
	}
	return jsonContent, nil
}
