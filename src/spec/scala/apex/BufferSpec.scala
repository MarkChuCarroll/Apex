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


import org.specs._

object BufferSpec extends Specification {

  "an empty buffer" should {
    "have length 0" in {
      val buf = new GapBuffer()
      buf.length must_== 0
    }

    "insert single strings" in {
      val buf = new GapBuffer()
      buf.insertString("hello there")
      buf.toString() must_== "{hello there}GAP{}"
    }

    "insert multiple strings the same as a single string concatenated" in {
      val buf1 = new GapBuffer()
      buf1.insertString("first ")
      buf1.insertString("second")
      val buf2 = new GapBuffer()
      buf2.insertString("first second")
      buf1.toString() must_== buf2.toString()
    }
  }

  "a populated buffer" should {
    "move the gap when the cursor is moved" in {
       val buf = new GapBuffer()
       buf.insertString("1234567890")
       buf.moveBy(-3)
       buf.toString() must_== "{1234567}GAP{890}" 
      
    }

    "expand its capacity if an insert makes it too large" in {
      val buf = new GapBuffer(10)
      buf.insertString("12345678\n")
      buf.insertString("9,10,11,12,13")
      buf.preChars.length must be_>(10)
      buf.postChars.length must be_>(10)
      buf.toString() must_== "{12345678\n9,10,11,12,13}GAP{}"
    }

    "track columns when the cursor is moved" in {
      val buf = new GapBuffer()
      buf.insertString("abcdef\nghijkl\nmnopqr\nstu")
      buf.moveTo(12)
      buf.currentColumn must_== 5
    }
  
    "inserts should be at the gap" in {
      val buf = new GapBuffer()
      buf.insertString("abcde")
      buf.moveBy(-3)
      buf.insertString("123")
      buf.moveBy(2)
      buf.toString() must_== "{ab123cd}GAP{e}"
    }

    "cursor position should change as cursor moves" in {
      val buf = new GapBuffer()
      buf.insertString("abcdefg\nhijklmnop")
      buf.moveTo(4)
      buf.currentColumn must_== 4
      buf.toString() must_==  "{abcd}GAP{efg\nhijklmnop}"

      buf.moveTo(8)
      buf.currentColumn must_== 0
      buf.currentLine must_== 2
    }
  }

  "a populated edit buffer supporting cuts" should {

    "do forwards cuts" in {
      val buf = new GapBuffer()
      buf.insertString("abcde\nfghijklm")
      buf.moveTo(4)
      val cut = new String(buf.delete(5))
      cut must_== "e\nfgh"
      buf.toString() must_=="{abcd}GAP{ijklm}"
    }

    "do backwards cuts" in {
      val buf = new GapBuffer()
      buf.insertString("abcde\nfghijklm")
      buf.moveTo(9)
      val cut = new String(buf.delete(-5))
      cut must_== "e\nfgh"
      buf.toString() must_== "{abcd}GAP{ijklm}"
    }

    "truncate cuts that try to go beyond buffer end" in {
      val buf = new GapBuffer()
      buf.insertString("abcdefg\nhijklmnop\nqrstuvwxyz\n")
      buf.moveTo(20)
      val cut = new String(buf.delete(20))
      cut.length() must_== 9
      cut must_== "stuvwxyz\n"
      buf.toString() must_== "{abcdefg\nhijklmnop\nqr}GAP{}"
    }
  }
}
