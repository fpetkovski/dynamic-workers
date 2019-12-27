package controller

import (
	"github.com/sirupsen/logrus"
)

type job struct {
	jobId   uint64
	payload []byte
}

func NewJob(id uint64, payload []byte) *job {
	return &job{
		jobId:   id,
		payload: payload,
	}
}

func (job job) GetId() uint64 {
	return job.jobId
}

func (job job) Execute() error {
	logrus.Infof("Executing job with id %d: %s\n", job.jobId, string(job.payload))
	//time.Sleep(1 * time.Second)

	return nil
}
