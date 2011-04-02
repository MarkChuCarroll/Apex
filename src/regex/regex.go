// Copyright 2010 Mark C. Chu-Carroll
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

// File: regex.go
// Author: Mark Chu-Carroll <markcc@gmail.com>
// Description: Apex regular expressions

// Syntax:
// character sets: [chars]
// repetition: re* re+
// Alternation: re | re
// Grouping: (re)
// character literals: abc
// binding: {($x)re}

package regex

import (
//  "buf"
  "language"
)

type RegexInput interface {
  CharAt(i int) uint8
  Size() int
}

type Regex interface {
  Exec(in RegexInput, pos int) (status bool, match_len int, result *language.Binding)
}

type CharsRegexInput struct {
  bytes []uint8
}

func (in *CharsRegexInput) Size() int { return len(in.bytes) }

func (in *CharsRegexInput) CharAt(i int) (c uint8) {
  if i >= in.Size() {
    c = 0	
  } else {
    c = in.bytes[i]
  }
  return
}

// ***************************************
// ** CharSets

type CharSetRegex struct {
  chars string
}

func NewCharSetRE(chars string) *CharSetRegex {
  return &CharSetRegex{chars}
}
  
func (self *CharSetRegex) Exec(in RegexInput, pos int) (status bool, match_len int, result *language.Binding) {
  status = false
  match_len = -1
  result = nil
  for c := range self.chars {
    if self.chars[c] == in.CharAt(pos) {
	  match_len = 1
	  status = true
    }
  }
  return
}

// ***************************************
// ** CharSets

type RepetitionRegex struct {
  min int
  regex Regex
}

func (self *RepetitionRegex) Exec(in RegexInput, pos int) (status bool, match_len int, result *language.Binding) {
  status = false
  curPos := pos
  numMatches := 0
  result = nil
  var stepResult *language.Binding = nil
  stepSuccess, stepLen, stepResult := self.regex.Exec(in, curPos)
  for stepSuccess {
    numMatches++
    curPos = curPos + stepLen
    if stepResult != nil {
      stepResult.Next = result
      result = stepResult
    }
    stepSuccess, stepLen, stepResult = self.regex.Exec(in, curPos)
  }
  if numMatches < self.min {
    status = false
    match_len = -1
    result = nil
  } else {
    status = true
    match_len = curPos - pos
  }
  return
}
