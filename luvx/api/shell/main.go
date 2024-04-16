package main

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "strings"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Print("$ ")
        cmdString, err := reader.ReadString('\n')
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
        }
        cmdString = strings.TrimSuffix(cmdString, "\n")

        cmd := exec.Command(cmdString)
        cmd.Stderr = os.Stderr
        cmd.Stdout = os.Stdout
        fmt.Println("执行:", cmdString)
        err = cmd.Run()
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
        }
    }
}
