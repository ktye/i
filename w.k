table:macro:funcs:extrn:mem:()
afun:{exported:1<+/":"=x; x:"."\:x[&x=":"]:"."; name:`$*x; rtyp:`$x 1; argt:,/`$'""\:x 2;`exported`name`rtyp`argt`body!(exported;name;rtyp;argt;y)}
class:{isnum:{x~^x:,/0+("0";x;"9")};$["!"~ ::x;mem,:`off`data!(.-1_x;."0x",y) ;isnum(*x);table,:`off`funcs!(.-1_x;";"\:y); ":"~ ::x;macro,:`name`repl!(-1_x;y); #y;funcs,:afun[x;y]; extrn,:afun[x;y]];}

trim:{(+/(_10 13 32)=*x)_x}
mf:{r:"{"\:x;(trim/r 0;r 1)}
compile:{class./mf'-1_"}"\:*(_10 92)\:x} 

/ k w.k -fsrc k.w -e 'compile src;'
