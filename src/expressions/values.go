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

// File: values.go
// Author: Mark Chu-Carroll <markcc@gmail.com>
// Description: Location and range values.

package expressions

import (
  "buf"
)


/////////////////////////////////////////////////
// Locations
/////////////////////////////////////////////////


type Location interface {
  GetBuffer() buf.Buffer
  GetAbsolute() int
  AsStartRelative() *StartRelativeLocation
  AsEndRelative() *EndRelativeLocation
}

// Check that the location specifies a valid position in
// the buffer.
func Validate(loc Location) buf.Status {
  pos := loc.GetAbsolute()
  if pos < 0 {
    return buf.NewFailure(buf.BEFORE_START, "Location points to before file start")
  }
  if pos > loc.GetBuffer().Length() {
    return buf.NewFailure(buf.PAST_END, "Location points to past file end")
  }
  return buf.NewSuccess()
}

/////////////////////////////////////////////////

type StartRelativeLocation struct {
  buffer buf.Buffer
  loc    int
}

func NewStartRelativeLocation(b buf.Buffer, pos int) *StartRelativeLocation {
  return &StartRelativeLocation{b, pos}
}

func (self *StartRelativeLocation) GetAbsolute() int {
  return self.loc
}

func (self *StartRelativeLocation) GetBuffer() buf.Buffer {
  return self.buffer
}

func (self *StartRelativeLocation) AsStartRelative() *StartRelativeLocation {
  return self
}

func (self *StartRelativeLocation) AsEndRelative() *EndRelativeLocation {
  return NewEndRelativeLocation(self.GetBuffer(), self.GetBuffer().Length()-self.GetAbsolute())
}

/////////////////////////////////////////////////

type EndRelativeLocation struct {
  buffer buf.Buffer
  offset int
}

func NewEndRelativeLocation(b buf.Buffer, offset int) *EndRelativeLocation {
  return &EndRelativeLocation{b, offset}
}

func (self *EndRelativeLocation) GetBuffer() buf.Buffer {
  return self.buffer
}

func (self *EndRelativeLocation) GetAbsolute() int {
  return self.buffer.Length() - self.offset
}

func (self *EndRelativeLocation) AsStartRelative() *StartRelativeLocation {
  return NewStartRelativeLocation(self.GetBuffer(), self.GetAbsolute())
}

func (self *EndRelativeLocation) AsEndRelative() *EndRelativeLocation {
  return self
}


/////////////////////////////////////////////////
// Ranges: the result of evaluating a range
// expression. These are concrete ranges.
/////////////////////////////////////////////////

type Range interface {
  GetBuffer() buf.Buffer
  GetStart() Location
  GetEnd() Location
  GetContents() ([]uint8, buf.Status)
  GetLength() int
  // Normalize converts a location into a form where its start
  // is represented relative to the buffer start, and its end
  // is represented relative to the buffer end.
  Normalize() Range
  Validate() buf.Status
}

///////////////////////////////////////////////////////////

// A SimpleRange is the basic version of a range,
// represented by the start and end-points.
type SimpleRange struct {
  buffer buf.Buffer
  start  Location
  end    Location
}

func NewSimpleRange(buffer buf.Buffer, start Location, end Location) *SimpleRange {
  return &SimpleRange{buffer, start, end}
}

func (self *SimpleRange) GetBuffer() buf.Buffer {
  return self.buffer
}

func (self *SimpleRange) GetStart() Location { return self.start }

func (self *SimpleRange) GetEnd() Location { return self.end }

func (self *SimpleRange) GetLength() int {
  return self.GetEnd().GetAbsolute() - self.GetStart().GetAbsolute()
}

func (self *SimpleRange) Validate() buf.Status {
  start_status := Validate(self.start)
  if start_status.GetResultCode() != buf.SUCCEEDED {
    return start_status
  }
  end_status := Validate(self.end)
  if end_status.GetResultCode() != buf.SUCCEEDED {
    return end_status
  }
  if self.start.GetAbsolute() > self.end.GetAbsolute() {
    return buf.NewFailure(buf.INVALID_RANGE, "Start comes after end")
  }
  return buf.NewSuccess()
}

func (self *SimpleRange) GetContents() (chars []uint8, status buf.Status) {
  status = self.Validate()
  if status.GetResultCode() != buf.SUCCEEDED {
    chars = nil
    return
  }
  chars, _ = self.buffer.GetChars(self.start.GetAbsolute(),
    self.end.GetAbsolute())
  status = buf.NewSuccess()
  return
}

