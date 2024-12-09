package main

import (
	"log"

	config "github.com/adamlahbib/go-ms-poc/common"
	"github.com/adamlahbib/go-ms-poc/spec"
	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"
)

type RabbitMsg struct {
	QueueName string                     `json:"queueName"`
	Message   spec.CreateDocumentMessage `json:"message"`
}

// channel to publish rabbit messages
var pchan = make(chan RabbitMsg, 10)

func initProducer() {
	// conn
	conn, err := amqp.Dial(config.RabbitConfig.Uri)
	if err != nil {
		log.Fatalf("ERROR: Failed to connect to RabbitMQ: %s", err)
	}

	log.Printf("INFO: Connected to RabbitMQ")

	// create channel
	amqpChannel, err := conn.Channel()
	if err != nil {
		log.Fatalf("ERROR: Failed to create RabbitMQ channel: %s", err)
	}

	for {
		select {
		case msg := <-pchan:
			// marshal
			data, err := proto.Marshal(&msg.Message)
			if err != nil {
				log.Printf("ERROR: Failed to marshal message: %s", err)
				continue
			}

			// publish message
			err = amqpChannel.Publish(
				"",            // exchange
				msg.QueueName, // routing key
				false,         // mandatory
				false,         // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        data,
				},
			)

			if err != nil {
				log.Printf("ERROR: Failed to publish message: %s", err.Error())
				continue
			}

			log.Printf("INFO: Published message: %v", msg.Message)
		}
	}

}
