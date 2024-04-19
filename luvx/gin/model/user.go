package model

import "time"

type User struct {
    Id         uint      `gorm:"primary_key" json:"id"`
    UserName   string    `gorm:"unique_index" json:"userName"`
    Password   string    `json:"password"`
    Age        int8      `gorm:"size:3" json:"age"`
    Birthday   time.Time `json:"birthday"`
    UpdateTime time.Time `json:"updateTime"`
}
