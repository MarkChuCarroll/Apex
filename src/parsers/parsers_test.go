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

// File: parsers_test.go
// Author: Mark Chu-Carroll <markcc@gmail.com>
// Description: tests for the parser combinator library.

package parsers

import (
  "container/vector"
  "testing"
  "fmt"
)

// Helper function for parsing.
// Given:
// - A parser
// - an input string
// - a position in the input string
// - the length of an expected match
// - the expected result value of a parse
//
// Runs the parser on the input starting at the specified
// position, and verifies that it matched the correct substring,
// and returns the correct value.
func ExpectSuccessfulParse(t *testing.T, p Parser, input string, pos int, exp_match int, exp_result ParseValue) bool {
  in := MakeStringParserInput(input)
  status, match_len, result := p.Parse(in, pos)
  if !status {
    t.Error("Parse failed")
    return false
  }
  if match_len != exp_match {
    t.Error(fmt.Sprintf("Expected match of length %v, but found %v",
      exp_match, match_len))
    return false
  }
  if result != exp_result {
    t.Error(fmt.Sprintf("Parse result did not match; "+
      "expected '%#v' but found '%#v'",
      exp_result, result))
    return false
  }
  return true
}

// Basically like ExpectSuccessful parse, except that it expects
// the parse to fail.
func ExpectFailedParse(t *testing.T, p Parser, input string, pos int) {
  in := MakeStringParserInput(input)
  status, match_len, _ := p.Parse(in, pos)
  if status {
    t.Error(fmt.Sprintf("Expected parse failure on input %v", in))
  }
  if match_len >= 0 {
    t.Error(fmt.Sprintf("Failed parse should have matchlen=-1, but found %v",
      match_len))
  }
}

func TestBasicCharSet(t *testing.T) {
  csp := CharSet("abc")

  ExpectSuccessfulParse(t, csp, "abcd", 0, 1, uint8('a'))
  ExpectSuccessfulParse(t, csp, "abcd", 1, 1, uint8('b'))
  ExpectSuccessfulParse(t, csp, "abcd", 2, 1, uint8('c'))
  ExpectFailedParse(t, csp, "abcd", 3)
}

func TestMany(t *testing.T) {
  csp := CharSet("abc")
  mp := Word(Many(csp, 2))
  ExpectSuccessfulParse(t, mp, "abcd", 0, 3, "abc")

  // should fail: match is only length 1, but min is 2.
  ExpectFailedParse(t, mp, "addu.", 0)

  mp2 := Word(Many(csp, 0))
  ExpectSuccessfulParse(t, mp2, "addu", 0, 1, "a")
  ExpectSuccessfulParse(t, mp2, "dde", 0, 0, "")
}

func TestAlts(t *testing.T) {
  as := CharSet("a")
  bs := CharSet("b")
  cs := CharSet("c")
  abcs := Word(Many(Alt([]Parser{as, bs, cs}), 1))
  ExpectSuccessfulParse(t, abcs, "abcd", 0, 3, "abc")
}

func TestSimpleSequence(t *testing.T) {
  a := CharSet("a")
  b := CharSet("b")
  c := CharSet("c")
  ab := Alt([]Parser{a, b})
  word := Word(Seq([]Parser{a, b, ab, ab, c}))

  ExpectSuccessfulParse(t, word, "abbac", 0, 5, "abbac")
  ExpectSuccessfulParse(t, word, "abaac", 0, 5, "abaac")
  ExpectFailedParse(t, word, "abcac", 0)
}

