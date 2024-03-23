package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	var wg sync.WaitGroup
	fmt.Println("Start process..")

	wg.Add(1)
	go func() {
		fmt.Println("Start consume..")
		topics := []string{"topic"}
		consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers": "localhost:9092",
			"group.id":          "foo",
			"auto.offset.reset": "earliest",
		})
		if err != nil {
			fmt.Printf("Failed to create producer: %s\n", err)
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
				fmt.Println(e.Value)
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
				run = false
				wg.Done()
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}

		consumer.Close()
	}()

	wg.Add(1)
	go func() {
		fmt.Println("Start produce..")
		for {
			p, err := kafka.NewProducer(&kafka.ConfigMap{
				"bootstrap.servers": "localhost:9092",
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
}
