package main

import (
	"log"
	"log/slog"
	"net"
	"time"

	"luvx_service/service"
	kv "luvx_service_sdk/proto_gen/proto_kv"

	"github.com/luvx21/coding-go/coding-common/common_x/runs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	svr := grpc.NewServer()
	reflection.Register(svr)

	kv.RegisterKVServer(svr, new(service.KVServiceImpl))

	lis, _ := net.Listen("tcp", ":18888")
	go func() {
		slog.Info("Server running...", "port", lis.Addr().String())
		if err := svr.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	runs.GracefulStop(nil, func() {
		gracefulStop(svr)
	})
	slog.Info("Server stopped gracefully")
}

func gracefulStop(grpcServer *grpc.Server) {
	// 1. 首先停止接受新连接
	grpcServer.GracefulStop()

	// 2. 执行你的清理逻辑
	cleanup()

	// 3. 如果需要，可以等待一段时间让现有请求完成
	// 注意: GracefulStop() 已经会等待现有 RPC 完成
	// 这里只是示例额外的等待
	time.Sleep(2 * time.Second)
}

func cleanup() {
	service.CleanUpKv()
}
