package apex_ecslogs

import (
	"encoding/json"
	"io"

	apex "github.com/apex/log"
	ecslogs "github.com/segmentio/ecs-logs-go"
)

type Config struct {
	Output      io.Writer
	Depth       int
	FuncInfo    func(uintptr) (ecslogs.FuncInfo, bool)
	MaxFieldLen int
}

func NewHandler(w io.Writer) apex.Handler {
	return NewHandlerWith(Config{Output: w})
}

func NewHandlerWith(c Config) apex.Handler {
	logger := ecslogs.NewLogger(c.Output)

	if c.FuncInfo == nil {
		return apex.HandlerFunc(func(entry *apex.Entry) error {
			return logger.Log(MakeEvent(entry, c.MaxFieldLen))
		})
	}

	return apex.HandlerFunc(func(entry *apex.Entry) error {
		var source string

		if pc, ok := ecslogs.GuessCaller(c.Depth, 10, "github.com/segmentio/ecs-logs", "github.com/apex/log"); ok {
			if info, ok := c.FuncInfo(pc); ok {
				source = info.String()
			}
		}

		return logger.Log(makeEvent(entry, source, c.MaxFieldLen))
	})
}

func MakeEvent(entry *apex.Entry, maxFieldLen int) ecslogs.Event {
	return makeEvent(entry, "", maxFieldLen)
}

func makeEvent(entry *apex.Entry, source string, maxFieldLen int) ecslogs.Event {
	var message string

	if maxFieldLen > 0 && len(entry.Message) > maxFieldLen {
		message = entry.Message[:maxFieldLen]
	} else {
		message = entry.Message
	}

	return ecslogs.Event{
		Level:   makeLevel(entry.Level),
		Info:    makeEventInfo(entry, source),
		Data:    makeEventData(entry, maxFieldLen),
		Time:    entry.Timestamp,
		Message: message,
	}
}

func makeEventInfo(entry *apex.Entry, source string) ecslogs.EventInfo {
	return ecslogs.EventInfo{
		Source: source,
		Errors: makeErrors(entry.Fields),
	}
}

func makeEventData(entry *apex.Entry, maxFieldLen int) ecslogs.EventData {
	data := make(ecslogs.EventData, len(entry.Fields))

	if maxFieldLen > 0 {
		for k, v := range entry.Fields {
			switch obj := v.(type) {
			case string:
				if len(obj) > maxFieldLen {
					data[k] = obj[:maxFieldLen]
				} else {
					data[k] = obj
				}
			case bool, int8, uint8, int16, uint16, int32, uint32, int64, uint64, int, uint, uintptr,
				float32, float64:
				data[k] = obj
			default:
				objBytes, err := json.Marshal(obj)
				if err != nil || len(objBytes) > maxFieldLen {
					// Just drop the field entirely, don't bother figuring out how to trim it
					continue
				}
				data[k] = obj
			}
		}
	} else {
		for k, v := range entry.Fields {
			data[k] = v
		}
	}

	return data
}

func makeLevel(level apex.Level) ecslogs.Level {
	switch level {
	case apex.DebugLevel:
		return ecslogs.DEBUG

	case apex.InfoLevel:
		return ecslogs.INFO

	case apex.WarnLevel:
		return ecslogs.WARN

	case apex.ErrorLevel:
		return ecslogs.ERROR

	case apex.FatalLevel:
		return ecslogs.CRIT

	default:
		return ecslogs.NONE
	}
}

func makeErrors(fields apex.Fields) (errors []ecslogs.EventError) {
	for k, v := range fields {
		if err, ok := v.(error); ok {
			errors = append(errors, ecslogs.MakeEventError(err))
			delete(fields, k)
		}
	}
	return
}
