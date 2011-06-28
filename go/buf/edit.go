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
// Description: Edit operations that work with the buffer cursor.

package buf

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

