package main

import (
	"math"
	"strconv"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type AvgOrderBookCalculator struct {
	act.Actor
	filewriter gen.PID
}

type UpdatesMessage struct {
	updates []CBUpdate
}

func avgOrderBookCalculatorFactory() gen.ProcessBehavior {
	return &AvgOrderBookCalculator{}
}

func (a *AvgOrderBookCalculator) Init(args ...any) error {

	productName := args[0].(string)

	filewriter, err := a.Spawn(fileWriterFactory, gen.ProcessOptions{}, "orderbookdata/"+productName+".txt")
	if err != nil {
		return err
	}
	a.filewriter = filewriter

	a.Log().Info("fileWriter for avgOrderBookCalculator started")

	return nil
}

func (a *AvgOrderBookCalculator) HandleMessage(from gen.PID, message any) error {
	// a.Log().Info("Received message from %s", from)
	switch msg := message.(type) {
	case UpdatesMessage:
		{
			var bids, offers []float64

			for _, update := range msg.updates {
				price, err := strconv.ParseFloat(update.Price_level, 64)
				if err != nil {
					continue
				}

				if update.Side == "bid" {
					bids = append(bids, price)
				} else {
					offers = append(offers, price)
				}
			}

			avgBids := avg(bids)
			avgOffers := avg(offers)

			if !math.IsNaN(avgBids) {
				a.Send(a.filewriter, WriteLineMessage{content: "Bids: " + strconv.FormatFloat(avg(bids), 'f', -1, 64)})
			}
			if !math.IsNaN(avgOffers) {
				a.Send(a.filewriter, WriteLineMessage{content: "Offers: " + strconv.FormatFloat(avg(offers), 'f', -1, 64)})
			}
		}
	case LastMessage:
		{
			a.Send(a.filewriter, LastMessage{})
		}
	case DoneMessage:
		{
			a.Send(a.Parent(), DoneMessage{})
		}
	default:
		{
			panic("AvgOrderBookCalculator received unknown message")
		}
	}
	return nil
}

func avg(elements []float64) float64 {
	var sum float64
	for _, bid := range elements {
		sum += bid
	}
	return sum / float64(len(elements))
}

func (a *AvgOrderBookCalculator) Terminate(reason error) {
	a.Log().Info("%s terminated with reason: %s", a.PID(), reason)
}
