package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"simple_example/proto"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial("localhost:8000", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := proto.NewUserClient(conn)

	// 调用rpc服务AddUser方法
	resp, err := client.AddUser(context.Background(), &proto.UserRequest{Name: "zhengwenfeng", Age: 18})
	if err != nil {
		log.Fatalf("fail to AddUser: %v", err)
	}
	fmt.Printf("AddUser, msg = %s, code = %d\n", resp.Msg, resp.Code)

	// 调用rpc服务GetUser方法
	getuserResp, err := client.GetUser(context.Background(), &proto.GetUserRequest{Name: "zhangsan"})
	if err != nil {
		log.Fatalf("fail to GetUser: %v", err)
	}
	fmt.Printf("GetUser, Name = %s, Age = %d", getuserResp.Name, getuserResp.Age)

}
