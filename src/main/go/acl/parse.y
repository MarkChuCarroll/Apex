// Copyright 2012 Mark C. Chu-Carroll
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
package acl
%}

%token <string> QUOTED_TEXT
%token <string> VAR IDENT
%token <int> NUMBER
%token LPAREN RPAREN LBRACE RBRACE LBRACK RBRACK 
%token COMMA BANG QUESTION STAR PLUS MINUS SLASH DOT
%token NEG OR AND EQUAL LTLT LT GT CARAT BAR BARBAR
%token UNDER
%token <rune> CHAR

%token CMD_STAR CMD_A CMD_C CMD_D CMD_G CMD_I 
%tokne CMD_J CMD_L CMD_M CMD_N CMD_O CMD_P
%token CMD_R CMD_S CMD_T CMD_W CMD_CAP_W CMD_X
%token EOF

%%

program	:	stmts EOF;

stmts:
  stmts stmt
| stmt
;

stmt:
  stmt CARAT seq_stmt
| seq_stmt   
;

seq_stmt:
  seq_stmt DOT base_stmt
| base_stmt  
;

base_stmt	:	
  CMD_A QUOTED_TEXT
| CMD_C LPAREN VAR RPAREN 
| CMD_D opt_var
| CMD_G LPAREN regex COMMA stmt RPAREN
| CMD_I QUOTED_TEXT
| CMD_J loc_with_opt_dir
| CMD_M loc_with_opt_dir
| CMD_P LPAREN loc_with_opt_dir COMMA loc_with_opt_dir RPAREN
| CMD_R QUOTED_TEXT
| CMD_S regex
| CMD_T loc_with_opt_dir
| CMD_W
| CMD_CAP_W QUOTED_TEXT
| CMD_X LPAREN block RPAREN
| CMD_STAR 
| LTLT QUOTED_TEXT
| LT QUOTED_TEXT
| BARBAR QUOTED_TEXT
| BAR QUOTED_TEXT	
| BANG LPAREN VAR COMMA expr RPAREN /* assignment  */
| AT LPAREN expr RPAREN
;

opt_var:
 LPAREN VAR RPAREN
|
;

dir:
  PLUS
| MINUS
;
	
loc_with_opt_dir:
 dir_opt loc
;

expr:
  VAR
| IDENT LPAREN expr_list_opt RPAREN
| NUMBER
| QUOTED_TEXT
| QUESTION LPAREN stmt COMMA stmt RPAREN
| block
;

loc:
  expr unit
| regex
;

dir_opt:
  dir
|
;

unit:
  CMD_L 
| CMD_C
| CMD_P
;

expr_list_opt:
  expr_list
|
;

expr_list:
  expr_list COMMA expr
| expr
;  	
	
block:
  LBRACE
  block_param_opt
  stmts
  RBRACE
;

block_param_opt:
  LPAREN var_list RPAREN
|
;

var_list:
  var_list COMMA VAR
| VAR
;

regex:
  SLASH re_el_list SLASH
;

re_el_list:
  re_el_list re_el
| re_el
;

re_el:
  choice_re_el rep_opt
;

rep_opt:
  STAR
| PLUS
;  

choice_re_el:
  choice_re_el simple_re_el
|
;

simple_re_el: 
  char
| DOT
| LBRACK neg_opt char_list RBRACK
| LPAREN var_assign_opt re_el_list RPAREN
;

var_assign_opt:
  VAR EQUAL
|
;  

char_list:
  char_list char
| char
;

char:
  CHAR
| ESC CHAR
;  

neg_opt:
  NEG
|
;  


/*
		
VAR	:	'$' ( 'A'..'Z' | 'a' .. 'z' | '_' )+
	;

IDENT: ('A' .. 'Z') ( 'A'..'Z' | 'a' .. 'z' | '_' )*

*/

