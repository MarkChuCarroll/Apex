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
package acl

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


// Scanning is messy, because of the quoting modes.
// 
// There are multiple modes determined by contexts, and the same
// thing is treated differently in different modes. The mode
// trigger is quoting. But a string following certain commands is
// automatically quoted, no matter what: the first character after the
// command is the quote character, and standard quotes don't mean anything.
// Auto-quoting applies to i, a, r.


//%token <string> QUOTED_TEXT
//%token <string> VAR IDENT
//%token <int> NUMBER
//%token LPAREN RPAREN COMMA STAR PLUS MINUS LBRACE RBRACE BANG QUESTION ESC SLASH
//%token LBRACK RBRACK NEG OR AND EQUAL DOT
//%token <rune> CHAR
//
//%token CMD_S CMD_M CMD_T CMD_J CMD_P CMD_STAR CMD_D CMD_C CMD_I CMD_A CMD_R CMD_G CMD_X
//%token CMD_L
//%token EOF
//
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
  NextToken() *Token
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
  case '!': // assignment command
    self.In.Advance()
    return self.NewToken(BANG, "!")
  case '<': // insert shell
    self.In.Advance()
    c2 := self.In.Current()
    if c2 == '<' {
      self.In.Advance()
      return self.NewToken(LTLT, "<<")
    } else {
      return self.NewToken(LT, "<")
    }
  case '|': // pipe
    self.In.Advance()
    c2 := self.In.Current()
    if c2 == '|' {
      self.In.Advance()
      return self.NewToken(BARBAR, "||")
      } else {
        return self.NewToken(BAR, "|")
      }
    }
  case '>': // comparison operator
    self.In.Advance()
    return self.NewToken(GT, ">")
  case '(':
    self.In.Advance()
    return self.NewToken(LPAREN, "(")
  case ')':
    self.In.Advance()
    return self.NewToken(RPAREN, ")")
  case '[':
    self.In.Advance()
    return self.NewToken(LBRACK, "[")
  case ']':
    self.In.Advance()
    return self.NewToken(RBRACK, "]")
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
  case '=':
    self.In.Advance()
    return self.NewToken(EQ, "=")
  case '?':
    self.In.Advance()
    return self.NewToken(QUESTION, "?")
  case '.':
    self.In.Advance()
    return self.NewToken(DOT, ".")
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
  case 'a': // append command
    self.In.Advance()
    self.SetQuotedMode()
    return self.NewToken(CMD_A, "a")
  case 'c': // copy command
    self.In.Advance()
    return self.NewToken(CMD_C, "c")
  case 'd': // delete command
    self.In.Advance()
    return self.NewToken(CMD_D, "d")
  case 'g': // global - iteration statement
    self.In.Advance()
    return self.NewToken(CMD_G, "g")
  case 'i': // insert statement
	 self.In.Advance()
	 self.SetQuotedMode()
	 return self.NewToken(CMD_I, "i")
  case 'j': // jump command
    self.In.Advance()
    return self.ParseJumpCommand()
  case 'l': // line unit
    self.In.Advance()
    return self.NewToken(CMD_L, "l")
  case 'm': // move command
    self.In.Advance()
    return self.ParseMoveCommand()
  case 'n': // new file command
    self.In.Advance()    
    return self.NewToken(CMD_N, "n")
  case 'o': // ??
    self.In.Advance()
    self.SetQuotedMode()
    return self.NewToken(CMD_O, "o")
  case 'p': // pick command
    self.In.Advance()
    return self.NewToken(CMD_P, "p")
  case 'r': // replace command
    self.In.Advance()
    self.SetQuotedMode()
    return self.NewToken(CMD_R, "r")
  case 's': // search command
    self.In.Advance()
    return self.NewToken(CMD_S, "s")
  case 't': // tail command?
    self.In.Advance()
    return self.NewToken(CMD_T, "t")
  case 'w': // write file command
    self.In.Advance()
    return self.NewToken(CMD_W, "w")
  case 'W': // write to named file
    self.In.Advance()
    return self.NewToken(CMD_CAP_W, "W")
  case 'x': // execute block command
    self.In.Advance()
    return self.NewToken(CMD_X, "x")
  case '$': // variable
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
