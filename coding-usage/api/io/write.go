package io

import (
    "bufio"
    "fmt"
    "github.com/luvx21/coding-go/coding-usage/api/common"
    "io"
    "io/ioutil"
    "os"
    "time"
)

var s = time.Now().String() + "\nfoo\nbar\n"
var fileName = "1.txt"

func fileExist(filename string) bool {
    if _, err := os.Stat(filename); os.IsNotExist(err) {
        return false
    }
    return true
}

/**
每次写入会清空
*/
func w1() {
    content := []byte(s)
    err := ioutil.WriteFile(fileName, content, 0644)
    if err != nil {
        panic(err)
    }
}

func w2() {
    var f *os.File
    var err error
    f, err = os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
    defer f.Close()
    common.CheckErr(err)

    n, err := io.WriteString(f, s)
    common.CheckErr(err)
    fmt.Printf("写入 %d 个字节\n", n)
}

func w3() {
    var f *os.File
    var err error
    f, err = os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
    defer f.Close()
    common.CheckErr(err)

    n, err := f.Write([]byte(s))
    common.CheckErr(err)
    fmt.Printf("写入 %d 个字节\n", n)

    n, err = f.WriteString(s)
    common.CheckErr(err)
    fmt.Printf("写入 %d 个字节\n", n)
    f.Sync()
}

func w4() {
    var f *os.File
    var err error
    if fileExist(fileName) {
        f, err = os.OpenFile(fileName, os.O_APPEND, 0666)
    } else {
        f, err = os.Create(fileName)
    }
    defer f.Close()
    common.CheckErr(err)

    w := bufio.NewWriter(f)
    n, err := w.WriteString(s)
    common.CheckErr(err)
    fmt.Printf("写入 %d 个字节\n", n)
    w.Flush()
}
