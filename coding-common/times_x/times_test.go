package times_x

import (
    "fmt"
    "testing"
)

func Test_00(t *testing.T) {
    for _, v := range TimeFormats {
        fmt.Println(v.Format)
    }
}
