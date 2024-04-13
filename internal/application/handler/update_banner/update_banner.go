//go:generate mockgen -source=update_banner.go -destination=mocks/mock.go
package update_banner

import (
	"encoding/json"
	"net/http"
	"strconv"

	bannerService "github.com/CyberPiess/banner_service/internal/domain/banner"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type putBannerService interface {
	SearchBanner(bannerFilter bannerService.GetFilter, user bannerService.User) (bannerService.BannerEntity, bool, error)
	SearchAllBanners(bannerFilter bannerService.GetAllFilter, user bannerService.User) ([]bannerService.BannerEntity, bool, error)
	PostBanner(newPostBanner bannerService.BannerEntity, user bannerService.User) (int, bool, error)
	PutBanner(newPutBanner bannerService.BannerEntity, user bannerService.User) (bool, bool, error)
	DeleteBanner(newPutBanner bannerService.BannerEntity, user bannerService.User) (bool, bool, error)
}

type logger interface {
	WithFields(fields logrus.Fields) *logrus.Entry
}

type PutBanner struct {
	service putBannerService
	logger  logger
}

func NewPutBannerHandler(service putBannerService, logger logger) *PutBanner {
	return &PutBanner{service: service, logger: logger}
}

func (ptB *PutBanner) PutBanner(w http.ResponseWriter, r *http.Request) {
	var dataFromPath UpdateBannerDTO
	var err error
	bannerID := mux.Vars(r)["id"]
	dataFromPath.ID, err = strconv.Atoi(bannerID)
	if err != nil {
		ptB.logger.WithFields(logrus.Fields{
			"package":  "update_banner",
			"function": "PutBanner",
			"error":    err,
		}).Warn("Error parsing path")

		response := ErrorBody{
			Error: err.Error(),
		}
		responseBody, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseBody)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&dataFromPath)
	if err != nil {
		ptB.logger.WithFields(logrus.Fields{
			"package":  "update_banner",
			"function": "PutBanner",
			"error":    err,
		}).Warn("Error decoding")

		response := ErrorBody{
			Error: err.Error(),
		}
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

	updateBanner := createEntityFromDTO(dataFromPath)
	found, accessPermited, err := ptB.service.PutBanner(updateBanner, user)
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
