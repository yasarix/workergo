# WorkerGo

[![Build Status](https://travis-ci.org/yasarix/workergo.svg?branch=master)](https://travis-ci.org/yasarix/workergo)
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/yasarix/workergo)
[![Go Report Card](https://goreportcard.com/badge/github.com/yasarix/workergo)](https://goreportcard.com/report/github.com/yasarix/workergo)

WorkerGo is an MIT licensed worker pool implementation that can be used in any Go program to handle tasks with workers.

WorkerGo is heavily influenced by the Marcio Catilho's post here: [http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/)

I was trying to write a worker pool implementation that I could use in a program with different portions of it will require parallel processing. I found his post while researching, created a new package using portions of his samples in my project. Since the package I created can be used for calling any struct with a method, I thought it would be good to share, so, it can be used in any program that needs a worker pool implementation.

## Installation

    go get gopkg.in/yasarix/workergo.v2

## Usage

First, a `Dispatcher` needs to be created.

    maxWorkers := 5 // Maximum number of workers
    queueBufferSize := 200 // Buffer size for job queue
    d := workergo.New(maxWorkers, queueBufferSize)

    // Start dispatcher
    d.Start()

Workergo accepts any struct that implements `Job` interface that you can see below:

    type Job interface {
        Run()
    }

Then, a struct that satisfies `Job` interface can be created and sent to job queue:

    type MyJob struct {
        id int
    }

    func (j *MyJob) Run() {
        fmt.Printf("This is job id: %d", id)
        time.Sleep(5 * time.Second)
    }

    ...
    ...
    func main() {
        d := workergo.New(5, 100)
        d.Start()
        for i := 0; i < 13; i++ {
            j := MyJob{id: i}
            d.Submit(&j)
        }

        // Wait for workers to finish executing `j.Run()` methods
        d.Wait()
    }

## Setting a rate limiter

You can also run the dispatcher with a rate limiter value. Simply, pass a `time.Duration` while creating the dispatcher:

    d := workergo.New(5, 100, workergo.RateLimit(time.Millisecond * 500))

Now, each job that you have submitted will be dispatched with 0.5 seconds delay.

## Stopping the dispatcher

Dispatcher can be stopped gracefully by calling `Stop()` method of `Dispatcher`. Dispatcher will wait until the actively running jobs to finish, and then stop dispatching rest of the jobs in the queue, and will shut down the workers.

## Code Documentation

[https://godoc.org/github.com/yasarix/workergo](https://godoc.org/github.com/yasarix/workergo)
