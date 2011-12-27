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

package apex
import scala.io.BufferedSource
import java.io.FileInputStream
import java.io.File

class GapBuffer(var file: File, initialSize: Int) {
  var preChars = new Array[Char](initialSize)
  var postChars = new Array[Char](initialSize)
  var capacity = initialSize
  var preIdx = 0
  var postIdx = 0
  var currentLine = 1
  var currentColumn = 0
  val undoStack = new java.util.Stack[UndoOperation]
  var undoing = false

  def this(file: File) =
    this(file, 65536)

  def this(size: Int) = this(new File("/tmp/scratch"), size)

  def this() = this(new File("/tmp/scratch"))

  def size = preIdx + postIdx

  // Position-based operations: operations which perform edits relative to
  // an implicit cursor point.

  /**
   * move the cursor to a character index.
   */
  def moveTo(pos: Int) = {
    moveBy(pos - preIdx)
  }

  /**
   * Move the cursor to the beginning of a line
   */
  def moveToLine(line: Int) = {
    // This could really use some optimization.
    moveTo(0)
    while (currentLine < line && postIdx > 0) {
      stepCursorForward
    }
  }

  /**
   * Move the cursor to a particular column on the current line.
   */
  def moveToColumn(col: Int) = {
    val curcol = currentColumn
    if (curcol > col) {
      val distance = curcol - col
      moveBy(distance)
      if (currentColumn != col) {
        throw new BufferPositionError(this, col,
          "line didn't contain enough columns")
      }
    }
  }

  /**
   * Move the cursor by a number of lines.
   */
  def moveByLine(l: Int) = {
    val (line, col) = currentLineAndColumn
    moveToLine(line + l)
  }

  /**
   * Move the cursor by a number of characters.
   */
  def moveBy(distance: Int) = {
    if (distance > 0) {
      for (i <- 0 until distance) {
        stepCursorForward
      }
    } else if (distance < 0) {
      for (i <- 0 until (-distance)) {
        stepCursorBackward
      }
    }
  }

  /**
   * Gets the line and column number of the current cursor position.
   */
  def currentLineAndColumn: (Int, Int) =
    (currentLine, currentColumn)

  /**
   * Insert a string at the current cursor position.
   */
  def insertString(str: String) = {
    val undo = InsertOperation(this, preIdx, str.length)
    pushUndo(undo)
    str foreach primInsertChar
  }

  /**
   * Insert a character at the current cursor position.
   */
  def insertChar(c: Char) = {
    val undo = InsertOperation(this, preIdx, 1)
    primInsertChar(c)
    pushUndo(undo)
  }

  /**
   * Insert an array of characters at the current cursor position.
   */
  def insertChars(cs: Array[Char]) = {
    val undo = InsertOperation(this, preIdx, cs.length)
    pushUndo(undo)
    for (c <- cs)
      primInsertChar(c)
  }

  /**
   * Gets the character at a position.
   */
  def charAt(pos: Int): Char = {
    if (pos < 0 || pos >= capacity) {
      return 0
    }
    var c = ' '
    if (pos < preIdx) {
      c = preChars(pos)
    } else {
      val postOffset = pos - preIdx;
      // post is in reverse order - so we need to reverse the index.
      c = postChars(postIdx - postOffset - 1)
    }
    return c
  }

  /**
   * Deletes the character before the cursor
   */
  def deleteCharBackwards = {
    if (preIdx > 0) {
      val pos = preIdx
      val c = popPre
      reverseUpdatePosition(c)
      val undo = DeleteOperation(this, pos, Array(c))
      pushUndo(undo)
    }
  }

  /**
   * Delete a string of characters starting at the current cursor
   * position. If the number of characters to delete is negative, it will
   * delete characters behind the cursor.
   * @return the deleted characters.
   */
  def delete(s: Int): Array[Char] = {
    if (s == 0) {
      return null
    }
    if (s > 0) {
      var realsize = s
      if (realsize > postIdx) {
        realsize = postIdx
      }
      val result = new Array[Char](realsize)
      for (i <- 0 until realsize) {
        result(i) = popPost
      }
      val undo = DeleteOperation(this, preIdx, result)
      pushUndo(undo)
      return result
    } else {
      var realsize = -s
      if (realsize > preIdx) {
        realsize = preIdx
      }
      moveBy(-realsize)
      return delete(realsize)
    }
  }

