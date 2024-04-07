//go:generate mockgen -source=banner.go -destination=mocks/mock.go
package adminbanner

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/CyberPiess/banner_sevice/internal/domain/banner"
)

type bannerService interface {
	SearchBanner(bannerFilter banner.Filter, user banner.User) (banner.BannerEntity, bool, error)
	SearchAllBanners(bannerFilter banner.Filter, user banner.User) ([]banner.BannerEntity, bool, error)
}

type Banner struct {
	service bannerService
}

func NewBannerHandler(service bannerService) *Banner {
	return &Banner{service: service}
}

func (b *Banner) GetAllBanners(w http.ResponseWriter, r *http.Request) {

	bannerFilter := banner.Filter{
		TagId:           r.FormValue("tag_id"),
		FeatureId:       r.FormValue("feature_id"),
		UseLastRevision: r.FormValue("use_last_revision"),
		Limit:           r.FormValue("limit"),
		Offset:          r.FormValue("offset"),
	}

	user := banner.User{
		Token: r.Header.Get("token"),
	}

	foundBanners, accessPermited, err := b.service.SearchAllBanners(bannerFilter, user)
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

	jsonContent, err := b.createFromEntity(foundBanners)
	if err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonContent)
}

func (b *Banner) createFromEntity(entityList []banner.BannerEntity) ([]byte, error) {
	type response struct {
		ID        int
		Content   map[string]interface{}
		TagId     []int
		FeatureId int
		IsActive  bool
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	var result []response

	for _, entity := range entityList {
		partOfSlice := response{ID: entity.ID,
			Content:   entity.Content,
			TagId:     entity.TagId,
			FeatureId: entity.FeatureId,
			IsActive:  entity.IsActive,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt}
		result = append(result, partOfSlice)
	}

	jsonContent, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	return jsonContent, nil
}
