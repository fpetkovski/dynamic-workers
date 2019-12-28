package worker_pool_test

import (
	"fpetkovski/worker_pool/internal/worker_pool"
	log "github.com/sirupsen/logrus"
	"sync"
	"testing"
)

var mux sync.Mutex

type jobStub struct {
	index   int
	counter *int
}

func (stub jobStub) Execute() error {
	defer mux.Unlock()
	mux.Lock()

	*stub.counter += 1
	return nil
}

func (stub jobStub) GetId() uint64 {
	return uint64(stub.index)
}

func TestPool_Start_ShouldGracefullyShutDown(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	counter := 0
	mux = sync.Mutex{}

	numberOfJobs := 10

	pool := makePool(2, numberOfJobs)
	go pool.Start()

	for i := 0; i < numberOfJobs; i++ {
		pool.AddJob(&jobStub{i, &counter})
	}

	pool.Close()
	pool.Wait()

	if counter != numberOfJobs {
		t.Errorf("Expected to process %d jobs, processed %d", numberOfJobs, counter)
	}
}

func makePool(numberOfWorkers int, numberOfJobs int) worker_pool.WorkerPool {
	complete := func(uint64) {}
	pool := worker_pool.NewWorkerPool(numberOfWorkers, complete)

	return pool
}
