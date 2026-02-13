package main

import (
	"context"
	"log"

	userService "github.com/souvikjs01/auth-microservice/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	grpcConn, err := grpc.NewClient("localhost:5000", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("cant connect to grpc server: %v", err)
	}
	defer grpcConn.Close()

	client := userService.NewUserServiceClient(grpcConn)

	ctx := context.Background()
	md := metadata.Pairs(
		"session_id", "ca63b7ce-6dfc-4ef5-9d27-72a964a285d2",
		"subsystem", "cli",
	)

	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := client.GetMe(ctx, &userService.GetMeRequest{})
	if err != nil {
		log.Fatalf("error getting user: %v", err)
	}
	log.Printf("response: %s", resp.String())
}
