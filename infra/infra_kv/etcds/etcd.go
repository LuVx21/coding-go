package etcds

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/luvx21/coding-go/coding-common/common_x"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

// ServiceRegister 服务注册结构体
type ServiceRegister struct {
	keyPrefix   string // key前缀, 前后不带"/"
	serviceName string // 服务名, 完整的key: keyPrefix/serviceName/addr
	addr        string // 服务注册的addr
	ttl         int64  // 租约时间(秒)

	client        *clientv3.Client
	leaseID       clientv3.LeaseID
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
}

// NewServiceRegister 创建服务注册实例
func NewServiceRegister(etcdEndpoints []string, keyPrefix, serviceName, addr string, ttl int64) (*ServiceRegister, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdEndpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	ser := &ServiceRegister{
		keyPrefix:   common_x.IfThen(keyPrefix != "", keyPrefix, "service"),
		serviceName: serviceName,
		addr:        addr,
		ttl:         ttl,
		client:      cli,
	}
	return ser, nil
}
func (s *ServiceRegister) serviceKey() string {
	return ServiceKey(s.keyPrefix, s.serviceName)
}
func (s *ServiceRegister) appKey() string {
	return AppKey(s.keyPrefix, s.serviceName, s.addr)
}

func (s *ServiceRegister) RegisterEndpoint() error {
	// 创建租约
	leaseResp, err := s.client.Grant(context.Background(), s.ttl)
	if err != nil {
		return err
	}

	em, err := endpoints.NewManager(s.client, s.serviceKey())
	if err != nil {
		return err
	}
	// 注册服务端点
	endpoint := endpoints.Endpoint{
		Addr: s.addr,
		Metadata: map[string]any{
			"keyPrefix":   s.keyPrefix,
			"serviceName": s.serviceName,
			"version":     "1.0.0",
		},
	}
	err = em.AddEndpoint(context.Background(), s.appKey(), endpoint, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}

	// 设置续租 定期发送心跳请求
	keepAliveChan, err := s.client.KeepAlive(context.Background(), leaseResp.ID)
	if err != nil {
		return err
	}

	s.leaseID = leaseResp.ID
	s.keepAliveChan = keepAliveChan

	slog.Info("服务注册成功", "key", s.appKey(), "address", s.addr)
	return nil
}

// Register 注册服务并设置租约
func (s *ServiceRegister) Register() error {
	// 创建租约
	leaseResp, err := s.client.Grant(context.Background(), s.ttl)
	if err != nil {
		return err
	}

	// 注册服务并绑定租约
	_, err = s.client.Put(context.Background(), s.appKey(), s.addr, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}

	// 设置续租 定期发送心跳请求
	keepAliveChan, err := s.client.KeepAlive(context.Background(), leaseResp.ID)
	if err != nil {
		return err
	}

	s.leaseID = leaseResp.ID
	s.keepAliveChan = keepAliveChan

	slog.Info("服务注册成功", "key", s.appKey(), "address", s.addr)
	return nil
}

// ListenLeaseRespChan 监听续租应答
func (s *ServiceRegister) ListenLeaseRespChan() {
	go func() {
		for leaseKeepResp := range s.keepAliveChan {
			if leaseKeepResp == nil {
				slog.Info("租约失效, 续约......")
				if err := s.Register(); err != nil {
					slog.Warn("租约续期失败", "err", err)
				}
				return
			}
		}
	}()
}

// Close 关闭服务
func (s *ServiceRegister) Close() error {
	// 撤销租约
	if _, err := s.client.Revoke(context.Background(), s.leaseID); err != nil {
		return err
	}
	// s.client.Delete(context.TODO(), s.key)
	return s.client.Close()
}

func DiscoverServiceV1(endpoints []string, serviceNamePrefix string) ([]string, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	return DiscoverService(cli, serviceNamePrefix)
}

func DiscoverService(cli *clientv3.Client, serviceNamePrefix string) ([]string, error) {
	resp, err := cli.Get(context.TODO(), serviceNamePrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	var addrs []string
	for _, kv := range resp.Kvs {
		addrs = append(addrs, string(kv.Value))
	}

	return addrs, nil
}

func WatchService(etcdClient *clientv3.Client, serviceNamePrefix string, onPut, onDelete func()) {
	watchCh := etcdClient.Watch(context.Background(), serviceNamePrefix, clientv3.WithPrefix())
	for resp := range watchCh {
		for _, ev := range resp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				if onPut != nil {
					onPut()
				}
			case clientv3.EventTypeDelete:
				if onDelete != nil {
					onDelete()
				}
			}
		}
	}
}

func ServiceKey(keyPrefix, serviceName string) string {
	return keyPrefix + "/" + serviceName
}
func AppKey(keyPrefix, serviceName, addr string) string {
	return ServiceKey(keyPrefix, serviceName) + "/" + addr
}
