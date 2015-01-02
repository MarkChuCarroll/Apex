
//line parse.y:19
package acl
const QUOTED_TEXT = 57346
const VAR = 57347
const IDENT = 57348
const NUMBER = 57349
const LPAREN = 57350
const RPAREN = 57351
const COMMA = 57352
const STAR = 57353
const PLUS = 57354
const MINUS = 57355
const LBRACE = 57356
const RBRACE = 57357
const BANG = 57358
const QUESTION = 57359
const ESC = 57360
const SLASH = 57361
const LBRACK = 57362
const RBRACK = 57363
const NEG = 57364
const OR = 57365
const AND = 57366
const EQUAL = 57367
const DOT = 57368
const LTLT = 57369
const LT = 57370
const GTGT = 57371
const GT = 57372
const CARAT = 57373
const BAR = 57374
const BARBAR = 57375
const CHAR = 57376
const CMD_S = 57377
const CMD_M = 57378
const CMD_T = 57379
const CMD_J = 57380
const CMD_P = 57381
const CMD_STAR = 57382
const CMD_D = 57383
const CMD_C = 57384
const CMD_I = 57385
const CMD_A = 57386
const CMD_R = 57387
const CMD_G = 57388
const CMD_X = 57389
const CMD_L = 57390
const EOF = 57391

