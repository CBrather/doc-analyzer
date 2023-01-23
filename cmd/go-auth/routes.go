package main

import (
	"context"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/riandyrn/otelchi"

	"go.uber.org/zap"

	"github.com/CBrather/go-auth/internal/api"
	"github.com/CBrather/go-auth/internal/config"
	"github.com/CBrather/go-auth/pkg/telemetry"
)

func SetupHttpRoutes() {
	logger := httplog.NewLogger("go-auth", httplog.Options{JSON: true, Concise: true})

	config := config.GetEnvironment()

	traceShutdown := telemetry.InitTracer(config)
	defer traceShutdown(context.Background())

	router := chi.NewRouter()

	router.Use(
		otelchi.Middleware("go-auth", otelchi.WithChiRoutes(router)),
		httplog.RequestLogger(logger),
		middleware.Recoverer,
	)

	router.Handle("/metrics", telemetry.NewMetricsHandler())
	api.SetupProbeRoutes(router)

	api.SetupAlbumRoutes(router)

	zap.L().Info("Server listening on :8080")
	http.ListenAndServe("0.0.0.0:8080", router)
}
