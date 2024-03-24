#!/bin/bash

# Wait for LocalStack to be ready
sleep 5

# Create SQS queue
awslocal sqs create-queue --queue-name myqueue

# Create SNS topic
awslocal sns create-topic --name mytopic