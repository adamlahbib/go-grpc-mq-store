package main

import (
	"log"

	config "github.com/adamlahbib/go-ms-poc/common"
	"github.com/adamlahbib/go-ms-poc/spec"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
)

func initConsumer() {
	// conn
	conn, err := amqp.Dial(config.RabbitConfig.Uri)
	if err != nil {
		log.Fatalf("ERROR: Failed to connect to RabbitMQ: %s", err.Error())
	}

	log.Printf("INFO: Connected to RabbitMQ")

	// create channel
	amqpChannel, err := conn.Channel()
	if err != nil {
		log.Fatalf("ERROR: Failed to create RabbitMQ channel: %s", err.Error())
	}

	// create queue
	queue, err := amqpChannel.QueueDeclare(
		"storage", // channel name
		true,      // durable
		false,     // autoDelete
		false,     // exclusive
		false,     // noWait
		nil,       // args
	)

	if err != nil {
		log.Fatalf("ERROR: Failed to declare queue: %s", err.Error())
	}

	// channel
	msgChannel, err := amqpChannel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // autoAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // args
	)
	if err != nil {
		log.Fatalf("ERROR: Failed to consume messages: %s", err.Error())
	}

	// consume
	for {
		select {
		case msg := <-msgChannel:
			// unmarshal
			docMsg := &spec.CreateDocumentMessage{}
			err := proto.Unmarshal(msg.Body, docMsg)
			if err != nil {
				log.Printf("ERROR: Failed to unmarshal message: %s", err.Error())
				continue
			}
			log.Printf("INFO: Received message: %s", docMsg.String())

			// ack for message
			err = msg.Ack(true)
			if err != nil {
				log.Printf("ERROR: Failed to ack message: %s", err.Error())
				continue
			}

			// handle docMsg
			handleMsg(docMsg)
		}
	}
}

func handleMsg(docMsg *spec.CreateDocumentMessage) {
	// reply
	docReply := spec.CreateDocumentReply{
		Uid:    docMsg.Uid,
		Status: "Created",
	}

	msg := RabbitMsg{
		QueueName: docMsg.ReplyTo,
		Reply:     docReply,
	}

	rchan <- msg
}
