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

// File: insert_actions.go
// Author: Mark Chu-Carroll <markcc@gmail.com>
// Description: Edit command actions that insert text into a buffer.


package actions

import (
  "buf"
  "expressions"
  "strings"
)

type TextActionKind int

const (
  INSERT TextActionKind = iota
  APPEND
  REPLACE
)

type TextAction struct {
  data string 
  kind TextActionKind
}


func NewInsertAction(s string) *TextAction {
  return &TextAction{ s, INSERT }
}

// An action which appends a string of text to the end
// of a range. (Since the range should have been normalized 
// before the action was invoked, the inserted text will 
// still be in the range.)
func NewAppendAction(s string) *TextAction {
  return &TextAction{ s, APPEND }
}

// An edit action that replaces a range with a new string of
// text.
func NewReplaceAction(s string) *TextAction {
  return &TextAction{  s, REPLACE }
}

func (self *TextAction) Execute(r expressions.Range) (result buf.Status) {
  buffer := r.GetBuffer()
  range_valid := r.Validate()
  // If the range is invalid, then return status indicating
  // what's wrong with the range.
  if range_valid.GetResultCode() != buf.SUCCEEDED {
    return range_valid	
  }
  text, text_status := self.ResolveText(r)
  if text_status.GetResultCode() != buf.SUCCEEDED {
    return text_status	
  }

  switch self.kind {
    case INSERT:
      buffer.MoveTo(r.GetStart().GetAbsolute())
    case APPEND:
	  buffer.MoveTo(r.GetEnd().GetAbsolute())
	case REPLACE:
	  buffer.MoveTo(r.GetStart().GetAbsolute())
	  buffer.Cut(r.GetLength())
  }
  buffer.InsertChars(text)
  result = buf.NewSuccess()
  return
}

func (self *TextAction) ResolveText(r expressions.Range) (txt []uint8, status buf.Status) {
  patrange, ok := r.(*expressions.PatternRange)
  if !ok {
	// not an insert of a successful pattern match
    txt = strings.Bytes(self.data)
    status = buf.NewSuccess()
  } else {
    txt, status = patrange.InstantiateTemplate(self.data)
  }
  return
}
