package main

import (
	"github.com/tuannguyenandpadcojp/go-training/lqm/week2/goroutine/worker"
)

func main() {
	const numJobs = 5
	const numWorkers = 3
	worker.StartWokrerPool(numJobs, numWorkers)
}
