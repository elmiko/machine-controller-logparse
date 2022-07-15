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
	"testing"
	"time"
)

func TestBuildLogEntries(t *testing.T) {
	testcases := []struct {
		name           string
		logLines       []string
		expectedLength int
	}{
		{
			name:           "Single entry",
			logLines:       []string{"I0101 2000-01-01T00:11:22.0000Z fake line"},
			expectedLength: 1,
		},
		{
			name:           "Single entry, date first, with extra line",
			logLines:       []string{"2000-01-01T00:11:22.0000Z I0101 fake line", "followup line"},
			expectedLength: 1,
		},
		{
			name:           "Single entry, date first",
			logLines:       []string{"2000-01-01T00:11:22.0000Z I0101 fake line"},
			expectedLength: 1,
		},
		{
			name: "Multiple entry, date first, with one extra line",
			logLines: []string{
				"2000-01-01T00:11:22.0000Z I0101 fake line",
				"extra line",
				"2000-01-01T00:11:22.0000Z I0101 fake line",
			},
			expectedLength: 2,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			observed := BuildLogEntries(tc.logLines)
			if len(observed) != tc.expectedLength {
				t.Errorf("Observed length %d, expected %d", len(observed), tc.expectedLength)
			}
		})
	}
}

func TestNewLogEntry(t *testing.T) {
	expectedTime, err := time.Parse(time.RFC3339, "2022-01-13T12:35:00.844253372Z")
	if err != nil {
		t.Error("Unable to construct timestamp")
	}
	testcases := []struct {
		name              string
		logline           string
		extralogline      string
		expectedTimestamp time.Time
		expectedMessage   string
		expectedError     bool
	}{
		{
			name:              "Well formatted entry",
			logline:           "2022-01-13T12:35:00.844253372Z I0113 12:35:00.844081       1 main.go:106] Watching machine-api objects only in namespace \"openshift-machine-api\" for reconciliation.",
			expectedTimestamp: expectedTime,
			expectedMessage:   " Watching machine-api objects only in namespace \"openshift-machine-api\" for reconciliation.",
			expectedError:     false,
		},
		{
			name:              "Well formatted entry, with extra log line",
			logline:           "2022-01-13T12:35:00.844253372Z I0113 12:35:00.844081       1 main.go:106] Watching machine-api objects only in namespace \"openshift-machine-api\" for reconciliation.",
			extralogline:      " A little more.",
			expectedTimestamp: expectedTime,
			expectedMessage:   " Watching machine-api objects only in namespace \"openshift-machine-api\" for reconciliation. A little more.",
			expectedError:     false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			observed, err := newLogEntry(tc.logline)
			if err != nil && !tc.expectedError {
				t.Errorf("Got error when not expected %v", err)
			}

			if len(tc.extralogline) > 0 {
				observed.AppendLine(tc.extralogline)
			}

			if observed.timestamp != expectedTime {
				t.Errorf("Observed timestamp %v, expected %v", observed.timestamp, expectedTime)
			}
			if observed.message != tc.expectedMessage {
				t.Errorf("Observed message %s, expected %s", observed.message, tc.expectedMessage)
			}
		})
	}

}
