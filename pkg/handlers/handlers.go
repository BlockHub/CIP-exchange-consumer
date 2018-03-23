// handlers used by consumers for dealing with incoming data

package handlers

import (
	"fmt"
	"CIP-exchange-consumer/internal/db"
	"github.com/jinzhu/gorm"
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
	db.AddTicker(*h.Db, h.Market, data[2], data[0])
}

