// Copyright 2011 Mark C. Chu-Carroll
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
// Description: Implementation of the basic primitive implementation of
//   a gap buffer.

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
  dirty		 bool
  filename	 string
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

// TODO: this should be replaced by FileManager.OpenFile
func NewFileBuffer(filename string) (buf *GapBuffer, result ResultCode) {
  stat, err := os.Stat(filename)
  if err != nil {
	buf = nil
	result = IO_ERROR
  } else {
	buf = NewBuffer(int(stat.Size) * 2)
	buf.filename = filename
	result, _ = buf.Read()
  }
  return
}

/////////////////////////////////////////////////////////////////
// First, the cursorless primitives. These are the methods that
// work only on the contents of the buffer, without relying on
// any editor state.

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

func (self *GapBuffer) GetCharAt(pos int) (c uint8, success ResultCode) {
  if pos > self.PreLength()+self.PostLength() {
    c = 0
    success = PAST_END
    return
  } else {
    success = SUCCEEDED
    if pos < self.PreLength() {
      c = self.prechars[pos]
    } else {
      offset := pos - self.PreLength()
      c = self.postchars[self.PostLength()-offset-1]
    }
  }
  return
}

func (self *GapBuffer) GetPositionOfLine(linenum int) (pos int, success ResultCode) {
  linepos := 1
  charpos := 0
  for charpos < self.Length() && linepos < linenum {
    charpos++
    c, _ := self.GetCharAt(charpos) // we know that charpos is in the buffer.
    if c == '\n' {
      linepos++
      charpos++ // skip past the newline The pointer should be placed
      // *after* the newline character, so the charpos needs to be one greater -
      // if the newline is the last character in the buffer, then this will
      // correctly produce an error below because charpos >= self.Length()
    }
  }
  if charpos >= self.Length() {
    // The line wasn't found
    pos = 0
    success = PAST_END
  } else {
    pos = charpos
    success = SUCCEEDED
  }
  return
}

func (self *GapBuffer) GetPositionOfLineAndColumn(linenum int, colnum int) (pos int, status ResultCode) {
  lpos, l_ok := self.GetPositionOfLine(linenum)
  if l_ok != SUCCEEDED {
    pos = 0
    status = INVALID_LINE
    return
  }
  pos = lpos
  status = SUCCEEDED
  for i := 0; i < colnum; i++ {
    if self.PeekPost() != '\n' {
      pos++
    } else {
      status = INVALID_COLUMN
      return
    }
  }
  return
}

func (self *GapBuffer) GetRange(start int, end int) (chars []uint8, success ResultCode) {
  if start >= self.Length() || end > self.Length() {
    chars = nil
    success = PAST_END
    return
  } else if start < 0 || end < 0 {
    chars = nil
    success = BEFORE_START
  } else {
    chars = make([]uint8, end-start)
    l := end - start
    for i := 0; i < l; i++ {
      chars[i], _ = self.GetCharAt(start + i)
    }
    success = SUCCEEDED
  }
  return
}

// This is a naive implementation, which is likely to be fairly slow. But slow
// is likely to be fast enough. If it turns out to be a problem, we can easily
// add some caching to the GapBuffer implementation to optimize it (ie, set up
// a table of line/column locations, so that we have a starting point every 500
// characters.
func (self *GapBuffer) GetCoordinates(pos int) (line int, col int, success ResultCode) {
  if pos > self.Length() {
    col = 0
    success = PAST_END
  } else {
    col = 0
    line = 1
    for i := 0; i < pos; i++ {
      c, _ := self.GetCharAt(i)
      if c == '\n' {
        line++
        col = 0
      } else {
        col++
      }
    }
    success = SUCCEEDED
  }
  return
}

func (self *GapBuffer) IsDirty() bool { return self.dirty }

func (self *GapBuffer) Length() int { return self.PreLength() + self.PostLength() }

func (self *GapBuffer) String() string { return string(self.Bytes()) }

// For debugging purposes: return two strings, one consisting
// of the characters before the gap, and one of the characters
// after.
func (self *GapBuffer) StringPair() (before string, after string) {
  before = string(self.prechars)

  after_bytes := make([]uint8, 0, self.PostLength())
  for i := self.PostLength(); i > 0; i-- {
    after_bytes = append(after_bytes, self.postchars[i-1])
  }
  after = string(after_bytes)
  return
}

func (self *GapBuffer) Bytes() []uint8 {
  result := make([]uint8, 0, self.Length())
  for i := 0; i < self.PreLength();  i++ {
    result = append(result, self.prechars[i])
  }
  for i := self.PostLength(); i > 0; i-- {
    result = append(result, self.postchars[i-1])
  }
  return result
}

func (self *GapBuffer) AllText() []uint8 {
  len := self.Length()
  text := make([]uint8, len)
  for i := 0; i < len; i++ {
    text[i], _ = self.GetCharAt(i)
  }
  return text
}


////////////////////////////////////////////////////////////////
// Primitives used for building the cursored methods.

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

////////////////////////////////////////////////////////////////
// The cursor-based editor methods, themselves.

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

func (self *GapBuffer) MoveCursorTo(pos int) (result ResultCode) {
  if self.PreLength() > pos {
    dist := self.PreLength() - pos
    result = self.MoveCursorBy(- dist)
  } else {
    dist := pos - self.PreLength()
    result = self.MoveCursorBy(dist)
  }
  return
}

func (self *GapBuffer) MoveCursorBy(dist int) (result ResultCode) {
  if dist < 0 {
	dist = -dist
	for i := int(0); i < dist; i++ {
		result = self.StepCursorBackward()
		if result != SUCCEEDED {
          return
		}
	}
  } else {
    for i := int(0); i < dist; i++ {
      result = self.StepCursorForward()
      if result != SUCCEEDED {
        return
      }
    }
  }
  return SUCCEEDED
}

func (self *GapBuffer) MoveToLine(linenum int) (result ResultCode) {
  self.MoveCursorTo(0)
  for self.line < linenum && self.PostLength() > 0 {
    result = self.StepCursorForward()
    if result != SUCCEEDED {
      return result
    }
  }
  return SUCCEEDED
}

func (self *GapBuffer) MoveToColumn(col int) (result ResultCode) {
  if col > self.column {
    dist := col - self.column
    for i := int(0); i < dist; i++ {
      if self.PeekPost() != '\n' {
        self.StepCursorForward()
      } else {
	    result = INVALID_COLUMN
	    return
      }
    }
  } else {
    dist := self.column - col
    result = self.MoveCursorBy(- int(dist))
  }
  return
}

func (self *GapBuffer) Cut(dist int) (cutbuf []uint8, result ResultCode) {
  self.dirty = true
  result = SUCCEEDED
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

func (self *GapBuffer) Copy(dist int) (copybuf []uint8, result ResultCode) {
  result = SUCCEEDED
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
        result = INVALID
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

func (self *GapBuffer) GetCurrentPosition() int { return len(self.prechars) }

func (self *GapBuffer) GetCurrentLine() int { return self.line }

func (self *GapBuffer) GetCurrentColumn() int { return self.column }

