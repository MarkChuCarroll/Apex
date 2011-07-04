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

// File: undo.go
// Author: Mark Chu-Carroll <markcc@gmail.com>
// Description: undo objects and methods.

package buf

func (self *GapBuffer) Undo() ResultCode {
  self.undoing = true
  undo := self.undo_stack.Pop().(UndoOperation)
  code := undo.Undo()
  self.undoing = false
  return code
}

func (self *GapBuffer) pushUndo(u UndoOperation) {
  self.undo_stack.Push(u)
}


// A basic insert operation
type InsertOperation struct {
  buf    *GapBuffer
  start  int
  length int
}

func RecordInsert(b *GapBuffer, start int, length int) (result *InsertOperation) {
  result = &InsertOperation{b, start, length}
  return
}


func (self *InsertOperation) GetBuffer() *GapBuffer {
  return self.buf
}

func (self *InsertOperation) Undo() ResultCode {
  self.buf.MoveCursorTo(self.start)
  _, status := self.buf.Cut(int(self.length))
  return status
}


// A cut operation.
type DeleteOperation struct {
  buf      *GapBuffer
  position int
  chars    []uint8
}

func RecordDelete(b *GapBuffer, pos int, chars []uint8) (result *DeleteOperation) {
  result = &DeleteOperation{b, pos, chars}
  return
}

func (self *DeleteOperation) GetBuffer() *GapBuffer {
  return self.buf
}

func (self *DeleteOperation) Undo() (result ResultCode) {
  result = self.buf.MoveCursorTo(self.position)
  if result != SUCCEEDED {
    return result
  }
  self.buf.InsertChars(self.chars)
  result = SUCCEEDED
  return
}


