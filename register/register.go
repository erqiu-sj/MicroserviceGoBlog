package main

import (
	"MicroserviceGoBlog/monomerService/globalVariable"
	"MicroserviceGoBlog/register/model"
	"MicroserviceGoBlog/register/protocol"
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"net/http"
	"time"
)

const (
	PORT = 8082
)

var (
	registerServiceHandle *grpc.Server
)

type RegisterServiceImpl struct {
}

// IsRegistered 判断是否注册user
func (registerInfo *RegisterServiceImpl) IsRegistered(_ context.Context, req *protocol.IsRegisterRequest) (*protocol.RegisterResponse, error) {
	response := &protocol.RegisterResponse{}
	fmt.Println("start response!")
	var result *model.Register
	result = new(model.Register)
	fmt.Println("start result!")
	fmt.Println(globalVariable.Db, "DB TYPE!")
	globalVariable.Db.Debug().Find(&result, "username = ?", req.Username)
	fmt.Println("awsome result！")
	if result.Username != "" {
		response.Status = true
		response.Message = "用户已注册"
	} else {
		response.Status = false
		response.Message = "用户未注册"
	}
	response.HttpCode = http.StatusOK
	fmt.Println("注册", response)
	return response, nil
}

// ReadyRegister 准备注册user
func (registerInfo *RegisterServiceImpl) ReadyRegister(_ context.Context, req *protocol.RegisterRequest) (*protocol.RegisterResponse, error) {
	response := &protocol.RegisterResponse{}
	register := model.Register{Username: req.Username, Password: req.Password, Email: req.Email, Birthday: req.Birthday, Gender: req.Gender}
	if register.Email != "" {
		emailState, emailMessage := register.VerifyUserOrEmail(true, 0, 0)
		if !emailState {
			// 邮箱验证错误，返回204，未返回内容。在未更新网页的情况下，可确保浏览器继续显示当前文档
			response.Status = emailState
			response.Message = emailMessage
			response.HttpCode = http.StatusNoContent
			fmt.Println("邮箱错误", response)
			return response, nil
		}
	}
	usrNameState, usrMessage := register.VerifyUserOrEmail(false,
		model.PswOptLower&model.PswOptNumber&model.PswOptSpecial&model.PswOptUpper,
		model.PswOptSpecial)
	if !usrNameState {
		// 用户名验证失败
		response.Status = usrNameState
		response.Message = usrMessage
		response.HttpCode = http.StatusNoContent
		fmt.Println("用户名错误", response)
		return response, nil
	}
	globalVariable.Db.Table("register").Create(&req)
	response.Status = true
	response.HttpCode = http.StatusOK
	response.Message = model.REGISTRATIONSUCCESS
	fmt.Println("创建成功", response)
	return response, nil
}

// StopRegister 关闭注册服务
func (registerInfo *RegisterServiceImpl) StopRegister(_ context.Context, req *protocol.StopRequest) (*protocol.StopResponse, error) {
	if req.StopMethod == "stop" {
		// 关闭
		registerServiceHandle.Stop()
	}
	if req.StopMethod == "gracefulStop" {
		// 优雅的关闭
		registerServiceHandle.GracefulStop()
	}
	return &protocol.StopResponse{Message: req.StopMethod, Status: true}, nil
}

func main() {
	sql, _ := globalVariable.DbInit().DB()
	sql.SetMaxOpenConns(20)
	sql.SetMaxIdleConns(3)
	sql.SetConnMaxLifetime(5 * time.Minute)
	defer sql.Close()
	// 往consul中注册服务
	registerServiceConsulConf := api.DefaultConfig()
	registerClient, _ := api.NewClient(registerServiceConsulConf)
	registerReg := api.AgentServiceRegistration{
		ID:      "registerMain",
		Name:    "register",
		Tags:    []string{"register", "verifyRegister"},
		Address: globalVariable.TCP,
		Port:    PORT,
		Check: &api.AgentServiceCheck{
			CheckID: "checkRegister",
			Name:    "registerUser",
			//HTTP:    fmt.Sprint(globalVariable.HTTP, ":", PORT),
			TCP:      fmt.Sprint(globalVariable.TCP, ":", PORT),
			Interval: "10s",
			Timeout:  "5s",
		},
	}
	_ = registerClient.Agent().ServiceRegister(&registerReg)

	globalVariable.MicroserviceInit(func(handle *grpc.Server) {
		registerServiceHandle = handle
		protocol.RegisterRegisterServiceServer(registerServiceHandle, new(RegisterServiceImpl))
	}, func() {
		fmt.Println("register service start!")
	}, func(err string) {
		panic(err)
	}, 8082)

}
