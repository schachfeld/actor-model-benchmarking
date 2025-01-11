package main

import (
	"encoding/json"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type JsonInterpreter struct {
	act.Actor
}

type ParseJsonMessage struct {
	json string
}

func jsonInterpreterFactory() gen.ProcessBehavior {
	return &JsonInterpreter{}
}

func (a *JsonInterpreter) Init(args ...any) error {
	return nil
}

func (a *JsonInterpreter) HandleMessage(from gen.PID, message any) error {
	a.Log().Info("Received message %v from %s", message, from)
	switch msg := message.(type) {
	case ParseJsonMessage:
		var result map[string]any
		json.Unmarshal([]byte(msg.json), &result)
	}
	return nil
}

func (a *JsonInterpreter) Terminate(reason error) {
	a.Log().Info("%s terminated with reason: %s", a.PID(), reason)
}
