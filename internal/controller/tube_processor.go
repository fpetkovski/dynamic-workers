package controller

import (
	"context"
	"fpetkovski/worker_pool/internal/beanstalkd_client"
	"fpetkovski/worker_pool/internal/worker_pool"
	"github.com/beanstalkd/go-beanstalk"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

type deregisterFunc func()

type tubeProcessor struct {
	beanstalkdClient beanstalkd_client.BeanstalkdClient
	tube             *beanstalk.TubeSet
	deregister       deregisterFunc
}

func newTubeProcessor(
	beanstalkdClient beanstalkd_client.BeanstalkdClient,
	tubeSet *beanstalk.TubeSet,
	deregister deregisterFunc,
) *tubeProcessor {
	return &tubeProcessor{
		beanstalkdClient: beanstalkdClient,
		tube:             tubeSet,
		deregister:       deregister,
	}
}

func (tubeProcessor tubeProcessor) process() {
	log.Debugln("Starting tube processor")

	ctx, cancel := context.WithCancel(context.Background())
	concurrency := rand.Intn(3) + 1
	workerPool := worker_pool.NewWorkerPool(concurrency, tubeProcessor.onComplete)
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
	log.Debugf("Getting job from tube %s", tubeProcessor.tube.Name)
	id, body, err := tubeProcessor.tube.Reserve(2 * time.Second)
	if err != nil {
		return err
	}

	log.Debugf("Adding job to pool %s", tubeProcessor.tube.Name)
	workerPool.AddJob(NewJob(id, body))
	return nil
}

func (tubeProcessor *tubeProcessor) onComplete(jobId uint64) {
	log.Infof("Deleting job %d", jobId)
	tubeProcessor.beanstalkdClient.DeleteJob(jobId)
}

func (tubeProcessor tubeProcessor) stop(err error, cancel context.CancelFunc) {
	log.Errorf("Error getting job from tube %s: %s", tubeProcessor.tube.Name, err.Error())

	cancel()
	tubeProcessor.deregister()
}
