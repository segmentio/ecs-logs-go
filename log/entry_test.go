package log_ecslogs

import (
	"bytes"
	"log"
	"testing"
	"time"
)

func TestTimeFormat(t *testing.T) {
	tests := []struct {
		flags  int
		format string
	}{
		{
			flags:  0,
			format: "",
		},
		{
			flags:  log.Ldate,
			format: "2006/01/02",
		},
		{
			flags:  log.Ltime,
			format: "15:04:05",
		},
		{
			flags:  log.Lmicroseconds,
			format: "15:04:05.999999",
		},
		{
			flags:  log.LstdFlags,
			format: "2006/01/02 15:04:05",
		},
		{
			flags:  log.Ldate | log.Ltime,
			format: "2006/01/02 15:04:05",
		},
		{
			flags:  log.Ltime | log.Lmicroseconds,
			format: "15:04:05.999999",
		},
		{
			flags:  log.Ldate | log.Lmicroseconds,
			format: "2006/01/02 15:04:05.999999",
		},
	}

	for _, test := range tests {
		if s := TimeFormat(test.flags); s != test.format {
			t.Errorf("invalid time format: %#v != %#v", s, test.format)
		}
	}
}

func TestParseEntry(t *testing.T) {
	tests := []struct {
		s string
		p string
		f int
		e Entry
	}{
		{
			s: "\n",
			p: "",
			f: 0,
			e: Entry{},
		},
		{
			s: "[prefix] \n",
			p: "[prefix] ",
			f: 0,
			e: Entry{
				Prefix: "[prefix] ",
			},
		},
		{
			s: "[prefix] 2016/07/07 12:06:25\n",
			p: "[prefix] ",
			f: log.LstdFlags,
			e: Entry{
				Prefix: "[prefix] ",
				Time:   time.Date(2016, 7, 7, 12, 6, 25, 0, time.Local),
			},
		},
		{
			s: "2016/07/07 12:06:25\n",
			p: "",
			f: log.LstdFlags,
			e: Entry{
				Time: time.Date(2016, 7, 7, 12, 6, 25, 0, time.Local),
			},
		},
		{
			s: "2016/07/07 12:06:25.745609\n",
			p: "",
			f: log.Ldate | log.Lmicroseconds,
			e: Entry{
				Time: time.Date(2016, 7, 7, 12, 6, 25, 745609000, time.Local),
			},
		},
		{
			s: "Hello World!\n",
			p: "",
			f: 0,
			e: Entry{
				Message: "Hello World!",
			},
		},
		{
			s: "[prefix] Hello World!\n",
			p: "[prefix] ",
			f: 0,
			e: Entry{
				Prefix:  "[prefix] ",
				Message: "Hello World!",
			},
		},
		{
			s: "[prefix] 2016/07/07 12:06:25 Hello World!\n",
			p: "[prefix] ",
			f: log.LstdFlags | log.LUTC,
			e: Entry{
				Prefix:  "[prefix] ",
				Message: "Hello World!",
				Time:    time.Date(2016, 7, 7, 12, 6, 25, 0, time.UTC),
			},
		},
		{
			s: "[prefix] 2016/07/07 12:06:25 entry_test.go:88: Hello World!\n",
			p: "[prefix] ",
			f: log.LstdFlags | log.Lshortfile,
			e: Entry{
				Prefix:  "[prefix] ",
				Message: "Hello World!",
				File:    "entry_test.go",
				Line:    88,
				Time:    time.Date(2016, 7, 7, 12, 6, 25, 0, time.Local),
			},
		},
		{
			s: "[prefix] 2016/07/07 12:06:25 /home/local/dev/src/github.com/segmentio/ecs-logs-go/log/entry_test.go:88: Hello World!\n",
			p: "[prefix] ",
			f: log.LstdFlags | log.Llongfile,
			e: Entry{
				Prefix:  "[prefix] ",
				Message: "Hello World!",
				File:    "/home/local/dev/src/github.com/segmentio/ecs-logs-go/log/entry_test.go",
				Line:    88,
				Time:    time.Date(2016, 7, 7, 12, 6, 25, 0, time.Local),
			},
		},
	}

	for _, test := range tests {
		if entry, err := ParseEntry(test.s, test.p, test.f); err != nil {
			t.Error("error parsing entry:", err)
		} else if entry != test.e {
			t.Errorf("invalid parsed entry: %#v != %#v", entry, test.e)
		}
	}
}

func TestParseEntryLogger(t *testing.T) {
	tests := []struct {
		p string
		m string
		f int
	}{
		{
			p: "",
			m: "",
			f: 0,
		},
	}

	for _, test := range tests {
		buffer := &bytes.Buffer{}
		logger := log.New(buffer, test.p, test.f)
		logger.Println(test.m)

		if entry, err := ParseEntryLogger(buffer.String(), logger); err != nil {
			t.Error("error parsing entry:", err)
		} else {
			if entry.Prefix != test.p {
				t.Errorf("invalid entry prefix: %#v != %#v", entry.Prefix, test.p)
			}
			if entry.Message != test.m {
				t.Errorf("invalid entry message: %#v != %#v", entry.Message, test.m)
			}
		}
	}
}
