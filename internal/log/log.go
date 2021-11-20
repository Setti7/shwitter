package log

import (
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()
}

func LogError(funcName string, msg string, err error) {
	Log.WithError(err).WithFields(logrus.Fields{"func": funcName}).Error(msg)
}
