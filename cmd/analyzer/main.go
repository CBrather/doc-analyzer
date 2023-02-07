package main

import (
	_ "github.com/lib/pq"

	"go.uber.org/zap"

	"github.com/CBrather/analyzer/internal/config"
	"github.com/CBrather/analyzer/pkg/log"
)

func main() {
	env := config.GetEnvironment()

	if err := log.Initialize(env.LogLevel); err != nil {
		zap.L().Fatal("Failed to setup logger")
	} else {
		zap.L().Info("Logger was successfully setup")
	}

	SetupHttpRoutes()
}
