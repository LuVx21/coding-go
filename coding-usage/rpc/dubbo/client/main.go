package main

// import (
// 	"context"
// 	"log/slog"

// 	"dubbo.apache.org/dubbo-go/v3/client"
// 	_ "dubbo.apache.org/dubbo-go/v3/imports"
// 	greet "github.com/luvx21/coding-go/coding-usage/rpc/dubbo/proto"
// )

// func main() {
// 	cli, err := client.NewClient(
// 		client.WithClientURL("127.0.0.1:20000"),
// 	)
// 	if err != nil {
// 		panic(err)
// 	}

// 	svc, err := greet.NewGreeter(cli)
// 	if err != nil {
// 		panic(err)
// 	}

// 	resp, err := svc.SayHello(context.Background(), &greet.HelloRequest{Name: "hello world"})
// 	if err != nil {
// 		slog.Error("异常", "err:", err)
// 	}
// 	slog.Info("Greet response: %s", resp.Message)
// }
