package scheduler

import (
	"github.com/nokka/go-scheduler/job"
	"github.com/nokka/go-scheduler/worker"
)

// Scheduler provides operations on the worker queue, and determines
// who runs their job at any given time.
type Scheduler interface {
	// Run will start the scheduler.
	Run()
}

type scheduler struct {
	// workerPool holds all the workers available to perform a job.
	workerPool chan chan job.Performer

	// maxWorkers is the allowed size of the worker pool.
	maxWorkers int

	// queue is a channel that we can send work on, and it will be performed
	// when a worker is available to perform the job.
	queue chan job.Performer
}

// Run will setup as many workers as allowed.
func (s *scheduler) Run() {
	for i := 0; i < s.maxWorkers; i++ {
		worker := worker.New(s.workerPool)
		worker.Start()
	}

	go s.dispatch()
}

func (s *scheduler) dispatch() {
	for {
		select {
		// Listen for jobs on the queue.
		case j := <-s.queue:
			// Got a job request on the job queue.
			go func(j job.Performer) {
				// Obtain a worker channel from the pool, will stall
				// until a worker is available.
				worker := <-s.workerPool

				// Send the job to be executed on the worker we received.
				worker <- j
			}(j)
		}
	}
}

// New initializes a scheduler with all dependencies.
func New(queue chan job.Performer, maxWorkers int) Scheduler {
	pool := make(chan chan job.Performer, maxWorkers)
	return &scheduler{
		workerPool: pool,
		maxWorkers: maxWorkers,
		queue:      queue,
	}
}
