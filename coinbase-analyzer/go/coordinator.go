package main

import (
	"time"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type Coordinator struct {
	act.Supervisor
	fileReader gen.PID
	startTime  time.Time
}

type CoordinatorStartMessage struct {
}

type LastMessage struct {
}

type DoneMessage struct {
}

type RouteJsonMessage struct {
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
	switch message.(type) {
	case CoordinatorStartMessage:
		{ // start the child actors
			s.startTime = time.Now()
			pid, err := s.Spawn(fileReaderFactory, gen.ProcessOptions{})
			if err != nil {
				return err
			}
			s.fileReader = pid
			s.Log().Info("fileReader started")

			// Start the test
			s.Send(s.fileReader, ReadFileMessage{filename: "../messages.log"})
		}
	case DoneMessage:
		{
			elapsed := time.Since(s.startTime)
			s.Log().Info("Coordinator received DoneMessage")
			s.Log().Info("Elapsed time: %s", elapsed)
		}
	default:
		{
			panic("Coordinator received unknown message")
		}

	}
	return nil
}

func (s *Coordinator) Terminate(reason error) {
	s.Log().Info("%s terminated with reason: %s", s.PID(), reason)
}
