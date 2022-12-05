package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	CLIENTS = 4
	SUSHIS = 10
)

var (
	names = [CLIENTS] string {"Cayetano", "Juan", "Victoria", "Francisca"}
	sushi = [SUSHIS] string {"", "", "", "", "", "", "", "", "", "", }
	sushi_plate = [] string {}
)

type Empty struct {}

type sushi_piece struct {
	t string
	n int
}

func client(id int, client string) {

}

func server(done chan Empty) {

}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	// Channel to finish
	done := make(chan Empty)

	// Filling the sushi plate

	// Init sushi's channel

	// Start go rutines
	go server(done)
	
	for i := 0; i < CLIENTS; i++ {
		go client(i, names[i])
	}

	fmt.Printf("END OF SIMULATION")
}