/k -fsrc k.w w.k
/   compile src

compile:{macro:funcs:extrn:(0#`)!();table:mem:();class./:mf'-1_"}"\:*(_10 92)\:x;p:parse[;macro;!funcs];`macro`funcs`extrn`table`mem!(macro;p'funcs;extrn;table;mem)}
afun:{exported:1<+/":"=x; x:"."\:@[x;&x=":";"."]; name:`$*x; rtyp:`$x 1; argt:,/`$'""\:x 2;name!,`exported`rtyp`argt`body!(exported;rtyp;argt;y)}
class:{isnum:{x~^x:,/0+("0";x;"9")};$["!"~(:x);mem,:`off`data!(.-1_x;."0x",y) ;isnum(*x);table,:`off`funcs!(.-1_x;";"\:y); ":"~(:x);macro,:(`$-1_x)!,y; #y;funcs,:afun[x;y]; extrn,:afun[x;y]];}
trim:{(+/(_10 13 32)=*x)_x}
mf:{r:"{"\:x;(trim/:r 0;r 1)}

parse:{[x;macro;funcs]x,`ast!type ex token"(",(x`body),")"}
token:{r:0#`!0;{r::r,$[#t:tok x;(t 1)!*t;0#r]; :t}/:x;r}
tok:  {t:toks;n:r:0;x:({$[~#t;0;n::*r::(*t)x]};{t::1_t;x})/:xpnd/:ws/:x;(n#x;!r;n_x)}

first:{0+*x}
toks:(tSem:{`sem!";"~*x}
 tOpa:{`op !$[~+/n:,/(x{y~(#y)#x}/:ops);0;":"~x l:#ops n?1;1+l;l]}
 tFun:{`fun!$n*|/funcs=`$(n:*tSym x)#x}
 tSym:{`sym!(alpha first x)*(alphanum 0+x)?0}
 tNlp:{`nlp!"/"~*x}
 tCon:{`con!$[~num first x;0;"0x"~2#x;18;"j"~t:x n*(#x)>n:(num 0+x)?0;1+n;"."~t;m+**tCon(m:1+n)_x;n]}
 tBra:{`bra!"("~*x}
 tClo:{`clo!clo first x})

unhex:{,/{+/_16 1*"0123456789abcdef"?x}'(,/(2\#x;2))#x}
ws:  {$[space first x;1_x;x]}
xpnd:{$[0~k: :pSym x;x;(#macro)>(!macro)?s:k 1;(macro s),(#$s)_x;x]}
ops: " "\:(">=' %' \' << >> <' <= >= ?/ ?' $[ I. F. C? I? J? F? : + - * % ~ _ * | \ & ^ ! ? < > = C I J F")
set: {@[256#0;x;1]};num:set"0"+!10;alphanum:num+alpha:set"a"+!26;clo:set")]"+0;ctype:set"VCIJF"+0;space:set@!33

/ast
ex:{e:ey x;$[#n:e 1;((`cal;.n;*t);(t:ex 1_n)1);e]}                                   /(e;tail)
ey:{p:*!x;$[`bra~p;seq 1_x;+/`sym`con=p;(p,*x;1_x);`fun~p;fun x;`op~p;mon x;(0;x)]}  /(e;tail)

seq:{r:,`seq;e:0;($[2~#r;r 1;r];1_({~0~*e::ex x};ser)/:x)}  /(`seq;e0;e1;..)
ser:{r::r,,*e;x:e 1;$[`sem~*!x;1_x;x]}
mon:{((`cal;**x;*a);(a:ex 1_x)1)}                           /(`cal;"+";arg)
fun:{((`cal;**x),$[`seq~**a;1_*a;,*a];(a:ex 1_x)1)}         /(`cal;"f";a;b;..)

/type
type:{x}

/\src:"f:I:II{x+y}g:I:II{f x}"
/\a:compile src
/ \a[`funcs;`g]

\
wasm:{}

leb:{a:128;r:!0;({c:a/x;s:c>63;r::r,c+w:$[s+x%a;a;0];w};{x%a})/:x;_r} /signed leb128 for x>=0
whdr:0x0061736d01000000
/wsig:0x60,...

