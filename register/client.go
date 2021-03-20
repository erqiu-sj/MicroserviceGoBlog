package main

import (
	"MicroserviceGoBlog/register/protocol"
	"errors"
	"google.golang.org/grpc"
)

type RegisterService struct {
}

func main() {

	conn, err := grpc.Dial("localhost:8081", grpc.WithInsecure())
	if err != nil {
		panic(errors.New("连接注册服务出错"))
	}
	defer conn.Close()
	protocol.NewRegisterServiceClient(conn)
	//registerService.IsRegistered(context.Context())
}
