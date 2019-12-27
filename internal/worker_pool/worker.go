package worker_pool

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type worker struct {
	index        int
	queue        *queue
	completeFunc completeFunc
}

var workerIndex int

func newWorker(queue *queue, completeFunc completeFunc) *worker {
	worker := &worker{
		index:        workerIndex,
		queue:        queue,
		completeFunc: completeFunc,
	}

	workerIndex++

	return worker
}

func (worker worker) Start(ctx context.Context, done func()) {
	defer done()

	log.Debugf("Starting worker %d", worker.index)

	for {
		log.Tracef("Worker %d waiting for job", worker.index)
		select {
		case job := <-worker.queue.Dequeue():
			log.Infof("Worker %d executing job %d", worker.index, job.GetId())

			err := job.Execute()
			if err != nil {
				log.Infof("Job failed")
			} else {
				log.Infof("Processed job %d, marking as done", job.GetId())
				worker.completeFunc(job.GetId())
			}
			continue
		default:
		}

		select {
		case <-ctx.Done():
			log.Debugf("Shutting down worker %d", worker.index)
			return
		default:
		}
	}
}
