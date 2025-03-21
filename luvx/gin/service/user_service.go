package service

import (
	"github.com/luvx21/coding-go/infra/logs"
	"luvx/gin/dao"
	"luvx/gin/model"
)

func GetUserByUsername(username string) (*model.User, error) {
	logs.Log.Infoln("username:", username)
	return dao.GetUserByUsername(username)
}
