package worker

import (
	"context"
	"log"
	"sync"
	"time"
)

type Result struct {
	JobID int
	State int // 0: Failed, 1: Success
}

type JobHandler func(ctx context.Context) Result

type Job struct {
	ID      int
	Payload string
	Handler JobHandler
}

const time_limit = 5 * time.Second

func Worker(ctx context.Context, workerID int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	timer := time.NewTimer(time_limit)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				// Job channels are closed
				return
			}
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(time_limit)
			log.Printf("Worker %d processed job %d", workerID, job.ID)
			results <- job.Handler(ctx)
		case <-timer.C:
			// No job received for 5 seconds, free the worker
			log.Printf("Worker %d is freed due to inactivity", workerID)
			return
		}
	}
}
