package main

import (
	"fmt"
	"github.com/IBM/sarama" // // "github.com/segmentio/kafka-go"
	"log"
)

const (
	topic = "top"
)

func main() {
	config := sarama.NewConfig()
	config.Metadata.AllowAutoTopicCreation = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	// 连接kafka
	client, _ := sarama.NewSyncProducer([]string{"localhost:19092"}, config)
	defer client.Close()

	msg := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder("topicKey"),
		Topic: topic,
		Value: sarama.StringEncoder("Golang 测试消息"),
	}
	// 发送消息
	partition, offset, _ := client.SendMessage(msg)
	fmt.Printf("partition:%v offset:%v\n", partition, offset)

	// 创建Kafka消费者
	consumer, _ := sarama.NewConsumer([]string{"localhost:19092"}, nil)
	defer consumer.Close()

	partitionList, _ := consumer.Partitions(topic)
	for partition := range partitionList {
		pc, _ := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetOldest)
		defer pc.Close()

		// 异步消费
		//go func(sarama.PartitionConsumer) {
		for msg := range pc.Messages() {
			log.Printf("Partition: %d Offset: %d Key: %s Value: %s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
		}
		//}(pc)
	}
}
