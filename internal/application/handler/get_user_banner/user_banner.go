//go:generate mockgen -source=user_banner.go -destination=mocks/mock.go
package userbanner

import (
	"encoding/json"
	"net/http"

	"github.com/CyberPiess/banner_sevice/internal/domain/banner"
)

type bannerService interface {
	SearchBanner(bannerFilter banner.Filter, user banner.User) (banner.BannerEntity, bool, error)
	SearchAllBanners(bannerFilter banner.Filter, user banner.User) ([]banner.BannerEntity, bool, error)
	PostBanner(newBanner banner.BannerEntity, user banner.User) (int64, bool, error)
}

type Banner struct {
	service bannerService
}

func NewBannerHandler(service bannerService) *Banner {
	return &Banner{service: service}
}

func (b *Banner) GetUserBanner(w http.ResponseWriter, r *http.Request) {

	bannerFilter := banner.Filter{
		TagId:           r.FormValue("tag_id"),
		FeatureId:       r.FormValue("feature_id"),
		UseLastRevision: r.FormValue("use_last_revision"),
	}

	user := banner.User{
		Token: r.Header.Get("token"),
	}

	newBanner, accessPermited, err := b.service.SearchBanner(bannerFilter, user)
	switch {
	case err != nil && err.Error() == "unauthorized user":
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	case err != nil && err.Error() == "wrong data supplied":
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	case err != nil:
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	case !accessPermited:
		http.Error(w, "Пользователь не имеет доступа", http.StatusForbidden)
		return
	case newBanner.Content == nil:
		http.Error(w, "Баннер не найден", http.StatusNotFound)
		return
	default:
	}

	jsonContent, err := b.createFromEntity(newBanner)
	if err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonContent)
}

func (b *Banner) createFromEntity(entity banner.BannerEntity) ([]byte, error) {
	jsonContent, err := json.Marshal(entity.Content)
	if err != nil {
		return nil, err
	}
	return jsonContent, nil
}
