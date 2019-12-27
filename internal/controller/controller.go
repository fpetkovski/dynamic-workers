package controller

import (
	"fpetkovski/worker_pool/internal/beanstalkd_client"
	"github.com/sirupsen/logrus"
	"time"
)

type controller struct {
	beanstalkdClient beanstalkd_client.BeanstalkdClient
	processorsBag    *tubeProcessorBag
}

func NewController(client beanstalkd_client.BeanstalkdClient) *controller {
	return &controller{
		beanstalkdClient: client,
		processorsBag:    newTubeProcessorBag(),
	}
}

func (controller *controller) Loop() {
	for {
		logrus.Debugf("Processing tubes")
		tubes := controller.beanstalkdClient.GetTubes()
		controller.createNewProcessors(tubes)

		time.Sleep(1 * time.Second)
	}
}

func (controller *controller) createNewProcessors(tubes []string) {
	for _, tubeName := range tubes {
		tubeProcessor := newTubeProcessor(
			controller.beanstalkdClient,
			controller.beanstalkdClient.GetTubeSet(tubeName),
			controller.removeProcessor(tubeName),
		)
		controller.processorsBag.Add(tubeName, tubeProcessor)
	}
}

func (controller *controller) removeProcessor(tubeName string) func() {
	return func() {
		logrus.Debugf("Deregistering processor %s", tubeName)
		controller.processorsBag.remove(tubeName)
	}
}
