package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/CBrather/go-auth/internal/api"
	"github.com/CBrather/go-auth/internal/api/middleware"
	"github.com/CBrather/go-auth/pkg/log"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		// Only logging on Warn level, as this might well be intentional, especially in production environments
		zap.L().Warn("Unable to load a .env file, will execute with environment as-is", zap.Error(err))
	}

	if err := log.Initialize(os.Getenv("LOGLEVEL")); err != nil {
		zap.L().Fatal("Failed to setup logger")
	} else {
		zap.L().Info("Logger was successfully setup")
	}

	db, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_SSLMODE")))
	if err != nil {
		zap.L().Fatal("Unable to open a Postgres connection", zap.Error(err))
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.JsonLoggerMiddleware())

	api.SetupAlbumRoutes(router, db)
	api.SetupProbeRoutes(router)

	err = router.Run("0.0.0.0:8080")
	if err != nil {
		zap.L().Fatal("Error on setting up HTTP server", zap.Error(err))
	}
}
