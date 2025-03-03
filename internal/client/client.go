package client

import (
	"log"
	"strings"
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
	reel := make([]string, 0, 10000)

	var i, max int
	for _, coeff := range client.Coefficients {

		d := int(coeff.Distribution)
		max = i + d
		for i < max {
			reel = append(reel, coeff.Symbol)
			i++
		}
	}

	for i < 10000 {
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

	lettersCount := len(client.Coefficients)
	letters := make([]string, lettersCount)

	for i, coeff := range client.Coefficients {
		letters[i] = coeff.Symbol
	}

	var reelResult string
	for i := 0; i < spins; i++ {
		spent += 1
		win = 0

		reelValue := client.SpinReel()

		if reelValue == "0" {
			reelResult = generateRandomResult(letters, lettersCount)
		} else {
			for _, coeff := range client.Coefficients {
				if reelValue == coeff.Symbol {
					win = coeff.Cost
					reelResult = strings.Repeat(coeff.Symbol, lettersCount)
					break
				}
			}
		}

		wins += win

		result := models.Result{
			Player: client.Name,
			Spin:   i,
			Result: reelResult,
			Win:    win,
		}

		if iBuf == 50 {
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

func generateRandomResult(letters []string, lettersCount int) string {
	var letterIndex int
	var nextLetter string

	result := make([]string, 0, lettersCount)
	for i := 0; i < lettersCount; i++ {
		letterIndex = rand.IntN(lettersCount)
		nextLetter = letters[letterIndex]
		result = append(result, nextLetter)
	}

	if strings.Count(strings.Join(result, ""), nextLetter) != lettersCount {
		return strings.Join(result, "")
	}

	if letterIndex == lettersCount-1 {
		result[0] = letters[0]
	} else {
		result[0] = letters[letterIndex+1]
	}

	return strings.Join(result, "")
}

func (client *Client) SpinReel() string {
	len := len(client.Reel)
	reelIndex := (client.generator.Uint64()) % uint64(len)

	return client.Reel[reelIndex]
}
