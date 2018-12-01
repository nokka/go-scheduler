package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nokka/go-scheduler/job"

	"github.com/nokka/go-scheduler/scheduler"
)

// Job is the actual work we want to perform.
type Job struct {
	Name    string
	Payload []int
}

// Run performs work.
func (j Job) Run() {
	time.Sleep(1 * time.Second)
	fmt.Printf("Performing job with name %s\n", j.Name)
	fmt.Printf("Payload was %v\n", j.Payload)
}

func main() {
	var (
		maxWorkers   = flag.Int("max_workers", 1, "The number of workers to start")
		maxQueueSize = flag.Int("max_queue_size", 100, "The size of job queue")
	)

	flag.Parse()

	// Channel to receive errors on.
	errors := make(chan error)

	// Create the job queue we'll use through the dispatcher.
	queue := make(chan job.Performer, *maxQueueSize)

	// Create the scheduler, that will keep track of all jobs.
	scheduler := scheduler.New(queue, *maxWorkers)

	// Start the scheduler,
	scheduler.Run()

	// Create 100 jobs, and let it go through the queue.
	for i := 0; i < 10; i++ {
		// Create a job, payload is an interface{} so can be anything you need.
		job := Job{Name: fmt.Sprintf("job-%d", i), Payload: []int{0, 1, 2}}

		// Queue the job on our work queue.
		queue <- job
	}

	// Capture interrupts.
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errors <- fmt.Errorf("Got signal: %s", <-c)
	}()

	if err := <-errors; err != nil {
		log.Printf("Got error: %+v\n", err)
	}

	log.Println("terminated")
}
