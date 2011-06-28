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

// File: cursorless.go
// Author: Mark Chu-Carroll <markcc@gmail.com>
// Description: Implementation of the cursorless methods - that is,
//   the methods that operate on the buffer by fully specifying
//   the locations at which things should be done, rather than relying
//   on a cursor point.

package buf

/////////////////////////////////////////////////
// Query operations
//

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
    text[i] = self.GetCharAt(i)
  }
}
