package main

import (
    "context"
    "flag"
    "fmt"
    "github.com/luvx21/coding-go/coding-common/logs"
    "google.golang.org/grpc/reflection"
    "log"
    "net"

    pb "github.com/luvx21/coding-go/coding-usage/rpc/grpc/proto"
    "google.golang.org/grpc"
)

var (
    port = flag.Int("port", 50051, "The server port")
)

type server struct {
    pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
    logs.Log.Printf("Received: %v", in.GetName())
    return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
    flag.Parse()
    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    reflection.Register(s)
    pb.RegisterGreeterServer(s, &server{})
    log.Printf("server listening at %v", lis.Addr())
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
