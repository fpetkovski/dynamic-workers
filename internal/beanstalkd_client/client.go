package beanstalkd_client

import (
	"github.com/beanstalkd/go-beanstalk"
	"time"
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

func (client beanstalkdClient) GetTubes() []string {
	tubes, err := client.connection.ListTubes()
	if err != nil {
		panic(err)
	}

	return removeItem(tubes, "default")
}

func removeItem(items []string, item string) []string {
	for i, v := range items {
		if v == item {
			return append(items[:i], items[i+1:]...)
		}
	}

	return items
}

func (client beanstalkdClient) GetFromTube(tubeName string) (id uint64, body []byte, err error) {
	tubeSet := beanstalk.NewTubeSet(client.connection, tubeName)

	return tubeSet.Reserve(50 * time.Second)
}

func (client beanstalkdClient) DeleteJob(jobId uint64) {
	err := client.connection.Delete(jobId)
	if err != nil {
		panic(err)
	}
}

type BeanstalkdClient interface {
	GetTubes() []string

	GetFromTube(tubeName string) (id uint64, body []byte, err error)

	DeleteJob(jobId uint64)
}
