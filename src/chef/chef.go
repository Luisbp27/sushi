package main

// Authors: Lluis Barca & Alejandro Medina
// Link Video: https://drive.google.com/file/d/1XPPSNVxzOBjA7LfmyorhkpbOEHHpwzxO/view?usp=sharing
import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	amqp "github.com/streadway/amqp"
)

const (
	CLIENTS = 4
	SUSHIS = 10
)

var (
	sushi_types = []string{"Nigiri de salmó", "Sashimi de tonyina", "Maki de cranc"}
	p int
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

	fmt.Printf("Connecting to RabbitMQ...\n")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare the sushi plate queue
	queue, err := ch.QueueDeclare(
        "sushi", // name
        true,   // durable
        false,   // delete when unused
        false,   // exclusive
        false,   // no-wait
		nil,     // arguments
	)
    failOnError(err, "Failed to declare the sushi queue")

	// Declare the consumer permisions to eat
	consume, err := ch.QueueDeclare(
        "consume", 	// name
        true,		// durable
        false,		// delete when unused
		false,		// exclusive
		false,      // no-wait
        nil,        // arguments
    )
	failOnError(err, "Failed to declare the consume queue")

	// Fill the queue with sushis
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < SUSHIS; i++ {
		sushi := random.Intn(len(sushi_types))
		err = ch.Publish(
			"",				// exchange
			queue.Name,		// routing key
			false,			// mandatory
			false,			// immediate

			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
                ContentType:  "text/plain",
                Body:         []byte(sushi_types[sushi]),
			},
		)
		failOnError(err, "Failed to publish a message to the queue")
		fmt.Printf(" [x] %s it's cooked", sushi_types[sushi])
		time.Sleep(400)
	}

	// Fill the permisions queue
	p := strconv.Itoa(SUSHIS)
	err = ch.Publish(
		"",					// exchange
        consume.Name,      	// routing key
        false,            	// mandatory
        false,            	// immediate

        amqp.Publishing{
            DeliveryMode: amqp.Persistent,
            ContentType:  "text/plain",
			Body:         []byte(p),
		},
	)
	failOnError(err, "Failed to publish a message to the consume queue")
	fmt.Printf(" [x] %s", p)
}