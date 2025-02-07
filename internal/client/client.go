package client

import (
	"encoding/json"
	"fmt"
	"time"

	"log"

	"github.com/lvestera/slot-machine/internal/models"

	"github.com/go-resty/resty/v2"

	rand2 "crypto/rand"
	"math/rand/v2"
)

type Client struct {
	Coefficients []models.Coefficient
	Reel         []string
}

func NewClient() *Client {

	coefficients := getConfig()

	reel := fillReel(coefficients)

	return &Client{
		Coefficients: coefficients,
		Reel:         reel,
	}
}
func (client *Client) Play(spins int, player int, ch chan float64) {

	log.Println("Play spins ", spins)
	var spent uint64

	var win, wins float64

	b := make([]byte, 32)
	_, err := rand2.Read(b)
	if err != nil {
		fmt.Println("error:", err)

	}
	generator := rand.NewChaCha8([32]byte(b))

	len := len(client.Reel)

	buffer := make(map[int]models.Result)
	iBuf := 0
	for i := 0; i < spins; i++ {
		spent += 1
		win = 0

		reel1 := client.SpinReel(generator, len)
		reel2 := client.SpinReel(generator, len)
		reel3 := client.SpinReel(generator, len)

		for _, coeff := range client.Coefficients {

			if reel1 == reel2 && reel1 == reel3 && reel1 == coeff.Symbol {
				win = coeff.Cost
			}
		}

		wins += win

		result := models.Result{
			Player: player,
			Spin:   i,
			Result: fmt.Sprint(reel1, reel2, reel3),
			Win:    win,
		}

		if iBuf < 20 {
			buffer[iBuf] = result
			iBuf++
		} else {
			iBuf = 0
			client.Save(buffer)
			time.Sleep(50 * time.Millisecond)
		}
	}

	ch <- wins
}

func (client *Client) SpinReel(generator *rand.ChaCha8, len int) string {
	reelIndex := (generator.Uint64()) % uint64(len)

	return client.Reel[reelIndex]

}

func (client *Client) Save(results map[int]models.Result) {
	var err error
	var body []byte

	if body, err = json.Marshal(results); err != nil {
		log.Fatal(err.Error())
	}

	clientR := resty.New()

	_, err = clientR.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post("http://localhost:8081/send-result")

	if err != nil {
		log.Fatal(err.Error())
	}
}

func getConfig() []models.Coefficient {

	log.Println("Request config")
	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		Post("http://localhost:8081/get-config")

	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Got config ")
	var coefficients []models.Coefficient

	err = json.Unmarshal(resp.Body(), &coefficients)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Unmarshal config")
	return coefficients

}

func fillReel(coefficients []models.Coefficient) []string {

	reel := make([]string, 0, 1000)

	var i, max int
	for _, coeff := range coefficients {

		d := int(coeff.Distribution * 1000)
		log.Println("Dist ", coeff.Symbol, " ", coeff.Distribution, " ", i)
		max = i + d
		for i < max {
			reel = append(reel, coeff.Symbol)
			i++
		}
	}

	for i < 1000 {
		reel = append(reel, "0")
		i++
	}

	return reel
}
