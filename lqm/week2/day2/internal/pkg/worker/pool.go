package worker

import (
	w "github.com/tuannguyenandpadcojp/go-training/lqm/week2/day2/pkg/worker"
)

type Config struct {
	PoolSize    int
	PoolMin     int
	MaxJobs     int
	NonBlocking bool
}

type pool struct {
	*w.Pool
}

func NewPool(cfg Config) (WorkerPool, error) {
	opts := []w.PoolOpt{}
	if cfg.NonBlocking {
		opts = append(opts, w.WithNonBlocking)
	}

	newPool, err := w.NewWorkerPool(cfg.MaxJobs, cfg.PoolSize, cfg.PoolMin, opts...)
	if err != nil {
		return nil, err
	}
	return &pool{newPool}, nil
}

func (p *pool) Submit(job Job) error {
	h := func() w.Result {
		r := job.Handler()()
		state := 0
		if r.Success() {
			state = 1
		}
		return w.Result{
			JobID: r.JobID(),
			State: state,
		}
	}
	j := w.Job{
		ID:      job.ID(),
		Handler: w.JobHandler(h),
	}
	return p.Pool.Submit(j)
}
