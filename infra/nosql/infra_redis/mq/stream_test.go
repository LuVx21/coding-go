package mq

import (
	"context"
	"log"
	"math/rand"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/luvx21/coding-go/coding-common/fmt_x"
	"github.com/luvx21/coding-go/coding-common/os_x"
	"github.com/luvx21/coding-go/coding-common/times_x"
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
		fmt_x.Debugln("发送消息: id=", msgId.Load(), t)
	}
}

func Test_stream_00(t *testing.T) {
	go date()

	group := NewConsumerGroup(rdb, stream, group)
	group.InitGroup()

	consumerNames := []string{"consumer-1", "consumer-2", "consumer-3"}
	for _, name := range consumerNames {
		group.CreateConsumer(name)
	}

	line := strings.Repeat("-", 80)
	messageHandler := func(msg Message) error {
		fmt_x.Infof("[%s] 处理消息: ID=%s, 内容=%v\n%s\n", msg.Consumer, msg.main.ID, msg.main.Values, line)
		// 模拟业务处理
		time.Sleep(time.Millisecond * 50)
		return nil
	}

	group.StartAllConsumers(messageHandler)
	log.Printf("消费组启动完成，有 %d 个消费者\n%s\n", group.GetConsumerCount(), line)

	// 模拟运行一段时间
	time.Sleep(time.Second * 30)

	// 7. 优雅关闭
	group.StopAll()
	log.Println("程序退出")
}
