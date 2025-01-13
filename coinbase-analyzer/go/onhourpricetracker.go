package main

import (
	"bufio"
	"os"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type OneHourPriceTracker struct {
	act.Actor
}

type XYZMessage struct {
	filename string
}

func oneHourPriceTrackerFactory() gen.ProcessBehavior {
	return &OneHourPriceTracker{}
}

func (a *OneHourPriceTracker) Init(args ...any) error {
	return nil
}

func (a *OneHourPriceTracker) HandleMessage(from gen.PID, message any) error {
	a.Log().Info("Received message %v from %s", message, from)
	switch msg := message.(type) {
	case XYZMessage:
		file, err := os.Open(msg.filename)

		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		buf := make([]byte, 0, 64*64*1024)
		scanner.Buffer(buf, 1024*1024*1024)

		// TODO: remove this, it's only for testing purpouses
		for scanner.Scan() {
			a.Send(a.Parent(), ParseJsonMessage{json: scanner.Text()})
		}

		if err := scanner.Err(); err != nil {
			return err
		}
	}
	return nil
}

func (a *OneHourPriceTracker) Terminate(reason error) {
	a.Log().Info("%s terminated with reason: %s", a.PID(), reason)
}
