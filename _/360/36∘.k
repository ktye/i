`<"      A P L \\ 3 6 0\n      "
O:1
APL:{$[")"~*x;SYS 1_x;RUN TOK x]}
SYS:{$[x~"ORIGIN";$O;x~"ORIGIN 1";[O::1;,"1"];x~"ORIGIN 0";[O::0;,"0"];"INCORRECT COMMAND"]}
UTF:{(&x>-65)^x}
TOK:{}
RUN:{}

MAP:{,/@[r;&(r:UTF z)~\:x;y]}

TOK:{[x]
 x:MAP["¯";"&";x]
 c:"cuao'"
 m:256#" "
 @[`m;128+!64;"c"]           /unicode-continuation
 @[`m;192+!64;"u"]           /unicode-start
 @[`m;0+"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.&";"a"]  
 @[`m;0+"~+_*<>?,/()[]";"o"]
 @[`m;0+"'";"'"]
}


/ /tokenize k 6class 8state
/ C:";*aq/:"
/ 
/ A:-127#C?"*q*****;;****a*aaaaaaaaaa:;*****aaaaaaaaaaaaaaaaaaaaaaaaaa;/;**aaaaaaaaaaaaaaaaaaaaaaaaaaa;*;*"
/ 
/ M:(C,"br")?m:1_'(";;*aq**";"*;*aq*:";"a;*bq**";"qrrr:/r";"/rrrrrr";":;*aq**";"b;*bq**";"rrrr:/r")
/ 
/ t:&5>M\A@
/ 
/ x:"ab+:4"
/ t x
/ 
/ 0 2 4