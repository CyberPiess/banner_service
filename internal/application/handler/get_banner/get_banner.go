//go:generate mockgen -source=get_banner.go -destination=mocks/mock.go
package adminbanner

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/CyberPiess/banner_sevice/internal/domain/banner"
	"github.com/gorilla/schema"
)

type bannerService interface {
	SearchBanner(bannerFilter banner.GetFilter, user banner.User) (banner.BannerEntity, bool, error)
	SearchAllBanners(bannerFilter banner.GetAllFilter, user banner.User) ([]banner.BannerEntity, bool, error)
	PostBanner(newPostBanner banner.BannerEntity, user banner.User) (int, bool, error)
	PutBanner(newPutBanner banner.BannerEntity, user banner.User) (bool, bool, error)
	DeleteBanner(newPutBanner banner.BannerEntity, user banner.User) (bool, bool, error)
}

type Banner struct {
	service bannerService
}

func NewBannerHandler(service bannerService) *Banner {
	return &Banner{service: service}
}

type ErrorBody struct {
	Error string `json:"error"`
}

func (b *Banner) GetAllBanners(w http.ResponseWriter, r *http.Request) {

	var decoder = schema.NewDecoder()

	err := r.ParseForm()
	if err != nil {
		response := ErrorBody{
			Error: err.Error(),
		}
		responseBody, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseBody)
		return
	}

	var bannerFilter banner.GetAllFilter
	err = decoder.Decode(&bannerFilter, r.Form)
	if err != nil {
		response := ErrorBody{
			Error: err.Error(),
		}
		responseBody, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseBody)
		return
	}

	user := banner.User{
		Token: r.Header.Get("token"),
	}
	if user.Token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	foundBanners, accessPermited, err := b.service.SearchAllBanners(bannerFilter, user)
	switch {
	case err != nil:
		response := ErrorBody{
			Error: err.Error(),
		}
		responseBody, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseBody)
		return
	case !accessPermited:
		w.WriteHeader(http.StatusForbidden)
		return
	default:
	}

	jsonContent, err := b.createFromEntity(foundBanners)
	if err != nil {
		response := ErrorBody{
			Error: err.Error(),
		}
		responseBody, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseBody)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonContent)
}

func (b *Banner) createFromEntity(entityList []banner.BannerEntity) ([]byte, error) {
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

	jsonContent, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	return jsonContent, nil
}
