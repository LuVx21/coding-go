package service

import (
    "github.com/luvx21/coding-go/coding-common/logs"
    "luvx/gin/dao"
    "luvx/gin/model"
)

func GetUserByUsername(username string) (*model.User, error) {
    logs.Log.Infoln("username:", username)
    return dao.GetUserByUsername(username)
}
