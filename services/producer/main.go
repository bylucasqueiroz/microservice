package main

import (
	"fmt"
	"os"

	kfk "github.com/bylucasqueiroz/microservice/libs/kafka"
)

func main() {
	fmt.Println("Start process..")
	host := os.Getenv("KAFKA_ENDPOINT")

	producer := kfk.NewKafkaProducer(&kfk.KafkaProducerConfig{
		Host:     host,
		ClientId: "group1",
		Acks:     "all",
	})

	count := 10
	for i := 0; i < count; i++ {
		producer.Produce("teste", []byte(fmt.Sprintf("TESTE %v", i)))
	}
}
