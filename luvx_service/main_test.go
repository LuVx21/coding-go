package main

import (
	"context"
	"fmt"
	kv "luvx_service_sdk/proto_gen/proto_kv"
	"testing"

	"github.com/luvx21/coding-go/coding-common/sets"

	"github.com/luvx21/coding-go/coding-common/test"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

var conn *grpc.ClientConn
var kvClient kv.KVClient

var beforeAfter = func(name string) func() {
	return test.BeforeAfterTest(name, func() {
		conn, _ = grpc.NewClient("localhost:18888", grpc.WithTransportCredentials(insecure.NewCredentials()))
		kvClient = kv.NewKVClient(conn)
	}, func() { conn.Close() })
}

func Test_00(t *testing.T) {
	defer beforeAfter("Test_00")()

	pr, _ := kvClient.Put(context.Background(), &kv.PutRequest{Entry: &kv.Entry{Key: "foo", Value: []byte("bar")}})
	fmt.Println(pr.Message)

	gr, _ := kvClient.Get(context.Background(), &kv.Key{Key: "foo"})
	fmt.Println("响应: ", string(gr.GetValue()))

	lr, _ := kvClient.List(context.Background(), &emptypb.Empty{})
	fmt.Println("---------------------------")
	for k, v := range lr.GetEntryMap() {
		fmt.Println(k, "=", string(v))
	}
	fmt.Println("---------------------------")
}

func Test_01(t *testing.T) {
	defer beforeAfter("Test_01")()

	lr, _ := kvClient.ListKeyByPrefix(context.Background(), &kv.Key{Key: ""})

	fmt.Println("---------------------------")
	set := sets.NewSet(lr.Keys...)
	fmt.Println(set.Len())
	fmt.Println("---------------------------")
}
func Test_02(t *testing.T) {
	defer beforeAfter("Test_02")()
	keys := []string{""}
	kvClient.Delete(t.Context(), &kv.Keys{Keys: keys})
}
