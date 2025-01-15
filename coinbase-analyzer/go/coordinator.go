package main

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type Coordinator struct {
	act.Supervisor
	fileReader             gen.PID
	jsonInterpreter        gen.PID
	avgOrderBookCalculator gen.PID
}

type CoordinatorStartMessage struct {
}

type DistributeJsonMessage struct {
	cbMessage CBMessage
}

func coordinatorFactory() gen.ProcessBehavior {
	return &Coordinator{}
}

func (s *Coordinator) Init(args ...any) (act.SupervisorSpec, error) {
	s.Log().Info("initialize supervisor...")
	spec := act.SupervisorSpec{
		Children: []act.SupervisorChildSpec{
			{
				Name:    "fileReader",
				Factory: fileReaderFactory,
			},
			{
				Name:    "jsonInterpreter",
				Factory: jsonInterpreterFactory,
			},
		},
	}
	return spec, nil
}

func (s *Coordinator) HandleMessage(from gen.PID, message any) error {
	// s.Log().Info("Received message from %s", from)
	switch message := message.(type) {
	case CoordinatorStartMessage:
		{ // start the child actors
			pid, err := s.Spawn(fileReaderFactory, gen.ProcessOptions{})
			if err != nil {
				return err
			}
			s.fileReader = pid
			s.Log().Info("fileReader started")

			pid, err = s.Spawn(jsonInterpreterFactory, gen.ProcessOptions{})
			if err != nil {
				return err
			}
			s.jsonInterpreter = pid
			s.Log().Info("jsonInterpreter started")

			pid, err = s.Spawn(avgOrderBookCalculatorFactory, gen.ProcessOptions{})
			if err != nil {
				return err
			}
			s.avgOrderBookCalculator = pid
			s.Log().Info("avgOrderBookCalculator started")

			// Start the test
			s.Send(s.fileReader, ReadFileMessage{filename: "messages.log"})
		}

	case ParseJsonMessage:
		{
			s.Send(s.jsonInterpreter, message)
		}
	case DistributeJsonMessage:
		{
			if len(message.cbMessage.Events) > 0 {
				updates := message.cbMessage.Events[0].Updates

				s.Send(s.avgOrderBookCalculator, UpdatesMessage{updates: updates})
			}
		}
	}
	return nil
}

func (s *Coordinator) Terminate(reason error) {
	s.Log().Info("%s terminated with reason: %s", s.PID(), reason)
}
