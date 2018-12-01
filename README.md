
# Go scheduler

Go scheduler lets you put work on a queue that will perform the job given when there's
a free worker in the worker pool.

#### Missing things
- Handling panics

#### Examples
Examples can be found [here](examples/scheduler/main.go) as well.

```
// Job implements the Performer interface.
type Job struct {
	Name    string
	Payload []int
}

// Run performs the given work.
func (j Job) Run() {
	time.Sleep(1 * time.Second)
	fmt.Printf("Performing job with name %s\n", j.Name)
	fmt.Printf("Payload was %v\n", j.Payload)
}

func main() {
	var (
		maxWorkers   = 5
		maxQueueSize = 100
	)

	// Channel to receive errors on.
	errors := make(chan error)

	// Create the job queue we'll use through the dispatcher.
	queue := make(chan job.Performer, *maxQueueSize)

	// Create the scheduler, that will keep track of all jobs.
	scheduler := scheduler.New(queue, *maxWorkers)

	// Start the scheduler.
	scheduler.Run()

	// Create 10 jobs, and let it go through the queue.
	for i := 0; i < 10; i++ {
		// Create a job.
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
```