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
import java.io.{File, FileInputStream, FileWriter}
import java.util.Stack

class BufferStackException extends Exception { }

class BufferStackUnderflowException extends BufferStackException { }

class GapBuffer(var file: File, initialSize: Int) {
  val pre = new BufferStack
  val post = new BufferStack
  var currentLine = 1
  var currentColumn = 0
  val undoStack = new Stack[UndoOperation]
  var undoing = false

  def this(file: File) =
    this(file, 65536)

  def this(size: Int) = this(new File("/tmp/scratch"), size)

  def this() = this(new File("/tmp/scratch"))

  // Position-based operations: operations which perform edits relative to
  // an implicit cursor point.
  
  /** Move the cursor forward one position.
    */
  def stepCursorForward = {
    if (post.length > 0) {
      val c = post.pop
      pre.push(c)
      forwardUpdatePosition(c)
    }
  }  
  
  /** Move the cursor backwards one step.
    */
  def stepCursorBackward = {
    if (pre.length > 0) {
      val c = pre.pop
      post.push(c)
      reverseUpdatePosition(c)
    }
  }
  
   /** Gets the number of characters in this buffer.
    * @return the number of characters
    */
  def length: Int = pre.length + post.length  

  /** Moves the cursor to a character index.
    */
  def moveTo(pos: Int) = {
    moveBy(pos - pre.length)
  }

  /** Moves the cursor to the beginning of a line
    */
  def moveToLine(line: Int) = {
    // This could really use some optimization.
    moveTo(0)
    while (currentLine < line && post.length > 0) {
      stepCursorForward
    }
  }

  /** Moves the cursor to a particular column on the current line.
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

  /** Moves the cursor by a number of characters.
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
  
    /** Moves the cursor by a number of lines.
    * @param numberOfLines
    * @return Some(the final position) or None if the position isn't in the buffer.
    */
  def moveByLines(numberOfLines: Int): Option[Int] = {
    val targetLine = currentLine + numberOfLines 
    if (targetLine < 0) {
      moveTo(0)
      None
    } else if (targetLine > currentLine) {
      while (post.length > 0 && currentLine != targetLine) { stepCursorForward }
      if (currentLine != targetLine) Some(pre.length) else None
    } else {
      while (pre.length > 0 && currentLine != targetLine) { stepCursorBackward }
      moveToColumn(0)
      Some(0)
    }
  }

  /** Gets the line and column number of the current cursor position.
    */
  def currentLineAndColumn: (Int, Int) = (currentLine, currentColumn)

  /** Inserts a string at the current cursor position.
    */
  def insertString(str: String) = {
    val undo = InsertOperation(this, pre.length, str.length)
    pushUndo(undo)
    str foreach primInsertChar
  }

  /** Inserts a character at the current cursor position.
    */
  def insertChar(c: Char) = {
    val undo = InsertOperation(this, pre.length, 1)
    primInsertChar(c)
    pushUndo(undo)
  }

  /** Insert an array of characters at the current cursor position.
    */
  def insertChars(cs: Array[Char]) = {
    val undo = InsertOperation(this, pre.length, cs.length)
    pushUndo(undo)
    for (c <- cs)
      primInsertChar(c)
  }

  /** Gets the character at a position.
    */
  def charAt(pos: Int): Option[Char] = {
    if (pos > length) {
      None
    } else if (pos < pre.length) {
      pre.at(pos)
    } else {
      val idx = post.length - (pos - pre.length) - 1
      post.at(idx)
    }
  }

  /** Deletes the character before the cursor
    */
  def deleteCharBackwards = {
    if (pre.length > 0) {
      val c = pre.pop
      reverseUpdatePosition(c)
      val undo = DeleteOperation(this, pre.length, Array(c))
      pushUndo(undo)
    }
  }

  /** Deletes a string of characters starting at the current cursor position. If the number of
    * characters to delete is negative, it will delete characters behind the cursor.
    * @return the deleted characters.
    */
  def delete(s: Int): Array[Char] = {
    if (s == 0) {
      return null
    }
    if (s > 0) {
      var realsize = s
      if (realsize > post.length) {
        realsize = post.length
      }
      val result = new Array[Char](realsize)
      for (i <- 0 until realsize) {
        result(i) = post.pop
      }
      val undo = DeleteOperation(this, pre.length, result)
      pushUndo(undo)
      return result
    } else {
      var realsize = -s
      if (realsize > pre.length) {
        realsize = pre.length
      }
      moveBy(-realsize)
      return delete(realsize)
    }
  }

