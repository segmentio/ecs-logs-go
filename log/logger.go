package log_ecslogs

import (
	"bytes"
	"io"
	"log"
	"strconv"

	"github.com/segmentio/ecs-logs-go"
)

const (
	DefaultLevel = ecslogs.INFO
)

type Handler interface {
	HandleEntry(Entry) error
}

type HandlerFunc func(Entry) error

func (h HandlerFunc) HandleEntry(entry Entry) error {
	return h(entry)
}

func New(out io.Writer, prefix string, flags int) *log.Logger {
	return NewWithLevel(DefaultLevel, out, prefix, flags)
}

func NewWithLevel(level ecslogs.Level, out io.Writer, prefix string, flags int) *log.Logger {
	return log.New(NewWriter(prefix, flags, NewHandlerWithLevel(level, out)), prefix, flags)
}

func NewOutput(out io.Writer, prefix string, flags int) io.Writer {
	return NewOutputWithLevel(DefaultLevel, out, prefix, flags)
}

func NewOutputWithLevel(level ecslogs.Level, out io.Writer, prefix string, flags int) io.Writer {
	return NewWriter(prefix, flags, NewHandlerWithLevel(level, out))
}

func NewHandler(out io.Writer) Handler {
	return NewHandlerWithLevel(DefaultLevel, out)
}

func NewHandlerWithLevel(level ecslogs.Level, out io.Writer) Handler {
	logger := ecslogs.NewLogger(out)
	return HandlerFunc(func(entry Entry) error {
		return logger.Log(makeEvent(level, entry))
	})
}

func NewWriter(prefix string, flags int, handler Handler) io.Writer {
	return newLineWriter(writer(func(b []byte) (n int, err error) {
		var entry Entry

		if entry, err = ParseEntry(string(b), prefix, flags); err == nil {
			err = handler.HandleEntry(entry)
		}

		n = len(b)
		return
	}))
}

type writer func(b []byte) (int, error)

func (f writer) Write(b []byte) (int, error) {
	return f(b)
}

func newLineWriter(w io.Writer) io.Writer {
	buffer := &bytes.Buffer{}
	return writer(func(b []byte) (n int, err error) {
		if n, err = buffer.Write(b); err != nil {
			return
		}

		for {
			if i := bytes.IndexByte(buffer.Bytes(), '\n'); i < 0 {
				break
			} else if _, err = w.Write(buffer.Next(i + 1)); err != nil {
				break
			}
		}

		return
	})
}

func makeEvent(level ecslogs.Level, entry Entry) ecslogs.Event {
	return ecslogs.Event{
		Level:   level,
		Info:    makeEventInfo(entry),
		Data:    makeEventData(entry),
		Message: entry.Message,
		Time:    entry.Time,
	}
}

func makeEventInfo(entry Entry) (info ecslogs.EventInfo) {
	if len(entry.File) != 0 {
		info.Source = entry.File + ":" + strconv.Itoa(entry.Line)
	}
	return
}

func makeEventData(entry Entry) (data ecslogs.EventData) {
	data = ecslogs.EventData{}
	if len(entry.Prefix) != 0 {
		data["prefix"] = entry.Prefix
	}
	return data
}
