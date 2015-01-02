Apogee: a text editing programming language
==============================================

Introduction
---------------

The core idea of Apogee is _text scanning_. Apogee programs always operate on a
_text buffer_. The programs work by moving a cursor around the buffer, and 
looking at or modifying the text at the cursor point. 

The scanning model is combined with _goal directed evaulation_: every construct
in the language can either _succeed_ or _fail_. Control flow constructs are
built on the success or failure of commands, rather than on true/false comparisons.

For example, to replace all instances of "hello" in a buffer with goodbye, in
Apogee, you'd write:

    every(find("hello"),replace("goodbye"))

* `every` is a looping construct, which continues to repeat its body until it fails.
* `find` scans forward until the cursor covers text matching its parameter.
* `replace` removes the text under the cursor, and replaces it with its parameter.

In more concise form, this could be written;
    !(/"hello", r"goodbye")
    

Language Constructs
---------------------

* Sequencing: any group of statements separated by commas are executed in sequence. If any element
  of the sequence fails, then the entire sequence fails.
  
* Alternation: any group of statements separated by `|`s are alternatives. The statements are
  executed left-to-right until one succeeds; it fails if all alternatives fail.

* Repetition: "every" or "!" followed by a command will repeatedly execute the command until it
   fails.
   
* Motion:
   * `move(dist)`: move a set distance; succeeds if the move ends at a location in the document. (At position 10,
     `move(-12c)` will fail.) Can also be written `+dist` or `-dist`.
   * `to(statement)`: executes the statement at the current cursor position. If it succeeds, then
     the entire statement succeeds. If not, it moves the cursor forward one step, and tries again. Keeps
     stepping cursor forward until there's no way the statement can succeed.
     Can also be written `/statement`.
   * `backto(statement)`: same as `to`, except that the cursor steps backward.
   * bare location: a location outside of any statement means move the cursor to that location.
   * `location:location`: two locations mean the cursor should be moved to the first location,
     and set to span the region to the second location. If the locations are moves, then the
     current cursor selection is extended/contracted by those distances.

* Distances: distances are written is a number followed by a unit. The unit can be
  `c` for characters or `l` for lines.

* Matching:
  * `"string"`: a literal string succeeds if the text starting at the cursor is identical to the
    string.
  * `many(stmt)`: matches as many characters at the cursor as match the condition. Takes an optional
    second and third argument: second arg is how many matches must succeed before the condition
    is treated as successful; the third arg is the maximum number of matches; the many will terminate
    successfully before trying any more matches after the maximum. By default, minimum is 0, and maximum
    is infinite.    
  * `not(stmt)`: not succeeds whenever its parameter fails. `many(not("a"))` will succeed and match
    any sequence of characters that are not "a".
  * `set(str)`: character set: matches any character which is in the character set. 
    Equivalent to typical regexp character sets, or to `str[0] | str[1] | str[2] | str[3]`.

