package models

type StockC struct {
	Code   string       `json:"code"`
	Prices []StockPrice `json:"prices"`
}

type StockPrice struct {
	Date   string  `json:"date"`
	Value  float64 `json:"value"`
	Volume float64 `json:"volume"`
}
