package mq

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

type (
	ConsumerGroup struct {
		rdb       *redis.Client
		stream    string
		group     string
		consumers map[string]*Consumer
		mu        sync.RWMutex
		wg        sync.WaitGroup

		ctx    context.Context
		cancel context.CancelFunc
	}
)

// NewConsumerGroup 创建消费组管理器
func NewConsumerGroup(rdb *redis.Client, stream, group string) *ConsumerGroup {
	ctx, cancel := context.WithCancel(context.Background())
	return &ConsumerGroup{
		rdb:       rdb,
		stream:    stream,
		group:     group,
		consumers: make(map[string]*Consumer),
		ctx:       ctx,
		cancel:    cancel,
	}
}

// InitGroup 初始化消费者组（服务启动时调用一次）
func (cg *ConsumerGroup) InitGroup() error {
	// XGROUP CREATE stream group 0 MKSTREAM
	// 0：从最早消息开始；$：只处理新消息
	err := cg.rdb.XGroupCreateMkStream(context.Background(), cg.stream, cg.group, "0").Err()
	if err != nil && !IsGroupExistsError(err) {
		return nil
	}
	return err
}

// CreateConsumer 创建消费者（不启动消费）
func (cg *ConsumerGroup) CreateConsumer(consumerID string) *Consumer {
	cg.mu.Lock()
	defer cg.mu.Unlock()

	if _, exists := cg.consumers[consumerID]; exists {
		log.Printf("消费者 %s 已存在", consumerID)
		return cg.consumers[consumerID]
	}

	consumer := newConsumer(cg.rdb, cg.stream, cg.group, consumerID, cg.ctx)
	cg.consumers[consumerID] = consumer
	return consumer
}

// StartConsumer 启动消费者消费
func (cg *ConsumerGroup) StartConsumer(consumerID string, onMessage func(msg Message) error) error {
	consumer, exists := func(i string) (*Consumer, bool) {
		cg.mu.RLock()
		defer cg.mu.RUnlock()
		a, b := cg.consumers[i]
		return a, b
	}(consumerID)
	if !exists {
		return fmt.Errorf("消费者 %s 不存在", consumerID)
	}

	consumer.onMessage = onMessage

	cg.wg.Go(func() {
		log.Printf("启动消费者消费: %s", consumer.name)
		consumer.consume()
	})

	return nil
}

// StopConsumer 停止指定消费者
func (cg *ConsumerGroup) StopConsumer(consumerID string) {
	consumer, exists := func(i string) (*Consumer, bool) {
		cg.mu.RLock()
		defer cg.mu.RUnlock()
		a, b := cg.consumers[i]
		return a, b
	}(consumerID)

	if exists {
		consumer.stop()
		log.Printf("消费者 %s 已停止", consumerID)
	}
}

// StartAllConsumers 启动所有消费者
func (cg *ConsumerGroup) StartAllConsumers(onMessage func(msg Message) error) {
	cg.mu.RLock()
	defer cg.mu.RUnlock()

	for _, consumer := range cg.consumers {
		consumer.onMessage = onMessage
		cg.wg.Go(func() {
			log.Printf("启动消费者消费: %s", consumer.name)
			consumer.consume()
		})
	}
}

// StopAll 停止所有消费者
func (cg *ConsumerGroup) StopAll() {
	cg.cancel()

	cg.mu.RLock()
	defer cg.mu.RUnlock()

	for _, consumer := range cg.consumers {
		consumer.stop()
	}

	cg.wg.Wait()
	log.Println("所有消费者已停止")
}

// GetConsumerCount 获取消费者数量
func (cg *ConsumerGroup) GetConsumerCount() int {
	cg.mu.RLock()
	defer cg.mu.RUnlock()
	return len(cg.consumers)
}

// GetPendingMessages 获取未确认消息数量
// func (cg *ConsumerGroup) GetPendingMessages(consumerID string) (int64, error) {
// 	consumer, exists := func(i string) (*Consumer, bool) {
// 		cg.mu.RLock()
// 		defer cg.mu.RUnlock()
// 		a, b := cg.consumers[i]
// 		return a, b
// 	}(consumerID)
// 	if !exists {
// 		return 0, fmt.Errorf("消费者 %s 不存在", consumerID)
// 	}

// 	pending, err := consumer.rdb.XPendingExt(cg.ctx, &redis.XPendingExtArgs{
// 		Stream: cg.stream,
// 		Group:  cg.group,
// 		Start:  "-",
// 		End:    "+",
// 		Count:  0, // 获取所有
// 	}).Result()
// 	if err != nil {
// 		return 0, err
// 	}

// 	var count int64
// 	for _, p := range pending {
// 		if p.Consumer == consumerID {
// 			count++
// 		}
// 	}

// 	return count, nil
// }
