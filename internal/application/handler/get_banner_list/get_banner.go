//go:generate mockgen -source=get_banner.go -destination=mocks/mock.go
package getbannerlist

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
	DeleteBanner(newDeleteBanner banner.BannerEntity, user banner.User) (bool, bool, error)
}

type Banner struct {
	service bannerService
}

func NewGetAllBannersHandler(service bannerService) *Banner {
	return &Banner{service: service}
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

	var dataFromQuery GetAllBannersDTO
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

	allBannersFilter := createAllBannersFilterFromDTO(dataFromQuery)

	foundBanners, accessPermited, err := b.service.SearchAllBanners(allBannersFilter, user)
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

	jsonContent, err := createFromEntity(foundBanners)
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
