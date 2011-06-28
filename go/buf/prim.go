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

// File: prim.go
// Author: Mark Chu-Carroll <markcc@gmail.com>
// Description: Implementation of the basic primitives that are
//   used to build and test buffer operations.

package buf

import (
	"container/vector"
	"os"
)

type GapBuffer struct {
  prechars   []uint8 
  postchars  []uint8
  line       int
  column     int
  undo_stack *vector.Vector
  undoing    bool
}

// Create a new gap buffer with a specified capacity.
func NewBuffer(size int) *GapBuffer {
  result := new(GapBuffer)
  result.prechars = make([]uint8, 0, size)
  result.postchars = make([]uint8, 0, size)
  result.line = 1
  result.column = 0
  result.undo_stack = new(vector.Vector)
  result.undoing = false
  result.dirty = false
  result.filename = ""
  return result
}

func NewFileBuffer(filename string) (buf *GapBuffer, result ResultCode) {
  stat, err := os.Stat(filename)
  if err != nil {
	buf = nil
	result = IO_ERROR
  } else {
	buf = NewBuffer(int(stat.Size) * 2)
	buf.filename = filename
	result = buf.Read()
  }
  return
}

////////////////////////////////////////////////////////////////
// First, the stateless methods. That is, the methods that
// work only on the contents of the buffer, without any other
// editor state.

func (self *GapBuffer) Clear() {
  self.MoveCursorTo(0)
  self.Cut(int(self.Length()))
}

func (self *GapBuffer) DeleteRange(start int, end int) (chars []uint8, result ResultCode) {
  if end < start {
    result = INVALID_RANGE
    chars = nil
    return
  }
  self.MoveCursorTo(start)
  chars, result = self.Cut(int(end - start))
  return
}

func (self *GapBuffer) InsertStringAt(pos int, s string) {
  self.MoveCursorTo(pos)
  self.InsertString(s)
}

func (self *GapBuffer) InsertCharsAt(pos int, cs []uint8) {
  self.MoveCursorTo(pos)
  self.InsertChars(cs)
}

////////////////////////////////////////////////////////////////
// Primitives used for building the stateful methods.

func (self *GapBuffer) PreLength() int { return int(len(self.prechars)) }

func (self *GapBuffer) PostLength() int { return int(len(self.postchars)) }

func (self *GapBuffer) PushPre(c uint8) { self.prechars = append(self.prechars, c) }

func (self *GapBuffer) PopPre() (result uint8) {
  if self.PreLength() > 0 {
    result = self.prechars[self.PreLength()-1]
    self.prechars = self.prechars[0 : self.PreLength()-1]
  } else {
    result = 0
  }
  return
}

func (self *GapBuffer) PushPost(c uint8) { self.postchars = append(self.postchars, c) }

func (self *GapBuffer) PopPost() (result uint8) {
  if self.PostLength() > 0 {
    result = self.postchars[self.PostLength()-1]
    self.postchars = self.postchars[0 : self.PostLength()-1]
  } else {
    result = 0
  }
  return
}

func (self *GapBuffer) PeekPost() (result uint8) {
  if self.PostLength() > 0 {
    result = self.postchars[self.PostLength()-1]
  } else {
    result = 0
  }
  return
}

func (self *GapBuffer) StepCursorForward() ResultCode {
  if self.PostLength() > 0 {
    c := self.PopPost()
    self.PushPre(c)
    if c == '\n' {
      self.line++
      self.column = 0
    }
  } else if self.PostLength() == 0 {
    return PAST_END
  }
  return SUCCEEDED
}

func (self *GapBuffer) StepCursorBackward() ResultCode {
  if self.PreLength() > 0 {
    c := self.PopPre()
    self.PushPost(c)
    if c == '\n' {
      self.line--
      i := 1
      for i < self.PreLength() && self.prechars[self.PreLength()-i] != '\n' {
        i++
      }
      self.column = i - 1
    } else {
      self.column--
      if self.column < 0 {
        return BEFORE_START
      }
    }
  }
  return SUCCEEDED
}

func (self *GapBuffer) countCurrentColumn() int {
  pos := self.PreLength() - 1
  column := int(0)
  // Count the number of characters going back until
  // you hit either a newline, or the beginning of a file.
  for self.prechars[pos] != '\n' && pos > 0 {
    column++
    pos--
  }
  return column
}


