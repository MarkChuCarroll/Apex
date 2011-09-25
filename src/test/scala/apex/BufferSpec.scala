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


import org.specs2.mutable._

object BufferSpec extends Specification {

  "an empty buffer" should {
    "have length 0" in {
      val buf = new GapBuffer()
      buf.length must_== 0
    }

    "insert single strings" in {
      val buf = new GapBuffer()
      buf.insert_string("hello there")
      buf.toString() must_== "{hello there}GAP{}"
    }

    "insert multiple strings the same as a single string concatenated" in {
      val buf1 = new GapBuffer()
      buf1.insert_string("first ")
      buf1.insert_string("second")
      val buf2 = new GapBuffer()
      buf2.insert_string("first second")
      buf1.toString() must_== buf2.toString()
    }
  }

  "a populated buffer" should {
    "move the gap when the cursor is moved" in {
       val buf = new GapBuffer()
       buf.insert_string("1234567890")
       buf.move_by(-3)
       buf.toString() must_== "{1234567}GAP{890}" 
      
    }

    "expand its capacity if an insert makes it too large" in {
      val buf = new GapBuffer(10)
      buf.insert_string("12345678\n")
      buf.insert_string("9,10,11,12,13")
      buf._prechars.length must be_>(10)
      buf._postchars.length must be_>(10)
      buf.toString() must_== "{12345678\n9,10,11,12,13}GAP{}"
    }

    "track columns when the cursor is moved" in {
      val buf = new GapBuffer()
      buf.insert_string("abcdef\nghijkl\nmnopqr\nstu")
      buf.move_to(12)
      buf.current_column should_== 5
    }
  
    "inserts should be at the gap" in {
      val buf = new GapBuffer()
      buf.insert_string("abcde")
      buf.move_by(-3)
      buf.insert_string("123")
      buf.move_by(2)
      buf.toString() must_== "{ab123cd}GAP{e}"
    }

    "cursor position should change as cursor moves" in {
      val buf = new GapBuffer()
      buf.insert_string("abcdefg\nhijklmnop")
      buf.move_to(4)
      buf.current_column must_== 4
      buf.toString() must_==  "{abcd}GAP{efg\nhijklmnop}"

      buf.move_to(8)
      buf.current_column must_== 0
      buf.current_line must_== 2
    }
  }

  "a populated edit buffer supporting cuts" should {

    "do forwards cuts" in {
      val buf = new GapBuffer()
      buf.insert_string("abcde\nfghijklm")
      buf.move_to(4)
      val cut = new String(buf.delete(5))
      cut must_== "e\nfgh"
      buf.toString() must_=="{abcd}GAP{ijklm}"
    }

    "do backwards cuts" in {
      val buf = new GapBuffer()
      buf.insert_string("abcde\nfghijklm")
      buf.move_to(9)
      val cut = new String(buf.delete(-5))
      cut must_== "e\nfgh"
      buf.toString() must_== "{abcd}GAP{ijklm}"
    }

    "truncate cuts that try to go beyond buffer end" in {
      val buf = new GapBuffer()
      buf.insert_string("abcdefg\nhijklmnop\nqrstuvwxyz\n")
      buf.move_to(20)
      val cut = new String(buf.delete(20))
      cut.length() must_== 9
      cut must_== "stuvwxyz\n"
      buf.toString() must_== "{abcdefg\nhijklmnop\nqr}GAP{}"
    }
  }
}
