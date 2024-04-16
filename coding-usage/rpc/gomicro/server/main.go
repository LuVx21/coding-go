package main

import (
    "context"
    "github.com/luvx21/coding-go/coding-common/logs"
    "log"

    pb "github.com/luvx21/coding-go/coding-usage/rpc/gomicro/proto"
    "go-micro.dev/v4"
)

type Greeter struct{}

func (g *Greeter) SayHello(ctx context.Context, req *pb.HelloRequest, rsp *pb.HelloReply) error {
    logs.Log.Printf("收到请求参数-> %s", req.Name)
    rsp.Message = "Hello " + req.Name
    return nil
}

func main() {
    service := micro.NewService(
        micro.Name("go-micro-service"),
    )

    service.Init()

    _ = pb.RegisterGreeterHandler(service.Server(), new(Greeter))

    if err := service.Run(); err != nil {
        log.Fatal(err)
    }
}
