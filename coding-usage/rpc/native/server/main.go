package main

import (
	"fmt"
	"github.com/luvx21/coding-go/coding-usage/rpc/native/service"
	"net"
	"net/rpc"
)

type Hello struct{}

func (h Hello) SayHi(req service.Person, res *service.Res) error {
	*res = service.Res{Msg: "success:" + fmt.Sprint(req), Code: 200}
	return nil
}

func main() {
	if err := rpc.RegisterName("hello", Hello{}); err != nil {
		fmt.Println(err)
		return
	}
	lister, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer lister.Close()
	for {
		fmt.Println("等待连接。。。。")
		// 3. 建立连接
		conn, err := lister.Accept()
		if err != nil {
			fmt.Println(err)
		}
		// 4. 绑定服务
		rpc.ServeConn(conn)
	}

}
