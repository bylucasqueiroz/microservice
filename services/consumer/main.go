package main

import (
	"fmt"
	"os"

	kfk "github.com/bylucasqueiroz/microservice/libs/kafka"
)

func main() {
	fmt.Println("Start process..")

	host := os.Getenv("KAFKA_ENDPOINT")

	consume := kfk.NewKafkaConsumer(&kfk.KafkaConsumerConfig{
		Host:            host,
		GroupId:         "group1",
		AutoOffsetReset: "earliest",
		Topics:          []string{"teste"},
	})

	go consume.Consume()

	consume.WaitConsumerDone()
}
