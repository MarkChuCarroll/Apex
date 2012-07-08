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

import scala.collection.mutable.MutableList
import java.io.File

// Next steps:
// - create a subclass of BufferView which uses cursor control codes to render on screen.
// - Set up some key bindings, to allow simple testing.

object ScreenGrid {
  val ATTR_PLAIN = '0'
  val CHAR_NULL = '\0'
}

object ANSI_CODES {
  // control sequence initiator: the master control code that starts most terminal
  // commands.
  val CSI = "\033["
  def moveCursorUp(lines: Int = 1): String = {
    CSI + lines.toString + "A"
  }

  def moveCursorDown(lines: Int = 1): String = {
    CSI + lines.toString + "B"
  }

  def moveCursorForward(cols: Int = 1): String = {
    CSI + cols.toString + "C"
  }

  def moveCursorBackward(col: Int = 1): String = {
    CSI + col.toString + "D"
  }

  def cursorNextLine(lines: Int = 1): String = {
    CSI + lines.toString + "E"
  }

  def cursorPrevLine(lines: Int = 1): String = {
    CSI + lines.toString + "E"
  }

  def cursorToLine(line: Int): String = {
    CSI + line.toString + "F"
  }

  def cursorToColumn(col: Int): String = {
    CSI + col.toString + "G"
  }

  def cursorToPosition(line: Int, col: Int): String = {
    CSI + line.toString + ";" + col.toString + "H"
  }

  def clearScreen: String = {
    CSI + "2J"
  }

  def setAttributes(att: String): String = {
    CSI + att + "m"
  }

  def RESET: String = "0"
  def ITALIC: String = "3"
  def BOLD: String = "21"
  def NORMAL: String = "22"
  def UNDERLINE_ON: String = "4"
  def UNDERLINE_OFF: String = "24"
  // 0 to 7
  def textColor(num: Int): String = (30 + num).toString
  def defaultTextColor: String = "39"
  def bgColor(num: Int) = (40 + num).toString
  def defaultBgColor = 49
}

/** The model of a displayed editor buffer window.
  * 
  * Each position in the grid has a display character, and an attribute. The
  * attribute is a value which indicates how the position shtould be displayed - essentially
  * a color code.
  * @param lines the number of lines in the window
  * @param columns the number of columns in the window
  */
trait ScreenGrid {
  def lines: Int
  def columns: Int

  /** Sets the character at a screen position.
    */
  def setChar(line: Int, col: Int)(c: Char)
  
  /** Gets the character at a screen position.
    */
  def getChar(line: Int, col: Int): Char

  /** Sets the attribute of a character at a screen position.
    */
  def setAttr(line: Int, col: Int)(a: Char)

  /** Gets the attribute of a character at a screen position.
    */
  def getAttr(line: Int, col: Int): Char

  /** Gets the character and the attribute of a screen position.
    */
  def get(line: Int, col: Int): (Char, Char)

  /** Clears the characters and attributes of all screen positions.
    */
  def clear = {
    for { l <- 0 until lines
          c <- 0 until columns } {
      setChar(l, c)(ScreenGrid.CHAR_NULL)
      setAttr(l, c)(ScreenGrid.ATTR_PLAIN)
    }
  }


  /** Clears the characters and attributes of all screen positions on a line
    */
  def clearLine(line: Int) {
    for { c <- 0 until columns } {
      setChar(line, c)(ScreenGrid.CHAR_NULL)
      setAttr(line, c)(ScreenGrid.ATTR_PLAIN)
    }
  }

  /** Clears the characters and attributes of all screen positions on a line
    * after a column.
    */
  def clearToEndOfLine(line: Int, col: Int) {
    for { c <- col until columns } {
      setChar(line, c)(ScreenGrid.CHAR_NULL)
      setAttr(line, c)(ScreenGrid.ATTR_PLAIN)
    }
  }


  /** Copy the characters and attributes from one line to another.
    */
  def copyLine(from: Int, to: Int) {
    for (col <- 0 until columns) {
      setChar(to, col)(getChar(from, col))
      setAttr(to, col)(getAttr(from, col))
    }
  }

  /** Inserts a blank line - the rest of the lines in the view are shifted
    * downward by one.
    */
  def insertBlankLine(line: Int) {
    if (line < lines) {
      if (line != lines - 1) {
        for { lineToCopy <- (lines - 2) to line by -1 } {
          copyLine(lineToCopy, lineToCopy + 1)
        }
      }
      clearLine(line)
    }
  }

  def removeLine(line: Int) {
    if (line < lines) {
      for (lineToCopy <- line + 1 until lines) {
        // Clear the line we're copying onto.
        clearLine(lineToCopy - 1)
        copyLine(lineToCopy, lineToCopy - 1)
      }
      // Wipe the last line of the screen, which was copied upwards.
      clearLine(lines - 1)
    }
  }

