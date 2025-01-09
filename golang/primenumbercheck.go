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

type PrimeCoordinator struct {
	act.Actor
	allPrimes   []int
	workerCount int
	received    int
	startTime   time.Time
}

type PrimeResult struct {
	primes []int
}

type PrimeWorker struct {
	act.Actor
}

type PrimeWorkerStartMessage struct {
	startNumber int
	endNumber   int
}

type PrimeCoordinatorStartMessage struct {
	workerCount int
	maxN        int
}

func factoryPrimeWorker() gen.ProcessBehavior {
	return &PrimeWorker{}
}

func (a *PrimeWorker) Init(args ...any) error {
	return nil
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}
	for i := 3; i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func (a *PrimeWorker) HandleMessage(from gen.PID, message any) error {
	switch m := message.(type) {
	case PrimeWorkerStartMessage:
		{
			primes := []int{}
			for i := m.startNumber; i < m.endNumber; i++ {
				if isPrime(i) {
					primes = append(primes, i)
				}
			}

			a.Send(from, PrimeResult{primes})
		}
	}

	return nil
}

func (a *PrimeWorker) Terminate(reason error) {
	a.Log().Info("%s terminated with reason: %s", a.PID(), reason)
}

func factoryPrimeCoordinator() gen.ProcessBehavior {
	return &PrimeCoordinator{}
}

func (a *PrimeCoordinator) Init(args ...any) error {
	return nil
}

func (a *PrimeCoordinator) HandleMessage(from gen.PID, message any) error {
	switch m := message.(type) {
	case PrimeResult:
		{
			a.allPrimes = append(a.allPrimes, m.primes...)
			a.received++
			if a.received == a.workerCount {
				elapsed := time.Since(a.startTime)
				a.Log().Info("All workers finished")
				a.Log().Info("Total primes found: %d", len(a.allPrimes))
				a.Log().Info("Total time: %s", elapsed)
				a.Terminate(nil)
			}
		}
	case PrimeCoordinatorStartMessage:
		{
			a.workerCount = m.workerCount
			a.startTime = time.Now()
			workers := make([]gen.PID, a.workerCount)
			for i := 0; i < a.workerCount; i++ {
				pid, err := a.Spawn(factoryPrimeWorker, gen.ProcessOptions{})
				if err != nil {
					panic(err)
				}
				workers[i] = pid
			}

			a.Log().Info("Starting prime number calculation")
			a.startTime = time.Now()
			for i := 0; i < a.workerCount; i++ {
				start := m.maxN / a.workerCount * i
				end := m.maxN / a.workerCount * (i + 1)
				if i == a.workerCount-1 {
					end = m.maxN
				}
				a.Send(workers[i], PrimeWorkerStartMessage{start, end})
			}
		}

	}

	return nil
}

func (a *PrimeCoordinator) Terminate(reason error) {
	a.Log().Info("%s terminated with reason: %s", a.PID(), reason)
}

func checkPrimes() {
	N := 10_000_000
	workers := 10
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

	node, err := ergo.StartNode("local@localhost", options)
	if err != nil {
		panic(err)
	}

	node.Log().Info("-------------------------- LOCAL 1-1 (start) ----------------------------------")
	node.Log().Info("Go Version : %s", runtime.Version())
	node.Log().Info("CPU: %s (Physical Cores: %d)", CPU.BrandName, CPU.PhysicalCores)
	node.Log().Info("Runtime CPUs: %d", runtime.NumCPU())

	// starting prime coordinator
	pid, err := node.Spawn(factoryPrimeCoordinator, gen.ProcessOptions{})
	if err != nil {
		panic(err)
	}

	node.Send(pid, PrimeCoordinatorStartMessage{workers, N})

	node.Wait()

	// nodeping.Log().Info("Total time: %s", elapsed)

	// nodeping.Log().Info("-------------------------- LOCAL 1-1 (end) ----------------------------------")
}
