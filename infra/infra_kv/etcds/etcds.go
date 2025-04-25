package etcds

import (
	"context"
	"time"

	mvccpb "go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func PutIfAbsent(cli *clientv3.Client, key, value string) (*clientv3.TxnResponse, bool) {
	tr, err := cli.Txn(context.Background()).
		If(clientv3.Compare(clientv3.Version(key), "=", 0)). // 判断 key 不存在（版本号为 0）
		Then(clientv3.OpPut(key, value)).                    // 不存在则 Put
		Else(clientv3.OpGet(key)).                           // 存在则忽略（这里可以选择不操作或返回当前值）
		Commit()
	return tr, err == nil && tr.Succeeded
}

func Set(cli *clientv3.Client, key, value string, ttl int64) error {
	leaseResp, err := cli.Grant(context.Background(), 10)
	if err != nil {
		return err
	}
	_, err = cli.Put(context.Background(), key, value, clientv3.WithLease(leaseResp.ID))
	return err
}

func IterateAll(cli *clientv3.Client) (map[string][]byte, error) {
	ctx := context.Background()
	var lastKey []byte
	r := make(map[string][]byte)

	for {
		resp, err := cli.Get(ctx, string(lastKey), clientv3.WithFromKey(), clientv3.WithLimit(100))
		if err != nil {
			return nil, err
		}
		for _, kv := range resp.Kvs {
			r[string(kv.Key)] = kv.Value
			lastKey = kv.Key
		}
		if !resp.More {
			break
		}
	}
	return r, nil
}
func IterateRange(cli *clientv3.Client, startKey, endKey string) (map[string][]byte, error) {
	return get(cli, startKey, clientv3.WithRange(endKey), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
}
func IterateWithPrefix(cli *clientv3.Client, prefix string) (map[string][]byte, error) {
	return get(cli, prefix, clientv3.WithPrefix())
}

func get(cli *clientv3.Client, key string, opts ...clientv3.OpOption) (map[string][]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := cli.Get(ctx, key, opts...)
	if err != nil {
		return nil, err
	}
	return convert(resp.Kvs), nil
}
func convert(kvs []*mvccpb.KeyValue) map[string][]byte {
	r := make(map[string][]byte, len(kvs))
	for _, kv := range kvs {
		r[string(kv.Key)] = kv.Value
	}
	return r
}
