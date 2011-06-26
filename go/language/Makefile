include $(GOROOT)/src/Make.inc

TARG=language
GOFILES=\
	symtable.go\
	parse.go\
	lex.go\
	ast.go\
	ast_cmd.go\
	ast_fun.go\
	ast_re.go\


include $(GOROOT)/src/Make.pkg

parse.go: parse.y
	goyacc -o parse.go parse.y
