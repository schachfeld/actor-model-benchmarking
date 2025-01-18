package main

import (
	"fmt"
	"os"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type FileWriter struct {
	filename string
	isOpen   bool
	file     *os.File
	act.Actor
}

type WriteLineMessage struct {
	content string
}

func fileWriterFactory() gen.ProcessBehavior {
	return &FileWriter{}
}

func (a *FileWriter) Init(args ...any) error {
	if len(args) > 0 {
		if filename, ok := args[0].(string); ok {
			a.filename = filename
			a.isOpen = false

			file, err := os.Create(a.filename)
			if err != nil {
				return err
			}
			a.file = file
			a.isOpen = true

		} else {
			return fmt.Errorf("expected a string for filename, got %T", args[0])
		}
	} else {
		return fmt.Errorf("filename argument is required")
	}
	return nil
}

func (a *FileWriter) HandleMessage(from gen.PID, message any) error {
	switch msg := message.(type) {
	case WriteLineMessage:
		if a.isOpen {
			_, err := a.file.WriteString(msg.content + "\n")
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("file %s is not open", a.filename)
		}
	case LastMessage:
		{
			if a.isOpen {
				a.file.Close()
				a.isOpen = false
			}
			a.Send(a.Parent(), DoneMessage{})
			return nil
		}
	default:
		return fmt.Errorf("unknown message: %T", message)
	}

	return nil
}

func (a *FileWriter) Terminate(reason error) {
	a.Log().Info("%s terminated with reason: %s", a.PID(), reason)

	if a.isOpen {
		a.file.Close()
	}
}
