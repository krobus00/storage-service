package infrastructure

import (
	"fmt"
	"net"

	"github.com/krobus00/storage-service/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func JaegerTraceProvider() (tp *sdktrace.TracerProvider, err error) {
	if config.DisableTracing() {
		return nil, nil
	}

	var exp *jaeger.Exporter

	jaegerURL := net.JoinHostPort(config.JaegerHost(), config.JaegerPort())
	switch config.JaegerProtocol() {
	case "http":
		exp, err = jaeger.New(
			jaeger.WithCollectorEndpoint(
				jaeger.WithEndpoint(fmt.Sprintf("http://%s/api/traces", jaegerURL)),
			),
		)
	case "grpc":
		exp, err = jaeger.New(
			jaeger.WithAgentEndpoint(
				jaeger.WithAgentHost(config.JaegerHost()),
				jaeger.WithAgentPort(config.JaegerPort()),
			),
		)
	default:
		exp, err = jaeger.New(
			jaeger.WithCollectorEndpoint(
				jaeger.WithEndpoint(fmt.Sprintf("http://%s/api/traces", jaegerURL)),
			),
		)
	}

	if err != nil {
		return nil, err
	}

	tp = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithSampler(sdktrace.ParentBased(
			sdktrace.TraceIDRatioBased(config.JaegerSampleRate())),
		),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(config.ServiceName()),
			semconv.ServiceVersionKey.String(config.ServiceVersion()),
			semconv.DeploymentEnvironmentKey.String(config.Env()),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp, nil
}
