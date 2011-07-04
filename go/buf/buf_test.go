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

// File: buf_test.go
// Author: Mark Chu-Carroll <markcc@gmail.com>
// Description: Tests of edit buffers.

package buf

import (
  "fmt"
  "testing"
)

func ExpectBufferValue(t *testing.T, b *GapBuffer, before string, after string) {
  pre, post := b.StringPair()
  if pre != before {
    t.Error(fmt.Sprintf("Buffer value '%v' before gap did not match expected '%v'",
      pre, before))
  }
  if post != after {
    t.Error(fmt.Sprintf("Buffer value '%v' after gap did not match expected '%v'",
      post, after))
  }
}

func ExpectStringEquals(t *testing.T, name string, expected string, actual string) {
  if expected != actual {
    t.Error(fmt.Sprintf("Expected %v to be '%v', but found '%v'", name,
      expected, actual))
  }
}

func TestSingleInsert(t *testing.T) {
  b := NewBuffer(10)
  b.InsertChars([]uint8{'a', 'b', 'c', 'd'})
  ExpectBufferValue(t, b, "abcd", "")
}

func TestInserts(t *testing.T) {
  b := NewBuffer(100)
  b.InsertChars([]uint8{'a', 'b', 'c', 'd'})
  b.InsertChar('q')
  b.InsertChar('r')
  ExpectBufferValue(t, b, "abcdqr", "")
}

func TestExpand(t *testing.T) {
  b := NewBuffer(100)
  b.InsertString("abcdefghijklmnopqrstuvwxyz\n")
  b.InsertString("abcdefghijklmnopqrstuvwxyz\n")
  b.InsertString("abcdefghijklmnopqrstuvwxyz\n")
  b.InsertString("abcdefghijklmnopqrstuvwxyz\n")
  ExpectBufferValue(t, b,
    "abcdefghijklmnopqrstuvwxyz\nabcdefghijklmnopqrstuvwxyz\n"+
      "abcdefghijklmnopqrstuvwxyz\nabcdefghijklmnopqrstuvwxyz\n",
    "")

  if cap(b.prechars) <= 100 {
    t.Error("Expected buffer to expand")
  }
  b.MoveCursorTo(50)
  b.InsertString("abcdefghijklmnopqrstuvwxyz\n")
  b.InsertString("abcdefghijklmnopqrstuvwxyz\n")
  b.InsertString("abcdefghijklmnopqrstuvwxyz\n")
  b.InsertString("abcdefghijklmnopqrstuvwxyz\n")
  b.InsertString("abcdefghijklmnopqrstuvwxyz\n")
  b.InsertString("abcdefghijklmnopqrstuvwxyz\n")
  b.InsertString("abcdefghijklmnopqrstuvwxyz\n")
  b.InsertString("abcdefghijklmnopqrstuvwxyz\n")
  b.InsertString("abcdefghijklmnopqrstuvwxyz\n")
  b.InsertString("abcdefghijklmnopqrstuvwxyz\n")
  ExpectBufferValue(t, b,
    "abcdefghijklmnopqrstuvwxyz\nabcdefghijklmnopqrstuvw"+
      "abcdefghijklmnopqrstuvwxyz\n"+
      "abcdefghijklmnopqrstuvwxyz\n"+
      "abcdefghijklmnopqrstuvwxyz\n"+
      "abcdefghijklmnopqrstuvwxyz\n"+
      "abcdefghijklmnopqrstuvwxyz\n"+
      "abcdefghijklmnopqrstuvwxyz\n"+
      "abcdefghijklmnopqrstuvwxyz\n"+
      "abcdefghijklmnopqrstuvwxyz\n"+
      "abcdefghijklmnopqrstuvwxyz\n"+
      "abcdefghijklmnopqrstuvwxyz\n",
    "xyz\nabcdefghijklmnopqrstuvwxyz\nabcdefghijklmnopqrstuvwxyz\n")
}

