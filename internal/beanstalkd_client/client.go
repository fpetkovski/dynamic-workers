package beanstalkd_client

import (
	"fpetkovski/worker_pool/internal/controller"
	"github.com/beanstalkd/go-beanstalk"
	"strconv"
)

type beanstalkdClient struct {
	connection *beanstalk.Conn
}

func NewBeanstalkdClient(host string) *beanstalkdClient {
	connection, err := beanstalk.Dial("tcp", host)
	if err != nil {
		panic(err)
	}

	return &beanstalkdClient{
		connection: connection,
	}
}

func (client beanstalkdClient) GetQueues() []string {
	tubes, err := client.connection.ListTubes()
	if err != nil {
		panic(err)
	}

	return removeItem(tubes, "default")
}

func (client beanstalkdClient) GetJob(tubeName string) (*controller.Job, error) {
	tubeSet := beanstalk.NewTubeSet(client.connection, tubeName)
	id, body, err := tubeSet.Reserve(0)
	if err != nil {
		return nil, err
	}

	return controller.NewJob(id, body), nil
}

func (client beanstalkdClient) GetReadyJobCount(tubeName string) int {
	tube := beanstalk.Tube{
		Conn: client.connection,
		Name: tubeName,
	}

	stats, err := tube.Stats()
	if err != nil {
		return 0
	}

	jobCount, err := strconv.Atoi(stats["current-jobs-ready"])
	if err != nil {
		return 0
	}

	return jobCount
}

func (client beanstalkdClient) DeleteJob(jobId uint64) {
	err := client.connection.Delete(jobId)
	if err != nil {
		panic(err)
	}
}

func removeItem(items []string, item string) []string {
	for i, v := range items {
		if v == item {
			return append(items[:i], items[i+1:]...)
		}
	}

	return items
}
