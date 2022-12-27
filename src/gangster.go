package main

// Authors: Lluis Barca & Alejandro Medina
import (
	"fmt"
	"strconv"

	amqp "github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err!= nil {
        panic(fmt.Errorf("%s: %s", msg, err))
    }
}

func main() {
	// Connect to RabbitMQ
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare the sushi plate queue
	queue, err := ch.QueueDeclare(
        "sushi",	// name
        true,   	// durable
        false,   	// delete when unused
        false,   	// exclusive
        false,   	// no-wait
		nil,     	// arguments
	)

	// Instructions to confirm the message are consumed
	err = ch.Qos(
		1,     	// prefetch count
        0,     	// prefetch size
        false,	// global
	)
	
	failOnError(err, "Failed to declare the consume queue")
	fmt.Printf("I wanna all the sushi!\n")

	// Declare the consumer permisions to eat
	consume, err := ch.QueueDeclare(
        "consume",	// name
        true,		// durable
        false,		// delete when unused
		false,		// exclusive
		false,      // no-wait
        nil,        // arguments
    )

	// Instructions to confirm the message are consumed
	err = ch.Qos(
		1,     	// prefetch count
        0,     	// prefetch size
        false,	// global
	)
	// failOnError(err, "Failed to declare the consume queue")
	// fmt.Printf("I wanna all the sushi!")

	// Channel for the client to warn they finished
	warn := make(chan int)

	go func() {
		// Channel to consume sushi
		cons, err := ch.Consume(
			queue.Name,	// queue
			"",      	// consumer
            true,    	// auto-ack
			false,     	// exclusive
            false,    	// no-local
			false,    	// no-wait
            nil,      	// arguments
		)

		// Channel to consume perms
		perms, err := ch.Consume(
			consume.Name,	// queue
			"",             // consumer
            true,          	// auto-ack
			false,          // exclusive
            false,          // no-local
			false,          // no-wait
            nil,            // arguments
		)

		p := 0

		// Take a permission to eat
		for d := range perms {
			p, err = strconv.Atoi(string(d.Body))
			failOnError(err, "Failed to parse the permissions")
			fmt.Printf("Permissions: %d\n", p)

			break		
		}
		fmt.Printf("Permissions: %d", p)
		// The consumer will consume 
		for d := range cons {
			fmt.Printf("Sushi piece number: %s\n", string(d.Body))

			break
		}

		warn <- 0
	}()

	<- warn
	fmt.Printf("Good night!")
}