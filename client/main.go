package main

import (
	"context"
	"log"

	entity "github.com/WenChunTech/Microservices-Template/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
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
