package main

import (
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

var (
	MaxWorker = os.Getenv("MAX_WORKERS")
	MaxQueue  = os.Getenv("MAX_QUEUE")
)

type Job struct {
	id int
}

var JobQueue = make(chan Job)

type Worker struct {
	ID         int
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

func NewWorker(id int, workerPool chan chan Job) Worker {
	return Worker{
		ID:         id,
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				log.Printf("work %+v processing job %v start\n", w.ID, job.id)
				log.Printf("work %+v processing job %v finished\n", w.ID, job.id)
			case <-w.quit:
				return
			}
		}
	}()
}
func (w *Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

type Dispatcher struct {
	WorkerPool chan chan Job
}

func NewDispatcher(maxworkers int) *Dispatcher {
	pool := make(chan chan Job, maxworkers)
	return &Dispatcher{WorkerPool: pool}
}

func (d *Dispatcher) Run() {
	for i := 0; i < MAX_WORKERS; i++ {
		worker := NewWorker(i, d.WorkerPool)
		worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		log.Printf("want received job ")
		select {
		case job := <-JobQueue:
			log.Printf("received job %+v", job)
			// a job request has been received
			go func(job Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.WorkerPool

				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		}
	}
}

func (d *Dispatcher) AddJob(job Job) {
	JobQueue <- job
}

var MAX_WORKERS = 5

func main() {
	var taskid int64

	log.Printf("maxworkers: %v, maxquer:  %v", MaxWorker, MaxQueue)

	dispath := NewDispatcher(MAX_WORKERS)
	dispath.Run()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		atomic.AddInt64(&taskid, 1)

		job := Job{id: int(taskid)}
		log.Printf("generate job: %v", job)
		JobQueue <- Job{id: int(taskid)}

		//go dispath.AddJob(job)

		c.JSON(http.StatusOK, gin.H{})
	})
	r.Run(":18890")
}