func (self *SimpleRange) Normalize() Range {
  return NewSimpleRange(self.GetBuffer(), self.start.AsStartRelative(), self.end.AsEndRelative())
}

///////////////////////////////////////////////
func TestCutPastEnd(t *testing.T) {
  b := New(100)
  b.InsertString(
  b.MoveTo(20)
  cutbuf, _ := b.Cut(20)
  if len(cutbuf) != 9 {
    t.Error(fmt.Sprintf("Expected cutbuf length = 9, but found %v", len(cutbuf)))
  }
  ExpectStringEquals(t, "cut buffer", "stuvwxyz\n", string(cutbuf))
  ExpectBufferValue(t, b, 
}

func TestCutPastEnd(t *testing.T) {
  b := New(100)
  b.InsertString(
  b.MoveTo(20)
  cutbuf, _ := b.Cut(20)
  if len(cutbuf) != 9 {
    t.Error(fmt.Sprintf("Expected cutbuf length = 9, but found %v", len(cutbuf)))
  }
  ExpectStringEquals(t, "cut buffer", "stuvwxyz\n", string(cutbuf))
  ExpectBufferValue(t, b, 
}

func TestCutPastEnd(t *testing.T) {
  b := New(100)
  b.InsertString(
  b.MoveTo(20)
  cutbuf, _ := b.Cut(20)
  if len(cutbuf) != 9 {
    t.Error(fmt.Sprintf("Expected cutbuf length = 9, but found %v", len(cutbuf)))
  }
  ExpectStringEquals(t, "cut buffer", "stuvwxyz\n", string(cutbuf))
  ExpectBufferValue(t, b, 
}

func TestCutPastEnd(t *testing.T) {
  b := New(100)
  b.InsertString(
  b.MoveTo(20)
  cutbuf, _ := b.Cut(20)
  if len(cutbuf) != 9 {
    t.Error(fmt.Sprintf("Expected cutbuf length = 9, but found %v", len(cutbuf)))
  }
  ExpectStringEquals(t, "cut buffer", "stuvwxyz\n", string(cutbuf))
  ExpectBufferValue(t, b, 
}

func TestCutPastEnd(t *testing.T) {
  b := New(100)
  b.InsertString(
  b.MoveTo(20)
  cutbuf, _ := b.Cut(20)
  if len(cutbuf) != 9 {
    t.Error(fmt.Sprintf("Expected cutbuf length = 9, but found %v", len(cutbuf)))
  }
  ExpectStringEquals(t, "cut buffer", "stuvwxyz\n", string(cutbuf))
  ExpectBufferValue(t, b, 
}

func TestCutPastEnd(t *testing.T) {
  b := New(100)
  b.InsertString(
  b.MoveTo(20)
  cutbuf, _ := b.Cut(20)
  if len(cutbuf) != 9 {
    t.Error(fmt.Sprintf("Expected cutbuf length = 9, but found %v", len(cutbuf)))
  }
  ExpectStringEquals(t, "cut buffer", "stuvwxyz\n", string(cutbuf))
  ExpectBufferValue(t, b, 
}

func TestCutPastEnd(t *testing.T) {
  b := New(100)
  b.InsertString(
  b.MoveTo(20)
  cutbuf, _ := b.Cut(20)
  if len(cutbuf) != 9 {
    t.Error(fmt.Sprintf("Expected cutbuf length = 9, but found %v", len(cutbuf)))
  }
  ExpectStringEquals(t, "cut buffer", "stuvwxyz\n", string(cutbuf))
  ExpectBufferValue(t, b, 
}

func TestCutPastEnd(t *testing.T) {
  b := New(100)
  b.InsertString(
  b.MoveTo(20)
  cutbuf, _ := b.Cut(20)
  if len(cutbuf) != 9 {
    t.Error(fmt.Sprintf("Expected cutbuf length = 9, but found %v", len(cutbuf)))
  }
  ExpectStringEquals(t, "cut buffer", "stuvwxyz\n", string(cutbuf))
  ExpectBufferValue(t, b, 
}

func TestCutPastEnd(t *testing.T) {
  b := New(100)
  b.InsertString(
  b.MoveTo(20)
  cutbuf, _ := b.Cut(20)
  if len(cutbuf) != 9 {
    t.Error(fmt.Sprintf("Expected cutbuf length = 9, but found %v", len(cutbuf)))
  }
  ExpectStringEquals(t, "cut buffer", "stuvwxyz\n", string(cutbuf))
  ExpectBufferValue(t, b, 
}

func TestCutPastEnd(t *testing.T) {
  b := New(100)
  b.InsertString(
  b.MoveTo(20)
  cutbuf, _ := b.Cut(20)
  if len(cutbuf) != 9 {
    t.Error(fmt.Sprintf("Expected cutbuf length = 9, but found %v", len(cutbuf)))
  }
  ExpectStringEquals(t, "cut buffer", "stuvwxyz\n", string(cutbuf))
  ExpectBufferValue(t, b, 
}

func TestCutPastEnd(t *testing.T) {
  b := New(100)
  b.InsertString(
  b.MoveTo(20)
  cutbuf, _ := b.Cut(20)
  if len(cutbuf) != 9 {
    t.Error(fmt.Sprintf("Expected cutbuf length = 9, but found %v", len(cutbuf)))
  }
  ExpectStringEquals(t, "cut buffer", "stuvwxyz\n", string(cutbuf))
  ExpectBufferValue(t, b, 
}

func TestCutPastEnd(t *testing.T) {
  b := New(100)
  b.InsertString(
  b.MoveTo(20)
  cutbuf, _ := b.Cut(20)
  if len(cutbuf) != 9 {
    t.Error(fmt.Sprintf("Expected cutbuf length = 9, but found %v", len(cutbuf)))
  }
  ExpectStringEquals(t, "cut buffer", "stuvwxyz\n", string(cutbuf))
  ExpectBufferValue(t, b, 
}

func TestCutPastEnd(t *testing.T) {
  b := New(100)
  b.InsertString(
  b.MoveTo(20)
  cutbuf, _ := b.Cut(20)
  if len(cutbuf) != 9 {
    t.Error(fmt.Sprintf("Expected cutbuf length = 9, but found %v", len(cutbuf)))
  }
  ExpectStringEquals(t, "cut buffer", "stuvwxyz\n", string(cutbuf))
  ExpectBufferValue(t, b, 
}

func TestCutPastEnd(t *testing.T) {
  b := New(100)
  b.InsertString(
  b.MoveTo(20)
  cutbuf, _ := b.Cut(20)
  if len(cutbuf) != 9 {
    t.Error(fmt.Sprintf("Expected cutbuf length = 9, but found %v", len(cutbuf)))
  }
  ExpectStringEquals(t, "cut buffer", "stuvwxyz\n", string(cutbuf))
  ExpectBufferValue(t, b, 
}

func TestCutPastEnd(t *testing.T) {
  b := New(100)
  b.InsertString(
  b.MoveTo(20)
  cutbuf, _ := b.Cut(20)
  if len(cutbuf) != 9 {
    t.Error(fmt.Sprintf("Expected cutbuf length = 9, but found %v", len(cutbuf)))
  }
  ExpectStringEquals(t, "cut buffer", "stuvwxyz\n", string(cutbuf))
  ExpectBufferValue(t, b, 
}
////////////

// A PatternRange is a range augmented with information
// about a regular expression match.
type PatternRange struct {
  SimpleRange
  sel          []uint8 // the range that was used as the selection to run the match
  match_ranges []int
}

func NewPatternRange(buffer buf.Buffer, start Location, end Location, sel []uint8, matches []int) *PatternRange {
  // We assume that this will only be called with a selection that
  // has already been tested for validity.
  return &PatternRange{SimpleRange{buffer, start, end}, sel, matches}
}

func (self *PatternRange) GetNumberOfMatches() int {
  return len(self.match_ranges)/2 - 1
}

func (self *PatternRange) GetMatch(i int) []uint8 {
  // This assumes that the range is valid!
  if i > self.GetNumberOfMatches() {
    return nil
  }
  range_start := self.match_ranges[i*2]
  range_end := self.match_ranges[i*2+1]
  return self.sel[range_start:range_end]
}

func (self *PatternRange) Normalize() Range {
  return NewPatternRange(self.GetBuffer(), self.start.AsStartRelative(), self.end.AsEndRelative(),
    self.sel, self.match_ranges)
}
