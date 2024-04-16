package service

import (
    "luvx/gin/dao"
    "luvx/gin/model"
)

func GetUserByUsername(username string) (*model.User, error) {
    return dao.GetUserByUsername(username)
}