func TestInsertAndMove(t *testing.T) {
  b := NewBuffer(100)
  b.InsertChars([]uint8{'a', 'b', 'c', 'd', 'e'})
  b.StepCursorBackward()
  b.StepCursorBackward()
  b.StepCursorBackward()
  ExpectBufferValue(t, b, "ab", "cde")
  b.InsertChars([]uint8{'1', '2', '3'})
  b.StepCursorForward()
  b.StepCursorForward()
  ExpectBufferValue(t, b, "ab123cd", "e")
}

func TestColumnTracking(t *testing.T) {
  b := NewBuffer(100)
  b.InsertChars([]uint8{'a', 'b', 'c', 'd', 'e', 'f', '\n',
    'g', 'h', 'i', 'j', 'k', 'l', '\n',
    'm', 'n', 'o', 'p', 'q', 'r', '\n',
    's', 't', 'u',
  })
  b.MoveCursorTo(12)
  if b.GetCurrentColumn() != 5 {
    t.Error(fmt.Sprintf("Expected column 5, but found %v", b.GetCurrentColumn()))
  }
}

func TestGotoPosition(t *testing.T) {
  b := NewBuffer(100)
  b.InsertChars([]uint8{'a', 'b', 'c', 'd', 'e', 'f', 'g', '\n',
    'h', 'i', 'j',
    'k', 'l', 'm', 'n', 'o', 'p',
  })
  b.MoveCursorTo(4)
  ExpectBufferValue(t, b, "abcd", "efg\nhijklmnop")
  if b.GetCurrentColumn() != 3 {
    t.Error(fmt.Sprintf("Expected buffer to be in column 3, but found column %v",
      b.GetCurrentColumn()))
  }
  b.MoveCursorTo(8)
  ExpectBufferValue(t, b, "abcdefg\n", "hijklmnop")
  if b.GetCurrentColumn() != 0 {
    t.Error(fmt.Sprintf("Expected buffer to be in column 0, but found column %v",
      b.GetCurrentColumn()))
  }
}

func TestCut(t *testing.T) {
  b := NewBuffer(100)
  b.InsertString("abcde\nfghijklm")
  b.MoveCursorTo(4)
  cutbuf, code := b.Cut(5)
  if code != SUCCEEDED {
    t.Error("Cut failed!")
  }
  if len(cutbuf) != 5 {
    t.Error(fmt.Sprintf("Expected cutbuf length = 5, but found %v", len(cutbuf)))
  }
  ExpectStringEquals(t, "cut buffer", "e\nfgh", string(cutbuf))
  ExpectBufferValue(t, b, "abcd", "ijklm")
}

func TestCutBackwards(t *testing.T) {
  b := NewBuffer(100)
  b.InsertString("abcde\nfghijklm")
  b.MoveCursorTo(9)
  cutbuf, code := b.Cut(-5)
  if code != SUCCEEDED {
    t.Error("Cut failed!")
  }
  if len(cutbuf) != 5 {
    t.Error(fmt.Sprintf("Expected cutbuf length = 5, but found %v", len(cutbuf)))
  }
  ExpectStringEquals(t, "cut buffer", "e\nfgh", string(cutbuf))
  ExpectBufferValue(t, b, "abcd", "ijklm")
}

func TestCutPastEnd(t *testing.T) {
  b := NewBuffer(100)
  b.InsertString("abcdefg\nhijklmnop\nqrstuvwxyz\n")
  b.MoveCursorTo(20)
  cutbuf, _ := b.Cut(20)
  if len(cutbuf) != 9 {
    t.Error(fmt.Sprintf("Expected cutbuf length = 9, but found %v", len(cutbuf)))
  }
  ExpectStringEquals(t, "cut buffer", "stuvwxyz\n", string(cutbuf))
  ExpectBufferValue(t, b, "abcdefg\nhijklmnop\nqr", "")
}

func TestCutPastStart(t *testing.T) {
  b := NewBuffer(100)
  b.InsertString("abcdefg\nhijklmnop\nqrstuvwxyz\n")
  b.MoveCursorTo(20)
  cutbuf, _ := b.Cut(-30)
  if len(cutbuf) != 20 {
    t.Error(fmt.Sprintf("Expected cutbuf length = 20, but found %v", len(cutbuf)))
  }
  ExpectStringEquals(t, "cut buffer", "abcdefg\nhijklmnop\nqr", string(cutbuf))
  ExpectBufferValue(t, b, "", "stuvwxyz\n")
}

