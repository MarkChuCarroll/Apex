
//line parse.y:3
package language

import (
  "fmt"
)

//line parse.y:26
type	yySymType	struct {
	yys	int;
	strval string
	intval int
	node *AstNode
	token *Token
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
const	DOLLAR	= 57358
const	STAR	= 57359
const	COLON	= 57360
const	EQ	= 57361
const	QUESTION	= 57362
const	DOT	= 57363
const	BAR	= 57364
const	UNDER	= 57365
const	SLASH	= 57366
const	PLUS	= 57367
const	EPLUS	= 57368
const	MINUS	= 57369
const	EMINUS	= 57370
const	FUN	= 57371
const	BANG	= 57372
const	LT	= 57373
const	LTLT	= 57374
const	GT	= 57375
const	GTGT	= 57376
const	CMD_CAP_A	= 57377
const	CMD_CAP_I	= 57378
const	CMD_CAP_O	= 57379
const	CMD_A	= 57380
const	CMD_C	= 57381
const	CMD_D	= 57382
const	CMD_MC	= 57383
const	CMD_EMC	= 57384
const	CMD_ML	= 57385
const	CMD_EML	= 57386
const	CMD_MP	= 57387
const	CMD_EMP	= 57388
const	CMD_JC	= 57389
const	CMD_EJC	= 57390
const	CMD_JL	= 57391
const	CMD_EJL	= 57392
const	CMD_JP	= 57393
const	CMD_EJP	= 57394
const	CMD_L	= 57395
const	CMD_R	= 57396
const	CMD_X	= 57397
const	CMD_W	= 57398
const	CMD_WF	= 57399
const	CMD_O	= 57400
const	CMD_T	= 57401
const	CMD_V	= 57402
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
	"DOLLAR",
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
}
var	yyStatenames	 =[]string {
}
const	yyEofCode	= 1
const	yyErrCode	= 2
const	yyMaxDepth	= 200

