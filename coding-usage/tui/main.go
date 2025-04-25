package main

import (
	"fmt"

	"github.com/c-bata/go-prompt"
)

func executor(in string) {
	fmt.Println("你输入的是:", in)
	p := prompt.New(func(ii string) {
		fmt.Println("内层你输入的是:", ii)
	}, completer, prompt.OptionPrefix("$$$ "), prompt.OptionTitle("中文输入测试"))
	p.Run()
}

func completer(in prompt.Document) []prompt.Suggest {
	return nil
}

func main() {
	m1()
	m2()
}

func m1() {
	fmt.Println("Please select table.")
	in := prompt.Input("> ", completer)
	fmt.Println("You selected " + in)
}

func m2() {
	p := prompt.New(executor, completer, prompt.OptionPrefix(">>> "), prompt.OptionTitle("中文输入测试"))
	p.Run()
}
