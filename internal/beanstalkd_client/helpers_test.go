package beanstalkd_client_test

import (
	"github.com/beanstalkd/go-beanstalk"
)

func putInTube(connection *beanstalk.Conn, tubeName string, message string) {
	tube := beanstalk.Tube{
		Conn: connection,
		Name: tubeName,
	}
	_, _ = tube.Put([]byte(message), 0, 0, 0)
}

func emptyTube(connection *beanstalk.Conn, tubeName string) {
	tubeSet := beanstalk.NewTubeSet(connection, tubeName)

	for {
		id, _, err := tubeSet.Reserve(0)
		if err != nil {
			return
		}
		err = connection.Delete(id)
		if err != nil {
			panic(err)
		}
	}
}

func emptyTubes(connection *beanstalk.Conn) {
	tubes, _ := connection.ListTubes()

	for _, tubeName := range tubes {
		emptyTube(connection, tubeName)
	}
}
