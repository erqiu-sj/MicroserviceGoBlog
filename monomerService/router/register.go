package router

import (
	"MicroserviceGoBlog/monomerService/globalVariable"
	"MicroserviceGoBlog/monomerService/utils"
	"MicroserviceGoBlog/register/protocol"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"net/http"
	"strconv"
)

type param struct {
	Username   string
	Password   string
	Email      string
	Birthday   string
	Gender     int64
	IsRegister bool
	IsStop     bool   // 是否停止服务
	StopMethod string // 停止服务的方式是否优雅？
}

func Register(C *gin.Context) {
	var p param
	if err := C.ShouldBindJSON(&p); err != nil {
		// 服务器出错
		utils.ReturnData(C, http.StatusInternalServerError, globalVariable.ERRORREADINGREGISTERTIONINFOMATION)
		return
	}
	// consul 寻找服务
	registerDefaultConf := api.DefaultConfig()
	registerClient, _ := api.NewClient(registerDefaultConf)
	servers, _, _ := registerClient.Health().Service("register", "verifyRegister", true, nil)
	target := servers[0].Service.Address + ":" + strconv.Itoa(servers[0].Service.Port)
	// 调用服务
	registerClientConn, _ := grpc.Dial(target, grpc.WithInsecure())
	defer registerClientConn.Close()

	registerServiceClient := protocol.NewRegisterServiceClient(registerClientConn)
	if p.IsStop && p.StopMethod != "" {
		// 关闭注册服务
		stopStatus, _ := registerServiceClient.StopRegister(context.TODO(), &protocol.StopRequest{StopMethod: p.StopMethod})
		utils.ReturnData(C, http.StatusOK, stopStatus.GetMessage())
		return
	}
	if p.IsRegister {
		// 查询是否注册
		isRegisterResult, _ := registerServiceClient.IsRegistered(context.TODO(), &protocol.IsRegisterRequest{
			Username: p.Username,
		})
		if isRegisterResult.GetStatus() {
			// 已注册，逻辑推出
			utils.ReturnData(C, http.StatusOK, isRegisterResult.GetMessage())
			return
		} else {
			// 未注册，调起注册服务
			registerResult, _ := registerServiceClient.ReadyRegister(context.TODO(), &protocol.RegisterRequest{
				Username: p.Username,
				Password: p.Password,
				Email:    p.Email,
				Birthday: p.Birthday,
				Gender:   p.Gender},
			)
			if !registerResult.GetStatus() {
				// 注册失败
				utils.ReturnData(C, http.StatusOK, registerResult.GetMessage())
				return
			} else {
				// 注册成功
				utils.ReturnData(C, http.StatusOK, registerResult.GetMessage())
				return
			}
		}
	}

}
