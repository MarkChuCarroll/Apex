package apex

import scala.collection.mutable.Stack
import scala.io.Source
import scala.io.BufferedSource
import java.io.Reader
import scala.io.BufferedSource
import java.io.InputStream
import java.io.OutputStream
import java.io.OutputStreamWriter
import scala.collection.mutable.ArraySeq

class Mark(val buffer: GapBuffer, var position: Int) extends BufferMark {
  var invalidated = false
  def valid: Boolean = !invalidated
  def invalidate {
    invalidated = true
  }
}

class GapBuffer(initialSize: Int) extends Buffer {
  val pre = new BufferStack("pre")
  val post = new BufferStack("post")
  var currentLine = 1
  var currentColumn = 0
  val undoStack = new Stack[UndoOperation]
  var undoing = false
  var marks: Seq[Mark] = new ArraySeq[Mark](0)
  
  def this() {
    this(65536)
  }
  
  def new_mark: BufferMark = {
    val m = new Mark(this, currentPosition)
    marks = marks :+ m
    m
  }
  
  def remove_invalid_marks {
    marks = marks.filter(m => m.valid)
  }

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
    if (pos < 0 || pos > length) {
      throw new BufferPositionError(this, pos, "Move to invalid location")
    }
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
    if (currentLine != line) {
      throw new BufferPositionError(this, line, "Move to invalid line")
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
      if (distance > post.length) {
        throw new BufferPositionError(this, currentPosition + distance, "Tried to move past end of buffer") 
      }
      for (i <- 0 until distance) {
        stepCursorForward
      }
    } else if (distance < 0) {
      if (-distance > pre.length) {
        throw new BufferPositionError(this, currentPosition + distance, "Tried to move past start of buffer") 
      }      
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
    val pos = currentPosition
    val undo = InsertOperation(this, currentPosition, str.length)
    pushUndo(undo)
    str foreach primInsertChar
    marks.foreach(m => if (m.position > pos) { m.position += str.length })
  }

  /** Inserts a character at the current cursor position.
    */
  def insertChar(c: Char) = {
    val pos = currentPosition
    val undo = InsertOperation(this, currentPosition, 1)
    primInsertChar(c)
    pushUndo(undo)
    marks.foreach(m => if (m.position > pos) { m.position += 1 })
  }

  /** Insert an array of characters at the current cursor position.
    */
  def insertChars(cs: Seq[Char]) = {
    val pos = currentPosition
    val undo = InsertOperation(this, currentPosition, cs.length)
    pushUndo(undo)
    for (c <- cs)
      primInsertChar(c)
    marks.foreach(m => if (m.position > pos) { m.position += cs.length })      
  }

  /** Gets the character at a position.
    */
  def charAt(pos: Int): Option[Char] = {
    if (pos < pre.length) {
      pre.at(pos)
    } else if (pos - pre.length < post.length) {
      val idx = post.length - (pos - pre.length) - 1
      post.at(idx)
    } else {
      None
    }
  }

  /** Deletes the character before the cursor
    */
  def deleteCharBackwards: Char = {
    if (pre.length > 0) {
      val pos = currentPosition    
      val c = pre.pop
      reverseUpdatePosition(c)
      val undo = DeleteOperation(this, pre.length, Array(c))
      pushUndo(undo)
      marks.foreach(m => 
        if (m.position == pos) { m.invalidate }
        else if (m.position > pos) {
          m.position -= 1
        }
      )
      c
    } else 0
  }
  
  /** Deletes a string of characters starting at the current cursor position. 
    * @return the deleted characters.
    */
  def delete(numChars: Int): Option[Seq[Char]] = {
    if (numChars == 0 || numChars > post.length) {
      None
    } else {
      val pos = currentPosition
      val result = (0 until numChars).map(_ => post.pop)
      val undo = DeleteOperation(this, pre.length, result)
      pushUndo(undo)
      marks.foreach(m =>
        if (m.position >= pos && m.position < pos + numChars) {
          m.invalidate
        } else if (m.position >= pos + numChars) {
          m.position -= numChars
        }
      )
      Some(result)
    }
  }

