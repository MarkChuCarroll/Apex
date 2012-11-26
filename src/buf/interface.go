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

// File: buf/interface.go
// Author: Mark Chu-Carroll <markcc@gmail.com>
// Description: The interface for edit buffers.

package buf

// Every buffer operation that can fail should return a code
// indicating whether the operation succeeded or not, and providing
// an error code if they failed.

/////////////////////////////////////////////////
// Result codes
/////////////////////////////////////////////////

type ResultCode int

const (
  SUCCEEDED ResultCode = iota
  PAST_END
  BEFORE_START
  TOO_LONG
  INVALID
  INVALID_RANGE
  INVALID_REPLACEMENT
  MATCH_FAILED
  INVALID_LINE
  INVALID_COLUMN
  IO_ERROR
)



// A generic buffer interface. This is the primitive API for
// accessing text buffers. It's agnostic about what the underlying
// implementation is: it provides both stateless access methods and
// cursor-based accessors.
type EditBuffer interface {
	// stateless interface methods
	Length() int
	Clear() 
	GetCharAt(pos int) (uint8, ResultCode)
	GetRange(start int, end int) ([]uint8, ResultCode)
	GetPositionOfLine(linenum int) (int, ResultCode)
	GetPositionOfLineAndColumn(linenum int, colnum int) (pos int, result ResultCode);
	GetCoordinates(pos int) (line int, col int, status ResultCode)
	
	// cursor-based interface methods	
	MoveCursorTo(pos int)
	MoveToLine(linenum int)
	MoveCursorBy(distance int)
	StepCursorBackward() ResultCode
	StepCursorForward() ResultCode
	GetCurrentPosition() int
	GetCurrentLine() int
	GetCurrentColumn() int
	InsertChar(c uint8)
	InsertChars(cs []uint8)
	InsertString(s string)
	Cut(numChars int) ([]uint8)
	Copy(numChars int) ([]uint8)
}

type UndoOperation interface {
  Undo()
}

