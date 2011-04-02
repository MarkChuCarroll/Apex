// Copyright 2009 Mark C. Chu-Carroll
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// File: base.go
// Author: Mark Chu-Carroll <markcc@gmail.com>
// Description: Range expressions

package expressions

import (
  "buf"
  "os"
  "regexp"
)

/////////////////////////////////////////////////
// Range expressions
/////////////////////////////////////////////////

type RangeExpr interface {
  Eval(selection Range) (Range, buf.Status)
}

///////////////////////////////////////////////////////////

// A two-point range expression is the basic range
// expression, defined by two location expressions, for
// its start and its end.

type TwoPointRangeExpr struct {
  start RangeExpr
  end   RangeExpr
}

func NewTwoPointRangeExpr(l1 RangeExpr, l2 RangeExpr) *TwoPointRangeExpr {
  return &TwoPointRangeExpr{l1, l2}
}

func (self *TwoPointRangeExpr) Eval(sel Range) (result Range, status buf.Status) {
  status = sel.Validate()
  if status.GetResultCode() != buf.SUCCEEDED {
    result = nil
    return
  }
  start, start_ok := self.start.Eval(sel)
  if start_ok.GetResultCode() != buf.SUCCEEDED {
    status = start_ok
    result = nil
    return
  }
  end, end_ok := self.end.Eval(sel)
  if end_ok.GetResultCode() != buf.SUCCEEDED {
    status = end_ok
    result = nil
    return
  }
  buffer := sel.GetBuffer()
  result = NewSimpleRange(buffer, start.GetStart(), end.GetEnd())
  status = buf.NewSuccess()
  return
}

///////////////////////////////////////////////////////////

// A range expression which is defined by a regular expression
// match.
type PatternRangeExpr struct {
  pattern string
  re      *regexp.Regexp
}

func NewPatternRangeExpr(pat string) (result *PatternRangeExpr, err os.Error) {
  re, error := regexp.Compile(pat)
  if err != nil {
    result = nil
    err = error
    return
  }
  result = &PatternRangeExpr{pat, re}
  err = nil
  return
}

func (self *PatternRangeExpr) Eval(sel Range) (result Range, status buf.Status) {
  status = sel.Validate()
  if status.GetResultCode() != buf.SUCCEEDED {
    result = nil
    return
  }
  slice, _ := sel.GetContents()
  matches := self.re.FindIndex(slice)
  if len(matches) == 0 {
    result = nil
    status = buf.NewFailure(buf.MATCH_FAILED, "Regular expression match failed")
    return
  }
  buffer := sel.GetBuffer()
  start := sel.GetStart().GetAbsolute() + matches[0]
  end := sel.GetStart().GetAbsolute() + matches[1]
  result = NewPatternRange(sel.GetBuffer(), NewStartRelativeLocation(buffer, start),
    NewStartRelativeLocation(buffer, end), slice, matches)
  status = buf.NewSuccess()
  return
}

// An after expression specifies a range within the current input range
// *after* the occurence of another range. This is equivalent
// to taking the original range, finding the first range, and changing
// the selection so that its start point is immediately after the 
// end of the range found by the previous location. The results of this when
// used with absolute locations outside of the input range are peculiar
// at best.
type AfterRangeExpr struct {
  before RangeExpr
  after  RangeExpr
}

func NewAfterRangeExpr(before, after RangeExpr) *AfterRangeExpr {
  return &AfterRangeExpr{before, after}
}

func (self *AfterRangeExpr) Eval(sel Range) (result Range, status buf.Status) {
  buffer := sel.GetBuffer()
  r1, r1status := self.before.Eval(sel)
  if r1status.GetResultCode() != buf.SUCCEEDED {
    status = r1status
    result = nil
    return
  }
  newrange := NewSimpleRange(buffer, r1.GetEnd(), sel.GetEnd()).Normalize()
  result, status = self.after.Eval(newrange)
  return
}
