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

// File: parse.y
// Author: Mark Chu-Carroll <markcc@gmail.com>
// Description: The Parser for the Apex language
%{
package language

import (

  "container/vector"
)
%}

%token <string> QUOTED_STRING
%token <string> VAR
%token <string> FIDENT
%token <int> NUMERIC
%token LPAREN RPAREN LSQUARE RSQUARE LBRACE RBRACE
%token COMMA CARAT STAR COLON EQ QUESTION DOT BAR
%token UNDER SLASH
%token PLUS EPLUS MINUS EMINUS
%token FUN
%token BANG LT LTLT GT GTGT
%token CMD_CAP_A CMD_CAP_I CMD_CAP_O
%token CMD_A CMD_C CMD_D CMD_G CMD_I CMD_MC CMD_EMC CMD_ML CMD_EML CMD_MP 
%token CMD_EMP CMD_JC CMD_EJC CMD_JL CMD_EJL CMD_JP CMD_EJP
%token CMD_L CMD_R CMD_X CMD_W CMD_WF CMD_O CMD_T CMD_V	
%token EOF

%union {
	strval string
	intval int
	node AstNode
	nodes []AstNode
	token *Token
    tokens []Token
}

%type <[]AstNode> stmts
%type <AstNode> stmt
%type <AstNode> fun_stmt
%type <AstNode> command_stmt
%type <*[]Token> ident_list_opt
%type <[]Token> ident_list
%type <*AstNode> conditional_opt
%type <AstNode> choice_command
%type <AstNode> seq_command
%type <AstNode> simple_command
%type <AstNode> atomic_command
%type <[]AstNode> pre_param

%%
prog: stmts EOF;

stmts:
  stmts DOT stmt { v := $1
                   v = append(v, $3)
                   $$ = v 
	            }
| stmt { v := make([]AstNode, 0, 10)
	     $$ = append(v, $1)
	   }
;

stmt:
  fun_stmt { $$ = $1 }
| command_stmt { $$ = $1 }
;

fun_stmt: 
  FUN LPAREN ident_list_opt RPAREN 
  FIDENT
  LPAREN ident_list_opt RPAREN 
  LBRACE stmts RBRACE 
;

ident_list_opt: 
  ident_list { $$ = $1 } // (ident...)
| { $$ = make([]Token, 0, 10) } // ()
;

ident_list:
  ident_list COMMA VAR { v := $1
	                     v = append(v, $3)
	                     $$ = v
	                   }
| VAR { v := make([]Token, 0,10)
	    v = append(v, $1)
	    $$ = v
	  }
;

command_stmt:
  choice_command conditional_opt 
  {
    if $2 == nil {
      $$ = $1
    } else {
      cond := $2
      cond.left = make([]AstNode, 1)
      cond.left[0] = $1
      $$ = cond
    }
  }
;

conditional_opt:
  QUESTION simple_command COLON simple_command
        { result := NewAstNode(NODE_COND)
          result.mid = make([]AstNode, 1)
          result.mid[0] =  $2
          result.right = make([]AstNode, 1)
          result.right[0] = $4
          $$ = result
        }
  |
        { $$ = nil }
;

choice_command:
  choice_command BAR seq_command
	{ result := NewAstNode(NODE_CHOICE)
	  result.left = make([]AstNode, 1)
	  result.left[0] = $1
	  result.right = make([]AstNode, 1)
	  result.right[0] = $3
	  $$ = result
	}
| seq_command { $$ = $1 }
;

seq_command:
  seq_command DOT simple_command 
	{ result := NewAstNode(NODE_SEQ)
      result.left = make([]AstNode, 1)
      result.left[0] = $1
      result.right = make([]AstNode, 1)
      result.right[0] = $3
      $$ = result
    }
| simple_command { $$ = $1 }
;

simple_command:
  atomic_command { $$ = $1 }
| LPAREN command_stmt RPAREN { $$ = $2 }
;

atomic_command:
  command_with_prearg { $$ = $1 }
| post_arg_command { $$ = $1 }
;

command_with_prearg:
  pre_param pre_command 
  { node := $2
    node.left = $1
    $$ = node
  }
;

pre_command:
  CMD_CAP_A // append with input from expression
    
| CMD_CAP_I // insert with input from expression
| CMD_C var_opt   // copy
| CMD_D var_opt   // delete
| CMD_MC 	// move char
| CMD_EMC	// extend and move char
| CMD_ML 	// move line
| CMD_EML
| CMD_MP 	// move page
| CMD_EMP
| CMD_JC 	// jump char
| CMD_EJC
| CMD_JL	// jump line
| CMD_EJL
| CMD_JP	// jump page
| CMD_EJP
| CMD_CAP_O // open file from expression
| CMD_R QUOTED_STRING		// replace
| CMD_X		// execute block on sel
| CMD_T // print
;

var_opt:
  VAR
|
;


post_arg_command:
  PLUS regex
| EPLUS regex
| MINUS regex
| EMINUS regex
| LT QUOTED_STRING
| LTLT QUOTED_STRING
| GT QUOTED_STRING
| GTGT QUOTED_STRING
| CMD_A QUOTED_STRING
| CMD_L	block	// loop
| CMD_W 		// write file
| CMD_WF QUOTED_STRING
| CMD_O	QUOTED_STRING	// open file
| CMD_V		// revert buffer
;

pre_param :
  expr { result := make([]AstNode, 1)
	     result[0] = $1
	     $$ = result }
| LPAREN params RPAREN { $$ = $2 }
;

params: 
  params COMMA expr { $$ = append($1, $3) }
| expr { result := make([]AstNode, 0, 5) // check order
	     result = append(result, $1)
	     $$ = result }
;

post_params:
 LPAREN post_param_list_opt RPAREN
;

post_param_list_opt:
  post_param_list
|
;

post_param_list: 
   post_param_list COMMA post_param
| post_param
;

post_param:
  expr
| regex
;

regex: 
  LSQUARE re_el RSQUARE
;

/* repetition */
re_el:
  re_el choice_el 
| choice_el
;

choice_el:
  choice_el BAR binding_el
| binding_el
;

binding_el: 
  re_bind_opt rep_el
;

re_bind_opt:
  VAR EQ { $$ = &$1 }
|        { $$ = nil }
;

rep_el:
  simple_el re_exp_opt
;

re_exp_opt:
  CARAT ( STAR
	    | PLUS
	    | NUMERIC )
|
;

rep_factor:
  STAR 		{ }
| PLUS 		{ }
| NUMERIC	{ }
;

simple_el:
   UNDER QUOTED_STRING 
|  SLASH QUOTED_STRING 
|  LPAREN re_el RPAREN 
;

expr:
  NUMERIC
| QUOTED_STRING
| funcall
| block
;

block:
  LBRACE block_params_opt stmts RBRACE
;

block_params_opt:
  BAR ident_list BAR
|
;

funcall:
  LPAREN params RPAREN FIDENT post_params
;