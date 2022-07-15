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

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/elmiko/machine-controller-logparse/pkg/analyzers"
	"github.com/elmiko/machine-controller-logparse/pkg/logentry"
)

func print_usage() {
	fmt.Println("Usage: machine-controller-logparse FILE")
}

func main() {
	var help bool
	flag.BoolVar(&help, "h", false, "Print usage help")
	flag.BoolVar(&help, "help", false, "Print usage help")
	flag.Parse()
	if help {
		print_usage()
		os.Exit(0)
	}
	logfile := flag.Arg(0)
	if len(logfile) == 0 {
		print_usage()
		os.Exit(1)
	}

	data, err := os.ReadFile(logfile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	lines := strings.Split(string(data), "\n")
	entries := logentry.BuildLogEntries(lines)
	runner := analyzers.BuildAnalyzerRunner(entries)
	runner.RunAllAnalyzers()
	runner.PrintAllSummaries()
}
