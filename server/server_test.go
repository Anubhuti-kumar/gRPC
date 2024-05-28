package main

import (
	"context"
	"log"
	"net"
	"testing"

	pb "user-service/user-service/proto"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{})

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestAddUser(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	assert.NoError(t, err)
	defer conn.Close()
	client := pb.NewUserServiceClient(conn)

	_, err = client.AddUser(ctx, &pb.User{Id: 1, Fname: "Steve", City: "LA", Phone: 1234567890, Height: 5.8, Married: true})
	assert.NoError(t, err)
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	assert.NoError(t, err)
	defer conn.Close()
	client := pb.NewUserServiceClient(conn)

	client.AddUser(ctx, &pb.User{Id: 1, Fname: "Steve", City: "LA", Phone: 1234567890, Height: 5.8, Married: true})

	user, err := client.GetUser(ctx, &pb.UserIDRequest{Id: 1})
	assert.NoError(t, err)
	assert.Equal(t, "Steve", user.Fname)
}

func TestSearchUsers(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	assert.NoError(t, err)
	defer conn.Close()
	client := pb.NewUserServiceClient(conn)

	client.AddUser(ctx, &pb.User{Id: 1, Fname: "Steve", City: "LA", Phone: 1234567890, Height: 5.8, Married: true})
	client.AddUser(ctx, &pb.User{Id: 2, Fname: "John", City: "NY", Phone: 9876543210, Height: 5.7, Married: false})

	users, err := client.SearchUsers(ctx, &pb.SearchRequest{City: "LA", Married: true})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(users.Users))
	assert.Equal(t, "Steve", users.Users[0].Fname)
}
