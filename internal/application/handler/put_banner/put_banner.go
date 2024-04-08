//go:generate mockgen -source=put_banner.go -destination=mocks/mock.go
package putbanner

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/CyberPiess/banner_sevice/internal/domain/banner"
	"github.com/gorilla/mux"
)

type putBannerService interface {
	SearchBanner(bannerFilter banner.Filter, user banner.User) (banner.BannerEntity, bool, error)
	SearchAllBanners(bannerFilter banner.Filter, user banner.User) ([]banner.BannerEntity, bool, error)
	PostBanner(newPostBanner banner.BannerEntity, user banner.User) (int64, bool, error)
	PutBanner(newPutBanner banner.BannerEntity, user banner.User) (bool, bool, error)
}

type PutBanner struct {
	service putBannerService
}

func NewPutBannerHandler(service putBannerService) *PutBanner {
	return &PutBanner{service: service}
}

func (ptB *PutBanner) PutBanner(w http.ResponseWriter, r *http.Request) {
	var postBanner banner.BannerEntity
	var err error
	bannerID := mux.Vars(r)["id"]
	postBanner.ID, err = strconv.Atoi(bannerID)
	if err != nil {
		http.Error(w, "некорректные данные", http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&postBanner)
	if err != nil {
		http.Error(w, "некорректные данные", http.StatusBadRequest)
		return
	}
	user := banner.User{
		Token: r.Header.Get("token"),
	}

	found, accessPermited, err := ptB.service.PutBanner(postBanner, user)
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
	case !found:
		http.Error(w, "Баннер не найден", http.StatusNotFound)
		return
	default:
	}

	w.WriteHeader(http.StatusOK)
}
