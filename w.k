macro:funcs:extrn:(0#`)!();table:mem:()
afun:{exported:1<+/":"=x; x:"."\:x[&x=":"]:"."; name:`$*x; rtyp:`$x 1; argt:,/`$'""\:x 2;name!,`exported`rtyp`argt`body!(exported;rtyp;argt;y)}
class:{isnum:{x~^x:,/0+("0";x;"9")};$["!"~ ::x;mem,:`off`data!(.-1_x;."0x",y) ;isnum(*x);table,:`off`funcs!(.-1_x;";"\:y); ":"~ ::x;macro,:(`$-1_x)!,y; #y;funcs,:afun[x;y]; extrn,:afun[x;y]];}
trim:{(+/(_10 13 32)=*x)_x}
mf:{r:"{"\:x;(trim/r 0;r 1)}
compile:{class./mf'-1_"}"\:*(_10 92)\:x}

expand:{k:!macro;$[(#k)>k?x;macro x;""]}
symbol:...
noun:{x:$[#s:symbol(x);(expand s),x;x]
c:~x;$["("~c;          seq[x;")"]    /(..)
             ~c^"VCIJF";     (`$c;1_x)     /type
	     #s:sym(x);   x:_#s
;]}

cmap:256#0
cmap[&2=((1+"z")>!256)+(-1+"a")<!256]: 1 /[a-z]
cmap[0+"0123456789"]+:                 2 /[0-9]
cmap[0+"VCIJF"]+:                      4 /type
cmap[0+"+-*%."]+:                      8 /op
cmap[!33]+:                           16 /space

/ k w.k -fsrc k.w -e 'compile src;'
/bug: \:2 3 4