func TestCopy(t *testing.T) {
  b := NewBuffer(100)
  b.InsertString("abcde\nfghijklm")
  b.MoveCursorTo(4)
  copybuf, code := b.Copy(5)
  if code != SUCCEEDED {
    t.Error("Copy failed!")
  }
  if len(copybuf) != 5 {
    t.Error(fmt.Sprintf("Expected copybuf length = 5, but found %v", len(copybuf)))
  }
  ExpectStringEquals(t, "copy buffer", "e\nfgh", string(copybuf))
  ExpectBufferValue(t, b, "abcd", "e\nfghijklm")
}

func TestCopyBackwards(t *testing.T) {
  b := NewBuffer(100)
  b.InsertString("abcde\nfghijklm")
  b.MoveCursorTo(9)
  copybuf, code := b.Copy(-5)
  if code != SUCCEEDED {
    t.Error("Copy failed!")
  }
  if len(copybuf) != 5 {
    t.Error(fmt.Sprintf("Expected copybuf length = 5, but found %v",
      len(copybuf)))
  }
  ExpectStringEquals(t, "copy buffer", "e\nfgh", string(copybuf))
  ExpectBufferValue(t, b, "abcde\nfgh", "ijklm")
}

func TestCopyPastEnd(t *testing.T) {
  b := NewBuffer(100)
  b.InsertString("abcdefg\nhijklmnop\nqrstuvwxyz\n")
  b.MoveCursorTo(20)
  copybuf, _ := b.Copy(20)
  if len(copybuf) != 9 {
    t.Error(fmt.Sprintf("Expected copybuf length = 9, but found %v", len(copybuf)))
  }
  ExpectStringEquals(t, "copy buffer", "stuvwxyz\n", string(copybuf))
  ExpectBufferValue(t, b, "abcdefg\nhijklmnop\nqr", "stuvwxyz\n")
}

func TestUndoInsert(t *testing.T) {
  b := NewBuffer(100)
  b.InsertString("123456789\n123456789\n")
  b.MoveCursorTo(6)
  b.InsertString("abcd")
  ExpectBufferValue(t, b, "123456abcd", "789\n123456789\n")
  b.Undo()
  ExpectBufferValue(t, b, "123456", "789\n123456789\n")
}

func TestUndoCut(t *testing.T) {
  b := NewBuffer(100)
  b.InsertString("123456789\n123456789\n")
  b.MoveCursorTo(6)
  cutbuf, _ := b.Cut(10)
  ExpectStringEquals(t, "cut buffer", "789\n123456", string(cutbuf))
  ExpectBufferValue(t, b, "123456", "789\n")
  b.Undo()
  // TODO: fill in?????
  //  b.InsertString("123456|789\n123456789\n")
}

//
// Test query methods.
//
func ExpectCharValue(t *testing.T, b *GapBuffer, pos int, expected uint8) {
  c, success := b.GetCharAt(pos)
  if success != SUCCEEDED {
    t.Error(fmt.Sprintf("Retrieving char at position %v failed", pos))
  }
  if c != expected {
    t.Error(fmt.Sprintf("Character at position %v should have been '%v', but found '%v'", pos,
      expected, c))
  }
}


func TestGetCharAt(t *testing.T) {
  b := NewBuffer(100)
  b.InsertString("123456789\n123456789\n")
  ExpectCharValue(t, b, 3, '4')
  ExpectCharValue(t, b, 13, '4')
  b.MoveCursorTo(2)
  ExpectCharValue(t, b, 3, '4')
  ExpectCharValue(t, b, 13, '4')
  b.MoveCursorTo(12)
  ExpectCharValue(t, b, 3, '4')
  ExpectCharValue(t, b, 13, '4')
  _, success := b.GetCharAt(100)
  if success == SUCCEEDED {
    t.Error("Retrieving a character beyond buffer end should have failed")
  }
}

