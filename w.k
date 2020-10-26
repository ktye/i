/k -fsrc k.w w.k
/   compile src

ast:macro:argt:rtyp:body:extrn:0#(funcs:0#`)!table:mem:()

compile:{module x;ast::funcs!parse'funcs}
module:{class./:mf'-1_"}"\:*(_10 92)\:x}
afun:{ex:1<+/":"=x; x:"."\:@[x;&x=":";"."];name:`$*x;rtyp,:name!`$x 1; argt,:name!,/`$'""\:x 2; body,:name!y;name}
class:{ \(`class;x);$["!"~(:x);mem,:`off`data!(.-1_x;."0x",y);num(0+*x);table,:`off`funcs!(.-1_x;";"\:y); ":"~(:x);macro,: \(`$-1_x)!,y; #y;funcs,:afun[x;y]]}
mf:{r:"{"\:x;(trim/:r 0;r 1)}
trim:{(+/(_10 13 32)=*x)_x}

parse:{ex token"(",(body x),")"}
token:{r:0#`!0;{r::r,$[#*t:tok x;(t 1)!*t;0#r]; :t}/:x;r}
tok:  {t:toks;n:r:0;x:({$[~#t;0;~n::*r::(*t)x]};{t::1_t;x})/:xpnd/:ws/:x;(n#x;!r;n_x)}

first:{0+*x}
toks:(tSem:{`sem!";"~*x}
 tOpa:{`op !$[~+/n:,/(x{y~(#y)#x}/:ops);0;":"~x l:#ops n?1;1+l;l]}
 tFun:{`fun!n*|/funcs=`$(n:*tSym x)#x}
 tSym:{`sym!(alpha first x)*(alphanum 0+x)?0}
 tNlp:{`nlp!"/"~*x}
 tCon:{`con!$[~num first x;0;"0x"~2#x;18;"j"~t:x n*(#x)>n:(num 0+x)?0;1+n;"."~t;m+**tCon(m:1+n)_x;n]}
 tBra:{`bra!"("~*x}
 tClo:{`clo!clo first x})

unhex:{,/{+/_16 1*"0123456789abcdef"?x}'(,/(2\#x;2))#x}
ws:  {$[space first x;1_x;x]}
xpnd:{$[~#s:(*tSym x)#x;x;(#macro)>(!macro)?v:`$s;(macro v),(#s)_x;x]}
ops: " "\:(">=' %' \' << >> <' <= >= ?/ ?' $[ I. F. C? I? J? F? : + - * % ~ _ * | \ & ^ ! ? < > = C I J F")
set: {@[256#0;x;1]};num:set"0"+!10;alphanum:num+alpha:set"a"+!26;clo:set")]"+0;ctype:set"VCIJF"+0;space:set@!33

/ast
ex:{u:ey x;$[0~**u;u;0~#u 1;u;|/")];"=o:**u 1;u;(($[|/funcs=`$o;`cal;`op2];o;*u;*v);(v:ex 1_u 1)1)]}     /(e;tail)
ey:{p:*!x;$[`bra~p;seq 1_x;+/`sym`con=p;(p,,*x;1_x);`fun~p;fun x;`op~p;mon x;(0;x)]}  /(e;tail)

seq:{r:,`seq;e:0;($[2~#r;r 1;r];1_({~0~*e::ex x};ser)/:x)}  /(`seq;e0;e1;..)
ser:{r::r,,*e;x:e 1;$[`sem~*!x;1_x;x]}
mon:{((`op1;*x;*a);(a:$["["~ :*x;seq 1_x;ex 1_x])1)}        /(`op1;"+";arg)
fun:{((`cal;**x),$[`seq~**a;1_*a;,*a];(a:ex 1_x)1)}         /(`cal;"f";a;b;..)


/typify ast
type:{ (`type;*x);(tp(*x))x}
tp :`con!{$["0x"~2#x;(`F`con;-3$[2_x 1]);|/"."=s:x 1;(`F`con;&x 1);"j"~*|s;(`J`con;&-1_s);(`I`con;&s)]}
tp,:`cal!{(,(funcs`$x 1),`cal),(x 1),type'2_x}
tp,:`op1!{((**a),`op1;a:type x 2)}
tp,:`op2!{((**a),`op2;a:type x 2;type x 3)}
tp,:`seq!{l:type'1_x;(,((** :l),`seq)),l}

/macro:`v!,"w"
/src:"f:I:II{x+y}g:I:II{f x}"
/funcs:`f`g!`I`I

a:compile src
/ \a[`funcs;`ini;`ast]

\
wasm:{}

leb:{a:128;r:!0;({c:a/x;s:c>63;r::r,c+w:$[s+x%a;a;0];w};{x%a})/:x;_r} /signed leb128 for x>=0
whdr:0x0061736d01000000
/wsig:0x60,...

