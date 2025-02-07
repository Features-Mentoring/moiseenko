package models

type Coefficient struct {
	Symbol       string
	Distribution float64
	Cost         float64
}

type Result struct {
	Player int
	Spin   int
	Result string
	Win    float64
}
