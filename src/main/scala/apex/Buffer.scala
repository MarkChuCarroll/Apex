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

trait Buffer {
  // Position-based operations: operations which perform edits relative to
  // an implicit cursor point.

  /** moves the cursor to a character index.
    */
  def moveTo(pos: Int)

  /** Moves the cursor to the beginning of a line
    */
  def moveToLine(line: Int) 

  /** Moves the cursor to a particular column on the current line.
    */
  def moveToColumn(col: Int) 

  /** Moves the cursor by a number of lines.
    */
  def moveByLine(l: Int)

  /** Moves the cursor by a number of characters.
    */
  def moveBy(distance: Int)

  /** Returns the line and column number of the current cursor position.
    */
  def currentLineAndColumn: (Int, Int)

  /** Inserts a string at the current cursor position.
    */
  def insertString(str: String)

  /** Insert a character at the current cursor position.
    */
  def insertChar(c: Char) 

  /** Inserts an array of characters at the current cursor position.
    */
  def insertChars(cs: Array[Char])

  /** Returns the character at a position.
    */
  def charAt(pos: Int): Char

  /** Deletes the character before the cursor.
    */
  def deleteCharBackwards() 

  /** Deletes a string of characters starting at the current cursor
    * position. If the number of characters to delete is negative, it will
    * delete characters behind the cursor.
    * @return the deleted characters.
    */
  def delete(s: Int): Array[Char] 

  /** Copies a range of characters starting at the current cursor position.
    *  @return the copied characters
    */
  def copy(size: Int): Array[Char] 

  def currentColumn: Int

  def currentLine: Int

  // absolute position operations: methods that operate on a buffer
  // without any notion of "current point". The underlying
  // implementation still uses a cursor, and these methods do
  // *not* guarantee that the cursor will be unchanged.

  /** Get the number of characters in this buffer.
    * @return the number of characters
    */
  def length: Int

  def clear: Unit

  /** Inserts a string at an index
    * @param pos the character index where the insert should be performed
    * @param str a string containing the characters to insert
    */
  def insertStringAt(pos: Int, str: String) 

  /** Inserts a single character.
    * @param pos the character index where the insert should be performed
    * @param c the character to insert
    */
  def insertCharAt(pos: Int, c: Char): Unit

  /** Inserts an array of characters.
    * @param pos the character index where the insert should be performed
    * @param c the character to insert
    */
  def insertCharsAt(pos: Int, cs: Array[Char]): Unit

  /** Delete a range of characters.
    * @param start the character index of the beginning of the range to delete
    * @param end the character index of the end of the range to delete
    * @return an array containing the delete characters.
    */
  def deleteRange(start: Int, end: Int): Array[Char]

  /** Retrieve a range of characters.
    * @param start the character index of the start of the range
    * @param end the character index of the end of the range
    * @return an array containing the characters in the range
    */
  def copyRange(start: Int, end: Int): Array[Char]
  
  def getLineLength(linenum: Int): Option[Int] = {
    getPositionOfLine(linenum).map({ start_pos =>
      getPositionOfLine(linenum + 1) match {
        case Some(next) => next - start_pos - 1 // exclude the newline 
        case None => length - start_pos
      }      
    })
  }

  /** Gets the contents of a line. Returns None if there is no such line. 
    */
  def getLine(linenum: Int): Option[Array[Char]] 
  
  def getPosition(): Int
  
  /** Convert a line/column to a character index.
    */
  def getPosition(linenum: Int, colnum: Int): Int 

  def getPositionOfLine(line: Int): Option[Int]

  /** Convert a character index to line/column
    */
  def getLineAndColumn(pos: Int): (Int, Int)

  def undo(): Unit

  def contents: Array[Char]

  def readFromFile(file: java.io.File, fail_if_not_found: Boolean): Unit

  def writeToFile(file: java.io.File, create: Boolean): Unit
}

