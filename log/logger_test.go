package log_ecslogs

import (
	"log"
	"reflect"
	"testing"
	"time"
)

func TestLineWriter(t *testing.T) {
	tests := []struct {
		chunks []string
		lines  []string
	}{
		{
			chunks: []string{},
			lines:  []string{},
		},
		{
			chunks: []string{"A", "B", "C", "\n"},
			lines:  []string{"ABC\n"},
		},
		{
			chunks: []string{"A\n", "B", "C\n\n"},
			lines:  []string{"A\n", "BC\n", "\n"},
		},
	}

	for i, test := range tests {
		s := []string{}
		w := newLineWriter(writer(func(b []byte) (n int, err error) {
			s = append(s, string(b))
			n = len(b)
			return
		}))

		for _, chunk := range test.chunks {
			if n, err := w.Write([]byte(chunk)); err != nil {
				t.Errorf("test#%d: failed writing chunk: %#v: %s", i, chunk, err)
			} else if n != len(chunk) {
				t.Errorf("test#%d: invalid byte count returned: %d != %d", i, n, len(chunk))
			}
		}

		if !reflect.DeepEqual(s, test.lines) {
			t.Errorf("test#%d: %#v != %#v", i, s, test.lines)
		}
	}
}

func TestNewWriter(t *testing.T) {
	tests := []struct {
		entries []Entry
		content string
		prefix  string
		flags   int
	}{
		{
			entries: []Entry{},
			content: "",
			prefix:  "",
			flags:   0,
		},
		{
			entries: []Entry{
				Entry{
					Prefix:  "[12345] ",
					Message: "Hello World!",
					File:    "logger_test.go",
					Line:    21,
					Time:    time.Date(2016, 7, 7, 12, 6, 25, 0, time.UTC),
				},
				Entry{
					Prefix:  "[12345] ",
					Message: "How are you?",
					File:    "logger_test.go",
					Line:    42,
					Time:    time.Date(2016, 7, 7, 12, 6, 26, 0, time.UTC),
				},
				Entry{
					Prefix:  "[12345] ",
					Message: "Another log message...",
					File:    "logger_test.go",
					Line:    84,
					Time:    time.Date(2016, 7, 7, 12, 6, 27, 0, time.UTC),
				},
				Entry{
					Prefix:  "[12345] ",
					Message: "Hello World!",
					File:    "logger_test.go",
					Line:    168,
					Time:    time.Date(2016, 7, 7, 12, 6, 28, 0, time.UTC),
				},
				Entry{
					Prefix:  "[12345] ",
					Message: "Hello World!",
					File:    "logger_test.go",
					Line:    336,
					Time:    time.Date(2016, 7, 7, 12, 6, 29, 0, time.UTC),
				},
			},
			content: `[12345] 2016/07/07 12:06:25 logger_test.go:21: Hello World!
[12345] 2016/07/07 12:06:26 logger_test.go:42: How are you?
[12345] 2016/07/07 12:06:27 logger_test.go:84: Another log message...
[12345] 2016/07/07 12:06:28 logger_test.go:168: Hello World!
[12345] 2016/07/07 12:06:29 logger_test.go:336: Hello World!
`,
			prefix: "[12345] ",
			flags:  log.LstdFlags | log.Lshortfile | log.LUTC,
		},
	}

	for i, test := range tests {
		entries := []Entry{}
		writer := NewWriter(test.prefix, test.flags, HandlerFunc(func(entry Entry) (err error) {
			entries = append(entries, entry)
			return
		}))

		if n, err := writer.Write([]byte(test.content)); err != nil {
			t.Errorf("test#%d: failed writing content: %#v: %s", i, test.content, err)
		} else if n != len(test.content) {
			t.Errorf("test#%d: invalid byte count returned: %d != %d", i, n, len(test.content))
		} else if !reflect.DeepEqual(entries, test.entries) {
			t.Errorf("test#%d: the writer produced invalid entries:\n%#v\n%#v", entries, test.entries)
		}
	}
}
