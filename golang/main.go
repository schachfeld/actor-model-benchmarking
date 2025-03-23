package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {

	workerCount := 1
	if len(os.Args) > 1 {
		var err error
		workerCount, err = strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Println("Invalid worker count, using default of 1")
			workerCount = 1
		}
	}
	checkPrimes(workerCount)
	// testThroughput()
	// testLatency()
}