func ExpectLinePosition(t *testing.T, b *GapBuffer, line int, expected int) {
  pos, success := b.GetPositionOfLine(line)
  if success != SUCCEEDED {
    t.Error(fmt.Sprintf("Line position of line '%v' failed", line))
  }
  if pos != expected {
    t.Error(fmt.Sprintf("Expected line '%v' to start at position '%v', but found '%v'",
      line, expected, pos))
  }
}

func TestGetPositionOfLine(t *testing.T) {
  b := NewBuffer(100)
  b.InsertString("123456789\n123456789\n")
  b.InsertString("123456789\n123456789\n")
  b.InsertString("123456789\n123456789\n")
  b.InsertString("123456789\n123456789\n")
  ExpectLinePosition(t, b, 1, 0)
  ExpectLinePosition(t, b, 7, 60)
  b.MoveCursorTo(12)
  ExpectLinePosition(t, b, 1, 0)
  ExpectLinePosition(t, b, 7, 60)
}

func ExpectChars(t *testing.T, b *GapBuffer, start int, end int, expected string) {
  bytes, success := b.GetRange(start, end)
  if success != SUCCEEDED {
    t.Error(fmt.Sprintf("Get Chars %v-%v failed", start, end))
  }
  if string(bytes) != expected {
    t.Error(fmt.Sprintf("Expected string '%v' but found '%v'",
      expected, string(bytes)))
  }
}

func TestGetChars(t *testing.T) {
  b := NewBuffer(100)
  b.InsertString("aaaaabbbb\ncccccdddd\neeeeeffff\nggggghhhh\niiiiijjjj\n")
  ExpectChars(t, b, 10, 20, "cccccdddd\n")
  ExpectChars(t, b, 27, 32, "ff\ngg")
  b.MoveCursorTo(31)
  ExpectChars(t, b, 27, 32, "ff\ngg")
}


func ExpectLineAndColumn(t *testing.T, b *GapBuffer, pos int, e_line int, e_col int) {
  l, c, success := b.GetCoordinates(pos)
  if success != SUCCEEDED {
    t.Error(fmt.Sprintf("Line/col of  position '%v' failed", pos))
  }
  if l != e_line {
    t.Error(fmt.Sprintf("Expected line '%v' for position '%v', but found '%v'",
      e_line, pos, l))
  }
  if c != e_col {
    t.Error(fmt.Sprintf("Expected column '%v' for position '%v', but found '%v'",
      e_col, pos, c))
  }
}

func TestGetLineAndColumn(t *testing.T) {
  b := NewBuffer(100)
  b.InsertString("aaaaabbbb\ncccccdddd\neeeeeffff\nggggghhhh\niiiiijjjj\n")
  ExpectLineAndColumn(t, b, 23, 3, 3)
  ExpectLineAndColumn(t, b, 35, 4, 5)
  ExpectLineAndColumn(t, b, 5, 1, 5)
  ExpectLineAndColumn(t, b, 10, 2, 0)
  b.MoveCursorTo(25)
  ExpectLineAndColumn(t, b, 23, 3, 3)
  ExpectLineAndColumn(t, b, 35, 4, 5)
  ExpectLineAndColumn(t, b, 5, 1, 5)
  ExpectLineAndColumn(t, b, 10, 2, 0)
}


func TestRead(t *testing.T) {
  f, status := NewFileBuffer("tests/foo")
  if status != SUCCEEDED {
    t.Error(fmt.Sprintf("Expected to be able to read file \"tests/foo\"; error '%v'.",
      status))
    return
  }
  ExpectBufferValue(t, f, "Hello world.\nThis is the second line.\nStuff and contents.\n\n", "")
}

/*
func TestWrite(t *testing.T) {
  f, status := NewFileBuffer("tests/foo")
  if status != SUCCEEDED {
    t.Error(fmt.Sprintf("Expected to be able to read file %v; error '%v'.", f.GetFilename(),
      status))
  }
  f.filename = "tests/foo2"
  s := f.Write()
  if s != SUCCEEDED {
    t.Error(fmt.Sprintf("Error writing file '%v', error was '%v'", f.filename,
      s))
  }
}

*/
