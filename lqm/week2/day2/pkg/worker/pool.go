package worker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"sync/atomic"
)

type Pool struct {
	// Mutex to protect the pool state
	mutex sync.Mutex

	// Context and cancel function to notify all workers to stop
	ctx        context.Context
	cancelFunc context.CancelFunc

	// Defines the number of workers and job queue
	workers       sync.Map
	jobs          chan Job
	maxWorkers    int
	minWorkers    int
	activeWorkers int32

	// Wait group to wait for all workers to finish
	workerWaitGroup sync.WaitGroup

	// Wait group to wait for the goroutine aggregate the job result
	resultWaitGroup sync.WaitGroup
	submitLock      sync.Mutex

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

func NewWorkerPool(maxJobs, maxWorkers, minWorkers int, opts ...PoolOpt) (*Pool, error) {
	if maxWorkers < 1 {
		return nil, ErrNegativeInput(maxWorkers)
	}
	if maxJobs < 1 {
		return nil, ErrNegativeInput(maxJobs)
	}
	p := &Pool{
		activeWorkers:   int32(maxWorkers),
		maxWorkers:      maxWorkers,
		minWorkers:      minWorkers,
		jobs:            make(chan Job, maxJobs),
		results:         make(chan Result, maxJobs),
		workerWaitGroup: sync.WaitGroup{},
		resultWaitGroup: sync.WaitGroup{},
		mutex:           sync.Mutex{},
	}
	for i := range maxWorkers {
		p.workers.Store(strconv.Itoa(i), Worker{ID: strconv.Itoa(i), IdleState: true})
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
	p.workerWaitGroup.Add(p.maxWorkers)
	atomic.StoreInt32(&p.activeWorkers, int32(p.maxWorkers))
	log.Printf("Init: %d active, %d min", atomic.LoadInt32(&p.activeWorkers), p.minWorkers)

	// spawn worker goroutine
	for i := range p.maxWorkers {
		id := strconv.Itoa(i)
		go StartWorker(p.ctx, id, p.minWorkers, &p.activeWorkers, p.jobs, p.results, &p.workerWaitGroup, &p.workers, &p.submitLock)
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

// / Submit is a dispatcher feeds jobs for the worker
func (p *Pool) Submit(job Job) error {
	if !p.running {
		return fmt.Errorf("pool is closed")
	}

	p.submitLock.Lock()
	defer p.submitLock.Unlock()

	idleWorkerID := p.FindIdleWorker()
	
	if idleWorkerID == "-1" && atomic.LoadInt32(&p.activeWorkers) < int32(p.maxWorkers) {
		newWorkerID := p.NewWorker()
		log.Printf("Not enough workers to handle jobs. Adding worker ID %v", newWorkerID)
		atomic.AddInt32(&p.activeWorkers, 1)
		p.workerWaitGroup.Add(1)
		go StartWorker(p.ctx, newWorkerID, p.minWorkers, &p.activeWorkers, p.jobs, p.results, &p.workerWaitGroup, &p.workers, &p.submitLock)
	} else if idleWorkerID != "-1" {
		value, ok := p.workers.Load(idleWorkerID)
		if !ok {
			log.Printf("Worker %v not found", idleWorkerID)
			return errors.New("internal error")
		}
		worker := value.(Worker)
		worker.IdleState = false
		p.workers.Store(idleWorkerID, worker)
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

func (p *Pool) FindIdleWorker() string {
	var idleWorkerID = "-1"
	p.workers.Range(func(key, value any) bool {
		worker := value.(Worker)
		if worker.IdleState {
			idleWorkerID = key.(string)
			log.Printf("Found worker ID %v", idleWorkerID)
			return false // Found an idle worker, break the loop
		}
		return true
	})
	return idleWorkerID
}

func (p *Pool) NewWorker() string {
	for i := range p.maxWorkers {
		id := strconv.Itoa(i)
		if _, ok := p.workers.Load(id); !ok {
			newWorker := Worker{
				ID:        id,
				IdleState: false,
			}
			p.workers.Store(id, newWorker)
			return id
		}
	}
	return "-1"
}

func PrintMap(sm *sync.Map) {
	sm.Range(func(key, value any) bool {
		log.Printf("Worker ID: %v, Value: %v", key, value)
		return true
	})
}
