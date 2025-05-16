package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/WenChunTech/Microservices-Template/entity"
	"github.com/WenChunTech/Microservices-Template/logger"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	provider "github.com/WenChunTech/Microservices-Template/otel"
)

var grpcPort = flag.Int("port", 8081, "the port to grpc serve on")
var restfulPort = flag.Int("restful", 8080, "the port to restful serve on")

type SpecificEntity struct {
	entity.UnimplementedEntityServiceServer
}

func (e SpecificEntity) GetEntity(context.Context, *entity.EntityRequest) (*entity.Entity, error) {
	log.Println("Server: GetEntity called")
	res := entity.Entity{
		Id: "1",
	}
	return &res, nil
}

func main() {
	// 设置OpenTelemetry。
	serviceName := "dice"
	serviceVersion := "0.1.0"

	otelShutdown, err := provider.SetupOTelSDK(context.Background(), serviceName, serviceVersion)
	if err != nil {
		log.Fatalln("Failed to setup OpenTelemetry: ", err)
	}
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	tp := otel.GetTracerProvider()
	mp := otel.GetMeterProvider()
	p := otel.GetTextMapPropagator()

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			otelgrpc.UnaryServerInterceptor(
				otelgrpc.WithTracerProvider(tp),
				otelgrpc.WithMeterProvider(mp),
				otelgrpc.WithPropagators(p),
			),
		),
	)

	entity.RegisterEntityServiceServer(server, SpecificEntity{})

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
	if err != nil {
		logger.Error("Failed to listen: ", err)
	}

	go func() {
		server.Serve(listener)
		logger.Info("Server started on port: ", *grpcPort)
	}()

	conn, err := grpc.NewClient(
		fmt.Sprintf(":%d", *grpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error("Failed to create client connection: ", err)
	}

	gwmux := runtime.NewServeMux()
	err = entity.RegisterEntityServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		logger.Error("Failed to register gateway: ", err)
	}

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", *restfulPort),
		Handler: gwmux,
	}

	err = gwServer.ListenAndServe()
	if err != nil {
		logger.Error("Failed to start restful server: ", err)
	}

}
