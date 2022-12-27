package main

// Authors: Lluis Barca & Alejandro Medina
import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

var (
	sushi []int
	p int = 0
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
	fmt.Printf(" [*] Waiting for messages. To exit press CTRL+C\n")
	 
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
	// fmt.Printf(" [*] Waiting for messages. To exit press CTRL+C\n")

	// Channel for the client to warn they finished
	warn := make(chan int)

	// Start consuming shushi
	go func() {
		// Channel to consume sushi
		cons, err := ch.Consume(
			queue.Name,	// queue
			"",      	// consumer
            false,    	// auto-ack
			false,     	// exclusive
            false,    	// no-local
			false,    	// no-wait
            nil,      	// arguments
		)

		// Channel to consume perms
		perms, err := ch.Consume(
			consume.Name,	// queue
			"",             // consumer
            false,          // auto-ack
			false,          // exclusive
            false,          // no-local
			false,          // no-wait
            nil,            // arguments
		)

		// Number of shushi pieces to eat
		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		sushi := rand.New(random).Intn(20)
		fmt.Printf("Client will consume %d shushi pieces\n", sushi)
		
		for i := 0; i < sushi; i++ {
			// Take a permission to eat
			for d := range perms {
				p, err = strconv.Atoi(string(d.Body))
				failOnError(err, "Failed to convert body to int")
				
				d.Ack(false)

				break
			}

			p = p - 1

			// The consumer will consume 
			for d := range cons {
				fmt.Printf("Client will eat a  %s\n", string(d.Body))
				time.Sleep(200)
				d.Ack(false)

				break
			}

			time.Sleep(600)
			fmt.Printf("Sushi pieces at the table %d\n", i + 1)

			// If there are sushi pieces to eat
			if p > 0 {
				perms := strconv.Itoa(p)
				err = ch.Publish(
					"", 			// exchange
					consume.Name,	// routing key
					false, 			// mandatory
					false,			// immediate

					amqp.Publishing{
						DeliveryMode	: amqp.Persistent,
						ContentType 	: "text/plain",
						Body 			: []byte(perms),
					})
				failOnError(err, "Failed to publish the message")
			}
		}
		warn <- 0
	}()
	<- warn
    fmt.Printf("Good night\n")
}