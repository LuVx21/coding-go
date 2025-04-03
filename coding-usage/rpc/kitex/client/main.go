package main

import (
	"context"
	"fmt"

	"github.com/cloudwego/kitex/client"
	"github.com/luvx21/coding-go/coding-usage/rpc/kitex/kitex_gen"
	"github.com/luvx21/coding-go/coding-usage/rpc/kitex/kitex_gen/greeter"
)

func main() {
	c, _ := greeter.NewClient("rpc.hello", client.WithHostPorts(":18888"))
	r, _ := c.SayHello(context.Background(), &kitex_gen.HelloRequest{Name: "kitex"})
	message := r.Message
	fmt.Println(message)
}
