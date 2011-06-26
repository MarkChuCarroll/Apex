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

// File: parsers.go
// Author: Mark Chu-Carroll <markcc@gmail.com>
// Description: A library of Parser combinators.

package parsers

import (
  "container/vector"
  "fmt"
)

///////////////////////////////////////////////////////////
// Set up basic types: result values, inputs, etc.

// ParseValue represents the type of values returned by
// successful parses.
type ParseValue interface{}

// ParserInput represents an input source readable by a
// parser.
type ParserInput interface {
  // Get the character at an arbitrary position in the
  // input source.
  CharAt(i int) uint8

  // Get the number of characters in the input source.
  Size() int
}

// A Parser is an object which parses input sources. The
// framework of parser combinators provides a very general,
// backtracking parser.
type Parser interface {

  // Run the parser on an input source, starting with
  // the character at a specified position. If the parse
  // succeeds, it returns "true" as the status, the
  // number of characters matched by the parser as the match_len,
  // and an arbitrary parser-specified, return value as the result.
  // If the parse fails, then it returns false as the status, -1 as
  // match_len, and nil as the result.
  Parse(in ParserInput, pos int) (status bool, match_len int, result ParseValue)
}

// A ParseAction is a callback which is executed on
// the results of a successful parse to transform
// it into some other value. It provides parser combinators
// with the ability to execute semantic actions.
type ParseAction interface {
  // Execute the callback.
  Execute(val ParseValue) ParseValue
}


// StringParserInput is an implementation of a ParserInput
// which uses a string value as a backing store.
type StringParserInput struct {
  value string
}

func (in *StringParserInput) Size() int { return len(in.value) }

func (in *StringParserInput) CharAt(i int) (result uint8) {
  if i >= in.Size() {
    result = 0
  } else {
    result = in.value[i]
  }
  return
}

func MakeStringParserInput(s string) *StringParserInput {
  return &StringParserInput{s}
}

// CharSetParser parses a single character from among a
// specified set.
type CharSetParser struct {
  chars string
}

func (self *CharSetParser) Parse(in ParserInput, pos int) (status bool, match_len int, result ParseValue) {
  status = false
  match_len = -1
  for c := range self.chars {
    if self.chars[c] == in.CharAt(pos) {
      result = self.chars[c]
      match_len = 1
      status = true
    }
  }
  return
}

// Create a CharSetParser which accepts any one character
// from a specified string.
func CharSet(s string) *CharSetParser { return &CharSetParser{s} }

// ManyParser is a parser that parses a repeated syntax
// element. It succeeds if the sub-parser succeeds at least
// a specified minimum number of times. Returns a list
// of the results of the sub-parses.
type ManyParser struct {
  min    int
  parser Parser
}

// Create a ManyParser which matches min or more
// repetitions of the sequence parsed by p.
func Many(p Parser, min int) *ManyParser {
  result := &ManyParser{min, p}
  return result
}

func (self *ManyParser) Parse(in ParserInput, pos int) (status bool, match_len int, results ParseValue) {
  status = false
  curPos := pos
  numMatches := 0
  stepResults := new(vector.Vector)
  stepSuccess, stepLen, stepResult := self.parser.Parse(in, curPos)
  for stepSuccess {
    numMatches++
    curPos = curPos + stepLen
    stepResults.Push(stepResult)
    stepSuccess, stepLen, stepResult = self.parser.Parse(in, curPos)
  }
  if numMatches < self.min {
    stepSuccess = false
    match_len = -1
    results = nil
    return
  }
  status = true
  results = stepResults
  match_len = curPos - pos
  return
}

// AltParser is a parser that parses one of a list of
// alternatives. Succeeds if any of its sub-parsers
// succeeds. Returns the results from the first of the
// alternatives that is successful.
type AltParser struct {
  alts []Parser
}

// Creates an AltParser for a list of alternatives.
func Alt(alts []Parser) (result *AltParser) {
  result = new(AltParser)
  result.alts = alts
  return
}

func (self *AltParser) Parse(in ParserInput, pos int) (status bool, match_len int, results ParseValue) {
  status = false
  match_len = -1
  results = nil
  for i := range (self.alts) {
    parser := self.alts[i]
    s, l, r := parser.Parse(in, pos)
    if s {
      status = true
      match_len = l
      results = r
      return
    }
  }
  return
}

// OptParser is a parser that a parser that parses an
// optional element. OptParser of a parser p is equivalent to
// an AltParser of p and an empty-string-matcher. OptParsers
// always succeed; match length is zero, and result value
// is nil if the opt was empty.
type OptParser struct {
  parser Parser
}

func Opt(p Parser) (result *OptParser) {
  result = new(OptParser)
  result.parser = p
  return
}

func (self *OptParser) Parse(in ParserInput, pos int) (status bool, match_len int, results ParseValue) {
  status = true
  match_len = -1
  results = nil
  s, l, r := self.parser.Parse(in, pos)
  if s {
    match_len = l
    results = r
  }
  return
}

// A SeqParser is a parser that parses an collection of
// elements in sequence. It succeeds if and only if *all*
// of its sub-parsers succeed. Returns a list of the
// parse results of its sub-parsers.
type SeqParser struct {
  parsers []Parser
}

func Seq(parsers []Parser) (result *SeqParser) {
  result = new(SeqParser)
  result.parsers = parsers
  return result
}

