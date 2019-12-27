package worker_pool_test

import (
	"context"
	"fpetkovski/worker_pool/internal/worker_pool"
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

type jobStub struct {
	counter *int
}

func (stub jobStub) Execute() error {
	*stub.counter += 1
	return nil
}

func (stub jobStub) GetId() uint64 {
	return 1
}

func TestPool_Start_ShouldGracefullyShutDown(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	numberOfJobs := 10
	counter := 0

	pool := makePool(2, numberOfJobs)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	go pool.Start(ctx)

	for i := 0; i < numberOfJobs; i++ {
		job := &jobStub{&counter}
		pool.AddJob(job)
	}

	cancel()
	pool.Wait()

	if counter != numberOfJobs {
		t.Errorf("Expected to process %d jobs, processed %d", numberOfJobs, counter)
	}
}

func makePool(numberOfWorkers int, numberOfJobs int) worker_pool.WorkerPool {
	doneChannel := make(chan<- uint64, numberOfJobs)
	pool := worker_pool.NewWorkerPool(numberOfWorkers, doneChannel)

	return pool
}
