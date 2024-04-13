//go:generate mockgen -source=create_banner.go -destination=mocks/mock.go
package create_banner

import (
	"encoding/json"
	"fmt"
	"net/http"

	bannerService "github.com/CyberPiess/banner_service/internal/domain/banner"
	"github.com/sirupsen/logrus"
)

type postBannerService interface {
	SearchBanner(bannerFilter bannerService.GetFilter, user bannerService.User) (bannerService.BannerEntity, bool, error)
	SearchAllBanners(bannerFilter bannerService.GetAllFilter, user bannerService.User) ([]bannerService.BannerEntity, bool, error)
	PostBanner(newPostBanner bannerService.BannerEntity, user bannerService.User) (int, bool, error)
	PutBanner(newPutBanner bannerService.BannerEntity, user bannerService.User) (bool, bool, error)
	DeleteBanner(newPutBanner bannerService.BannerEntity, user bannerService.User) (bool, bool, error)
}

type logger interface {
	WithFields(fields logrus.Fields) *logrus.Entry
}

type PostBanner struct {
	service postBannerService
	logger  logger
}

func NewPostBannerHandler(service postBannerService, logger logger) *PostBanner {
	return &PostBanner{service: service, logger: logger}
}

func (pb *PostBanner) PostBanner(w http.ResponseWriter, r *http.Request) {
	var dataFromBody CreateDTO
	err := json.NewDecoder(r.Body).Decode(&dataFromBody)
	if err != nil {
		pb.logger.WithFields(logrus.Fields{
			"package":  "create_banner",
			"function": "PostBanner",
			"error":    err,
		}).Warn("Error decoding body")
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
		pb.logger.WithFields(logrus.Fields{
			"package":  "create_banner",
			"function": "PostBanner",
			"error":    err,
		}).Warn("Error verifying data")
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
		pb.logger.WithFields(logrus.Fields{
			"package":  "create_banner",
			"function": "PostBanner",
			"error":    err,
		}).Warn("Error marshalling content")
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
