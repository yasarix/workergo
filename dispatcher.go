package workergo

import (
	"sync"
	"time"
)

// Dispatcher is the main code that runs and starts workers. Dispatcher is
// resposible of managing workers, dispatching them when needed
type Dispatcher struct {
	maxWorkers int
	queueSize  int
	workerWG   sync.WaitGroup
	jobWG      sync.WaitGroup
	jobChan    chan Job
	workerPool chan chan Job
	workers    []*worker
	quit       chan bool
	limiter    <-chan time.Time
	running    bool
}

// New Creates a new dispatcher instance with given maximum number of workers
func New(maxWorkers, queueSize int, options ...func(*Dispatcher)) *Dispatcher {
	d := Dispatcher{
		maxWorkers: maxWorkers,
		queueSize:  queueSize,
		workerPool: make(chan chan Job, maxWorkers),
		quit:       make(chan bool),
		workers:    make([]*worker, maxWorkers),
	}

	for _, opt := range options {
		opt(&d)
	}

	d.jobChan = make(chan Job, queueSize)

	return &d
}

// RateLimit Sets duration for a rate limiter to dispatch workers
func RateLimit(dur time.Duration) func(*Dispatcher) {
	return func(d *Dispatcher) {
		d.startTicker(dur)
	}
}

func (d *Dispatcher) startTicker(dur time.Duration) {
	d.limiter = time.Tick(dur)
}

// Start Starts the dispatcher
func (d *Dispatcher) Start() {
	d.running = true
	// Start all workers
	d.workerWG.Add(d.maxWorkers)
	for i := 0; i < d.maxWorkers; i++ {
		d.workers[i] = newWorker(d.workerPool, &d.jobWG, &d.workerWG)
		d.workers[i].start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	defer d.shutdown()

	for {
		select {
		case job := <-d.jobChan:
			// Get available workers channel
			workerjobChannel := <-d.workerPool

			if d.limiter != nil {
				<-d.limiter
			}

			// Send job to that channel
			workerjobChannel <- job
		case <-d.quit:
			return
		}
	}
}

// Submit Submits given job into job queue
func (d *Dispatcher) Submit(j Job) error {
	if !d.running {
		return ErrNotRunning
	}
	d.jobWG.Add(1)
	d.jobChan <- j
	return nil
}

// Stop Stops dispatcher
func (d *Dispatcher) Stop() {
	d.running = false
	d.quit <- true
	d.workerWG.Wait()
}

func (d *Dispatcher) shutdown() {
	for _, worker := range d.workers {
		for !worker.stop() {
		}
	}

	close(d.workerPool)
	close(d.jobChan)
}

// QueueSize returns current size of the job queue
func (d *Dispatcher) QueueSize() int {
	return len(d.jobChan)
}

// Wait for all workers to finish their jobs
func (d *Dispatcher) Wait() {
	if !d.running {
		return
	}
	d.jobWG.Wait()
}
