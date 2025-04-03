package main

// import (
// 	"context"
// 	"log/slog"

// 	_ "dubbo.apache.org/dubbo-go/v3/imports"
// 	"dubbo.apache.org/dubbo-go/v3/protocol"
// 	"dubbo.apache.org/dubbo-go/v3/server"
// 	greet "github.com/luvx21/coding-go/coding-usage/rpc/dubbo/proto"
// )

// type GreeterService struct {
// }

// func (srv *GreeterService) SayHello(ctx context.Context, req *greet.HelloRequest) (*greet.HelloReply, error) {
// 	resp := &greet.HelloReply{Message: req.Name}
// 	return resp, nil
// }

// func main() {
// 	srv, err := server.NewServer(
// 		server.WithServerProtocol(
// 			protocol.WithPort(20000),
// 			protocol.WithTriple(),
// 		),
// 	)
// 	if err != nil {
// 		panic(err)
// 	}

// 	if err := greet.RegisterGreeterHandler(srv, &GreeterService{}); err != nil {
// 		panic(err)
// 	}

// 	if err := srv.Serve(); err != nil {
// 		slog.Error("异常", "err:", err)
// 	}
// }