  /** Copies a range of characters starting at the current cursor position.
    */
  def copy(size: Int): Array[Char] = {
    if (size == 0) {
      return null
    }
    var realsize = size
    var startpos = pre.length
    if (size > 0) {
      if (realsize > post.length) {
        realsize = post.length
      }
    } else {
      realsize = -size
      if (realsize > pre.length) {
        realsize = pre.length
      }
      startpos -= realsize
    }
    val result = new Array[Char](realsize)
    for (i <- 0 until realsize) {
      result(i) = charAt(startpos + i).get
    }
    return result
  }
  
  

  // absolute position operations: methods that operate on a buffer
  // without any notion of "current point". The underlying
  // implementation still uses a cursor, and these methods do
  // *not* guarantee that the cursor will be unchanged.

  def getLine(lineNum: Int): Option[Array[Char]] = {
    positionOfLine(lineNum).map({ start =>
      val endMaybeWithNewline = positionOfLine(lineNum + 1).map(_ - 1).getOrElse(length)
      val end = if (charAt(endMaybeWithNewline) == '\n') {
          endMaybeWithNewline - 1
        } else {
          endMaybeWithNewline
        }
      val len = end - start
      val result = new Array[Char](len)
      for (i <- 0 until len) {
        result(i) = charAt(start + i).get
      }
      result
    })
  }

  def copyLines(startLine: Int, numLines: Int): Option[Array[Char]] = {
    moveToLine(startLine) 
    if (currentLine != startLine) None
    else {
      val endLine = startLine + numLines - 1
      val buf =  new StringBuilder
      while (currentLine <= endLine && post.length > 0) {
        buf.append(charAt(pre.length).get)
        stepCursorForward
      }
      Some(buf.toArray)
    }
  }

  def deleteLines(startLine: Int, numLines: Int): Option[Array[Char]] = {
    positionOfLine(startLine).map({ startPos =>
      val endPos = positionOfLine(startLine + numLines).getOrElse(length)
      deleteRange(startPos, endPos)
    })
  }
  
  /** Inserts a string at an index
    * @param pos the character index where the insert should be performed
    * @param str a string containing the characters to insert
    */
  def insertStringAt(pos: Int, str: String) {
    moveTo(pos)
    insertString(str)
  }

  /** Inserts a single character.
    * @param pos the character index where the insert should be performed
    * @param c the character to insert
    */
  def insertCharAt(pos: Int, c: Char) {
    moveTo(pos)
    insertChar(c)
  }

  /** Inserts an array of characters.
   * @param pos the character index where the insert should be performed
   * @param c the character to insert
   */
  def insertCharsAt(pos: Int, cs: Array[Char]) {
    moveTo(pos)
    insertChars(cs)
  }

  /** Deletes a range of characters.
    * @param start the character index of the beginning of the range to delete
    * @param end the character index of the end of the range to delete
    * @return an array containing the delete characters.
    */
  def deleteRange(start: Int, end: Int): Array[Char] = {
    moveTo(start)
    delete(end - start)
  }

  /** Retrieves a range of characters.
    * @param start the character index of the start of the range
    * @param end the character index of the end of the range
    * @return an array containing the characters in the range
    */
  def copyRange(start: Int, end: Int): Array[Char] = {
    if (start > length) {
      throw new BufferPositionError(
        this, start,
        "Start of requested range past end of buffer")
    }
    if (end > length) {
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
      result(i) = charAt(start + i).get
    }
    result
  }

  /** Converts a line/column to a character index.
   */
  def positionOf(linenum: Int, colnum: Int): Option[Int] = {
    moveToLine(linenum)
    moveToColumn(colnum)
    if (currentLine == linenum && currentColumn == colnum) {
      Some(currentPosition)
    } else None
  }

