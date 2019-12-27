package main

import (
	"fpetkovski/worker_pool/internal/beanstalkd_client"
	"fpetkovski/worker_pool/internal/controller"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	beanstalkClient := beanstalkd_client.NewBeanstalkdClient("127.0.0.1:11300")
	ctrl := controller.NewController(beanstalkClient)

	ctrl.Loop()
}
