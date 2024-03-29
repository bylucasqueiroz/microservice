package kafkalib

import (
	"fmt"
	"os"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaProducer struct {
	wg  sync.WaitGroup
	pdc *kafka.Producer
}

type KafkaProducerConfig struct {
	Host     string
	Acks     string
	ClientId string
}

func NewKafkaProducer(cfg *KafkaProducerConfig) *KafkaProducer {
	fmt.Println("Start produce..")
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Host,
		"client.id":         cfg.ClientId,
		"acks":              cfg.Acks,
	})
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}
	return &KafkaProducer{
		pdc: p,
	}
}

func (k *KafkaProducer) Produce(topic string, value []byte) error {
	k.wg.Add(1)
	delivery_chan := make(chan kafka.Event, 10000)

	err := k.pdc.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
	},
		delivery_chan,
	)
	if err != nil {
		k.wg.Done()
		return fmt.Errorf("error %s", err)
	}

	e := <-delivery_chan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		k.wg.Done()
		close(delivery_chan)
		return fmt.Errorf("delivery failed: %v", m.TopicPartition.Error)
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
		k.wg.Done()
		close(delivery_chan)
		return nil
	}
}
