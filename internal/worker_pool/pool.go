package worker_pool

import (
	"context"
	"sync"
)

type pool struct {
	queue    *queue
	workers  []*worker
	doneJobs chan<- uint64
	wg       sync.WaitGroup
}

type completeFunc func(uint64)

func NewWorkerPool(concurrency int, completeFunc completeFunc) *pool {
	jobQueue := NewJobQueue(concurrency)

	workers := make([]*worker, concurrency)
	for i := 0; i < concurrency; i++ {
		workers[i] = newWorker(jobQueue, completeFunc)
	}

	wg := sync.WaitGroup{}
	wg.Add(concurrency)

	return &pool{
		queue:    jobQueue,
		workers:  workers,
		wg: wg,
	}
}

func (pool *pool) AddJob(job Job) {
	pool.queue.Enqueue(job)
}

func (pool *pool) Start(ctx context.Context) {
	for _, worker := range pool.workers {
		go worker.Start(ctx, pool.wg.Done)
	}
}

func (pool *pool) Wait() {
	pool.wg.Wait()
}

type WorkerPool interface {
	AddJob(job Job)
	Start(ctx context.Context)
	Wait()
}
