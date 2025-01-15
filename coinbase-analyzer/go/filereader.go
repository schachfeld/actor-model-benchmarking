package main

import (
	"bufio"
	"os"

	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
)

type FileReader struct {
	act.Actor
}

type ReadFileMessage struct {
	filename string
}

func fileReaderFactory() gen.ProcessBehavior {
	return &FileReader{}
}

func (a *FileReader) Init(args ...any) error {
	return nil
}

func (a *FileReader) HandleMessage(from gen.PID, message any) error {
	switch msg := message.(type) {
	case ReadFileMessage:
		file, err := os.Open(msg.filename)

		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		buf := make([]byte, 0, 64*64*1024)
		scanner.Buffer(buf, 1024*1024*1024)

		for scanner.Scan() {
			a.Send(a.Parent(), ParseJsonMessage{json: scanner.Text()})
		}

		if err := scanner.Err(); err != nil {
			return err
		}

		a.Log().Info("FileReader done!")
	}
	return nil
}

func (a *FileReader) Terminate(reason error) {
	a.Log().Info("%s terminated with reason: %s", a.PID(), reason)
}
