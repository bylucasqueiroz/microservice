package main

import (
	"fmt"
	"os"
	"sync"

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

	wg.Add(1)
	go func() {
		fmt.Println("Start consume..")
		topics := []string{"topic"}
		consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers": host,
			"group.id":          "teste",
			"auto.offset.reset": "earliest",
		})
		if err != nil {
			fmt.Printf("Failed to create consumer: %s\n", err)
			os.Exit(1)
		}

		err = consumer.SubscribeTopics(topics, nil)
		if err != nil {
			fmt.Printf("Subscribed topics error: %s\n", err)
			os.Exit(1)
		}

		run := true
		for run {
			ev := consumer.Poll(100)
			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("Consumed %v\n", string(e.Value))
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
				run = false
				wg.Done()
			}
		}

		consumer.Close()
	}()

	wg.Wait()
}
