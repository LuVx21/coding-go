package ios

import (
    "fmt"
    "os"
    "testing"
)

var path = "/Users/renxie/OneDrive/Code/coding-go/coding-common/ios/ios.go"

func Test_io_00(t *testing.T) {
    lines, _ := ReadLines(path)
    fmt.Println(len(lines), lines)
    lines, _ = ReadLines1(path)
    fmt.Println(len(lines), lines)
}

func Test_io_01(t *testing.T) {
    b, _ := os.ReadFile(path)
    fmt.Println(string(b))
}
