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

package org.scientopia.goodmath.apex

import junit.framework.TestCase
import org.scalatest.junit.AssertionsForJUnit
import scala.collection.mutable.ListBuffer
import junit.framework.Assert._

class BufferTest extends TestCase with AssertionsForJUnit {
  def testEmpty() = {
    val buf = new GapBuffer(100)
    assert(buf.length() == 0)
  }
	
  def testSingleInsert() = {
    val buf = new GapBuffer(100)
    buf.insert_string("hello there")
    assertEquals("{hello there}GAP{}", buf.toString())
  }

  def testMultipleInsert() = {
    // multiple inserts should be the same as a single insert.
    val buf1 = new GapBuffer(100)
    buf1.insert_string("first ")
    buf1.insert_string("second")
    val buf2 = new GapBuffer(100)
    buf2.insert_string("first second")
    assertEquals(buf2.toString(), buf1.toString());
  }

  def testCursorMotion() = {
    val buf = new GapBuffer(100)
    buf.insert_string("1234567890")
    buf.move_by(-3)
    assertEquals("{1234567}GAP{890}", buf.toString())
  }

  def testExpandCapacity() {
    val buf = new GapBuffer(10)
    buf.insert_string("12345678\n")
    buf.insert_string("9,10,11,12,13")
    assert(buf._prechars.length > 10)
    assert(buf._postchars.length > 10)
    assertEquals("{12345678\n9,10,11,12,13}GAP{}", buf.toString())
  }

  def testMoveWithinLine() {
    val buf = new GapBuffer(100)
    buf.insert_string("abcdef\nghijkl\nmnopqr\nstu")
    buf.move_to(12)
    assertEquals(5, buf.get_current_column())
  }

  def testMoveAroundAndInsert() {
    val buf = new GapBuffer(100)
    buf.insert_string("abcde")
    buf.move_by(-3)
    buf.insert_string("123")
    buf.move_by(2)
    assertEquals("{ab123cd}GAP{e}", buf.toString())
  }

  def testMoveCursorAndCursorPosition() {
    val buf = new GapBuffer(100)
    buf.insert_string("abcdefg\nhijklmnop")
    buf.move_to(4)
    assertEquals(4, buf.get_current_column())
    assertEquals("{abcd}GAP{efg\nhijklmnop}", buf.toString())
  }

  def testMoveCursorMultipleTimesAndCheckColumn() {
    val buf = new GapBuffer(100)
    buf.insert_string("abcdefg\nhijklmnop")
    buf.move_to(4)
    buf.move_to(8)
    assertEquals(0, buf.get_current_column())
    assertEquals("{abcdefg\n}GAP{hijklmnop}", buf.toString())
  }

  def testCutForwards() {
    val buf = new GapBuffer(100)
    buf.insert_string("abcde\nfghijklm")
    buf.move_to(4)
    val cut = new String(buf.delete(5))
    assertEquals("e\nfgh", cut)
	assertEquals("{abcd}GAP{ijklm}", buf.toString())
  }

  def testCutBackwards() {
    val buf = new GapBuffer(100)
    buf.insert_string("abcde\nfghijklm");
    buf.move_to(9)
    val cut = new String(buf.delete(-5))
    assertEquals("e\nfgh", cut)
	assertEquals("{abcd}GAP{ijklm}", buf.toString())
  }

  def testCutPastEnd() {
    val buf = new GapBuffer(100)
    buf.insert_string("abcdefg\nhijklmnop\nqrstuvwxyz\n")
    buf.move_to(20)
    val cut = new String(buf.delete(20))
    assertEquals(9, cut.length())
    assertEquals("stuvwxyz\n", cut)
    assertEquals("{abcdefg\nhijklmnop\nqr}GAP{}", buf.toString())
  }
}
