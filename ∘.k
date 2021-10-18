`<"      A P L \\ 3 6 0\n      "

O:1; C:""; I:!0; P:0

APL:{C::x;$[")"~*x;SYS 1_x;RUN TOK x]}
SYS:{$[x~"ORIGIN 1";*|(O::1;"WAS ",$O);x~"ORIGIN 0";*|(O::0;"WAS ",$O);"INCORRECT COMMAND"]}
UTF:{(&x>-65)^x}
RUN:{"> "," | "/:x}

UTF:{(&x>-65)^x}
ABC:UTF"¨¯<≤=≥>≠∨^-÷+×?⍵∊⍴~↑↓⍳○*→←⍺⌈⌊_∇∆∘'⎕⊂⊃∩∪⊥⊤|;:,./\\()[] ",_("0"+!9),"A"+!26

err:((`$_32+3#)'e)!e:" "\:"CHARACTER DEPTH DOMAIN DEFN INDEX LABEL LENGTH RANK SYNTAX SYMBOL VALUE"
ERR:{(err x)," ERROR\n      ",(P#C),"\n",((6+P)#" "),"^"}

MAP:{,/@[r;&(r:UTF z)~\:x;y]}


TOK:{[x]x:MAP["¯";"_";x]
 c:"_+u'ace1"
 a:256#                       "e"   /illegal or quoted
 @[`a;32;                     "_"] 
 @[`a;0+"~+-*<>^?,/()[]";     "+"]
 @[`a;192+!64;                "u"]  /utf8 start
 @[`a;39;                     "'"]
 @[`a;0+"_.",_(48+!10),65+!26;"a"]
 @[`a;128+!64;                "c"]  /utf8 continuation
 
 /state matrix _+u'ace  _+u'ace  _+u'ace  _+u'ace; _+u'ace  _+u'ace  _+u'ace  _+u'ace
 M:c?1_'";"\:"__+u'aee;+_+u'ace;ueeeeece;''''_''';a_+u'1ee;c_+u'ace;eeeeeeee;1_+u'1ce"
 (&5>0{M[x;y]}\(c?a)@256/256+x)^x
}

