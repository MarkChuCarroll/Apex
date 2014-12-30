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

import java.io.File
import scala.collection.mutable.MutableList
import scala.collection.mutable.ArraySeq

trait ViewCommand {  
}

case class DisplayChars(
		chars: Seq[Char]) extends ViewCommand
case class DisplayString(s: String) extends ViewCommand	
case class DisplayChar(c: Char) extends ViewCommand
case class SetForegroundColor(color: Int)
case class SetBackgroundColor(color: Int)
case class MoveCursor(line: Int, col: Int) extends ViewCommand
case class ClearToEndOfScreen() extends ViewCommand
case class ClearToEndOfLine() extends ViewCommand
case class InsertBlankLineBefore() extends ViewCommand
case class DeleteLine() extends ViewCommand

/** A controller that queues up commands until the current command
  * queue is requested. 
  * 
  */
trait EditorController {
  /** Resets the command sequence.
    */
  def cleanState
  
  def getCommands: List[ViewCommand]
  
  /** The editor action performed when a user types a character
    * If it's not a newline, inserts the character into the current line.
    * If it is a newline, then breaks the current line.  
    */
  def typeChar(c: Char)
  
  /** The editor action performed when a user types backspace 
   */
  def backspace
  
  /** The editor action performed when the user moves the cursor backwards
    */  
  def back

  /** The editor action performed when the user moves the cursor forwards
    */  
  def forward

  /** The editor action performed when the user moves the cursor up 
    */
  def up

  /** The editor action performed when the user moves the cursor down
   */
  def down

  /** The editor action performed when the user pages down
    */
  def pageDown
  
  /** The editor action performed when the user pages up
    */
  def pageUp
  
  /** The editor action performed when the user scrolls up
    */
  def scrollUp(lines: Int)

  /** The editor action performed when the user scrolls up
    */
  def scrollDown(lines: Int)
  
  /** The editor action to jump the cursor and the view to a line
    */
  def jumpToLine(line: Int)
  
  def save
  def saveAs(file: File)
  
  /** Mark a range of text as being highlighted as a selection.
    * This means changing the background of the characters.  
    */
  def select(start: Int, end: Int)
  
  def selectLines(startLine: Int, endLine: Int)
  def cut
  def copy
  def paste
  def toStart
  def toEnd
  def refresh
}

abstract class EditorServerController(val buf: GapBuffer, val view: ScreenGrid)
    extends EditorController {
  
  val commands: MutableList[ViewCommand] = new MutableList[ViewCommand]
  
  /** the line number in the buffer which is shown on the top line
    * of the view.
    */
  var viewPosition: Int = 1
  var cursorLine: Int = 0
  var cursorColumn: Int = 0
  
  def getTextForLine(displayLine: Int): Seq[Char] = {
    val lineNum = viewPosition + displayLine
    buf.copyLine(lineNum).getOrElse(ArraySeq[Char]())
  }
  
  /** Resets the command sequence.
    */
  def cleanState {
    commands.clear
  }
  
  def getCommands: List[ViewCommand] = {
    commands.toList
  }
  
  /** The editor action performed when a user types a character
    */
  def typeChar(c: Char) {
    if (c == '\n') {
      commands += ClearToEndOfScreen()
      // redraw lines to the end of the screen
    } else {
      commands += ClearToEndOfLine()
      commands += DisplayChar(c)
      val rest = ""
      for {
        c <- rest
      } {
        // insert the new character
        // redraw the rest of the line
      }
    }
  }
  
  def backspace {
    val c: Char = buf.deleteCharBackwards
    if (c == '\n') {
      refresh
    } else {
      val column = buf.currentColumn 
      commands += (MoveCursor(cursorLine, 0))
      val chars = buf.copyLine(buf.currentLine).get      
      commands += (DisplayChars(chars))
      commands += MoveCursor(cursorLine, column)
    }
  }

  def back {
    
  }

  def forward = {
    
  }

  def up

  def down

  def pageDown
  
  def pageUp
  
  def scrollUp(lines: Int)

  def scrollDown(lines: Int)
  
  def jumpToLine(line: Int)
  def save
  def saveAs(file: File)
  def select(start: Int, end: Int)
  def selectLines(startLine: Int, endLine: Int)
  def cut
  def copy
  def paste
  def toStart
  def toEnd
  def refresh  
}