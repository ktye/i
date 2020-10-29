/k -fsrc k.w w.k
/   compile src

compile:{reset[];module x;tokens:funcs!token'funcs;ast::funcs!parse'funcs;type'funcs;"ok"}
reset:{ast::tokens::macro::argt::rtyp::locl::body::0#(exprt::extrn::funcs::0#`)!table::mem::()}
module:{class./:mf'-1_"}"\:*(_10 92)\:x}
afun:{ex:1<+/":"=x; x:"."\:@[x;&x=":";"."];name:`$*x; exprt,:$[ex;name;0#`]; rtyp,:name!`$x 1; argt,:name!,/`$'""\:x 2; body,:name!y;name}
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
ex:{u:ey x;$[0~**u;u;0~#u 1;u;(#o)*|/")];"=*o:*u 1;u;"("~o;icl u;(nlp($[|/(extrn,funcs)=`$o;`cal;`op2];o;*u;*v);(v:ex 1_u 1)1)]}  /(e;tail)
ey:{p:*!x;$[`bra~p;seq 1_x;+/`sym`con=p;(p,,*x;1_x);`fun~p;fun x;`op~p;mon x;(0;x)]}                                              /(e;tail)

seq:{r:,`seq;e:0;($[2~#r;r 1;r];1_({~0~*e::ex x};ser)/:x)}                /(`seq;e0;e1;..)
ser:{r::r,,*e;x:e 1;$[`sem~*!x;1_x;x]}
mon:{((`op1;*x;*a);(a:$["["~ :*x;seq 1_x;ex 1_x])1)}                      /(`op1;"+";arg)
fun:{((`cal;*x),$[`seq~**a;1_*a;,*a];(a:ex 1_x)1)}                        /(`cal;"f";a;b;..)                    
icl:{(((`icl;u 1;(u:*x)2),$[`seq~*w:*v;1_w;,w]);(v:ex x 1)1)}             /(`icl;"I.";f;x;y..) indirect call

nlp:{$[(`op2;"/")~x 0 1;$[(`con;"1")~x 2; (`op2;"?/"),2_x ;wnlp 2_x];x]}  /rewrite x/y(n-loop) as ?/(while)
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
typ:{(tp(*x))x}                                                                  /typed ast:
tp :`con!{$["0x"~2#x;(`F`con;-3$[2_x 1]);|/"."=s:x 1;(`F`con;&s);(`I`con;&s)]}   /(`I`con;3)  (`F`con;3.14)
tp,:`cal!{(,(rtyp s),`cal),(s:`$x 1),typ'2_x}                                    /(`I`cal;`f;x;y;..)
tp,:`icl!{((`$*x 1),`ical),2_x}
tp,:`op1!{((**a),`op1;x 1;a:typ x 2)}                                            /(`I`op1;"+";x)
tp,:`op2!{vo2((**a),`op2;x 1;a:$[":"~x 1;asn[x 2;**t];typ x 2];t:typ x 3)}       /(`I`op2;"+";x;y)
tp,:`sym!{$["!"~x 1;(,`V`trp);((l s),`sym;s:`$x 1)]}                             /(`I`sym;`a)   ,`V`trp
tp,:`seq!{;r:typ'1_x;(,((** :r),`seq)),r}                                        /(`I`seq;a;b;..)
asn:{l,:$[(#l)~(!l)?s:`$x 1;s!y;`!0];(tp`sym)x}
vo2:{$[|/(`$'" "\:"? ::")=`$o:x 1;((,`V`op2),1_x);x]} /force void rtyp


compile src

/web-wide assembler
wasm:{sigs:?,/fsig'(funcs,extrn); 0x0061736d01000000,,/{wsec[*x;(x 1)[]]}'((1
     {(leb@#sigs),,/$sigs})                                    /sec1(function signature definitions)
 ( 2;{(leb@#extrn),,/wext'extrn})                              /sec2(import)
 ( 3;{(leb@#funcs),,/leb'(sigs?,/fsig'funcs)})                 /sec3(function signature indexes)
 ( 4;{0x017000,(leb@|/{(x`off)+#x`funcs}'table)})              /sec4(indirect call table size)
 ( 5;{0x010001})                                               /sec5(linear memory)
 ( 7;{(leb 1+#exprt),(leb@#"mem"),"mem",0x0200,,/wexp'exprt})  /sec7(exports)
 ( 9;{(leb@#table),/wtab'table})                               /sec9(elements/call table)
 (10;{(leb@#funcs),,/{(leb@#b),b:wfcn x}'funcs})               /sec10(code)
 (11;{0x010041,(leb mem`off),0x0b,(leb@#mem`data),mem`data}))} /sec11(data)

wext:{m:"ext"; (leb@#m),m,(leb@#$x),($x),0x00,leb sigs?fsig x}
wexp:{(leb@$x),($x),0x00,leb funcs?x}                                   /todo: or (funcs,extrn)?x
wtab:{0x0041,(leb x`off),0x0b,(leb@#x`funcs),,/leb'funcs?,/`$'x`funcs}  / "
fsig:{`$0x60,(leb@#t),(,/wtyp t:argt x),$[`V~r:rtyp x;_0;0x01,wtyp r]}
wtyp:`V`C`I`J`F!0x00017f7e7c
wsec:{(_x),(leb@#y),y}
leb:{a:128;r:!0;({c:a/x;s:c>63;r::r,c+w:$[s+x%a;a;0];w};{x%a})/:x;_r} /signed leb128 for x>=0

wfcn:{[x]"nyi"} /todo

\
hex:  {,/{"0123456789abcdef"@,/(x%16;16/x)}'0+x}  /hex"alpha"
