// Simple routines handling a Websocket channel


package consumer

import "CIP-exchange-consumer/pkg/handlers"

func Consumer(in chan []float64, handler handlers.Handler) {
	for {
		msg := <-in
		handler.Handle(msg)
	}
}

