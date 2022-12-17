package main

import (
	"fmt"
	"math/rand"
	"time"
	"github.com/streadway/amqp"
)

var (
	sushi []int
	p := 0
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
    defer failOnError(conn.Close(), "Failed to close the connection")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer failOnError(ch.Close(), "Failed to close the channel")

	// Declare the sushi plate queue
	queue, err := ch.QueueDeclare(
        "sushi", // name
        true,   // durable
        false,   // delete when unused
        false,   // exclusive
        false,   // no-wait
		nil,     // arguments
	)

	// Instructions to confirm the message are consumed
	err = ch.Qos(
		1,     // prefetch count
        0,     // prefetch size
        false, // global
	)

	// Declare the consumer permisions to eat
	consume, err := ch.QueueDeclarePassive(
        "consume", 	// name
        true,		// durable
        false,		// delete when unused
		false,		// exclusive
		false,      // no-wait
        nil,        // arguments
    )

	// Instructions to confirm the message are consumed
	err = ch.Qos(
		1,     // prefetch count
        0,     // prefetch size
        false, // global
	)
	failOnError(err, "Failed to declare the consume queue")
	fmt.Printf(" [*] Waiting for messages. To exit press CTRL+C\n")

	// Channel for the client to warn they finished
	warn := make(chan int)

	// Start consuming shushi
	go func() {
		// Channel to consume sushi
		cons, err := ch.Consume(
			queue.name,		// queue
			"",            	// consumer
            false,        	// auto-ack
			false,        	// exclusive
            false,        	// no-local
			false,        	// no-wait
            nil,          	// arguments
		)

		// Channel to consume perms
		perms, err := ch.Consume(
			consume.name, 	// queue
			"",             // consumer
            false,          // auto-ack
			false,          // exclusive
            false,          // no-local
			false,          // no-wait
            nil,            // arguments
		)

		// Number of shushi pieces to eat
		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		sushi := rand.Intn(20)
		fmt.Printf("Client will consume %d shushi pieces\n", sushi)

		for i := 0; i < sushi; i++ {
			if d = range perms {
				perms, err = strconv.Atoi(string(d.Body))
				failOnError(err, "Failed to convert body to int")
				d.ack(false)
			}

			p := p - 1

			if d = range cons {
				fmt.Printf("Client will eat a  %s\n", string(d.Body))
				time.Sleep(200)
				d.ack(false)
			}

			time.Sleep(600)
			fmt.Printf("Sushi pieces at the table %d\n", i + 1)

			// If there are sushi pieces to eat
			if p > 0 {
				perms := strconv.Itoa(p)
				err = ch.Publish{
					"", // exchange
					consume.name // routing key
					false, // mandatory
					false,	// immediate
					amqp.Publishing{
						perms := strconv.Itoa(p)
						err = ch.Publish{
							"", // exchange
							consume.name, // routing key
                            false, // mandatory
							false, // immediate
                            amqp.Publishing{
								DeliveryMode : amqp.Persistent,
								ContentType : "text/plain",
								Body : []byte(p),
							}
						}
						failOnError(err, "Failed to publish the message")
					}
				}
			}
		}
		warn <- 0
	}()
	<- warn
    fmt.Printf("Good night\n")
}