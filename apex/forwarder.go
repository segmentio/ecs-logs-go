package apex_ecslogs

import (
	"github.com/apex/log"
	"github.com/segmentio/ecs-logs-go/log"
)

type Forwarder struct {
	Level  log.Level
	Logger log.Interface
}

func (f *Forwarder) HandleEntry(entry log_ecslogs.Entry) error {
	logger := f.Logger
	if logger == nil {
		logger = log.Log
	}

	msg := entry.Message

	switch f.Level {
	case log.DebugLevel:
		logger.Debug(msg)
	case log.InfoLevel:
		logger.Info(msg)
	case log.WarnLevel:
		logger.Warn(msg)
	case log.ErrorLevel:
		logger.Error(msg)
	case log.FatalLevel:
		logger.Fatal(msg)
	}

	return nil
}
