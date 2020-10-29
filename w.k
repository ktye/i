/k -fsrc k.w w.k
/   compile src


compile:{reset[];module x;tokens:funcs!token'funcs;ast::funcs!parse'funcs;type'funcs;"ok"}
reset:{ast::tokens::macro::argt::rtyp::locl::body::0#(extrn::funcs::0#`)!table::mem::()}
module:{class./:mf'-1_"}"\:*(_10 92)\:x}
afun:{ex:1<+/":"=x; x:"."\:@[x;&x=":";"."];name:`$*x; rtyp,:name!`$x 1; argt,:name!,/`$'""\:x 2; body,:name!y;name}
class:{$["!"~(:x);mem::`off`data!(.-1_x;."0x",y);num(0+*x);table,:`off`funcs!(.-1_x;trim'";"\:y); ":"~(:x);macro,:(`$-1_x)!,y;*(#y;s:afun[x;y]);funcs,:s;extrn,:s]}
mf:{r:"{"\:x;(trim/:r 0;r 1)}
trim:{$[blank 0+*x;1_x;x]}

parse:{*ex tokens x}
token:{x:"(",(body x),")";  ;r:0#`!0;{r::r,$[#*t:tok x;(t 1)!*t;0#r]; :t}/:x;r}
tok:  {t:toks;n:r:0;x:({$[~#t;0;~n::*r::(*t)x]};{t::1_t;x})/:xpnd/:ws/:x;(n#x;!r;n_x)}

first:{0+*x}
toks:(tSem:{`sem!";"~*x}
 tOpa:{`op !$[~+/n:,/(x{y~(#y)#x}/:ops);0;":"~x l:#ops n?1;1+l;l]}
 tFun:{`fun!n*|/(extrn,funcs)=`$(n:*tSym x)#x}
 tSym:{`sym!(alpha first x)*(alphanum 0+x)?0}
 tNlp:{`nlp!"/"~*x}
 tCon:{`con!$[~num first x;0;"0x"~2#x;18;"."~t:x n*(#x)>n:(num 0+x)?0;m+**tCon(m:1+n)_x;n]}
 tBra:{`bra!"("~*x}
 tTrp:{`sym!"!"~*x}
 tClo:{`clo!clo first x})

unhex:{,/{+/_16 1*"0123456789abcdef"?x}'(,/(2\#x;2))#x}
ws:  {$[space first x;1_x;x]}
xpnd:{$[~#s:(*tSym x)#x;x;(#macro)>(!macro)?v:`$s;(macro v),(#s)_x;x]}
ops: " "\:(">=' I?' F?' %' \' << >> <' >' <= >= ?/ ?' $[ I. F. C? I? F? : + - * % ~ _ * | \ & ^ ? < > = C I F")
set: {@[256#0;x;1]};num:set"0"+!10;alphanum:num+alpha:set"a"+!26;clo:set")]"+0;ctype:set"VCIF"+0;blank:set 10 13 32;space:set@!33

/ast
ex:{u:ey x;$[0~**u;u;0~#u 1;u;(#o)*|/")];"=*o:*u 1;u;"("~o;icl u;(nlp($[|/(extrn,funcs)=`$o;`cal;`op2];o;*u;*v);(v:ex 1_u 1)1)]}     /(e;tail)
ey:{p:*!x;$[`bra~p;seq 1_x;+/`sym`con=p;(p,,*x;1_x);`fun~p;fun x;`op~p;mon x;(0;x)]}  /(e;tail)

seq:{r:,`seq;e:0;($[2~#r;r 1;r];1_({~0~*e::ex x};ser)/:x)}     /(`seq;e0;e1;..)
ser:{r::r,,*e;x:e 1;$[`sem~*!x;1_x;x]}
mon:{((`op1;*x;*a);(a:$["["~ :*x;seq 1_x;ex 1_x])1)}           /(`op1;"+";arg)
fun:{((`cal;*x),$[`seq~**a;1_*a;,*a];(a:ex 1_x)1)}             /(`cal;"f";a;b;..)                    
icl:{(((`icl;u 1;(u:*x)2),$[`seq~*w:*v;1_w;,w]);(v:ex x 1)1)}  /(`icl;"I.";f;x;y..) indirect call

/i
nlp:{$[(`op2;"/")~x 0 1;$[(`con;"1")~x 2; (`op2;"?/"),2_x ;wnlp 2_x];x]}                        /rewrite x/y(n-loop) as ?/(while)
wnlp:{h:(,`seq),$[`sym~**x;();,(`op2;":";(`sym;"n");*x)],,(`op2;":";(`sym;"i");(`con;"0"));h,,(`op2;"?/";(`op2;"<";(`sym;"i");(`sym;$[`sym~**x;(*x)1;"n"]));scat[x 1;(`op2;"+:";(`sym;"i");(`con;"1"))])}
scat:{$[`seq~*x;x,,y;(`seq;x;y)]}

/w string(from ast)
wst:{t:*x;$[`seq~t;,/"(",(";"/:wst'1_x),")";`sym~t;x 1;`con~t;x 1;`op1~t;(sOp1 x 1),wst x 2;`op2~t;,/(embr(x 2);sOp x 1;wst x 3);`cal~t;sCal x;`icl~t;sIcl x;'k(t;"??")]}
sOp:{x,$[|/alpha 0+-x;" ";""]}
sOp1:{$[":"~x;" :";sOp x]}
sCal:{$[3~n:#x;(x 1)," ",wst x 2;4~n;(embr x 2)," ",(x 1)," ",wst x 3;(x 1),"(",(";"/:wst'2_x),")"]}
sIcl:{"(",(x 1),(wst x 2),")(",(";"/:wst'3_x),")"}
embr:{t:*x;e:{"(",(wst x),")"};$[|/t=`cal`op1`op2;e x;wst x]}
cmp:{s:$[2<#s:wst ast x;1_-1_s;s];$[s~b:body x;`ok;(s;b)]}


/typify ast, parse constants
type:{l:((#t)#`x`y`z`x3`x4`x5`x6`x7)!t:argt x;a:typ ast x;locl,:x!l;a}
typ:{(tp(*x))x}                                                         /typed ast:
tp :`con!{$["0x"~2#x;(`F`con;-3$[2_x 1]);|/"."=s:x 1;(`F`con;&s);(`I`con;&s)]}     /(`I`con;3)  (`F`con;3.14)
tp,:`cal!{(,(rtyp s),`cal),(s:`$x 1),typ'2_x}                                      /(`I`cal;`f;x;y;..)
tp,:`icl!{((`$*x 1),`ical),2_x}
tp,:`op1!{((**a),`op1;x 1;a:typ x 2)}                                              /(`I`op1;"+";x)
tp,:`op2!{vo2((**a),`op2;x 1;a:$[":"~x 1;asn[x 2;**t];typ x 2];t:typ x 3)}         /(`I`op2;"+";x;y)
tp,:`sym!{$["!"~x 1;(,`V`trp);((l s),`sym;s:`$x 1)]}                               /(`I`sym;`a)   ,`V`trp
tp,:`seq!{;r:typ'1_x;(,((** :r),`seq)),r}                                          /(`I`seq;a;b;..)
asn:{l,:$[(#l)~(!l)?s:`$x 1;s!y;`!0];(tp`sym)x}
vo2:{$[|/(`$'" "\:"? ::")=`$o:x 1;((,`V`op2),1_x);x]} //force void rtyp

/macro:`v!,"w"
/src:"f:I:II{x+y}g:I:II{f x}"
/funcs:`f`g


/compile src


\
(xt:tp x;xn:nn x;xp:8+x;yt:tp y;yn:nn y;yp:8+y;
((~xt~yt)+xt>5)?!;
(xn<yn)?(dx x;dx y; :2 mk 0);
(~yn)?(dx x;dx y; :(seq(0;xn;1)) drop 1);
r:2 mk 0;w:C xt;eq:8+xt;
i:0;(i:0;(i<xn)?/(a:0;(i:0;(i<yn)?/(k:w*j;a+:((I. eq)(xp+k;yp+k);i+:1));i+:1));(a~yn)?(r:r ucat mki i;i+:yn-1;xp+:w*yn-1);xp+:w)

(`seq;(`op2;":";(`sym;"xt");(`cal;"tp";(`sym;"x")));(`op2;":";(`sym;"xn");(`cal;"nn";(`sym;"x")));(`op2;":";(`sym;"xp");(`op2;"+";(`con;"8");(`sym;"x")));(`op2;":";(`sym;"yt");(`cal;"tp";(`sym;"y")));(`op2;":";(`sym;"yn");(`cal;"nn";(`sym;"y")));(`op2;":";(`sym;"yp");(`op2;"+";(`con;"8");(`sym;"y")));(`op2;"?";(`op2;"+";(`op1;"~";(`op2;"~";(`sym;"xt");(`sym;"yt")));(`op2;">";(`sym;"xt");(`con;"5")));(`sym;"!"));(`op2;"?";(`op2;"<";(`sym;"xn");(`sym;"yn"));(`seq;(`cal;"dx";(`sym;"x"));(`cal;"dx";(`sym;"y"));(`op1;":";(`cal;"mk";(`con;"2");(`con;"0")))));(`op2;"?";(`op1;"~";(`sym;"yn"));(`seq;(`cal;"dx";(`sym;"x"));(`cal;"dx";(`sym;"y"));(`op1;":";(`cal;"drop";(`cal;"seq";(`con;"0");(`sym;"xn");(`con;"1"));(`con;"1")))));(`op2;":";(`sym;"r");(`cal;"mk";(`con;"2");(`con;"0")));(`op2;":";(`sym;"w");(`op1;"C";(`sym;"xt")));(`op2;":";(`sym;"eq");(`op2;"+";(`con;"8");(`sym;"xt")));(`op2;":";(`sym;"i");(`con;"0"));(`seq;(`op2;":";(`sym;"i");(`con;"0"));(`op2;"?/";(`op2;"<";(`sym;"i");(`sym;"xn"));(`seq;(`op2;":";(`sym;"a");(`con;"0"));(`seq;(`op2;":";(`sym;"i");(`con;"0"));(`op2;"?/";(`op2;"<";(`sym;"i");(`sym;"yn"));(`seq;(`op2;":";(`sym;"k");(`op2;"*";(`sym;"w");(`sym;"j")));(`op2;"+:";(`sym;"a");(`seq;(`op2;"(";(`op1;"I.";(`sym;"eq"));(`op2;"+";(`sym;"xp");(`sym;"k")));(`op2;"+";(`sym;"yp");(`sym;"k"))));(`op2;"+:";(`sym;"i");(`con;"1")))));(`op2;"+:";(`sym;"i");(`con;"1")))));(`op2;"?";(`op2;"~";(`sym;"a");(`sym;"yn"));(`seq;(`op2;":";(`sym;"r");(`cal;"ucat";(`sym;"r");(`cal;"mki";(`sym;"i"))));(`op2;"+:";(`sym;"i");(`op2;"-";(`sym;"yn");(`con;"1")));(`op2;"+:";(`sym;"xp");(`op2;"*";(`sym;"w");(`op2;"-";(`sym;"yn");(`con;"1"))))));(`op2;"+:";(`sym;"xp");(`sym;"w")))

\
wasm:{}

leb:{a:128;r:!0;({c:a/x;s:c>63;r::r,c+w:$[s+x%a;a;0];w};{x%a})/:x;_r} /signed leb128 for x>=0
whdr:0x0061736d01000000
/wsig:0x60,...

