package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"luvx_service/service"
	kv "luvx_service_sdk/proto_gen/proto_kv"
	"net"
	"time"

	"github.com/luvx21/coding-go/coding-common/common_x/runs"
	"github.com/luvx21/coding-go/infra/infra_kv/etcds"
	"google.golang.org/grpc"
)

const (
	serviceName = "rpc-luvx"
	keyPrefix   = "rpc-services"
)

var (
	rpc_port  = flag.Int("port", 28888, "listening port")
	endpoints = []string{"etcd1:2379", "etcd2:2379", "etcd3:2379"}
)

func main() {
	flag.Parse()

	svr := grpc.NewServer()

	kv.RegisterKVServer(svr, new(service.KVServiceImpl))

	lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", *rpc_port))
	go func() {
		slog.Info("RPC Server running...", "port", lis.Addr().String())
		if err := svr.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	addr := fmt.Sprintf("localhost:%d", *rpc_port)
	ser, _ := etcds.NewServiceRegister(endpoints, keyPrefix, serviceName, addr, 5)
	if err := ser.Register(); err != nil {
		log.Fatalf("服务注册失败: %v", err)
	}
	ser.ListenLeaseRespChan()

	runs.GracefulStop(nil, func() {
		gracefulStop(ser, svr)
	})
	slog.Info("RPC Server stopped gracefully")
}

func gracefulStop(ser *etcds.ServiceRegister, grpcServer *grpc.Server) {
	if err := ser.Close(); err != nil {
		slog.Warn("服务关闭失败", "err", err)
	}

	// 停止接受新连接
	grpcServer.GracefulStop()

	// 执行你的清理逻辑
	cleanup()

	// 如果需要，可以等待一段时间让现有请求完成
	// 注意: GracefulStop() 已经会等待现有 RPC 完成
	time.Sleep(2 * time.Second)
}

func cleanup() {
	service.CleanUpKv()
}
