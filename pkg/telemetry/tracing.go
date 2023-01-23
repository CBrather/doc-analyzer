package telemetry

import (
	"context"

	"github.com/CBrather/go-auth/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"go.uber.org/zap"
	"google.golang.org/grpc/credentials"
)

func InitTracer(config *config.EnvConfig) func(context.Context) error {
	securityOption := getGrpcSecurityOption(config)

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			securityOption,
			otlptracegrpc.WithEndpoint(config.OTelExporter.OTLPEndpoint),
		),
	)

	if err != nil {
		zap.L().Fatal("Tracing :: Failed setting up the OTLP trace exporter", zap.Error(err))
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", "go-auth"),
		),
	)

	if err != nil {
		zap.L().Error("Tracing :: Could not set resources", zap.Error(err))
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	zap.L().Info("Tracing set up successfully")

	return exporter.Shutdown
}

func getGrpcSecurityOption(config *config.EnvConfig) otlptracegrpc.Option {
	exportInsecure := config.OTelExporter.InsecureMode

	if exportInsecure {
		return otlptracegrpc.WithInsecure()
	}

	return otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
}
