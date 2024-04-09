package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	deleteBanner "github.com/CyberPiess/banner_sevice/internal/application/handler/delete_banner"
	adminBanner "github.com/CyberPiess/banner_sevice/internal/application/handler/get_banner"
	userBanner "github.com/CyberPiess/banner_sevice/internal/application/handler/get_user_banner"
	postBanner "github.com/CyberPiess/banner_sevice/internal/application/handler/post_banner"
	putBanner "github.com/CyberPiess/banner_sevice/internal/application/handler/put_banner"

	bannerService "github.com/CyberPiess/banner_sevice/internal/domain/banner"
	"github.com/gorilla/mux"

	"github.com/CyberPiess/banner_sevice/internal/infrastructure/postgres"
	bannerStorage "github.com/CyberPiess/banner_sevice/internal/infrastructure/postgres/banner"

	"github.com/joho/godotenv"
)

func main() {

	currentDir, _ := os.Getwd()
	envFilePath := filepath.Join(currentDir, "..", "build\\.env")
	err := godotenv.Load(envFilePath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mux := mux.NewRouter()

	db, err := postgres.NewPostgresDb(postgres.Config{
		Host:     os.Getenv("PG_HOST"),
		Port:     os.Getenv("PG_PORT"),
		Username: os.Getenv("POSTGRES_USER"),
		DBName:   os.Getenv("DBNAME"),
		SSLMode:  os.Getenv("SSLMODE"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
	})
	defer db.Close()

	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	bannerStore := bannerStorage.NewBannerRepository(db)
	bannerService := bannerService.NewBannerService(bannerStore)

	bannerHandler := userBanner.NewBannerHandler(bannerService)
	adminBannerHandler := adminBanner.NewBannerHandler(bannerService)
	postBannerHandler := postBanner.NewPostBannerHandler(bannerService)
	putBannerHandler := putBanner.NewPutBannerHandler(bannerService)
	deleteBannerHandler := deleteBanner.NewDeleteBannerHandler(bannerService)

	mux.HandleFunc("/user_banner", bannerHandler.GetUserBanner).Methods(http.MethodGet)
	mux.HandleFunc("/banner", adminBannerHandler.GetAllBanners).Methods(http.MethodGet)
	mux.HandleFunc("/banner", postBannerHandler.PostBanner).Methods(http.MethodPost)
	mux.HandleFunc("/banner/{id}", putBannerHandler.PutBanner).Methods(http.MethodPut)
	mux.HandleFunc("/banner/{id}", deleteBannerHandler.DeleteBanner).Methods(http.MethodDelete)

	err = http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