  def render(): String = {
    val buf = new StringBuffer
    for { l <- 0 until lines } {
      for { c <- 0 until columns if getChar(l, c) != ScreenGrid.CHAR_NULL } {
        buf.append(getChar(l, c))
      }
      buf.append('\n')
    }
    buf.toString
  }

  def renderDebug(): String = {
    val buf = new StringBuffer
    for { l <- 0 until lines } {
      buf.append(l)
      buf.append(":|")
      for { c <- 0 until columns if getChar(l, c) != ScreenGrid.CHAR_NULL } {
        buf.append(getChar(l, c))
      }
      buf.append("|\n")
    }
    buf.toString
  }
}

class SimpleScreenGrid(override val lines: Int, override val columns: Int)
    extends ScreenGrid {

  private val cells: Array[Char] = new Array[Char](columns * lines)
  private val attrs: Array[Char] = new Array[Char](columns * lines)

  /** Sets the character at a scren position.
    */
  override def setChar(line: Int, col: Int)(c: Char) {
    cells(line * columns + col) = c
  }
  
  /** Gets the character at a screen position.
    */
  override def getChar(line: Int, col: Int): Char = cells(line * columns + col)

  /** Sets the attribute of a character at a screen position.
    */
  override def setAttr(line: Int, col: Int)(a: Char) {
    attrs(line * columns + col) = a
  }

  /** Gets the attribute of a character at a screen position.
    */
  override def getAttr(line: Int, col: Int): Char = attrs(line * columns + col)

  /** Gets the character and the attribute of a screen position.
    */
  override def get(line: Int, col: Int): (Char, Char) = (getChar(line, col), getAttr(line, col))
}

class AnsiScreenGrid(override val lines: Int, override val columns: Int)
    extends SimpleScreenGrid(lines, columns) {

  private val cells: Array[Char] = new Array[Char](columns * lines)
  private val attrs: Array[Char] = new Array[Char](columns * lines)

  /** Sets the character at a scren position.
    */
  override def setChar(line: Int, col: Int)(c: Char) {
    super.setChar(line, col)(c)
//    System.out.print(ANSI_CODES.cursorToPosition(line, col) + c)
  }
  

  /** Sets the attribute of a character at a screen position.
    */
  override def setAttr(line: Int, col: Int)(a: Char) {
    attrs(line * columns + col) = a
  }

  /** Gets the attribute of a character at a screen position.
    */
  override def getAttr(line: Int, col: Int): Char = attrs(line * columns + col)

  /** Gets the character and the attribute of a screen position.
    */
  override def get(line: Int, col: Int): (Char, Char) = (getChar(line, col), getAttr(line, col))

  override def render(): String = {
    val buf = new StringBuffer
    buf.append(ANSI_CODES.clearScreen)
    buf.append(ANSI_CODES.cursorToPosition(0, 0))
    for { l <- 0 until lines } {
      buf.append(ANSI_CODES.cursorToPosition(l, 0))
      for { c <- 0 until columns if getChar(l, c) != ScreenGrid.CHAR_NULL } {
        buf.append(getChar(l, c))
      }
      buf.append(ANSI_CODES.cursorNextLine(1))
    }
    buf.toString
  }
}

/** A screen grid implementation that maintains two versions of the screen:
  * both a complete grid, and a delta list that contains the sequence of screen
  * update commands necessary to produce the current screen grid knowing the state
  * of the last grid. The first that that displayAt is called, if the viewLine is
  * the same as the last update, then the screen is updated using the delta commands.
  * Otherwise, it does a fill screen update from the grid. 
  */
class DeltaScreenGrid(override val lines: Int, override val columns: Int) extends SimpleScreenGrid(lines, columns) {
  val commands = new MutableList

  override def render(): String = {
    val buf = new StringBuffer
    buf.append("!clear\n")
    buf.append("!goto(0, 0)\n")
    for { l <- 0 until lines } {
      buf.append("!goto(" + l + ", 0)\n")
      for { c <- 0 until columns if getChar(l, c) != ScreenGrid.CHAR_NULL } {
        buf.append("!\"" + getChar(l, c) + "\"\n")
      }
    }
    buf.toString
  }
}

class BufferView(val buffer: Buffer, val grid: ScreenGrid) {
  
  /** Update the display to show the buffer starting at the specified line
    */
  def displayAt(viewLine: Int) {
    for { linenum <- 0 until grid.lines } {
      grid.clearLine(linenum)
      buffer.copyLine(viewLine + linenum).map({ line =>
        (0 until line.length).foreach({ i =>
          if (i < grid.columns) {
            grid.setChar(linenum, i)(line(i))
          }
        })
      })
    }      
  }
}




