// Copyright 2012 Mark C. Chu-Carroll
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
import scala.collection.mutable.ArraySeq
import java.io.OutputStream
import java.io.InputStream

class BufferException(str: String) extends Exception(str)

class BufferStackException(val msg: String) extends BufferException(msg)

class BufferPositionError(b: Buffer, pos: Int, msg: String)
    extends Exception(msg) {
  val buffer = b
  val requestedPosition = pos
  val message = msg
}


case class Selection(val buf: Buffer, val start: Int, val end: Int)


trait EditAction {
  def exec(sel: Selection): Boolean
}

trait Buffer {
  def currentLine: Int
  def currentColumn: Int

  /** Move the cursor forward one position.
    */
  def stepCursorForward
  
  /** Move the cursor backwards one step.
    */
  def stepCursorBackward
  
   /** Gets the number of characters in this buffer.
    * @return the number of characters
    */
  def length: Int

  /** Moves the cursor to a character index.
    */
  def moveTo(pos: Int)

  /** Moves the cursor to the beginning of a line
    */
  def moveToLine(line: Int)

  /** Moves the cursor to a particular column on the current line.
    */
  def moveToColumn(col: Int)

  /** Moves the cursor by a number of characters.
    */
  def moveBy(distance: Int)
  
    /** Moves the cursor by a number of lines.
    * @param numberOfLines
    * @return Some(the final position) or None if the position isn't in the buffer.
    */
  def moveByLines(numberOfLines: Int): Option[Int]

  /** Gets the line and column number of the current cursor position.
    */
  def currentLineAndColumn: (Int, Int)

  /** Inserts a string at the current cursor position.
    */
  def insertString(str: String)

  /** Inserts a character at the current cursor position.
    */
  def insertChar(c: Char)

  /** Insert an array of characters at the current cursor position.
    */
  def insertChars(cs: Seq[Char])

  /** Gets the character at a position.
    */
  def charAt(pos: Int): Option[Char]

  /** Deletes the character before the cursor
    */
  def deleteCharBackwards: Char

  /** Deletes a string of characters starting at the current cursor position. If the number of
    * characters to delete is negative, it will delete characters behind the cursor.
    * @return the deleted characters.
    */
  def delete(s: Int): Option[Seq[Char]]

  /** Copies a range of characters starting at the current cursor position.
    */
  def copy(size: Int): Seq[Char] 
  
  def copyLine(lineNum: Int): Option[Seq[Char]]

  def copyLines(startLine: Int, numLines: Int): Option[Seq[Char]]

  def deleteLines(startLine: Int, numLines: Int): Option[Seq[Char]]
  
  /** Inserts a string at an index
    * @param pos the character index where the insert should be performed
    * @param str a string containing the characters to insert
    */
  def insertStringAt(pos: Int, str: String)

  /** Inserts a single character.
    * @param pos the character index where the insert should be performed
    * @param c the character to insert
    */
  def insertCharAt(pos: Int, c: Char)
  
  /** Inserts an array of characters.
   * @param pos the character index where the insert should be performed
   * @param c the character to insert
   */
  def insertCharsAt(pos: Int, cs: Seq[Char])

  /** Deletes a range of characters.
    * @param start the character index of the beginning of the range to delete
    * @param end the character index of the end of the range to delete
    * @return an array containing the delete characters.
    */
  def deleteRange(start: Int, end: Int): Option[Seq[Char]]

  /** Retrieves a range of characters.
    * @param start the character index of the start of the range
    * @param end the character index of the end of the range
    * @return an array containing the characters in the range
    */
  def copyRange(start: Int, end: Int): Option[Seq[Char]]

  /** Converts a line/column to a character index.
   */
  def positionOf(linenum: Int, colnum: Int): Option[Int]

  /** Returns Some of the character position of the first character in the specified line,
    * or None if the file doesn't have that many lines. 
    */
  def positionOfLine(line: Int): Option[Int]

  /** Convert a character index to line/column
    */
  def lineAndColumnOf(pos: Int): (Int, Int)

  def undo

  def currentPosition: Int

  def contents: Seq[Char]

  /** Read the contents of the buffer from a file or other input source.  
    */
  def read(in: InputStream)

  /** write the contents of the buffer to a file.
    * @param file the filename to write to
    * @param backup true if the current file should be renamed and saved as a backup.
    */
  def write(out: OutputStream)
  
  def new_mark: BufferMark
  
  def remove_invalid_marks: Unit
  
}

/** Marks for buffer locations.
 * A mark is a persistent pointer to a bit of text in a buffer. The mark moves as edits occur, to ensure that
 * it always points to the same text. If the text that it points to is deleted, then the mark becomes
 * invalid.
 */
trait BufferMark {
  def valid: Boolean
  def invalidate: Unit
  def buffer: Buffer
  def position: Int
}
