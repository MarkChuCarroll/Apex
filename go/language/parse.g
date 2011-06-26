program : ( stmt )+;

stmt : fun_decl | command;

fun_decl : 
    FUN 
    param_list
    IDENT
    '{' ( command )* '}'
;

param_list :
  '('
  ( IDENT (',' IDENT)*  )?
  ')'
;

command : choice ( conditional )? ;

conditional: '?' simple_cmd ':' simple_cmd ;

choice : seq ( '|' seq  )*;

seq : simple ( '.' simple )*;

simple :
     atomic
|   '(' command ')'
;


atomic : ( pre_param )? CMD ( post_param  )?
   ;

post_param : regex | block | var | expr_list;

pre_param :
  expr
| expr_list
;

expr_list :
  '(' ( expr (',' expr )* )?  ')'
;



/////////////////////////
// Regular expressions

regex : '[' re_el ']';

re_el : ( choice_el )+
      ;

choice_el : binding_el ( '|' binding_el )* ;

binding_el : (IDENT '=')? rep_el ;

rep_el : simple_el ('^' ('*' | '+' | NUMERIC )? ;

simple_el :
  '_' QUOTED
| '/' QUOTED
| '(' re_el  ')'
;

expr :
  NUMERIC
| QUOTED
| funcall
| block
;

block :
	  '{' ( '|' IDENT ( ',' IDENT)* '|' )? ( stmt )? '}'
;

funcall:
  expr_list IDENT
;

po