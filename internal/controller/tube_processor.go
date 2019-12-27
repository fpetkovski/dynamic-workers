package controller

import (
	"context"
	"fpetkovski/worker_pool/internal/beanstalkd_client"
	"fpetkovski/worker_pool/internal/worker_pool"
	log "github.com/sirupsen/logrus"
)

type deregisterFunc func()

type tubeProcessor struct {
	beanstalkdClient beanstalkd_client.BeanstalkdClient
	tubeName         string
	deregister       deregisterFunc
}

func newTubeProcessor(
	beanstalkdClient beanstalkd_client.BeanstalkdClient,
	tubeName string,
	deregister deregisterFunc,
) *tubeProcessor {
	return &tubeProcessor{
		beanstalkdClient: beanstalkdClient,
		tubeName:         tubeName,
		deregister:       deregister,
	}
}

func (tubeProcessor tubeProcessor) process() {
	ctx, cancel := context.WithCancel(context.Background())
	workerPool := worker_pool.NewWorkerPool(1, tubeProcessor.onComplete)
	go workerPool.Start(ctx)

	for {
		err := tubeProcessor.feedJob(workerPool)
		if err != nil {
			tubeProcessor.stop(err, cancel)
			workerPool.Wait()
			return
		}
	}
}

func (tubeProcessor tubeProcessor) feedJob(workerPool worker_pool.WorkerPool) error {
	log.Debugf("Getting job from tube %s", tubeProcessor.tubeName)
	id, body, err := tubeProcessor.beanstalkdClient.GetFromTube(tubeProcessor.tubeName)
	if err != nil {
		return err
	}

	log.Debugf("Adding job to pool %s", tubeProcessor.tubeName)
	workerPool.AddJob(NewJob(id, body))
	return nil
}

func (tubeProcessor *tubeProcessor) onComplete(jobId uint64) {
	log.Infof("Deleting job %d", jobId)
	tubeProcessor.beanstalkdClient.DeleteJob(jobId)
}

func (tubeProcessor tubeProcessor) stop(err error, cancel context.CancelFunc) {
	log.Errorf("Error getting job from tube %s: %s", tubeProcessor.tubeName, err.Error())

	cancel()
	tubeProcessor.deregister()
}
