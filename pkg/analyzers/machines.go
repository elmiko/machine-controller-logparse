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

type MachinesAnalyzer struct {
	reconciledMachines map[string][]time.Time
}

func (a MachinesAnalyzer) Analyze(context AnalyzerContext) string {
	for _, e := range context.LogEntries() {
		// check for reconcile message
		if strings.Contains(e.Message(), ": reconciling Machine") {
			info := strings.SplitN(e.Message(), ":", 2)
			if len(info) != 2 {
				// this means it did not parse the message into 2 parts, we should continue
				// TODO add a log message here
				continue
			}
			machineName := strings.TrimSpace(info[0])
			_, found := a.reconciledMachines[machineName]
			if !found {
				a.reconciledMachines[machineName] = []time.Time{}
			}
			a.reconciledMachines[machineName] = append(a.reconciledMachines[machineName], e.Timestamp())
		}
	}

	summary := "Machines reconciled:\n"
	for k, _ := range a.reconciledMachines {
		summary += fmt.Sprintf("%s\n", k)
	}

	return summary
}

func (a MachinesAnalyzer) Name() string {
	return "Machines Analyzer"
}

func buildMachinesAnalyzer() Analyzer {
	return MachinesAnalyzer{
		reconciledMachines: map[string][]time.Time{},
	}
}
