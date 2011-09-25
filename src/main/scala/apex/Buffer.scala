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

class GapBuffer(file: File, initial_size: Int) {
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

  /**
   * move the cursor to a character index.
   */
  def move_to(pos: Int) = {
    move_by(pos - _pre)
  }

  /**
   * Move the cursor to the beginning of a line
   */
  def move_to_line(line: Int) = {
    // This could really use some optimization.
    move_to(0)
    while (_line < line && _post > 0) {
      advance_cursor
    }
  }

  /**
   * Move the cursor to a particular column on the current line.
   */
  def move_to_column(col: Int) = {
    val curcol = current_column
    if (curcol > col) {
      val distance = curcol - col
      move_by(distance)
      if (current_column != col) {
        throw new BufferPositionError(this, col,
          "line didn't contain enough columns")
      }
    }
  }

  /**
   * Move the cursor by a number of lines.
   */
  def move_by_line(l: Int) = {
    val (line, col) = current_line_and_column
    move_to_line(line + l)
  }

  /**
   * Move the cursor by a number of characters.
   */
  def move_by(distance: Int) = {
    if (distance > 0) {
      for (i <- 0 until distance) {
        advance_cursor
      }
    } else if (distance < 0) {
      for (i <- 0 until (-distance)) {
        retreat_cursor
      }
    }
  }

  /**
   * Get the line and column number of the current cursor position.
   */
  def current_line_and_column: (Int, Int) =
    (current_line, current_column)

  /**
   * Insert a string at the current cursor position.
   */
  def insert_string(str: String) = {
    val undo = InsertOperation(this, _pre, str.length)
    push_undo(undo)
    str foreach prim_insert_char
  }

  /**
   * Insert a character at the current cursor position.
   */
  def insert_char(c: Char) = {
    val undo = InsertOperation(this, _pre, 1)
    prim_insert_char(c)
    push_undo(undo)
  }

  /**
   * Insert an array of characters at the current cursor position.
   */
  def insert_chars(cs: Array[Char]) = {
    val undo = InsertOperation(this, _pre, cs.length)
    push_undo(undo)
    for (c <- cs)
      prim_insert_char(c)
  }

  /**
   * Get the character at a position.
   */
  def char_at(pos: Int): Char = {
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

  def delete_char_backwards() = {
    if (_pre > 0) {
      val pos = _pre
      val c = pop_pre()
      reverse_update_position(c)
      val undo = DeleteOperation(this, pos, Array(c))
      push_undo(undo)
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
      if (realsize > _post) {
        realsize = _post
      }
      val result = new Array[Char](realsize)
      for (i <- 0 until realsize) {
        result(i) = pop_post()
      }
      val undo = DeleteOperation(this, _pre, result)
      push_undo(undo)
      return result
    } else {
      var realsize = -s
      if (realsize > _pre) {
        realsize = _pre
      }
      move_by(-realsize)
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
      result(i) = char_at(startpos + i)
    }
    return result
  }

  def current_column = _column

  def current_line = _line

  // absolute position operations: methods that operate on a buffer
  // without any notion of "current point". The underlying
  // implementation still uses a cursor, and these methods do
  // *not* guarantee that the cursor will be unchanged.

  /**
   * Get the number of characters in this buffer.
   * @return the number of characters
   */
  def length: Int = _pre + _post

  def clear = {
    _pre = 0
    _post = 0
  }

  /**
   * Insert a string at an index
   * @param pos the character index where the insert should be performed
   * @param str a string containing the characters to insert
   */
  def insert_string_at(pos: Int, str: String) {
    move_to(pos)
    insert_string(str)
  }

  /**
   * Insert a single character.
   * @param pos the character index where the insert should be performed
   * @param c the character to insert
   */
  def insert_char_at(pos: Int, c: Char) {
    move_to(pos)
    insert_char(c)
  }

  /**
   * Insert an array of characters.
   * @param pos the character index where the insert should be performed
   * @param c the character to insert
   */
  def insert_chars_at(pos: Int, cs: Array[Char]) {
    move_to(pos)
    insert_chars(cs)
  }

  /**
   * Delete a range of characters.
   * @param start the character index of the beginning of the range to delete
   * @param end the character index of the end of the range to delete
   * @return an array containing the delete characters.
   */
  def delete_range(start: Int, end: Int): Array[Char] = {
    move_to(start)
    delete(end - start)
  }

  /**
   * Retrieve a range of characters.
   * @param start the character index of the start of the range
   * @param end the character index of the end of the range
   * @return an array containing the characters in the range
   */
  def copy_range(start: Int, end: Int): Array[Char] = {
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
      result(i) = char_at(start + i)
    }
    result
  }

