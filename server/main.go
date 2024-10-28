package main

import (
	"context"
	"log"
	"net"

	entity "github.com/WenChunTech/Microservices-Template/gen"
	"google.golang.org/grpc"
)

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
}
