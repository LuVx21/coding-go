package etcds

import (
	"fmt"
	"log"
	"log/slog"
	"testing"
	"time"

	"github.com/luvx21/coding-go/coding-common/common_x/runs"
	"github.com/luvx21/coding-go/coding-common/test"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	serviceName = "user-service"
	addr        = "127.0.0.1:8080"
)

var (
	etcdEndpoints = []string{"etcd1:2379", "etcd2:2379", "etcd3:2379"}
	cli           *clientv3.Client
	beforeAfter   = func(name string) func() {
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
)

func Test_etcd_00(t *testing.T) {
	ser, _ := NewServiceRegister(etcdEndpoints, "services", serviceName, addr, 5)
	if err := ser.RegisterEndpoint(); err != nil {
		log.Fatalf("服务注册失败: %v", err)
	}
	ser.ListenLeaseRespChan()

	runs.GracefulStop(nil, func() {
		slog.Info("接收到退出信号")
	})

	// 关闭服务
	if err := ser.Close(); err != nil {
		slog.Warn("服务关闭失败", "err", err)
	}
	slog.Info("服务已关闭")
}

func Test_etcd_01(t *testing.T) {
	defer beforeAfter("Test_etcd_01")()
	cli.Delete(t.Context(), "a")

	r, _ := PutIfAbsent(cli, "a", "a")
	gr, _ := cli.Get(t.Context(), "a")
	fmt.Println(r, string(gr.Kvs[0].Value))

	r, _ = PutIfAbsent(cli, "a", "aa")
	gr, _ = cli.Get(t.Context(), "a")
	fmt.Println(string(r.Responses[0].GetResponseRange().Kvs[0].Value), string(gr.Kvs[0].Value))
}

func Test_lock_00(t *testing.T) {
	defer beforeAfter("Test_lock_00")()

	sleep := time.Second * 5
	locker := NewLocker[string](cli)
	go locker.LockRun("lock_foo", time.Second*30, func() {
		fmt.Println("加锁成功，执行任务1")
		time.Sleep(sleep)
		fmt.Println("任务执行完成1")
	})
	time.Sleep(time.Second * 1)
	go locker.LockRun("lock_foo", time.Second*30, func() {
		fmt.Println("加锁成功，执行任务2")
		time.Sleep(sleep)
		fmt.Println("任务执行完成2")
	})

	time.Sleep(sleep * 2)
}
