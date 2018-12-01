package job

// Performer is the interface that can run a job.
type Performer interface {
	// Run will perform the work.
	Run()
}