  /**
   * Convert a line/column to a character index.
   *
   */
  def get_position(linenum: Int, colnum: Int): Int = {
    move_to_line(linenum)
    move_to_column(colnum)
    get_position()
  }

  def get_position_of_line(line: Int): Int = {
    if (line == 1) {
      return 0
    } else {
      var current = 1
      for (i <- 0 until _size) {
        if (char_at(i) == '\n') {
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
  def get_line_and_column(pos: Int): (Int, Int) = {
    if (pos > _size) {
      throw new BufferPositionError(this, pos, "Position past end of buffer")
    }
    var line = 1
    var column = 0
    for (i <- 0 until pos) {
      if (char_at(i) == '\n') {
        line = line + 1
        column = 0
      } else {
        column = column + 1
      }
    }
    (line, column)
  }

  // Undo operations
  private def push_undo(u: UndoOperation) {
    if (!_undoing) {
      _undo_stack.push(u);
    }
  }

  def undo() {
    val u = _undo_stack.pop()
    _undoing = true
    u.execute()
    _undoing = false;
  }

  // Internal primitives
  private def push_pre(c: Char) = {
    check_capacity()
    _prechars(_pre) = c
    _pre += 1
  }

  private def push_post(c: Char) = {
    check_capacity()
    _postchars(_post) = c
    _post += 1
  }

  private def pop_pre(): Char = {
    val result = _prechars(_pre - 1)
    _pre -= 1
    result
  }

  private def pop_post(): Char = {
    val result = _postchars(_post - 1)
    _post -= 1
    result
  }

  // Methods for use in testing and debugging.

  def get_pre(): String = {
    var result = ""
    for (i <- 0 until _pre) {
      result += _prechars(i)
    }
    return result
  }

  def get_post(): String = {
    var result = ""
    for (i <- 0 until _post) {
      result += _postchars(i)
    }
    return result
  }

  def get_position(): Int = _pre

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

  def contents: Array[Char] = {
    val buffer_contents = new Array[Char](length)
    for (i <- 0 until length) {
      buffer_contents(i) = char_at(i)
    }
    buffer_contents
  }

  def read_from_file(file: java.io.File, fail_if_not_found: Boolean) {
    clear
    if (!fail_if_not_found && !file.exists()) {
      return
    }
    val in = new BufferedSource(new FileInputStream(file))
    in.getLines() foreach (line => insert_string(line + '\n'))
    in.close()
  }

  def write_to_file(file: java.io.File, create: Boolean) {

  }

  // --------------------------------
  // Primitives

  private def retreat_cursor = {
    if (_pre > 0) {
      val c = pop_pre()
      push_post(c)
      reverse_update_position(c)
    }
  }

  private def check_capacity() = {
    if ((_pre + _post) >= _size) {
      expand_capacity(2 * _size)
    }
  }

  def expand_capacity(newsize: Int) = {
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
    push_pre(c)
    forward_update_position(c)
  }

  /**
   * For any method that moves the cursor forward - whether by inserting or
   * by simple cursor motion - update the line and column positions.
   */
  private def forward_update_position(c: Char) = {
    if (c == '\n') {
      _line += 1
      _column = 0
    } else {
      _column += 1
    }
  }

  private def advance_cursor = {
    if (_post > 0) {
      val c = pop_post()
      push_pre(c)
      forward_update_position(c)
    }
  }

  /**
   * For any operation that steps the cursor backward, update the
   * line and column positions.
   */
  private def reverse_update_position(c: Char) = {
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
    buf.move_to(pos)
    buf.delete(len)
  }
}

case class DeleteOperation(buf: GapBuffer, pos: Int, dels: Array[Char])
  extends UndoOperation() {
  def execute() {
    buf.move_to(pos)
    buf.insert_chars(dels)
  }
}

class BufferPositionError(b: GapBuffer, pos: Int, msg: String)
  extends Exception {
  val buffer = b
  val requested_position = pos
  val message = msg
}
