package main

import (
    "context"
    "fmt"

    pb "github.com/luvx21/coding-go/coding-usage/rpc/gomicro/proto"
    "go-micro.dev/v4/client"
)

func main() {
    cl := pb.NewGreeterService("go-micro-service", client.DefaultClient)

    rsp, err := cl.SayHello(context.TODO(), &pb.HelloRequest{Name: "foo-bar"})
    if err != nil {
        fmt.Println(err)
    }

    fmt.Println(rsp.Message)
}
