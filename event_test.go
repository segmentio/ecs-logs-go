package ecslogs

import (
	"encoding/json"
	"fmt"
	"io"
	"syscall"
	"testing"
)

func TestEvent(t *testing.T) {
	unserializable := make(chan int) // a value that cannot be JSON-marshalled to ensure failure doesn't break everything

	tests := []struct {
		e Event
		s string
	}{
		{
			e: Eprintf(INFO, "answer = %d", 42),
			s: `{"level":"INFO","time":"0001-01-01T00:00:00Z","info":{},"data":{},"message":"answer = 42"}`,
		},
		{
			e: Eprintf(WARN, "an error was raised (%s)", syscall.Errno(2)),
			s: `{"level":"WARN","time":"0001-01-01T00:00:00Z","info":{"errors":[{"type":"syscall.Errno","error":"no such file or directory","errno":2,"origError":2}]},"data":{},"message":"an error was raised (no such file or directory)"}`,
		},
		{
			e: Eprint(ERROR, "an error was raised:", io.EOF),
			s: `{"level":"ERROR","time":"0001-01-01T00:00:00Z","info":{"errors":[{"type":"*errors.errorString","error":"EOF","origError":{}}]},"data":{},"message":"an error was raised: EOF"}`,
		},
		{
			e: Eprint(ERROR, "an error was raised:", fakeError{"fail", unserializable}), // value that cannot be JSON-marshalled
			s: fmt.Sprintf(`{"level":"ERROR","time":"0001-01-01T00:00:00Z","info":{},"data":{},"message":"an error was raised: {fail %v}"}`, unserializable),
		},
	}

	for _, test := range tests {
		if s := test.e.String(); s != test.s {
			t.Errorf("\n- expected: %s\n- found:    %s", test.s, s)
		}
	}
}

type fakeError struct {
	msg   string
	value interface{}
}

func (e *fakeError) Error() string {
	return e.msg
}

func TestRoundtrip(t *testing.T) {
	marshalled := Eprintf(ERROR, "error was raised %s", io.EOF).String()
	var unmarshalled Event
	err := json.Unmarshal([]byte(marshalled), &unmarshalled)
	if err != nil {
		t.Errorf("Couldn't unmarshal %s; %s", marshalled, err)
	}
}
