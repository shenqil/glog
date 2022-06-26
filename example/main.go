package main

import (
	"github.com/shenqil/glog"
	"github.com/sirupsen/logrus"
)

func main() {
	glog.InitLog(&glog.Config{
		IsWriteToFile: false,
	})

	l := glog.Logger
	l.Infof("this is %v demo", "glog")

	lWebServer := l.WithField("component", "web-server")
	lWebServer.Info("starting...")

	lWebServerReq := lWebServer.WithFields(logrus.Fields{
		"req":   "GET /api/stats",
		"reqId": "#1",
	})

	lWebServerReq.Info("params: startYear=2048")
	lWebServerReq.Error("response: 400 Bad Request")

	lDbConnector := l.WithField("category", "db-connector")
	lDbConnector.Info("connecting to db on 10.10.10.13...")
	lDbConnector.Warn("connection took 10s")

	l.Info("demo end.")
}
