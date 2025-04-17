package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"testing"
	"time"
)

func Test01(t *testing.T) {
	write()
	readByReader()
}

func write() {
	// 创建一个writer 向top发送消息
	w := &kafka.Writer{
		Addr:                   kafka.TCP("localhost:19092"),
		Topic:                  "top",
		Balancer:               &kafka.LeastBytes{}, // 指定分区的balancer模式为最小字节分布
		RequiredAcks:           kafka.RequireAll,    // ack模式
		Async:                  true,                // 异步
		AllowAutoTopicCreation: true,
	}

	messages := []kafka.Message{
		{
			Key:   []byte("Key-A"),
			Value: []byte("Hello World!"),
		},
		{
			Key:   []byte("Key-B"),
			Value: []byte("One!"),
		},
		{
			Key:   []byte("Key-C"),
			Value: []byte("Two!"),
		},
	}

	var err error
	const retries = 3
	// 重试3次
	for range retries {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = w.WriteMessages(ctx, messages...)
		if errors.Is(err, context.DeadlineExceeded) {
			time.Sleep(time.Millisecond * 250)
			continue
		}

		break
	}

	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}

func readByReader() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:19092"},
		GroupID:   "consumer-group-id",
		Topic:     "top",
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})
	r.SetOffset(42) // 设置Offset

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("消息->offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
