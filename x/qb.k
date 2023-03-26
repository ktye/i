/compile to qbe (c9x.me/compile)   untested
F:`!()                                         /function map: type to stringer
ftyp:`!()
gtyp:`!()
ifun:&#T                                       /function index per node
ltyp:0                                         /ifun!name!type
ssan:ssl:0

qb:{[o]opts::o
 F[`prg`mem`con`var`tab`fun`arg`sym`res`loc`ast`lod`sto]:(qbprg;qbmem;qbcon;qbvar;qbtab;qbfun;qbarg;qbsym;qbres;qbloc;qbast;qblod;qbsto)
 F[`add`and`bnd`ant`asn`cal`cli`cnd`dfr`div`drp`eql`cst]:(qbadd;qband;qbbnd;qbant;qbasn;qbcal;qbcli;qbcnd;qbdfr;qbdiv;qbdrp;qbeql;qbcst)
 F[`for`get`Get`gte`jmp`les`lit`lte`mod`mor`mul`neg    ]:(qbfor;qbget;qbGet;qbgte;qbjmp;qbles;qblit;qblte;qbmod;qbmor;qbmul;qbneg)
 F[`neq`nop`not`orr`bor`xor`ret`shl`shr`stm`sub`swc`typ]:(qbneq;qbnop;qbnot;qborr;qbbor;qbxor;qbret;qbshl;qbshr;qbstm;qbsub;qbswc;qbtyp)

 T::@[T;i@&(1=I j)&`for=T j:P P i:i@&`jmp=T i:lp P;`nop] /delete continue in simple loop
 qbren'`main`ldexp`hypot                       /rename libm
 /i:i@&{(0N~|/(T 1_c)?`swc`stm)&3=#c:&P=x}'i:&T=`cnd  /keep ternary
 rwcnd'|&(~`=S)&(T=`cnd)|T=`swc                /rewrite cnd/swc expressions 
 S[i]:`$"f",/:$i:i@&`=S i:&T=`for              /label for loops
 rwlab'i@&`=S i:&T=`jmp                        /set jmp labels
 S[0]:$[0N~i:o?`lib;S 0;o 1+i]                 /emit`lib`name (instead of main + remove func main)
 ftyp::(S i)!qbsig'i:&T=`fun                   /function signatures "rxyz"
 ftyp::ftyp,deftyp
 ifun::-1+\T=`fun                              /function index for each node
 ltyp::(!g)!((S 1+i)@.g)!'(S i)@.g:=ifun i:&(T=`arg)|T=`loc
 n:(&#P){x[y]+:1}/P;n[0]-:1;                   /number of children per node
 s:{c:(0,(#x)-z)^x                             /cut args from stack
   (c 0),,,/(F T y)[y;|c 1]}/[();|!#P;|n]      /fold over the tree in reverse order
*s}


hx8:{[i]"0x",`x@|C i+!8}
lp:{(~~)#@[&#x;x i;i:!#x]}
em:{"(",x,")"}

INS:{[i;x]T::ins[i;T;*x];P::ins[i;@[P;&P>i-1;+;1];x 1];P[1+i]:i;I::ins[i;I;x 2];S::ins[i;S;x 3]}
ins:{[i;x;y](i#x),y,i_x}
del:{T::T n:(!#P)_x;P::@[P n;&(P n)>x-1;-;1];I::I n;S::S n}
rwcnd:{c:|1_&P=x;S[x]:`;P[x]:P x-1;c INS\:(T x;1+x;I x; S x:x-1);del x;restm'|i@&(T P i:&`stm=T)':`ret`asn}
restm:{c:*|&P=x;INS[c;(T x-1;x;I x-1;S x-1)];P[x]:P[x-1];del[x-1]} /fix ret stm

qbren:{S[&(S=x)&T':`fun`cal]:s:`$"_",$x;}
rwlab:{ S[x]:S@{$[`for~T x;x;P x]}/:x }


NL:{x,"\n"}
nl:{$[#x:"\n"/:x;x,"\n";x]}
nyi:{y;em x}

qbtype:{,(`i`u`j`k`f!"wwlld")x}

qbprg:{
 fncs: ,/y@&`fun=T c:1_&P=0
 cons:nl y i:&(`con=T c)|`var=T c
 gtyp::(S i)!(S 1+i)
 runt:$[`nort':opts;"";qbrnt]
 main:$[|/`nort`lib':opts;"";qbmain]
 data:$[#D;"%t =l call memcpy(%_M, $data0, ",($#D), ")\n";""]
 dat0:$[#D;"data $data0 = { b ",("\"\\x","\\x"/:(#D)^`x D),"\" }\n";""]
 tabl:,/y@&`tab=T c
 tabl:$[#tabl;"%_F =l call $malloc(",($8*1+|/I@&`tab=T),");\n";""],tabl
 mems:nl y@&`mem=T c
 init:$[|/(#mems;#data;#tabl);nl("function init(){\n@start\n";mems;data;tabl;,"}");""]
 (runt;cons;dat0),(,fncs),init,main}
 
qbfun:{
 $[(`main~S x)&`lib':opts; :""; 0]
 ssan::ssl::0
 args:", "/:y@&`arg=T c:&x=P
 ret:$[#r:", "/:y@&`res=T c;r;""]
 /locs:nl@{(qbtype x)," ",y}'[!g;.g:","/:'(y i)@=(S c)i:&`loc=T c]
 defr:     *y@&`dfr=T c
 setj:$[#defr;"if(!setjmp(_jb_)){";""]
 body:qbfmt$[#r:y@*&`ast=T c;r;""]
 ,/("function ",ret;" $";$S x;"(";args;"){\n@start\n";setj;body;defr,"\n}\n")}

qbsig:{,/(qbtype 1#S@*&(T=`res)&P=x),qbtype S@&(T=`arg)&P=x} /"rxyz"
qbmem:{y;("Memory";$[`b~S x;"2";""];"(";$I x;");")}
qbtab:{y;"%fi =l add %_F, ",($8*I x),"\nstorel $",($S x),", %fi\n"}

qbflts:{d:{(0=*x)_x:@[@[x;i;-;10];-1+i:&9<x:@[0,2*x;#x;+;y];+;1]};,/$d/,/+(8#2)\\256/256+|x}
qbflt:{x:C x+!8;$[#r:qbflts x;r;"d_",$*`f x]}
qblit:{hi:{*`i C 4+x+!4};lo:{*`i C x+!4};y;i:I x
 $[`f~t:S x;qbflt i;(`u~t)|`i~t;$[0N~i;"2147483648";$i];i~0;,"0";qbflts C i+!8]}

qbvar:qbcon:{"data $",($S x)," = { ",(qbtype S 1+x)," ",(*y)," }"}

/          i:iujkf   u:iujkf   j:i     u     jkf   k:i     u     jkf   f:i u j     k     f
qbcast:","\:",,,,dtosi,,,,,dtoui,extsw,extsw,,,dtosi,extuw,extuw,,,dtoui,s,u,sltof,ultof,"
qbcst:{c:qbcast(5*`i`u`j`k`f?S x)+`i`u`j`k`f?S 1+x; t:qbtype S x
 $[""~c;y:y 1;("u"~*c)|"s"~*c;ssa1["d ",(*c),"ltof";ssa1["l ext",(*c),"w";y 1]];ssa1[t," ",c;y 1]]}


qbtyp:{y;qbtype S x} 
qbres:{y;qbtype S x}
qbloc:{y}
qbsym:{y;"%",$S x}
qbarg:{(qbtype S x;" ";*y)}


qbstm:{,/y}
qbast:qbstm
qbget:{y;"%",$S x}
qbGet:qbget
qbret:{("ret "/:ssa y 0),"\n"}

qbsiz:`b`i`j`f!`I8`I32`I64`F64
qblod:{("?=";(`b`i`j`f!"wwlf")S x;" call ";$qbsiz S x;,"(";y 0;")\n")}
qbsto:{y:ssa'y;(,/*'y),"call $Set",($qbsiz S x),"(",(", "/:y[;1]),")\n"}

/todo: ftyp for builtins?
qbdrp:{*ssa@*y}
qbcl:{[f;rt;at;y] y:ssa'y; a:{x," ",y 1}'[at;y];  (,/*'y;$[" "~rt;"call ";"?=",rt," call "];f;,"("; ", "/:a ;")\n")}
qbcal:{$[`I32B~S x;*y;qbcl["$",$S x;*t;1_t:ftyp S x;y]]}
qbcli:{a:ssa@"?=l loadl "/:ssa NL@"?=l add %_F, "/:ssa NL@"?=l mul 8, "/:ssa y 0; n:I x;  (*a),"\n",qbcl[*|a;qbtype S x;(1+n)_*'y;n#1_y]}

qbasn:{ ("%",($S x)," "),/ $[1=#y:"?"\:y 0;("";"=",$[I x;gtyp S x;qbltyp x]," copy ",y,"\n");y] }
qbltyp:{qbtype(ltyp ifun x)S x}

qbcnd:{l:$ssl+:1;("jnz "/:ssa@*y),", @i",l,", @e",l,"\n@i",l,"\n",(y 1),$[3~#y;"\n"/:("jmp @f",l;"@e",l;(y 2),"@f",l);"@e",l],"\n"}
qbnop:{y;""}

qbswc:{e:ssa@*y;d:$[I x;*|y;""];n:#y:(-I x)_1_y;  c:{"?=w eq ",x,", ",($y),"\n"}/:[e 1;!n]; 
 (*e),d{qbcnd[0; y,,x ]}/|{(,x),,y}'[c;y]    }

qbdfr:{("else{";*y;";}\n}")}

qbfor:{l:"@",($S x);le:l,"e"; l,"\n",$[#*y; ("jnz "/:ssa y 0),", ",lc,", ",le,"\n",(lc:l,"c"),"\n";""],(y 2),(y 1),"jmp ",l,"\n",le,"\n"}
qbjmp:{y;"jmp @",($S x),(" e"I x),"\n"}

qb1:{x,*z}
qb2:{z:ssa'z;(,/*'z),("?=",qbtype S y)," ",x," ",(", "/:(*1_)'z),"\n"}
qbc:{z:ssa'z;(,/*'z),"?=w ",("su"@(S y)':`u`k),x," ",(", "/:(*1_)'z),"\n"}

qbneg:{ NL@("?=",(qbtype S x)," sub 0, ")/:ssa y 0 }
qbnot:qb1"!"
qbadd:qb2"add"; qbsub:qb2"sub"; qbmul:qb2"mul"; qbdiv:qb2"div"; qbmod:qb2"rem"
qbeql:qb2"eq";qbgte:qbc"ge";qbles:qbc"lt"; qblte:qbc"le";qbmor:qbc"gt"; qbneq:qb2"ne"
qband:qb2"&"; qbbnd:qb2"&&";qbant:qb2"&~";qborr:qb2"|"; qbbor:qb2"||";qbxor:qb2"^"; qbshl:qb2"<<";qbshr:qb2">>"

qbnd:{$["!=0)"~t:-4#x;(-4_x),")";"==0)"~t;"(!",(1_-4_x),")";x]}
qbs:{x; \`QBS}

ssa:{$["?"':x;((r," ")/:"?"\:x; r:"%.",$ssan+:1);("";x)]}
ssa1:{ (("?=",x," ")/:ssa y),"\n" }

qbfmt:{"\n"/:" ",/:"\n"\:-1_x}

deftyp:`Memorycopy`Memorycopy2`Memorycopy3`Memoryfill!4#," www"
deftyp,:`Memorysize`Memorysize2`Memorygrow`Memorygrow2`Exit!(,"w";,"w";"ww";"ww";" w")
deftyp,:`F64abs`F64sqrt`F64ceil`F64floor`F64nearest!5#,"ff"
deftyp,:`F64min`F64max`F64copysign!3#,"fff"
deftyp,:`"F64reinterpret_i64"`"I64reinterpret_f64"`I32clz`Native`panic!("fl";"lf";"ww";"lll";" w")
deftyp,:`Args`Arg`Read`Write`ReadIn!1 3 4 5 3#'"w"

/ run time (missing)
qbrnt:""
qbmain:""

