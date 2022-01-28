package log

import (
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()
}

func LogError(funcName string, msg string, err error) {
	// TODO get current function name automatically
	// https://stackoverflow.com/questions/25927660/how-to-get-the-current-function-name
	Log.WithError(err).WithFields(logrus.Fields{"func": funcName}).Error(msg)
}
