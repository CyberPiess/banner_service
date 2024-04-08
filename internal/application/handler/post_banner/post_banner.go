//go:generate mockgen -source=post_banner.go -destination=mocks/mock.go
package postbanner

import (
	"encoding/json"
	"net/http"

	"github.com/CyberPiess/banner_sevice/internal/domain/banner"
)

type postBannerService interface {
	SearchBanner(bannerFilter banner.Filter, user banner.User) (banner.BannerEntity, bool, error)
	SearchAllBanners(bannerFilter banner.Filter, user banner.User) ([]banner.BannerEntity, bool, error)
	PostBanner(newBanner banner.BannerEntity, user banner.User) (int64, bool, error)
}

type PostBanner struct {
	service postBannerService
}

func NewPostBannerHandler(service postBannerService) *PostBanner {
	return &PostBanner{service: service}
}

func (pb *PostBanner) PostBanner(w http.ResponseWriter, r *http.Request) {
	var postBanner banner.BannerEntity
	err := json.NewDecoder(r.Body).Decode(&postBanner)
	if err != nil {
		http.Error(w, "некорректные данные", http.StatusBadRequest)
		return
	}
	user := banner.User{
		Token: r.Header.Get("token"),
	}

	createdID, accessPermited, err := pb.service.PostBanner(postBanner, user)
	switch {
	case err != nil && err.Error() == "unauthorized user":
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	case err != nil:
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	case !accessPermited:
		http.Error(w, "Пользователь не имеет доступа", http.StatusForbidden)
		return
	default:
	}
	bannerID := struct {
		ID int64 `json:"banner_id"`
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
