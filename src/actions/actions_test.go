// Copyright 2009 Mark C. Chu-Carroll
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

// File: actions_test.go
// Author: Mark Chu-Carroll <markcc@gmail.com>
// Description: Tests of edit actions.

package actions

import (
  "buf"
  "expressions"
  "gospec"
  "testing"
)

/* Create a relatively short buffer with a simple, predictable content
 * which makes it easy to test actions
 */
func CreateActionTestBuffer() buf.Buffer {
  b := buf.New(500)
  b.InsertString("1abcdedgh\n2ijklmnop\n3qrstuvwx\n")
  b.InsertString("4yz abcde\n5fghijklm\n6nopqrstu\n")
  return b
}

func CreateTestRange(b buf.Buffer, start int, end int) expressions.Range {
  return expressions.NewSimpleRange(b,
    expressions.NewStartRelativeLocation(b, start),
    expressions.NewStartRelativeLocation(b, end))
}

func TextActionSpec(c *gospec.Context) {
  buffer := CreateActionTestBuffer()
  testrange := CreateTestRange(buffer, 30, 50)
  badrange := CreateTestRange(buffer, 300, 500)
  c.Specify("An insert action", func() {
    insert := NewInsertAction("abc")
    c.Specify("should insert text at the beginning of a range", func() {
      status := insert.Execute(testrange)
      c.Then(status.GetResultCode()).Should.Equal(buf.SUCCEEDED)
      c.Then(buffer.String()).Should.Equal("1abcdedgh\n2ijklmnop\n3qrstuvwx\n" +
        "abc4yz abcde\n5fghijklm\n6nopqrstu\n")
    })
    c.Specify("should return an error code for an invalid range", func() {
      status := insert.Execute(badrange)
      c.Then(status.GetResultCode()).ShouldNot.Equal(buf.SUCCEEDED)
    })
  })
  c.Specify("An append action should ", func() {
    append := NewAppendAction("abc")
    c.Specify("add text at the end of a valid range", func() {
      status := append.Execute(testrange)
      c.Then(status.GetResultCode()).Should.Equal(buf.SUCCEEDED)
      c.Then(buffer.String()).Should.Equal("1abcdedgh\n2ijklmnop\n3qrstuvwx\n" +
        "4yz abcde\n5fghijklm\nabc6nopqrstu\n")
    })
    c.Specify("return an error if the range goes past the end of the buffer", func() {
      status := append.Execute(badrange)
      c.Then(status.GetResultCode()).ShouldNot.Equal(buf.SUCCEEDED)
    })
  })
  c.Specify("A replace action should ", func() {
    replace := NewReplaceAction("abc")
    c.Specify("replace the selection", func() {
      status := replace.Execute(testrange)
      c.Then(status.GetResultCode()).Should.Equal(buf.SUCCEEDED)
      c.Then(buffer.String()).Should.Equal("1abcdedgh\n2ijklmnop\n3qrstuvwx\n" +
        "abc6nopqrstu\n")
    })
  })
  c.Specify("A text action with a pattern should ", func() {
    ptestrange := expressions.NewSimpleRange(buffer,
      expressions.NewStartRelativeLocation(buffer, 0),
      expressions.NewEndRelativeLocation(buffer, 1))
    pat, _ := expressions.NewPatternRangeExpr("5([f-m]+)")
    patrange, _ := pat.Eval(ptestrange)
    repl := NewReplaceAction("|||$1|||")
    c.Specify("insert text with valid substititions from its match", func() {
      status := repl.Execute(patrange)
      c.Then(status.GetResultCode()).Should.Equal(buf.SUCCEEDED)
      c.Then(buffer.String()).Should.Equal("1abcdedgh\n2ijklmnop\n3qrstuvwx\n4yz abcde\n|||fghijklm|||\n6nopqrstu\n")
    })
  })
}

