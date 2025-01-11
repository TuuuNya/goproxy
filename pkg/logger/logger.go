package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	log.SetLevel(logrus.DebugLevel)

	log.SetOutput(os.Stdout)
}

func SetLogLevel(debug bool) {
	if debug {
		log.SetLevel(logrus.DebugLevel)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}
}

func GetLogger() *logrus.Logger {
	return log
}
