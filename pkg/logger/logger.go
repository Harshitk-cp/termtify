package logger

import (
	"github.com/sirupsen/logrus"
)

func Setup() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)
}
