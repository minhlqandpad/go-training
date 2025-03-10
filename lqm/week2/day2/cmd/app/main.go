package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/tuannguyenandpadcojp/go-training/lqm/week2/day2/pkg/worker"
)

func main() {
	const numJobs = 100
	const numWorkers = 20
	const maxJobs = 10

	pool, err := worker.NewWorkerPool(maxJobs, numWorkers)
	if err != nil {
		log.Println("errors: ", err)
		return
	}
	pool.Start(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Submit jobs
	go func() {
		for i := range numJobs {
			go func() {
				var jobHandler = func(ctx context.Context) worker.Result {
					return worker.Result{JobID: 1, State: 1}
				}
				job := worker.Job{ID: i, Payload: fmt.Sprintf("Job %d", i), Handler: jobHandler}
				if err := pool.Submit(job); err != nil {
					log.Printf("Failed to submit job with id: %d", job.ID)
				}
			}()
		}
	}()

	<-c
	log.Println("Total Goroutine before release the pool: ", runtime.NumGoroutine())
	pool.Release()
	log.Printf("Jobs success:%d - failed:%d", pool.TotalSucceed, pool.TotalFailed)
	time.Sleep(5 * time.Second)
	log.Println("Total Goroutine after release the pool: ", runtime.NumGoroutine())
}
