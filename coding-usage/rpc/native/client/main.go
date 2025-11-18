package main

import (
	"fmt"
	"github.com/luvx21/coding-go/coding-usage/rpc/native/service"
	"net/rpc"
)

func main() {
	// 1. 与rpc 服务器建立连接
	rclient, err := rpc.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rclient.Close()
	// 2. 调用远程方法
	var rep service.Res
	req := service.Person{Name: "名字", Age: 18}
	err = rclient.Call(service.HelloSayHi, req, &rep)
	if err != nil {
		fmt.Println("call fail:", err)
		return
	}
	fmt.Printf("%#v \n", rep)
}