  /** Returns Some of the character position of the first character in the specified line,
    * or None if the file doesn't have that many lines. 
    */
  def positionOfLine(line: Int): Option[Int] = {
    val newlineAsInt: Int = ('\n'.toInt)
    if (line == 1) {
      Some(0)
    } else {
      var current = 1
      for (i <- 0 until length) {
        if (charAt(i) == Some('\n')) {
          current = current + 1
          if (current == line) {
            return Some(i + 1)
          }
        }
      }
      None
    }
  }

  /** Convert a character index to line/column
    */
  def lineAndColumnOf(pos: Int): (Int, Int) = {
    if (pos > length) {
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

  def currentPosition: Int = pre.length


  def contents: Array[Char] = {
    val bufferContents = new Array[Char](length)
    for (i <- 0 until length) {
      bufferContents(i) = charAt(i).get
    }
    bufferContents
  }

  /** Read the contents of the buffer from a file.  
    */
  def readFromFile(file: java.io.File, failIfNotFound: Boolean) {
    clear
    if (!failIfNotFound && !file.exists()) {
      return
    }
    val in = new BufferedSource(new FileInputStream(file))
    for (c <- in.toSeq) {
      pre.push(c)
    }
    in.close()
  }

  /** write the contents of the buffer to a file.
    * @param file the filename to write to
    * @param backup true if the current file should be renamed and saved as a backup.
    */
  def writeToFile(file: java.io.File, backup: Boolean) {
    if (file.exists && backup) {
      val backupFile = new File(file.getName + ".BAK")
      if (backupFile.exists) {
        backupFile.delete
      }
      file.renameTo(backupFile)
    }
    val out = new FileWriter(file)
    out.write(contents)
    out.close
  }
  

  /** Resets the buffer state to match the file.
    */
  private def clear = {
    pre.reset
    post.reset
  }

  private def primInsertChar(c: Char) = {
    pre.push(c)
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

  /**
   * For any operation that steps the cursor backward, update the
   * line and column positions.
   */
  private def reverseUpdatePosition(c: Char) = {
    if (currentLine == 1) {
      currentColumn = pre.length
    } else if (c == '\n') {
      currentLine -= 1
      // Look back to either the beginning of the buffer, or the previous
      // newline. The distance from there to the current position is the
      // current column number. (Not sure if this works correctly on the first
      // line.
      var i = 0
      while ((pre.length - i - 1) > 0 && pre.at(pre.length - i - 1) != Some('\n')) {
        i += 1
      }
      currentColumn = i
    } else {
      currentColumn -= 1
    }
  }

  // Debugging only
  override def toString(): String = {
    var result = "{"
    for (i <- 0 until pre.length) {
      result += pre.at(i).get
    }
    result += "}GAP{"
    for (i <- 0 until post.length) {
      result += post.at(post.length - i - 1).get
    }
    result += "}"
    result
  }

}


/** A utility class which represents one of the two contiguous text segments of
  * a gap buffer. The main behavior of this resembles a stack, where you push or
  * pop off of the active end. 
  */
class BufferStack {
  val initLength = 65536

  /** Returns the current number of characters in this buffer.
    */
  var length: Int = 0
  
  private var chars: Array[Char] = new Array[Char](initLength)
  
  /** Pushes a character onto the end of the buffer.
    */
  def push(c: Char) {
    if (length >= chars.length - 1) {
      val newChars = new Array[Char](chars.length * 2)
      for (i <- 0 until length) {
        newChars(i) = chars(i)
        chars = newChars
      }
    }
    chars(length) = c
    length = length + 1
  }
  
  /** Pops character off the active end of the buffer. Throws an
    * exception if the buffer is empty. 
    */
  def pop: Char = {
    if (length > 0) {
      length = length - 1
      chars(length)
    } else {
      throw new BufferStackUnderflowException
    } 
  }
  
  /** Get the character at a position.
    * @param idx the character position 
    * @return Some of the character, or else None.
    */ 
  def at(idx: Int): Option[Char] = {
    if (idx >= length) {
      None
    } else {
      Some(chars(idx))
    }
  }
  
  /** Get the entire buffer 
    */
  def all: String = new String(chars.slice(0, length))

  /** Reset the buffer, deleting all of its contents.
    */
  def reset {
    length = 0
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
