package workergo

import (
	"sync"
)

// Dispatcher is the main code that runs and starts workers. Dispatcher is
// resposible of managing workers, dispatching them when needed
type Dispatcher struct {
	maxWorkers int
	wg         *sync.WaitGroup
	wait       bool
	JobQueue   chan Job
	workerPool chan chan Job
	workers    []Worker
	quit       chan bool
}

// NewDispatcher Creates a new dispatcher instance with given maximum number of
// workers
func NewDispatcher(maxWorkers int) *Dispatcher {
	return &Dispatcher{
		maxWorkers: maxWorkers,
		workerPool: make(chan chan Job, maxWorkers),
		JobQueue:   make(chan Job),
		wait:       false,
		quit:       make(chan bool),
		workers:    make([]Worker, maxWorkers),
	}
}

// NewDispatcherWG Creates a new Dispatcher instance with given maximum number of
// workers and uses the given wait group to wait for workers to finish their jobs
func NewDispatcherWG(maxWorkers int, exWg *sync.WaitGroup) *Dispatcher {
	d := NewDispatcher(maxWorkers)
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

func (d *Dispatcher) dispatch() {
	defer d.shutdown()

	for {
		select {
		case job := <-d.JobQueue:
			go func(job Job) {
				// Get available workers channel
				workerJobChannel := <-d.workerPool

				// Send job to that channel
				workerJobChannel <- job
			}(job)
		case <-d.quit:
			return
		}
	}
}

// SubmitJob Submits given Job into job queue
func (d *Dispatcher) SubmitJob(job Job) {
	d.wg.Add(1)

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
}
