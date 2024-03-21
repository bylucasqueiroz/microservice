#!/bin/bash

# Wait for LocalStack to be ready
sleep 5

# Create SQS queue
aws --endpoint-url=http://localhost:4566 sqs create-queue --queue-name myqueue

# Create SNS topic
aws --endpoint-url=http://localhost:4566 sns create-topic --name mytopic