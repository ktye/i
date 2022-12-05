/ ik transforms the IR (a k table) into an "i" program.
/ "i" is a low level input language with k syntax
/  it can be parsed by k but not evaluated (see ki.k).
/ types
/  `i`u`j`k`f:  i32 u32 i64 u64 f64
/ heap
/  I x  U x  J x  K x  F x  read. x is byte address(not vector index) in linear memory
/  I[x]:y ..                write. alignment is assumed, e.g. avoid I[1+2]:..
/ functions
/  f[i]:{[xi;yf]..}      f(name), i(return type), args&locals have type suffix in the name
/                        locals are not declared, they are detected by k on assignment, global: zi::..
/  f[i]:{[]..; :x; ..}   return x, automatic on last statement if f has a return value.
/  f[ii]:{[]..;(,a;,b)}  multi-return
/ globals
/  ai:3                  global var (mutable)
/  bi=4                  const
/ literals
/  xi:1; yf:2.           int float parsed by k
/  2 5                   (two int array) is int64 [high low]
/  0x123...              precise float constants binary 64bit
/ cast
/  `i x  `u x  ..
/ call
/  f x  f[x;y]           normal k syntax
/  g f x                 f returns same number of values as g has args
/  `a`b:f x              destructing assignment
/  `a``c:f x             empty var drops a return value
/ call-indirect          
/  .[int-expr;ret-type;arg-types][args]    .[5;i;iif][1;2;3.0]
/ loops
/  (expr;post)'[statement-list]            (x<9;x:x+1)'[a;b;c]
/  (;)'[..]              while true
/  (;;label)'[..]        jump label for nested loops
/  (;)'[..;<;..]         jump backwards(continue)
/  (;)'[..;>;..]         jump forwards(break inner loop)
/  (;)'[..;>label;..]    break/continue label (outer loop)
/  (;)'[..;<]            simple loop (continue is last in list)
/  simple loops don't check cond on first entry (repeat-until)
/ cond
/  $[x;y]                if x then y
/  $[x;y;z]              if x then y else z
/  $[a;b;c;d]            if a then b else if c then d
/ switch
/  ?[int-expr;a;b;..]    jump table: switch(int-expr)case 0:a; case 1:b; default:c
/  ?[int-expr;a;b;]      last expr is empty: no default
/  cond and switch may be expressions or statements.
/  as expressions, e.g. 1+r:$[x;y;z] cond needs an else clause and switch a default,
/  and branched must evaluate to the same type.
/ data
/  0x1234..              first statements are data section(copied to initial memory at 0)
/ function table
/  0:`a`b`c              table for indirect calls start-index:names
/  50:`f`g ";"/:

F:`!()
hx8:{[i]"0x",`x@|C i+!8}
lp:{(~~)#@[&#x;x i;i:!#x]}


ik:{[o]opts::o
 /T::@[T;i@&(1=I j)&`for=T j:P P i:i@&`jmp=T i:lp P;`nop] /delete continue in simple loops
 
 {[d;i;t]$[t':`arg`loc;d[S 1+i]:S i;t':`sym`get;S[i]:`$,/$(S i),d S i;0];d}/[(0#`)!0#`;!#T;T] /append type to symbols
 
 F[`prg`mem`con`var`tab`fun`arg`sym`res`loc`ast`lod`sto]:(ikprg;ikmem;ikcon;ikvar;iktab;ikfun;ikarg;iksym;ikres;ikloc;ikast;iklod;iksto)
 F[`add`and`bnd`ant`asn`cal`cli`cnd`dfr`div`drp`eql`cst]:(ikadd;ikand;ikbnd;ikant;ikasn;ikcal;ikcli;ikcnd;ikdfr;ikdiv;ikdrp;ikeql;ikcst)
 F[`for`get`Get`gte`jmp`les`lit`lte`mod`mor`mul`neg    ]:(ikfor;ikget;ikGet;ikgte;ikjmp;ikles;iklit;iklte;ikmod;ikmor;ikmul;ikneg)
 F[`neq`nop`not`orr`bor`xor`ret`shl`shr`stm`sub`swc`typ]:(ikneq;iknop;iknot;ikorr;ikbor;ikxor;ikret;ikshl;ikshr;ikstm;iksub;ikswc;iktyp)
 
 S[0]:$[0N~i:o?`lib;S 0;o 1+i]
 n:(&#P){x[y]+:1}/P;n[0]-:1;
 s:{c:(0,(#x)-z)^x
   (c 0),,,/(F T y)[y;|c 1]}/[();|!#P;|n]
*s}


ikprg:{
 c:1_&P=0
 excl:!0
 prog:"`",($S 0)," /program\n"
 fncs:"\n"/:y@&`fun=T c
 cons:$[#r:ikfmt@";"/:y@&`con=T c;"/constants\n",r;""]
 vars:$[#r:ikfmt@";"/:y@&`var=T c;"/globals\n",r;""]
 data:$[#D;"/data\n","\n"/:("0x",`x@)'-32^D,_&31;""]
 tabl:iktbl y t:&`tab=T c
 "\n"/:(""~/:)_(prog;data;tabl;cons;vars;"\n/functions\n",fncs)}
 
ikdat:{,/@["\\",/:256^`x@_!256;0+ikc;ikc:_("a"+!26),("A"+!26),("0"+!10),"~!@#$%^&*()_+{}[]|;:<>,./?"]x}
ikmem:{y;""}
ikcon:{($S x),"=",*y}
ikvar:{($S x),":", *y}
iktab:{y;(I x;S x)}
iktbl:{"\n"/:{($x),":",y}'[(*x)i;,/'"`",''$(i:&~1=-':*x)^(x:+x)1]}
ikfun:{s:$S x
 args:";"/:y@&`arg=T c:&x=P
 rets:   ,/y@&`res=T c
 locs:" "/:y@&`loc=T c
 body:    y@*&`ast=T c
 expo:$[I x;"(export \"",s,"\")";""]
 s,"[",rets,"]:{[",args,"]",ikfmt body,"}"}
 
 
ikarg:{$[#y;$*y;$S x]}
iksym:{y;$[`"_"~s:S x;"";$s]}
ikres:{y;,/$S x}
ikloc:{"(local ",(*y)," ",(iktype S x),")"}
ikstm:{";"/:y}
ikast:{";"/:$[":"~**|y;y[l]:1_y[l:-1+#y];y]}

iklod:{$[`b~s:S x;"B";_($s)-32]," ",*y}
iksto:{$[`b~s:S x;"B";_($s)-32],"[",(*y),"]::",y 1}

ikant:{(y,'"\n"),t,".const -1\n",t,".xor\n",(t:iktype S x),".and"}

ikasn:{$[1~i:(I x);*y;,/"`",'i#y],":",*|y}
ikcal:{n:#y:";"/:y;$[`panic~s:S x;"!",y;($s),$[1~n;" ",y;"[",y,"]"]]}

ikcli:{n:I x;".[",($S x),";",(","/:n#(1+n)_y),";",(*y),"]",$[1~n;y 1;"[",(";"/:n#1_y),"]"]}


ikdfr:{y;""}
ikdrp:{y}
ikcst:{$[(S x)~y 0;y 1;"`",($S x)," ",y 1]}

ikfor:{"(",(";"/:2#y),$[s:S x;";",$s;""],")'[",(y 2),"]"}

ikget:{y;$S x}
ikGet:{y;$S x}
ikjmp:{y;(,$[I x;"<";">"]),$[`~S x;"";($S x)]}

ikflt:{$[f~0.+_f:*`f@c:C x+!8;$f;"0x",`x@`c@,f]}
iklit:{y;i:I x;$[`f~t:S x;ikflt i;t':`i`u;$i;i~0;"0 0";" "/:$`i C i+!8]}

ikneg:{"-",*y}
iknot:{(*y),"\ni32.eqz"}
iknop:{y;""}
ikret:{$[1~#y; ":",*y; ":(,",(";,"/:y),")"]}

ikswc:{"?[",(";"/:y,$[I x;();,""]),"]"}
ikcnd:{"$[",(";"/:y),"]"}

iktyp:{y;S x} /keep as symbol

ik2: {$[(T 1+y)':`get`Get`lit;z 0;"(",(z 0),")"],x,z 1}

ikadd:ik2"+";iksub:ik2"-"; ikmul:ik2"*";ikdiv:ik2"/"; ikmod:ik2"%"
ikand:ik2"&";ikbnd:ik2"&&";ikorr:ik2"|";ikbor:ik2"||";ikxor:ik2"^"
ikneq:ik2"~";ikeql:ik2"="
ikmor:ik2">";ikles:ik2"<";ikgte:ik2">=";iklte:ik2"<="
ikshl:{ik2["<<";x;|y]};ikshr:{ik2[">>";x;|y]}

iktype:`i`u`j`k`f!("i32";"i32";"i64";"i64";"f64")
iksign:`i`u`j`k`f!("_s";  "_u"; "_s"; "_u";   "")
ikinst:`I32B`Memorysize`Memorygrow`Memorycopy`Memoryfill!("nop";"memory.size";"memory.grow";"memory.copy";"memory.fill")
ikinst,:`Memorysize2`Memorygrow2`Memorycopy2`Memorycopy3!((,"unreachable\ni32.const 0")0 0),(,"drop\ndrop\ndrop\nunreachable")0 0
ikinst,:{x!{(_32+*x),(1_3#x),".",3_x}'$x}`I32clz`F64abs`F64sqrt`F64floor`F64copysign`F64min`F64max`"F64reinterpret_i64"`"I64reinterpret_f64"


/format: if longest line is too long, try to split it at lowest level (until fixed)
ikfmt:{l:1+(+\("["=x)|"("=x)-+\("]"=x)|")"=x  /indentation level
 x:@[x;i;_l i:&";"=x]                         /replace ";" with level(byte) spl
 "\n"/:iksem'ikspl/:,x}
ikspl:{$[100<n:n i:*&n=|/n:#'x;(x@!i),(ikf2 x i),(1+i)_x;x]}
ikf2:{i:&x=l:&/x i:&" ">x; (,(*i)#x), (_32+&0+l),/:1_'(i)^x}
iksem:{@[x;&" ">x;";"]}                       /remove binary levels

/ikfmt:{l:1+(+\("["=x)|"("=x)-+\("]"=x)|")"=x:";",x;"\n"/:((_(32)+&)'l i),'1_'(i:&";"=x)^x} /split everything and indent


/ `<ikfmt"alpha;1+(gamma;delta)*2;epsilon"
