package workergo

import (
	"reflect"
	"sync"
)

// Worker Structure of worker
type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	wait       bool
	wg         *sync.WaitGroup
	quit       chan bool
	idle       bool
}

// NewWorker Creates a new Worker instance
func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		wait:       false,
		quit:       make(chan bool),
		idle:       true,
	}
}

// NewWorkerWG Creates a new Worker instance with pointer to a sync.WaitGroup
// instance to handle wait groups
func NewWorkerWG(workerPool chan chan Job, wg *sync.WaitGroup) Worker {
	w := NewWorker(workerPool)
	w.wait = true
	w.wg = wg

	return w
}

// Run Starts worker
func (w *Worker) Run() {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				w.idle = false

				// call target method
				if job.Type == TASK {
					reflect.ValueOf(*job.Payload).MethodByName(job.TargetFunc).Call([]reflect.Value{})
				}

				w.idle = true
				if w.wait {
					w.wg.Done()
				}
			case <-w.quit:
				return
			}
		}
	}()
}

// Stop Stops worker instance if it is idle
func (w *Worker) Stop() bool {
	if w.idle {
		go func() {
			w.quit <- true
		}()

		return true
	}

	return false
}
