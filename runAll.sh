#bin/bash

# 启动注册服务
# shellcheck disable=SC2164
cd ./register
go run register.go

# 启动登陆服务
cd ./article
go run main.go