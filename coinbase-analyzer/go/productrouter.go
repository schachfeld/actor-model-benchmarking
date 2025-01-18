package main

import (
	"fmt"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type ProductRouter struct {
	act.Actor
	productDistributors map[string]gen.PID
	doneDistributors    map[string]bool
}

func productRouterFactory() gen.ProcessBehavior {
	return &ProductRouter{}
}

func (a *ProductRouter) Init(args ...any) error {
	a.productDistributors = make(map[string]gen.PID)
	a.doneDistributors = make(map[string]bool)
	return nil
}

func (a *ProductRouter) HandleMessage(from gen.PID, message any) error {
	// a.Log().Info("Received message from %s", from)
	switch msg := message.(type) {
	case RouteJsonMessage:
		{
			for _, event := range msg.cbMessage.Events {
				// if product does not exist, create an actor tree for that product
				if _, ok := a.productDistributors[event.Product_id]; !ok {
					pid, err := a.Spawn(productDistributorFactory, gen.ProcessOptions{}, event.Product_id)
					if err != nil {
						return err
					}
					a.productDistributors[event.Product_id] = pid
				}

				a.Send(a.productDistributors[event.Product_id], event)
			}
		}
	case LastMessage:
		{
			for _, pid := range a.productDistributors {
				a.Send(pid, LastMessage{})
			}
		}
	case DoneMessage:
		{
			if _, ok := a.doneDistributors[from.String()]; ok {
				return fmt.Errorf("parser %s is already done", from.String())
			}
			a.doneDistributors[from.String()] = true
			if len(a.doneDistributors) == len(a.productDistributors) {
				a.Send(a.Parent(), DoneMessage{})
			}
		}

	}
	return nil
}

func (a *ProductRouter) Terminate(reason error) {
	a.Log().Info("%s terminated with reason: %s", a.PID(), reason)
}
