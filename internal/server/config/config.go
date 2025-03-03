package config

type Config struct {
	Coefficients map[string][2]float64
	DBConnection string
	Host         string
}

func NewConfig() *Config {
	//symbol: {frequency, cost}
	coeff := map[string][2]float64{
		"A": {0.4, 2},   //40% of wins
		"B": {0.25, 5},  //25% of wins
		"C": {0.15, 10}, //15% of wins
		"D": {0.1, 20},  //10% of wins
		"E": {0.05, 30}, //5% of wins
	}

	return &Config{
		Coefficients: coeff,
		Host:         "localhost:8081",
		DBConnection: "host=localhost user=gouser password=gouser dbname=gouser_db sslmode=disable",
	}
}
