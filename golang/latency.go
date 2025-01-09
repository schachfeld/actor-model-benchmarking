package main

import (
	"fmt"
	"time"

	"ergo.services/ergo"
	"ergo.services/ergo/gen"
)

func testLatency() {
	// prepare node
	options := gen.NodeOptions{}
	options.Network.Cookie = "cookie"
	// loggercolored, err := colored.CreateLogger(colored.Options{
	// 	TimeFormat: time.DateTime,
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// options.Log.DefaultLogger.Disable = true
	// options.Log.Loggers = append(
	// 	options.Log.Loggers,
	// 	gen.Logger{Name: "colored", Logger: loggercolored},
	// )

	nodeping, err := ergo.StartNode("local@localhost", options)
	if err != nil {
		panic(err)
	}

	// nodeping.Log().Info("-------------------------- LOCAL 1-1 (start) ----------------------------------")
	// nodeping.Log().Info("Go Version : %s", runtime.Version())
	// nodeping.Log().Info("CPU: %s (Physical Cores: %d)", CPU.BrandName, CPU.PhysicalCores)
	// nodeping.Log().Info("Runtime CPUs: %d", runtime.NumCPU())

	// var times []time.Duration
	for i := 0; i < 1_000_000; i++ {
		start := time.Now()
		_, err := nodeping.Spawn(factoryPrimeWorker, gen.ProcessOptions{})
		if err != nil {
			panic(err)
		}

		// times = append(times, time.Since(start))

		fmt.Printf("%v,", time.Since(start).Nanoseconds())
	}

}
