//go:generate mockgen -source=put_banner.go -destination=mocks/mock.go
package updatebanner

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/CyberPiess/banner_sevice/internal/domain/banner"
	"github.com/gorilla/mux"
)

type putBannerService interface {
	SearchBanner(bannerFilter banner.GetFilter, user banner.User) (banner.BannerEntity, bool, error)
	SearchAllBanners(bannerFilter banner.GetAllFilter, user banner.User) ([]banner.BannerEntity, bool, error)
	PostBanner(newPostBanner banner.BannerEntity, user banner.User) (int, bool, error)
	PutBanner(newPutBanner banner.BannerEntity, user banner.User) (bool, bool, error)
	DeleteBanner(newPutBanner banner.BannerEntity, user banner.User) (bool, bool, error)
}

type PutBanner struct {
	service putBannerService
}

func NewPutBannerHandler(service putBannerService) *PutBanner {
	return &PutBanner{service: service}
}

type ErrorBody struct {
	Error string `json:"error"`
}

func (ptB *PutBanner) PutBanner(w http.ResponseWriter, r *http.Request) {
	var putBanner banner.BannerEntity
	var err error
	bannerID := mux.Vars(r)["id"]
	putBanner.ID, err = strconv.Atoi(bannerID)
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

	err = json.NewDecoder(r.Body).Decode(&putBanner)
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

	found, accessPermited, err := ptB.service.PutBanner(putBanner, user)
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
	case !found:
		w.WriteHeader(http.StatusNotFound)
		return
	default:
	}

	w.WriteHeader(http.StatusOK)
}
