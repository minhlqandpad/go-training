package worker

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Println("worker", id, "started job", job.ID)
		response := fmt.Sprintf("This is %v", job.payload)
		time.Sleep(time.Second)
		fmt.Println("workder", id, "finished job", job.ID)
		results <- Result{workerID: id, jobID: job.ID, response: response}
	}
}

func StartWokrerPool(numJobs, numWorkers int) {
	var wg sync.WaitGroup

	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	go func() {
		for i := 1; i <= numJobs; i++ {
			jobs <- Job{ID: i, payload: fmt.Sprintf("job id: %d", i)}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		fmt.Println(res.response)
	}

	fmt.Println("All jobs done")
}
