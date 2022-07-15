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
	"strings"
	"time"
)

type AverageReconcileDurationAnalyzer struct{}

func (a AverageReconcileDurationAnalyzer) Analyze(context AnalyzerContext) string {
	timestamps := []time.Time{}
	for _, e := range context.LogEntries() {
		// check for reconcile message
		if strings.Contains(e.Message(), reconcilingMachine) {
			timestamps = append(timestamps, e.Timestamp())
		}
	}
	prevstamp := timestamps[0]
	diffs := []time.Duration{}
	for _, t := range timestamps[1:] {
		diff := t.Sub(prevstamp)
		diffs = append(diffs, diff)
		prevstamp = t
	}
	var diffsum time.Duration = 0
	for _, d := range diffs {
		diffsum += d
	}
	average := int64(diffsum) / int64(len(diffs))

	summary := fmt.Sprintf("Average reconcile duration: %s\n", time.Duration(average).String())
	return summary
}

func (a AverageReconcileDurationAnalyzer) Name() string {
	return "Average Reconcile Duration Analyzer"
}

func buildAverageReconcileDurationAnalyzer() Analyzer {
	return AverageReconcileDurationAnalyzer{}
}
