/w compiler
/
/write wasm binary to stdout
/ k -fsrc k.w w.k -e wasm compile src

compile:{reset[];module x;tokens:funcs!token'funcs;ast::funcs!parse'funcs;tast::funcs!dce'type'funcs}
reset:{tast::ast::tokens::macro::argt::rtyp::locl::body::0#(exprt::extrn::funcs::0#`)!table::mem::()}
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
ops: " "\:(">=' I?' F?' %' \' << >> <' >' <= >= ?/ ?' $[ V. I. F. C? I? F? : + - * % ~ _ * | \ & ^ ? < > = C I F")
set: {@[256#0;x;1]};num:set"0"+!10;alphanum:num+alpha:set"a"+!26;clo:set")]"+0;ctype:set"VCIF"+0;blank:set 10 13 32;space:set@!33

/ast
ex:{u:ey x;$[0~**u;u;0~#u 1;u;(#o)*|/")];"=*o:*u 1;u;"("~o;icl u;(($[|/(extrn,funcs)=`$o;`cal;`op2];o;*u;*v);(v:ex 1_u 1)1)]}  /(e;tail)
ey:{p:*!x;$[`bra~p;seq 1_x;+/`sym`con=p;(p,,*x;1_x);`fun~p;fun x;`op~p;mon x;(0;x)]}                                           /(e;tail)

seq:{r:,`seq;e:0;($[2~#r;r 1;r];1_({~0~*e::ex x};ser)/:x)}                /(`seq;e0;e1;..)
ser:{r::r,,*e;x:e 1;$[`sem~*!x;1_x;x]}
mon:{((`op1;*x;*a);(a:$["["~ :*x;seq 1_x;ex 1_x])1)}                      /(`op1;"+";arg)
fun:{((`cal;*x),$[`seq~**a;1_*a;,*a];(a:ex 1_x)1)}                        /(`cal;"f";a;b;..)                    
icl:{(((`icl;u 1;(u:*x)2),$[`seq~*w:*v;1_w;,w]);(v:ex x 1)1)}             /(`icl;"I.";f;x;y..) todo infix: (I.x)(a;b;c) -> x I.(a;b;c)


/typify ast, parse constants
type:{l:((#t)#`x`y`z`x3`x4`x5`x6`x7)!t:argt x;a:typ ast x;locl,:x!l;a}
typ:{(tp(*x))x}                                                                      /typed ast:
tp :`con!{$["0x"~2#x 1;(`F`con;-3$unhex 2_x 1);|/"."=s:x 1;(`F`con;&s);(`I`con;&s)]} /(`I`con;3)  (`F`con;3.14)  constant
tp,:`cal!{(,(rtyp s),`cal),(s:`$x 1),typ'2_x}                                        /(`I`cal;`f;x;y;..)         function call
tp,:`icl!{(,(`$*x 1),`icl),typ'2_x}
tp,:`op1!{a:typ x 2;($[3>`C`I`F?t:`$o;iif t;1~|/"?"=o;iif`$*o;(**a)],`op1;o:x 1;a)}  /(`I`op1;"+";x)             monadic op
tp,:`op2!{$[":"~x 1;tasn x;"/"~x 1;tnlp x;fcmp top2 x]}                              /(`I`op2;"+";x;y)           dyadic op
tp,:`sym!{$["!"~x 1;(,`V`trp);((l s),`sym;s:`$x 1)]}                                 /(`I`sym;`a)   ,`V`trp      local variable
tp,:`seq!{;r:typ'1_x;(,((** :r),`seq)),r}                                            /(`I`seq;a;b;..)            sequence
loca:{l,:$[(#l)~(!l)?!x;x;0#`!0]}
top2:{a:typ x 2;vo2((**a),`op2;x 1;a;typ x 3)}
frel:`$'" "\:"~ < > <= >= ~"
flop:{$[~|/frel~\:`$x 1;0;~~/(`F;`F)~{**x}'2_x;0;1]}
fcmp:{$[~`F~**x;x;flop x;(,(`I`op2)),1_x;x]}
tasn:{b:typ x 3;loca(`$(x 2)1)!**b;((,/(**a;`op2));":";a:(tp`sym)x 2;b)}
tnlp:{$[`I`sym~*a:typ x 2;0;loca`n!`I];loca`i!`I;(`V`op2;"/";a;typ x 3)}
vo2:{$[|/(`$'" "\:"? :: /")=`$o:x 1;((,`V`op2),1_x);x]} /force void rtyp
iif:{$[`C~x;`I;x]}
unhex:{,/{+/_16 1*"0123456789abcdef"?x}'(,/(2\#x;2))#x}

/dce (e.g. v1/v2 macros)
dce :{elim[x;(flas x)^fsym x]}
flas:{$[(`I`op2;":")~2#x;?,/(x 2)1,flas x 3;6~@x;?,/flas'x;0#`]}             /find local assign
fsym:{$[`I`sym~*x;(x 1);(`I`op2;":")~2#x;fsym x 3;6~@x;?,/fsym'x;0#`]}       /find symbols, ignore asn
elim:{$[~#y;x;~`I`seq~*x;x;x@&~|/,/'{(`I`op2;":";(`I`sym;y))~3#x}\:/:[x;y]]} /!

/pretty-print (typed)ast
AST:{f:{i:,/(+\o:"("=x)-+\")"=x;l:0,i@&o;(_10)/:(l#'" "),'(0,&o)_x}; ,/{($x),":",(f@ 'ky),_10}'[!tast;.tast]}

/web-wide assembler
wasm:{sigs:?,/wfsg'allf:(extrn,funcs); 0x0061736d01000000,,/{wsec[*x;(x 1)[]]}'((1
     {(leb@#sigs),,/$sigs})                                                 /sec1(function signature definitions)
 ( 2;{$[#extrn;(leb@#extrn),,/wext'extrn;""]})                              /sec2(import)
 ( 3;{(leb@#funcs),,/leb'(sigs?,/wfsg'funcs)})                              /sec3(function signature indexes)
 ( 4;{$[#table;0x017000,(leb@|/{(x`off)+#x`funcs}'table);""]})              /sec4(indirect call table size)
 ( 5;{0x010001})                                                            /sec5(linear memory)
 ( 7;{(leb 1+#exprt),(leb@#"mem"),"mem",0x0200,$[#exprt;,/wexp'exprt;""]})  /sec7(exports)
 ( 9;{$[#table;(leb@#table),/wtab'table;""]})                               /sec9(elements/call table)
 (10;{(leb@#funcs),,/{(leb@#b),b:wfcn x}'funcs})                            /sec10(code)
 (11;{$[#mem;0x010041,(leb mem`off),0x0b,(leb@#mem`data),mem`data;""]}))}   /sec11(data)

wext:{m:"ext"; (leb@#m),m,(leb@#$x),($x),0x00,leb sigs?wfsg x}
wexp:{(leb@#$x),($x),0x00,leb allf?x}
wtab:{0x0041,(leb x`off),0x0b,(leb@#x`funcs),,/leb'allf?,/`$'x`funcs}
wfsg:{wsig[rtyp x;argt x]}
wsig:{`$0x60,(leb@#y),(,/wtyp y),$[`V~x;_0;0x01,wtyp x]}
wtyp:`V`C`I`J`F!0x00017f7e7c
wsec:{$[#y;(_x),(leb@#y),y;""]}
leb:{a:128;r:!0;({c:a/x;s:c>63;r::r,c+w:$[s+x%a;a;0];w};{x%a})/:x;_r} /signed leb128 for x>=0
rle:{(x i;,/1_-':(i:&1,1_~,/=':x),#x)}
lup:{(!loc)?x}

wfcn:{emit:{(wemt@(*x)1)x};loc:locl x;(wloc x),(emit tast x),0x0b}
wloc:{l:,/.(#argt x)_locl x; $[#*l;(leb@#*u),,/(*u){(leb y),wtyp x}'((u:rle l)1);_0]} /locals
topi:{`$($[flop x;"F";$**x]),x 1}
syms:{,/`$'" "\:x}

wemt :`seq`trp!({,/emit'1_x};{0x00})
wemt,:`con!{$[`I~t:**x;0x41,leb x 1;0x44,8$x 1]}                                                 /const
wemt,:`sym!{0x20,leb lup x 1}                                                                    /get local
wemt,:`cal!{(,/emit'2_x),0x10,leb allf?x 1}
wemt,:`icl!{(,/emit'2_x),(emit x 1),0x11,(leb sigs?wsig[**x;,/{**x}'2_x]),0x00}                  /call indirect
wemt,:`op1!{$[|/"CIF"~\:x 1;wlod x;wtvc x;wcvt x;(`I`op1;"-")~x 0 1;0x4100,(emit x 2),0x6b;"$["~x 1;wcnd x;(emit x 2),wop1 topi x]}
wemt,:`op2!{$[(#wspf)>(!wspf)?s:`$x 1;(wspf s)x;(":"~ :x 1)&1<#?x 1;wasn x;(emit x 2),(emit x 3),wop2 topi x]}

wlod:{a:_alin t:`$x 1;(emit x 2),((`C`I`F!0x2d282b)t),a,0x00}                                    /load(heap)
wasn:{emit@[x;1 3;(":";@[x;1;*x 1])]}                                                            /modified assign
wtvc:{$[2>#a:x 1;0;"?"~a 1;1;0]}
wcvt:{t:iif`$**x 1;xt:**x 2;sn:_3~#x 1;(emit x 2),$[xt~t;"";t~`I;0xab-sn;0xb8-sn]}               /convert
wcnd:{j:1+i:2*!2\-1+#a:1_x 2;f:{(emit x),0x0440,(emit y),0x05}
    $[1~#i;f[a i;a j];,/f'[a i;a j]],(emit@ :a),(2\#a)#0x0b}                                     /cnd $[x;y;..]

wspf:(syms": :: ?/ / ?")!({(emit x 3),0x21,leb lup(x 2)1}                                        /set local
 {a:_alin t:$[(`I`op1;"C?")~2#x 3;`C;**x 3];(emit x 2),(emit x 3),((`C`I`F!0x3a3639)t),a,0x00}   /store(heap)
 {0x02400340,$[`I`con~*x 2;"";(emit x 2),0x450d01],(emit x 3),0x0c000b0b}                        /while  x?/y
 {i:leb lup`i;n:leb lup`n;r:(emit x 2);s:$[`I`sym~*x 2;*("";n:leb lup(x 2)1);0x22,n]             /n-loop
  r,s,0x0440410021,i,0x0340,(emit x 3),0x20,i,0x41016a22,i,0x20,n,0x490d000b0b}
 {(emit x 2),0x0440,(emit x 3),0x0b})                                                            /if

alin:`C`I`F!0 2 3                                    /wasm opcodes
wop1 :(syms"I~ I* F- F+ F_ F% I: F:"                )!0x45679a999c9f0f0f                   
wop2 :(syms"I+ I- I* I% I%' I\ I\' I% I| I^ I<< I>>")!0x6a6b6c6e6d706f7172737476
wop2,:(syms"I< I<' I> I>' I<= I<=' I>= I>=' I~ I&"  )!0x49484b4a4d4c4f4e4671
wop2,:(syms"F+ F- F* F% F< F> F<= F>= F~"           )!0xa0a1a2a36364656661


