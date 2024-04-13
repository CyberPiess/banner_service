//go:generate mockgen -source=get_banner.go -destination=mocks/mock.go
package get_banner_list

import (
	"encoding/json"
	"net/http"

	bannerService "github.com/CyberPiess/banner_service/internal/domain/banner"
	"github.com/gorilla/schema"
	"github.com/sirupsen/logrus"
)

type getAllBannerService interface {
	SearchBanner(bannerFilter bannerService.GetFilter, user bannerService.User) (bannerService.BannerEntity, bool, error)
	SearchAllBanners(bannerFilter bannerService.GetAllFilter, user bannerService.User) ([]bannerService.BannerEntity, bool, error)
	PostBanner(newPostBanner bannerService.BannerEntity, user bannerService.User) (int, bool, error)
	PutBanner(newPutBanner bannerService.BannerEntity, user bannerService.User) (bool, bool, error)
	DeleteBanner(newDeleteBanner bannerService.BannerEntity, user bannerService.User) (bool, bool, error)
}

type logger interface {
	WithFields(fields logrus.Fields) *logrus.Entry
}

type Banner struct {
	service getAllBannerService
	logger  logger
}

func NewGetAllBannersHandler(service getAllBannerService, logger logger) *Banner {
	return &Banner{service: service, logger: logger}
}

func (b *Banner) GetAllBanners(w http.ResponseWriter, r *http.Request) {

	var decoder = schema.NewDecoder()

	err := r.ParseForm()
	if err != nil {
		b.logger.WithFields(logrus.Fields{
			"package":  "get_banner_list",
			"function": "GetAllBanners",
			"error":    err,
		}).Warn("Error parsing form")

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
		b.logger.WithFields(logrus.Fields{
			"package":  "get_banner_list",
			"function": "GetAllBanners",
			"error":    err,
		}).Warn("Error decoding form")

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
		b.logger.WithFields(logrus.Fields{
			"package":  "get_banner_list",
			"function": "GetAllBanners",
			"error":    err,
		}).Error("Error creating from entity")

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
