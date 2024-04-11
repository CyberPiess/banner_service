//go:generate mockgen -source=user_banner.go -destination=mocks/mock.go
package userbanner

import (
	"encoding/json"
	"net/http"

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

func (b *Banner) GetUserBanner(w http.ResponseWriter, r *http.Request) {
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

	var dataFromQuery GetUserBannerDTO
	err = decoder.Decode(&dataFromQuery, r.Form)
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

	bannerFilter := createFilterFromDTO(dataFromQuery)

	newBanner, accessPermited, err := b.service.SearchBanner(bannerFilter, user)
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
	case newBanner.Content == nil:
		w.WriteHeader(http.StatusNotFound)
		return
	default:
	}

	jsonContent, err := createFromEntity(newBanner)
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
