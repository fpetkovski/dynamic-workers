package main

import (
	"fmt"
	"github.com/beanstalkd/go-beanstalk"
	"time"
)

func main() {

	connection, _ := beanstalk.Dial("tcp", "127.0.0.1:11300")
	for i := 1; i <= 10; i++ {
		for j := 0; j < 200; j++ {
			tube := fmt.Sprintf("tube-%d", i)
			message := fmt.Sprintf("hello from %s, message %d", tube, j)
			putInTube(connection, tube, message)
		}
	}
}

func putInTube(connection *beanstalk.Conn, tubeName string, message string) {
	tube := beanstalk.Tube{
		Conn: connection,
		Name: tubeName,
	}
	_, _ = tube.Put([]byte(message), uint32(2^31), 0, 100*time.Second)
}