  /**
   * copy a range of characters starting at the current cursor position.
   */
  def copy(size: Int): Array[Char] = {
    if (size == 0) {
      return null
    }
    var realsize = size
    var startpos = preIdx
    if (size > 0) {
      if (realsize > postIdx) {
        realsize = postIdx
      }
    } else {
      realsize = -size
      if (realsize > preIdx) {
        realsize = preIdx
      }
      startpos -= realsize
    }
    val result = new Array[Char](realsize)
    for (i <- 0 until realsize) {
      result(i) = charAt(startpos + i)
    }
    return result
  }

  // absolute position operations: methods that operate on a buffer
  // without any notion of "current point". The underlying
  // implementation still uses a cursor, and these methods do
  // *not* guarantee that the cursor will be unchanged.

  /**
   * Gets the number of characters in this buffer.
   * @return the number of characters
   */
  def length: Int = preIdx + postIdx

  def clear = {
    preIdx = 0
    postIdx = 0
  }

  /**
   * Insert a string at an index
   * @param pos the character index where the insert should be performed
   * @param str a string containing the characters to insert
   */
  def insertStringAt(pos: Int, str: String) {
    moveTo(pos)
    insertString(str)
  }

  /**
   * Insert a single character.
   * @param pos the character index where the insert should be performed
   * @param c the character to insert
   */
  def insertCharAt(pos: Int, c: Char) {
    moveTo(pos)
    insertChar(c)
  }

  /**
   * Insert an array of characters.
   * @param pos the character index where the insert should be performed
   * @param c the character to insert
   */
  def insertCharsAt(pos: Int, cs: Array[Char]) {
    moveTo(pos)
    insertChars(cs)
  }

  /**
   * Delete a range of characters.
   * @param start the character index of the beginning of the range to delete
   * @param end the character index of the end of the range to delete
   * @return an array containing the delete characters.
   */
  def deleteRange(start: Int, end: Int): Array[Char] = {
    moveTo(start)
    delete(end - start)
  }

  /**
   * Retrieve a range of characters.
   * @param start the character index of the start of the range
   * @param end the character index of the end of the range
   * @return an array containing the characters in the range
   */
  def copyRange(start: Int, end: Int): Array[Char] = {
    if (start > capacity) {
      throw new BufferPositionError(
        this, start,
        "Start of requested range past end of buffer")
    }
    if (end > capacity) {
      throw new BufferPositionError(
        this, end,
        "End of requested range past end of buffer")
    }

    if (end < start) {
      throw new BufferPositionError(
        this, end,
        "End of requested range is greater than start")
    }
    val size = end - start;
    val result = new Array[Char](end - start)
    for (i <- 0 until size) {
      result(i) = charAt(start + i)
    }
    result
  }

  /**
   * Convert a line/column to a character index.
   *
   */
  def positionOf(linenum: Int, colnum: Int): Int = {
    moveToLine(linenum)
    moveToColumn(colnum)
    currentPosition
  }

  def positionOfLine(line: Int): Int = {
    if (line == 1) {
      return 0
    } else {
      var current = 1
      for (i <- 0 until size) {
        if (charAt(i) == '\n') {
          current = current + 1
          if (current == line) {
            return i + 1
          }
        }
      }
      return -1
    }
  }

  /**
   * Convert a character index to line/column
   */
  def lineAndColumnOf(pos: Int): (Int, Int) = {
    if (pos > capacity) {
      throw new BufferPositionError(this, pos, "Position past end of buffer")
    }
    var line = 1
    var column = 0
    for (i <- 0 until pos) {
      if (charAt(i) == '\n') {
        line = line + 1
        column = 0
      } else {
        column = column + 1
      }
    }
    (line, column)
  }

  // Undo operations
  private def pushUndo(u: UndoOperation) {
    if (!undoing) {
      undoStack.push(u);
    }
  }

  def undo {
    val u = undoStack.pop()
    undoing = true
    u.execute
    undoing = false;
  }

  // Internal primitives
  private def pushPre(c: Char) = {
    checkCapacity
    preChars(preIdx) = c
    preIdx += 1
  }

