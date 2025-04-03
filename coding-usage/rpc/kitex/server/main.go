package main

import (
	"context"
	"log"
	"net"

	"github.com/cloudwego/kitex/server"
	"github.com/luvx21/coding-go/coding-usage/rpc/kitex/kitex_gen"
	"github.com/luvx21/coding-go/coding-usage/rpc/kitex/kitex_gen/greeter"
)

type GreeterServiceImpl struct {
}

func (s *GreeterServiceImpl) SayHello(ctx context.Context, req *kitex_gen.HelloRequest) (resp *kitex_gen.HelloReply, err error) {
	resp = new(kitex_gen.HelloReply)
	resp.Message = "hello " + req.Name
	return
}

func main() {
	addr, _ := net.ResolveTCPAddr("tcp", ":18888")
	svr := greeter.NewServer(new(GreeterServiceImpl), server.WithServiceAddr(addr))

	err := svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
