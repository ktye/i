`<                 "      A P L \\ 3 6 0\n      "
/      ¨¯<≤=≥>≠∨^-÷+×?⍵∊⍴~↑↓⍳○*→←⌈⌊_∇∆∘'⎕⊂⊃∩∪⊥⊤|;:,./\\()[]
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

O:1; C:""; I:!0; P:0; V:`!()

APL:{C::x;$[")"~*x;SYS 1_x;RUN TOK x]}
SYS:{$[x~"ORIGIN 1";*|(O::1;"WAS ",$O);x~"ORIGIN 0";*|(O::0;"WAS ",$O);"INCORRECT COMMAND"]}
UTF:{(&x>-65)^x}
MAP:{,/@[r;&(r:UTF z)~\:x;y]}


err:((`$_32+3#)'e)!e:" "\:"CHARACTER DEPTH DOMAIN DEFN INDEX LABEL LENGTH RANK SYNTAX SYMBOL VALUE"
ERR:{(err x)," ERROR\n      ",C,"\n",((6+P)#" "),"^"}

NM:"_.0123456789"               /¯ will be mapped to _
AZ:_"A"+!26
OP:"¨<≤=≥>≠∨^-÷+×?⍵∊⍴~↑↓⍳○*→←⌈⌊∇∆∘⎕⊂⊃∩∪⊥⊤|;:,/\\()[]"

A:256#                    "e"   /illegal but allowed in quotation
@[`A;32;                  "_"]  /blank
@[`A;0+"~+-*<>^?,/()[]";  "+"]  /ascii punctuation
@[`A;192+!64;             "u"]  /utf8 start
@[`A;39;                  "'"]  /quote
@[`A;0+NM,AZ;             "a"]  /alphanum
@[`A;128+!64;             "c"]  /utf8 continuation

TOK:{[x]x:MAP["¯";"_";x]; c:"_+u'aceqb1"    /q,b(start,stop quotation) 1(continue alphanum)
 M:c?1_'";"\:"__+u'aee;+_+u'ace;ueeeeece;'qqq_qqq;a_+u'1ee;c_+u'ace;eeeeeeee;qqqqbqqq;b___q___;1_+u'1ce"  /state matrix
 r:(&i:5>s:{M[x;y]}\(c?A)@256/256+x)^x      /r(token list): cut input at i (where token class is low, e.g. not a continuation
 P::0,+\#'r                                 /save token start positions within input C
 b:{P::P i:&~x~\:," ";x i}                  /blanks result in their own tokens, remove them
 r:b VEC b r                                /merge successive numbers (strong vector binding)   
 $[(!0)~e:&s=6;$[+/p:0N~/:r:PRS'r;ERR@*(`val;P::P@*&p);r];ERR@*(`cha;P::*e)]}

VEC:{f:{(*x)':NM};x:{$[(f x)&f y;x," ",y;y]}\x;i:-1+&b&0b~':b:f'x;@[x;i;(,," ")@&,#i]}
 
PRS:{                                                   / apl token           k value             k-type
 ({x}                                                   / "⍳"   ,"+"   "⍰"  → "⍳"   ,"+"   "⍰"    C
  {$[r:`F$" "/:x@&~""~/:x:" "\:x:@[x;&x="_";"-"];r;0N]} / "_1.23E6"    "1"  → ,-1.23E6   ,1.0     F
  `$                                                    / "X1"              → `X1                 s  
  {(#x;x:"'"/:v@&~""~/:v:"'"\:1_-1_x)}                  / "'⍺ isn''t ⍵'"    → (14;"⍺ isn't ⍵")    L
 )[*&~0N=(OP;NM;AZ;,"'")?\:*x]x}


F1:{(*y;x@*|y)}
F2:{[f;x;y]$[(!0)~a:*x;(*y;f[**|x;*|y]);(!0)~b:*y;(a;f[*|x;**|y]);a~b;(a;f[**|x;**|y]);`len]}
NEGATE: F1(-)                            /     - X
PLUS:   F2(+)                            /   X + Y
MINUS:  F2(-)                            /   X - Y
SHAPE:  {(#*x;*x)}                       /     ⍴ X
IOTA:   {(#x;x:O+!**|x)}                 /     ⍳ X
RESHAPE:{(s;(*/s:$[0~#s:*x;,s;s])#*|y)}  /   X ⍴ Y
OVER:   {[f;x]REDUCE[f;x;-1+#*x]}        /    F/ X
REDUCE: {[f;x;a](*x;1#*|x)}              / F[a]/ X  nyi
DOT:    {[f;g;x;y]x}                     / X F.G Y

RUN:{q:s:-1+#x                                  /x[q] last queue token, x[s] stack top value (gap buffer)
 lu:{$[""~r:V x;`val;r]}                        /lookup symbol



 mv:{[]x[1+s-:1]:ev x 1+q-:1}                   /move token from queue to stack
 co:{[]0b}                                      /top stack is assign or index assign
 ev:{$[co[]&`s~t:@x;$[`val~x:lu x;*(x:l::0);x]  /evaluate token on the move, l::0 ends execution on error
       `f~t; (!0;,x); `F~t; (#x;x) ];x}
}
RUN:{$[`C~@x;x;`k x]}

 
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
/   +1    /     ?     ?     mo         0,1   →  + (monadic)
/   +     ⌷     ?     ?     fix        0,1   →  + (monadic)
/   1     ⌷     ?     ?     idx        0,1   →  1
/   $     +     1     ?     mf         1,2   →  1
/   1     +     1     ?     df         0,2   →  1
/   +     .     +     ?     do         0,2   →  + (dyadic)
/   s     ←     1     ?     asn        0,2   →  1
/   (     1     )     ?     inp        0,2   →  1
/   ?     +     +     1     mf         2,3   →  1
/   s     ⌷     ←     1     asi        0,3   →  1
/   [     ?     ?     ?     ma         0,n   →  ⌷ 
/
/ ACTIONS
/  mf  apply monadic function
/  df  apply dyadic  function
/  mo  apply monadic operator
/  do  apply dyadic  operator
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
