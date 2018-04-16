// handlers used by consumers for dealing with incoming data

package handlers

import (
	"fmt"
	"CIP-exchange-consumer-bitfinex/internal/db"
	"github.com/jinzhu/gorm"
	"time"
)

type Handler interface {
	Handle([]float64)
}

// simple handler for printing the data
type PrintHandler struct{}
func (PrintHandler) Handle(data []float64){
	fmt.Println(data)
}

// saves each order as a single transaction in the DB (postgress takes a decent hit at approx 1000 writes/s)
type OrderDbHandler struct{
	Db *gorm.DB
	Orderbook db.BitfinexOrderBook
}
func (h OrderDbHandler) Handle(data []float64){
	if len(data) == 1 {
		return
	}
	db.AddOrder(*h.Db, data[0], data[2], h.Orderbook)
}

// saves a ticker to the DB
type TickerDbHandler struct{
	Db *gorm.DB
	Market db.BitfinexMarket
}
func (h TickerDbHandler) Handle(data []float64){
	fmt.Println(h.Market, "ticker", data)
	db.AddTicker(*h.Db, h.Market, data[1], data[7])
}

// saves a ticker to the DB
type TradeDbHandler struct{
	Db *gorm.DB
	Market db.BitfinexMarket
}
func (h TradeDbHandler) Handle(data []float64){
	if len(data) == 3{
		// sometimes frigging bitfinex does not give us an ID: we use zero which means
		// the DB should auto increment.
		data = append([]float64{float64(0)}, data...)
	}
	fmt.Println(data)
	db.AddTrade(*h.Db, h.Market, data[2], data[3], time.Unix(int64(data[1]), 0))
}

