package etcds

import (
	"log"
	"log/slog"
	"testing"

	"github.com/luvx21/coding-go/coding-common/common_x/runs"
)

const (
	serviceName = "user-service"
	addr        = "127.0.0.1:8080"
)

var (
	etcdEndpoints = []string{"etcd1:2379", "etcd2:2379", "etcd3:2379"}
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
