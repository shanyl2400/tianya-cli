package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

var _log = logrus.New()

func init() {
	// Log as JSON instead of the default ASCII formatter.
	_log.SetFormatter(&logrus.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	// You could set this to any `io.Writer` such as a file
	file, err := os.OpenFile("tianya.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		_log.Out = file
	} else {
		_log.Info("Failed to log to file, using default stderr")
	}

	// Only log the warning severity or above.
	_log.SetLevel(logrus.DebugLevel)

}

func WithField(key string, value interface{}) *logrus.Entry {
	return _log.WithField(key, value)
}

func Debug(args ...interface{}) {
	_log.Debug(args...)
}

func Info(args ...interface{}) {
	_log.Info(args...)
}

func Warn(args ...interface{}) {
	_log.Warn(args...)
}

func Error(args ...interface{}) {
	_log.Error(args...)
}

func Fatal(args ...interface{}) {
	_log.Fatal(args...)
}
