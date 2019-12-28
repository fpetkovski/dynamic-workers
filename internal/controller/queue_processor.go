package controller

import (
	"fpetkovski/worker_pool/internal/worker_pool"
	log "github.com/sirupsen/logrus"
	"math/rand"
)

type deregisterFunc func()

type queueProcessor struct {
	queue      *queue
	deregister deregisterFunc
}

func newQueueProcessor(queue *queue, deregister deregisterFunc) *queueProcessor {
	return &queueProcessor{
		queue:      queue,
		deregister: deregister,
	}
}

func (queueProcessor queueProcessor) process() {
	log.Debugln("Starting tube processor")

	concurrency := rand.Intn(10) + 1
	workerPool := worker_pool.NewWorkerPool(concurrency, queueProcessor.onComplete)
	go workerPool.Start()

	for {
		err := queueProcessor.feedJob(workerPool)
		if err != nil {
			queueProcessor.stop(err)
			workerPool.Wait()
			return
		}
	}
}

func (queueProcessor queueProcessor) feedJob(workerPool worker_pool.WorkerPool) error {
	log.Debugf("Getting Job from tube %s", queueProcessor.queue.GetName())
	job, err := queueProcessor.queue.GetJob()

	if err != nil {
		return err
	}

	log.Debugf("Adding Job to pool %s", queueProcessor.queue.GetName())
	workerPool.AddJob(job)
	return nil
}

func (queueProcessor *queueProcessor) onComplete(jobId uint64) {
	log.Infof("Deleting Job %d", jobId)
	queueProcessor.queue.DeleteJob(jobId)
}

func (queueProcessor queueProcessor) stop(err error) {
	log.Errorf("Error getting Job from tube %s: %s", queueProcessor.queue.GetName(), err.Error())

	queueProcessor.deregister()
}
