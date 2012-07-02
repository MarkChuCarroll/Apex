//  Copyright 2011 Mark C. Chu-Carroll
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
import org.junit.Test
import org.junit.Assert._
import org.junit.Before

class BufferSpec() {
  var buf: GapBuffer = null
  
  @Before
  def initBuffer() {
    buf = new GapBuffer()
  }
  

  @Test
  def testEmptyBufferProperties() {
    assertEquals(0, buf.length)
  }

  @Test
  def testInsertString() {
    buf.insertString("hello there")
    assertEquals("{hello there}GAP{}", buf.toString())
  }

  @Test
  def testInsertMultipleStrings() {
    buf.insertString("first ")
    buf.insertString("second")
    val buf2 = new GapBuffer()
    buf2.insertString("first second")
    assertEquals(buf.toString(), buf2.toString())
  }

  @Test
  def testMoveCursor() {
     buf.insertString("1234567890")
     buf.moveBy(-3)
     assertEquals("{1234567}GAP{890}", buf.toString()) 
  }


  @Test
  def testColumnTracking() {
    System.err.println("======== starting column trackingtest")
    buf.insertString("abcdef\nghijkl\nmnopqr\nstu")
    buf.moveTo(12)
    assertEquals(5, buf.currentColumn)
  }
  
  @Test
  def testInsertAtGap() {
    buf.insertString("abcde")
    buf.moveBy(-3)
    buf.insertString("123")
    buf.moveBy(2)
    assertEquals("{ab123cd}GAP{e}", buf.toString())  
  }

  @Test
  def testCursorPositionTracking() {
    buf.insertString("abcdefg\nhijklmnop")
    buf.moveTo(4)
    assertEquals(4, buf.currentColumn)
    assertEquals("{abcd}GAP{efg\nhijklmnop}", buf.toString())  
    buf.moveTo(8)
    assertEquals(0, buf.currentColumn)
    assertEquals(2, buf.currentLine)
  }

  @Test
  def testCutForward() {
    buf.insertString("abcde\nfghijklm")
    buf.moveTo(4)
    val cut = new String(buf.delete(5))
    assertEquals("e\nfgh", cut)
    assertEquals("{abcd}GAP{ijklm}", buf.toString())
  }
  
  @Test
  def testPositionMovement {
    buf.insertString("abcdefg\nhij\nklmnop")
    buf.moveTo(4)
    assertEquals(buf.pre.all, "abcd")    
    assertEquals("efg\nhij\nklmnop", buf.post.all.reverse)
    assertEquals(4, buf.currentColumn)
    assertEquals(1, buf.currentLine)
    buf.moveTo(8)
    assertEquals("abcdefg\n", buf.pre.all)
    assertEquals(0, buf.currentColumn)
    assertEquals(2, buf.currentLine)
  }
  
  @Test
  def testStepForwardAndBack {
    buf.insertString("abcd\nefgh\nijkl\nmnop\n")
    buf.moveTo(12)
    assertEquals("abcd\nefgh\nij", buf.pre.all)
    (1 to 4).foreach(i => buf.stepCursorBackward)
    assertEquals("abcd\nefg", buf.pre.all)
    assertEquals(3, buf.currentColumn)
    assertEquals(2, buf.currentLine)
    for (i <- 1 to 8) {
      buf.stepCursorForward
    }
    assertEquals(16, buf.pre.length)
    assertEquals(4, buf.currentLine)
    assertEquals(1, buf.currentColumn)
  }

  @Test
  def testCutBackward() {
    buf.insertString("abcde\nfghijklm")
    buf.moveTo(9)
    val cut = new String(buf.delete(-5))
    assertEquals("e\nfgh", cut)
    assertEquals("{abcd}GAP{ijklm}", buf.toString()) 
  }
  
  @Test
  def testCutPastEnd() {
    buf.insertString("abcdefg\nhijklmnop\nqrstuvwxyz\n")
    buf.moveTo(20)
    val cut = new String(buf.delete(20))
    assertEquals(9, cut.length())
    assertEquals("stuvwxyz\n", cut)
    assertEquals("{abcdefg\nhijklmnop\nqr}GAP{}", buf.toString()) 
  }
  
  @Test
  def testGetLine() {
    buf.insertString("abcdefg\nhijklmnop\nqrstuvwxyz\n")
//    assertEquals("abcdefg", buf.getLine(1).map(new String(_)).getOrElse("wrong"))
  }
}