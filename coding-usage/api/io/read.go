package io

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "strings"
)

var path = "/Users/renxie/OneDrive/Code/doc-jupyter/Go/gists/api/io/read_test.go"

func r1() {
    file, _ := os.Open(path)
    defer file.Close()
    content, _ := io.ReadAll(file)
    fmt.Println(string(content))
}

func r2() {
    content, _ := os.ReadFile(path)
    fmt.Println(string(content))
}

func r3() {
    fi, _ := os.Open(path)
    defer fi.Close()

    chunks := make([]byte, 0)
    buf := make([]byte, 1024)
    r := bufio.NewReader(fi)
    for {
        n, err := r.Read(buf)
        if err != nil && err != io.EOF {
            panic(err)
        }
        if n == 0 {
            break
        }
        chunks = append(chunks, buf[:n]...)
    }
    fmt.Println(string(chunks))
}

func r4() {
    f, _ := os.Open(path)
    defer f.Close()

    chunks := make([]byte, 0)
    buf := make([]byte, 1024)
    for {
        n, err := f.Read(buf)
        if err != nil && err != io.EOF {
            panic(err)
        }
        if n == 0 {
            break
        }
        chunks = append(chunks, buf[:n]...)
    }
    fmt.Println(string(chunks))
}

func r5() {
    file, _ := os.OpenFile(path, os.O_RDWR, 0666)
    defer file.Close()

    stat, _ := file.Stat()
    var size = stat.Size()
    fmt.Println("file size=", size)

    r := bufio.NewReader(file)
    for {
        line, err := r.ReadString('\n')
        if err != nil && err != io.EOF {
            panic(err)
        }
        line = strings.TrimSpace(line)
        fmt.Println(line)
    }
}
