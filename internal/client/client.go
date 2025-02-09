package client

import (
	"fmt"
	"log"
	"time"

	"github.com/lvestera/slot-machine/internal/client/requests"
	"github.com/lvestera/slot-machine/internal/models"

	rand2 "crypto/rand"
	"math/rand/v2"
)

type Client struct {
	Name          int
	RequestClient *requests.RequestClient
	Coefficients  []models.Coefficient
	Reel          []string
	generator     *rand.ChaCha8
}

func NewClient(host string, name int) (*Client, error) {

	rc := requests.GetRequestClient(host)

	client := &Client{
		Name:          name,
		RequestClient: rc,
	}

	coefficients, err := rc.GetConfig()
	if err != nil {
		return nil, err
	}
	client.Coefficients = coefficients

	client.fillReel()

	b := make([]byte, 32)
	_, err = rand2.Read(b)
	if err != nil {
		return nil, err

	}
	client.generator = rand.NewChaCha8([32]byte(b))

	return client, nil
}

func (client *Client) fillReel() {
	reel := make([]string, 0, 1000)

	var i, max int
	for _, coeff := range client.Coefficients {

		d := int(coeff.Distribution * 1000)
		//log.Println("Dist ", coeff.Symbol, " ", coeff.Distribution, " ", i)
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

	client.Reel = reel
}

func (client *Client) Play(spins int) (uint64, float64, error) {

	log.Println("Play spins ", spins)

	var err error
	var spent uint64
	var win, wins float64

	buffer := make(map[int]models.Result)
	iBuf := 0
	for i := 0; i < spins; i++ {
		spent += 1
		win = 0

		reel1 := client.SpinReel()
		reel2 := client.SpinReel()
		reel3 := client.SpinReel()

		for _, coeff := range client.Coefficients {
			if reel1 == reel2 && reel1 == reel3 && reel1 == coeff.Symbol {
				win = coeff.Cost
			}
		}

		wins += win

		result := models.Result{
			Player: client.Name,
			Spin:   i,
			Result: fmt.Sprint(reel1, reel2, reel3),
			Win:    win,
		}

		if iBuf == 30 {
			err = client.RequestClient.SaveResults(buffer)
			if err != nil {
				return spent, wins, err
			}
			iBuf = 0
			for key := range buffer {
				delete(buffer, key)
			}

			time.Sleep(20 * time.Millisecond)
		}

		buffer[iBuf] = result
		iBuf++
	}
	if iBuf != 1 {
		err = client.RequestClient.SaveResults(buffer)
		if err != nil {
			return spent, wins, err
		}
	}

	return spent, wins, nil
}

func (client *Client) SpinReel() string {
	len := len(client.Reel)
	reelIndex := (client.generator.Uint64()) % uint64(len)

	return client.Reel[reelIndex]
}
