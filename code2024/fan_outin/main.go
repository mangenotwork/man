package main

import (
	"fmt"
	"sync"
	"time"
)

/*

为什么要使用  Fan-Out/Fan-In 模式
扇出/扇入模式通过将复杂的并发任务分解成较小、可管理的部分，简化了任务。它提供了一种结构化的方式来在多个 goroutine 之间分配工作，
并高效地收集结果。当有大量任务需要同时执行，并希望确保结果正确聚合时，这种模式尤其有用。

*/

func worker(id int, jobs <-chan int, results chan<- int) {
	for job := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, job)
		// Simulate some work
		time.Sleep(time.Millisecond * 500)
		results <- job * 2 // Send the result to the results channel
	}
}

func main() {
	const numJobs = 10
	const numWorkers = 3

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	var wg sync.WaitGroup

	// Start workers
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(id, jobs, results)
		}(i)
	}

	// Send jobs to workers
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	// Collect results from workers
	go func() {
		wg.Wait()
		close(results)
	}()

	// Print results
	for r := range results {
		fmt.Printf("Result: %d\n", r)
	}
}
