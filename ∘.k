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

O:1; C:""; I:!0; E:0; P:!0; V:`!()

APL:{C::x;$[")"~*x;SYS 1_x;PRINT RUN TOK x]}
SYS:{$[x~"ORIGIN 1";*|(O::1;"WAS ",$O);x~"ORIGIN 0";*|(O::0;"WAS ",$O);"INCORRECT COMMAND"]}
UTF:{(&x>-65)^x}


err:((`$_32+3#)'e)!e:" "\:"CHARACTER DEPTH DOMAIN DEFN INDEX LABEL LENGTH RANK SYNTAX SYMBOL VALUE"
ERR:{(err x)," ERROR\n      ",C,"\n",((6+E)#" "),"^"}

NM:"_.0123456789"                                               /¯ is mapped to _
AZ:_("A"+!26),"a"+!26
OP:"¨<≤=≥>≠∨^⍲⍱-÷+×?⍵∊⍴~↑↓⍳○⌽⊖⍉*⍟→⌈⌊∇∆⍋⍒∘⎕⊂⊃∩∪⊥⊤!|,/⌿⍀{\\;:.()[]"   /← is mapped to :
SY:UTF OP;                                       SE:SY?,";"


A:256#                       "e"   /illegal but allowed in quotation
@[`A;32;                     "_"]  /blank
@[`A;0+"~+-*<>^?!,:;/\\()[]";"+"]  /ascii punctuation
@[`A;192+!64;                "u"]  /utf8 start
@[`A;39;                     "'"]  /quote
@[`A;0+NM,AZ;                "a"]  /alphanum
@[`A;128+!64;                "c"]  /utf8 continuation

TOK:{[x]m:{y/:x\:z};x:m["¯";"_"]m["←";":"]m["⍀";"\\[",(_O+"0"),"]"]m["⌿";"/[",(_O+"0"),"]";x]
 c:"_+u'aceqb1"                             /q,b(start,stop quotation) 1(continue alphanum)
 M:c?1_'";"\:"__+u'aee;+_+u'ace;ueeeeece;'qqqqqqq;a_+u'1ee;c_+u'ace;eeeeeeee;qqqqbqqq;b__uq___;1_+u'1ce"  /state matrix
 r:(&i:5>s:{M[x;y]}\(c?A)@256/256+x)^x      /r(token list): cut input at i (where token class is low, e.g. not a continuation
 P::0,+\#'r                                 /save token start positions within input C
 b:{P::P i:&~" "=*'x;x i}                   /blanks result in their own tokens, remove them
 r:b VEC b r                                /merge successive numbers (strong vector binding)   
 $[(!0)~e:&s=6;$[+/p:0N~/:r:PRS'r;ERR@*(`val;E::P@*&p);r];ERR@*(`cha;E::*e)]}

VEC:{x:,/|'{x," ",y}\'(&~b&0b~':b:{(*x)':NM}'x)^x;x[i]:(#i:&b&0b~':b)#,""}
 
PRS:{                                                   / apl token           k value               k-type
 ({x}                                                   / "⍳"   ,"+"   "⍰"  → "⍳"   ,"+"   "⍰"      C
  {$[r:`F$" "/:x@&~""~/:x:" "\:x:@[x;&x="_";"-"];r;0N]} / "_1.23E6"    "1"  → ,-1.23E6   ,1.0       F
  `$                                                    / "X1"              → `X1                   s  
  {(,#x;x:"'"/:v@&~""~/:v:"'"\:1_-1_x)}                 / "'⍺ isn''t ⍵'"    → ("1";,14;"⍺ isn't ⍵") L
 )[*&~0N=(OP;NM;AZ;,"'")?\:*x]x}



F1:{[f;x](*x;f x 1)}
F2:{[f;x;y]$[(!0)~a:*x;(*y;f[*x 1;y 1]);(!0)~b:*y;(a;f[x 1;*y 1]);a~b;(a;f[x 1;y 1]);[L::-1;`len]]}
F:(#OP)#"";Q:(#OP)#"";J:(#OP)#""          /Q stores reductions and J outer product
f:{F[SY?*x]:(F1 x 1;F2 x 2);Q[SY?*x]:(x 2)/;J[SY?*x]:,/(x 2)/:\:}
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

/ mixed functions
up:{[]$[`F':@'(*|x;*|y);0.;_0]}
ls:{@[z;i;+;y i:&x]}
PERAXIS: {[x;a]$[1<#*x;(x 1)@.=+(!*x)_a;*|x]}
IOTA:    {$[(!0)~*x;($[n;,n;!0];O+!n:_0.+*x 1);[L::0;`ran]]}
RESHAPE: {(r;(y 1)@(#y 1)/!*/r:_0.+x 1)}
CATENATE:{((*x)+(*y);(x 1),y 1)}                         /only vectors
COMPRESS:{$[1<r:#*y;CMPRAXIS[O+r-1;x;y];(,#i;(y 1)@i:&(_0.+x 1))]}
CMPRAXIS:{[a;x;y]a:a-O;s:@[*y;a;#i:&_0.+*|x];(s;(*|y)@(*y)//@[p;a;i(p:!s)a])}
EXPAND:  {EXPANDAX[O+(#*y)-1;x;y]}
EXPANDAX:{[a;x;y]a:a-O;r:@[s:*y;a;#v:1=*|x];e:(*/r)#o:$[`C~t:@*|y;" ";`B~t;0b;`I~t;0;0.];$[~(s a)~+/v;[I::0;`ran];(r;@[e;r//@[i;a;(&v)(i:!s)a];*|y])]}
INDEXOF: {o:up[];$[~1~#*x;[I::0;`ran];(*y;O+@[r;&0N=r:(o+*|x)?(o+*|y);#*|x])]}
MEMBER:  {o:up[];(*x;(o+*|x)':(o+*|y))}
TAKEDROP:{[o;n;x]*|x;(n;(*|x)@(*x)//o+!n)}
TAKE:    {$[~(#*|x)~(#*y);[I::0;`ran];0=*/s:abs i:_0.+*|x;(s;0#*|y);TAKEDROP[(i<0)*i+*y;s;y]]}
DROP:    {$[~(#*|x)~(#*y);[I::0;`ran];TAKEDROP[s*i>0;(*y)-s:abs i:_0.+*|x;y]]}
GRADEUP: {(,#r;r:O+<PERAXIS[x;-1+#*x])}
GRADEDN: {(,#r;r:O+>PERAXIS[x;-1+#*x])}
 


f:{F[SY?*x]:1_x}
f(  "⍴";{(,#*x;*x)}    ; RESHAPE                         )
f( ,",";{(,#(x 1);x 1)}; CATENATE                        )
f(  "⍳"; IOTA          ; INDEXOF                         )
f(  "↑"; ""            ; TAKE                            )
f(  "↓"; ""            ; DROP                            )
f(  "⍋"; GRADEUP       ;         ; GRADEUPAX             )
f(  "⍒"; GRADEDN       ;         ; GRADEDNAX             )
f( ,"/"; ""            ; COMPRESS; ""       ; CMPRAXIS   )
f(,"\\"; ""            ; EXPAND  ; ""       ; EXPANDAX   )
f(  "⌽"; REVERSE       ; ROTATE  ; REVERSEAX; ROTATEAX   )
f(  "⍉"; TRANSPOSE     ; TRANSPOSE2                      )
f(  "∊"; ""            ; MEMBER                          )
f(  "⊥"; ""            ; DECODE                          )
f(  "⊤"; ""            ; ENCODE                          )
f( ,"?"; ""            ; DEAL                            )


DOT:{[f;g;x;y]$[f~"∘";JOT[g;x;y];REDUCE[f;DF[g;x;y]]]}
JOT:{[f;x;y]((*x),*y;J[SY?f][x 1;y 1])}
AXIS:{[f;i]$[~((!0)~*i:*i)&(1~#i);[L::0;`dom];(2_F[SY?f])@\:**|i]}
INDEX:{[x;y]$[~(#y)~#*x;*(`ran;L::0);(,/*'y;(*|x)@(*x)//i@'!#'i:((*|)'y:ifill[y;*x])-O)]}
AMEND:{[s;i;y]x:1_lu s;i:(*x)//i@'!#'i:((*|)'j:ifill[i;*x])-O;V[s]:r:(*x;@[$[t~`F;0.+*|x;*|x];i;*(*|y;t:@*|y)]);r}
REDUCE:{[f;x]REDAXE[f;-1+O+#*x;x]}
REDAXE:{[f;a;x]$[1~#*x;(!0;,Q[SY?f]x 1);~(#*x)>a:a-O;[L::0;`ran];((*x)_a;Q[SY?f]'PERAXIS[x;a])]}
MF:{[f;x]($[~`C~@f;*f;F[SY?f;0]])x}
DF:{[f;x;y]($[`C~@f;F[SY?f;1];f 1])[x;y]}


split:{x[0,-1+#x]:";";x[&(0b~':s)&s:";"=*'x]:0;x@&~";"=*'x} /split index expression, fill space for empty index
ifill:{[x;s]x[i]:{(,#x;x)}'O+!'s i:&x~\:0}                  /fill 0s in index vector, given shape s
drop:{x[i]:_0.+1_'x i:&~0~/:x}
mo:{[y]n-:1;x[2+R]:("1",MF[1_y 1;1_y 2]);*y}                /apply monadic function (correct stack for leading +)

/   case   action       PARSE TABLE
PT:( ,"x"; {*x}                                             /no match, keep shifting
    "1+1"; {("1",DF[1_x 1;1_*x;1_x 2])}                     /apply dyadic function
    "++1"; mo                                               /apply monadic
    "(+1"; mo                                               /apply monadic
    "[+1"; mo                                               /apply monadic
    ":+1"; mo                                               /apply monadic
    "+{1"; {("1",REDUCE[1_*x;1_x 2])}                       /reduction
   "+{i1"; {("1",REDAXE[1_*x;**|*1_x 2;1_x 3])} )           /reduce with axis
PT,:("1i"; {("1",INDEX[1_*x;1_x 1])}                        /index array
     "+i"; {("+",AXIS[1_*x;1_x 1])}                         /fix axis to (monadic) function
    "(1)"; {x 1}                                            /unbrace
     ,"["; {[y]n::1+(s:R_x)?"]";("i",drop split n#s)}       /indexing expression: count to ], remove brackets and semicolons
    "+.+"; {("+",(`val;DOT[1_*x;1_x 2]))}                   /dyadic operator
    "s:1"; {("1",*(r;V[(*x)1]:r:1_x 2))}                    /assign
   "si:1"; {("1",AMEND[(*x)1;1_x 1;1_x 3])})                /indexed assignment

PN:#'PC:PT[i:2*!2\#PT]
PF:     PT[1+i]
        

RUN:{L:-1+R:#x   ;                                              /gap buffer: x[L] right most input token, x[R] stack top
 lu:{$[""~r:V x;`val;("1",r)]}                                  /lookup variable
 co:{[](":"~*c)|"i:"~c:*'x 0 1+R}                               /top stack is assign or index assign
 ov:{[y]$[(`C~@s:x L)&("+/"~y);$[SE>SY?s;"{";y];y]}             /change / to { if next token is a verb (compress vs over)
 ev:{[y]                                                        /evaluate token on the move
  $[`s~t:@y;$[co[];("s";y);lu y]                                /lookup name if it's not an assignment
    `F~t   ;("1";$[1~#y;!0;,#y];y)                              /array from numeric literal, single number is atom
    `C~t   ;$[SE>SY?y;ov("+",y);*y];"1",y]}                     /build verbs, the rest are atoms ; . / \ ( ) [ ] :(assign)
 s:{x[R-:1]:ev x 1+L-:1}                                        /shift from left(queue) to right(stack) within x
 r:{E::0|P L;x[R+:n-1]:PF[i]r:(n:#PC i:0|*&PC~'PN#\:*'r)#r:R_x} /reduce right side: match case, call action, update R pointer
 while[L>-1;x:r/:s x]                                           /shift once, reduce until fixpoint
 $["+1"~*'r:R_x;x[R+:1]:("1",MF[1_r 0;1_r 1]);0]                /fix last monadic application
 $[~1~#x:R_x;ERR`syn;*x]}                                       /return the single noun left on the stack otherwise a syntax error

PRINT:{`k x}

/ \APL"2+3"
/ \APL"+/1 2 3"
/ \APL"+/[1]2 3⍴⍳6"
/ \APL"X←⍳6"
/ \APL"2 3 4⍳3"
/ \APL"1 2 3⍳3 4 5"
/ \APL"'ABCDEFGH'⍳'GAFFE'"
/ \APL"1 2 3 4⍳2 3⍴⍳6"
/ \APL"3 4 7∊⍳5"
/ \APL"2↑3 4 5 6"
/ \APL"2 3↑3 4⍴⍳12"
/ \APL"2 ¯3↑3 4⍴⍳12"
/ \APL"2↓3 4 5 6 7"
/ \APL"2 ¯3↑3 4⍴⍳12"
/ \APL"⍋3 5 3 2"
/ \APL"⍒3 5 3 2"
/ \APL"1 0 1 0/3 4⍴⍳12"
/ \APL"1 0 1/[1]3 4⍴⍳12"
/ \APL"1 0 1\\⍳2"
/ \APL"1 0 1 1 1\\3 4⍴'ABCDEFGHIJKL'"
 