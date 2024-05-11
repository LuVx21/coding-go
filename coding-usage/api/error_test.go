package api

import (
    "fmt"
    "testing"
)

func Test_error_00(t *testing.T) {
    err := fmt.Errorf("异常1")
    err = fmt.Errorf("2: %w", err)
    fmt.Println(err)
}