class GapBuffer(file: File, initial_size: Int) extends Buffer {
  var _prechars = new Array[Char](initial_size)
  var _postchars = new Array[Char](initial_size)
  var _size = initial_size
  var _pre = 0
  var _post = 0
  var _line = 1
  var _column = 0
  val _undo_stack = new java.util.Stack[UndoOperation]
  var _undoing = false
  var _file = file

  def this(file: File) =
    this(file, 65536)

  def this(size: Int) = this(new File("/tmp/scratch"), size)

  def this() = this(new File("/tmp/scratch"))

  // Position-based operations: operations which perform edits relative to
  // an implicit cursor point.

  /** moves the cursor to a character index.
    */
  override def moveTo(pos: Int) {
    moveBy(pos - _pre)
  }

  /** Moves the cursor to the beginning of a line
    */
  override def moveToLine(line: Int)  {
    // This could really use some optimization.
    moveTo(0)
    while (_line < line && _post > 0) {
      advanceCursor
    }
  }

  /** Moves the cursor to a particular column on the current line.
    */
  override def moveToColumn(col: Int) {
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

  /** Moves the cursor by a number of lines.
    */
  override def moveByLine(l: Int) {
    val (line, col) = currentLineAndColumn
    moveToLine(line + l)
  }

  /** Moves the cursor by a number of characters.
    */
  override def moveBy(distance: Int) {
    if (distance > 0) {
      for (i <- 0 until distance) {
        advanceCursor
      }
    } else if (distance < 0) {
      for (i <- 0 until (-distance)) {
        retreatCursor
      }
    }
  }

  /** Returns the line and column number of the current cursor position.
    */
  override def currentLineAndColumn: (Int, Int) =
    (currentLine, currentColumn)

  /** Inserts a string at the current cursor position.
    */
  override def insertString(str: String) = {
    val undo = InsertOperation(this, _pre, str.length)
    pushUndo(undo)
    str foreach prim_insert_char
  }

  /** Insert a character at the current cursor position.
    */
  override def insertChar(c: Char) = {
    val undo = InsertOperation(this, _pre, 1)
    prim_insert_char(c)
    pushUndo(undo)
  }

  /** Inserts an array of characters at the current cursor position.
    */
  override def insertChars(cs: Array[Char]) = {
    val undo = InsertOperation(this, _pre, cs.length)
    pushUndo(undo)
    for (c <- cs)
      prim_insert_char(c)
  }

  /** Returns the character at a position.
    */
  override def charAt(pos: Int): Char = {
    if (pos < 0 || pos >= _size) {
      return 0
    }
    var c = ' '
    if (pos < _pre) {
      c = _prechars(pos)
    } else {
      val post_offset = pos - _pre;
      // post is in reverse order - so we need to reverse the index.
      c = _postchars(_post - post_offset - 1)
    }
    return c
  }

  /** Deletes the character before the cursor.
    */
  override def deleteCharBackwards() = {
    if (_pre > 0) {
      val pos = _pre
      val c = popPre()
      reverseUpdatePosition(c)
      val undo = DeleteOperation(this, pos, Array(c))
      pushUndo(undo)
    }
  }

  /** Deletes a string of characters starting at the current cursor
    * position. If the number of characters to delete is negative, it will
    * delete characters behind the cursor.
    * @return the deleted characters.
    */
  override def delete(s: Int): Array[Char] = {
    if (s == 0) {
      return null
    }
    if (s > 0) {
      var realsize = s
      if (realsize > _post) {
        realsize = _post
      }
      val result = new Array[Char](realsize)
      for (i <- 0 until realsize) {
        result(i) = popPost()
      }
      val undo = DeleteOperation(this, _pre, result)
      pushUndo(undo)
      return result
    } else {
      var realsize = -s
      if (realsize > _pre) {
        realsize = _pre
      }
      moveBy(-realsize)
      return delete(realsize)
    }
  }

  /** Copies a range of characters starting at the current cursor position.
    *  @return the copied characters
    */
  override def copy(size: Int): Array[Char] = {
    if (size == 0) {
      return null
    }
    var realsize = size
    var startpos = _pre
    if (size > 0) {
      if (realsize > _post) {
        realsize = _post
      }
    } else {
      realsize = -size
      if (realsize > _pre) {
        realsize = _pre
      }
      startpos -= realsize
    }
    val result = new Array[Char](realsize)
    for (i <- 0 until realsize) {
      result(i) = charAt(startpos + i)
    }
    return result
  }

  override def currentColumn = _column

  override def currentLine = _line

  // absolute position operations: methods that operate on a buffer
  // without any notion of "current point". The underlying
  // implementation still uses a cursor, and these methods do
  // *not* guarantee that the cursor will be unchanged.

  /** Get the number of characters in this buffer.
    * @return the number of characters
    */
  override def length: Int = _pre + _post

  override def clear = {
    _pre = 0
    _post = 0
  }

  /** Inserts a string at an index
    * @param pos the character index where the insert should be performed
    * @param str a string containing the characters to insert
    */
  override def insertStringAt(pos: Int, str: String) {
    moveTo(pos)
    insertString(str)
  }

  /** Inserts a single character.
    * @param pos the character index where the insert should be performed
    * @param c the character to insert
    */
  override def insertCharAt(pos: Int, c: Char) {
    moveTo(pos)
    insertChar(c)
  }

  /** Inserts an array of characters.
    * @param pos the character index where the insert should be performed
    * @param c the character to insert
    */
  override def insertCharsAt(pos: Int, cs: Array[Char]) {
    moveTo(pos)
    insertChars(cs)
  }

  /** Delete a range of characters.
    * @param start the character index of the beginning of the range to delete
    * @param end the character index of the end of the range to delete
    * @return an array containing the delete characters.
    */
  override def deleteRange(start: Int, end: Int): Array[Char] = {
    moveTo(start)
    delete(end - start)
  }

  /** Retrieve a range of characters.
    * @param start the character index of the start of the range
    * @param end the character index of the end of the range
    * @return an array containing the characters in the range
    */
  override def copyRange(start: Int, end: Int): Array[Char] = {
    if (start > _size) {
      throw new BufferPositionError(
        this, start,
        "Start of requested range past end of buffer")
    }
    if (end > _size) {
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

  /** Convert a line/column to a character index.
    */
  override def getPosition(linenum: Int, colnum: Int): Int = {
    moveToLine(linenum)
    moveToColumn(colnum)
    getPosition()
  }

  override def getPositionOfLine(line: Int): Option[Int] = {
    if (line == 1) {
      return Some(0)
    } else {
      var current = 1
      for (i <- 0 until _size) {
        if (charAt(i) == '\n') {
          current = current + 1
          if (current == line) {
            return Some(i + 1)
          }
        }
      }
      return None
    }
  }

  /** Convert a character index to line/column
    */
  override def getLineAndColumn(pos: Int): (Int, Int) = {
    if (pos > _size) {
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
    if (!_undoing) {
      _undo_stack.push(u);
    }
  }

  override def undo() {
    val u = _undo_stack.pop()
    _undoing = true
    u.execute()
    _undoing = false;
  }

  // Internal primitives
  private def pushPre(c: Char) = {
    checkCapacity()
    _prechars(_pre) = c
    _pre += 1
  }

  private def pushPost(c: Char) = {
    checkCapacity()
    _postchars(_post) = c
    _post += 1
  }

  private def popPre(): Char = {
    val result = _prechars(_pre - 1)
    _pre -= 1
    result
  }

  private def popPost(): Char = {
    val result = _postchars(_post - 1)
    _post -= 1
    result
  }

  // Methods for use in testing and debugging.

  def getPre(): String = {
    var result = ""
    for (i <- 0 until _pre) {
      result += _prechars(i)
    }
    return result
  }

  def getPost(): String = {
    var result = ""
    for (i <- 0 until _post) {
      result += _postchars(i)
    }
    return result
  }

  override def getPosition(): Int = _pre

  override def toString(): String = {
    var result = "{"
    for (i <- 0 until _pre) {
      result += _prechars(i)
    }
    result += "}GAP{"
    for (i <- 0 until _post) {
      result += _postchars(_post - i - 1)
    }
    result += "}"
    result
  }

  override def contents: Array[Char] = {
    val buffer_contents = new Array[Char](length)
    for (i <- 0 to length) {
      buffer_contents(i) = charAt(i)
    }
    buffer_contents
  }
  
  def getLine(linenum: Int): Option[Array[Char]] = {
    getPositionOfLine(linenum) map { startPos =>
      // .get is safe, because it only returns None when getPositionOfLine is also None.
      val lineLen = getLineLength(linenum).get
      System.err.println("Line length = " + lineLen)
      val lineChars = new Array[Char](lineLen)
      for (i <- 0 to lineLen - 1) {
        lineChars(i) = charAt(startPos + i) 
      }
      lineChars
    }
  }  

  override def readFromFile(file: java.io.File, fail_if_not_found: Boolean) {
    clear
    if (!fail_if_not_found && !file.exists()) {
      return
    }
    val in = new BufferedSource(new FileInputStream(file))
    in.getLines() foreach (line => insertString(line + '\n'))
    in.close()
  }

  override def writeToFile(file: java.io.File, create: Boolean) {

  }

  // --------------------------------
  // Primitives

  private def retreatCursor = {
    if (_pre > 0) {
      val c = popPre()
      pushPost(c)
      reverseUpdatePosition(c)
    }
  }

  private def checkCapacity() = {
    if ((_pre + _post) >= _size) {
      expandCapacity(2 * _size)
    }
  }

  def expandCapacity(newsize: Int) = {
    val newpre = new Array[Char](newsize)
    val newpost = new Array[Char](newsize)
    for (i <- 0 until _pre) {
      newpre(i) = _prechars(i)
    }
    for (i <- 0 until _post) {
      newpost(newsize - i) = _postchars(_size - i)
    }
    _prechars = newpre
    _postchars = newpost
    _size = newsize
  }

  private def prim_insert_char(c: Char) = {
    pushPre(c)
    forwardUpdatePosition(c)
  }

  /** For any method that moves the cursor forward - whether by inserting or
    * by simple cursor motion - update the line and column positions.
    */
  private def forwardUpdatePosition(c: Char) = {
    if (c == '\n') {
      _line += 1
      _column = 0
    } else {
      _column += 1
    }
  }

  private def advanceCursor = {
    if (_post > 0) {
      val c = popPost()
      pushPre(c)
      forwardUpdatePosition(c)
    }
  }

  /** For any operation that steps the cursor backward, update the
    * line and column positions.
    */
  private def reverseUpdatePosition(c: Char) = {
    if (_line == 1) {
      _column = _pre
    } else if (c == '\n') {
      _line -= 1
      // Look back to either the beginning of the buffer, or the previous
      // newline. The distance from there to the current position is the
      // current column number. (Not sure if this works correctly on the first
      // line.
      var i = 0
      while ((_pre - i - 1) > 0 && _prechars(_pre - i - 1) != '\n') {
        i += 1
      }
      _column = i
    } else {
      _column -= 1
    }
  }
}

abstract class UndoOperation() {
  def execute()
}

case class InsertOperation(buf: GapBuffer, pos: Int, len: Int)
  extends UndoOperation() {
  def execute() {
    buf.moveTo(pos)
    buf.delete(len)
  }
}

case class DeleteOperation(buf: GapBuffer, pos: Int, dels: Array[Char])
  extends UndoOperation() {
  def execute() {
    buf.moveTo(pos)
    buf.insertChars(dels)
  }
}

class BufferPositionError(b: GapBuffer, pos: Int, msg: String)
  extends Exception {
  val buffer = b
  val requested_position = pos
  val message = msg
}

