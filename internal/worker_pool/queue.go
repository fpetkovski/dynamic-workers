package worker_pool

type queue struct {
	queue chan Job
}

func NewJobQueue(size int) *queue {
	return &queue{
		queue: make(chan Job, size),
	}
}

func (q *queue) Enqueue(job Job) {
	q.queue <- job
}

func (q *queue) Dequeue() <-chan Job {
	return q.queue
}
