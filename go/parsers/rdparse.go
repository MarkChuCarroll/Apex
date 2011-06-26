package language

type TokenType int32
const (
 QUOTED_STRING TokenType = iota
 VAR
 FIDENT
 NUMERIC
 LPAREN
 RPAREN 
 LSQUARE
 RSQUARE
 LBRACE
 RBRACE
 COMMA
 CARAT
 STAR
 COLON
 EQ
 QUESTION
 DOT
 BAR
 UNDER
 SLASH
 PLUS
 EPLUS
 MINUS
 EMINUS
 FUN
 BANG
 LT 
 LTLT
 GT
 GTGT
 CMD
 EOF
)

type Parser struct {
  IScanner* scanner
}

func (self *Parser) ExpectToken(t TokenType) bool {
  if self.scanner.CurrentToken().Type == t {
	self.scanner.Advance()
    return true
  } else {
    self.SignalError(fmt.Sprintf("Expected token type %s, but found %s",
                                 t, self.scanner.CurrentToken().Str))
  }
}

func (self *Parser) ParseProgram() (bool success, result []ASTNode) {
  result = make([]AstNode, 10, 0)
  success, stmt := ParseStatement()
  for success {
    result = append(result, stmt)
    success, stmt = ParseStatement()
  }
}

