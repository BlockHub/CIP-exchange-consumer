package db

import "time"

type BitfinexOrderBook struct {
	ID uint 			`gorm:"primary_key"`
	MarketID uint
	Time time.Time		`gorm:"primary_key"`
}

type BitfinexOrder struct {
	ID uint 			`gorm:"primary_key"`
	OrderbookID uint
	Rate float64
	//bitfinex supports giving the total number of sell/buyorders.
	//however we should skimp on memory and not add those
	//count int64
	Quantity float64
	Time time.Time		`gorm:"primary_key"`
	}

type BitfinexMarket struct {
	ID uint 			`gorm:"primary_key"`
	Ticker string		`gorm:"unique_index:idx_market"`
	Quote string		`gorm:"unique_index:idx_market"`
}

type BitfinexTicker struct {
	ID  uint 			`gorm:"primary_key"`
	MarketID uint
	Price float64
	Volume float64
	Time time.Time		`gorm:"primary_key"`
}