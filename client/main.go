package main

import (
	"context"
	"errors"
	"log"

	"github.com/WenChunTech/Microservices-Template/entity"
	provider "github.com/WenChunTech/Microservices-Template/otel"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 设置OpenTelemetry。
	serviceName := "dice"
	serviceVersion := "0.1.0"

	otelShutdown, err := provider.SetupOTelSDK(context.Background(), serviceName, serviceVersion)
	if err != nil {
		log.Fatalln("Failed to setup OpenTelemetry: ", err)
	}
	// 适当处理关闭，以避免泄漏。
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	tp := otel.GetTracerProvider()
	mp := otel.GetMeterProvider()
	p := otel.GetTextMapPropagator()

	conn, err := grpc.NewClient("localhost:8081",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			otelgrpc.UnaryClientInterceptor(
				otelgrpc.WithTracerProvider(tp),
				otelgrpc.WithMeterProvider(mp),
				otelgrpc.WithPropagators(p),
			),
		),
	)
	if err != nil {
		log.Fatalln("Failed to dial server: ", err)
	}

	defer conn.Close()

	client := entity.NewEntityServiceClient(conn)

	ctx := context.Background()

	res, err := client.GetEntity(ctx, &entity.EntityRequest{Id: "1"})

	if err != nil {
		log.Fatalln("Failed to get entity: ", err)
	}

	log.Println("Client: GetEntity response: ", res)
}
