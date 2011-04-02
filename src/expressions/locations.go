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

// File: locexprs.go
// Author: Mark Chu-Carroll <markcc@gmail.com>
// Description: Location expressions.

package expressions

import (
  "buf"
  _ "fmt"
)

/////////////////////////////////////////////////
// Location Expressions
/////////////////////////////////////////////////


// An expression specifying a location as an absolute
// position in a file. If the offset is positive, it's
// specified in terms of the offset from the beginning
// of the buffer; if it's negative, then the position
// is relative to the end of the buffer.
type CharLocExpr struct {
  offset int
}

func NewCharLocExpr(loc int) *CharLocExpr { return &CharLocExpr{loc} }

func (self *CharLocExpr) Eval(selection Range) (r Range, status buf.Status) {
  var loc Location
  if self.offset >= 0 {
    loc = NewStartRelativeLocation(selection.GetBuffer(), self.offset)
  } else {
	loc = NewStartRelativeLocation(selection.GetBuffer(),
	                               self.buffer.Length() + self.offset)
  }
  status = Validate(loc)
  r = NewSimpleRange(selection.GetBuffer(), loc, loc)
  return
}

/////////////////////////////////////////////////

type LineLocExpr struct {
  line int
}

func NewLineLocExpr(line int) *LineLocExpr { return &LineLocExpr{line} }

func (self *LineLocExpr) Eval(selection Range) (r Range, status buf.Status) {
  var loc Location
  buffer := selection.GetBuffer()
  if self.line > 0 {
    pos, ok := buffer.GetPositionOfLine(self.line)
    if ok != buf.SUCCEEDED {
      loc = nil
      status = buf.NewFailure(buf.INVALID, "Line specifier was invalid")
      return
    }
    loc = NewStartRelativeLocation(buffer, pos)
  } else {
    // Getting the position of the end of the buffer should never fail.
    endline, _, _ := buffer.GetLineAndColumnOf(buffer.Length())
    line := endline + self.line
    pos, ok := buffer.GetPositionOfLine(line)
    if ok != buf.SUCCEEDED {
      loc = nil
      status = buf.NewFailure(buf.INVALID, "Line specifier was invalid")
      return
    }
    loc = NewStartRelativeLocation(buffer, pos)
  }
  status = Validate(loc)
  r = NewSimpleRange(selection.GetBuffer(), loc, loc)
  return
}

/////////////////////////////////////////////////

// An expression specifying a buffer location as a character
// offset relative to the start of a selection
type StartRelCharLocExpr struct {
  offset int
}

func NewStartRelCharLocExpr(loc int) *StartRelCharLocExpr {
  return &StartRelCharLocExpr{loc}
}

func (self *StartRelCharLocExpr) Eval(selection Range) (r Range, status buf.Status) {
  var loc Location
  loc = NewStartRelativeLocation(selection.GetBuffer(),
    selection.GetStart().GetAbsolute()+self.offset)
  status = Validate(loc)
  r = NewSimpleRange(selection.GetBuffer(), loc, loc)
  return
}

/////////////////////////////////////////////////

// An expression specifying a buffer location as a character
// offset relative to the end of a selection
type EndRelCharLocExpr struct {
  offset int
}

func NewEndRelCharLocExpr(loc int) *EndRelCharLocExpr {
  return &EndRelCharLocExpr{loc}
}

func (self *EndRelCharLocExpr) Eval(selection Range) (r Range, status buf.Status) {
  loc := NewStartRelativeLocation(selection.GetBuffer(),
    selection.GetEnd().GetAbsolute()+self.offset)
  status = Validate(loc)
  return
}

/////////////////////////////////////////////////

// An expression specifying a buffer location by an offset
// from the beginning of the selection in lines.
type StartRelLineLocExpr struct {
  line int
}

func NewStartRelLineLocExpr(line int) *StartRelLineLocExpr {
  return &StartRelLineLocExpr{line}
}

func (self *StartRelLineLocExpr) Eval(selection Range) (r Range, status buf.Status) {
  var loc Location
  buffer := selection.GetBuffer()
  base := selection.GetStart().GetAbsolute()
  base_line, _, line_ok := buffer.GetLineAndColumnOf(base)
  if line_ok != buf.SUCCEEDED {
    status = buf.NewFailure(buf.INVALID_LINE, "Line specifier was invalid")
    loc = nil
    return
  }
  line := base_line + self.line
  pos, pos_ok := buffer.GetPositionOfLine(line)
  if pos_ok != buf.SUCCEEDED {
    status = buf.NewFailure(buf.INVALID_LINE, "After offset, line specifier was invalid")
    return
  }
  loc = NewStartRelativeLocation(buffer, pos)
  status = Validate(loc)
  return
}

// An expression specifying a buffer location by an offset
// from the beginning of the selection in lines.
type EndRelLineLocExpr struct {
  line int
}

func NewEndRelLineLocExpr(line int) *EndRelLineLocExpr {
  return &EndRelLineLocExpr{line}
}

func (self *EndRelLineLocExpr) Eval(selection Range) (r Range, status buf.Status) {
  var loc Location
  buffer := selection.GetBuffer()
  base := selection.GetEnd().GetAbsolute()
  base_line, _, line_ok := buffer.GetLineAndColumnOf(base)
  if line_ok != buf.SUCCEEDED {
    status = buf.NewFailure(buf.INVALID_LINE, "Line specifier was invalid")
    r = nil
    return
  }
  line := base_line + self.line
  pos, pos_ok := buffer.GetPositionOfLine(line)
  if pos_ok != buf.SUCCEEDED {
    status = buf.NewFailure(buf.INVALID_LINE, "After offset, line specifier was invalid")
    r = nil
    return
  }
  loc = NewStartRelativeLocation(buffer, pos)
  status = Validate(loc)
  r = NewSimpleRange(selection.GetBuffer(), loc, loc)
  return
}

/////////////////////////////////////////////////

type GridLocationExpression struct {
  line   int
  column int
}

func NewGridLocationExpression(line int, column int) *GridLocationExpression {
  return &GridLocationExpression{line, column}
}

func (self *GridLocationExpression) Eval(sel Range) (r Range, status buf.Status) {
  var loc Location
  pos, ok := sel.GetBuffer().GetPositionOfLineAndColumn(self.line, self.column)
  if ok != buf.SUCCEEDED {
    r = nil
    status = buf.NewFailure(ok, "Error evaluating grid location")
  }
  loc = NewStartRelativeLocation(sel.GetBuffer(), pos)
  status = buf.NewSuccess()
  r = NewSimpleRange(sel.GetBuffer(), loc, loc)
  return
}