//line yacctab:1
var	yyExca = []int {
-1, 1,
	1, -1,
	-2, 0,
-1, 43,
	9, 61,
	14, 61,
	-2, 58,
}
const	yyNprod	= 94
const	yyPrivate	= 57344
var	yyTokenNames []string
var	yyStates []string
const	yyLast	= 185
var	yyAct	= []int {

  29,   2, 134,  95,  81,  94,  96,  42,   9,  82,
  31,  90,  43,  30,  11,   8,   3, 102,  34,  39,
  33,  38, 150, 109,  79, 117, 114, 127, 116,  40,
  35,  15,  16,  17,  18,   6,  35,  19,  20,  21,
  22, 112, 113, 141,  23,  75,  35,  31,  85,  86,
  30,  11,  80, 139,  84,  34, 115, 103, 145,  24,
  92, 140,  25,  26,  27, 102,  28, 136,  15,  16,
  17,  18, 148, 137,  19,  20,  21,  22,  89,  34,
  99,  23,  88,  65,  67,  68,  69,  89,  98, 100,
 105,  98,   5, 146, 107, 142,  24,  66, 108,  25,
  26,  27,  31,  28,  41,  30, 106, 124,  66, 144,
  34, 101, 120,  87, 123, 131, 125, 122,  31,  36,
 130,  30, 106, 135, 104, 118,  34,  83,  98, 119,
  91, 129, 128,  93, 108,  77, 143,  76,  74,  73,
  72,  71,  70,  78,  32, 138, 135, 126, 147, 111,
 149,  45,  46,  61, 110,  47,  48,  49,  50,  51,
  52,  53,  54,  55,  56,  57,  58,  59,  60,  97,
  62,  63, 133, 132, 121,  64,  44,  14,  13,  12,
  10,  37,   7,   4,   1,
};
var	yyPact	= []int {

   6,-1000,  25,-1000,-1000,-1000, 111,  -1,   8,-1000,
-1000,  43,-1000,-1000, 116,  87,  87,  87,  87, 138,
 137, 136, 135, 134,  67,-1000, 133, 131,-1000,-1000,
-1000,-1000,-1000,-1000,   2,   6, 122,-1000,  43,  43,
  43, 104,  73,-1000,-1000,-1000,-1000, 125, 125,-1000,
-1000,-1000,-1000,-1000,-1000,-1000,-1000,-1000,-1000,-1000,
-1000,-1000, 129,-1000,-1000,-1000, 123,-1000,-1000,-1000,
-1000,-1000,-1000,-1000,-1000,-1000,-1000,-1000,   6, 122,
-1000, 102,  51,-1000,   8,  39,-1000,-1000, 118, 114,
-1000,-1000,-1000,-1000,  83,   1,-1000,  18,  37,  15,
   3, 119, 124,  43, 109,-1000, 114,-1000,   1, 123,
-1000,  12, 128, 127, 123,-1000,-1000,-1000, 107,-1000,
-1000,-1000,  98,  64,-1000,-1000,-1000,  36,-1000,-1000,
  86, 122, 100,  44,-1000,-1000,-1000, 118,-1000,-1000,
-1000,-1000,-1000,  84,-1000,  98,  60,-1000,   6,   9,
-1000,
};
var	yyPgo	= []int {

   0, 184,   1,  16, 183,  92,   4,   9, 182, 181,
   8,  15, 180, 179, 178, 177, 176,  11,  67,  20,
   0,   7, 174, 173, 172,   2,   5,   3,   6, 169,
 154, 149, 147, 145, 144, 143,
};
var	yyR1	= []int {

   0,   1,   2,   2,   3,   3,   4,   6,   6,   7,
   7,   5,   9,   8,   8,  11,  11,  10,  10,  12,
  12,  13,  16,  16,  16,  16,  16,  16,  16,  16,
  16,  16,  16,  16,  16,  16,  16,  16,  16,  16,
  16,  16,  17,  17,  14,  14,  14,  14,  14,  14,
  14,  14,  14,  14,  14,  14,  14,  14,  15,  15,
  21,  21,  22,  23,  23,  24,  24,  25,  25,  18,
  26,  26,  27,  27,  28,  29,  29,  30,  32,  32,
  33,  33,  33,  31,  31,  31,  20,  20,  20,  20,
  19,  35,  35,  34,
};
var	yyR2	= []int {

   0,   1,   3,   1,   1,   1,  11,   1,   0,   3,
   1,   2,   4,   3,   1,   3,   1,   1,   3,   1,
   1,   2,   1,   1,   2,   2,   1,   1,   1,   1,
   1,   1,   1,   1,   1,   1,   1,   1,   1,   2,
   1,   1,   1,   0,   2,   2,   2,   2,   2,   2,
   2,   2,   2,   2,   1,   2,   2,   1,   1,   3,
   3,   1,   3,   1,   0,   3,   1,   1,   1,   3,
   2,   1,   3,   1,   2,   2,   0,   2,   2,   0,
   1,   1,   1,   2,   2,   3,   1,   1,   1,   1,
   4,   3,   0,   5,
};
var	yyChk	= []int {

-1000,  -1,  -2,  -3,  -4,  -5,  29,  -8, -11, -10,
 -12,   8, -13, -14, -15,  25,  26,  27,  28,  31,
  32,  33,  34,  38,  53,  56,  57,  58,  60, -20,
   7,   4, -34, -19,  12,  21,   8,  -9,  22,  20,
  21,  -5, -21, -20, -16,  35,  36,  39,  40,  41,
  42,  43,  44,  45,  46,  47,  48,  49,  50,  51,
  52,  37,  54,  55,  59, -18,  10, -18, -18, -18,
   4,   4,   4,   4,   4, -19,   4,   4, -35,  22,
  -3,  -6,  -7,   5, -11, -10, -10,   9,   9,  14,
 -17,   5, -17,   4, -26, -27, -28, -29,   5,  -2,
  -7,   9,  14,  18,   6, -20,   8,  11, -27,  22,
 -30, -31,  23,  24,   8,  19,  13,  22,   6,   5,
 -10, -22,   8, -21, -20, -28, -32,  15,   4,   4,
 -26,   8, -23, -24, -25, -20, -18,   9, -33,  17,
  25,   7,   9,  -6,   9,  14,   9, -25,  12,  -2,
  13,
};
var	yyDef	= []int {

   0,  -2,   1,   3,   4,   5,   0,   0,  14,  16,
  17,   0,  19,  20,   0,   0,   0,   0,   0,   0,
   0,   0,   0,   0,   0,  54,   0,   0,  57,  58,
  86,  87,  88,  89,  92,   0,   8,  11,   0,   0,
   0,   0,   0,  -2,  21,  22,  23,  43,  43,  26,
  27,  28,  29,  30,  31,  32,  33,  34,  35,  36,
  37,  38,   0,  40,  41,  44,  76,  45,  46,  47,
  48,  49,  50,  51,  52,  53,  55,  56,   0,   0,
   2,   0,   7,  10,  13,   0,  15,  18,  59,   0,
  24,  42,  25,  39,  76,  71,  73,   0,   0,   0,
   0,   0,   0,   0,   0,  60,   0,  69,  70,  76,
  74,  79,   0,   0,  76,  75,  90,  91,   0,   9,
  12,  93,  64,   0,  61,  72,  77,   0,  83,  84,
  76,   8,   0,  63,  66,  67,  68,   0,  78,  80,
  81,  82,  85,   0,  62,   0,   0,  65,   0,   0,
   6,
};
var	yyTok1	= []int {

   1,
};
var	yyTok2	= []int {

   2,   3,   4,   5,   6,   7,   8,   9,  10,  11,
  12,  13,  14,  15,  16,  17,  18,  19,  20,  21,
  22,  23,  24,  25,  26,  27,  28,  29,  30,  31,
  32,  33,  34,  35,  36,  37,  38,  39,  40,  41,
  42,  43,  44,  45,  46,  47,  48,  49,  50,  51,
  52,  53,  54,  55,  56,  57,  58,  59,  60,
};
var	yyTok3	= []int {
   0,
 };

