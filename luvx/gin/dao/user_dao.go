package dao

import (
    "fmt"
    "luvx/gin/db"
    "luvx/gin/model"
)

func GetUserByUsername(username string) (*model.User, error) {
    var user model.User
    if err := db.MySQLClient.Where("user_name = ?", username).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func crud() {
    // 插入记录
    client := db.MySQLClient
    client.Create(&model.User{UserName: "xxx", Age: 18, Password: "xxx"})

    var users []model.User
    client.Find(&users)
    fmt.Println(users)

    // 查询一条记录
    var user model.User
    client.First(&user, "name = ? and password = ?", "xxx", "xxx")
    fmt.Println("查询:", user)

    // 更新记录(基于查出来的数据进行更新)
    client.Model(&user).Update("name", "biaoge")
    fmt.Println("更新后的记录:", user)

    client.Delete(&user)

    client.Find(&users)
    fmt.Println(users)
}
