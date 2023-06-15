package handlers

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

func GetRabbitResponse(name, to string, req []byte) ([]byte, error) {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		name,  // name
		false, // durable
		true,  // delete when unused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	//reqByte, err := proto.Marshal(&brand.ShowBrandRequest{
	//	Id: id,
	//})
	//failOnError(err, "Failed to marshal")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",    // exchange
		to,    // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			ReplyTo:     q.Name,
			Body:        req,
		})
	failOnError(err, "Failed to publish a message")

	for d := range msgs {
		return d.Body, nil
	}
	return nil, nil
}
