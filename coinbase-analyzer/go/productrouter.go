package main

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type ProductRouter struct {
	act.Actor
	productParsers map[string]gen.PID
}

func productRouterFactory() gen.ProcessBehavior {
	return &ProductRouter{}
}

func (a *ProductRouter) Init(args ...any) error {
	a.productParsers = make(map[string]gen.PID)
	return nil
}

func (a *ProductRouter) HandleMessage(from gen.PID, message any) error {
	// a.Log().Info("Received message from %s", from)
	switch msg := message.(type) {
	case RouteJsonMessage:
		{
			for _, event := range msg.cbMessage.Events {
				// if product does not exist, create an actor tree for that product
				if _, ok := a.productParsers[event.Product_id]; !ok {
					pid, err := a.Spawn(productDistributorFactory, gen.ProcessOptions{}, event.Product_id)
					if err != nil {
						return err
					}
					a.productParsers[event.Product_id] = pid
				}

				a.Send(a.productParsers[event.Product_id], event)
			}
		}
	}
	return nil
}

func (a *ProductRouter) Terminate(reason error) {
	a.Log().Info("%s terminated with reason: %s", a.PID(), reason)
}
