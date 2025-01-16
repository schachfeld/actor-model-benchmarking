package main

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type ProductDistributor struct {
	act.Actor
	product_id             string
	avgOrderBookCalculator gen.PID
}

func productDistributorFactory() gen.ProcessBehavior {
	return &ProductDistributor{}
}

func (a *ProductDistributor) Init(args ...any) error {
	a.product_id = args[0].(string)

	pid, err := a.Spawn(avgOrderBookCalculatorFactory, gen.ProcessOptions{}, a.product_id)
	if err != nil {
		return err
	}
	a.avgOrderBookCalculator = pid
	a.Log().Info("avgOrderBookCalculator started")
	return nil
}

func (a *ProductDistributor) HandleMessage(from gen.PID, message any) error {
	// a.Log().Info("Received message from %s", from)
	switch msg := message.(type) {
	case CBEvent:
		{
			a.Send(a.avgOrderBookCalculator, UpdatesMessage{updates: msg.Updates})
		}
	default:
		{
			panic("ProductDistributor received unknown message")
		}
	}
	return nil
}

func (a *ProductDistributor) Terminate(reason error) {
	a.Log().Info("%s terminated with reason: %s", a.PID(), reason)
}
