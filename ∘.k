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
/ └───────────┴───┴───┴───┴───┴───┴───┴───┴───┴───┴───┴─────────┘ *&

O:1; C:""; I:!0; E:0; P:!0; V:`!()

APL:{C::x;$[")"~*x;SYS 1_x;RUN TOK x]}
SYS:{$[x~"ORIGIN 1";*|(O::1;"WAS ",$O);x~"ORIGIN 0";*|(O::0;"WAS ",$O);"INCORRECT COMMAND"]}
UTF:{(&x>-65)^x}
MAP:{,/@[r;&(r:UTF z)~\:x;y]}


err:((`$_32+3#)'e)!e:" "\:"CHARACTER DEPTH DOMAIN DEFN INDEX LABEL LENGTH RANK SYNTAX SYMBOL VALUE"
ERR:{(err x)," ERROR\n      ",C,"\n",((6+E)#" "),"^"}

NM:"_.0123456789"                                       /¯ is mapped to _
AZ:_"A"+!26
OP:"¨<≤=≥>≠∨^⍲⍱-÷+×?⍵∊⍴~↑↓⍳○*⍟→⌈⌊∇∆∘⎕⊂⊃∩∪⊥⊤!|,;:./\\()[]"   /← is mapped to :
SY:UTF OP;                            SE:SY?,";"


A:256#                    "e"   /illegal but allowed in quotation
@[`A;32;                  "_"]  /blank
@[`A;0+"~+-*<>^?!,:/()[]";"+"]  /ascii punctuation
@[`A;192+!64;             "u"]  /utf8 start
@[`A;39;                  "'"]  /quote
@[`A;0+NM,AZ;             "a"]  /alphanum
@[`A;128+!64;             "c"]  /utf8 continuation

TOK:{[x]x:MAP["¯";"_";MAP["←";":";x]];
 c:"_+u'aceqb1"                             /q,b(start,stop quotation) 1(continue alphanum)
 M:c?1_'";"\:"__+u'aee;+_+u'ace;ueeeeece;'qqq_qqq;a_+u'1ee;c_+u'ace;eeeeeeee;qqqqbqqq;b___q___;1_+u'1ce"  /state matrix
 r:(&i:5>s:{M[x;y]}\(c?A)@256/256+x)^x      /r(token list): cut input at i (where token class is low, e.g. not a continuation
 P::0,+\#'r                                 /save token start positions within input C
 b:{P::P i:&~" "=*'x;x i}                   /blanks result in their own tokens, remove them
 r:b VEC b r                                /merge successive numbers (strong vector binding)   
 $[(!0)~e:&s=6;$[+/p:0N~/:r:PRS'r;ERR@*(`val;E::P@*&p);r];ERR@*(`cha;E::*e)]}

VEC:{x:,/|'{x," ",y}\'(&~b&0b~':b:{(*x)':NM}'x)^x;x[i]:(#i:&b&0b~':b)#,""}
 
PRS:{                                                   / apl token           k value              k-type
 ({x}                                                   / "⍳"   ,"+"   "⍰"  → "⍳"   ,"+"   "⍰"     C
  {$[r:`F$" "/:x@&~""~/:x:" "\:x:@[x;&x="_";"-"];r;0N]} / "_1.23E6"    "1"  → ,-1.23E6   ,1.0      F
  `$                                                    / "X1"              → `X1                  s  
  {(#x;x:"'"/:v@&~""~/:v:"'"\:1_-1_x)}                  / "'⍺ isn''t ⍵'"    → ("1";14;"⍺ isn't ⍵") L
 )[*&~0N=(OP;NM;AZ;,"'")?\:*x]x}



F1:{[f;x](*x;f x 1)}
F2:{[f;x;y]$[(!0)~a:*x;(*y;f[*x 1;y 1]);(!0)~b:*y;(a;f[x 1;*y 1]);a~b;(a;f[x 1;y 1]);[L::-1;`len]]}
F:SE#"";f:{F[SY?*x]:(F1 x 1;F2 x 2)}
binom:{f:*/1+!;(f y)%(f x)*f y-x}

/ scalar functions
f(,"+";   {x}     ;   +   )     /dex         plus
f(,"-";    -      ;   -   )     /negate      minus
f( "×";{(x>0)-x<0};   *   )     /signum      times
f( "÷";{ 1. % x  }; 0.+%  )     /reciprocal  divide
f( "⌈";   -_-     ;   |   )     /ceil        max
f( "⌊";    _      ;   &   )     /floor       min
f(,"*";   exp     ;{x^/:y})     /exponential power
f( "⍟";   log     ;  log  )     /logarithm   base
f(,"|";   abs     ; {x/y} )     /absolute    residue
f(,"!";  */1+!    ; binom )     /factorial   binomial
f(,"~";    ~      ;  ""   )     /not
f(,"^";    ""     ;   &   )     /            and
f( "∨";    ""     ;   |   )     /            or
f( "⍲";    ""     ;  ~&   )     /            nand
f( "⍱";    ""     ;  ~|   )     /            nor
f(,"<";    ""     ;   <   )     /            less
f( "≤";    ""     ;  ~>   )     /            not greater
f(,"=";    ""     ;   =   )     /            equal
f( "≥";    ""     ;  ~<   )     /            not less
f(,">";    ""     ;   >   )     /            greater
f( "≠";    ""     ;  ~=   )     /            not equal



/   case   action       PARSE TABLE
PT:( ,"x"; {*x}                                        /no match
    "1+1"; {("1",(((x 1)1)1)[1_x 0;1_x 2])}            /apply dyadic function
    "++1"; {[y]n-:1;x[2+R]:("1",(((y 1)1)0)1_y 2);*y}  /apply monadic function (correct stack for leading +)
    "(1)"; {x 1}                                       /unbrace
    "s:1"; {("1",*(r;V[(*x)1]:r:1_x 2))})              /assign

PN:#'PC:PT[i:2*!2\#PT]
PF:     PT[1+i]


RUN:{L:-1+R:#x                                                  /gap buffer: x[L] right most input token, x[R] stack top
 lu:{$[""~r:V x;`val;("1",r)]}                                  /lookup variable
 co:{[]":"~*x R}                                                /top stack is assign or index assign
 ev:{[y]                                                        /evaluate token on the move
  $[`s~t:@y;$[co[];("s";y);lu y]                                /lookup name if it's not an assignment
    `F~t   ;("1";$[1~#y;!0;#y];y)                               /array from numeric literal, single number is atom
    `C~t   ;$[SE>i:SY?y;mo("+";F@i);*y];y]}                     /build verbs, the rest are atoms ; . / \ ( ) [ ] :(assign)
 s:{x[R-:1]:ev x 1+L-:1}                                        /shift from left(queue) to right(stack) within x
 r:{E::0|P L;x[R+:n-1]:PF[i]r:(n:#PC i:0|*&PC~'PN#\:*'r)#r:R_x} /reduce right side: match case, call action, update R pointer
 while[L>-1;x:r/:s x]                                           /shift once, reduce until fixpoint
 $[~1~#x:R_x;ERR`syn;$[`s~@x:*|*x;ERR x;x]]}

  \APL"2+3"
  \APL"(1-2)"
  \APL"1+-3"
  \APL"A←4"
  \APL"A+2+A←4"

 
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
/ INPUT CLASSES (from the left queue) $ / . : + 1 s ( ) [ ; ]
/  $ start of the sentence
/  / monadic operator (only /)
/  . dyadic operator  (only .)
/  + any function symbol
/  1 array(noun)
/  s symbol
/  ; ( ) [ ] / . :    tokens represented by themselves
/  i index specification results from resolving [..;..]
/  ? anything
/
/ TABLE
/   0     1     2     3     ACTION     TO       REPLACE BY
/   +     /                 mo         0,1   →  + (monadic)
/   1     /                 mo         0,1   →  + (monadic)
/   +     i                 vi         0,1   →  + (monadic)
/   1     i                 ix         0,1   →  1
/   1     +     1           df         0,2   →  1
/   +     .     +           do         0,2   →  + (dyadic)
/   s     :     1           as         0,2   →  1
/   (     1     )           pa         0,2   →  1
/   +     +     1           mf         1,2   →  1
/   s     i     :     1     ax         0,3   →  1
/   [                       br         0,n   →  ⌷ 
/   $     +     1           mf         1,2   →  1 (handled later)
/
/   " "\:"+/  1/ +i 1i 1+1 +.+ s:1 (1) ++1 si:1 ["
/         0   0  0  0  0   0   0   0   1   0    1
/
/ ACTIONS
/  mf  apply monadic function
/  df  apply dyadic  function
/  mo  apply monadic operator
/  do  apply dyadic  operator
/  as  assign
/  ax  indexed assign
/  pa  remove parens
/  br  match closing bracket, count semicolons
/  vi  specify axis to a verb(transpose,grade), reduction or compression
/  ix  index array
/ 
/ REF
/  rapl  rationlized apl - iversion 83
/  dapl  dicitionary of apl - iverson 87
/  ioj   an implmentation of j - hui 9x
/  jdict dictionary - hui/iverson 91
/  2by2  pairwise reduction - bunda/gerth 84
