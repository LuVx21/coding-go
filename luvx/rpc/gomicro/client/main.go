package main

import (
    "context"
    "fmt"

    "go-micro.dev/v4/client"
    pb "luvx/rpc/gomicro/proto"
)

func main() {
    cl := pb.NewGreeterService("go-micro-service", client.DefaultClient)

    rsp, err := cl.SayHello(context.TODO(), &pb.HelloRequest{Name: "foo-bar"})
    if err != nil {
        fmt.Println(err)
    }

    fmt.Println(rsp.Message)
}
