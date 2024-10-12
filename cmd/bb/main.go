package main

import (
    "flag"
    "fmt"
    "os"
)

func m1() {
    if len(os.Args) > 0 {
        for index, arg := range os.Args {
            fmt.Printf("args[%d]=%v\n", index, arg)
        }
    }
}

var (
    host = *flag.String("h", "127.0.0.1", "host(127.0.0.1)")
)

func m2() {
    flag.Parse()
    fmt.Println(host)
}

func main() {
    //m1()
    m2()
}
