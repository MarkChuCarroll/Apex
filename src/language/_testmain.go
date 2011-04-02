package main

import "language"
import "testing"
import __regexp__ "regexp"

var tests = []testing.InternalTest{
	{"language.TestScannerInput", language.TestScannerInput},
	{"language.TestScanner", language.TestScanner},
}
var benchmarks = []testing.InternalBenchmark{}

func main() {
	testing.Main(__regexp__.MatchString, tests)
	testing.RunBenchmarks(__regexp__.MatchString, benchmarks)
}
