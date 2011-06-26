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

// File: lex.go
// Author: Mark Chu-Carroll <markcc@gmail.com>
// Description: The scanner for the Apex language
package language

import (
  _ "io"
  "fmt"
)

type ScannerInput interface {
  Peek() uint8
  Current() uint8
  Advance() bool
  Line() int
}

type StringScannerInput struct {
  str string
  pos int
  curtok []uint8
  line int
}

func NewStringInput(s string) (result *StringScannerInput) {
  result = &StringScannerInput{s, 0, make([]uint8, 0, 32), 1}
  result.curtok = append(result.curtok, result.Current())
  return
}

func (self *StringScannerInput) Peek() uint8 {
  if (self.pos + 1) >= len(self.str) {
    return 0	
  }
  return self.str[self.pos + 1]
}

func (self *StringScannerInput) Advance() bool {
  self.pos++
  if self.pos >= len(self.str) {
    return false
  }
  c := self.Current()
  if c == 0 {
	return false
  } else {
    if c == '\n' {
      self.line++	
    }
    return true
  }
  return false
}

func (self *StringScannerInput) Current() uint8 {
  if self.pos >= len(self.str) {
    return 0
  }
  return self.str[self.pos]
}

func (self *StringScannerInput) Line() int {
  return self.line	
}


// Scanning is messy. There are multiple modes contexts, and the same
// thing is treated differently in different modes. The mode
// trigger is quoting. But a string following certain commands is
// automatically quoted, no matter what: the first character after the
// command is the quote character, and standard quotes don't mean anything.

// Commands that auto-quote:
// "_" (regexp match string, should only occur inside of a regexp)
// "/" (regexp charset string; should only occur in regexp)
// "+" (regexp search forward) (does this really need to be this way? The
//      regexp has its own quote.)
// "-" (regexp search backward)
// "i" (insert text)
// "a" (append text)
// "r" (replace text)
// "<" (command shell invocation)
// "<<"
// ">"
// ">>"
// "w" (write file)
// "o" (open file)
type ScannerMode int

const (
  MODE_NORMAL ScannerMode = iota
  MODE_QUOTED
)

type Scanner struct {
  In   ScannerInput
  mode ScannerMode
  error  string
}

type Token struct {
  Str    string
  Strval string
  Type   int
  Line   int
}

type IScanner interface {
  CurrentToken() *Token
  Advance()
  PushBack(t *Token)
  
  SetError(err string)
  GetLastError() string
}

func (self *Scanner) SetError(err string) {
  self.error = fmt.Sprintf("Scan error at line %d: %v", self.In.Line(), err)
}

func (self *Scanner) GetLastError() string {
  return self.error
}

func (self *Scanner) SetQuotedMode() {
  self.mode = MODE_QUOTED
}

func (self *Scanner) SetNormalMode() {
  self.mode = MODE_NORMAL
}

var identchars string = "abcdefghijklmnopqrstuvwxyz1234567890_+-*/^%#"
func isIdentChar(target uint8) bool {
  for _, c := range(identchars) {
    if uint8(c) == target {
	  return true
    }	
  }
  return false
}

var numchars string = "0123456789"
func isNumeric(target uint8) bool {
  for _, c := range(numchars) {
    if uint8(c) == target {
	  return true
    }	
  }
  return false
}

func NewScanner(in ScannerInput) *Scanner {
  return &Scanner{in, MODE_NORMAL, "" }
}

func (self *Scanner) NewToken(t int, s string) *Token {
  return &Token{s, "", t, self.In.Line()}
}

func (self *Scanner) NextToken() *Token {
  if self.mode == MODE_NORMAL {
    return self.ParseNextStandardToken()
  } else {
    return self.ParseQuotedString()	
  }
  return nil
}

