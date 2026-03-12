package rpc

import (
	"log/slog"
	"slices"

	"github.com/luvx21/coding-go/luvx_service_sdk/proto_gen/proto_kv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	RpcConn     *grpc.ClientConn
	KvRpcClient *proto_kv.KVClient
)

func init() {
	RpcConn, err := grpc.NewClient("rpc_service:18888", grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(50*1024*1024),
			grpc.MaxCallSendMsgSize(50*1024*1024),
		),
	)
	work := []connectivity.State{connectivity.Connecting, connectivity.Ready}
	if err != nil || !slices.Contains(work, RpcConn.GetState()) {
		slog.Error("RPC连接失败", "err", err, "状态", RpcConn.GetState().String())
		return
	}
	a := proto_kv.NewKVClient(RpcConn)
	KvRpcClient = &a
}
