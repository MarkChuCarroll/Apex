
//line parse.y:19
package language

import (

  "container/vector"
)

//line parse.y:43
type	yySymType	struct {
	yys	int;
	strval string
	intval int
	node AstNode
	token *Token
    v vector.Vector
    tokens []Token
    nodes []AstNode
}
const	QUOTED_STRING	= 57346
const	VAR	= 57347
const	FIDENT	= 57348
const	NUMERIC	= 57349
const	LPAREN	= 57350
const	RPAREN	= 57351
const	LSQUARE	= 57352
const	RSQUARE	= 57353
const	LBRACE	= 57354
const	RBRACE	= 57355
const	COMMA	= 57356
const	CARAT	= 57357
const	STAR	= 57358
const	COLON	= 57359
const	EQ	= 57360
const	QUESTION	= 57361
const	DOT	= 57362
const	BAR	= 57363
const	UNDER	= 57364
const	SLASH	= 57365
const	PLUS	= 57366
const	EPLUS	= 57367
const	MINUS	= 57368
const	EMINUS	= 57369
const	FUN	= 57370
const	BANG	= 57371
const	LT	= 57372
const	LTLT	= 57373
const	GT	= 57374
const	GTGT	= 57375
const	CMD_CAP_A	= 57376
const	CMD_CAP_I	= 57377
const	CMD_CAP_O	= 57378
const	CMD_A	= 57379
const	CMD_C	= 57380
const	CMD_D	= 57381
const	CMD_G	= 57382
const	CMD_I	= 57383
const	CMD_MC	= 57384
const	CMD_EMC	= 57385
const	CMD_ML	= 57386
const	CMD_EML	= 57387
const	CMD_MP	= 57388
const	CMD_EMP	= 57389
const	CMD_JC	= 57390
const	CMD_EJC	= 57391
const	CMD_JL	= 57392
const	CMD_EJL	= 57393
const	CMD_JP	= 57394
const	CMD_EJP	= 57395
const	CMD_L	= 57396
const	CMD_R	= 57397
const	CMD_X	= 57398
const	CMD_W	= 57399
const	CMD_WF	= 57400
const	CMD_O	= 57401
const	CMD_T	= 57402
const	CMD_V	= 57403
const	EOF	= 57404
var	yyToknames	 =[]string {
	"QUOTED_STRING",
	"VAR",
	"FIDENT",
	"NUMERIC",
	"LPAREN",
	"RPAREN",
	"LSQUARE",
	"RSQUARE",
	"LBRACE",
	"RBRACE",
	"COMMA",
	"CARAT",
	"STAR",
	"COLON",
	"EQ",
	"QUESTION",
	"DOT",
	"BAR",
	"UNDER",
	"SLASH",
	"PLUS",
	"EPLUS",
	"MINUS",
	"EMINUS",
	"FUN",
	"BANG",
	"LT",
	"LTLT",
	"GT",
	"GTGT",
	"CMD_CAP_A",
	"CMD_CAP_I",
	"CMD_CAP_O",
	"CMD_A",
	"CMD_C",
	"CMD_D",
	"CMD_G",
	"CMD_I",
	"CMD_MC",
	"CMD_EMC",
	"CMD_ML",
	"CMD_EML",
	"CMD_MP",
	"CMD_EMP",
	"CMD_JC",
	"CMD_EJC",
	"CMD_JL",
	"CMD_EJL",
	"CMD_JP",
	"CMD_EJP",
	"CMD_L",
	"CMD_R",
	"CMD_X",
	"CMD_W",
	"CMD_WF",
	"CMD_O",
	"CMD_T",
	"CMD_V",
	"EOF",
}
var	yyStatenames	 =[]string {
}
				