package etcds

import (
	"context"
	"time"

	mvccpb "go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

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