  private def pushPost(c: Char) = {
    checkCapacity
    postChars(postIdx) = c
    postIdx += 1
  }

  private def popPre: Char = {
    val result = preChars(preIdx - 1)
    preIdx -= 1
    result
  }

  private def popPost: Char = {
    val result = postChars(postIdx - 1)
    postIdx -= 1
    result
  }

  // Methods for use in testing and debugging.

  def preString: String = {
    var result = ""
    for (i <- 0 until preIdx) {
      result += preChars(i)
    }
    return result
  }

  def postString: String = {
    var result = ""
    for (i <- 0 until postIdx) {
      result += postChars(i)
    }
    return result
  }

  def currentPosition: Int = preIdx

  override def toString(): String = {
    var result = "{"
    for (i <- 0 until preIdx) {
      result += preChars(i)
    }
    result += "}GAP{"
    for (i <- 0 until postIdx) {
      result += postChars(postIdx - i - 1)
    }
    result += "}"
    result
  }

  def contents: Array[Char] = {
    val bufferContents = new Array[Char](length)
    for (i <- 0 until length) {
      bufferContents(i) = charAt(i)
    }
    bufferContents
  }

  def readFromFile(file: java.io.File, failIfNotFound: Boolean) {
    clear
    if (!failIfNotFound && !file.exists()) {
      return
    }
    val in = new BufferedSource(new FileInputStream(file))
    in.getLines() foreach (line => insertString(line + '\n'))
    in.close()
  }

  def writeToFile(file: java.io.File, create: Boolean) {

  }

  // --------------------------------
  // Primitives

  private def stepCursorBackward = {
    if (preIdx > 0) {
      val c = popPre
      pushPost(c)
      reverseUpdatePosition(c)
    }
  }

  private def checkCapacity = {
    if ((preIdx + postIdx) >= capacity) {
      expandCapacity(2 * capacity)
    }
  }

  def expandCapacity(newsize: Int) = {
    val newpre = new Array[Char](newsize)
    val newpost = new Array[Char](newsize)
    for (i <- 0 until preIdx) {
      newpre(i) = preChars(i)
    }
    for (i <- 0 until postIdx) {
      newpost(newsize - i) = postChars(capacity - i)
    }
    preChars = newpre
    postChars = newpost
    capacity = newsize
  }

  private def primInsertChar(c: Char) = {
    pushPre(c)
    forwardUpdatePosition(c)
  }

  /**
   * For any method that moves the cursor forward - whether by inserting or
   * by simple cursor motion - update the line and column positions.
   */
  private def forwardUpdatePosition(c: Char) = {
    if (c == '\n') {
      currentLine += 1
      currentColumn = 0
    } else {
      currentColumn += 1
    }
  }

  private def stepCursorForward = {
    if (postIdx > 0) {
      val c = popPost
      pushPre(c)
      forwardUpdatePosition(c)
    }
  }

  /**
   * For any operation that steps the cursor backward, update the
   * line and column positions.
   */
  private def reverseUpdatePosition(c: Char) = {
    if (currentLine == 1) {
      currentColumn = preIdx
    } else if (c == '\n') {
      currentLine -= 1
      // Look back to either the beginning of the buffer, or the previous
      // newline. The distance from there to the current position is the
      // current column number. (Not sure if this works correctly on the first
      // line.
      var i = 0
      while ((preIdx - i - 1) > 0 && preChars(preIdx - i - 1) != '\n') {
        i += 1
      }
      currentColumn = i
    } else {
      currentColumn -= 1
    }
  }
}

abstract class UndoOperation {
  def execute
}

case class InsertOperation(buf: GapBuffer, pos: Int, len: Int)
     extends UndoOperation {
  def execute {
    buf.moveTo(pos)
    buf.delete(len)
  }
}

case class DeleteOperation(buf: GapBuffer, pos: Int, dels: Array[Char])
     extends UndoOperation {
  def execute {
    buf.moveTo(pos)
    buf.insertChars(dels)
  }
}

class BufferPositionError(b: GapBuffer, pos: Int, msg: String)
    extends Exception {
  val buffer = b
  val requestedPosition = pos
  val message = msg
}
