package beanstalkd_client_test

import (
	"fpetkovski/worker_pool/internal/beanstalkd_client"
	"github.com/beanstalkd/go-beanstalk"
	"reflect"
	"sort"
	"testing"
)

func setUp() *beanstalk.Conn {
	connection, _ := beanstalk.Dial("tcp", "127.0.0.1:11300")
	emptyTubes(connection)

	err := connection.Close()
	if err != nil {
		panic(err)
	}

	connection, _ = beanstalk.Dial("tcp", "127.0.0.1:11300")
	return connection
}

func TestBeanstalkdClient_GetTubes(t *testing.T) {
	connection := setUp()

	expectedTubes := []string{"test-tube-1", "test-tube-2"}
	for _, tubeName := range expectedTubes {
		putInTube(connection, tubeName, "message 1")
		putInTube(connection, tubeName, "message 2")
	}

	actualTubes := beanstalkd_client.NewBeanstalkdClient("127.0.0.1:11300").GetQueues()
	sort.Strings(actualTubes)
	if !reflect.DeepEqual(expectedTubes, actualTubes) {
		t.Errorf("Expected %s, got %s", expectedTubes, actualTubes)
	}
}

func TestBeanstalkdClient_GetFromTube_NoTubeExists_ShouldReturnError(t *testing.T) {
	client := beanstalkd_client.NewBeanstalkdClient("127.0.0.1:11300")
	setUp()

	_, err := client.GetJob("random-tube")
	if err == nil {
		t.Error("Expected to receive an error")
	}
}

func TestBeanstalkdClient_GetFromTube_TubeIsEmpty_ShouldReturnError(t *testing.T) {
	connection := setUp()

	emptyTubes(connection)
	putInTube(connection, "random-tube", "message 1")

	client := beanstalkd_client.NewBeanstalkdClient("127.0.0.1:11300")

	_, _ = client.GetJob("random-tube")
	_, err := client.GetJob("random-tube")

	if err == nil {
		t.Error("Expected to receive an error")
	}
}

func TestBeanstalkdClient_GetReadyJobCount_EmptyTube_ShouldReturnZero(t *testing.T) {
	connection := setUp()

	emptyTubes(connection)
	client := beanstalkd_client.NewBeanstalkdClient("127.0.0.1:11300")

	jobCount := client.GetReadyJobCount("random-tube")
	if jobCount != 0 {
		t.Errorf("Expected job count to equal 0, got %d", jobCount)
	}
}

func TestBeanstalkdClient_GetReadyJobCount_TubeHasJobs_ShouldReturnCount(t *testing.T) {
	connection := setUp()

	emptyTubes(connection)
	putInTube(connection, "random-tube", "message 1")
	putInTube(connection, "random-tube", "message 2")

	client := beanstalkd_client.NewBeanstalkdClient("127.0.0.1:11300")

	jobCount := client.GetReadyJobCount("random-tube")
	if jobCount != 2 {
		t.Errorf("Expected job count to equal 2, got %d", jobCount)
	}
}
