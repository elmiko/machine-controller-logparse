/*
Copyright 2022 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package logentry

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

var (
	loglevel = regexp.MustCompile("^[A-Z][0-9]{4} ")
	logtime  = regexp.MustCompile(`^[\d]{4}-[\d]{2}-[\d]{2}`)
)

// a single log entry, may contain multiple lines
type LogEntry struct {
	timestamp time.Time
	message   string
	lines     []string
}

func (l *LogEntry) AppendLine(line string) {
	l.lines = append(l.lines, line)
	// update the message as well
	l.message += line
}

func (l LogEntry) Message() string {
	return l.message
}

func (l LogEntry) Timestamp() time.Time {
	return l.timestamp
}

func newLogEntry(line string) (LogEntry, error) {
	ret := LogEntry{}
	var timestamp time.Time
	var err error

	data := strings.Split(string(line), " ")
	// date should be first or second entry
	if logtime.MatchString(data[0]) {
		timestamp, err = time.Parse(time.RFC3339, data[0])
	} else if logtime.MatchString(data[1]) {
		timestamp, err = time.Parse(time.RFC3339, data[1])
	} else {
		err = errors.New("unable to detect time stamp in log line")
	}

	if err != nil {
		return ret, err
	}

	// if we have a message lets parse it out
	parts := strings.SplitN(line, "]", 2)
	if len(parts) == 2 {
		ret.message = parts[1]
	}
	ret.timestamp = timestamp
	ret.lines = []string{line}

	return ret, nil
}

func BuildLogEntries(log_lines []string) []LogEntry {
	logentries := []LogEntry{}
	var prev_entry *LogEntry

	prev_entry = nil
	for _, line := range log_lines {
		if loglevel.MatchString(line) || logtime.MatchString(line) {
			// new log line if it starts with the log level or time
			entry, err := newLogEntry(line)
			if err != nil {
				continue
			}
			prev_entry = &entry
			logentries = append(logentries, *prev_entry)
		} else {
			// possible continuation of log line
			if prev_entry != nil {
				prev_entry.AppendLine(line)
			}
		}
	}

	return logentries
}
