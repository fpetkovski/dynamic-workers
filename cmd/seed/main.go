package main

import "github.com/beanstalkd/go-beanstalk"

func main() {

	connection, _ := beanstalk.Dial("tcp", "127.0.0.1:11300")
	putInTube(connection, "tube-1", "hello from tube-1 1")
	putInTube(connection, "tube-1", "hello from tube-1 2")
	putInTube(connection, "tube-1", "hello from tube-1 3")
	putInTube(connection, "tube-1", "hello from tube-1 4")
	putInTube(connection, "tube-1", "hello from tube-1 5")

	//putInTube(connection, "tube-2", "hello from tube-2 1")
	//putInTube(connection, "tube-2", "hello from tube-2 2")
	//putInTube(connection, "tube-2", "hello from tube-2 3")
	//putInTube(connection, "tube-2", "hello from tube-2 4")
	//putInTube(connection, "tube-2", "hello from tube-2 5")
}

func putInTube(connection *beanstalk.Conn, tubeName string, message string) {
	tube := beanstalk.Tube{
		Conn: connection,
		Name: tubeName,
	}
	_, _ = tube.Put([]byte(message), uint32(2^31), 0, 0)
}