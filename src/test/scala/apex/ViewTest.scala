//  Copyright 2012 Mark C. Chu-Carroll
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

import org.junit.{Before, Test}
import org.junit.Assert._

class ScreenGridTest {
  val grid = new SimpleScreenGrid(6, 10)

  @Before
  def setupScreenGrid() {
    grid.clear
    grid.setChar(0, 0)(' ')
    grid.setChar(0, 1)(' ')
    grid.setChar(0, 2)('h')
    grid.setChar(0, 3)('e')
    grid.setChar(0, 4)('l')
    grid.setChar(0, 5)('l')
    grid.setChar(0, 6)('o')

    grid.setChar(1, 0)('w')
    grid.setChar(1, 1)('o')
    grid.setChar(1, 2)('r')
    grid.setChar(1, 3)('l')
    grid.setChar(1, 4)('d')

    grid.setChar(3, 0)('f')
    grid.setChar(3, 1)('o')
    grid.setChar(3, 2)('o')
  }

  @Test
  def testRender() {
    assertEquals("  hello\nworld\n\nfoo\n\n\n", grid.render())
  }

  @Test
  def testRemoveLine() {
    grid.removeLine(1)
    assertEquals("  hello\n\nfoo\n\n\n\n", grid.render())
    grid.removeLine(1)
    assertEquals("  hello\nfoo\n\n\n\n\n", grid.render())
    grid.removeLine(0)
    assertEquals("foo\n\n\n\n\n\n", grid.render())
    // Put some text on the last line.
    grid.setChar(5, 0)('a')
    grid.setChar(5, 1)('b')
    grid.setChar(5, 2)('c')
    // Remove the last line.
    grid.removeLine(5)
    assertEquals("foo\n\n\n\n\n\n", grid.render())
    // Remove the empty last line.
    grid.removeLine(5)
    assertEquals("foo\n\n\n\n\n\n", grid.render())
  }

  @Test
  def testInsertBlankLine() {
    grid.insertBlankLine(1)
    assertEquals("  hello\n\nworld\n\nfoo\n\n", grid.render())
    grid.insertBlankLine(0)
    assertEquals("\n  hello\n\nworld\n\nfoo\n", grid.render())
    grid.removeLine(0)
    grid.insertBlankLine(5)
    assertEquals("  hello\n\nworld\n\nfoo\n\n", grid.render())
    grid.insertBlankLine(4)
    assertEquals("  hello\n\nworld\n\n\nfoo\n", grid.render())
    grid.insertBlankLine(5)
    assertEquals("  hello\n\nworld\n\n\n\n", grid.render())
  }

}


class ViewTest {
  var grid: ScreenGrid = null
  var buf: GapBuffer = null
  var view: BufferView = null

  @Before
  def initBufferAndView() {
    grid = new SimpleScreenGrid(6, 10)
    buf = new GapBuffer()
    buf.insertString("1111\n2222\n3333\n4444\n5555\n6666\n7777\n8888\n")
    view = new BufferView(buf, grid)
  }

  @Test
  def testDisplayAt() {
    view.displayAt(2)
    assertEquals("0:|2222|\n1:|3333|\n2:|4444|\n3:|5555|\n4:|6666|\n5:|7777|\n",
                 grid.renderDebug())
    view.displayAt(4)
    assertEquals("0:|4444|\n1:|5555|\n2:|6666|\n3:|7777|\n4:|8888|\n5:||\n",
                 grid.renderDebug())
    
  }

}
