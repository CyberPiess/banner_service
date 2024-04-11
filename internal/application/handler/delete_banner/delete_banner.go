//go:generate mockgen -source=delete_banner.go -destination=mocks/mock.go
package deletebanner

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/CyberPiess/banner_sevice/internal/domain/banner"
	"github.com/gorilla/mux"
)

type deleteBannerService interface {
	SearchBanner(bannerFilter banner.GetFilter, user banner.User) (banner.BannerEntity, bool, error)
	SearchAllBanners(bannerFilter banner.GetAllFilter, user banner.User) ([]banner.BannerEntity, bool, error)
	PostBanner(newPostBanner banner.BannerEntity, user banner.User) (int, bool, error)
	PutBanner(newPutBanner banner.BannerEntity, user banner.User) (bool, bool, error)
	DeleteBanner(newDeleteBanner banner.BannerEntity, user banner.User) (bool, bool, error)
}

type DeleteBanner struct {
	service deleteBannerService
}

func NewDeleteBannerHandler(service deleteBannerService) *DeleteBanner {
	return &DeleteBanner{service: service}
}

func (dB *DeleteBanner) DeleteBanner(w http.ResponseWriter, r *http.Request) {
	var bannerFromPath BannerDeleteDTO
	var err error
	var response ErrorBody
	bannerID := mux.Vars(r)["id"]
	bannerFromPath.ID, err = strconv.Atoi(bannerID)
	if err != nil {
		response.Error = err.Error()
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

	bannerForDelete := createEntityFromDTO(bannerFromPath)

	found, accessPermited, err := dB.service.DeleteBanner(bannerForDelete, user)
	switch {
	case err != nil:
		response.Error = err.Error()
		responseBody, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseBody)
		return
	case !accessPermited:
		w.WriteHeader(http.StatusForbidden)
		return
	case !found:
		w.WriteHeader(http.StatusNotFound)
		return
	default:
	}

	w.WriteHeader(http.StatusNoContent)
}
