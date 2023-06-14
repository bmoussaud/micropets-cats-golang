package internal

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func initTracerAuto() func(context.Context) error {
	fmt.Printf("initTracerAuto//////////////////\n")
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint("wavefront-proxy.observability-system.svc.cluster.local:4317"),
		),
	)

	if err != nil {
		log.Fatal("Could not set exporter: ", err)
	}

	fmt.Printf("exporter %+v\n", exporter)

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", LoadConfiguration().Observability.Service),
			attribute.String("application", LoadConfiguration().Observability.Application),
		),
	)
	fmt.Printf("resources %+v\n", resources)
	if err != nil {
		log.Fatal("Could not set resources: ", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter)),
			sdktrace.WithSyncer(exporter),
			sdktrace.WithResource(resources),
		),
	)

	return exporter.Shutdown
}

func NewGlobalTracer() io.Closer {
	fmt.Printf("NewGlobalTracer//////////////////\n")
	cleanup := initTracerAuto()
	defer cleanup(context.Background())

	return ioutil.NopCloser(nil)

}

func NewServerSpan(req *http.Request, spanName string) io.Closer {

	return ioutil.NopCloser(nil)
}

func NewTrace(ctx context.Context, traceName string) {
	_, span := otel.Tracer(traceName).Start(ctx, "Run")
	defer span.End()
}