var yyToknames = []string{
	"QUOTED_TEXT",
	"VAR",
	"IDENT",
	"NUMBER",
	"LPAREN",
	"RPAREN",
	"COMMA",
	"STAR",
	"PLUS",
	"MINUS",
	"LBRACE",
	"RBRACE",
	"BANG",
	"QUESTION",
	"ESC",
	"SLASH",
	"LBRACK",
	"RBRACK",
	"NEG",
	"OR",
	"AND",
	"EQUAL",
	"DOT",
	"LTLT",
	"LT",
	"GTGT",
	"GT",
	"CARAT",
	"BAR",
	"BARBAR",
	"CHAR",
	"CMD_S",
	"CMD_M",
	"CMD_T",
	"CMD_J",
	"CMD_P",
	"CMD_STAR",
	"CMD_D",
	"CMD_C",
	"CMD_I",
	"CMD_A",
	"CMD_R",
	"CMD_G",
	"CMD_X",
	"CMD_L",
	"EOF",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 74
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 206

var yyAct = []int{

	66, 90, 65, 3, 38, 24, 33, 2, 99, 4,
	95, 98, 93, 43, 44, 88, 89, 97, 95, 5,
	116, 131, 95, 126, 92, 30, 94, 108, 34, 35,
	91, 130, 128, 113, 94, 84, 107, 37, 94, 28,
	25, 26, 27, 31, 63, 69, 36, 105, 34, 31,
	71, 23, 29, 34, 34, 64, 41, 42, 110, 111,
	103, 100, 17, 18, 80, 79, 85, 20, 19, 81,
	6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
	16, 21, 22, 75, 32, 33, 70, 28, 25, 26,
	27, 28, 25, 26, 27, 133, 62, 31, 129, 127,
	29, 31, 37, 74, 29, 117, 106, 118, 104, 102,
	101, 119, 121, 120, 124, 60, 59, 125, 28, 25,
	26, 27, 122, 58, 57, 132, 85, 56, 31, 109,
	23, 29, 48, 47, 45, 115, 83, 76, 73, 72,
	55, 17, 18, 54, 53, 52, 20, 19, 51, 6,
	7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	21, 22, 28, 25, 26, 27, 50, 49, 114, 123,
	112, 87, 31, 86, 23, 29, 67, 82, 61, 78,
	96, 77, 68, 39, 40, 17, 18, 46, 1, 0,
	20, 19, 0, 6, 7, 8, 9, 10, 11, 12,
	13, 14, 15, 16, 21, 22,
}
var yyPact = []int{

	158, -1000, 35, -3, 3, -1000, 18, 44, 44, 44,
	126, -1000, 125, 124, 163, 162, 144, 141, 140, 139,
	136, 119, 116, 115, -1000, -1000, 108, -1000, -1000, 107,
	-1000, 88, -1000, -3, 158, 158, -1000, -1000, -1000, 83,
	-1000, -1000, -1000, -1000, -1000, 44, -1000, 134, 133, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, 18, 29, 132, 87,
	158, 158, 131, 3, -1000, 16, -1000, 4, -1000, -31,
	-1000, 51, 101, 100, 50, 99, 37, 97, 26, -1000,
	17, 114, 49, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, 11, 130, -1000, -14, -1000, -1000, -1000, -1000,
	44, -1000, -1000, 158, -1000, 87, -1000, 87, 158, -1000,
	-1000, 117, -8, -1000, -1000, -2, -1000, 90, 23, 89,
	-1000, 22, -1000, 0, -1000, 86, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000,
}
var yyPgo = []int{

	0, 188, 7, 3, 9, 19, 46, 4, 187, 25,
	5, 184, 183, 182, 181, 180, 179, 178, 177, 2,
	0, 176, 173, 171, 1, 170, 169, 168,
}
var yyR1 = []int{

	0, 1, 2, 2, 3, 3, 4, 4, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 8, 8, 11,
	11, 7, 10, 10, 10, 10, 10, 10, 13, 13,
	12, 12, 15, 15, 15, 14, 14, 16, 16, 9,
	17, 17, 18, 18, 6, 19, 19, 20, 22, 22,
	21, 21, 23, 23, 23, 23, 27, 27, 26, 26,
	24, 24, 25, 25,
}
var yyR2 = []int{

	0, 2, 2, 1, 3, 1, 3, 1, 2, 2,
	2, 2, 6, 1, 2, 4, 2, 2, 2, 2,
	2, 2, 2, 6, 4, 6, 1, 3, 0, 1,
	1, 2, 1, 4, 1, 1, 6, 1, 2, 1,
	1, 0, 1, 1, 1, 1, 0, 3, 1, 4,
	3, 0, 3, 1, 3, 2, 1, 2, 1, 1,
	2, 0, 1, 1, 4, 4, 2, 0, 2, 1,
	1, 2, 1, 0,
}
var yyChk = []int{

	-1000, -1, -2, -3, -4, -5, 35, 36, 37, 38,
	39, 40, 41, 42, 43, 44, 45, 27, 28, 33,
	32, 46, 47, 16, -10, 5, 6, 7, 4, 17,
	-9, 14, 49, -3, 31, 26, -6, 19, -7, -12,
	-11, 12, 13, -7, -7, 8, -8, 8, 8, 4,
	4, 4, 4, 4, 4, 4, 8, 8, 8, 8,
	8, -17, 8, -4, -5, -19, -20, -21, -13, -10,
	-6, -7, 5, 5, -6, -9, 5, -14, -16, -10,
	-3, -2, -18, 5, 19, -20, -22, -23, 11, 12,
	-24, 26, 20, 8, 34, 18, -15, 48, 42, 39,
	10, 9, 9, 10, 9, 10, 9, 10, 10, 15,
	9, 10, -25, 22, -27, 5, 34, -7, -3, -10,
	-10, -3, 5, -26, -24, -19, 25, 9, 9, 9,
	9, 21, -24, 9,
}
var yyDef = []int{

	0, -2, 0, 3, 5, 7, 0, 41, 41, 41,
	0, 13, 28, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 26, 32, 0, 34, 35, 0,
	37, 51, 1, 2, 0, 0, 8, 61, 9, 0,
	40, 29, 30, 10, 11, 41, 14, 0, 0, 16,
	17, 18, 19, 20, 21, 22, 0, 0, 0, 46,
	0, 0, 0, 4, 6, 61, 56, 0, 31, 0,
	39, 0, 0, 0, 0, 0, 0, 0, 45, 48,
	0, 0, 0, 53, 54, 55, 57, 60, 58, 59,
	62, 63, 73, 67, 70, 0, 38, 42, 43, 44,
	41, 27, 15, 0, 24, 0, 33, 0, 0, 49,
	50, 0, 0, 72, 61, 0, 71, 0, 0, 0,
	47, 0, 52, 0, 69, 61, 66, 12, 23, 25,
	36, 64, 68, 65,
}
var yyTok1 = []int{

	1,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 47, 48, 49,
}
var yyTok3 = []int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var yyDebug = 0

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c > 0 && c <= len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return fmt.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return fmt.Sprintf("state-%v", s)
}

func yylex1(lex yyLexer, lval *yySymType) int {
	c := 0
	char := lex.Lex(lval)
	if char <= 0 {
		c = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		c = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			c = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		c = yyTok3[i+0]
		if c == char {
			c = yyTok3[i+1]
			goto out
		}
	}

out:
	if c == 0 {
		c = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		fmt.Printf("lex %U %s\n", uint(char), yyTokname(c))
	}
	return c
}

func yyParse(yylex yyLexer) int {
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	yyS := make([]yySymType, yyMaxDepth)

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
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

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
		yyVAL = yylval
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
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yychar {
				break
			}
		}
		yyn = yyExca[xi+1]
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
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					fmt.Printf("error recovery pops state %d\n", yyS[yyp].yys)
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
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

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
