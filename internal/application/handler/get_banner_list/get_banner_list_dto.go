package getbannerlist

import (
	"encoding/json"
	"time"

	"github.com/CyberPiess/banner_sevice/internal/domain/banner"
)

type GetAllBannersDTO struct {
	TagId     int `schema:"tag_id"`
	FeatureId int `schema:"feature_id"`
	Limit     int `schema:"limit"`
	Offset    int `schema:"offset"`
}

func createAllBannersFilterFromDTO(dataFromQuery GetAllBannersDTO) banner.GetAllFilter {
	return banner.GetAllFilter{
		TagId:     dataFromQuery.TagId,
		FeatureId: dataFromQuery.FeatureId,
		Limit:     dataFromQuery.Limit,
		Offset:    dataFromQuery.Offset,
	}
}

type ErrorBody struct {
	Error string `json:"error"`
}

func createFromEntity(entityList []banner.BannerEntity) ([]byte, error) {
	type response struct {
		ID        int                    `json:"banner_id"`
		TagId     []int                  `json:"tag_ids"`
		FeatureId int                    `json:"feature_id"`
		Content   map[string]interface{} `json:"content"`
		IsActive  bool                   `json:"is_active"`
		CreatedAt time.Time              `json:"created_at"`
		UpdatedAt time.Time              `json:"updated_at"`
	}

	var result []response
	for _, entity := range entityList {
		partOfSlice := response{ID: entity.ID,
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
