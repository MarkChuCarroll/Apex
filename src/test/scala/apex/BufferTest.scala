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

import org.junit.{Before, Test}
import org.junit.Assert._

class BufferTest {
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
  def testMarksWithInsertAtGap() {
    buf.insertString("abcde")
    buf.moveBy(-1)   
    val m1 = buf.new_mark
    val p1 = m1.position
    buf.moveTo(1)
    val m2 = buf.new_mark
    val p2 = m2.position
    buf.moveTo(p1)
    buf.moveBy(-2)
    buf.insertString("123")
    buf.moveBy(2)
    assertEquals("{ab123cd}GAP{e}", buf.toString())    
    assertNotEquals(p1, m1.position)
    assertEquals(buf.length - 1, m1.position)
    assertEquals(p2, m2.position)
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
    val cut = buf.delete(5).map(seq_to_string).get
    assertEquals("e\nfgh", cut)
    assertEquals("{abcd}GAP{ijklm}", buf.toString())
  }

  @Test
  def testMarksWithCutForward() {
    buf.insertString("abcde\nfghijklm")
    buf.moveTo(7)    
    val mark_in_cut = buf.new_mark
    buf.moveTo(12)
    val mark_after_cut = buf.new_mark
    val orig_mark_after_cut = mark_after_cut.position
    buf.moveTo(2)
    val mark_before_cut = buf.new_mark
    val pos_mark_before_cut = mark_before_cut.position
    buf.moveTo(4)
    val cut = buf.delete(5).map(seq_to_string).get
    assertEquals("e\nfgh", cut)
    assertEquals("{abcd}GAP{ijklm}", buf.toString())
    assertEquals(pos_mark_before_cut, mark_before_cut.position)
    assertNotEquals(orig_mark_after_cut, mark_after_cut.position)
    assertFalse(mark_in_cut.valid)
    assertTrue(mark_after_cut.valid)
    assertTrue(mark_before_cut.valid)
    assertEquals(orig_mark_after_cut - 5, mark_after_cut.position)
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

//  @Test
//  def testCutBackward() {
//    buf.insertString("abcde\nfghijklm")
//    buf.moveTo(9)
//    val cut = (buf.delete(-5))
//    assertEquals("e\nfgh", cut)
//    assertEquals("{abcd}GAP{ijklm}", buf.toString()) 
//  }
  
  @Test
  def testCutPastEnd() {
    buf.insertString("abcdefg\nhijklmnop\nqrstuvwxyz\n")
    buf.moveTo(20)
    val cut = buf.delete(20)
    assertEquals(None, cut)
    //assertEquals(9, cut.length)
    //assertEquals("stuvwxyz\n", cut)
    //assertEquals("{abcdefg\nhijklmnop\nqr}GAP{}", buf.toString()) 
  }
  
  @Test
  def testGetPositionOfLine() {
    buf.insertString("abcdefg\nhijklmnop\nqrstuvwxyz\n")
    assertEquals(Some(8), buf.positionOfLine(2))
  } 
  
  @Test
  def testCopyLine() {
    buf.insertString("abcdefg\nhijklmnop\nqrstuvwxyz\n")
    assertEquals("abcdefg", buf.copyLine(1).map(seq_to_string).getOrElse("wrong"))
    assertEquals("hijklmnop", buf.copyLine(2).map(seq_to_string).getOrElse("wrong"))
    assertEquals("qrstuvwxyz", buf.copyLine(3).map(seq_to_string).getOrElse("wrong"))
    assertEquals("wrong", buf.copyLine(4).map(seq_to_string).getOrElse("wrong"))
    assertEquals(None, buf.copyLine(5))
  }
  
  def seq_to_string(s: Seq[Char]): String = {
    val sb = new StringBuffer
    s.map(sb.append(_))
    sb.toString
  }

  @Test
  def testCopyLines() {
    buf.insertString("1111\n2222\n3333\n4444\n5555\n6666\n7777\n8888\n")
    assertEquals("3333\n4444\n", seq_to_string(buf.copyLines(3, 2).get))
    assertEquals("7777\n8888\n", seq_to_string(buf.copyLines(7, 2).get))
    // Behavior change: make fetching a range that goes past the end of the buffer return None
    // assertEquals("7777\n8888\n", seq_to_string(buf.copyLines(7, 10).get))
    assertEquals(None, buf.copyLines(10, 12))
    assertEquals("1111\n", buf.copyLines(1, 1).map(seq_to_string).get)
  }

  @Test
  def testDeleteLines() {
    buf.insertString("1111\n2222\n3333\n4444\n5555\n6666\n7777\n8888\n")
    assertEquals("3333\n4444\n", buf.deleteLines(3, 2).map(seq_to_string).get)
    assertEquals("1111\n2222\n5555\n6666\n7777\n8888\n", seq_to_string(buf.contents))
    assertEquals("7777\n8888\n", buf.deleteLines(5, 3).map(seq_to_string).get)
    assertEquals("1111\n2222\n5555\n6666\n", seq_to_string(buf.contents))
  }
}
