package controller

import (
	"github.com/sirupsen/logrus"
	"time"
)

type controller struct {
	queueAdapter  QueueAdapter
	processorsBag *queueProcessorSet
}

func NewController(stats QueueAdapter) *controller {
	return &controller{
		queueAdapter:  stats,
		processorsBag: newTubeProcessorBag(),
	}
}

func (controller *controller) Loop() {
	for {
		logrus.Debugf("Processing queues")
		tubes := controller.queueAdapter.GetQueues()
		controller.createNewProcessors(tubes)

		time.Sleep(1 * time.Second)
	}
}

func (controller *controller) createNewProcessors(queueNames []string) {
	for _, queueName := range queueNames {
		jobCount := controller.queueAdapter.GetReadyJobCount(queueName)
		if jobCount == 0 {
			continue
		}

		tubeProcessor := newQueueProcessor(
			newQueue(controller.queueAdapter, queueName),
			controller.removeProcessor(queueName),
		)
		controller.processorsBag.Add(queueName, tubeProcessor)
	}
}

func (controller *controller) removeProcessor(queueName string) func() {
	return func() {
		logrus.Debugf("Removing processor %s", queueName)
		controller.processorsBag.remove(queueName)
	}
}
