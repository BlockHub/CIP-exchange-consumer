package main

import (
	 "github.com/getsentry/raven-go"
	"github.com/bitfinexcom/bitfinex-api-go/v1"
	"log"
	"strings"
	"CIP-exchange-consumer-bitfinex/pkg/handlers"
	"CIP-exchange-consumer-bitfinex/pkg/consumer"
	"github.com/jinzhu/gorm"
	 _ "github.com/jinzhu/gorm/dialects/postgres"

	"os"
	"CIP-exchange-consumer-bitfinex/internal/db"
	"github.com/joho/godotenv"
)

func init(){
	useDotenv := true
	if os.Getenv("PRODUCTION") == "true"{
		useDotenv = false
	}

	// this loads all the constants stored in the .env file (not suitable for production)
	// set variables in supervisor then.
	if useDotenv {
		err := godotenv.Load()
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
	}
	raven.SetDSN(os.Getenv("RAVEN_DSN"))
}


func main() {
	c := bitfinex.NewClient()

	pairs, err := c.Pairs.All()
	if nil != err {
		raven.CaptureErrorAndWait(err, nil)
	}

	err = c.WebSocket.Connect()
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
	}
	defer c.WebSocket.Close()

	gormdb, err := gorm.Open(os.Getenv("DB"), os.Getenv("DB_URL"))
	if err != nil{
		raven.CaptureErrorAndWait(err, nil)
	}
	defer gormdb.Close()

	// migrations are only performed by GORM if a table/column/index does not exist.
	err = gormdb.AutoMigrate(&db.BitfinexMarket{}, &db.BitfinexTicker{}, &db.BitfinexOrder{}, &db.BitfinexOrderBook{}).Error
	if err != nil{
		raven.CaptureErrorAndWait(err, nil)
	}
	err = gormdb.Exec("CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;").Error
	if err != nil{
		raven.CaptureErrorAndWait(err, nil)
	}
	err = gormdb.Exec("SELECT create_hypertable('bitfinex_orders', 'time',  'orderbook_id', if_not_exists => TRUE)").Error
	if err != nil{
		raven.CaptureErrorAndWait(err, nil)
	}
	err = gormdb.Exec("SELECT create_hypertable('bitfinex_tickers', 'time', 'market_id', if_not_exists => TRUE)").Error
	if err != nil{
		raven.CaptureErrorAndWait(err, nil)
	}
	err =gormdb.Exec("SELECT create_hypertable('bitfinex_order_books', 'time', 'market_id', if_not_exists => TRUE)").Error
	if err != nil{
		raven.CaptureErrorAndWait(err, nil)
	}
	for _, pair := range pairs {
		// if the market already exists, this fails (with a warning, but no error, and the market is returned
		market := db.CreateGetMarket(*gormdb, pair[0:3], pair[len(pair)-3:])
		//a new orderbook is created at each disconnect/startup. Orderbooks are continuous chained orders
		orderbook := db.CreateOrderBook(*gormdb, market)

		bookChannel := make(chan []float64)
		trades_chan := make(chan []float64)

		c.WebSocket.AddSubscribe(bitfinex.ChanBook, strings.ToUpper(pair), bookChannel)
		c.WebSocket.AddSubscribe(bitfinex.ChanTrade, strings.ToUpper(pair), trades_chan)

		orderhandler := handlers.OrderDbHandler{gormdb, orderbook}
		tickerhandler := handlers.TickerDbHandler{gormdb, market}
		//tickerhandler := handlers.PrintHandler{}

		go consumer.Consumer(bookChannel, orderhandler)
		go consumer.Consumer(trades_chan, tickerhandler)
	}

	err = c.WebSocket.Subscribe()
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
	}
}


