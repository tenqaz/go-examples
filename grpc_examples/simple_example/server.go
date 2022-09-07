package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"simple_example/proto"
)

// UserService 定义结构体，实现UserServer
type UserService struct {
	proto.UnimplementedUserServer
}

func NewUserService() *UserService {
	return &UserService{}
}

// AddUser 实现rpc方法
func (us *UserService) AddUser(ctx context.Context, request *proto.UserRequest) (*proto.UserResponse, error) {
	fmt.Printf("add user success. name = %s, age = %d\n", request.GetName(), request.GetAge())
	return &proto.UserResponse{Msg: "success", Code: 0}, nil
}

func (us *UserService) GetUser(ctx context.Context, request *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	fmt.Printf("get user success. name = %s\n", request.GetName())
	return &proto.GetUserResponse{Name: request.GetName(), Age: 1999}, nil
}

func main() {
	// 监听端口
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	// 注册rpc服务
	proto.RegisterUserServer(grpcServer, NewUserService())
	grpcServer.Serve(lis)
}
