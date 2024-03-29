module consumer

go 1.22.0

require (
    github.com/confluentinc/confluent-kafka-go v1.9.2
    github.com/bylucasqueiroz/microservice/libs/kafka v0.0.0
)

replace github.com/bylucasqueiroz/microservice/libs/kafka v0.0.0 => ../../libs/kafka
