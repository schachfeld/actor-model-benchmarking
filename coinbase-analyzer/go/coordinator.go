package main

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type Coordinator struct {
	act.Supervisor
	fileReader      gen.PID
	jsonInterpreter gen.PID
}

type CoordinatorStartMessage struct {
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
	switch message.(type) {
	case CoordinatorStartMessage:
		{ // start the child actors
			pid, err := s.Spawn(fileReaderFactory, gen.ProcessOptions{})
			if err != nil {
				return err
			}
			s.fileReader = pid

			pid, err = s.Spawn(jsonInterpreterFactory, gen.ProcessOptions{})
			if err != nil {
				return err
			}
			s.jsonInterpreter = pid

			s.Send(s.fileReader, ReadFileMessage{filename: "Output1.txt"})
		}

	case ParseJsonMessage:
		{
			s.Send(s.jsonInterpreter, message)
		}
	}
	return nil
}

func (s *Coordinator) Terminate(reason error) {
	s.Log().Info("%s terminated with reason: %s", s.PID(), reason)
}
