package workergo

import (
	"sync"
)

// worker Structure of worker
type worker struct {
	workerPool chan chan Job
	jobChan    chan Job
	jobWG      *sync.WaitGroup
	workerWG   *sync.WaitGroup
	quit       chan bool
	idle       bool
}

// newWorker Creates a new Worker instance
func newWorker(workerPool chan chan Job, jobWG, workerWG *sync.WaitGroup) *worker {
	return &worker{
		workerPool: workerPool,
		jobChan:    make(chan Job),
		quit:       make(chan bool),
		idle:       true,
		jobWG:      jobWG,
		workerWG:   workerWG,
	}
}

// start Starts worker
func (w *worker) start() {
	go func() {
		defer w.workerWG.Done()
		for {
			w.workerPool <- w.jobChan

			select {
			case job := <-w.jobChan:
				w.idle = false
				job.Run()
				w.idle = true
				w.jobWG.Done()
			case <-w.quit:
				return
			}
		}
	}()
}

// stop Stops worker instance if it is idle
func (w *worker) stop() bool {
	if w.idle {
		w.quit <- true
		return true
	}

	return false
}
