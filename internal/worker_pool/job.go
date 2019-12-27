package worker_pool

type Job interface {
	GetId() uint64
	Execute() error
}
