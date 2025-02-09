package config

type Config struct {
	Rtp          float64
	WinAAAFreq   float64
	WinBBBFreq   float64
	WinCCCFreq   float64
	WinAAACost   float64
	WinBBBCost   float64
	WinCCCCost   float64
	DBConnection string
	Host         string
}

const rtp = 0.95
const winAAAFreq = 0.5  //50% of wins
const winBBBFreq = 0.3  //30% of wins
const winCCCFreq = 0.15 //15% of wins
const winAAACost = 5
const winBBBCost = 10
const winCCCCost = 20

func NewConfig() *Config {
	return &Config{
		Rtp:          rtp,
		WinAAAFreq:   winAAAFreq,
		WinBBBFreq:   winBBBFreq,
		WinCCCFreq:   winCCCFreq,
		WinAAACost:   winAAACost,
		WinBBBCost:   winBBBCost,
		WinCCCCost:   winCCCCost,
		Host:         "localhost:8081",
		DBConnection: "host=localhost user=gouser password=gouser dbname=gouser_db sslmode=disable",
	}
}
