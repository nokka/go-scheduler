package worker

import (
	"github.com/nokka/go-scheduler/job"
)

// Worker represents a working unit that can perform work.
type Worker struct {
	WorkerPool chan chan job.Performer
	JobChannel chan job.Performer
	quit       chan bool
}

// New returns a new worker with all dependencies.
func New(workerPool chan chan job.Performer) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan job.Performer),
		quit:       make(chan bool),
	}
}

// Start will start the worker to listen for work on its job channel.
func (w Worker) Start() {
	go func() {
		for {
			// Register the work channel on the worker pool.
			w.WorkerPool <- w.JobChannel

			select {
			// Received a job to run.
			case job := <-w.JobChannel:
				job.Run()
			case <-w.quit:
				return
			}
		}
	}()
}

// Stop will send the quit signal to the worker,
// and stop it from taking any more work.
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
