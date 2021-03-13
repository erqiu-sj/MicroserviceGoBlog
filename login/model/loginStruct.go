package model

import (
	"MicroserviceGoBlog/publicModel"
	"gorm.io/gorm"
)

type Login struct {
	publicModel.BasicLoginAndRegisterField
	gorm.Model
}
