package workergo

import (
	"sync"
	"time"
)

// Dispatcher is the main code that runs and starts workers. Dispatcher is
// resposible of managing workers, dispatching them when needed
type Dispatcher struct {
	maxWorkers      int
	queueBufferSize int
	wg              *sync.WaitGroup
	wait            bool
	JobQueue        chan Job
	workerPool      chan chan Job
	workers         []Worker
	quit            chan bool
	limiter         <-chan time.Time
}

// NewDispatcher Creates a new dispatcher instance with given maximum number of
// workers
func NewDispatcher(maxWorkers int, queueBufferSize int) *Dispatcher {
	d := &Dispatcher{
		maxWorkers:      maxWorkers,
		queueBufferSize: queueBufferSize,
		workerPool:      make(chan chan Job, maxWorkers),
		wait:            false,
		quit:            make(chan bool),
		workers:         make([]Worker, maxWorkers),
	}

	if queueBufferSize != 0 {
		d.JobQueue = make(chan Job, queueBufferSize)
	} else {
		d.JobQueue = make(chan Job)
	}

	return d
}

// NewDispatcherWG Creates a new Dispatcher instance with given maximum number of
// workers and uses the given wait group to wait for workers to finish their jobs
func NewDispatcherWG(maxWorkers int, queueBufferSize int, exWg *sync.WaitGroup) *Dispatcher {
	d := NewDispatcher(maxWorkers, queueBufferSize)
	d.wg = exWg
	d.wait = true

	return d
}

// Run Starts the dispatcher
func (d *Dispatcher) Run() {
	// Start all workers
	for i := 0; i < d.maxWorkers; i++ {
		if d.wait {
			d.workers[i] = NewWorkerWG(d.workerPool, d.wg)
		} else {
			d.workers[i] = NewWorker(d.workerPool)
		}

		d.workers[i].Run()
	}

	go d.dispatch()
}

// RunWithLimiter Adds delays between jobs
func (d *Dispatcher) RunWithLimiter(limiterGap time.Duration) {
	d.limiter = time.Tick(limiterGap)
	d.Run()
}

func (d *Dispatcher) dispatch() {
	defer d.shutdown()

	for {
		select {
		case job := <-d.JobQueue:
			// Get available workers channel
			workerJobChannel := <-d.workerPool

			if d.limiter != nil {
				<-d.limiter
			}

			// Send job to that channel
			workerJobChannel <- job
		case <-d.quit:
			return
		}
	}
}

// SubmitJob Submits given Job into job queue
func (d *Dispatcher) SubmitJob(job Job) {
	if d.wait {
		d.wg.Add(1)
	}

	// Send to job queue
	d.JobQueue <- job
}

// Stop Stops dispatcher
func (d *Dispatcher) Stop() {
	go func() {
		d.quit <- true
	}()
}

func (d *Dispatcher) shutdown() {
	for _, worker := range d.workers {
		for !worker.Stop() {

		}
	}

	close(d.workerPool)
	close(d.JobQueue)
}

// GetQueueSize returns current size of the job queue
func (d *Dispatcher) GetQueueSize() int {
	return len(d.JobQueue)
}
