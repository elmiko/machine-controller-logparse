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

	"github.com/elmiko/machine-controller-logparse/pkg/logentry"
)

// This is the list of functions which create all the analyzers that a runner
// will use. When creating a new analyzer, it must have a builder function
// `func () Analyzer` which is present in this list for it to be utilized.
// The ordering of this list is also important as each Analyzer can modify
// the AnalyzerContext for subsequent Analyzers to use.
var analyzerBuilders []func() Analyzer = []func() Analyzer{
	buildGeneralInfoAnalyzer,
	buildMachinesAnalyzer,
}

// An interface for individual analyzers
type Analyzer interface {
	Analyze(AnalyzerContext) string
	Name() string
}

// A context to hold information that can be shared between analyzers
// the idea here is that this will contain the source loglines and then
// have some arbitrary information (eg Machines detected) that can be
// passed between Analyzers. In this manner, Analyzers can use the work
// of their previous siblings to enhance their activity.
type AnalyzerContext struct {
	logentries []logentry.LogEntry
}

func (ac AnalyzerContext) LogEntries() []logentry.LogEntry {
	return ac.logentries
}

type AnalyzerRunner struct {
	context   AnalyzerContext
	analyzers []Analyzer
	results   map[string]string
}

func (r AnalyzerRunner) RunAllAnalyzers() {
	for _, analyzer := range r.analyzers {
		r.results[analyzer.Name()] = analyzer.Analyze(r.context)
	}
}

func (r AnalyzerRunner) PrintAllSummaries() {
	for k, v := range r.results {
		fmt.Printf("Results for %s\n", k)
		fmt.Println("====================")
		fmt.Printf("%s\n", v)
	}
}

func BuildAnalyzerRunner(logentries []logentry.LogEntry) AnalyzerRunner {
	ctx := AnalyzerContext{
		logentries: logentries,
	}

	runner := AnalyzerRunner{
		context: ctx,
		results: map[string]string{},
	}
	for _, builder := range analyzerBuilders {
		analyzer := builder()
		runner.analyzers = append(runner.analyzers, analyzer)
	}

	return runner
}
