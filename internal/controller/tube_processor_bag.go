package controller

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type tubeProcessorBag struct {
	mu         sync.Mutex
	processors map[string]*tubeProcessor
}

func newTubeProcessorBag() *tubeProcessorBag {
	return &tubeProcessorBag{
		processors: make(map[string]*tubeProcessor),
	}
}

func (bag *tubeProcessorBag) Add(name string, processor *tubeProcessor) {
	defer bag.mu.Unlock()
	bag.mu.Lock()

	if _, ok := bag.processors[name]; ok {
		return
	}

	logrus.Debugf("Registering processor %s", name)
	bag.processors[name] = processor
	go processor.process()
}

func (bag *tubeProcessorBag) remove(name string) {
	defer bag.mu.Unlock()
	bag.mu.Lock()

	delete(bag.processors, name)
}
