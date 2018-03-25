package db

import (
	"github.com/jinzhu/gorm"
	"time"
	"strings"
)

func CreateOrderBook (db gorm.DB, market BitfinexMarket) BitfinexOrderBook{
	// since ID is zero, GORM will override the value and auto increment it.
	orderbook := BitfinexOrderBook{0,market.ID, time.Now()}
	err := db.Create(&orderbook).Error
	if err != nil{
		panic(err)
	}
	return orderbook
}


func AddOrder(db gorm.DB, rate float64, quantity float64, orderbook BitfinexOrderBook) BitfinexOrder{
	order := BitfinexOrder{0, orderbook.ID, rate, quantity, time.Now()}
	err := db.Create(&order).Error
	if err != nil{
		panic(err)
	}
	return order
}

func CreateGetMarket(db gorm.DB, ticker string, quote string) BitfinexMarket{
	market := BitfinexMarket{0, ticker, quote}
	err := db.Create(&market).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			db.Where(map[string]interface{}{"ticker": ticker, "quote": quote}).Find(&market)
		}
	}
	return market
}

func AddTicker(db gorm.DB, market BitfinexMarket, price float64, volume float64){
	ticker := BitfinexTicker{0, market.ID, price, volume,time.Now()}
	err := db.Create(&ticker).Error
	if err != nil {
		panic(err)
	}
}




