package main

import "parsers"
import "testing"
import __regexp__ "regexp"

var tests = []testing.InternalTest{
	{"parsers.TestBasicCharSet", parsers.TestBasicCharSet},
	{"parsers.TestMany", parsers.TestMany},
	{"parsers.TestAlts", parsers.TestAlts},
	{"parsers.TestSimpleSequence", parsers.TestSimpleSequence},
	{"parsers.TestManySequence", parsers.TestManySequence},
	{"parsers.TestLetterBuiltin", parsers.TestLetterBuiltin},
	{"parsers.TestSpaceBuiltin", parsers.TestSpaceBuiltin},
	{"parsers.TestTokenBuiltin", parsers.TestTokenBuiltin},
	{"parsers.TestParensParser", parsers.TestParensParser},
	{"parsers.TestLispParser", parsers.TestLispParser},
	{"parsers.TestVectorToSExpr", parsers.TestVectorToSExpr},
	{"parsers.TestSexprParseAndBuild", parsers.TestSexprParseAndBuild},
}
var benchmarks = []testing.InternalBenchmark{}

func main() {
	testing.Main(__regexp__.MatchString, tests)
	testing.RunBenchmarks(__regexp__.MatchString, benchmarks)
}
