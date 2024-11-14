package main

import (
	"runtime"
	"time"

	"ergo.services/ergo"
	"ergo.services/ergo/act"
	"ergo.services/ergo/gen"
	"ergo.services/logger/colored"
	. "github.com/klauspost/cpuid/v2"
)

type MyActor struct {
	i int
	act.Actor
}

func factoryMyActor() gen.ProcessBehavior {
	return &MyActor{}
}

func (a *MyActor) Init(args ...any) error {
	a.i = 0
	return nil
}

func (a *MyActor) HandleMessage(from gen.PID, message any) error {
	a.i++
	return nil
}

func (a *MyActor) Terminate(reason error) {
	a.Log().Info("received %d messages", a.i)
	a.Log().Info("%s terminated with reason: %s", a.PID(), reason)
}

func main() {
	N := 100_000_000
	// prepare node
	options := gen.NodeOptions{}
	options.Network.Cookie = "cookie"
	loggercolored, err := colored.CreateLogger(colored.Options{
		TimeFormat: time.DateTime,
	})
	if err != nil {
		panic(err)
	}
	options.Log.DefaultLogger.Disable = true
	options.Log.Loggers = append(
		options.Log.Loggers,
		gen.Logger{Name: "colored", Logger: loggercolored},
	)

	nodeping, err := ergo.StartNode("local@localhost", options)
	if err != nil {
		panic(err)
	}

	nodeping.Log().Info("-------------------------- LOCAL 1-1 (start) ----------------------------------")
	nodeping.Log().Info("Go Version : %s", runtime.Version())
	nodeping.Log().Info("CPU: %s (Physical Cores: %d)", CPU.BrandName, CPU.PhysicalCores)
	nodeping.Log().Info("Runtime CPUs: %d", runtime.NumCPU())

	// starting 1 ping process
	pid, err := nodeping.Spawn(factoryMyActor, gen.ProcessOptions{})
	if err != nil {
		panic(err)
	}
	nodeping.Log().Info("BENCHMARK: 1 process sends %d messages to 1 process", N)

	start := time.Now()
	for i := 0; i < N; i++ {
		nodeping.Send(pid, nil)
	}

	elapsed := time.Since(start)

	nodeping.Log().Info("received %d messages. %f msg/sec", N, float64(N)/elapsed.Seconds())
	nodeping.Log().Info("Total time: %s", elapsed)

	// nodeping.Log().Info("-------------------------- LOCAL 1-1 (end) ----------------------------------")
}
