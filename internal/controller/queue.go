package controller

type queue struct {
	adapter QueueAdapter
	name    string
}

func newQueue(queueAdapter QueueAdapter, queueName string) *queue {
	return &queue{
		adapter: queueAdapter,
		name:    queueName,
	}
}

func (queue *queue) GetName() string {
	return queue.name
}

func (queue *queue) GetJob() (*Job, error) {
	return queue.adapter.GetJob(queue.name)
}

func (queue *queue) DeleteJob(jobId uint64) {
	queue.adapter.DeleteJob(jobId)
}