  /** Copies a range of characters starting at the current cursor position.
    */
  def copy(size: Int): Seq[Char] = {
    val realsize = if (size < post.length) size else post.length
    (0 until realsize).map(pos => charAt(currentPosition + pos).getOrElse('\0'))
  }
  
  // absolute position operations: methods that operate on a buffer
  // without any notion of "current point". The underlying
  // implementation still uses a cursor. Side-effecting operations do
  // *not* guarantee that the cursor will be unchanged.

  def copyLine(lineNum: Int): Option[Seq[Char]] = {
    copyLines(lineNum, 1) map { line =>
      line.slice(0, line.length - 1)
    }
  }
    
  def copyLines(startLine: Int, numLines: Int): Option[Seq[Char]] = {
    val current = currentPosition
   val result = positionOfLine(startLine).flatMap( startPos => 
      positionOfLine(startLine + numLines).flatMap( endPos =>
        copyRange(startPos, endPos)      
      ))
    moveTo(current)
    result
  } 
  
  def deleteLines(startLine: Int, numLines: Int): Option[Seq[Char]] = {
    positionOfLine(startLine).flatMap({ startPos =>
      val endPos = positionOfLine(startLine + numLines).getOrElse(length)
      deleteRange(startPos, endPos)
    })
  }
  
  /** Inserts a string at an index. 
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
  def insertCharsAt(pos: Int, cs: Seq[Char]) {
    moveTo(pos)
    insertChars(cs)
  }

  /** Deletes a range of characters.
    * @param start the character index of the beginning of the range to delete
    * @param end the character index of the end of the range to delete
    * @return an array containing the delete characters.
    */
  def deleteRange(start: Int, end: Int): Option[Seq[Char]] = {
    moveTo(start)
    delete(end - start)
  }

  /** Retrieves a range of characters.
    * @param start the character index of the start of the range
    * @param end the character index of the end of the range
    * @return an array containing the characters in the range
    */
  def copyRange(start: Int, end: Int): Option[Seq[Char]] = {
    if (start > length) {
      None
    }
    if (end > length) {
      None
    }
    if (end < start) {
      throw new BufferPositionError(
        this, end,
        "End of requested range is greater than start")
    }
    val p = currentPosition
    val size = end - start;
    val result = (0 until size).flatMap(i => charAt(start + i))
    moveTo(p)
    Some(result)
  }

  /** Converts a line/column to a character index.
   */
  def positionOf(linenum: Int, colnum: Int): Option[Int] = {
    val p = currentPosition
    moveToLine(linenum)
    moveToColumn(colnum)
    val result =
      if (currentLine == linenum && currentColumn == colnum) {
        Some(currentPosition)
      } else None
    moveTo(p)
    result
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
      if (charAt(i) == Some('\n')) {
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

  def contents: Seq[Char] = {
    val bufferContents = new Array[Char](length)
    for (i <- 0 until length) {
      bufferContents(i) = charAt(i).get
    }
    bufferContents
  }

  /** Read the contents of the buffer from a file.  
    */
  def read(in: InputStream) {
    clear    
    val source = new BufferedSource(in)
    for (c <- source.toSeq) {
      pre.push(c)
    }
    in.close()
  }

  /** write the contents of the buffer
    */
  def write(out: OutputStream) {
    val writer = new OutputStreamWriter(out)
    writer.write(contents.toArray)
    writer.close
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
class BufferStack(val name: String) {
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
      throw new BufferStackException(s"Text segment $name empty")
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

case class DeleteOperation(buf: GapBuffer, pos: Int, dels: Seq[Char])
     extends UndoOperation {
  def execute {
    buf.moveTo(pos)
    buf.insertChars(dels)
  }
}


