package globalVariable

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
)

//  MicroserviceInit 初始化微微服务
//  serviceStartFn  服务启动前函数
//  beforeListen 监听前函数
//  ListenError  错误处理函数
//  port 端口
func MicroserviceInit(serviceStartFn func(handle *grpc.Server), beforeListen func(), ListenError func(err string), port int) {
	serviceHandle := grpc.NewServer()
	serviceStartFn(serviceHandle)
	monitor, listenErr := net.Listen("tcp", fmt.Sprint(":", string(port)))
	beforeListen()
	if listenErr != nil {
		ListenError(listenErr.Error())
	}
	serviceHandle.Serve(monitor)
}
