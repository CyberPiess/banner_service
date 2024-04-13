//go:generate mockgen -source=delete_banner.go -destination=mocks/mock.go
package delete_banner

import (
	"encoding/json"
	"net/http"
	"strconv"

	bannerService "github.com/CyberPiess/banner_service/internal/domain/banner"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type deleteBannerService interface {
	SearchBanner(bannerFilter bannerService.GetFilter, user bannerService.User) (bannerService.BannerEntity, bool, error)
	SearchAllBanners(bannerFilter bannerService.GetAllFilter, user bannerService.User) ([]bannerService.BannerEntity, bool, error)
	PostBanner(newPostBanner bannerService.BannerEntity, user bannerService.User) (int, bool, error)
	PutBanner(newPutBanner bannerService.BannerEntity, user bannerService.User) (bool, bool, error)
	DeleteBanner(newDeleteBanner bannerService.BannerEntity, user bannerService.User) (bool, bool, error)
}

type logger interface {
	WithFields(fields logrus.Fields) *logrus.Entry
}

type DeleteBanner struct {
	service deleteBannerService
	logger  logger
}

func NewDeleteBannerHandler(service deleteBannerService, logger logger) *DeleteBanner {
	return &DeleteBanner{service: service, logger: logger}
}

func (dB *DeleteBanner) DeleteBanner(w http.ResponseWriter, r *http.Request) {
	var bannerFromPath BannerDeleteDTO
	var err error
	var response ErrorBody
	bannerID := mux.Vars(r)["id"]
	bannerFromPath.ID, err = strconv.Atoi(bannerID)
	if err != nil {
		dB.logger.WithFields(logrus.Fields{
			"package":  "create_banner",
			"function": "DeleteBanner",
			"error":    err,
		}).Warn("Error decoding path")

		response.Error = err.Error()
		responseBody, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseBody)
		return
	}
	user := bannerService.User{
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
