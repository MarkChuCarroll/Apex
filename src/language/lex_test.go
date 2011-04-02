package language

import (
  "fmt"
  "testing"
)

func ExpectChar(t *testing.T, actual uint8, expected uint8) {
  if actual != expected {
    t.Error("Incorrect input from scanner.")	
  }
}

func ExpectTrue(t *testing.T, b bool) {
  if !b {
    t.Error("Expected true")
  }
}

func ExpectToken(t *testing.T, tok *Token, ttype int, val string) {
  if tok == nil {
    t.Error("Expected token value, but found nil")	
    return
  }
  if tok.Type != ttype {
    t.Error(fmt.Sprintf("Expected token type %v, but found %v", tok.Type, ttype))
  }
  if tok.Str != val {
    t.Error(fmt.Sprintf("Expected token string '%v', but found '%v'", val, tok.Str))
  }
}

func ExpectQuotedToken(t *testing.T, tok *Token, ttype int, val string, qval string) {
  ExpectToken(t, tok, ttype, val)
  if tok.Strval != qval {
    t.Error(fmt.Sprintf("Expected qval string '%v', but found '%v'", qval, tok.Strval))	
  }
}

func TestScannerInput(t *testing.T) {
  in := NewStringInput("hello world, \n you stupid tester!")
  ExpectChar(t, 'h', in.Current())
  ExpectTrue(t, in.Advance())
  ExpectChar(t, 'e', in.Current())
  ExpectTrue(t, in.Advance())
  ExpectChar(t, 'l', in.Current())
  ExpectTrue(t, in.Advance())
  ExpectChar(t, 'l', in.Current())
  ExpectTrue(t, in.Advance())
  ExpectChar(t, 'o', in.Current())
  ExpectTrue(t, in.Advance())
  ExpectChar(t, ' ', in.Current())
  ExpectChar(t, 'w', in.Peek())
  ExpectTrue(t, in.Advance()) // w
  ExpectTrue(t, in.Advance()) // o
  ExpectTrue(t, in.Advance()) // r
  ExpectTrue(t, in.Advance()) // l
  ExpectTrue(t, in.Advance()) // d
  ExpectTrue(t, in.Advance()) // ,
  ExpectTrue(t, in.Advance()) // ' '
  if in.Line() != 1 {
    t.Error("Expected line 1")
  }
  ExpectTrue(t, in.Advance()) // '\n'
  if in.Line() != 2 {
    t.Error("Expected line 2")
  }
  ExpectTrue(t, in.Advance()) // ' '
  ExpectTrue(t, in.Advance()) // y
  ExpectTrue(t, in.Advance()) // o
  ExpectTrue(t, in.Advance()) // u
  ExpectTrue(t, in.Advance()) // ' '
  ExpectTrue(t, in.Advance()) // s
  ExpectTrue(t, in.Advance()) // t
  ExpectTrue(t, in.Advance()) // u
  ExpectTrue(t, in.Advance()) // p
  ExpectTrue(t, in.Advance()) // i
  ExpectTrue(t, in.Advance()) // d
  ExpectTrue(t, in.Advance()) // ' '
  ExpectTrue(t, in.Advance()) // t
  ExpectTrue(t, in.Advance()) // e
  ExpectTrue(t, in.Advance()) // s
  ExpectTrue(t, in.Advance()) // t
  ExpectTrue(t, in.Advance()) // e
  ExpectTrue(t, in.Advance()) // r
  ExpectTrue(t, in.Advance()) // !
  ExpectTrue(t, !in.Advance())
  if in.Line() != 2 {
    t.Error("Expected line 2")
  }
}


func TestScanner(t *testing.T) {
  in := NewStringInput("27jc18ti'bef'a/faz'oom!/@v$x")
  scanner := NewScanner(in)
  tok := scanner.NextToken()
  ExpectToken(t, tok, NUMERIC, "27")
  tok = scanner.NextToken()
  ExpectToken(t, tok, CMD_JC, "jc")
  tok = scanner.NextToken()
  ExpectToken(t, tok, NUMERIC, "18")
  tok = scanner.NextToken()
  ExpectToken(t, tok, CMD_T, "t")
  tok = scanner.NextToken()
  ExpectToken(t, tok, CMD_I, "i")
  tok = scanner.NextToken()
  ExpectQuotedToken(t, tok, QUOTED_STRING, "q'bef'", "bef")
  tok = scanner.NextToken()
  ExpectToken(t, tok, CMD_A, "a")
  tok = scanner.NextToken()
  ExpectQuotedToken(t, tok, QUOTED_STRING, "q/faz'oom!/", "faz'oom!")
  tok = scanner.NextToken()
  ExpectToken(t, tok, FIDENT, "@v")
  tok = scanner.NextToken()
  ExpectToken(t, tok, VAR, "$x")
  tok = scanner.NextToken()
  ExpectToken(t, tok, EOF, "")
}


