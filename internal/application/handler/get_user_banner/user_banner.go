//go:generate mockgen -source=user_banner.go -destination=mocks/mock.go
package userbanner

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/CyberPiess/banner_sevice/internal/domain/banner"
)

type bannerService interface {
	SearchBanner(ctx context.Context, bannerFilter banner.Filter, user banner.User) (banner.BannerEntity, bool, error)
}

type Banner struct {
	service bannerService
}

func NewBannerHandler(service bannerService) *Banner {
	return &Banner{service: service}
}

func (b *Banner) GetUserBanner(w http.ResponseWriter, r *http.Request) {

	bannerFilter, err := b.createFilterFromRequest(r)
	if err != nil {
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}

	user := b.recieveUserTokenFromRequest(r)
	if user.Token == "" {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	var bnnr banner.BannerEntity
	ctx := context.WithValue(context.Background(), "banner", bnnr)
	newBanner, accessPermited, err := b.service.SearchBanner(ctx, bannerFilter, user)
	switch {
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

func (b *Banner) createFilterFromRequest(r *http.Request) (banner.Filter, error) {
	var bannerFilter banner.Filter
	var err error

	bannerFilter.TagId, err = strconv.Atoi(r.FormValue("tag_id"))
	if err != nil {
		return banner.Filter{}, err
	}

	bannerFilter.FeatureId, err = strconv.Atoi(r.FormValue("feature_id"))
	if err != nil {
		return banner.Filter{}, err
	}

	useLastRevision := r.FormValue("use_last_revision")
	if useLastRevision == "" {
		bannerFilter.UseLastRevision = false
		return bannerFilter, nil
	}

	bannerFilter.UseLastRevision, err = strconv.ParseBool(useLastRevision)
	if err != nil {
		return banner.Filter{}, err
	}

	return bannerFilter, nil
}

func (b *Banner) recieveUserTokenFromRequest(r *http.Request) banner.User {
	var user banner.User
	user.Token = r.Header.Get("token")
	return user
}

func (b *Banner) createFromEntity(entity banner.BannerEntity) ([]byte, error) {
	jsonContent, err := json.Marshal(entity.Content)
	if err != nil {
		return nil, err
	}
	return jsonContent, nil
}
