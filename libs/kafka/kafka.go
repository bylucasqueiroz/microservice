package kafkalib

import (
	"fmt"
	"os"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaConsumer struct {
	wg     sync.WaitGroup
	csm    *kafka.Consumer
	topics []string
}

type KafkaConsumerConfig struct {
	Host            string
	GroupId         string
	AutoOffsetReset string
	Topics          []string
}

func NewKafkaConsumer(cfg *KafkaConsumerConfig) *KafkaConsumer {
	csm, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Host,
		"group.id":          cfg.GroupId,
		"auto.offset.reset": cfg.AutoOffsetReset,
	})
	if err != nil {
		fmt.Printf("Failed to create consumer: %s\n", err)
		os.Exit(1)
	}
	err = csm.SubscribeTopics(cfg.Topics, nil)
	if err != nil {
		fmt.Printf("Subscribed topics error: %s\n", err)
		os.Exit(1)
	}
	return &KafkaConsumer{
		csm:    csm,
		topics: cfg.Topics,
	}
}

func (c *KafkaConsumer) Consume() {
	c.wg.Add(1)
	run := true
	for run {
		ev := c.csm.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			fmt.Printf("Consumed %v\n", string(e.Value))
		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			run = false
			c.wg.Done()
		}
	}
	c.csm.Close()
}

func (c *KafkaConsumer) WaitConsumerDone() {
	c.wg.Wait()
}
