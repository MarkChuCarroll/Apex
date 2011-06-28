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
package apex;

import static org.junit.Assert.*;

import org.junit.Before;
import org.junit.Test;

public class GapBufferTest {
   GapBuffer _buf;
   
   @Before
   public void setUp() {
      _buf = new GapBuffer(100);
   }

   @Test
   public void testClear() {
      _buf.insert("abcdefg");
      assertEquals(7, _buf.length());
      _buf.clear();
      assertEquals(0, _buf.length());
      assertEquals(0, _buf.currentColumn());
      assertEquals(1, _buf.currentLine());
   }

   @Test
   public void testInsert() {
      _buf.insert("abcdef\nghijk");
      assertEquals("{abcdef\nghijk}GAP{}", _buf.debugString());
      _buf.moveBy(-3);
      assertEquals("{abcdef\ngh}GAP{ijk}", _buf.debugString());
   }

   
   @Test
   public void testMoveTo() {
      _buf.insert("123456789\n123456789\n123456789\n");
      _buf.moveTo(12);
      assertEquals(12, _buf.currentPosition());
      assertEquals("{123456789\n12}GAP{3456789\n123456789\n}", _buf.debugString());
      _buf.moveTo(16);
      assertEquals(16, _buf.currentPosition());
      assertEquals("{123456789\n123456}GAP{789\n123456789\n}", _buf.debugString());      
   }

   @Test
   public void testMoveToLine() {
      _buf.insert("123456789\n123456789\n123456789\n");
      _buf.moveToLine(3);
      assertEquals(20, _buf.currentPosition());
   }
   
   @Test
   public void testCut() {
      _buf.insert("123456789\n123456789\n123456789\n");
      _buf.moveTo(12);
      char[] cuts = _buf.cut(10);
      assertEquals("3456789\n12", new String(cuts));
      assertEquals("{123456789\n12}GAP{3456789\n}", _buf.debugString());
      cuts = _buf.cut(10);
      assertEquals("3456789\n", new String(cuts));
      assertEquals("{123456789\n12}GAP{}", _buf.debugString());
   }

   @Test
   public void testCopy() {
      _buf.insert("123456789\n123456789\n123456789\n");
      _buf.moveTo(12);
      char[] copys = _buf.copy(10);
      assertEquals("3456789\n12", new String(copys));
      assertEquals("{123456789\n12}GAP{3456789\n123456789\n}", _buf.debugString());
      _buf.moveBy(12);
      assertEquals(24, _buf.currentPosition());
      copys = _buf.copy(20);
      assertEquals("56789\n", new String(copys));
   }
   
   @Test
   public void testColumnPosition() {
      _buf.insert("123456789\n123456789\n123456789\n");
      _buf.moveTo(12);
      assertEquals(2, _buf.currentColumn());
      _buf.moveTo(0);
      assertEquals(0, _buf.currentColumn());
      _buf.moveTo(14);
      assertEquals(4, _buf.currentColumn());
   }
   
   @Test
   public void testLinePosition() {
      _buf.insert("123456789\n123456789\n123456789\n123456789\n123456789\n");
      _buf.moveTo(12);
      assertEquals(2, _buf.currentLine());
      _buf.moveTo(28);
      assertEquals(3, _buf.currentLine());
      _buf.moveTo(2);
      assertEquals(1, _buf.currentLine());
   }
   
   @Test
   public void testGetLineAndColumnOf() {
      _buf.insert("123456789\n123456789\n123456789\n123456789\n123456789\n");
      int[] pos = _buf.getLineAndColumnOf(27);
      assertEquals(3, pos[0]);
      assertEquals(7, pos[1]);
   }

   @Test
   public void testGetPositionOfLine() {
      _buf.insert("123456789\n123456789\n123456789\n123456789\n123456789\n");
      assertEquals(20, _buf.getPositionOfLine(3));
   }

   @Test
   public void testUndoCut() {
      _buf.insert("123456789\n123456789\n123456789\n123456789\n123456789\n");
      _buf.moveTo(12);
      _buf.cut(5);
      assertEquals("{123456789\n12}GAP{89\n123456789\n123456789\n123456789\n}",
            _buf.debugString());
      _buf.undo();
      // Reposition the cursor: undo doesn't make any guarantees about cursor position.
      _buf.moveTo(12);
      assertEquals("{123456789\n12}GAP{3456789\n123456789\n123456789\n123456789\n}",
            _buf.debugString());
   }

   public void testUndoInsert() {
      _buf.insert("123456789\n123456789\n123456789\n");
      _buf.moveTo(12);
      _buf.insert("ABC");
      assertEquals("{123456789\n12ABC}GAP{3456789\n123456789\n}",
            _buf.debugString());
      assertTrue(_buf.undo());
      // Reposition the cursor: undo doesn't make any guarantees about cursor position.
      _buf.moveTo(12);
      assertEquals("{123456789\n12}GAP{3456789\n123456789\n",
            _buf.debugString());
   }
   
   public void testUndoMultiple() {
      _buf.insert("123456789\n123456789\n123456789\n");
      _buf.moveTo(12);
      _buf.insert("ABC");
      assertEquals("{123456789\n12ABC}GAP{3456789\n123456789\n}",
            _buf.debugString());
      _buf.cut(5);
      assertEquals("{123456789\n12ABC}GAP{89\n123456789\n}",
            _buf.debugString());
      assertTrue(_buf.undo());
      _buf.moveTo(12);
      assertEquals("{123456789\n12ABC}GAP{3456789\n123456789\n}",
            _buf.debugString());
      assertTrue(_buf.undo());
      _buf.moveTo(12);
      // Reposition the cursor: undo doesn't make any guarantees about cursor position.
      assertEquals("{123456789\n12}GAP{3456789\n123456789\n",
            _buf.debugString());
      assertFalse(_buf.undo());
   }
}
