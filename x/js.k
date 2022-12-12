
F:`!()

js:{[o]opts::o
 F[`prg`mem`con`var`tab`fun`arg`sym`res`loc`ast`lod`sto]:(jsprg;jsmem;jscon;jsvar;jstab;jsfun;jsarg;jssym;jsres;jsloc;jsast;jslod;jssto)
 F[`add`and`bnd`ant`asn`cal`cli`cnd`dfr`div`drp`eql`cst]:(jsadd;jsand;jsbnd;jsant;jsasn;jscal;jscli;jscnd;jsdfr;jsdiv;jsdrp;jseql;jscst)
 F[`for`get`Get`gte`jmp`les`lit`lte`mod`mor`mul`neg    ]:(jsfor;jsget;jsGet;jsgte;jsjmp;jsles;jslit;jslte;jsmod;jsmor;jsmul;jsneg)
 F[`neq`nop`not`orr`bor`xor`ret`shl`shr`stm`sub`swc`typ]:(jsneq;jsnop;jsnot;jsorr;jsbor;jsxor;jsret;jsshl;jsshr;jsstm;jssub;jsswc;jstyp)

 i:i@&{(0N~|/(T 1_c)?`swc`stm)&3=#c:&P=x}'i:&T=`cnd  /keep ternary
 rwcnd'|j@&~(j:&(~`=S)&(T=`cnd)|T=`swc)':i     /rewrite cnd/swc expressions 
 S[0]:$[0N~i:o?`lib;S 0;o 1+i]                 /emit`lib`name (instead of main + remove func main)
 n:(&#P){x[y]+:1}/P;n[0]-:1;                   /number of children per node
 s:{c:(0,(#x)-z)^x                             /cut args from stack
   (c 0),,,/(F T y)[y;|c 1]}/[();|!#P;|n]      /fold over the tree in reverse order
*s}


hx8:{[i]"0x",`x@|C i+!8}
em:{"(",x,")"}

INS:{[i;x]T::ins[i;T;*x];P::ins[i;@[P;&P>i-1;+;1];x 1];P[1+i]:i;I::ins[i;I;x 2];S::ins[i;S;x 3]}
ins:{[i;x;y](i#x),y,i_x}
del:{T::T n:(!#P)_x;P::@[P n;&(P n)>x-1;-;1];I::I n;S::S n}
rwcnd:{c:|1_&P=x;S[x]:`;P[x]:P x-1;c INS\:(T x;1+x;I x; S x:x-1);del x;restm'|i@&(T P i:&`stm=T)':`ret`asn}
restm:{c:*|&P=x;INS[c;(T x-1;x;I x-1;S x-1)];P[x]:P[x-1];del[x-1]} /fix ret stm

jsren:{S[&(S=x)&T':`fun`cal]:s:`$"_",$x;}

nl:{$[#x:"\n"/:x;x,"\n";x]}
nyi:{y;em x}

jsty:{(_($x)-32),em y}
jstype:`i`u`j`k`f!"IUJKF"

jsprg:{
 f:(S i)!i:&T=`fun
 fpre:nl@jssig'f@?S i@&i<f S i:&T=`cal
 fncs: ,/y@&`fun=T c:1_&P=0
 cons:nl y@&`con=T c
 vars:nl y@&`var=T c
 runt:$[`nort':opts;"";jsrnt]
 main:$[|/`nort`lib':opts;"";jsmain]
 data:$[#D;"memcpy(_M,\"",("\\x","\\x"/:(#D)^`x D),"\",",($#D),");";""]
 tabl:nl y@&`tab=T c
 tabl:$[#tabl;"_F=malloc(",($1+|/I@&`tab=T),"*sizeof(void*));\n";""],tabl
 mems:nl y@&`mem=T c
 init:$[|/(#mems;#data;#tabl);nl("void init(){\n";mems;data;tabl;"}\n");""]
 (runt;fpre;cons;vars),(,fncs),(jsfmt init),main}
 
jsfun:{
 $[(`main~S x)&`lib':opts; :""; 0]
 args:","/:y@&`arg=T c:&x=P
 locs:($[#r:","/:$S 1+c@&`loc=T c;"let ",r,"\n";""])
 defr:     *y@&`dfr=T c
 setj:$[#defr;"if(!setjmp(_jb_)){";""]
 body: 2_-1_y@*&`ast=T c
 jsfmt@,/("let ";$S x;"=function(";args;"){\n";setj;locs;body;defr)}

jssig:{
 args:","/:jstype S c@&`arg=T c:&x=P
 ret:$[#r:jstype S c@&`res=T c;*r;"void"]
 ret," ",($S x),"(",args,");"}

jsmem:{y;("Memory";$[`b~S x;"2";""];"(";$I x;");")}
jstab:{y;("_F[";$I x;"]=";$S x;";\n")}

/floating point literals that don't roundtrip (written by mk flt)
jsflts:(0x25fab843dda06440;0x2cde25af6a56113f;0x182d4454fb21e93f;0x485555555555c5bf;0xcd3b7f669ea0e63f;0x935555555555e53f;0xc64bac7e4f7e92be;0x4b5555555555a53f;0x5d1f29a9e5e55abe;0xf544c819a001fa3e;0x9b1a86a049faa8bd;0xe69d3f334f500340;0xec62bdfb9e516840;0xf16bd2c541bdbbbe;0xaf788e1dc571cc3f;0x182d4454fb210940;0x93bdbe166cc166bf;0x7a806b5b542830c0;0xcd9cd11ffdd8e53d;0xea49b0133f567e40;0xfe822b654715f73f;0x000000000000a03c;0x0000000000001000;0x075c143326a6913c;0x8eec28fd693650c0;0x00000000000030be;0x04fa97999999d93f;0xff799f501344d33f;0xa1487d56e31dc73e;0x03dfbf19a0012abf;0xd0f710111111813f;0x9425f7a17f00ecbf;0x730288368cc052c0;0x3bbee2d2180e7b40;0xef39fafe422e8640;0x44523edf12f1c23f;0x000000002d44643e;0x7051cc989846e83c;0x25ba052dbfb85ec0;0x51302dd5104987c0;0x0000000000003043;0x399d52a246df913f;0x00000040fb21e93f;0x555555555555c53f;0xde03cb966446c73f;0xffffffffffffef7f;0x0000e0fe422ee63f;0xd0a4be726937663e;0x9fc678d0099ac33f;0x000000000000e043;0x053f4e7b9dee213e;0x182d4454fb21f93f;0x3c60145bc4db3840;0x0ee526157bcbdb3f;0x914fc1166cc156bf;0x000000000000303e;0x83c8c96d305ff43f;0x763c7935ef39ea3d;0x599322942449d23f;0x8dedb5a0f7c6b03e;0xf8c1631adca54c40)!("165.02700983169885";"6.613756321437934e-05";"0.7853981633974483";"-0.1666666666666663";"0.7071067811865476";"0.6666666666666735";"-2.755731417929674e-07";"0.041666666666666595";"-2.5050747762857807e-08";"2.4801587288851704e-05";"-1.1358536521387682e-11";"2.414213562373095";"194.5506571482614";"-1.6533902205465252e-06";"0.22222198432149784";"3.141592653589793";"-0.0027777777777015593";"-16.157537187333652";"1.5896230157654656e-10";"485.3903996359137";"1.4426950408889634";"1.1102230246251565e-16";"2.2250738585072014e-308";"6.123233995736766e-17";"-64.85021904942025";"-3.725290298461914e-09";"0.3999999999940942";"0.3010299956639812";"2.7557313621385722e-06";"-0.0001984126982958954";"0.008333333333322118";"-0.8750608600031904";"-75.00855792314705";"432.88106049129027";"709.782712893384";"0.14798198605116586";"3.774894707930798e-08";"2.6951514290790595e-15";"-122.88666844901361";"-745.1332191019411";"4.503599627370496e+15";"0.017453292519943295";"0.7853981256484985";"0.16666666666666666";"0.1818357216161805";"1.7976931348623157e+308";"0.6931471803691238";"4.1381367970572385e-08";"0.15313837699209373";"9.223372036854776e+18";"2.087570084197473e-09";"1.5707963267948966";"24.858464901423062";"0.4342944819032518";"-0.0013888888888873056";"3.725290298461914e-09";"1.2732395447351628";"1.9082149292705877e-10";"0.2857142874366239";"1e-06";"57.29577951308232")
jsflt:{x:C x+!8;$[#r:jsflts x;r;$*`f x]}
jslit:{hi:{*`i C 4+x+!4};lo:{*`i C x+!4}
 y;i:I x;$[`f~t:S x;jsflt i
           `i~t;"(",$[0N~i;"-2147483648";$i],"|0)"
           `u~t;$[i<0;"0?",`x@|`c@,i;$i]
           `j~t;$[i~0;,"J(0n)";(-1~h:hi i)&0>l:lo i;$l;(h~0)&l>0;$l;"J(BigInt(\"",(hx8 i),"\"))"]
           `k~t;$[i~0;,"K(0n)";(0~h:hi i)&0<l:lo i;"K(",($l),"n)";"K(BigInt(\"",(hx8 i),"\"))"];"???"]}
 

jsvar:{" "/:(jstype S 1+x;  $S x;,"=";(*y),";")}
jscon:{"const ",(jstype S 1+x)," ",($S x),"=",(*y),";"}

jscst:{(jstype S x),(*y),em y 1}

jstyp:{y;jstype S x} 
jsres:{y;jstype S x}
jsloc:{y}
jssym:{y;$S x}
jsarg:{*y}

jsstm:{ "{\n",(,/y,\:";\n"),"}\n"}
/jsstm:{"{\n",(";\n"/:y),"}\n"}
jsast:jsstm
jsget:{y;$S x}
jsGet:jsget
jsret:{"return ",(", "/:y),";\n"}

jssiz:`b`i`j`f!`I8`I32`I64`F64
jslod:{($jssiz S x;"(";y 0;")")}
jssto:{("Set";$jssiz S x;"(";y 0;",";y 1;")")}

jscal:{($S x;,"("; ","/:y ;,")")}
jsdrp:{y}

jscli:{n:I x;jsty[S x;(","/:n#(1+n)_y),"F[",(*y),"])(",(","/:n#1_y),")"]}
jsasn:{r:(s:$S x),"=",(*|y),";\n";
 $[((n:3+2*#s)#r)~s,"=(",s,"+";$["+=1"~-3#r:s,"+=",(n_-3_r);s,"++";r],";\n";r]}
jscnd:{$[S x; (jsnd@*y;"?(";y 1;"):(";y 2;")") ;"if(",(jsnd y 0),")",(jsemb jssem y 1),$[3~#y;"else ",(jsemb jssem y 2);""]]}
jsnop:{y;""}
jsswc:{("switch(";*y;"){\n"),(1_(@[("case ",/:$-1+!#y); (-1+#y)*I x; "default"] ,\:":"),'(jssem'y),\:"break;\n"),"}" }
 
jsfor:{$[I x;jswhl[x;y];("for(;";jsnd y 0;";";$[";\n"~-2#r:y 1;-2_r;r];")";y 2)]}
jswhl:{((,"do"),(,-12_y 2),(,y 1)),("}while";jsnd y 0;";\n")}
jsjmp:{y;$[I x;"break";"continue"]}
jsdfr:{("else{";*y;";}\n}")}

js1:{r:x,*z;$[`i~S y;em r,"|0";r]}
js2:{z;jsty[S y;(z 0),x,z 1]}
jsneg:js1"-"; jsnot:js1"!"
jsadd:js2"+"; jssub:js2"-"; jsmul:js2"*"; jsdiv:js2"/"; jsmod:js2"%"
jseql:js2"==";jsgte:js2">=";jsles:js2"<"; jslte:js2"<=";jsmor:js2">"; jsneq:js2"!="
jsand:js2"&"; jsbnd:js2"&&";jsant:js2"&~";jsorr:js2"|"; jsbor:js2"||";jsxor:js2"^"; jsshl:js2"<<";jsshr:js2">>"

jsnd:{$["!=0)"~t:-4#x;(-4_x),")";"==0)"~t;"(!",(1_-4_x),")";x]}
jssem:{$[(*|x)':";}\n";x;x,";"]}
jsemb:{$["{"=*x;x;"{\n",x,"\n}\n"]}

jsfmt:{
 x:@[x;1+find["}\nelse";x];" "]
 x:@[x;1+find["}\nwhil";x];" "]
 x:@[x;1+find["};\n"   ;x];" "]
 l:"\n"\:x;
 l@:&~l~\:"r=r;"
 l@:&~(""~/:l)|(,";")~/:l;s:_32+&'(-"}"=*'l)+-1_0,+\+/'("{"=l)-"}"=l
 ,/(s,'l),\:"\n"}

/ run time
jsrnt:"
const I=x=>x|0,J=x=>BigInt.asIntN(64,x),K=x=>BigInt.asUintN(64,x),F=x=>+x
"
