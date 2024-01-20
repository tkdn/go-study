package telemetry

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/tkdn/go-study/log"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

const (
	serviceName       = "go-study"
	deployEnvironment = "local"
	serviceVersion    = "0.0.1" // Commit hash etc.
)

func NewOtelHttpMiddleware() func(http.Handler) http.Handler {
	return otelhttp.NewMiddleware("",
		otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
		otelhttp.WithPublicEndpoint(),
		otelhttp.WithSpanNameFormatter(func(_ string, r *http.Request) string {
			return fmt.Sprintf("%s %s", r.Method, r.URL.Path)
		}))
}

func Do(ctx context.Context) (func(), error) {
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	resc, err := resource.New(ctx,
		resource.WithAttributes(newResource().Attributes()...),
		resource.WithContainer(),
		resource.WithContainerID(),
		resource.WithTelemetrySDK())
	if err != nil && !errors.Is(err, resource.ErrPartialResource) {
		return nil, err
	}
	provider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(resc))
	otel.SetTracerProvider(provider)
	// Propagator については基本 TextMapPropagator のみ
	// see: https://opentelemetry.io/docs/specs/otel/context/api-propagators/
	// otel.SetTextMapPropagator(awesome.Propagator)
	cleanup := func() {
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		if err := provider.ForceFlush(ctx); err != nil {
			log.Logger.Error(err.Error())
		}
		defer cancel()
		if err := provider.Shutdown(shutdownCtx); err != nil {
			log.Logger.Warn("failed to shutdown otel tracerprovider")
		}
	}
	return cleanup, nil
}

func newResource() *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
		semconv.ServiceVersion(serviceVersion),
		semconv.DeploymentEnvironment(deployEnvironment),
	)
}
