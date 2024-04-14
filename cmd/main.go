package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	createBanner "github.com/CyberPiess/banner_service/internal/application/handler/create_banner"
	deleteBanner "github.com/CyberPiess/banner_service/internal/application/handler/delete_banner"
	getBannerList "github.com/CyberPiess/banner_service/internal/application/handler/get_banner_list"
	getUserBanner "github.com/CyberPiess/banner_service/internal/application/handler/get_user_banner"
	updateBanner "github.com/CyberPiess/banner_service/internal/application/handler/update_banner"
	"github.com/sirupsen/logrus"

	bannerService "github.com/CyberPiess/banner_service/internal/domain/banner"
	"github.com/gorilla/mux"

	"github.com/CyberPiess/banner_service/internal/infrastructure/logging"
	"github.com/CyberPiess/banner_service/internal/infrastructure/postgres"
	bannerStorage "github.com/CyberPiess/banner_service/internal/infrastructure/postgres/banner"
	"github.com/CyberPiess/banner_service/internal/infrastructure/redis"
	redisCache "github.com/CyberPiess/banner_service/internal/infrastructure/redis/cache"

	"github.com/joho/godotenv"
)

func main() {

	logger, err := logging.LoggerCreate(logging.Config{
		LogLevel: "info",
		LogFile:  "logrus.log",
	})
	if err != nil {
		log.Fatal("Error starting log")
	}
	dir, err := os.Getwd()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"package":  "main",
			"function": "main",
			"error":    err,
		}).Error()
	}

	err = godotenv.Load(fmt.Sprintf("%s/build/.env", dir))
	if err != nil {
		logger.Fatal("Error loading .env file")
	}
	if err != nil {
		logger.Fatal("Error loading .env file")
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
	if err != nil {
		logger.Fatalf("failed to initialize db: %s", err.Error())
	}
	defer db.Close()

	client, err := redis.NewRedis(redis.Config{
		Addres:        os.Getenv("REDIS_ADDRESS"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
	})

	if err != nil {
		logger.Fatalf("failed to initialize redis: %s", err.Error())
	}
	defer client.Close()

	bannerStore := bannerStorage.NewBannerRepository(db, logger)
	redisCache := redisCache.NewBannerCache(client, logger)
	bannerService := bannerService.NewBannerService(bannerStore, redisCache, logger)

	userBannerHandler := getUserBanner.NewBannerHandler(bannerService, logger)
	adminBannerHandler := getBannerList.NewGetAllBannersHandler(bannerService, logger)
	postBannerHandler := createBanner.NewPostBannerHandler(bannerService, logger)
	putBannerHandler := updateBanner.NewPutBannerHandler(bannerService, logger)
	deleteBannerHandler := deleteBanner.NewDeleteBannerHandler(bannerService, logger)

	mux.HandleFunc("/user_banner", userBannerHandler.GetUserBanner).Methods(http.MethodGet)
	mux.HandleFunc("/banner", adminBannerHandler.GetAllBanners).Methods(http.MethodGet)
	mux.HandleFunc("/banner", postBannerHandler.PostBanner).Methods(http.MethodPost)
	mux.HandleFunc("/banner/{id}", putBannerHandler.PutBanner).Methods(http.MethodPut)
	mux.HandleFunc("/banner/{id}", deleteBannerHandler.DeleteBanner).Methods(http.MethodDelete)

	err = http.ListenAndServe(":8080", mux)
	logger.Fatal(err)
}
