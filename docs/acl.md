Apex Command Language
=======================

Cursor Motion Commands
-----------------------

### Position/Distance Units

The units:

- l: line
- c: character
- p: page (24 lines)

Number preceed the units: 24l, 3c, etc.

You can preceed it with a direction, either + or -, in which case, it becomes a
__relative__ position - that is, it moves relative to the current position.



### Commands


- s __[search]__ : takes a regexp, and moves the entire cursor to wrap the match.
- m __[move front]__
- t __[move tail]__
- j __[jump - move both]__
- p(pos1,pos2) __[pick]__ moves the front of the cursor to pos1 and the back of the cursor to pos2. If pos2 is a relative position, it's relative to the new position of the front of the cursor. (If you want to  move the back relative to its cursor position, you can do that
with an m/n sequence.)
- * __[select-all]__


Positions can be:
- Absolute position
- Relative position: pos1(+/-)pos2. If you omit the first pos, then it's relative to
  an contextually appropriate marker - either cursor front or cursor back.
- Search position: /regexp/ 
- 0c: first position in file.
- $: end of file. ($l is last line in file.)

Examples:
- M3l: moves the cursor to a point at the beginning of line 3.
- m-1l: moves the front of the cursor back by one line.
- m+4c: moves the front of the cursor forward by 4 characters
- s3l,+5c: move the front of the cursor to line 3, and and back to 5 characters away from its current position.


Edit Commands
--------------

- d(var) - delete text, and put the deletion into the variable. If the parens and variable
  are omitted, then cut text is discarded. (Or maybe put into OS cut-buf??)
- c(var) - copy text into variable.
- i'text' - insert to the front of the cursor
- a'text' - append to tail of cursor
- r'text' - replace cursor with text



Control Flow Commands
------------------------
- g(/pattern/,cmd_or_block) - for each match of the pattern, execute the command. 
- x(block, params) - execute the block, as if the entire text were the current contents of
   the cursor.

Blocks
--------

{(param, param) body}
{ body }




