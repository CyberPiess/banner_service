package get_banner_list

import (
	"encoding/json"
	"time"

	bannerService "github.com/CyberPiess/banner_service/internal/domain/banner"
)

type GetAllBannersDTO struct {
	TagId     int `schema:"tag_id"`
	FeatureId int `schema:"feature_id"`
	Limit     int `schema:"limit"`
	Offset    int `schema:"offset"`
}

func createAllBannersFilterFromDTO(dataFromQuery GetAllBannersDTO) bannerService.GetAllFilter {
	return bannerService.GetAllFilter{
		TagId:     dataFromQuery.TagId,
		FeatureId: dataFromQuery.FeatureId,
		Limit:     dataFromQuery.Limit,
		Offset:    dataFromQuery.Offset,
	}
}

type ErrorBody struct {
	Error string `json:"error"`
}

type responseBody struct {
	ID        int                    `json:"banner_id"`
	TagId     []int                  `json:"tag_ids"`
	FeatureId int                    `json:"feature_id"`
	Content   map[string]interface{} `json:"content"`
	IsActive  bool                   `json:"is_active"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

func createFromEntity(entityList []bannerService.BannerEntity) ([]byte, error) {

	var result []responseBody
	for _, entity := range entityList {
		partOfSlice := responseBody{ID: entity.ID,
			Content:   entity.Content,
			TagId:     entity.TagId,
			FeatureId: entity.FeatureId,
			IsActive:  *entity.IsActive,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt}
		result = append(result, partOfSlice)
	}

	if len(result) == 0 {
		return []byte("[]"), nil
	}

	jsonContent, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	return jsonContent, nil
}
