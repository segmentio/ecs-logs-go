package play_ecslogs

import (
	"io"

	"github.com/go-playground/log"
	"github.com/segmentio/ecs-logs-go"
)

type Config struct {
	Output   io.Writer
	Depth    int
	FuncInfo func(uintptr) (ecslogs.FuncInfo, bool)
}

func NewHandler(w io.Writer) log.Handler {
	return NewHandlerWith(Config{Output: w})
}

func NewHandlerWith(c Config) log.Handler {
	logger := ecslogs.NewLogger(c.Output)

	if c.FuncInfo == nil {
		return &handler{fn: func(entry log.Entry) {
			logger.Log(makeEvent(entry, ""))
		}}
	}

	return &handler{fn: func(entry log.Entry) {
		var source string

		if pc, ok := ecslogs.GuessCaller(c.Depth, 10, "github.com/segmentio/ecs-logs", "github.com/go-playground/log"); ok {
			if info, ok := c.FuncInfo(pc); ok {
				source = info.String()
			}
		}

		logger.Log(makeEvent(entry, source))
	}}
}

type handler struct {
	fn func(entry log.Entry)
}

func (h *handler) Log(entry log.Entry) {
	h.fn(entry)
}

func makeEvent(entry log.Entry, source string) ecslogs.Event {
	return ecslogs.Event{
		Level:   makeLevel(entry.Level),
		Info:    makeEventInfo(entry, source),
		Data:    makeEventData(entry),
		Time:    entry.Timestamp,
		Message: entry.Message,
	}
}

func makeEventInfo(entry log.Entry, source string) ecslogs.EventInfo {
	return ecslogs.EventInfo{
		Source: source,
		Errors: makeErrors(entry.Fields),
	}
}

func makeEventData(entry log.Entry) ecslogs.EventData {
	data := make(ecslogs.EventData, len(entry.Fields))

	for _, fld := range entry.Fields {
		// skip errors already handled in makeErrors func
		if _, ok := fld.Value.(error); ok {
			continue
		}
		data[fld.Key] = fld.Value
	}

	return data
}

func makeLevel(level log.Level) ecslogs.Level {
	switch level {
	case log.DebugLevel:
		return ecslogs.DEBUG

	case log.InfoLevel:
		return ecslogs.INFO

	case log.NoticeLevel:
		return ecslogs.NOTICE

	case log.WarnLevel:
		return ecslogs.WARN

	case log.ErrorLevel:
		return ecslogs.ERROR

	case log.AlertLevel:
		return ecslogs.ALERT

	case log.FatalLevel:
		return ecslogs.CRIT

	case log.PanicLevel:
		return ecslogs.EMERG

	default:
		return ecslogs.NONE
	}
}

func makeErrors(fields log.Fields) (errors []ecslogs.EventError) {
	for _, fld := range fields {
		if err, ok := fld.Value.(error); ok {
			errors = append(errors, ecslogs.MakeEventError(err))
		}
	}
	return
}
