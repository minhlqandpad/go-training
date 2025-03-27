package worker

import (
	"context"
	"log"
	"sync"
)

func worker(ctx context.Context, workerID int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				// Job channels are closed
				return
			}
			log.Printf("Worker %d processed job %d", workerID, job.ID)
			results <- job.Handler(ctx)
		}
	}
}