func DeleteActionSpec(c *gospec.Context) {
  buffer := CreateActionTestBuffer()
  testrange := CreateTestRange(buffer, 30, 50)
  badrange := CreateTestRange(buffer, 300, 500)
  c.Specify("A delete action", func() {
	delete := NewDeleteAction()
    c.Specify("on a valid range should delete the contents the range", func() {
      status := delete.Execute(testrange)
      c.Then(status.GetResultCode()).Should.Equal(buf.SUCCEEDED)
      c.Then(buffer.String()).Should.Equal("1abcdedgh\n2ijklmnop\n3qrstuvwx\n" +
        "6nopqrstu\n")
    })
    c.Specify("on an invalid range", func() {
      status := delete.Execute(badrange)
      c.Specify("fail", func() {
        c.Then(status.GetResultCode()).ShouldNot.Equal(buf.SUCCEEDED)
      })
      c.Specify("leave the buffer unmodified", func() {
	    c.Then(buffer.String()).Should.Equal(    "1abcdedgh\n2ijklmnop\n3qrstuvwx\n4yz abcde\n5fghijklm\n6nopqrstu\n")
      })
    })
  })
}

func LoopActionSpec(c *gospec.Context) {
  buffer := CreateActionTestBuffer()
  testrange := expressions.NewSimpleRange(buffer,
	expressions.NewStartRelativeLocation(buffer, 10),
    expressions.NewEndRelativeLocation(buffer, 0))
  c.Specify("A loop expression", func() {
	c.Specify("Which matches in its range", func() {
      pat, _ := expressions.NewPatternRangeExpr("qr")
      a := NewInsertAction("!!!")
      loop := NewLoopAction(pat, a)
      c.Specify("should successfully iterate over a valid range", func() {
        status := loop.Execute(testrange)
        c.Then(status.GetResultCode()).Should.Equal(buf.SUCCEEDED)
        c.Then(buffer.String()).Should.Equal("1abcdedgh\n2ijklmnop\n3!!!qrstuvwx\n4yz " +
          "abcde\n5fghijklm\n6nop!!!qrstu\n")
      })
    })
    c.Specify("Which doesn't match in its range", func() {
      pat, _ := expressions.NewPatternRangeExpr("[!@#$%^&]+")
	  a := NewInsertAction("!!!")
	  loop := NewLoopAction(pat, a)
      c.Specify("Should succeed without changing the buffer", func() {
        status := loop.Execute(testrange)
        c.Then(status.GetResultCode()).Should.Equal(buf.SUCCEEDED)
        c.Then(buffer.String()).Should.Equal("1abcdedgh\n2ijklmnop\n3qrstuvwx\n" +
            "4yz abcde\n5fghijklm\n6nopqrstu\n")
      })
    })
  })
}

func SequenceActionSpec(c *gospec.Context) {
  buffer := CreateActionTestBuffer()
  testrange := CreateTestRange(buffer, 30, 50).Normalize()
  insert := NewInsertAction("@@@")
  append := NewAppendAction("!!!")
  c.Specify("A sequence action ", func() {
    seq := NewSequenceAction()
    seq.AddAction(insert)
    seq.AddAction(append)
    c.Specify("Should execute its member actions", func() {
	  status := seq.Execute(testrange)
	  c.Then(status.GetResultCode()).Should.Equal(buf.SUCCEEDED)
	  c.Then(buffer.String()).Should.Equal("1abcdedgh\n2ijklmnop\n3qrstuvwx\n" +
          "@@@4yz abcde\n5fghijklm\n!!!6nopqrstu\n")
    })
  })
}

func TestActionSpecs(t *testing.T) {
  r := gospec.NewRunner()
  r.AddSpec("Text action spec", TextActionSpec)
  r.AddSpec("Delete action spec", DeleteActionSpec)
  r.AddSpec("Loop action spec", LoopActionSpec)
  r.AddSpec("Sequence action spec", SequenceActionSpec)
  gospec.MainGoTest(r, t)
}
