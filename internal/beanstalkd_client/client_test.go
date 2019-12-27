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

	actualTubes := beanstalkd_client.NewBeanstalkdClient("127.0.0.1:11300").GetTubes()
	sort.Strings(actualTubes)
	if !reflect.DeepEqual(expectedTubes, actualTubes) {
		t.Errorf("Expected %s, got %s", expectedTubes, actualTubes)
	}
}

func TestBeanstalkdClient_GetFromTube_NoTubeExists_ShouldReturnError(t *testing.T) {
	client := beanstalkd_client.NewBeanstalkdClient("127.0.0.1:11300")
	setUp()

	_, _, err := client.GetFromTube("random-tube")
	if err == nil {
		t.Error("Expected to receive an error")
	}
}

func TestBeanstalkdClient_GetFromTube_TubeIsEmpty_ShouldReturnError(t *testing.T) {
	connection := setUp()

	emptyTubes(connection)
	putInTube(connection, "random-tube", "message 1")

	client := beanstalkd_client.NewBeanstalkdClient("127.0.0.1:11300")

	_, _, _ = client.GetFromTube("random-tube")
	_, _, err := client.GetFromTube("random-tube")

	if err == nil {
		t.Error("Expected to receive an error")
	}
}
