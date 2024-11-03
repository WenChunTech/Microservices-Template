package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/WenChunTech/Microservices-Template/entity"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var port = flag.Int("port", 50051, "the port to serve on")
var restful = flag.Int("restful", 8080, "the port to restful serve on")

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

	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatalln("Failed to listen: ", err)
	}
	server.Serve(listener)

	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	err = entity.RegisterEntityServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", *restful),
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0" + fmt.Sprintf(":%d", *restful))
	log.Fatalln(gwServer.ListenAndServe())

}
