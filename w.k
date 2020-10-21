/k -fsrc k.w w.k
/   compile src

compile:{macro:funcs:extrn:(0#`)!();table:mem:();class./:mf'-1_"}"\:*(_10 92)\:x;p:parse[;macro];`macro`funcs`extrn`table`mem!(macro;p'funcs;extrn;table;mem)}
afun:{exported:1<+/":"=x; x:"."\:@[x;&x=":";"."]; name:`$*x; rtyp:`$x 1; argt:,/`$'""\:x 2;name!,`exported`rtyp`argt`body!(exported;rtyp;argt;y)}
class:{isnum:{x~^x:,/0+("0";x;"9")};$["!"~(:x);mem,:`off`data!(.-1_x;."0x",y) ;isnum(*x);table,:`off`funcs!(.-1_x;";"\:y); ":"~(:x);macro,:(`$-1_x)!,y; #y;funcs,:afun[x;y]; extrn,:afun[x;y]];}
trim:{(+/(_10 13 32)=*x)_x}
mf:{r:"{"\:x;(trim/:r 0;r 1)}
parse:{[x;macro]x}

a:compile src

noun:{p:(pTyp;pSym);r:0;({$[~#p;0;r::(*p)x;0;1]};{p:1_p})/:x;((*r)_x;:r)}

first:{0+*x}

pTyp:{h:first x;$[~ctype first x;0;2>#x;0;alphanum 0+x 1;0;(1;(`typ;`$*x))]}
pSym:{h:first x;$[~alpha first h;0;r:(alphanum 0+x)?0;(r;(`sym;`$r#x));0]}

set:{@[256#0;x;1]};num:set"0"+!10;alphanum:num+alpha:set"a"+!26;ctype:set"VCIJF"+0;space:set@!33;op:set"+-*%"+0


\
expand:{k:!macro;$[(#k)>k?x;macro x;""]}
symbol:...
noun:{x:$[#s:symbol(x);(expand s),x;x]
c:~x;i:0+c;$["("~i;          seq[x;")"]    /(..)
             ctype i;     (`$c;1_x)        /type
	     #s:sym(x);   x:_#s
;]}
and:{[f;x]({$[~#f;0;(*f)x;1;0]};{f::1_f;x})/:x;~#f}
set:{@[256#0;x;1]};num:set"0"+!10;alphanum:alpha+num;ctype:set"VCJF";op:set"+-*%";space:set(!33)

