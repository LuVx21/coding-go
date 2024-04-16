package common

import (
    "fmt"
)

type User struct {
    Id   uint64
    Name string `json:"name"`
    Age  uint8
}

func (u *User) String() string {
    return fmt.Sprintf("%d,%s,%d", u.Id, u.Name, u.Age)
}

func CheckErr(err error) {
    if err != nil {
        panic(err)
    }
}
