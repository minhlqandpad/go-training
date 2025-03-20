package worker

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type Pool struct {
	// Mutex to protect the pool state
	mutex sync.Mutex

	// Context and cancel function to notify all workers to stop
	ctx        context.Context
	cancelFunc context.CancelFunc

	// Defines the number of workers and job queue
	workers int
	jobs    chan Job

	// Wait group to wait for all workers to finish
	workerWaitGroup sync.WaitGroup

	// Wait group to wait for the goroutine aggregate the job result
	resultWaitGroup sync.WaitGroup

	// Pool state
	running bool

	// Non-blocking mode
	nonBlocking bool

	// Job result
	results      chan Result
	TotalSucceed int
	TotalFailed  int
}

type PoolOpt func(p *Pool)

func WithNonBlocking(p *Pool) {
	p.nonBlocking = true
}

type ErrNegativeInput int

func (e ErrNegativeInput) Error() string {
	return fmt.Sprintf("invalid input: input should be higher than 0: %d", e)
}

func NewWorkerPool(maxJobs, numWorkers int, opts ...PoolOpt) (*Pool, error) {
	if numWorkers < 1 {
		return nil, ErrNegativeInput(numWorkers)
	}
	if maxJobs < 1 {
		return nil, ErrNegativeInput(maxJobs)
	}
	p := &Pool{
		workers:         numWorkers,
		jobs:            make(chan Job, maxJobs),
		results:         make(chan Result, maxJobs),
		workerWaitGroup: sync.WaitGroup{},
		resultWaitGroup: sync.WaitGroup{},
		mutex:           sync.Mutex{},
	}
	for _, opt := range opts {
		opt(p)
	}
	return p, nil
}

func (p *Pool) Start(ctx context.Context) {
	p.mutex.Lock()

	// Unlock after done
	defer p.mutex.Unlock()

	if p.running {
		log.Println("Worker Pool is running")
		return
	}

	p.running = true
	p.ctx, p.cancelFunc = context.WithCancel(ctx)
	p.workerWaitGroup.Add(p.workers)

	// spawn worker goroutine
	for i := range p.workers {
		go Worker(p.ctx, i, p.jobs, p.results, &p.workerWaitGroup)
	}

	// aggregate job's result
	p.resultWaitGroup.Add(1)
	go func() {
		defer p.resultWaitGroup.Done()
		for result := range p.results {
			if result.State == 1 {
				p.TotalSucceed++
			} else {
				p.TotalFailed++
			}
		}
	}()
}

func (p *Pool) Release() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if !p.running {
		log.Println("Worker Pool is not running ...")
		return
	}
	// close the Jobs channel to prevent dispatcher send jobs
	close(p.jobs)
	// wait for all workers to finish processing the rest of jobs
	p.workerWaitGroup.Wait()
	// call context.CancelFunc to notify all workers to stop
	p.cancelFunc()
	// close the result channel to stop the goroutine aggregate the job's result
	close(p.results)
	// wait for the goroutine aggregate the job result is done
	p.resultWaitGroup.Wait()
	p.running = false
}

// Submit is a dispatcher feeds jobs for the worker
func (p *Pool) Submit(job Job) error {
	if !p.running {
		return fmt.Errorf("pool is closed")
	}
	if p.nonBlocking {
		select {
		case p.jobs <- job:
			// If the job channel has space, the job is sent and the function returns nil.
			return nil
		default:
			// If the job channel is full, the default case is executed and an error is returned.
			return fmt.Errorf("job queue is full")
		}
	} else {
		// blocking if the jobs channel is full
		p.jobs <- job
		return nil
	}
}

func (p *Pool) Results() (totalSucceed, totalFailed int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.TotalSucceed, p.TotalFailed
}