func TestManySequence(t *testing.T) {
  as := Word(Many(CharSet("a"), 1))
  bs := Word(Many(CharSet("b"), 0))
  cs := Word(Many(CharSet("c"), 1))
  abcs := Word(Seq([]Parser{as, bs, cs}))
  ExpectSuccessfulParse(t, abcs, "abc", 0, 3, "abc")
  ExpectSuccessfulParse(t, abcs, "ac", 0, 2, "ac")
  ExpectSuccessfulParse(t, abcs, "aabc", 0, 4, "aabc")
  ExpectSuccessfulParse(t, abcs, "abbc", 0, 4, "abbc")
  ExpectSuccessfulParse(t, abcs, "aabcc", 0, 5, "aabcc")
  ExpectSuccessfulParse(t, abcs, "aaaccc", 0, 6, "aaaccc")
  ExpectSuccessfulParse(t, abcs, "aaabbbbccc", 0, 10, "aaabbbbccc")
}

func TestLetterBuiltin(t *testing.T) {
  ExpectSuccessfulParse(t, Word(Letters), "abCDef", 0, 6, "abCDef")
  ExpectSuccessfulParse(t, Word(Letters), "abCD7ef", 0, 4, "abCD")
  ExpectFailedParse(t, Letters, "8aoeu", 0)
}

func TestSpaceBuiltin(t *testing.T) {
  ExpectSuccessfulParse(t, Word(Spaces), "    aoeu", 0, 4, "    ")
  ExpectSuccessfulParse(t, Word(Spaces), "  \n	  aoeu", 0, 6, "  \n	  ")
  ExpectSuccessfulParse(t, Word(Spaces), "aoeu", 0, 0, "")
}

func TestTokenBuiltin(t *testing.T) {
  ExpectSuccessfulParse(t, Token(Letters), "   aoeu  ", 0, 7, "aoeu")
  ExpectSuccessfulParse(t, Token(Letters), "   \n  \tAbC", 0, 10, "AbC")
  ExpectSuccessfulParse(t, Token(Letters), "qwerty", 0, 6, "qwerty")
  ExpectFailedParse(t, Token(Letters), "  123", 0)
  ExpectFailedParse(t, Token(Letters), "123", 0)
}


// Getting more interesting: Test matching parens using RefParsers.
//
// A := many B
// B := ( A )
func TestParensParser(t *testing.T) {
  NumParser := CharSet("123")
  OpenParen := CharSet("(")
  CloseParen := CharSet(")")
  b := MakeRef()
  a := Alt([]Parser{Many(b, 1), NumParser})
  b.SetTarget(Seq([]Parser{OpenParen, a, CloseParen}))
  in := MakeStringParserInput("((1)(2))")
  status, _, _ := a.Parse(in, 0)
  if !status {
    t.Error("Parse failed")
  }
}

// And even more interesting: Test sexpression syntax parser.
func TestLispParser(t *testing.T) {
  a := Cons(Sym("a"), Cons(Cons(Sym("b"), Cons(Sym("c"), nil)), Cons(Sym("d"), nil)))
  if a.String() != "(a . ((b . (c . nil)) . (d . nil)))" {
    t.Error("Cons lists aren't working; no point trying to do a real test")
    return
  }
  SymParser := Token(Letters)
  Lp := Fixed(Token(CharSet("(")), "(")
  Rp := Fixed(Token(CharSet(")")), ")")

  SexprParserRef := MakeRef()
  ManySexprs := Concat(Many(SexprParserRef, 1), " ")
  ListParser := Concat(Seq([]Parser{Lp, ManySexprs, Rp}), "")
  Sexpr := Alt([]Parser{ListParser, SymParser})
  SexprParserRef.SetTarget(Sexpr)

  ExpectSuccessfulParse(t, Sexpr, "(a (b  c) d e (f (g   i) ) )", 0, 28,
    "(a (b c) d e (f (g i)))")
}

// Now, we get to the really good one: parse sexpressions, returning as
// a result a cons-cell based list structure.

// Start by defining the cons-list data structures.

type SList struct {
  car SVal
  cdr *SList
}

func (list *SList) Car() SVal { return list.car }

func (list *SList) Cdr() *SList { return list.cdr }

func (list *SList) String() (result string) {
  if list.cdr == nil {
    result = "(" + list.car.String() + " . nil)"
  } else {
    result = "(" + list.car.String() + " . " + list.cdr.String() + ")"
  }
  return
}

