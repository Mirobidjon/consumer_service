# consumer_service

## Information questionnaire. Listen v1.phone.consumer route and get data from the rest service and send response to the logger service.
Listens to messages in the format:
```json
{
    "record_id": "string uuid v4",
}
```

## How to work 
    1. As soon as it receives a message, it logs it to Logger service (debug route). 
    2. Get data from the rest service
        1. If not found - logs to Logger service (error route): Not found :record_id 
        2. If found - logs to Logger service (info route): Found :record_id
        3. If internal service error - the initial message is sent for reprocessing but with a delay of retryCount * 10 seconds. (retryCount is the number of attempts to process the message). So the message will be processed again after 10, 20, 30, 40, 50 seconds. If the message is not processed after 5 attempts, this message deleted from the queue.
    3. Send All messages to the logger service (consumer4 route)

## Installation
Use the following command to run the application:

    docker run us.gcr.io/learn-cloud-0809/consumer_service:latest

## ENVIRONMENT VARIABLES
    1. RABBIT_MQ_URL - RabbitMQ url (default: amqp://guest:guest@localhost:5672)
    2. EXCHANGE_NAME - Exchange name (default: v1.phone)
    3. REST_SERVICE_URL - Rest service url (default: http://localhost:80)

