package main

import (
	"fmt"
	"log"
	"math"
	"sync"

	"github.com/lvestera/slot-machine/internal/client"
)

func main() {

	log.Println("Run client")
	client := client.NewClient()

	var wins float64

	var wg sync.WaitGroup
	resultClient := make(chan float64, 10)

	play := func(spins int, player int, ch chan float64) {

		client.Play(spins, player, ch)

		wg.Done()
	}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go play(200000, i, resultClient)
	}
	for i := 0; i < 5; i++ {
		wins += <-resultClient
	}

	wg.Wait()

	spent := 1000000

	fmt.Printf("Total spent: %v \n", spent)
	fmt.Printf("Total spins: %v \n", spent)
	fmt.Printf("Total wins: %v \n", wins)

	fmt.Printf("RTP: %v \n", math.Round(wins/float64(spent)*100)/100)
	//fmt.Printf("RTP: %v \n", wins/float64(spent)*100)
	// fmt.Println(client.Reel)
}
