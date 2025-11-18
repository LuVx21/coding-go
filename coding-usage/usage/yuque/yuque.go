package main

import (
	"fmt"
)

const root = "/Users/renxie/OneDrive/Code/doc/"
const indexFile = root + "yp.csv"
const modifiedFiles = root + ".git/hooks/modified_files.txt"

func init() {
	fmt.Println("init")
	//content, _ := ioutil.ReadFile(path)
	//fmt.Println(string(content))
}

func main() {
	fmt.Println("a")
	//fmt.Println(yq.NewCreateBody())
	//fmt.Println(yq.CreateBody{Body: "a"})
}
