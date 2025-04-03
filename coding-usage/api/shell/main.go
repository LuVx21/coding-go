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

		cmd := exec.Command("bash", "-c", cmdString)
		cmd.Stderr, cmd.Stdout = os.Stderr, os.Stdout
		fmt.Println("执行:", cmdString)
		err = cmd.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		bytes, err := exec.Command("bash", "-c", cmdString).
			// Output()
			CombinedOutput()
		fmt.Println("输出结果:", string(bytes))
	}
}
