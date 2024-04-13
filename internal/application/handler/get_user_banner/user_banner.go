//go:generate mockgen -source=user_banner.go -destination=mocks/mock.go
package user_banner

import (
	"encoding/json"
	"net/http"

	bannerService "github.com/CyberPiess/banner_service/internal/domain/banner"
	"github.com/gorilla/schema"
	"github.com/sirupsen/logrus"
)

type getUserBannerService interface {
	SearchBanner(bannerFilter bannerService.GetFilter, user bannerService.User) (bannerService.BannerEntity, bool, error)
	SearchAllBanners(bannerFilter bannerService.GetAllFilter, user bannerService.User) ([]bannerService.BannerEntity, bool, error)
	PostBanner(newPostBanner bannerService.BannerEntity, user bannerService.User) (int, bool, error)
	PutBanner(newPutBanner bannerService.BannerEntity, user bannerService.User) (bool, bool, error)
	DeleteBanner(newPutBanner bannerService.BannerEntity, user bannerService.User) (bool, bool, error)
}

type logger interface {
	WithFields(fields logrus.Fields) *logrus.Entry
}

type Banner struct {
	service getUserBannerService
	logger  logger
}

func NewBannerHandler(service getUserBannerService, logger logger) *Banner {
	return &Banner{service: service, logger: logger}
}

func (b *Banner) GetUserBanner(w http.ResponseWriter, r *http.Request) {
	var decoder = schema.NewDecoder()

	err := r.ParseForm()
	if err != nil {
		b.logger.WithFields(logrus.Fields{
			"package":  "user_banner",
			"function": "GetUserBanners",
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

	var dataFromQuery GetUserBannerDTO
	err = decoder.Decode(&dataFromQuery, r.Form)
	if err != nil {
		b.logger.WithFields(logrus.Fields{
			"package":  "user_banner",
			"function": "GetUserBanners",
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
		b.logger.WithFields(logrus.Fields{
			"package":  "user_banner",
			"function": "GetUserBanners",
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
