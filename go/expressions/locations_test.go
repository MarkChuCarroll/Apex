package expressions

import (
  "buf"
  "gospec"
  "testing"
)

/* Create a buffer with a simple, predictable content
 * which makes it easy to test ranges.
 */
func CreateLocationTestBuffer() buf.Buffer {
  b := buf.New(500)
  b.InsertString("1abcdedgh\n2ijklmnop\n3qrstuvwx\n")
  b.InsertString("4yz abcde\n5fghijklm\n6nopqrstu\n")
  b.InsertString("7vwxyz ab\n8cdefghij\n9klmnopqr\n")
  b.InsertString("10stuvwxy\n11z abcde\n12fghijkl\n")
  b.InsertString("13mnopqrs\n14tuvwxyz\n15 abcdef\n")
  b.InsertString("16ghijklm\n17opqrstu\n18vwxzy a\n")
  b.InsertString("19bcdefgh\n20ijklmno\n21pqrstuv\n")
  b.InsertString("22wxyz ab\n23cdefghi\n24jklmnop\n")
  b.InsertString("25qrstuvw\n26yz abcd\n27efghijk\n")
  b.InsertString("28lmnopqr\n29stuvwxy\n30z abcde\n")
  return b
}

func CreateTestRange(b buf.Buffer, start, end int) Range {
  return NewSimpleRange(b, NewStartRelativeLocation(b, start),
    NewStartRelativeLocation(b, end)).Normalize()
}

func BasicLocationSpec(c gospec.Context) {
  buffer := CreateLocationTestBuffer()
  testrange := CreateTestRange(buffer, 30, 70)

  c.Specify("The basic test buffer", func() {
    c.Specify("Should be addressable by simple positions", func() {
      loc_expr := NewCharLocExpr(18)

      c.Specify("Location 18 should evaluate to a valid location", func() {
        loc, status := loc_expr.Eval(testrange)
        c.Expect(int(status.GetResultCode()), gospec.Equals, int(buf.SUCCEEDED))
        status = Validate(loc.GetStart())
        c.Expect(status.GetResultCode(), gospec.Equals, buf.SUCCEEDED)
        c.Expect(loc.GetStart().GetAbsolute(), gospec.Equals, 18)
      })
    })
    c.Specify("Should be addressable by end-relative location expressions like -28", func() {
      loc_expr := NewCharLocExpr(-28)
      loc, status := loc_expr.Eval(testrange)

      c.Specify("Position -28 should evaluate without error", func() {
        c.Expect(status.GetResultCode(), gospec.Equals, buf.SUCCEEDED)
        status = Validate(loc.GetStart())
        c.Expect(status.GetResultCode(), gospec.Equals, buf.SUCCEEDED)
        c.Expect(loc.GetStart().GetAbsolute(), gospec.Equals, buffer.Length() - 28)
      })
    })
    c.Specify("Should be adressable by a expression like line(3)", func() {
      loc_expr := NewLineLocExpr(3)
      loc, status := loc_expr.Eval(testrange)

      c.Specify("line 3 should evaluate without error", func() {
        // Without the "int" conversion here, the comparison fails.
        c.Expect(status.GetResultCode(), gospec.Equals, buf.SUCCEEDED)
      })

      c.Specify("Line three should evaluate to a valid location", func() {
        loc_status := Validate(loc.GetStart())
        c.Expect(loc_status.GetResultCode(), gospec.Equals, buf.SUCCEEDED)
      })

      c.Specify("Line three should resolve to the beginning of the third line", func() {
        pos := loc.GetStart().GetAbsolute()
        c.Expect(pos, gospec.Equals, 20)
      })

    })
  })
}

