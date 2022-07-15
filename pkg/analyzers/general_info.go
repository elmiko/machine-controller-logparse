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

package analyzers

import (
	"fmt"
	"time"
)

type GeneralInfoAnalyzer struct {
	numEntries     int
	firstTimestamp time.Time
	lastTimestamp  time.Time
}

func (a GeneralInfoAnalyzer) Analyze(context AnalyzerContext) string {
	a.numEntries = len(context.LogEntries())
	a.firstTimestamp = context.LogEntries()[0].Timestamp()
	a.lastTimestamp = context.LogEntries()[a.numEntries-1].Timestamp()
	summary := fmt.Sprintf("Log entries processed: %d\n", a.numEntries)
	summary += fmt.Sprintf("First timestamp: %s\n", a.firstTimestamp)
	summary += fmt.Sprintf("Last timestamp: %s\n", a.lastTimestamp)
	return summary
}

func (a GeneralInfoAnalyzer) Name() string {
	return "Basic Info Analyzer"
}

func buildGeneralInfoAnalyzer() Analyzer {
	return GeneralInfoAnalyzer{}
}
