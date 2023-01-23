package config

import (
	"github.com/joho/godotenv"
	"github.com/vrischmann/envconfig"
	"go.uber.org/zap"

	"github.com/CBrather/go-auth/pkg/log"
)

var env *EnvConfig = nil

func GetEnvironment() *EnvConfig {
	if env == nil {
		initConfig()
	}

	return env
}

func initConfig() {
	logger, err := log.GetLoggerWithLevel("info")
	if err != nil {
		zap.L().Fatal("Failed to get a new logger", zap.Error(err))
	}

	err = godotenv.Load(".env")
	if err != nil {
		// Only logging on Info level, as a .env file not being present would be intentional in non-local environments
		logger.Info("Unable to load a .env file, will execute with environment as-is", zap.Error(err))
	}

	env = new(EnvConfig)

	if err = envconfig.Init(env); err != nil {
		logger.Fatal("Failed initializing the app config", zap.Error(err))
	}
}

type EnvConfig struct {
	Auth struct {
		Domain   string `envconfig:"AUTH_DOMAIN"`
		Audience string `envconfig:"AUTH_AUDIENCE"`
	}

	Database struct {
		Host     string `envconfig:"DB_HOST"`
		Port     string `envconfig:"DB_PORT"`
		Name     string `envconfig:"DB_NAME"`
		User     string `envconfig:"DB_USER"`
		Password string `envconfig:"DB_PASSWORD"`
		SslMode  string `envconfig:"DB_SSLMODE,default=require"`
	}

	LogLevel string `envconfig:"LOGLEVEL"`

	OTelExporter struct {
		InsecureMode bool   `envconfig:"OTEL_EXPORTER_INSECURE_MODE,default=false"`
		OTLPEndpoint string `envconfig:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	}
}
