package main

import (
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type Coordinator struct {
	act.Supervisor
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
			act.SupervisorChildSpec{
				Name:    "fileReader",
				Factory: fileReaderFactory,
			},
		},
	}
	return spec, nil
}

func (s *Coordinator) HandleMessage(from gen.PID, message any) error {
	switch message.(type) {
	case CoordinatorStartMessage:
		// start the child process
		pid, err := s.Spawn(fileReaderFactory, gen.ProcessOptions{})
		if err != nil {
			return err
		}
		s.Send(pid, ReadFileMessage{filename: "Output1.txt"})
	}
	return nil
}

func (s *Coordinator) Terminate(reason error) {
	s.Log().Info("%s terminated with reason: %s", s.PID(), reason)
}
