package main

import (
	"context"
	"fmt"
	"log"
	kv "luvx_service_sdk/proto_gen/proto_kv"
	"math/rand"
	"testing"
	"time"

	"github.com/luvx21/coding-go/coding-common/test"
	"github.com/luvx21/coding-go/infra/infra_kv/etcds"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	conn       *grpc.ClientConn
	kvClient   kv.KVClient
	serviceKey = etcds.ServiceKey(keyPrefix, serviceName)
)

var beforeAfter = func(name string) func() {
	return test.BeforeAfterTest(name, func() {
		addrs, err := etcds.DiscoverServiceV1(endpoints, serviceKey+"/")
		if err != nil || len(addrs) == 0 {
			log.Fatal("未找到可用服务")
		}

		conn, _ = grpc.NewClient(addrs[rand.Intn(len(addrs))], grpc.WithTransportCredentials(insecure.NewCredentials()))
		kvClient = kv.NewKVClient(conn)
	}, func() { conn.Close() })
}

var beforeAfter1 = func(name string) func() {
	return test.BeforeAfterTest(name, func() {
		cli, err := clientv3.New(clientv3.Config{
			Endpoints:   endpoints,
			DialTimeout: 5 * time.Second,
		})
		if err != nil {
			log.Fatal(err)
		}
		// defer cli.Close()

		// 服务必须是使用RegisterEndpoint方式注册
		resolverBuilder, err := resolver.NewBuilder(cli)
		if err != nil {
			log.Fatalf("create resolver failed: %v", err)
		}
		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithResolvers(resolverBuilder),
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		}

		conn, _ = grpc.NewClient(fmt.Sprintf("etcd:///%s", serviceKey), opts...)
		kvClient = kv.NewKVClient(conn)
	}, func() { conn.Close() })
}

func Test_00(t *testing.T) {
	defer beforeAfter("Test_00")()

	pr, _ := kvClient.Put(context.Background(), &kv.PutRequest{Entry: &kv.Entry{Key: "foo", Value: []byte("bar")}})
	fmt.Println(pr.Message)
}

func Test_01(t *testing.T) {
	defer beforeAfter("Test_01")()

	gr, _ := kvClient.Get(context.Background(), &kv.Key{Key: "foo"})
	fmt.Println("响应: ", string(gr.GetValue()))
}

func Test_02(t *testing.T) {
	defer beforeAfter1("Test_02")()

	gr, _ := kvClient.Get(context.Background(), &kv.Key{Key: "foo"})
	fmt.Println("响应: ", string(gr.GetValue()))
}
