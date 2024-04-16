package defaults

import (
    "fmt"
    "testing"
)

type User struct {
    Name  string  `default:"Goku"`
    Power float64 `default:"9000.01"`
}

func Test_01(t *testing.T) {
    var u User
    _ = Apply(&u)

    fmt.Println(u.Name)  // Goku
    fmt.Println(u.Power) // 9000.01
}
