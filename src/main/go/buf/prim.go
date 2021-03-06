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
	"os"
)

type GapBuffer struct {
  prechars   []uint8 
  postchars  []uint8
  line       int
  column     int
  undo_stack []UndoOperation
  undoing    bool
  dirty      bool
  filename   string	
}

// Create a new gap buffer with a specified capacity.
func NewBuffer(size int) *GapBuffer {
  result := new(GapBuffer)
  result.prechars = make([]uint8, 0, size)
  result.postchars = make([]uint8, 0, size)
  result.line = 1
  result.column = 0
  result.undo_stack = make([]UndoOperation, 0, 1000)
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
	buf = NewBuffer(int(stat.Size()) * 2)
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
