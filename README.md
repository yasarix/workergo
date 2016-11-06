# WorkerGo

[![Build Status](https://travis-ci.org/yasarix/workergo.svg?branch=master)](https://travis-ci.org/yasarix/workergo)
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/yasarix/workergo)
[![Go Report Card](https://goreportcard.com/badge/github.com/yasarix/workergo)](https://goreportcard.com/report/github.com/yasarix/workergo)

WorkerGo is a worker pool implementation that can be used in any Go program to handle tasks with workers. Workers created by WorkerGo calls the method of the structs sent them as a job. So, any struct with a method needs to be called in parallel can be sent to WorkerGo's job queue.

WorkerGo is heavily influenced by the Marcio Catilho's post here: [http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/)

I was trying to write a worker pool implementation that I could use in a program with different portions of it will require parallel processing. I found his post while researching, created a new package using portions of his samples in my project. Since the package I created can be used for calling any struct with a method, I thought it would be good to share, so, it can be used in any program that needs a worker pool implementation.

# Installation

```
go get gopkg.in/yasarix/workergo.v1
```

# Usage

First, a `Dispatcher` needs to be created.

```
maxWorkers := 5 // Maximum number of workers
queueBufferSize := 20 // Buffer size for job queue
d := workergo.Dispatcher(maxWorkers, queueBufferSize)
d.Run()
```

Then, a job can be created and sent to job queue:

```
work := NewLengthyWork(123, "Hello")
job := workergo.NewJob(workergo.TASK, work, "DoLengthyWork")
d.SubmitJob(job)
```

The work here should be a struct with `DoLengthyWork()` method:

```
type LengthyWork struct {
	number int
	message string
}

func NewLengthyWork(number int, message string) *LengthyWork {
	return &LengthyWork{
		number: number,
		message: message,
	}
}

func (w *LengthyWork) DoLengthyWork() {
	fmt.Println("Doing some lengthy work for", w.number, " - message:", w.message)
}
```

You can also pass the pointer of your existing `sync.waitGroup` into dispatcher to wait for workers to finish the jobs. Instead of calling `NewDispatcher()`, call `NewDispatcherWG()`:

```
var wg sync.WorkGroup
d := NewDispatcherWG(maxWorkers, queueBufferSize, &wg)
```

Whenever a new job has been submitted into the job queue, dispatcher will call `wg.Add(1)`, and once a worker finished that job, it will call `wg.Done()`.

# Setting a rate limiter

You can also run the dispatcher with a rate limiter value. Simply, call `Dispatcher.RunWithLimiter()` method and pass a desired `time.Duration` value instead of calling `Dispatcher.Run()`

Let's say, you want each job to be executed with 0.5 seconds delays:
```
d.RunWithLimiter(time.Millisecond * time.Duration(500))
```

Now, each job that you have submitted will be dispatched with 0.5 seconds delay.

# More Documentation
[https://godoc.org/github.com/yasarix/workergo](https://godoc.org/github.com/yasarix/workergo)
