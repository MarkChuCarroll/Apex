package main

import "buf"
import "testing"
import __regexp__ "regexp"

var tests = []testing.InternalTest{
	{"buf.TestSingleInsert", buf.TestSingleInsert},
	{"buf.TestInserts", buf.TestInserts},
	{"buf.TestExpand", buf.TestExpand},
	{"buf.TestInsertAndMove", buf.TestInsertAndMove},
	{"buf.TestColumnTracking", buf.TestColumnTracking},
	{"buf.TestGotoPosition", buf.TestGotoPosition},
	{"buf.TestCut", buf.TestCut},
	{"buf.TestCutBackwards", buf.TestCutBackwards},
	{"buf.TestCutPastEnd", buf.TestCutPastEnd},
	{"buf.TestCutPastStart", buf.TestCutPastStart},
	{"buf.TestCopy", buf.TestCopy},
	{"buf.TestCopyBackwards", buf.TestCopyBackwards},
	{"buf.TestCopyPastEnd", buf.TestCopyPastEnd},
	{"buf.TestUndoInsert", buf.TestUndoInsert},
	{"buf.TestUndoCut", buf.TestUndoCut},
	{"buf.TestGetCharAt", buf.TestGetCharAt},
	{"buf.TestGetPositionOfLine", buf.TestGetPositionOfLine},
	{"buf.TestGetChars", buf.TestGetChars},
	{"buf.TestGetLineAndColumn", buf.TestGetLineAndColumn},
	{"buf.TestRead", buf.TestRead},
	{"buf.TestWrite", buf.TestWrite},
}
var benchmarks = []testing.InternalBenchmark{}

func main() {
	testing.Main(__regexp__.MatchString, tests)
	testing.RunBenchmarks(__regexp__.MatchString, benchmarks)
}