func (self *Scanner) ParseNextStandardToken() *Token {
  c := self.In.Current()
  switch c {
  case 0:
	return self.NewToken(EOF, "")
  case '!':
    self.In.Advance()
    return self.NewToken(BANG, "!")
  case '<':
    self.In.Advance()
    c2 := self.In.Current()
    if c2 == '<' {
      self.In.Advance()
      return self.NewToken(LTLT, "<<")
    } else {
      return self.NewToken(LT, "<")
    }
  case '>':
    self.In.Advance()
    c2 := self.In.Current()
    if c2 == '>' {
      self.In.Advance()
      return self.NewToken(GTGT, ">>")
    } else {
      return self.NewToken(GT, ">")
    }
  case '(':
    self.In.Advance()
    return self.NewToken(LPAREN, "(")
  case ')':
    self.In.Advance()
    return self.NewToken(RPAREN, ")")
  case '[':
    self.In.Advance()
    return self.NewToken(LSQUARE, "[")
  case ']':
    self.In.Advance()
    return self.NewToken(RSQUARE, "]")
  case '{':
    self.In.Advance()
    return self.NewToken(LBRACE, "{")
  case '}':
    self.In.Advance()
    return self.NewToken(RBRACE, "}")
  case ',':
    self.In.Advance()
    return self.NewToken(COMMA, ",")
  case '^':
    self.In.Advance()
    return self.NewToken(CARAT, "^")
  case '*':
    self.In.Advance()
    return self.NewToken(STAR, "*")
  case ':':
    self.In.Advance()
    return self.NewToken(COLON, ":")
  case '=':
    self.In.Advance()
    return self.NewToken(EQ, "=")
  case '?':
    self.In.Advance()
    return self.NewToken(QUESTION, "?")
  case '.':
    self.In.Advance()
    return self.NewToken(DOT, ".")
  case '|':
    self.In.Advance()
    return self.NewToken(BAR, "|")
  case '_':
    self.In.Advance()
    return self.NewToken(UNDER, "_")
  case '/':
    self.In.Advance()
    return self.NewToken(SLASH, "/")
  case '+':
    self.In.Advance()
    return self.NewToken(PLUS, "+")
  case '-':
    self.In.Advance()
    if isNumeric(self.In.Current()) {
	  chars := make([]uint8, 0, 16)
	  chars = append(chars, '-')
	  for isNumeric(self.In.Current()) {
        chars = append(chars, self.In.Current())		
	    self.In.Advance()	
      }
      return self.NewToken(NUMERIC, string(chars))
    } else {
      return self.NewToken(MINUS, "-")
    }
  case 'A':
    self.In.Advance()
    return self.NewToken(CMD_CAP_A, "A")
  case 'a':
    self.In.Advance()
    self.SetQuotedMode()
    return self.NewToken(CMD_A, "a")
  case 'c':
    self.In.Advance()
    return self.NewToken(CMD_C, "c")
  case 'e':
    self.In.Advance()
    return self.ParseExtendCommand()
  case 'I':
	  self.In.Advance()
	  return self.NewToken(CMD_CAP_I, "I")
  case 'i':
	 self.In.Advance()
	 self.SetQuotedMode()
	 return self.NewToken(CMD_I, "i")
  case 'j':
    self.In.Advance()
    return self.ParseJumpCommand()
  case 'm':
    self.In.Advance()
    return self.ParseMoveCommand()
  case 'O':
    self.In.Advance()
    return self.NewToken(CMD_CAP_O, "O")
  case 'o':
    self.In.Advance()
    self.SetQuotedMode()
    return self.NewToken(CMD_O, "o")
  case 'r':
    self.In.Advance()
    self.SetQuotedMode()
    return self.NewToken(CMD_R, "r")
  case 't':
    self.In.Advance()
    return self.NewToken(CMD_T, "t")
  case 'x':
    self.In.Advance()
    return self.NewToken(CMD_X, "x")
  case 'w':
    self.In.Advance()
    if self.In.Current() == 'F' {
      self.In.Advance()
      self.SetQuotedMode()
      return self.NewToken(CMD_WF, "wf")
    } else {
      return self.NewToken(CMD_W, "w")
    }
  case 'v':
    self.In.Advance()
    return self.NewToken(CMD_T, "v")
  case 'q':
    return self.ParseQuotedString()
  case '$':
	self.In.Advance()
	str := make([]uint8, 0, 32)
	str = append(str, '$')
	for isIdentChar(self.In.Current()) {
      str = append(str, self.In.Current())
	  self.In.Advance()
	}
	return self.NewToken(VAR, string(str))
  case '@':
	str := make([]uint8, 0, 32)
	str = append(str, '@')
	self.In.Advance()
	for isIdentChar(self.In.Current()) {
      str = append(str, self.In.Current())
	  self.In.Advance()
	}
	return self.NewToken(FIDENT, string(str))
  case '0','1','2','3','4','5','6','7','8','9':
    str := make([]uint8, 0, 32)
	for isNumeric(self.In.Current()) {
	  str = append(str, self.In.Current())
	  self.In.Advance()
	}
	return self.NewToken(NUMERIC, string(str))
  }
  return nil
}

