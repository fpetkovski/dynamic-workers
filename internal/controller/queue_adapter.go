package controller

type QueueAdapter interface {
	GetQueues() []string
	GetJob(queueName string) (*Job, error)
	GetReadyJobCount(queueName string) int
	DeleteJob(uint64)
}
