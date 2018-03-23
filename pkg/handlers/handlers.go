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


type PrintHandler struct{}
func (PrintHandler) Handle(data []float64){
	fmt.Println(data)
}

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

type TickerHandler struct{
	Db *gorm.DB
	Market db.BitfinexMarket
}
func (h TickerHandler) Handle(data []float64){
	db.AddTicker(*h.Db, h.Market, data[2], data[0])
}

