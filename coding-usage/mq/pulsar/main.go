package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/luvx21/coding-go/infra/logs"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	uri   = "pulsar://localhost:6650"
	topic = "my-topic"
)

func main() {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               uri,
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})
	if err != nil {
		log.Fatalf("Could not instantiate Pulsar client: %v", err)
	}

	defer client.Close()

	go func() {
		prometheusPort := 2112
		logs.Log.Infof("开启Prometheus metrics接口-> http://localhost:%v/metrics\n", prometheusPort)
		http.Handle("/metrics", promhttp.Handler())
		_ = http.ListenAndServe(":"+strconv.Itoa(prometheusPort), nil)
	}()

	consumer(client)
	//reader(client)

	producer, _ := client.CreateProducer(pulsar.ProducerOptions{
		Topic: topic,
		Name:  "producer-name",
	})
	defer producer.Close()

	ctx := context.Background()
	webPort := 8082
	http.HandleFunc("/produce", func(w http.ResponseWriter, r *http.Request) {
		msg := fmt.Sprintf("%s\thello world", time.Now())
		msgId, err := producer.Send(ctx, &pulsar.ProducerMessage{
			Payload: []byte(msg),
		})
		if err != nil {
			log.Fatal(err)
		} else {
			log.Printf("发送消息: %v %s", msgId, msg)
			fmt.Fprintf(w, "发送消息: %v %s", msgId, msg)
		}
	})

	logs.Log.Infof("开启生产者接口-> http://localhost:%v/produce", webPort)
	err = http.ListenAndServe(":"+strconv.Itoa(webPort), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func reader(client pulsar.Client) {
	go func() {
		reader, _ := client.CreateReader(pulsar.ReaderOptions{
			Name:           "reader-name",
			Topic:          topic,
			StartMessageID: pulsar.LatestMessageID(),
		})
		defer reader.Close()

		for reader.HasNext() {
			if msg, err := reader.Next(context.Background()); err != nil {
				log.Fatal(err)
			} else {
				fmt.Printf("reader收到消息-> msgId: %s, 内容: %s\n", msg.ID().String(), string(msg.Payload()))
			}
		}
	}()
}

func consumer(client pulsar.Client) {
	go func() {
		options := pulsar.ConsumerOptions{
			Name:             "consumer-name-single",
			Topic:            topic,
			SubscriptionName: "my-sub",
			Type:             pulsar.Shared,
		}
		//options = pulsar.ConsumerOptions{
		//    Name:             "consumer-name-regex",
		//    TopicsPattern:    "persistent://public/default/my-*",
		//    SubscriptionName: "regex-sub",
		//}

		consumer, _ := client.Subscribe(options)
		defer consumer.Close()

		for i := 0; i < 100; i++ {
			if msg, err := consumer.Receive(context.Background()); err != nil {
				log.Fatal(err)
			} else {
				fmt.Printf("收到消息-> msgId: %s, 内容: %s\n", msg.ID().String(), string(msg.Payload()))
				consumer.Ack(msg)
			}
		}

		if err := consumer.Unsubscribe(); err != nil {
			log.Fatal(err)
		}
	}()
}
