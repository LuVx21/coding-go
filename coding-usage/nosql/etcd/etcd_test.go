package etcd

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/luvx21/coding-go/coding-common/test"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var cli *clientv3.Client

var beforeAfter = func(name string) func() {
	return test.BeforeAfterTest(name, func() {
		if cli == nil {
			var err error
			cli, err = clientv3.New(clientv3.Config{
				Endpoints:   []string{"etcd1:2379", "etcd2:2379", "etcd3:2379"},
				DialTimeout: 5 * time.Second,
			})
			if err != nil {
				log.Fatalf("failed to connect to etcd: %v", err)
			}
		}
	}, func() { cli.Close() })
}

func Test_etcd_00(t *testing.T) {
	defer beforeAfter("Test_etcd_00")()

	_, err := cli.Put(context.TODO(), "foo", "barbar")
	if err != nil {
		log.Fatal(err)
	}
}

func Test_etcd_01(t *testing.T) {
	defer beforeAfter("Test_etcd_01")()

	gr, err := cli.Get(context.TODO(), "foo")
	if err != nil {
		log.Fatal(err)
	}

	for _, kv := range gr.Kvs {
		fmt.Println(string(kv.Key), "=", string(kv.Value))
	}
}

func Test_etcd_02(t *testing.T) {
	defer beforeAfter("Test_etcd_02")()

	serviceNamePrefix := "/services/rpc-service/"
	watchCh := cli.Watch(context.Background(), serviceNamePrefix, clientv3.WithPrefix())
	for resp := range watchCh {
		for _, ev := range resp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				log.Printf("服务上线: %s", ev.Kv.Key)
			case clientv3.EventTypeDelete:
				log.Printf("服务下线: %s", ev.Kv.Key)
			}
		}
	}
}

func Test_etcd_03(t *testing.T) {
	defer beforeAfter("Test_etcd_03")()
	resp, _ := cli.Get(context.TODO(), "", clientv3.WithFromKey())
	for _, kv := range resp.Kvs {
		fmt.Printf("%s = %s\n", kv.Key, kv.Value)
	}
}
