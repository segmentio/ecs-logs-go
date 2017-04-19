package logrus_ecslogs

import (
	"github.com/Sirupsen/logrus"
	"github.com/segmentio/ecs-logs-go/log"
)

type Forwarder struct {
	Level  logrus.Level
	Logger logrus.FieldLogger
}

func (f *Forwarder) HandleEntry(entry log_ecslogs.Entry) error {
	logger := f.Logger
	if logger == nil {
		logger = logrus.StandardLogger()
	}

	msg := entry.Message

	switch f.Level {
	case logrus.DebugLevel:
		logger.Debug(msg)
	case logrus.InfoLevel:
		logger.Info(msg)
	case logrus.WarnLevel:
		logger.Warn(msg)
	case logrus.ErrorLevel:
		logger.Error(msg)
	case logrus.FatalLevel:
		logger.Fatal(msg)
	case logrus.PanicLevel:
		logger.Panic(msg)
	}

	return nil
}
