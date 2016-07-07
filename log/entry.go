package log_ecslogs

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type Entry struct {
	Prefix  string
	Message string
	File    string
	Line    int
	Time    time.Time
}

func ParseEntryLogger(s string, logger *log.Logger) (entry Entry, err error) {
	var prefix string
	var flags int

	if logger != nil {
		prefix, flags = logger.Prefix(), logger.Flags()
	}

	return ParseEntry(s, prefix, flags)
}

func ParseEntry(s string, prefix string, flags int) (entry Entry, err error) {
	if !strings.HasPrefix(s, prefix) {
		err = fmt.Errorf("missing prefix in log line: %#v", prefix)
	}

	entry.Prefix = strings.TrimSpace(prefix)
	s = s[len(prefix):]

	if format := TimeFormat(flags); len(format) != 0 {
		var tz *time.Location
		var ts string

		if n := len(format); len(s) >= n {
			ts, s = s[:n], s[n:]
		} else {
			ts, s = s, ""
		}

		if (flags & log.LUTC) != 0 {
			tz = time.UTC
		} else {
			tz = time.Local
		}

		entry.Time, err = time.ParseInLocation(format, ts, tz)
	}

	if (flags&log.Lshortfile) != 0 || (flags&log.Llongfile) != 0 {
		var file string
		var line string

		s = skip(s, ' ')
		file, s = parse(s, ':') // parse the file name or path
		line, s = parse(s, ':') // parse the line number

		entry.File = file
		entry.Line, err = strconv.Atoi(line)
	}

	entry.Message = skip(s, ' ')

	if n := len(entry.Message); n != 0 && entry.Message[n-1] == '\n' {
		entry.Message = entry.Message[:n-1]
	}

	return
}

func TimeFormat(flags int) string {
	var format string

	if (flags & log.Ldate) != 0 {
		format = "2006/01/02 "
	}

	if (flags & log.Lmicroseconds) != 0 {
		format += "15:04:05.999999"

	} else if (flags & log.Ltime) != 0 {
		format += "15:04:05"

	} else if len(format) != 0 {
		format = format[:len(format)-1]
	}

	return format
}

func skip(s string, b byte) string {
	if len(s) != 0 && s[0] == b {
		s = s[1:]
	}
	return s
}

func parse(s string, b byte) (left string, right string) {
	if index := strings.IndexByte(s, b); index >= 0 {
		left, right = s[:index], s[index+1:]
	} else {
		left = s
	}
	return
}
