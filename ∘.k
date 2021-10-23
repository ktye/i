`<"      A P L \\ 3 6 0\n      "
 / ┌────┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬────────┐
 / │    │ ¨ │ ¯ │ < │ ≤ │ = │ ≥ │ > │ ≠ │ ∨ │ ^ │ - │ ÷ │  OVER  │
 / │    │ 1 │ 2 │ 3 │ 4 │ 5 │ 6 │ 7 │ 8 │ 9 │ 0 │ + │ × │ STRIKE │
 / ├────┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬──────┤
 / │      │ ? │ ⍵ │ ∊ │ ⍴ │ ~ │ ↑ │ ↓ │ ⍳ │ ○ │ * │ → │   │      │
 / │      │ Q │ W │ E │ R │ T │ Y │ U │ I │ O │ P │ ← │   │      │
 / ├──────┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴──────┤
 / │        │ ⍺ │ ⌈ │ ⌊ │ _ │ ∇ │ ∆ │ ∘ │ ' │ ⎕ │ ( │ ) │        │
 / │        │ A │ S │ D │ F │ G │ H │ J │ K │ L │ [ │ ] │ RETURN │
 / ├────────┴──┬┴──┬┴──┬┴──┬┴──┬┴──┬┴──┬┴──┬┴──┬┴──┬┴──┬┴────────┤
 / │           │ ⊂ │ ⊃ │ ∩ │ ∪ │ ⊥ │ ⊤ │ | │ ; │ : │ \ │         │
 / │  SHIFT    │ Z │ X │ C │ V │ B │ N │ M │ , │ . │ / │  SHIFT  │
 / └───────────┴───┴───┴───┴───┴───┴───┴───┴───┴───┴───┴─────────┘
 /      ¨¯<≤=≥>≠∨^-÷+×?⍵∊⍴~↑↓⍳○*→←⌈⌊_∇∆∘'⎕⊂⊃∩∪⊥⊤|;:,./\\()[]

O:1; C:""; I:!0; P:0

APL:{C::x;$[")"~*x;SYS 1_x;RUN TOK x]}
SYS:{$[x~"ORIGIN 1";*|(O::1;"WAS ",$O);x~"ORIGIN 0";*|(O::0;"WAS ",$O);"INCORRECT COMMAND"]}
UTF:{(&x>-65)^x}
MAP:{,/@[r;&(r:UTF z)~\:x;y]}
RUN:{`k x}


err:((`$_32+3#)'e)!e:" "\:"CHARACTER DEPTH DOMAIN DEFN INDEX LABEL LENGTH RANK SYNTAX SYMBOL VALUE"
ERR:{(err x)," ERROR\n      ",(P#C),"\n",((6+P)#" "),"^"}

NM:"_.0123456789"
AB:"ABCDEFGHIJKLMNOPQRSTUVXWYZ"
OP:"¨<≤=≥>≠∨^-÷+×?⍵∊⍴~↑↓⍳○*→←⌈⌊∇∆∘⎕⊂⊃∩∪⊥⊤|;:,/\\()[]"

A:256#                    "e"   /illegal, allowed in quotation
@[`A;32;                  "_"] 
@[`A;0+"~+-*<>^?,/()[]";  "+"]
@[`A;192+!64;             "u"]  /utf8 start
@[`A;39;                  "'"]
@[`A;0+NM,AB;             "a"]  /alphanum
@[`A;128+!64;             "c"]  /utf8 continuation

TOK:{[x]x:MAP["¯";"_";x]; c:"_+u'aceqb1"    /q,b(start,stop quotation) 1(continue alphanum)
 M:c?1_'";"\:"__+u'aee;+_+u'ace;ueeeeece;'qqq_qqq;a_+u'1ee;c_+u'ace;eeeeeeee;qqqqbqqq;b___q___;1_+u'1ce"  /state matrix
 r:(&i:5>s:{M[x;y]}\(c?A)@256/256+x)^x      /token list is cut input
 P::0,+\#'r                                 /token (start)positions in input
 P::P b:&~r~\:," "                          /strip empty tokens
 r@:b
 $[(!0)~e:&s=6;$[(!0)~p:0N~/:r:PRS'r;ERR@*(`val;P::P@*&p);r];ERR@*(`cha;P::*e)]}      /return parsed token list, character or value error

PRS:{                                        / apl token           k value          k-type
 ({(UTF OP)?x}                               / "⍳"  (,"+")  "⍰"  → 20   12   0N     i
  {`f$x:@[x;&x="_";"-"]}                     / "_1.23E6"    "1"  → -1.23E6   1.0    f
  `$                                         / "X1"              → `X1              s  
  {"'"/:v@&~""~/:v:"'"\:1_-1_x}              / "'⍺ isn''t ⍵'"    → "⍺ isn't ⍵"      C
 )[*&~0N=(OP;NM;AB;,"'")?\:*x]x}
 
/ RUN (parse and execute)
/ - tokens are marked with $ on the left and form a queue
/ - an empty stack is initialized
/ - four top elements of the stack are matched against the columns
/   of the parse TABLE below
/   - for the first line that matches the action is applied to the
/     columns and they are replaced by the result of the action
/   - if no line matches, a token from right end of the the queue
/     is pushed to top of the stack
/     - while moving, variable names are resolved
/       except if the top of the stack is an assignment,
/       or the first two are an index-specification followed by an assignment
/ - finally there must be two tokens remaining: $ 1 
/   otherwise it's a SYNTAX ERROR
/
/ INPUT CLASSES (from the left queue) $ / . ← + 1 s ( ) [ ; ]
/  $ start of the sentence
/  / monadic operator (only /)
/  . dyadic operator  (only .)
/  + any function symbol
/  1 array(noun)
/  s symbol
/  ; ( ) [ ] / . ←    tokens represented by themselves
/  ⌷ index specification results from resolving [..;..]
/  ? anything
/
/ TABLE
/   0     1     2     3     ACTION     TO       REPLACE BY
/   $     +     1     ?     mf         1,2   →  1
/   ?     +     +     1     mf         2,3   →  1
/   1     +     1     ?     df         0,2   →  1
/   +1    /     ?     ?     mo         0,1   →  + (monadic)
/   +     ⌷     ?     ?     fix        0,1   →  + (monadic)
/   1     ⌷     ?     ?     idx        0,1   →  1
/   +     .     +     ?     do         0,2   →  + (dyadic)
/   s     ←     1     ?     asn        0,2   →  1
/   s     ⌷     ←     1     asi        0,3   →  1
/   (     1     )     ?     inp        0,2   →  1
/   [     ?     ?     ?     ma         0,n   →  ⌷ 
/
/ ACTIONS
/  mf  apply monadic function
/  df  apply monadic function
/  mo  apply monadic operator
/  do  apply monadic operator
/  asn assign
/  asi indexed assign
/  pa  remove parens
/  ma  match closing bracket, count semicolons
/  fix specify axis to a verb(transpose,grade), reduction or compression
/  idx index array
/ 
/ REF
/  rapl  rationlized apl - iversion 83
/  dapl  dicitionary of apl - iverson 87
/  ioj   an implmentation of j - hui 9x
/  jdict dictionary - hui/iverson 91
/  2by2  pairwise reduction - bunda/gerth 84
