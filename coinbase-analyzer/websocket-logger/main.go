package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "advanced-trade-ws.coinbase.com", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: *addr, Path: "/"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	f, err := os.OpenFile("messages.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("open file:", err)
		return
	}
	defer f.Close()

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			// log.Printf("recv: %s", message)

			if _, err := f.WriteString(string(message) + "\n"); err != nil {
				log.Println("write file:", err)
				return
			}
		}
	}()

	c.WriteMessage(websocket.TextMessage, []byte(`{"type": "subscribe", "channel": "level2", "product_ids": ["BTC-USD", 
	"ETH-USD", "ETH-BTC", "LTC-USD", "LTC-BTC", "BCH-USD", "BCH-BTC", "ETC-USD", "ETC-BTC", "ZRX-USD", "ZRX-BTC",
	"BAT-USD", "BAT-BTC", "LINK-USD", "LINK-BTC", "DAI-USD", "DAI-BTC", "REP-USD", "REP-BTC", "OMG-USD", "OMG-BTC"
	]}`))

	ticker := time.NewTicker(2 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			log.Println("2 hours passed, stopping the connection")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
