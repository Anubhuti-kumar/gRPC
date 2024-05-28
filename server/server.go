package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "user-service/user-service/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUserServiceServer
	users []*pb.User
	mu    sync.Mutex
}

// Get User by ID
func (s *server) GetUser(ctx context.Context, req *pb.UserIDRequest) (*pb.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, user := range s.users {
		if user.Id == req.Id {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

// Get Users
func (s *server) GetUsers(ctx context.Context, req *pb.UserIDsRequest) (*pb.UserList, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var users []*pb.User
	for _, id := range req.Ids {
		for _, user := range s.users {
			if user.Id == id {
				users = append(users, user)
			}
		}
	}
	return &pb.UserList{Users: users}, nil
}

// Search User
func (s *server) SearchUsers(ctx context.Context, req *pb.SearchRequest) (*pb.UserList, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var users []*pb.User
	for _, user := range s.users {
		if (req.City == "" || user.City == req.City) &&
			(req.Phone == 0 || user.Phone == req.Phone) &&
			(user.Married == req.Married) {
			users = append(users, user)
		}
	}
	return &pb.UserList{Users: users}, nil
}

// Add User
func (s *server) AddUser(ctx context.Context, req *pb.User) (*pb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.users = append(s.users, req)
	return &pb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