func Cons(h SVal, t *SList) *SList {
  result := new(SList)
  result.car = h
  result.cdr = t
  return result
}

func (list *SList) IsAtom() bool { return false }

func (list *SList) IsList() bool { return true }

func (list *SList) AsList() *SList { return list }

type Sym string

func (s Sym) IsAtom() bool { return true }

func (s Sym) IsList() bool { return false }

func (s Sym) AsList() *SList { return nil }

func (s Sym) String() string { return string(s) }


type SVal interface {
  IsAtom() bool
  IsList() bool
  AsList() *SList
  String() string
}

func VectorToSExpr(v *vector.Vector, pos int) (result *SList) {
  if pos == v.Len() {
    result = nil
  } else {
    car := SVal(nil)
    switch typed_el := v.At(pos).(type) {
    case *vector.Vector:
      car = VectorToSExpr(typed_el, 0)
    case SVal:
      car = typed_el
    default:
      car = Sym(fmt.Sprintf("%v", typed_el))
    }
    result = Cons(car, VectorToSExpr(v, pos+1))
  }
  return
}

func TestVectorToSExpr(t *testing.T) {
  v := new(vector.Vector)
  v.Push(Sym("ab"))
  w := new(vector.Vector)
  w.Push(Sym("c"))
  w.Push(Sym("d"))
  v.Push(w)
  v.Push(Sym("e"))
  s := VectorToSExpr(v, 0)
  if s.String() != "(ab . ((c . (d . nil)) . (e . nil)))" {
    t.Error(fmt.Sprintf("Expected  '(ab . ((c . (d . nil)) . (e . nil)))'"+
      " but found '%v'",
      s))
  }
}

type SymAction struct{}

func (self *SymAction) Execute(val ParseValue) ParseValue {
  return Sym(val.(string))
}

func Symbol(p Parser) Parser { return Action(p, new(SymAction)) }

type VectAction struct{}

func (self *VectAction) Execute(val ParseValue) ParseValue {
  return VectorToSExpr(val.(*vector.Vector), 0)
}


func TestSexprParseAndBuild(t *testing.T) {
  a := Cons(Sym("a"), Cons(Cons(Sym("b"), Cons(Sym("c"), nil)), Cons(Sym("d"),
    nil)))
  if a.String() != "(a . ((b . (c . nil)) . (d . nil)))" {
    t.Error("Cons lists aren't working; no point trying to do a real test")
    return
  }
  // The tokens
  SymParser := Symbol(Token(Letters))
  Lp := Fixed(Token(CharSet("(")), "(")
  Rp := Fixed(Token(CharSet(")")), ")")

  // ManySexprs : ( Sexpr )+ { action: generate List Sexpr for the list
  //                             of expressions }
  //  List : Lp ManySexprs Rp { action: return $2 }
  // Sexpr : List | Sym
  SexprParserRef := MakeRef()
  ManySexprs := Action(Many(SexprParserRef, 1), new(VectAction))
  ListParser := Second(Seq([]Parser{Lp, ManySexprs, Rp}))
  Sexpr := Alt([]Parser{ListParser, SymParser})
  SexprParserRef.SetTarget(Sexpr)

  in := MakeStringParserInput("(a (b  c) d e (f (g   i) ) )")
  success, _, val := Sexpr.Parse(in, 0)
  if !success {
    t.Error("Parse failed")
  }
  pval := val.(SVal)
  if pval.String() != "(a . ((b . (c . nil)) . (d . (e . ((f . ((g . "+
    "(i . nil)) . nil)) . nil)))))" {
    t.Error(fmt.Sprintf("Invalid parse result: expected: "+
      "'(a . ((b . (c . nil)) . (d . (e . ((f . ((g . "+
      "(i . nil)) . nil)) . nil)))))', but found %v",
      pval.String()))

  }
}
