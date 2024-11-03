package middlewares

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/WenChunTech/Microservices-Template/middlewares/auth"
	"github.com/WenChunTech/Microservices-Template/middlewares/recovery"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
)

func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Do something
		return handler(ctx, req)
	}
}

func GrpcServer() *grpc.Server {
	return grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				AuthInterceptor(),
				grpc_auth.UnaryServerInterceptor(auth.AuthInterceptor),
				grpc_recovery.UnaryServerInterceptor(recovery.RecoveryInterceptor()),
			),
		),
	)

}

func GrpcClient() *grpc.ClientConn {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithUnaryInterceptor(),
	}

	grpc.Dial("localhost:8080", opts...)

	return nil
}
