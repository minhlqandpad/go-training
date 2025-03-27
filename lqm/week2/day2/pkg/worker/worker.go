package worker

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type Result struct {
	JobID string
	State int // 0: Failed, 1: Success
}

type JobHandler func() Result

type Job struct {
	ID      string
	Payload string
	Handler JobHandler
}

type Worker struct {
	ID        string
	IdleState bool
}

const time_limit = 3 * time.Second

func StartWorker(ctx context.Context, workerID string, minWorkers int, activeWorkers *int32, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup, workersPtr *sync.Map, freeLock *sync.Mutex) {
	workers := workersPtr
	timer := time.NewTimer(time_limit)
	defer func() {
		wg.Done()
		atomic.AddInt32(activeWorkers, -1)
		timer.Stop()
	}()

	for {
		select {
		case <-ctx.Done():
			return

		case job, ok := <-jobs:
			if !ok {
				// Job channels are closed
				return
			}
			timer.Reset(time_limit)

			results <- job.Handler()
			log.Printf("Worker %v processed job %v", workerID, job.ID)

			if value, ok := workers.Load(workerID); ok {
				worker := value.(Worker)
				worker.IdleState = true
				workers.Store(workerID, worker)
			}

		case <-timer.C:
			// Time out, free the worker
			freeLock.Lock()
			// defer freeLock.Unlock() // Does not work if the worker is not free (stuck at timer.Reset)
			if atomic.LoadInt32(activeWorkers) > int32(minWorkers) {
				log.Printf("Worker %v is freed due to inactivity", workerID)
				workers.Delete(workerID)
				freeLock.Unlock()
				return
			}
			timer.Reset(time_limit) // Does not work with defer
			freeLock.Unlock()
		}
	}
}
