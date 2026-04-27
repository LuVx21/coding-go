package mq

import (
	"context"
	"fmt"
	"log"
	"maps"
	"time"

	"github.com/redis/go-redis/v9"
)

type (
	// Consumer 消费者定义
	Consumer struct {
		name      string // 同一组内要有唯一性
		stream    string
		group     string
		rdb       *redis.Client
		ctx       context.Context
		stopChan  chan struct{}
		onMessage func(msg Message) error
		idleTime  time.Duration // 判定消息超时的时间阈值
	}
	// Message 消息结构
	Message struct {
		main     redis.XMessage
		Stream   string
		Group    string
		Consumer string
	}
)

func newConsumer(rdb *redis.Client, stream, group, consumerID string, ctx context.Context) *Consumer {
	return &Consumer{
		name:      consumerID,
		stream:    stream,
		group:     group,
		rdb:       rdb,
		ctx:       ctx,
		stopChan:  make(chan struct{}),
		onMessage: nil,
		idleTime:  30 * time.Second,
	}
}

// ClaimStuckMessages 认领超时消息
func (c *Consumer) ClaimStuckMessages() ([]redis.XMessage, error) {
	ctx := context.Background()

	// 1. 查看 Pending 列表，找出超时的消息
	pending, err := c.rdb.XPendingExt(ctx, &redis.XPendingExtArgs{
		Stream: c.stream,
		Group:  c.group,
		Start:  "-",        // 最早的消息
		End:    "+",        // 最新的消息
		Count:  100,        // 每次最多检查100条
		Idle:   c.idleTime, // 只找出闲置超过阈值（未确认）的消息
	}).Result()

	if err != nil || len(pending) == 0 {
		return nil, err
	}

	// 2. 提取超时消息ID
	messageIDs := make([]string, len(pending))
	for i, p := range pending {
		messageIDs[i] = p.ID
		fmt.Printf("[扫描] 检测到超时消息: %s (闲置时间: %v)\n", p.ID, p.Idle)
	}

	// 3. XCLAIM 将这些消息转移给当前消费者
	claimed, err := c.rdb.XClaim(ctx, &redis.XClaimArgs{
		Stream:   c.stream,
		Group:    c.group,
		Consumer: c.name,
		MinIdle:  c.idleTime,
		Messages: messageIDs,
	}).Result()

	if err != nil {
		return nil, err
	}

	if len(claimed) > 0 {
		fmt.Printf("[认领] %s 成功认领 %d 条超时消息\n", c.name, len(claimed))
	}

	return claimed, nil
}

// Consume 主消费循环
// 0: 消费者组从最早的消息开始消费
// $: 消费者组只消费新消息（不处理历史）
// >: 读取时只获取未分配给任何消费者的新消息
func (c *Consumer) consume() {
	go c.claimLoop()

	for {
		select {
		case <-c.stopChan:
			log.Printf("消费者 %s 停止读取消息", c.name)
			return
		default:
			results, err := c.rdb.XReadGroup(c.ctx, &redis.XReadGroupArgs{
				Group:    c.group,
				Consumer: c.name,
				Streams:  []string{c.stream, ">"},
				Count:    10,
				Block:    time.Second * 2,
				NoAck:    false,
			}).Result()

			if err != nil || len(results) == 0 {
				if err == redis.Nil {
					continue
				}
				time.Sleep(time.Second)
				continue
			}

			for _, stream := range results {
				for _, msg := range stream.Messages {
					select {
					case <-c.stopChan:
						continue
					default:
						msg := Message{
							main:     msg,
							Stream:   c.stream,
							Group:    c.group,
							Consumer: c.name,
						}
						c.processMessage(msg)
					}
				}
			}
		}
	}
}

// claimLoop 定时检查并认领超时消息
func (c *Consumer) claimLoop() {
	ticker := time.NewTicker(10 * time.Second) // 每10秒检查一次
	defer ticker.Stop()

	for range ticker.C {
		if claimed, err := c.ClaimStuckMessages(); err == nil {
			// 认领到的超时消息需要重新处理
			for _, msg := range claimed {
				fmt.Printf("[重试] 重新处理认领的消息: %s\n", msg.ID)
				c.processMessage(Message{
					main:     msg,
					Stream:   c.stream,
					Group:    c.group,
					Consumer: c.name,
				})
			}
		}
	}
}

func (c *Consumer) processMessage(msg Message) {
	fmt.Printf("[%s] 开始消费: ID=%s, 内容=%v\n", msg.Consumer, msg.main.ID, msg.main.Values)
	// 从消息中获取重试次数（自定义字段）
	retryCount := 0
	if val, ok := msg.main.Values["retry_count"]; ok {
		if count, ok := val.(int64); ok {
			retryCount = int(count)
		}
	}

	// 调用用户自定义的消息处理函数
	if err := c.onMessage(msg); err != nil {
		retryCount++
		fmt.Printf("[失败] 消息 %s 处理失败，重试次数: %d，错误: %v\n", msg.main.ID, retryCount, err)

		if retryCount >= 3 {
			// 超过最大重试次数，发送到死信队列或记录日志
			fmt.Printf("[丢弃] 消息 %s 超过最大重试次数，移入死信队列\n", msg.main.ID)
			c.sendToDeadLetterQueue(msg.main)
			c.rdb.XAck(context.Background(), c.stream, c.group, msg.main.ID)
			return
		}

		c.updateRetryCount(msg.main, retryCount)
		return
	}
	// 处理成功，发送 ACK
	if err := c.rdb.XAck(c.ctx, c.stream, c.group, msg.main.ID).Err(); err != nil {
		log.Printf("消费者 %s 确认消息失败: %v", c.name, err)
	}

}

// 停止消费者
func (c *Consumer) stop() {
	close(c.stopChan)
}

// updateRetryCount 更新消息的重试次数
func (c *Consumer) updateRetryCount(msg redis.XMessage, retryCount int) {
	ctx := c.ctx

	values := make(map[string]any, len(msg.Values))
	maps.Copy(values, msg.Values)
	values["retry_count"] = retryCount

	newID, _ := c.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: c.stream,
		Values: values,
	}).Result()

	c.rdb.XAck(ctx, c.stream, c.group, msg.ID)
	c.rdb.XDel(ctx, c.stream, msg.ID)

	fmt.Printf("[重试更新] 消息 %s 已更新重试次数为 %d，新ID: %s\n", msg.ID, retryCount, newID)
}

// sendToDeadLetterQueue 死信队列处理
func (c *Consumer) sendToDeadLetterQueue(msg redis.XMessage) {
	c.rdb.XAdd(c.ctx, &redis.XAddArgs{
		Stream: c.stream + ":dead",
		Values: map[string]any{
			"original_id":   msg.ID,
			"original_data": msg.Values,
			"reason":        "max_retries_exceeded",
			"failed_at":     time.Now().Unix(),
		},
	})
}
