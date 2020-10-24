/k -fsrc k.w w.k
/   compile src

compile:{macro:funcs:extrn:(0#`)!();table:mem:();class./:mf'-1_"}"\:*(_10 92)\:x;p:parse[;macro;!funcs];`macro`funcs`extrn`table`mem!(macro;p'funcs;extrn;table;mem)}
afun:{exported:1<+/":"=x; x:"."\:@[x;&x=":";"."]; name:`$*x; rtyp:`$x 1; argt:,/`$'""\:x 2;name!,`exported`rtyp`argt`body!(exported;rtyp;argt;y)}
class:{isnum:{x~^x:,/0+("0";x;"9")};$["!"~(:x);mem,:`off`data!(.-1_x;."0x",y) ;isnum(*x);table,:`off`funcs!(.-1_x;";"\:y); ":"~(:x);macro,:(`$-1_x)!,y; #y;funcs,:afun[x;y]; extrn,:afun[x;y]];}
trim:{(+/(_10 13 32)=*x)_x}
mf:{r:"{"\:x;(trim/:r 0;r 1)}

parse:   {[x;macro;funcs]x,`ast!,ex tokenize"(",(x`body),")"}
tokenize:{r:0#`!0;{r::r,$[0~n:(l:tok x)2;0#r;(*n)!(,l 1),1_n];*l}/:x;r}
/tokenize:{r:0#`!0;{r::r,$[0~n:(l:tok x)1;0#r;( \*n)!(,l 1),1_n];*l}/:x;r}
tok:     {t:toks;r:0;x:({$[~#t;0;0~r::(*t)x;1;0]};{t::1_t;x})/:xpnd/:ws/:x;((*r)_x;(*r)#x; :r)}

first:{0+*x}
toks:(tSem:{`sem!";"~*x}
 tOpa:{`op !$[~+/n:,/(x{y~(#y)#x}/:ops);0;":"~x l:#ops n?1;1+l;l]}
 tFun:{`fun!$n*|/funcs=`$(n:*tSym x)#x}
 tSym:{`sym!(alpha first x)*(alphanum 0+x)?0}
 tNlp:{`nlp!"/"~*x}
 tCon:{`con!$[~num first x;0;"0x"~2#x;18;"j"~t:x n*(#x)>n:(num 0+x)?0;1+n;"."~t;m+**tCon(m:1+n)_x;n]}
 tCnd:{$[3>#x;0;~"$["~x 0 1;0;*(2+*r;(,`cnd),1_ :r;r:pSeq 2_x)]}
 tBra:{$["("~*x;(1;(`bra;"("));0]}
 tCbr:{$[")"~*x;(1;(`cbr;")"));0]}
 tCpa:{$["]"~*x;(1;(`cpa;"]"));0]})

unhex:{,/{+/_16 1*"0123456789abcdef"?x}'(,/(2\#x;2))#x}
ws:  {$[space first x;1_x;x]}
xpnd:{$[0~k: :pSym x;x;(#macro)>(!macro)?s:k 1;(macro s),(#$s)_x;x]}
ops: " "\:(">=' %' \' << >> <' <= >= ?/ ?' $[ : + - * % ~ _ * | \ & ^ ! ? < > = C I J F") /long-first
set: {@[256#0;x;1]};num:set"0"+!10;alphanum:num+alpha:set"a"+!26;ctype:set"VCIJF"+0;space:set@!33

seq:{r:,`seq;e:0;(r;1_({~0~*e::ex x};{r::r,,*e;x:e 1;$[`sem~*!x;1_x;x]})/:x)}  /(`seq;e0;e1;..)
cnd:{r:seq x;((`cnd,1_*r);r 1)}                                                /(`cnd;e1;e2;..)
mon:{$["$["~o:**x;cnd 1_x;((`mon;o;*a);(a:ex 1_x)1)]}                          /(`mon;"+";arg)
fun:{o:**x;($[`seq~**a;((`cal;o),1_**a);(`mon;o;*a)];(a:ex 1_x)1)}
dya:{a:*x;m:*n:mon x 1;((`dya;m 1;a;m 2);n 1)}                                 /(`dya;"+";a;b)

exx:{p:*!x;$[`bra~p;seq 1_x;+/`sym`con=p;(p,*x;1_x);`fun~p;fun x;`op~p;mon x;(0;x)]}        /(e;tail)
ex: {h:exx x;$[`op~*!(h 1);dya h;h]}

/ \e:*ex tokenize"$[3;4;5]"
/ \e:*ex tokenize"+4.5"
/ \e:ex tokenize"1+2"

/\src:"f:I:II{x+y}g:I:II{f x}"
/\a:compile src
/ \a[`funcs;`g]

ftCon:{(x;(*.tCon x)#x)}
 \ftCon"a"
 \ftCon"1"
 \ftCon"12"
 \ftCon"12+"
 \ftCon"1x"
 \ftCon"1.23"
 \ftCon"1000j+5"
 \ftCon"0x0123456789abcdefgji"


\
wasm:{}

leb:{a:128;r:!0;({c:a/x;s:c>63;r::r,c+w:$[s+x%a;a;0];w};{x%a})/:x;_r} /signed leb128 for x>=0
whdr:0x0061736d01000000
/wsig:0x60,...

