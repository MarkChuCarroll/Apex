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

// File: edit.go
// Author: Mark Chu-Carroll <markcc@gmail.com>
// Description: Implemenation of the edit methods.

package buf

/* Insert a char at the cursor.
 */
func (self *GapBuffer) InsertChar(c uint8) {
  self.dirty = true
  self.primInsertChar(c, !self.undoing)
}

func (self *GapBuffer) primInsertChar(c uint8, record_undo bool) {
  if record_undo {
    undo := RecordInsert(self, self.PreLength(), 1)
    self.pushUndo(undo)
  }
  self.PushPre(c)
  if c == '\n' {
    self.line++
    self.column = 0
  } else {
    self.column++
  }
}

// This may look inefficient - but since we need to walk through the string looking for
// newlines in order to keep the line and column markers correct, it's actually about 
// as good as you can really get. 
// Eventually, we can use some caching/memoization to improve it.
func (self *GapBuffer) InsertChars(cs []uint8) {
  self.dirty = true
  pos := self.PreLength()
  for i := range cs {
    self.primInsertChar(cs[i], false)
  }
  if !self.undoing {
    undo := RecordInsert(self, pos, len(cs))
    self.pushUndo(undo)
  }
}

func (self *GapBuffer) InsertString(s string) {
  self.dirty = true
  pos := self.PreLength()
  for i := range (s) {
    self.primInsertChar(s[i], false)
  }
  if !self.undoing {
    undo := RecordInsert(self, pos, len(s))
    self.pushUndo(undo)
  }
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

func (self *GapBuffer) MoveCursorTo(pos int) {
  if self.PreLength() > pos {
    dist := self.PreLength() - pos
    self.MoveCursorBy(- dist)
  } else {
    dist := pos - self.PreLength()
    self.MoveCursorBy(dist)
  }
}

func (self *GapBuffer) MoveCursorBy(dist int) {
  if dist < 0 {
    dist = -dist
    for i := int(0); i < dist; i++ {
		  self.StepCursorBackward()
    }
  } else {
    for i := int(0); i < dist; i++ {
      self.StepCursorForward()
    }
  }
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

func (self *GapBuffer) MoveToLine(linenum int) {
  self.MoveCursorTo(0)
  for self.line < linenum && self.PostLength() > 0 {
    self.StepCursorForward()
  }
}

func (self *GapBuffer) MoveToColumn(col int) {
  if col > self.column {
    dist := col - self.column
    for i := int(0); i < dist; i++ {
      if self.PeekPost() != '\n' {
        self.StepCursorForward()
      }
    }
  } else {
    dist := self.column - col
    self.MoveCursorBy(- int(dist))
  }
}

func (self *GapBuffer) Cut(dist int) (cutbuf []uint8) {
  self.dirty = true
  if dist >= 0 {
    realdist := int(dist)
    if realdist > self.PostLength() {
      realdist = self.PostLength()
    }
    cutbuf = make([]uint8, realdist)
    for i := int(0); i < realdist; i++ {
      c := self.PopPost()
      cutbuf[i] = c
    }
    if !self.undoing {
      undo := RecordDelete(self, self.PreLength(), cutbuf)
      self.pushUndo(undo)
    }
  } else {
    realdist := -dist
    if realdist > self.PreLength() {
      realdist = self.PreLength()
    }
    pos := self.PreLength() - realdist
    cutbuf = make([]uint8, realdist)
    for i := int(0); i < realdist; i++ {
      cutbuf[realdist-i-1] = self.PopPre()
    }
    if !self.undoing {
      undo := RecordDelete(self, pos, cutbuf)
      self.pushUndo(undo)
    }
  }
  return
}

func (self *GapBuffer) Copy(dist int) (copybuf []uint8) {
  pos := self.PreLength()
  if dist >= 0 {
    realdist := dist
    if realdist > self.PostLength() {
      realdist = self.PostLength()
    }
    copybuf = make([]uint8, realdist)
    for i := int(0); i < realdist; i++ {
      c := self.PopPost()
      if c == 0 {
        return
      }
      copybuf[i] = c
      self.PushPre(c)
    }
  } else {
    realdist := -dist
    if realdist > self.PreLength() {
      realdist = self.PreLength()
    }
    copybuf = make([]uint8, realdist)
    for i := int(0); i < realdist; i++ {
      c := self.PopPre()
      copybuf[realdist-i-1] = c
      self.PushPost(c)
    }
  }
  self.MoveCursorTo(pos)
  return
}

func (self *GapBuffer) Undo() ResultCode {
  self.undoing = true
  undo, s := self.undo_stack[len(self.undo_stack)-1], self.undo_stack[:len(self.undo_stack)-1]
  self.undo_stack = s
  undo.Undo()
  self.undoing = false
  return SUCCEEDED
}

func (self *GapBuffer) GetCurrentPosition() int { return len(self.prechars) }

func (self *GapBuffer) GetCurrentLine() int { return self.line }

func (self *GapBuffer) GetCurrentColumn() int { return self.column }

func RecordInsert(b *GapBuffer, start int, length int) (result *InsertOperation) {
  result = &InsertOperation{b, start, length}
  return
}

func RecordDelete(b *GapBuffer, pos int, chars []uint8) (result *DeleteOperation) {
  result = &DeleteOperation{b, pos, chars}
  return
}

func (self *InsertOperation) GetBuffer() EditBuffer {
  return self.buf
}

func (self *InsertOperation) Undo() {
  self.buf.MoveCursorTo(self.start)
  _ = self.buf.Cut(int(self.length))
}

func (self *DeleteOperation) GetBuffer() EditBuffer {
  return self.buf
}

func (self *DeleteOperation) Undo() {
  self.buf.MoveCursorTo(self.position)
  self.buf.InsertChars(self.chars)
}

//
// Undo-record types
//


// A basic insert operation
type InsertOperation struct {
  buf    *GapBuffer
  start  int
  length int
}


// A cut operation.
type DeleteOperation struct {
  buf      *GapBuffer
  position int
  chars    []uint8
}

func (self *GapBuffer) pushUndo(u UndoOperation) {
  self.undo_stack = append(self.undo_stack, u)
}
