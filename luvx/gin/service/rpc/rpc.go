package rpc

import (
	"log/slog"

	"luvx_service_sdk/proto_gen/proto_kv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	RpcConn     *grpc.ClientConn
	KvRpcClient proto_kv.KVClient
)

func init() {
	RpcConn, err := grpc.NewClient("rpc_service:18888", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("RPC连接失败", "err", err)
		return
	}
	KvRpcClient = proto_kv.NewKVClient(RpcConn)
}
