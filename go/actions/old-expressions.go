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

// File: expressions.go
// Author: Mark Chu-Carroll <markcc@gmail.com>
// Description: Implementation of expressions that operate on
//   selected ranges within a buffer. (tests in se_test.go)

package actions

import (
	"regexp"
	"os"
)

////////////////////////////////////////////////////////////
// Structured expressions: expressions that operate on
// selections.
////////////////////////////////////////////////////////////

type StructuredExpression interface {
  Eval(s Selection) (match Selection, success bool)
}

////////////////////////////////////////////////////////////

type LinePositionExpression struct {
  linenum int
  rel     PositionBase
}

func NewLinePositionExpression(base PositionBase, pos int) *LinePositionExpression {
  return &LinePositionExpression{pos, base}
}

func (self *LinePositionExpression) Eval(s Selection) (match Selection, success bool) {
  pos := 0
  ok := false
  switch self.rel {
  case REL_BUF_START:
    pos, ok = s.GetBuffer().GetPositionOfLine(self.linenum)
  case REL_SEL_START:
    // linenumber is relative. Find out what line the selection starts
    // on, and then convert it to an absolute  number.
    startline, _, sok :=
      s.GetBuffer().GetLineAndColumnOf(s.GetStart().GetStartRelativePosition())
    if sok {
      pos, ok = s.GetBuffer().GetPositionOfLine(startline + self.linenum)
    }
  case REL_SEL_END:
    endline, _, eok :=
      s.GetBuffer().GetLineAndColumnOf(s.GetEnd().GetStartRelativePosition())
    if eok {
      pos, ok = s.GetBuffer().GetPositionOfLine(endline + self.linenum)
    }
  case REL_BUF_END:
    endline, _, eok := s.GetBuffer().GetLineAndColumnOf(s.GetBuffer().Length())
    if eok {
      pos, ok = s.GetBuffer().GetPositionOfLine(endline + self.linenum)
    }
  default:
    ok = false
  }
  if ok {
    match = NewStartRelativePoint(s.GetBuffer(), pos)
    success = true
  } else {
    match = nil
    success = false
  }
  return
}

////////////////////////////////////////////////////////////

type CharPositionExpression struct {
  charnum int
  rel     PositionBase
}

func NewCharPositionExpression(base PositionBase, pos int) *CharPositionExpression {
  return &CharPositionExpression{pos, base}
}

func (self *CharPositionExpression) Eval(s Selection) (match Selection, success bool) {
  pos := 0
  switch self.rel {
  case REL_BUF_START:
    pos = self.charnum
  case REL_SEL_START:
    pos = s.GetStart().GetStartRelativePosition() + self.charnum
  case REL_SEL_END:
    pos = s.GetEnd().GetStartRelativePosition() + self.charnum
  case REL_BUF_END:
    pos = s.GetBuffer().Length() + self.charnum
  default:
    pos = -1
  }
  if pos > s.GetBuffer().Length() || pos < 0 {
    match = nil
    success = false
  } else {
    match = NewStartRelativePoint(s.GetBuffer(), pos)
    success = true
  }
  return
}

////////////////////////////////////////////////////////////

type RangeExpression struct {
  from StructuredExpression
  to   StructuredExpression
}

func NewRangeExpression(from, to StructuredExpression) StructuredExpression {
  return &RangeExpression{from, to}
}

func (self *RangeExpression) Eval(s Selection) (match Selection, success bool) {
  start, start_ok := self.from.Eval(s)
  end, end_ok := self.to.Eval(s)
  if !start_ok || !end_ok {
    match = nil
    success = false
  } else {
    if start.GetStart().GetStartRelativePosition() >
      end.GetStart().GetStartRelativePosition() {
      match = nil
      success = false
    } else {
      match = NewRange(s.GetBuffer(), start.GetStart(), end.GetStart())
      success = true
    }
  }
  return
}

////////////////////////////////////////////////////////////

type PatternExpression struct {
  pat string
  re  *regexp.Regexp
}

func NewPatternExpression(pat string) (result *PatternExpression, err os.Error) {
  re, error := regexp.Compile(pat)
  if err != nil {
    result = nil
    err = error
  } else {
    result = &PatternExpression{pat, re}
    err = nil
  }
  return
}

func (self *PatternExpression) Eval(s Selection) (result Selection, success bool) {
  slice, ok := s.GetContent()
  if !ok {
    result = nil
    success = false
  } else {
    matches := self.re.Execute(slice)
    if len(matches) == 0 {
      result = nil
      success = false
    } else {
      result = NewPatternRange(s.GetBuffer(),
        NewStartRelativePoint(s.GetBuffer(),
          s.GetStart().GetStartRelativePosition()+matches[0]),
        NewStartRelativePoint(s.GetBuffer(),
          s.GetStart().GetStartRelativePosition()+matches[1]),
        matches,
        slice)
      success = true
    }
  }
  return
}

////////////////////////////////////////////////////////////

type ChoiceExpression struct {
  one StructuredExpression
  two StructuredExpression
}

func NewChoiceExpression(one, two StructuredExpression) *ChoiceExpression {
  return &ChoiceExpression{one, two}
}

func (self *ChoiceExpression) Eval(s Selection) (match Selection, success bool) {
  one, one_ok := self.one.Eval(s)
  if one_ok {
    match = one
    success = one_ok
  } else {
    match, success = self.two.Eval(s)
  }
  return
}

////////////////////////////////////////////////////////////

type AfterExpression struct {
  before StructuredExpression
  after  StructuredExpression
}

func NewAfterExpression(b, a StructuredExpression) *AfterExpression {
  return &AfterExpression{b, a}
}

func (self *AfterExpression) Eval(s Selection) (match Selection, success bool) {
  before_sel, before_ok := self.before.Eval(s)
  if !before_ok {
    match = nil
    success = false
  } else {
    new_range := NewRange(s.GetBuffer(), before_sel.GetEnd(), s.GetEnd())
    if new_range.GetEnd().GetStartRelativePosition() <=
      new_range.GetStart().GetStartRelativePosition() {
      match = nil
      success = false
    } else {
      match, success = self.after.Eval(new_range)
    }
  }
  return
}