/*
func RelativePointSpec(c *gospec.Context) {
  buffer := CreateLocationTestBuffer()
  testrange := NewSimpleRange(buffer, NewStartRelativeLocation(buffer, 30), NewStartRelativeLocation(buffer, 70))

  c.Specify("Points should be specifyable by a positive offset relative to selection start", func() {
    positive_point := NewStartRelCharLocExpr(12)
    ppoint, ppok := positive_point.Eval(testrange)
    c.Then(ppok.GetResultCode()).Should.Equal(buf.SUCCEEDED)
    c.Then(ppoint.GetAbsolute()).Should.Equal(42)
  })

  c.Specify("Points should be specifyable by a negative offset relative to selection start", func() {
    negative_point := NewStartRelCharLocExpr(-12)
    npoint, npok := negative_point.Eval(testrange)
    c.Then(npok.GetResultCode()).Should.Equal(buf.SUCCEEDED)
    c.Then(npoint.GetAbsolute()).Should.Equal(18)
  })

  c.Specify("Points should be specifyable by a positive offset relative to selection end", func() {
    positive_point := NewEndRelCharLocExpr(12)
    ppoint, ppok := positive_point.Eval(testrange)
    c.Then(ppok.GetResultCode()).Should.Equal(buf.SUCCEEDED)
    c.Then(ppoint.GetAbsolute()).Should.Equal(82)
  })

  c.Specify("Points should be specifyable by a negative offset relative to selection end", func() {
    positive_point := NewEndRelCharLocExpr(-12)
    ppoint, ppok := positive_point.Eval(testrange)
    c.Then(ppok.GetResultCode()).Should.Equal(buf.SUCCEEDED)
    c.Then(ppoint.GetAbsolute()).Should.Equal(58)
  })
}

func LineLocationSpec(c *gospec.Context) {
  buffer := CreateLocationTestBuffer()
  testrange := NewSimpleRange(buffer, NewStartRelativeLocation(buffer, 30), NewStartRelativeLocation(buffer, 70))
  c.Specify("Locations should be specifyable using line numbers", func() {
    line_point := NewLineLocExpr(4)
    lpoint, lpok := line_point.Eval(testrange)
    c.Then(lpok.GetResultCode()).Should.Equal(buf.SUCCEEDED)
    c.Then(lpoint.GetAbsolute()).Should.Equal(30)
  })

  c.Specify("Locations should be specifyable using lines from end", func() {
    line_point := NewLineLocExpr(-4)
    lpoint, lpok := line_point.Eval(testrange)
    c.Then(lpok.GetResultCode()).Should.Equal(buf.SUCCEEDED)
    c.Then(lpoint.GetAbsolute()).Should.Equal(260)
  })
}

func LineLocationInSelectionSpec(c *gospec.Context) {
  buffer := CreateLocationTestBuffer()
  testrange := CreateTestRange(buffer, 32, 76)
  c.Specify("Locations can be specified using line-offsets from the start of a selection", func() {
    line_point := NewStartRelLineLocExpr(2)
    lp, lpok := line_point.Eval(testrange)
    c.Then(lpok.GetResultCode()).Should.Equal(buf.SUCCEEDED)
    c.Then(lp.GetAbsolute()).Should.Equal(50)
  })
}

func TwoPointRangeExprSpec(c *gospec.Context) {
  buffer := CreateLocationTestBuffer()
  testrange := CreateTestRange(buffer, 32, 76)
  c.Specify("Ranges specify a range of text between two locations", func() {
    start := NewStartRelCharLocExpr(5)
    end := NewStartRelCharLocExpr(15)
    range_expr := NewTwoPointRangeExpr(start, end)
    r, status := range_expr.Eval(testrange)
    c.Specify("A valid range should evaluate without any errors", func() {
      c.Then(status.GetResultCode()).Should.Equal(buf.SUCCEEDED)
      c.Then(r.GetStart().GetAbsolute()).Should.Equal(37)
      c.Then(r.GetEnd().GetAbsolute()).Should.Equal(47)
    })
    c.Specify("A valid range should be able to get its contents", func() {
      content, cok := r.GetContents()
      c.Then(cok.GetResultCode()).Should.Equal(buf.SUCCEEDED)
      c.Then(string(content)).Should.Equal("de\n5fghijk")
    })
  })
}

func PatternRangeExprSpec(c *gospec.Context) {
  buffer := CreateLocationTestBuffer()
  pattern_string := "[abcdefghijklmnopqrstuvwxyz]+\n([0123456789]+)"
  testrange := CreateTestRange(buffer, 32, 76)
  c.Specify("A pattern range", func() {
    pattern_expr, ok := NewPatternRangeExpr(pattern_string)
    c.Specify("created from a valid pattern is valid", func() { c.Then(ok).Should.Equal(nil) })
    prange, prok := pattern_expr.Eval(testrange)
    c.Specify("should find text that matches its pattern", func() {
      c.Then(prok.GetResultCode()).Should.Equal(buf.SUCCEEDED)
      c.Then(prange.GetStart().GetAbsolute()).Should.Equal(34)
      c.Then(prange.GetEnd().GetAbsolute()).Should.Equal(41)
    })
    c.Specify("should be able to get its contents", func() {
      content, cok := prange.GetContents()
      c.Then(cok.GetResultCode()).Should.Equal(buf.SUCCEEDED)
      c.Then(string(content)).Should.Equal("abcde\n5")
    })
    patrange := prange.(*PatternRange)
    c.Specify("should be able to access sub-matches", func() {
      c.Then(patrange.GetNumberOfMatches()).Should.Equal(1)
      c.Then(string(patrange.GetMatch(1))).Should.Equal("5")
    })
    c.Specify("should return nil if you try to access an invalid sub-match", func() { c.Then(string(patrange.GetMatch(8))).Should.Equal("") })
    c.Specify("should be able to use pattern matches in string interpolation", func() {
      templ := "hello $1there"
      chars, status := patrange.InstantiateTemplate(templ)
      c.Then(status.GetResultCode()).Should.Equal(buf.SUCCEEDED)
      c.Then(string(chars)).Should.Equal("hello 5there")
    })
  })
}

func PatternRangeTemplateSpec(c *gospec.Context) {
  b := buf.New(500)

  b.InsertString("12343aoe 89put 4945#@$%$a hs")
  num := "[0123456789]+"
  letter := "[abcdefghijklmnopqrstuvwxyz]+"
  whitespace := " +"
  sym := "[#@$%&]+"
  //    L W (N) W N (S)
  pattern_string := fmt.Sprintf("%v%v(%v)%v%v%v(%v)",
    letter, whitespace, num, letter, whitespace, num, sym)

  testrange := CreateTestRange(b, 0, b.Length())
  c.Specify("Pattern interpolation should", func() {
    pattern_expr, pat_ok := NewPatternRangeExpr(pattern_string)
    c.Specify("be based on a valid pattern", func() { c.Then(pat_ok).Should.Equal(nil) })
    prange, prok := pattern_expr.Eval(testrange)
    c.Specify("work on a valid match", func() { c.Then(prok.GetResultCode()).Should.Equal(buf.SUCCEEDED) })
    patrange := prange.(*PatternRange)
    c.Specify("fill in positions in a pattern string", func() {
      subst, subst_ok := patrange.InstantiateTemplate("($1)($2)")
      c.Then(subst_ok.GetResultCode()).Should.Equal(buf.SUCCEEDED)
      c.Then(string(subst)).Should.Equal("(89)(#@$%$)")
    })
    c.Specify("Return an error code for an invalid pattern string", func() {
      _, subst_ok := patrange.InstantiateTemplate("($1)($2)($9)")
      c.Then(subst_ok.GetResultCode()).Should.Equal(buf.INVALID_REPLACEMENT)
    })
  })
}

func AfterRangeSpec(c *gospec.Context) {
  buffer := CreateLocationTestBuffer()
  testrange := CreateTestRange(buffer, 30, 70)
  badtestrange := CreateTestRange(buffer, 50, 70)
  before, _ := NewPatternRangeExpr("5[abcdefghijklmnopqrstuvwxyz ]+")
  after, _ := NewPatternRangeExpr("[qrstuvw]+")
  c.Specify("An after range expression", func() {
    after_expr := NewAfterRangeExpr(before, after)
    c.Specify("Should match its after range after its before", func() {
      result, status := after_expr.Eval(testrange)
      c.Then(status.GetResultCode()).Should.Equal(buf.SUCCEEDED)
      content, cok := result.GetContents()
      c.Then(cok.GetResultCode()).Should.Equal(buf.SUCCEEDED)
      c.Then(string(content)).Should.Equal("qrstu")
    })
    c.Specify("Should fail when its before doesn't match anything", func() {
      _, status := after_expr.Eval(badtestrange)
      c.Then(status.GetResultCode()).ShouldNot.Equal(buf.SUCCEEDED)
    })
  })
}
*/
/*
  public void testAfterRangeAction() throws Exception {
    PatternRangeExpression before = new PatternRangeExpression("5");
    PatternRangeExpression after = new PatternRangeExpression("t");
    // first t after a 5
    AfterRangeExpression arange = new AfterRangeExpression(before, after);
    EditAction act = new InsertAction("$$$");
    Range range = arange.eval(_selection);
    act.execute(range);
    Assert.assertEquals("1abcdefgh\n2ijklmnop\n3qrstuvwx\n4yz " +
        "abcde\n5fghijklm\n6nopqrs$$$tu\n" +
        "7vwxyz ab\n8cdefghij\n9klmnopqr\n10stuvwxy\n", _buffer.toString());   
 } 
}
*/

// Should add some failure tests.

func TestSpecs(t *testing.T) {
  r := gospec.NewRunner()
  r.AddSpec(BasicLocationSpec)
//  r.AddSpec("Selection-relative Location Spec", RelativePointSpec)
//  r.AddSpec("Line Location Spec", LineLocationSpec)
//  r.AddSpec("Line Location in Selection Spec", LineLocationInSelectionSpec)
//  r.AddSpec("Two Point Range spec", TwoPointRangeExprSpec)
//  r.AddSpec("Pattern range expression spec", PatternRangeExprSpec)
//  r.AddSpec("Pattern interpolation spec", PatternRangeTemplateSpec)
//  r.AddSpec("After range spec", AfterRangeSpec)
  gospec.MainGoTest(r, t)
}
