package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"go.uber.org/zap"

	"github.com/CBrather/go-auth/internal/config"
	"github.com/CBrather/go-auth/pkg/log"
)

func main() {
	env := config.GetEnvironment()

	if err := log.Initialize(env.LogLevel); err != nil {
		zap.L().Fatal("Failed to setup logger")
	} else {
		zap.L().Info("Logger was successfully setup")
	}

	db, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		env.Database.Host, env.Database.Port, env.Database.Name, env.Database.User, env.Database.Password, env.Database.SslMode))

	if err != nil {
		zap.L().Fatal("Unable to open a Postgres connection", zap.Error(err))
	}

	SetupHttpRoutes(db)
}
