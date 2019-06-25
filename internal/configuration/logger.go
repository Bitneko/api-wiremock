package configuration

import (
	log "github.com/sirupsen/logrus"
)

func initLogger() {
	log.SetFormatter(&log.JSONFormatter{})

	// log an event as usual with logrus
	log.Info("Logger set up successfully")
}
