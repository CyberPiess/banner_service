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

type ResponseBody struct {
	Content string `json:"content"`
}

func createFilterFromDTO(dataFromSchema GetUserBannerDTO) banner.GetFilter {
	return banner.GetFilter{
		TagId:           dataFromSchema.TagId,
		FeatureId:       dataFromSchema.FeatureId,
		UseLastRevision: dataFromSchema.UseLastRevision}
}

type responseBody struct {
	Content map[string]interface{} `json:"content"`
}

func createFromEntity(entity banner.BannerEntity) ([]byte, error) {
	responseBody := responseBody{
		Content: entity.Content,
	}
	jsonContent, err := json.Marshal(responseBody.Content)
	if err != nil {
		return nil, err
	}
	return jsonContent, nil
}
