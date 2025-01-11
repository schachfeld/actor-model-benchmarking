package main

import (
	"runtime"
	"time"

	"ergo.services/ergo"
	"ergo.services/ergo/gen"
	"ergo.services/logger/colored"
	. "github.com/klauspost/cpuid/v2"
)

func main() {
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
	pid, err := nodeping.Spawn(coordinatorFactory, gen.ProcessOptions{})
	if err != nil {
		panic(err)
	}
	nodeping.Log().Info("Started File Reader Actor with PID: %s", pid)

	// send start message to the process
	err = nodeping.Send(pid, CoordinatorStartMessage{})
	if err != nil {
		panic(err)
	}

	nodeping.Wait()

}
