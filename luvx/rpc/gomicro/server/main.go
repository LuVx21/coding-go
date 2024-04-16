package main

import (
    "context"
    "log"
    "luvx/api/logs"

    "go-micro.dev/v4"
    pb "luvx/rpc/gomicro/proto"
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
