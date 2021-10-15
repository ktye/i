`<"      A P L \\ 3 6 0\n      "
O:1
I:""
P:0

APL:{I::x;$[")"~*x;SYS 1_x;$[#i:&x':"@&`\"",(_!32),_"a"+!26;*(ERR`chr;P::*i);RUN TOK x]]}
SYS:{$[x~"ORIGIN 1";*|(O::1;"WAS ",$O);x~"ORIGIN 0";*|(O::0;"WAS ",$O);"INCORRECT COMMAND"]}
UTF:{(&x>-65)^x}
TOK:{}
RUN:{}
ERR:{((`chr`dep`dom`def`idx`lab`len`rnk`syn`sym`val!(" "\:"CHARACTER DEPTH DOMAIN DEFN INDEX LABEL LENGTH RANK SYNTAX SYMBOL VALUE"),\:" ERROR")x),"\n",(6#" "),I,"\n",((6+P)#" "),"^"}

MAP:{,/@[r;&(r:UTF z)~\:x;y]}

TOK:{[x]
 x:MAP["¯";"_";x]
 c:"cuao'"
 m:256#" "
 @[`m;128+!64;"c"]           /unicode-continuation
 @[`m;192+!64;"u"]           /unicode-start
 @[`m;0+"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789._";"a"]  
 @[`m;0+"~+_*<>?,/()[]";"o"]
 @[`m;0+"'";"'"]
}


/ /tokenize k 6class 8state
/C:";*aq/:"
/ 
/A:(33#0),C?"*q*****;;****a*aaaaaaaaaa:;*****aaaaaaaaaaaaaaaaaaaaaaaaaa;/;**aaaaaaaaaaaaaaaaaaaaaaaaaaa;*;*"
/ 
/                  ;*aq/:    ;*aq/:    ;*aq/:    ;*aq/:    ;*aq/:    ;*aq/:    ;*aq/:    ;*aq/:
/M:(C,"br")?m:1_'(";;*aq**";"*;*aq*:";"a;*bq**";"qrrr:/r";"/rrrrrr";":;*aq**";"b;*bq**";"rrrr:/r")
/ 
/ t:&5>M\A@
/ 
/ x:"ab+:4"
/ t x
/ 
/ 0 2 4                 ab+:4
/                       aa*:a
/     ; * a q / :
/  ;  0 1 2 3 1 1
/  *  0 1 2 3 1 5
/  a  0 1 6 3 1 1
/  q  7 7 7 5 4 7
/  /  7 7 7 7 7 7
/  :  0 1 2 3 1 1
/  b  0 1 6 3 1 1
/  r  7 7 7 5 4 7
