package main

import (
	"log"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(client int, err error, msg string) {
	if err != nil {
		log.Panicf("Client ", client, "Failed to %s: %s", msg, err)
	} else {
		log.Println("Client ", client, "Succeeded to ", msg)
	}
}

func consumer(client int, conn *amqp.Connection) {
	// create virtual channel within the tcp connection
	ch, err := conn.Channel()
	failOnError(client, err, "open a channel")
	defer ch.Close()

	chname := "client" + strconv.Itoa(client)

	q, err := ch.QueueDeclare(
		chname, // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	failOnError(client, err, "declare a queue")

	for {
		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		failOnError(client, err, "register a consumer")

		for d := range msgs {
			log.Printf("Client %d received a message: %s", client, d.Body)
		}

		time.Sleep(1 * time.Second)
	}
}

func main() {
	// create tcp connection
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(0, err, "connect to RabbitMQ")
	defer conn.Close()

	for i := 1; i <= 1000; i++ {
		go consumer(i, conn)
	}

	var forever chan struct{}
	<-forever
}
