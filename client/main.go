package main

import (
	"context"
	"log"

	"github.com/WenChunTech/Microservices-Template/entity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
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
