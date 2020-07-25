package apex_ecslogs

import (
	"bytes"
	"io"
	"strings"
	"testing"

	apex "github.com/apex/log"
	ecslogs "github.com/segmentio/ecs-logs-go"
)

func TestHandler(t *testing.T) {
	buf := &bytes.Buffer{}
	log := &apex.Logger{
		Handler: NewHandlerWith(Config{
			Output:   buf,
			FuncInfo: testFuncInfo,
		}),
		Level: apex.DebugLevel,
	}

	log.
		WithField("error", io.EOF).
		WithField("hello", "world").
		Errorf("an error was raised (%s)", io.EOF)

	s := strings.TrimSpace(buf.String())

	// I wish we could make better testing here but the apex
	// API doesn't let us mock the timestamp so we can't really
	// predict what "time" is gonna be.
	if !strings.HasPrefix(s, `{"level":"ERROR","time":"`) || !strings.HasSuffix(s, `"errors":[{"type":"*errors.errorString","error":"EOF","origError":{}}]},"data":{"hello":"world"},"message":"an error was raised (EOF)"}`) {
		t.Error("apex handler failed:", s)
	}
}

func TestHandlerMaxFieldLength(t *testing.T) {
	buf := &bytes.Buffer{}
	log := &apex.Logger{
		Handler: NewHandlerWith(Config{
			Output:      buf,
			FuncInfo:    testFuncInfo,
			MaxFieldLen: 10,
		}),
		Level: apex.DebugLevel,
	}

	log.
		WithField("hello", 1234).
		WithField("key", "01234567890123456789").
		Info("abcdefghijklmnopqrstuvwxyz")

	s := strings.TrimSpace(buf.String())

	// I wish we could make better testing here but the apex
	// API doesn't let us mock the timestamp so we can't really
	// predict what "time" is gonna be.
	if !strings.HasPrefix(s, `{"level":"INFO","time":"`) || !strings.HasSuffix(s, `"info":{"source":"bytes/buffer.go:42:bytes.(*Buffer).String"},"data":{"hello":1234,"key":"0123456789"},"message":"abcdefghij"}`) {
		t.Error("apex handler failed:", s)
	}
}

func testFuncInfo(pc uintptr) (info ecslogs.FuncInfo, ok bool) {
	if info, ok = ecslogs.GetFuncInfo(pc); !ok {
		return
	}
	info.Line = 42
	return
}
