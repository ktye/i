
O:10/8\65536!
D:O@+/1 256*256!+-2^@[`c@&2^14;!#D;D]  /convert data section to uint16 (and fill)
F:(0#`)!()
ifun:&#T                /function index of each node
afun:(0#`)!!0           /fname!-addr  first is 0(main), never called
locs:0            
atyp:anam:()            /argtypes  argnames for each function
ltyp:lnam:()            /loctypes  locnames for each function (args+locs)
cons:(0#`)!()           /constant!literal-offset
offs:0                  /`symbol!offset for get/set args+locals for each function
lits:()
rtyp:(0#`)!!0

/print table
/print:{`<(`l@x_+`n`t`p`i`s!(!#T;T;P;I;S)),\"\n"}
/stack at 2559 word offset (5118 011776 bytes) within data

pd:{[o]opts::o
 F[`prg`mem`con`var`tab`fun`arg`sym`res`loc`ast`lod`sto]:(pdprg;pdmem;pdcon;pdvar;pdtab;pdfun;pdarg;pdsym;pdres;pdloc;pdast;pdlod;pdsto)
 F[`add`and`bnd`ant`asn`cal`cli`cnd`dfr`div`drp`eql`cst]:(pdadd;pdand;pdbnd;pdant;pdasn;pdcal;pdcli;pdcnd;pddfr;pddiv;pddrp;pdeql;pdcst)
 F[`for`get`Get`gte`jmp`les`lit`lte`mod`mor`mul`neg    ]:(pdfor;pdget;pdGet;pdgte;pdjmp;pdles;pdlit;pdlte;pdmod;pdmor;pdmul;pdneg)
 F[`neq`nop`not`orr`bor`xor`ret`shl`shr`stm`sub`swc`typ]:(pdneq;pdnop;pdnot;pdorr;pdbor;pdxor;pdret;pdshl;pdshr;pdstm;pdsub;pdswc;pdtyp)
 F[`i32b]:i32b
 
 rwcnd'|&(~`=S)&(T=`cnd)|T=`swc
 $[1~#i:i@&`store=S P i:&T=`ast;T[&P=*i]:`nop;0]  /remove store/catch
 $[1~#i:i@&`catch=S P i:&T=`ast;T[&P=*i]:`nop;0]

 cons::S@&T=`con
 afun::s!-!#s:S@&T=`fun   
 
 T[i:&(S=`I32B)&T=`cal]:`i32b
 ifun::1_+\-1,T=`fun                              /function index for each node
 perfunc:!1+*|ifun
 fill:{x[i]:(#i:&~#'x@!n:+/T=`fun)#,!0}
 a:fill i(=ifun i:i@&`fun=T P i:&T=`arg)
 atyp::(S   a)perfunc     /argtypes
 anam::(S 1+a)perfunc     /argnames
 a:fill i $[#j:ifun i:&T=`loc;=j;(!0)!()]
 ltyp::(S   a)perfunc     /loctypes
 lnam::(S 1+a)perfunc     /locnames                                     
 offs::({(|x)!$[#r:8*y+!#x;-8,1_r;!0]}'[anam;#'lnam]) /  main   ; xyz:args ab:locals
 `"k.off"<","/(($S@&T=`fun),\":"),'$#'anam            /k.off: number of args per func

 @[`offs;i;,';{x!8*!#x}'lnam(i:&~~#'lnam)]            /((0#`)!!0;`z`y`x`a`b!-8 16 24 0 8;..)

 rtyp:(!afun)!(#afun)#`
 rtyp[S P r]:S r:&T=`res                              /funcs with empty results have no `res node

 glob::$[#s:S@&T=`var;s!O 5118-8*!#s;(0#`)!!0]
 lits::$[#glob;-4^&4*#glob;()];T[1+&T=`con]:`nop      /initial values for globals are assigned later
 cons::$[#cons;cons!pdlit'[1+&T=`con;0];(0#`)!!0]

 n:{x[y]+:1}/(,&#P),P;n[0]-:1
 s:{c:(0,(#x)-z)^x;(c 0),,,/(F@ T y)[y;|c 1]}/[();|!#P;|n]
 lits[!#glob]:litl'[I i;S i:1+&T=`var]
"0",",\n0"/$imag stak link@*s}

link:{o:(!builtins)!512+2*-1_+\(#x),#'b:.builtins;@[x;i;O o x i:&x<0],,/b}
stak:{D[(2558)-!n:4*#lits]:,/lits;x}   /embed literals/constants/gobals in stack
imag:{@[&(2^15)-256;!`size \#x;x],D}   /fill disk image: code(32kw-01000) + data(32kw)

INS:{[i;x]T::ins[i;T;*x];P::ins[i;@[P;&P>i-1;+;1];x 1];P[1+i]:i;I::ins[i;I;x 2];S::ins[i;S;x 3];S}
ins:{[i;x;y](i#x),y,i_x}
del:{T::T n:(!#P)_x;P::@[P n;&(P n)>x-1;-;1];I::I n;S::S n;S}
rwcnd:{c:|1_&P=x;S[x]:`;P[x]:P x-1;c INS\(T x;1+x;I x; S x:x-1);del x;restm'|i@&(T P i:&`stm=T)?\`ret`asn}
restm:{c:*|&P=x;INS[c;(T x-1;x;I x-1;S x-1)];P[x]:P[x-1];del[x-1]} /fix ret stm


pdprg:{
 c:1_&P=0;f:y fn:&`fun=T c
 tabl:$[#i:&`tab=T c;62707,(O 32+4*#i),10004 16600 02 16601 04 16602 06 16603 10 11666 10 62706 10 6304 6304 60407,,/y i;!0]
 init: tabl,12706,O 5118-8*#lits             /mov #11776,sp (init sp: 1k reserved for kos) -literals
 
 d:(- 1_!#f)!O 512+2*1_-1_+\(#init),#'f      /-1 -2..!function start addresses
 ($1_(S c)fn),$.d
 `"k.map"<"({",(",\n"/(($1_(S c)fn),\":"),'$.d),"})"
 @[r;i;d r i:&(r>-#f)&0>r:init,,/f]}
 


pdsig:{*y}
pdmem:{*y}
pdcon:{*y}
pdvar:{*y}
pdtab:{y;127,afun S x}
pdtbl:{*y}

pdfun:{
 c:&x=P
 body:y@*&`ast=T c
 n:O 4+8*+/`loc=T c
 m:0|8*-1++/`arg=T c
 pro:10546 162706,n,10605            /mov fp,-(sp); sub #n(sp); mov sp,fp
 epi:10506 62706,n,12605 12604       /mov fp,sp   ; add #n,sp ; mov (sp)+,fp; mov (sp)+,r4; add 8*-1+#args,sp
 epi,:$[m;62706,O m;!0],104          /add #m,sp   ; jmp r4
 $[`main~S x;body,127 0;ret[pro,body;epi]]}
 
ret:{j:(#x)-2+i:&x=-66666; x[i]:62707;x[1+i]:O 2*j ;x,y}  /replace -66666 0 with add #n,pc (jump to epilog) 
 
pdarg:{*y}
pdsym:{y;S x}
pdres:{y;S x}
pdloc:{y;S x}
pdstm:{,/y}
pdast:pdstm

pdlod:{(*y),4727,-(70000+`b`i`j?S x)}
pdsto:{(,/|y),4727,-(60000+`b`i`j?S x)}

pdasn:{(y 0),$[I x;12746,(glob S x),4727 -75000;(set setbuiltin@(offs ifun x)S x)]}
pdget:{y;get getbuiltin@(offs ifun x)S x}
pdGet:{y;$[`Lt~S x;4727 -55523;(S x)?!cons;cons S x;10346 12703,(glob S x),4727 -77000]}

pdcal:{$[#y;(,/y);4727 -99999],4727,$[`panic~S x;0;0~r:pdinst S x;(afun S x),$[`~rtyp S x;po;!0];r]}
pdcli:{$[1~#y;4727 -99999;,/1_y],(*y),4727 1004}    /there are no indirect calls of funcs without results
pddfr:{y;!0}

pddrp:{(*y),po}
pdcst:{s:*y;d:S x;(y 1),$[(pdw d)>pdw s;4727,-50007-s~`u;!0]}

pdfor:{(pdlop;pdslp)[I x].(x;y)}                       /slp:do-while; lop:while-do
pdlop:{p:$[#*y;po;!0];(cnt brk $[#*y;(y 0),5700 1002 62707,(O@2*#b);!0],b:(p,,/|1_y),-44444 0),p}
pdslp:{c:(*y),5700 1402 162707;b:po,cnt brk(-2_y 2),y 1;  pu,b,c,O 2+2*(#b)+#c}

brk:{j:(#x)-2+i:&x=-55555;x[i]:62707;x[1+i]:O 2*j;x}   /replace -55555(break)    with add #n,pc
cnt:{j:2+i:&x=-44444;x[i]:162707;x[1+i]:O 2*j;x}       /replace -44444(continue) with sub #n,pc
pdjmp:{y;,$[I x;-55555 0;-44444 0]}                    /break/continue (no nested loops)


/literals, constants and globals are embedded in the stack
pdlit:{y;i:I x;$[2~pdw s:S x;lit1 I x;0x00000000~C i+4+!4;lit1@*`i C i+!4;10346 12703,lit[i;`k],4727 -77000]} /push r3;mov #literal-offset to r3
lit1:{$[0~x; pu;1~x; pu,5200;-1~x;pu,5300 5301;8~x;4727 -55508;16~x;4727 -55516;(x<32767)&x>0; pu,12700,O x; 10346 12703,lit[x;`i],4727 -77000]}
litl:{$[`i~y;O@,/(`i@,/(2^`c@,x),\0x0000),0 0;O@,/|'+65536 65536\(`i@C x+!4),`i@C x+4+!4]}
lit: {O 5118-8*$[0N~r:lits?x:litl[x;y]; -1+#lits,:,x ;r]}


pdneg:{(*y),5401 5400 5601}             /neg32:neg r1;neg r0;sbc r1
pdnot:{(*y),0+`not}                     /should not exist
pdnop:{y;!0}
pdret:{(*y),$[`ast~T P x;!0;-66666 0]}

pdswc:{b:(1_y),\-55555 0;brk@-2_,/(!0
 (*y)                                /switch-expr in r0
 22700,O@-1+#b                       /cmp #(n-1),r0
 101402 62707,O@2*#*|b               /jump over default
 *|b                                 /default
 70027 4 10100 60007                 /mul #4,r0;mov r1,r0;add r0,pc (jump to case offset)
 ,/62707,/'O(4*|-1_!#b)+2*-2_0,+\#'b /branch offsets
 ,/-1_b                              /cases
 po)}                                /pop switch expr


i32b:{(*y),5001}
pdcnd:{(iff;ife)[3~#y].(x;y)}
iff:{$[(T 1+x)?`bnd`bor;(*y),5700 1401 402;(-4_*y)],62707,(O@2*#y 1),(y 1),po}
ife:{$[(T 1+x)?`bnd`bor;(*y),5700 1401 402;(-4_*y)],62707,(O@2*#b),(b:(y 1),62707,O@2*#y 2),(y 2),po}

pd2:{(,/|z),(4727;-50000-x)}
pdand:pd2 13;pdorr:pd2 14;pdxor:pd2 15;pdant:pd2 16
pdbnd:{(*y),5700,(O 768+#r),r:12600 12601 12602 12603,y 1} /y0&&y1: y0;tst r0;bne +#r;po;y1
pdbor:{(*y),5700,(O 512+#r),r:12600 12601 12602 12603,y 1} /y0||y1: y0;tst r0;beq +#r;po;y1

pdc:{(,/z),(4727;$[4~pdw t;-50022;`u~t;-50023;-50021]),ib(x(t:S y)?`u`k)}
pdneq:pdc  512 512  ;pdeql:pdc  768 768                 /bne beq
pdles:pdc 1280 34560;pdmor:pdc 1536 33280               /blt blo  bgt bhi
pdlte:pdc 1792 33536;pdgte:pdc 1024 34304               /blt blos bge bhis
nc:{(-1_x),$[512!c:*|x;c-256;c+256]}                    /negate condition +/- 0400
ib:{(O x+2),5000 402 5000 5200}

pdtyp:{y;S x}

pdadd:{(,/y),4727 -50000}
pdsub:{(,/|y),4727 -50001}
pdmul:{$[2~pdw S x;(,/y),4727 -50002;0+`nyimul+ \y]}
pdshl:{$[2~pdw S x;(,/|y),4727 -50010; c32~y 1; (*y),4727 -50017; c59~y 1;(*y),4727 -50019;0+`nyishll+`shl \y]}
pdshr:{$[`i~S x;(,/|y),4727 -50011;c32~y 1;(*y),4727 -50018;c59~y 1;(*y),4727 -50020;0+`nyishrl+ \y]}
pddiv:{(,/|y),4727 -50003}
pdmod:{(,/|y),4727 -50004}

c59:4727 -99999 12700 73   /constant shift arguments
c32:4727 -99999 12700 40
c21:4727 -99999 12700 25

pdw:`i`u`j`k`f!2 2 4 4 4   /type widths
pdtype:`i`u`j`k`f`!0x7f7f7e7e7c00

pdinst:{$[0N~`Memorysize`Memorygrow`Memoryfill`Memorysize2`Memorygrow2`Memorycopy2`Memorycopy3`F64abs`F64sqrt`F64floor`F64copysign`F64min`F64max`"F64reinterpret_i64"`"I64reinterpret_f64"?x;0;0+`nyi]
 $[`Memorycopy~x;-55000;`I32clz~x;-55001;`Write~x;-55002;`ReadIn~x;-55003;0]}


pu:4727 -99999            /jsr pc,#push
po:4727 -88888            /jsr pc,#pop
get:4727,-10000-          /e.g. -9996..-10002
set:4727,-20000-

readin:12604 11600 10001 62706 10 105727 177560 100375 112721 177562 1375 160100 5400 5300 5001 104
write: 12604 11601 62706 20 5700 1404 112127 177566 77003 5000 5001 104
memcpy:12604 10002 12600 16601 6 5702 1403 112021 5302 773 62706 16 12600 12601 12602 12603 104
clz:   12604 10102 10003 5000 5702 1003 62700 20 10302 32702 100000 1003 5200 6302 772 5001 104

/               -44444(break) -55555(continue) -66666(return)
builtins:!/+-2^(-99999;12604,(10046+100*|!4),(5000+!4),104  /push4: mov (sp)+,r4;push r3..r0;jmp r4
                -88888;12604 12600 12601 12602 12603 104    /pop4 : mov (sp)+,r4;pop  r0..r3;jmp r4
		-77000;12604 10246 10146 10046 14300 14301 14302 14303 104 /get literal: store offset (rel to stack bottom) in r3;push r0..r3;load r0..r3 from stack
                -75000;16604 2 10044 10144 10244 10344 12604 62706 2 12600 12601 12602 12603 104  /global assign: mov 2(sp),r4;..
		-70000;12604 111000 6701 104                /load`b:movb (r0),r0;sxt r1
		-70001;12604 10001 11000 16101 2 104        /load`i:mov r0,r1;mov (r0),r0;mov 2(r1),r1
		-70002;12604 10003 11000 16301 2 16302 4 16303 6 104   /load`j
		-60000;12604 111610 62706 10 12600 12601 12602 12603 104  /store`b:movb (sp),(r0);add #8,sp
		-60001;12604 10003 12623 12623 62706 4 12600 12601 12602 12603 104       /store`i:mov r0,r3;mov (sp+),(r3)+;..;add #4;sp
		-60002;12604 10003 12623 12623 12623 12623 12600 12601 12602 12603 104   /store`j:mov r0,r3;mov (sp)+,(r3)+;..
		-55523;12604 10346 10246 10146 10046 12700 27 5001 5002 5003 104 /get constant Lt(23) used 70 times
		-55516;12604 10346 10246 10146 10046 12700 20 5001 5002 5003 104 /get literal 16 used 94 times
		-55508;12604 10346 10246 10146 10046 12700 10 5001 5002 5003 104 /get literal  8 used 90 times
		-55003;readin
		-55002;write
		-55001;clz
		-55000;memcpy
		-50000;12604 62600 5501 5502 5503 62601 5502 5503 62602 5503 62603 104  /add: mov (sp)+,r4;add (sp)+,r0;adc r1;add (sp)+,r1;add #4,sp;jmp r4
                -50001;12604 162600 5601 5602 5603 162601 5602 5603 162602 5603 162603 104 /sub: mov (sp)+,r4;sub (sp)+,r0;sbc r1;sub (sp)+,r1;add #4,sp;jmp r4
		-50002;12604 10102 10046 70066 2 70266 2 60300 16602 4 70226 60300 10002 10100 10201 62706 10 104 /mul32
		-50003;12604 12602 12603 62706 4 71002 104              /idiv(cheating slightly)
		-50004;12604 12602 12603 62706 4 71002 10200 10301 104  /irem
		-50007;12604 5701 6702 6703 104             /int64: mov (sp)+,r4;tst r1;sxt r2;sxt r3;jmp r4
		-50008;12604 5002 5003 104                  /int64u
		-50010;12604 10003 10102 12600 73200 12600 12600 12600 10300 10201 104  /x<<y (ashc)
		-50011;12604 10003 10102 12600 5400 73200 12600 12600 12600 10300 10201 104  /x>>y (neg;ashc)
		-50012;99999 99999 104                      /todo x>>y unsigned
		-50013;12604 5116 5166 2 5166 4 5166 6 42600 42601 42602 42603 104  /x&y: com (sp)..; bic (sp)+,r0..
		-50014;12604 52600 52601 52602 52603 104    /x|y: bis (sp)+,r0..
		-50015;16604 2 74400 16604 4 74401 16604 6 74402 16604 10 74403 12604 61627 10 104 /x^y: xor src must be a register
		-50016;12604 42600 42601 42603 42604 104    /x&^y (and not: bic..)
		-50017;10002 10103 5000 5001 207            /u64<<32
		-50018;10200 10301 5002 5003 207            /u64>>32
		-50019;12701 13 72001 10003 5000 5001 5002 207     /u64<<59
		-50020;12700 177765 72300 10300 5001 5002 5003 207 /u64>>59
		-50021;12604 12602 12603 62706 4 20200 5603 160103 1001 20200 104 /cmp4
		-50022;12604 162600 162601 162602 162603 1010 5702 1006 5701 1004 5700 1002 5000 104 5000 5200 104 /cmp8                      /todo cmp8 (only for == and !=)
		-50023;12604 12602 12603 62706 4 20301 1001 20200 104)   /cmpu 32 bit unsigned
getbuiltin:{$[#builtins i:-10000-x;x;builtins[i]:12604 10346 10246 10146 10046 16500 0 16501 0 16502 0 16503 0 104+0 0 0 0 0,(,/+2^(&4),O x+2*!4),0];x}  /e.g. mov 20(r5),r0 .. 26(r5),r3;jmp r4  (get x2)
setbuiltin:{$[#builtins i:-20000-x;x;builtins[i]:12604 10065 0 10165 0 10265 0 10365 0 12600 12601 12602 12603 104+0,(,/+2^(&4),O x+2*!4),0 0 0 0 0];x}  /e.g. mov r0,50(r5) .. r3,56(r5);jmp r4  (set x5)


/mul32 r0..r3:[a b][c d]     ab*cd = [a*c +a*d+b*c]
/  mov  r1,r2     10102      shift high bit
/  push r0        10046
/  mul  2(sp),r0  70066 2    a*c in r0,r1   pdp11 mul stores hi,lo words in r0,r1
/  mul  2(sp),r2  70266 2    b*c in r2,r3
/  add  r3,r0     60300      (hi ac)+lo bc
/  mov  4(sp),r2  16602 4    d in r2
/  mul  (sp+),r2  70226      a*d in r2,r3
/  add  r3,r0     60300      (hi ac)+(lo bc)+lo ad
/  mov  r0,r2     10002      swap r0 and r1
/  mov  r1,r0     10100
/  mov  r2,r1     10201
/  add  #8,sp     62706 10 
/
/memcpy:
/ mov (sp)+,r4
/                /r0 src       n mov r0,r2       10002
/ mov (sp+),r1   /r1 dst       s mov (sp)+,r0    12600
/ mov 6(sp),r2   /r2 n         d mov 6(sp),r1    16601 06
/ tst r2
/ beq +3
/ movb (r0)+,(r1)+
/ dec r2
/ br -5
/ add #16,sp
/ mov +(sp),r0
/ mov +(sp),r1
/ mov +(sp),r2
/ mov +(sp),r3
/ rts pc
/
/clz  x in r0,r1
/ mov(sp)+,r4
/ mov r1,r2
/ mov r0,r3
/ clr r0
/ tst r2
/ bne +3
/ add #20,r0
/ mov r3,r2
/ bit #100000,r2
/ bne +3
/ inc r0
/ asl r2
/ br -5
/ clr r1
/ jmp r4
