package router

import (
	"MicroserviceGoBlog/monomerService/globalVariable"
	"MicroserviceGoBlog/monomerService/utils"
	"MicroserviceGoBlog/register/protocol"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"net/http"
	"strconv"
)

type param struct {
	Username   string // Username 用户名
	Password   string // Password 密码
	Email      string // Email 邮箱
	Birthday   string // Birthday 生日
	Gender     int64  // Gender 性别
	IsRegister bool   // IsRegister 验证是否注册
	IsStop     bool   // IsStop 是否停止服务
	StopMethod string // StopMethod 停止服务的方式是否优雅？
}

func Register(C *gin.Context) {
	var p param
	if err := C.ShouldBindJSON(&p); err != nil {
		// 服务器出错
		utils.ReturnData(C, http.StatusInternalServerError, globalVariable.ERRORREADINGREGISTERTIONINFOMATION)
		return
	}
	fmt.Println(p, "params")
	fmt.Println(p.IsRegister, "isRegister")
	fmt.Println(p.Username, "user")
	// consul 寻找服务
	registerDefaultConf := api.DefaultConfig()
	registerClient, _ := api.NewClient(registerDefaultConf)
	servers, _, _ := registerClient.Health().Service("register", "verifyRegister", true, nil)
	target := servers[0].Service.Address + ":" + strconv.Itoa(servers[0].Service.Port)
	fmt.Println(target, "service consul")
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
		isRegisterResult, _ := registerServiceClient.IsRegistered(context.Background(), &protocol.IsRegisterRequest{
			Username: p.Username,
		})
		fmt.Println(isRegisterResult.GetStatus(), "isRegister")
		if isRegisterResult.GetStatus() {
			// 已注册，逻辑推出
			fmt.Println("no register!", isRegisterResult.GetStatus())
			utils.ReturnData(C, utils.Int64ToInt(isRegisterResult.GetHttpCode()), isRegisterResult.GetMessage())
			return
		} else {
			// 未注册，调起注册服务
			registerResult, _ := registerServiceClient.ReadyRegister(context.Background(), &protocol.RegisterRequest{
				Username: p.Username,
				Password: p.Password,
				Email:    p.Email,
				Birthday: p.Birthday,
				Gender:   p.Gender},
			)
			fmt.Println("registerStatus", registerResult.GetStatus())
			if !registerResult.GetStatus() {
				// 注册失败
				fmt.Println("注册失败!", registerResult.GetMessage(), registerResult.GetHttpCode())
				utils.ReturnData(C, utils.Int64ToInt(registerResult.GetHttpCode()), registerResult.GetMessage())
				return
			} else {
				// 注册成功
				fmt.Println("注册成功!", registerResult.GetMessage(), registerResult.GetHttpCode())
				utils.ReturnData(C, utils.Int64ToInt(registerResult.GetHttpCode()), registerResult.GetMessage())
				return
			}
		}
	}

}
