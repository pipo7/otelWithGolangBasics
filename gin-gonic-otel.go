package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	service     = "medium-gin-server-test"
	environment = "development"
	id          = 1
)

var tracer = otel.Tracer(service)

func tracerProvider() (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint())
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
		),
		),
	)

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp, nil
}

func main() {

	fmt.Println("initializing")

	tp, err := tracerProvider()
	if err != nil {
		log.Fatal(err)
	}

	// Cleanly shutdown and flush telemetry when the application exits.
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	r := gin.New()
	r.Use(otelgin.Middleware("my-server"))
	loadRoutes(r)

	r.Run()
}

func loadRoutes(r *gin.Engine) {
	r.GET("/ping", pingFunc)
}

func pingFunc(c *gin.Context) {

	ctx, span := tracer.Start(c.Request.Context(), "/ping", oteltrace.WithAttributes(attribute.String("hello", "the user")))
	defer span.End()

	bar(ctx)

	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func ping2Func(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "/ping-2", oteltrace.WithAttributes(attribute.String("hello2", "the user 2")))
	defer span.End()

	bar(ctx)

	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func bar(ctx context.Context) {
	fmt.Println("on bar")
	// Use the global TracerProvider.
	ct, span := tracer.Start(ctx, "bar")
	span.SetAttributes(attribute.Key("testset").String("value"))
	defer span.End()

	time.Sleep(1 * time.Millisecond)

	go bar3(ct)
}

func bar3(ctx context.Context) {
	fmt.Println("on bar 3")

	_, span := tracer.Start(ctx, "bar-3-on-goroutine")
	span.AddEvent("starting goroutine bar3")

	defer func() {
		span.End()
	}()
	span.AddEvent("executing logic")
	time.Sleep(1 * time.Second)

	span.AddEvent("completed goroutine bar3")
}
