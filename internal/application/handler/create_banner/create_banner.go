//go:generate mockgen -source=post_banner.go -destination=mocks/mock.go
package createbanner

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CyberPiess/banner_sevice/internal/domain/banner"
)

type postBannerService interface {
	SearchBanner(bannerFilter banner.GetFilter, user banner.User) (banner.BannerEntity, bool, error)
	SearchAllBanners(bannerFilter banner.GetAllFilter, user banner.User) ([]banner.BannerEntity, bool, error)
	PostBanner(newPostBanner banner.BannerEntity, user banner.User) (int, bool, error)
	PutBanner(newPutBanner banner.BannerEntity, user banner.User) (bool, bool, error)
	DeleteBanner(newPutBanner banner.BannerEntity, user banner.User) (bool, bool, error)
}

type PostBanner struct {
	service postBannerService
}

func NewPostBannerHandler(service postBannerService) *PostBanner {
	return &PostBanner{service: service}
}

func (pb *PostBanner) PostBanner(w http.ResponseWriter, r *http.Request) {
	var dataFromBody CreateDTO
	err := json.NewDecoder(r.Body).Decode(&dataFromBody)
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

	err = pb.verifyData(dataFromBody)
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

	postBanner := createEntityFromDTO(dataFromBody)

	createdID, accessPermited, err := pb.service.PostBanner(postBanner, user)
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
	bannerID := struct {
		ID int `json:"banner_id"`
	}{
		ID: createdID,
	}
	jsonContent, err := json.Marshal(bannerID)
	if err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonContent)
}

func (pb *PostBanner) verifyData(postBanner CreateDTO) error {
	if len(postBanner.TagId) == 0 {
		return fmt.Errorf("tag_ids is empty")
	}
	if postBanner.FeatureId == 0 {
		return fmt.Errorf("feature_id is empty")
	}
	if len(postBanner.Content) == 0 {
		return fmt.Errorf("content is empty")
	}
	if postBanner.IsActive == nil {
		return fmt.Errorf("is_active is empty")
	}
	return nil
}
