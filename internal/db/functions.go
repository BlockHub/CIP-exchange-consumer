package db

import (
	"github.com/jinzhu/gorm"
	"time"
)

func CreateOrderBook (db gorm.DB, market BitfinexMarket) BitfinexOrderBook{
	// since ID is zero, GORM will override the value and auto increment it.
	orderbook := BitfinexOrderBook{0,market.ID, int64(time.Now().Unix())}
	err := db.Create(&orderbook).Error
	if err != nil{
		panic(err)
	}
	return orderbook
}


func AddOrder(db gorm.DB, rate float64, quantity float64, orderbook BitfinexOrderBook) BitfinexOrder{
	order := BitfinexOrder{0, orderbook.ID, rate, quantity, int64(time.Now().Unix())}
	err := db.Create(&order).Error
	if err != nil{
		panic(err)
	}
	return order
}

func CreateOrGetMarket(db gorm.DB, ticker string, quote string) BitfinexMarket{
	market := BitfinexMarket{0, ticker, quote}
	err := db.Create(&market).Error
	if err != nil {
		err := db.Find(&market).Error
		if err != nil{
			panic(err)
		}
		return market
	}
	return market
}

func AddTicker(db gorm.DB, market BitfinexMarket, price float64, volume float64){
	ticker := BitfinexTicker{0, market.ID, price, volume,int64(time.Now().Unix())}
	err := db.Create(&ticker)
	if err != nil {
		panic(err)
	}
}




