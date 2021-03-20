package router

import (
	"MicroserviceGoBlog/register/protocol"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"strconv"
)

func Register(C *gin.Context) {
	// consul 寻找服务
	registerDefaultConf := api.DefaultConfig()
	registerClient, _ := api.NewClient(registerDefaultConf)
	servers, _, _ := registerClient.Health().Service("register", "verifyRegister", true, nil)
	target := servers[0].Service.Address + ":" + strconv.Itoa(servers[0].Service.Port)
	// 调用服务
	registerClientConn, _ := grpc.Dial(target, grpc.WithInsecure())
	defer registerClientConn.Close()
	registerServiceClient := protocol.NewRegisterServiceClient(registerClientConn)

	isRegisterResult, _ := registerServiceClient.IsRegistered(context.TODO(), &protocol.IsRegisterRequest{
		Username: "user",
	})
	fmt.Println(isRegisterResult, "isRegister")

}
