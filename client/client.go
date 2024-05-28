package main

import (
	"context"
	"log"
	"time"

	pb "user-service/user-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	// Set up a connection to the server
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a new gRPC client instance
	c := pb.NewUserServiceClient(conn)

	// Set up a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Add a user
	_, err = c.AddUser(ctx, &pb.User{Id: 1, Fname: "Steve", City: "LA", Phone: 1234567890, Height: 5.8, Married: true})
	if err != nil {
		log.Fatalf("could not add user: %v", err)
	}

	// Get a user by ID
	user, err := c.GetUser(ctx, &pb.UserIDRequest{Id: 1})
	if err != nil {
		log.Fatalf("could not get user: %v", err)
	}
	log.Printf("User: %v", user)

	// Search users
	users, err := c.SearchUsers(ctx, &pb.SearchRequest{City: "LA", Married: true})
	if err != nil {
		log.Fatalf("could not search users: %v", err)
	}
	log.Printf("Users: %v", users.Users)
}
