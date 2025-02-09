package main

import (
	"fmt"
	"log"
	"math"
	"sync"

	"github.com/lvestera/slot-machine/internal/client"
)

func main() {
	log.Println("Run clients")

	host := "http://localhost:8081"

	var wg sync.WaitGroup
	resultCh := make(chan float64)
	var wins float64

	for i := 0; i < 5; i++ {

		wg.Add(1)
		go func() {
			defer wg.Done()
			client, err := client.NewClient(host, i)
			if err != nil {
				log.Println(err)
				return
			}

			_, wins, err := client.Play(200000)
			resultCh <- wins
			if err != nil {
				log.Println(err)
				return
			}
		}()
	}

	go func() {
		for v := range resultCh {
			wins += v
		}
	}()
	wg.Wait()
	close(resultCh)
	spent := 1000000
	fmt.Printf("Total spent: %v \n", spent)
	fmt.Printf("Total spins: %v \n", spent)
	fmt.Printf("Total wins: %v \n", wins)

	fmt.Printf("RTP: %v \n", math.Round(wins/float64(spent)*100)/100)
}
