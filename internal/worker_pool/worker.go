package worker_pool

import (
	log "github.com/sirupsen/logrus"
)

type worker struct {
	index        int
	queue        <-chan Job
	completeFunc completeFunc
}

var workerIndex int

func newWorker(queue <-chan Job, completeFunc completeFunc) *worker {
	worker := &worker{
		index:        workerIndex,
		queue:        queue,
		completeFunc: completeFunc,
	}
	workerIndex++

	return worker
}

func (worker worker) Start(done func()) {
	defer done()

	log.Debugf("Starting worker %d", worker.index)

	for job := range worker.queue {
		log.Infof("Worker %d executing job %d", worker.index, job.GetId())

		err := job.Execute()
		if err != nil {
			log.Infof("Job failed")
		} else {
			log.Infof("Processed job %d, marking as done", job.GetId())
			worker.completeFunc(job.GetId())
		}
	}
}
