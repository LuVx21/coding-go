package main

import (
    "flag"
    "fmt"
    "os"
)

func main() {
    if len(os.Args) <= 1 {
        fmt.Println("无子命令, 可用子命令有: cmd1, cmd2")
        os.Exit(1)
    }

    switch os.Args[1] {
    case "cmd1":
        m1()
    case "cmd2":
        m2()
    default:
        fmt.Println("未知命令:", os.Args[1])
        os.Exit(1)
    }
}

func m1() {
    cmd1 := flag.NewFlagSet("cmd1", flag.ExitOnError)
    search := cmd1.String("search", "search默认值", "search参数说明")
    foo := cmd1.String("foo", "foo默认值", "foo参数说明")
    _ = cmd1.Parse(os.Args[2:])

    fmt.Printf("子命令: cmd1, 参数search: %s, foo: %s\n", *search, *foo)
}

func m2() {
    cmd2 := flag.NewFlagSet("cmd2", flag.ExitOnError)
    search := cmd2.String("search", "search默认值", "search参数说明")
    bar := cmd2.String("bar", "bar默认值", "bar参数说明")
    _ = cmd2.Parse(os.Args[2:])

    fmt.Printf("子命令: cmd2, 参数search: %s, foo: %s\n", *search, *bar)
}