//line yaccpar:1

/*	parser for yacc output	*/

var yyDebug = 0

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

const yyFlag = -1000

func yyTokname(yyc int) string {
	if yyc > 0 && yyc <= len(yyToknames) {
		if yyToknames[yyc-1] != "" {
			return yyToknames[yyc-1]
		}
	}
	return fmt.Sprintf("tok-%v", yyc)
}

func yyStatname(yys int) string {
	if yys >= 0 && yys < len(yyStatenames) {
		if yyStatenames[yys] != "" {
			return yyStatenames[yys]
		}
	}
	return fmt.Sprintf("state-%v", yys)
}

func yylex1(yylex yyLexer, lval *yySymType) int {
	var yychar int
	var c int

	yychar = yylex.Lex(lval)
	if yychar <= 0 {
		c = yyTok1[0]
		goto out
	}
	if yychar < len(yyTok1) {
		c = yyTok1[yychar]
		goto out
	}
	if yychar >= yyPrivate {
		if yychar < yyPrivate+len(yyTok2) {
			c = yyTok2[yychar-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		c = yyTok3[i+0]
		if c == yychar {
			c = yyTok3[i+1]
			goto out
		}
	}
	c = 0

out:
	if c == 0 {
		c = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		fmt.Printf("lex %U %s\n", uint(yychar), yyTokname(c))
	}
	return c
}

func yyParse(yylex yyLexer) int {
	var yyn int
	var yylval yySymType
	var YYVAL yySymType
	YYS := make([]yySymType, yyMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yychar := -1
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		fmt.Printf("char %v in %v\n", yyTokname(yychar), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(YYS) {
		nyys := make([]yySymType, len(YYS)*2)
		copy(nyys, YYS)
		YYS = nyys
	}
	YYS[yyp] = YYVAL
	YYS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yychar < 0 {
		yychar = yylex1(yylex, &yylval)
	}
	yyn += yychar
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yychar { /* valid shift */
		yychar = -1
		YYVAL = yylval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yychar < 0 {
			yychar = yylex1(yylex, &yylval)
		}

		/* look through exception table */
		yyxi := 0
		for {
			if yyExca[yyxi+0] == -1 && yyExca[yyxi+1] == yystate {
				break
			}
			yyxi += 2
		}
		for yyxi += 2; ; yyxi += 2 {
			yyn = yyExca[yyxi+0]
			if yyn < 0 || yyn == yychar {
				break
			}
		}
		yyn = yyExca[yyxi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error("syntax error")
			Nerrs++
			if yyDebug >= 1 {
				fmt.Printf("%s", yyStatname(yystate))
				fmt.Printf("saw %s\n", yyTokname(yychar))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[YYS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current yyp has no shift onn "error", pop stack */
				if yyDebug >= 2 {
					fmt.Printf("error recovery pops state %d, uncovers %d\n",
						YYS[yyp].yys, YYS[yyp-1].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				fmt.Printf("error recovery discards %s\n", yyTokname(yychar))
			}
			if yychar == yyEofCode {
				goto ret1
			}
			yychar = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		fmt.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt		// guard against "declared and not used"

	yyp -= yyR2[yyn]
	YYVAL = YYS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + YYS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	}
	goto yystack /* stack new state and value */
}
