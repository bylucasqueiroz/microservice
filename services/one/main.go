package main

import (
	"fmt"
	"os"
	"sync"

	csm "github.com/bylucasqueiroz/microservice/libs/kafka"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	fmt.Println("Start process..")

	var wg sync.WaitGroup
	host := os.Getenv("KAFKA_ENDPOINT")

	wg.Add(1)
	go func() {
		fmt.Println("Start produce..")
		for i := 0; i < 10; i++ {
			p, err := kafka.NewProducer(&kafka.ConfigMap{
				"bootstrap.servers": host,
				"client.id":         "teste",
				"acks":              "all",
			})
			if err != nil {
				fmt.Printf("Failed to create producer: %s\n", err)
				os.Exit(1)
			}
			delivery_chan := make(chan kafka.Event, 10000)

			topic := "topic"
			err = p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          []byte("TESTE"),
			},
				delivery_chan,
			)
			if err != nil {
				fmt.Errorf("error %s", err)
				wg.Done()
			}

			e := <-delivery_chan
			m := e.(*kafka.Message)

			if m.TopicPartition.Error != nil {
				fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
			} else {
				fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
					*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
			}
			close(delivery_chan)
		}
	}()
	wg.Wait()

	consume := csm.NewKafkaConsumer(&csm.KafkaConsumerConfig{
		Host:            host,
		GroupId:         "group1",
		AutoOffsetReset: "earliest",
		Topics:          []string{"teste"},
	})

	go consume.Consume()

	consume.WaitConsumerDone()
}
