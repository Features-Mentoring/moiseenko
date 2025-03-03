package main

import (
	"fmt"
	"log"
	"math"
	"sync"

	"github.com/lvestera/slot-machine/internal/client"
)

const host = "http://localhost:8081"
const perPlayer = 200_000

func main() {
	log.Println("Run clients")

	var wg sync.WaitGroup
	resultCh := make(chan float64)
	var wins float64
	var spins int

	for i := 0; i < 5; i++ {

		wg.Add(1)
		go func() {
			defer wg.Done()
			client, err := client.NewClient(host, i)
			if err != nil {
				log.Println(err)
				return
			}

			_, wins, err := client.Play(perPlayer)
			resultCh <- wins
			if err != nil {
				log.Println(err)
				return
			}
			spins += perPlayer
		}()
	}

	go func() {
		for v := range resultCh {
			wins += v
		}
	}()
	wg.Wait()
	close(resultCh)

	fmt.Printf("Total spent: %v \n", spins)
	fmt.Printf("Total spins: %v \n", spins)
	fmt.Printf("Total wins: %v \n", wins)

	fmt.Printf("RTP: %v \n", math.Round(wins/float64(spins)*100)/100)
}
