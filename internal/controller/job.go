package controller

import (
	"github.com/sirupsen/logrus"
	"time"
)

type Job struct {
	jobId   uint64
	payload []byte
}

func NewJob(id uint64, payload []byte) *Job {
	return &Job{
		jobId:   id,
		payload: payload,
	}
}

func (job Job) GetId() uint64 {
	return job.jobId
}

func (job Job) Execute() error {
	logrus.Infof("Executing Job with id %d: %s\n", job.jobId, string(job.payload))
	time.Sleep(1 * time.Second)

	return nil
}
