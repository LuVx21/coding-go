package mq

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/luvx21/coding-go/coding-common/os_x"
	"github.com/luvx21/coding-go/coding-common/times_x"
	"github.com/luvx21/coding-go/infra/nosql/infra_redis/mq"
	"github.com/redis/go-redis/v9"
)

var (
	ctx    = context.Background()
	uri, _ = os_x.Command("sh", "-c", "kv get redis_uri")
	opt, _ = redis.ParseURL(uri)
	rdb    = redis.NewClient(opt)

	stream = "test:task:queue"
	group  = "worker-group"

	msgId atomic.Int64
)

func date() {
	rdb.XTrimMaxLen(ctx, stream, 0)

	for {
		time.Sleep(time.Second * time.Duration(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(9)+1))

		t := times_x.TimeNowMicrosecond()
		rdb.XAdd(ctx, &redis.XAddArgs{
			Stream: stream,
			MaxLen: 10000,
			Values: map[string]any{
				"id":  msgId.Add(1),
				"now": t,
			},
		})
		fmt.Println("发送消息:", "id", msgId.Load(), t)
	}
}

func Test_stream_00(t *testing.T) {
	read := func(rdb *redis.Client, stream string) {
		lastID := "0" // "0" 表示从最早的消息开始, "$"表示从最新的消息开始
		for {
			results, err := rdb.XRead(ctx, &redis.XReadArgs{
				Streams: []string{stream, lastID},
				Count:   2,
				Block:   0, // 0 表示永久阻塞等待新消息
			}).Result()
			if err != nil || len(results) == 0 {
				if err == redis.Nil {
					continue
				}
				time.Sleep(time.Second)
				continue
			}
			for _, r := range results {
				for _, msg := range r.Messages {
					log.Printf("消费者 %s 消息Id:%v 消息内容: %v", "main", msg.ID, msg.Values)
					lastID = msg.ID
				}
			}
		}
	}

	go date()
	go read(rdb, stream)

	select {
	case <-time.After(time.Minute):
	}
}

func Test_stream_01(t *testing.T) {
	read := func(rdb *redis.Client, stream, consumer string) {
		if err := rdb.XGroupCreateMkStream(ctx, stream, group, "0").Err(); err != nil && !mq.IsGroupExistsError(err) {
			log.Fatal(err)
		}

		for {
			results, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    group,
				Consumer: consumer,
				Streams:  []string{stream, ">"}, // ">" 表示获取尚未分配给其他消费者的新消息
				Count:    5,                     // 每次最多读取5条
				Block:    0,                     // 0 表示永久阻塞等待新消息
			}).Result()

			if err != nil {
				if err == redis.Nil {
					continue
				}
				log.Printf("读取消息失败: %v", err)
				time.Sleep(time.Second)
				continue
			}

			for _, r := range results {
				for _, msg := range r.Messages {
					log.Printf("消费者 %s 消息Id:%v 消息内容: %v", group+"/"+consumer, msg.ID, msg.Values)

					if err := rdb.XAck(ctx, stream, group, msg.ID).Err(); err != nil {
						slog.Error("ACK 失败:", "err", err)
					}
				}
			}
		}
	}

	go date()
	go read(rdb, stream, fmt.Sprintf("consumer-%d", 1))
	go read(rdb, stream, fmt.Sprintf("consumer-%d", 2))

	select {
	case <-time.After(time.Minute):
	}
}
