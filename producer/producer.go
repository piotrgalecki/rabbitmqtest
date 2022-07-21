package main

import (
	"flag"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("Failed to %s: %s", msg, err)
	} else {
		log.Println("Succeeded to ", msg)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "open a channel")
	defer ch.Close()

	var queueName string
	flag.StringVar(&queueName, "queue", "none", "queue name")
	var msg string
	flag.StringVar(&msg, "message", "Hello!", "message")
	flag.Parse()
	log.Println("Queue=", queueName, " Msg=", msg)

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "declare a queue")
	log.Println("")

	body := msg
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "publish a message")
	log.Printf(" [x] Sent %s\n", body)
}
