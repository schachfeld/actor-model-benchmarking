package main

import (
	"encoding/json"
	"strings"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type JsonInterpreter struct {
	act.Actor
	productRouter gen.PID
}

type ParseJsonMessage struct {
	json string
}

type CBMessage struct {
	Channel      string    `json:"channel"`
	Client_id    string    `json:"client_id"`
	Timestamp    string    `json:"timestamp"`
	Sequence_num int       `json:"sequence_num"`
	Events       []CBEvent `json:"events"`
}

type CBEvent struct {
	Event_type string `json:"type"`
	Product_id string `json:"product_id"`
	Updates    []CBUpdate
}

type CBUpdate struct {
	Side         string `json:"side"`
	Event_time   string `json:"event_time"`
	Price_level  string `json:"price_level"`
	New_quantity string `json:"new_quantity"`
}

func jsonInterpreterFactory() gen.ProcessBehavior {
	return &JsonInterpreter{}
}

func (a *JsonInterpreter) Init(args ...any) error {
	pid, err := a.Spawn(productRouterFactory, gen.ProcessOptions{})
	if err != nil {
		return err
	}
	a.productRouter = pid
	return nil
}

func (a *JsonInterpreter) HandleMessage(from gen.PID, message any) error {
	switch msg := message.(type) {
	case ParseJsonMessage:
		{
			var result CBMessage

			decoder := json.NewDecoder(strings.NewReader(msg.json))
			decoder.DisallowUnknownFields()

			err := decoder.Decode(&result)
			if err != nil {
				// a.Log().Warning("Error parsing JSON. Skipping message")
				return nil
			}

			a.Send(a.productRouter, RouteJsonMessage{cbMessage: result})
		}
	default:
		{
			panic("JsonInterpreter received unknown message")
		}
	}
	return nil
}

func (a *JsonInterpreter) Terminate(reason error) {
	a.Log().Info("%s terminated with reason: %s", a.PID(), reason)
}
