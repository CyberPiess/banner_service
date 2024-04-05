package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	appUser "github.com/CyberPiess/banner_sevice/internal/application/user"
	userService "github.com/CyberPiess/banner_sevice/internal/domain/user"

	"github.com/CyberPiess/banner_sevice/internal/infrastructure/postgres"
	userStorage "github.com/CyberPiess/banner_sevice/internal/infrastructure/postgres/user"

	"github.com/joho/godotenv"
)

func main() {

	currentDir, _ := os.Getwd()
	envFilePath := filepath.Join(currentDir, "..", "build\\.env")
	err := godotenv.Load(envFilePath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mux := http.NewServeMux()

	db, err := postgres.NewPostgresDb(postgres.Config{
		Host:     os.Getenv("PG_HOST"),
		Port:     os.Getenv("PG_PORT"),
		Username: os.Getenv("POSTGRES_USER"),
		DBName:   os.Getenv("DBNAME"),
		SSLMode:  os.Getenv("SSLMODE"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
	})

	if err != nil {
		log.Fatal("failed to initialize db: %s", err.Error())
	}

	userS := userStorage.NewUserRepository(db)
	userService := userService.NewUserService(userS)

	userHandler := appUser.NewUserHandler(userService)

	mux.HandleFunc("/user_banner", userHandler.GetUserBanner)

	err = http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
