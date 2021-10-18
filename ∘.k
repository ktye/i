`<"      A P L \\ 3 6 0\n      "

O:1; C:""; I:!0; P:0

APL:{C::x;$[")"~*x;SYS 1_x;RUN TOK x]}
SYS:{$[x~"ORIGIN 1";*|(O::1;"WAS ",$O);x~"ORIGIN 0";*|(O::0;"WAS ",$O);"INCORRECT COMMAND"]}
UTF:{(&x>-65)^x}
MAP:{,/@[r;&(r:UTF z)~\:x;y]}
RUN:{`k x}


/ ¨¯<≤=≥>≠∨^-÷+×?⍵∊⍴~↑↓⍳○*→←⌈⌊_∇∆∘'⎕⊂⊃∩∪⊥⊤|;:,./\\()[] ",_("0"+!9),"A"+!26

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
 