func (self *Scanner) ParseQuotedString() *Token {
  quote := self.In.Current()
  newstr := make([]uint8, 0, 64) // just a guess at a good length
  self.In.Advance()
  for self.In.Current() != quote {
    if self.In.Current() == '\\' && self.In.Peek() == quote {
      self.In.Advance()
      newstr = append(newstr, quote)
    } else if self.In.Current() == 0 {
	  self.SetError("EOF in quoted string")
	  return nil
	} else {
      newstr = append(newstr, self.In.Current())
    }
    self.In.Advance()
  }
  // advance past the close quote
  self.In.Advance()
  quoted := fmt.Sprintf("q%s%s%s", string(quote), string(newstr), string(quote))
  self.SetNormalMode()
  return &Token{quoted, string(newstr), QUOTED_STRING, self.In.Line()}
}

func (self *Scanner) ParseExtendCommand() *Token {
  // current char is the "e" for extend.
  cmd := self.In.Current()
  switch cmd {
	case 'j':
	  self.In.Advance()
	  next := self.In.Current()
	  self.In.Advance()
      switch (next) {
        case 'c':
          return self.NewToken(CMD_EJC, "ejc")	
        case 'l':
          return self.NewToken(CMD_EJL, "ejl")
        case 'p':
          return self.NewToken(CMD_EJP, "ejp")
      }
      self.SetError(fmt.Sprintf("Unknown unit '%v' in extend-jump command", self.In.Current()))
      return nil
	case 'm':
	  self.In.Advance()
	  next := self.In.Current()
	  self.In.Advance()
      switch (next) {
        case 'c':
          return self.NewToken(CMD_EMC, "emc")	
        case 'l':
          return self.NewToken(CMD_EML, "eml")
        case 'p':
          return self.NewToken(CMD_EMP, "emp")
      }
      self.SetError(fmt.Sprintf("Unknown unit '%v' in extend-move command", self.In.Current()))
      return nil
  }
  self.SetError(fmt.Sprintf("Unknown command '%v' in extend command", self.In.Current()))
  return nil
}

func (self *Scanner) ParseMoveCommand() *Token {
  // Current char is the char that came after the "m", so it should
  // be a unit.
  next := self.In.Current()
  self.In.Advance()
  switch (next) {
    case 'c':
	  return self.NewToken(CMD_MC, "mc")	
	case 'l':
	  return self.NewToken(CMD_ML, "ml")
	case 'p':
	  return self.NewToken(CMD_MP, "mp")
  }
  self.SetError(fmt.Sprintf("Unknown unit '%v' in move command", self.In.Current()))
  return nil
}

func (self *Scanner) ParseJumpCommand() *Token {
  // Current char is the char that came after the "j", so it should
  // be a unit.
  next := self.In.Current()
  self.In.Advance()
  switch (next) {
    case 'c':
	  return self.NewToken(CMD_JC, "jc")
	case 'l':
	  return self.NewToken(CMD_JL, "jl")
	case 'p':
	  return self.NewToken(CMD_JP, "jp")
  }
  self.SetError(fmt.Sprintf("Unknown unit '%v' in jump command", self.In.Current()))
  return nil
}
