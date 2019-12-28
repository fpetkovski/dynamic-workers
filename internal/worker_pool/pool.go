package worker_pool

import (
	"sync"
)

type pool struct {
	queue    chan Job
	workers  []*worker
	doneJobs chan<- uint64
	wg       sync.WaitGroup
}

type completeFunc func(uint64)

func NewWorkerPool(concurrency int, completeFunc completeFunc) *pool {
	jobQueue := make(chan Job, concurrency)

	workers := make([]*worker, concurrency)
	for i := 0; i < concurrency; i++ {
		workers[i] = newWorker(jobQueue, completeFunc)
	}

	wg := sync.WaitGroup{}
	wg.Add(concurrency)

	return &pool{
		queue:   jobQueue,
		workers: workers,
		wg:      wg,
	}
}

func (pool *pool) AddJob(job Job) {
	pool.queue <- job
}

func (pool *pool) Start() {
	for _, worker := range pool.workers {
		go worker.Start(pool.wg.Done)
	}
}

func (pool *pool) Close() {
	close(pool.queue)
}

func (pool *pool) Wait() {
	pool.wg.Wait()
}

type WorkerPool interface {
	AddJob(job Job)
	Start()
	Close()
	Wait()
}
