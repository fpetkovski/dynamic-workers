package controller

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type queueProcessorSet struct {
	mu         sync.Mutex
	processors map[string]*queueProcessor
}

func newTubeProcessorBag() *queueProcessorSet {
	return &queueProcessorSet{
		processors: make(map[string]*queueProcessor),
	}
}

func (bag *queueProcessorSet) Add(name string, processor *queueProcessor) {
	defer bag.mu.Unlock()
	bag.mu.Lock()

	if _, ok := bag.processors[name]; ok {
		return
	}

	logrus.Debugf("Registering processor %s", name)
	bag.processors[name] = processor
	go processor.process()
}

func (bag *queueProcessorSet) remove(name string) {
	defer bag.mu.Unlock()
	bag.mu.Lock()

	delete(bag.processors, name)
}
