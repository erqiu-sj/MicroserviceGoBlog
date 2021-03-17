package main

import (
	"MicroserviceGoBlog/register/model"
	"MicroserviceGoBlog/register/protocol"
	"context"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net"
)

var (
	DB     *gorm.DB
	result *model.Register
)

type RegisterServiceImpl struct {
}

// TODO: RegisterResponse 类型缺少一个 message 类型
// IsRegistered 判断是否注册user
func (registerInfo *RegisterServiceImpl) IsRegistered(ctx context.Context, req *protocol.IsRegisterRequest) (*protocol.RegisterResponse, error) {
	response := &protocol.RegisterResponse{}
	DB.Find(&result, "username = ?", req.Username)
	if result.Username == "" {
		response.Status = false
	}
	return response, nil
}

func (registerInfo *RegisterServiceImpl) ReadyRegister(ctx context.Context, req *protocol.RegisterRequest) (*protocol.RegisterResponse, error) {
	response := &protocol.RegisterResponse{}
	register := model.Register{Username: req.Username, Password: req.Password, Email: req.Email, Birthday: req.Birthday, Gender: req.Gender}
	if register.Email != "" {
		emailState, _ := register.VerifyUserOrEmail(true, 0, 0)
		if !emailState {
			response.Status = emailState
			return response, nil
		}
	}
	usrNameState, _ := register.VerifyUserOrEmail(false, model.PswOptLower&model.PswOptNumber&model.PswOptSpecial&model.PswOptUpper, model.PswOptSpecial)
	if !usrNameState {
		response.Status = usrNameState
		return response, nil
	}
	DB.Table("register").Create(&req)
	response = new(protocol.RegisterResponse)
	response.Status = true
	return response, nil
}

func main() {
	dbConf := "root:Qsj.0228@tcp(127.0.0.1:3306)/Blog?charset=utf8&parseTime=True&loc=Local"

	DB, _ = gorm.Open(mysql.Open(dbConf), &gorm.Config{})
	DB.AutoMigrate(&model.Register{})
	service := grpc.NewServer()
	protocol.RegisterRegisterServiceServer(service, new(RegisterServiceImpl))
	monitor, _ := net.Listen("tcp", ":8081")
	service.Serve(monitor)
}
