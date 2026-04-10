package main

import (
	"context"
	"fmt"
	"sync"
)

type jobType string

const (
	multiply jobType = "multiply"
)

func (jt jobType) String() string {
	return string(jt)
}

type Job struct {
	id      int
	jobType jobType
	data    any // I am trying to make it like real production grade worker pool, where we expect some concrete dataType based on the job type
}

type Pool interface {
	Submit(job Job)
	ShutDown()
}
type WorkerPool struct {
	numberOfWorkers int
	jobChan         chan Job
	resultChan      chan any
	wg              sync.WaitGroup
	once            sync.Once
}

func NewWorkerPool(nOfW int, queueSize int) *WorkerPool {
	jobChan := make(chan Job, queueSize)
	resultChan := make(chan any, nOfW)
	wp := &WorkerPool{
		numberOfWorkers: nOfW,
		jobChan:         jobChan,
		resultChan:      resultChan,
	}
	for i := 0; i < nOfW; i++ {
		wp.wg.Add(1)
		go worker(i, jobChan, resultChan)
	}
	return wp
}

func (wp *WorkerPool) SubmitWait(ctx context.Context, job Job) {
	for {
		select {
		case wp.jobChan <- job:
			return
		case <-ctx.Done():
			fmt.Println("Time out")
		}
	}

}
func (wp *WorkerPool) Submit(job Job) {

	select {
	case wp.jobChan <- job:
	default:
		fmt.Println("Workers are busy at the moment")
	}

}

func (wp *WorkerPool) ShutDown() {
	wp.once.Do(func() {
		close(wp.jobChan) // signals workers to exit their range loop
	})
	wp.wg.Wait()         // wait for all workers to finish
	close(wp.resultChan) // now safe — no worker is writing
}

// Results returns the result channel for the caller to read from
func (wp *WorkerPool) Results() <-chan any {
	return wp.resultChan
}
func worker(workerID int, job <-chan Job, result chan<- any) {
	fmt.Println("This is worker", workerID)
	for job := range job {
		fmt.Println("Excuting Job ID", job.id, "Job Type", job.jobType.String())
		switch job.jobType {
		case multiply:
			jd := job.data.(int)
			result <- jd * jd
		default:
			fmt.Println("can handle other job types as well")
		}
	}
}
