table:()  /0:{add;sum}   numeric hdr
macro:()  /r8:{rp:r+8}   no types
funcs:()  /add:I:II{x+y} hdr with return type and value sum.I:I (unexported)
extrn:()  /sin:F:F{}     empty body

afun:{exported:1<+/":"=x; x:"."\:x[&x=":"]:"."; name:*x; rtyp:x 1; argt:""\:x 2
 `exported`name`rtyp`argt`body!(exported;name;rtyp;argt;y)}

class:{isnum:{x~^x:,/0+("0";x;"9")}
     $[isnum(*x);table,:`off`funcs!(.-1_x;";"\:y)
       ":"~::x  ;macro,:`name`repl!(-1_x;y)
       #y       ;funcs,:afun[x;y];extrn:afun[x;y]]}

trim:{(+/(_10 13 32)=*x)_x}
mf:{r:"{"\:x;`hdr`body!(trim/r 0;r 1)}
compile:{mf'-1_"}"\:*(_10 92)\:x} 
/wasm:{[x]0x0061736d01000000}
wasm:{x}