func (self *SeqParser) Parse(in ParserInput, pos int) (status bool, match_len int, results ParseValue) {
  result_values := new(vector.Vector)
  curPos := pos
  for p := range self.parsers {
    parser := self.parsers[p]
    s, l, r := parser.Parse(in, curPos)
    if s {
      curPos += l
      result_values.Push(r)
    } else {
      status = false
      match_len = -1
      results = nil
      return
    }
  }
  status = true
  match_len = (curPos - pos)
  results = result_values
  return
}

// ActionParser is a wrapper for executing semantic actions
// during a parse. It invokes a sub-parser; then, if the
// sub-parser succeeds, in invokes a callback on the
// sub-parser result, generating a new value which is
// used as the result of the ActionParser.
// Semantic actions should *not* perform global side-effects;
// it is possible that the semantic action is taking place as
// part of an attempted parse that will fail, and will be erased
// by backtracking.
type ActionParser struct {
  parser Parser
  action ParseAction
}

func Action(p Parser, a ParseAction) (result *ActionParser) {
  result = new(ActionParser)
  result.parser = p
  result.action = a
  return
}

func (self *ActionParser) Parse(in ParserInput, pos int) (status bool, match_len int, result ParseValue) {
  status, match_len, result = self.parser.Parse(in, pos)
  if status {
    result = self.action.Execute(result)
  }
  return
}

// A tracing parser is a parser which emits trace
// messages before and after parsing an element.
type TracingParser struct {
  parser Parser
  name   string
}

func (self *TracingParser) Parse(in ParserInput, pos int) (status bool, match_len int, result ParseValue) {
  fmt.Printf("Entering parser: %s\n", self.name)
  status, match_len, result = self.parser.Parse(in, pos)
  if status {
    fmt.Printf("Exiting parser(success): %s\n", self.name)
  } else {
    fmt.Printf("Exiting parser(failed): %s\n", self.name)
  }
  return
}

func Trace(p Parser, name string) Parser { return &TracingParser{p, name} }

// RefParser is a wrapper for producing mutually recursive grammar
// rules. You can first create an empty RefParser as a placeholder,
// and then update it.
//
// For example, given grammar rules:
// A : '(' B* ')' | 'X'
// B : '[' A* ']'
//
// you could:
//    B := MakeRef();
//    A := Alt([]*Parser{ Seq([]*Parser { OpenParen, B, CloseParen }), X })
//    B.SetTarget(Seq([] { OpenBracket, A, CloseBracket }))
type RefParser struct {
  target Parser
}

func MakeRef() *RefParser { return new(RefParser) }

func (self *RefParser) SetTarget(p Parser) { self.target = p }

func (self *RefParser) Parse(in ParserInput, pos int) (status bool, match_len int, result ParseValue) {
  status, match_len, result = self.target.Parse(in, pos)
  return
}

/////////////////////////////////////////////////////////////
// Some built-in basics: standard parsers and actions

type DiscardAction struct{}

func (self *DiscardAction) Execute(val ParseValue) ParseValue {
  return nil
}

type FixedValueAction struct {
  val ParseValue
}

func (self *FixedValueAction) Execute(val ParseValue) ParseValue {
  return self.val
}

// Parser actions for doing string concatenations of "Many" sub-parser results.
type ConcatAction struct {
  sep string
}

func (self *ConcatAction) Execute(val ParseValue) ParseValue {
  v, ok := val.(*vector.Vector)
  if !ok {
    return nil
  }
  s := ""
  for i := 0; i < v.Len(); i++ {
    if i != 0 {
      s = s + self.sep
    }
    el := v.At(i)
    switch typed_el := el.(type) {
    case byte:
      s = s + string(typed_el)
    case string:
      s = s + typed_el
      //    case Sym: s = s + string(typed_el)
    case fmt.Stringer:
      s = s + typed_el.String()
    default:
      s = fmt.Sprintf("%v %v", s, el)
    }
  }
  return s
}

type NthValueAction struct {
  num int
}


func (self *NthValueAction) Execute(val ParseValue) ParseValue {
  v, ok := val.(*vector.Vector)
  if !ok {
    return nil
  }
  if v.Len() < self.num {
    return nil
  }
  return v.At(self.num - 1)
}

// Take a parser which returns a value, and transform
// it to a parser that retuns nil
func Discard(p Parser) Parser { return Action(p, new(DiscardAction)) }

// Some vector-related combinators: take parsers which produce
// a vector of values, and do something to the vector in an action.

// Merge the text of a sequence of parse results, concatenating
// the results without any separator.
func Word(p Parser) Parser { return Action(p, &ConcatAction{""}) }

// Merge the text of a sequence of parse results, concatenating
// with a specified separator.
func Concat(p Parser, sep string) Parser { return Action(p, &ConcatAction{sep}) }

// Change the result of a parser from a vector to
// the Nth value of that vector. Useful for cases
// where you have things like a set of pure syntax tokens
// around a single semantic value.
func Nth(p Parser, n int) Parser { return Action(p, &NthValueAction{n}) }

func Second(p Parser) Parser { return Nth(p, 2) }

// Misc transformers

// Create a parser which always returns a fixed result
// value.
func Fixed(p Parser, v ParseValue) Parser { return Action(p, &FixedValueAction{v}) }

// generate a parser which strips all leading spaces, and then runs a particular parser
func Token(p Parser) Parser { return Second(Seq([]Parser{Spaces, Word(p)})) }


var (
  Space    = CharSet(" \t\n")
  Spaces   = Many(Space, 0)
  Letter   = CharSet("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_")
  Letters  = Many(Letter, 1)
  Digit    = CharSet("0123456789")
  AlphaNum = Alt([]Parser{Letter, Digit})
)
