package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/WenChunTech/Microservices-Template/entity"
	"github.com/WenChunTech/Microservices-Template/logger"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	server := grpc.NewServer()

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
