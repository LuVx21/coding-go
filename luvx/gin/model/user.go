package model

import "time"

type User struct {
    Id        uint      `gorm:"primary_key" json:"id"`
    Username  string    `gorm:"unique_index" json:"username"`
    Password  string    `json:"-"`
    Age       int8      `gorm:"size:3" json:"age"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
