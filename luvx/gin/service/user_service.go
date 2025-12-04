package service

import (
	"luvx/gin/dao"
	"luvx/gin/model"

	log "github.com/sirupsen/logrus"
)

func GetUserByUsername(username string) (*model.User, error) {
	log.Infoln("username:", username)
	return dao.GetUserByUsername(username)
}
