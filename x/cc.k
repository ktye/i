
F:`!()                                         /function map: type to stringer


cc:{[o]opts::o
 F[`prg`mem`con`var`tab`fun`arg`sym`res`loc`ast`lod`sto]:(ccprg;ccmem;cccon;ccvar;cctab;ccfun;ccarg;ccsym;ccres;ccloc;ccast;cclod;ccsto)
 F[`add`and`bnd`ant`asn`cal`cli`cnd`dfr`div`drp`eql`cst]:(ccadd;ccand;ccbnd;ccant;ccasn;cccal;cccli;cccnd;ccdfr;ccdiv;ccdrp;cceql;cccst)
 F[`for`get`Get`gte`jmp`les`lit`lte`mod`mor`mul`neg    ]:(ccfor;ccget;ccGet;ccgte;ccjmp;ccles;cclit;cclte;ccmod;ccmor;ccmul;ccneg)
 F[`neq`nop`not`orr`bor`xor`ret`shl`shr`stm`sub`swc`typ]:(ccneq;ccnop;ccnot;ccorr;ccbor;ccxor;ccret;ccshl;ccshr;ccstm;ccsub;ccswc;cctyp)

 i:i@&{(0N~|/(T 1_c)?`swc`stm)&3=#c:&P=x}'i:&T=`cnd
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

nl:{$[#x:"\n"/:x;x,"\n";x]}
nyi:{y;em x}

cctype:`i`u`j`k`f!("int32_t";"uint32_t";"int64_t";"uint64_t";"double")

ccprg:{
 fncs: ,/y@&`fun=T c:1_&P=0
 cons:nl y@&`con=T c
 vars:nl y@&`var=T c
 runt:$[`nort':opts;"";ccrnt]
 data:$[#D;"memcpy(_M,\"",("\\x","\\x"/:(#D)^`x D),"\",",($#D),");";""]
 tabl:nl y@&`tab=T c
 mems:nl y@&`mem=T c
 init:$[|/(#mems;#data;#tabl);nl("void init(){\n";mems;data;tabl;"}\n");""]
 (runt;cons;vars;ccfmt init),,fncs}
 
ccfun:{
 $[(`main~S x)&`lib':opts; :""; 0]
 args:", "/:y@&`arg=T c:&x=P
 ret:$[#r:", "/:y@&`res=T c;r;"void"]
 locs:nl@{(cctype x)," ",y,";"}'[!g;.g:","/:'(y i)@=(S c)i:&`loc=T c]
 defr:     *y@&`dfr=T c
 body: 2_-1_y@*&`ast=T c
 ccfmt@,/(ret;" ";$S x;"(";args;"){\n";defr;locs;body)}
 

ccmem:{y;("Memory";$[`b~S x;"2";""];"(";$I x;");")}
cctab:{y;("_F[";$I x;"]=";$S x;";\n")}

/floating point literals that don't roundtrip (written by mk flt)
ccflts:(0x25fab843dda06440;0x2cde25af6a56113f;0x182d4454fb21e93f;0x485555555555c5bf;0xcd3b7f669ea0e63f;0x935555555555e53f;0xc64bac7e4f7e92be;0x4b5555555555a53f;0x5d1f29a9e5e55abe;0xf544c819a001fa3e;0x9b1a86a049faa8bd;0xe69d3f334f500340;0xec62bdfb9e516840;0xf16bd2c541bdbbbe;0xaf788e1dc571cc3f;0x182d4454fb210940;0x93bdbe166cc166bf;0x7a806b5b542830c0;0xcd9cd11ffdd8e53d;0xea49b0133f567e40;0xfe822b654715f73f;0x000000000000a03c;0x0000000000001000;0x075c143326a6913c;0x8eec28fd693650c0;0x00000000000030be;0x04fa97999999d93f;0xff799f501344d33f;0xa1487d56e31dc73e;0x03dfbf19a0012abf;0xd0f710111111813f;0x9425f7a17f00ecbf;0x730288368cc052c0;0x3bbee2d2180e7b40;0xef39fafe422e8640;0x44523edf12f1c23f;0x000000002d44643e;0x7051cc989846e83c;0x25ba052dbfb85ec0;0x51302dd5104987c0;0x0000000000003043;0x399d52a246df913f;0x00000040fb21e93f;0x555555555555c53f;0xde03cb966446c73f;0xffffffffffffef7f;0x0000e0fe422ee63f;0xd0a4be726937663e;0x9fc678d0099ac33f;0x000000000000e043;0x053f4e7b9dee213e;0x182d4454fb21f93f;0x3c60145bc4db3840;0x0ee526157bcbdb3f;0x914fc1166cc156bf;0x000000000000303e;0x83c8c96d305ff43f;0x763c7935ef39ea3d;0x599322942449d23f;0x8dedb5a0f7c6b03e;0xf8c1631adca54c40)!("165.02700983169885";"6.613756321437934e-05";"0.7853981633974483";"-0.1666666666666663";"0.7071067811865476";"0.6666666666666735";"-2.755731417929674e-07";"0.041666666666666595";"-2.5050747762857807e-08";"2.4801587288851704e-05";"-1.1358536521387682e-11";"2.414213562373095";"194.5506571482614";"-1.6533902205465252e-06";"0.22222198432149784";"3.141592653589793";"-0.0027777777777015593";"-16.157537187333652";"1.5896230157654656e-10";"485.3903996359137";"1.4426950408889634";"1.1102230246251565e-16";"2.2250738585072014e-308";"6.123233995736766e-17";"-64.85021904942025";"-3.725290298461914e-09";"0.3999999999940942";"0.3010299956639812";"2.7557313621385722e-06";"-0.0001984126982958954";"0.008333333333322118";"-0.8750608600031904";"-75.00855792314705";"432.88106049129027";"709.782712893384";"0.14798198605116586";"3.774894707930798e-08";"2.6951514290790595e-15";"-122.88666844901361";"-745.1332191019411";"4.503599627370496e+15";"0.017453292519943295";"0.7853981256484985";"0.16666666666666666";"0.1818357216161805";"1.7976931348623157e+308";"0.6931471803691238";"4.1381367970572385e-08";"0.15313837699209373";"9.223372036854776e+18";"2.087570084197473e-09";"1.5707963267948966";"24.858464901423062";"0.4342944819032518";"-0.0013888888888873056";"3.725290298461914e-09";"1.2732395447351628";"1.9082149292705877e-10";"0.2857142874366239";"1e-06";"57.29577951308232")
ccflt:{x:C x+!8;$[#r:ccflts x;r;$*`f x]}
cclit:{y;i:I x;$[`f~t:S x;ccflt i
                 `i~t;$[0N~i;"-0x80000000";$i]
                 `u~t;$[i<0;"0x",`x@|`c@,i;$i]
		 `j~t;$[i~0;,"0";(-1~h:*`i C 4+i+!4)&0>j:*`i C i+!4; $j; (h~0)&j>0; $j; hx8 i]
		 `k~t;$[i~0;,"0";hx8 i];"???"]}

ccvar:{" "/:(cctype S 1+x;  $S x;,"=";(*y),";")}
cccon:{" "/:("const";cctype S 1+x;$S x;,"=";*y)}

cccst:{("(";cctype S x;")";,"(";y 1;,")")}

cctyp:{y;cctype S x} 
ccres:{y;cctype S x}
ccloc:{y}
ccsym:{y;$S x}
ccarg:{(cctype S x;" ";*y)}

ccstm:{ "{\n",(,/y,\:";\n"),"}\n"}
/ccstm:{"{\n",(";\n"/:y),"}\n"}
ccast:ccstm
ccget:{y;$S x}
ccGet:ccget
ccret:{"return ",(", "/:y),";\n"}

ccsiz:`b`i`j`f!`I8`I32`I64`F64
cclod:{($ccsiz S x;"(";y 0;")")}
ccsto:{("Set";$ccsiz S x;"(";y 0;",";y 1;")")}

cccal:{($S x;,"("; ","/:y ;,")")}
ccdrp:{y}

cccli:{n:I x;("((";cctype S x;"(*)("; ","/:n#(1+n)_y; "))Func[";*y;"])(";  ","/:n#1_y; ")")}
ccasn:{r:(s:$S x),"=",(*|y),";\n";
 $[((n:3+2*#s)#r)~s,"=(",s,"+";$["+=1"~-3#r:s,"+=",(n_-3_r);s,"++";r],";\n";r]}
cccnd:{$[S x; (ccnd@*y;"?(";y 1;"):(";y 2;")") ;"if",(ccnd y 0),(y 1),$[3~#y;"else ",(y 2);""]]}
ccnop:{y;""}
ccswc:{("switch(";*y;"){\n"),(1_(@[("case ",/:$-1+!#y); (-1+#y)*I x; "default"] ,\:":\n"),'y,\:"break;\n"),"}" }
 
ccfor:{$[I x;ccwhl[x;y];("for(;";ccnd y 0;";";$[";\n"~-2#r:y 1;-2_r;r];")";y 2)]}
ccwhl:{((,"do"),,-12_y 2),("}while";ccnd y 0;";\n")}
ccjmp:{y;$[I x;"break";"continue"]}
ccdfr:{("defer func(){if recover()!=nil{";*y;"}}()")}

cc1:{x,*z}
cc2:{z;("(";z 0;x;z 1;")")}
ccneg:cc1"-"; ccnot:cc1"!"
ccadd:cc2"+"; ccsub:cc2"-"; ccmul:cc2"*"; ccdiv:cc2"/"; ccmod:cc2"%"
cceql:cc2"==";ccgte:cc2">=";ccles:cc2"<"; cclte:cc2"<=";ccmor:cc2">"; ccneq:cc2"!="
ccand:cc2"&"; ccbnd:cc2"&&";ccant:cc2"&^";ccorr:cc2"|"; ccbor:cc2"||";ccxor:cc2"^"; ccshl:cc2"<<";ccshr:cc2">>"

ccnd:{$["!=0)"~t:-4#x;(-4_x),")";"==0)"~t;"(!",(1_-4_x),")";x]}

ccfmt:{
 x:@[x;1+find["}\nelse";x];" "]
 x:@[x;1+find["}\nwhil";x];" "]
 l:"\n"\:x;
 l@:&~(""~/:l)|(,";")~/:l;s:_32+&'(-"}"=*'l)+-1_0,+\+/'("{"=l)-"}"=l
 ,/(s,'l),\:"\n"}

/ run time
ccrnt:"import ( \"encoding/binary\"; \"math\"; \"math/bits\"; \"os\" )
var Bytes, Bytes2 []byte
var Func []interface{}
func Memory(blocks int) { if Bytes == nil { Bytes = make([]byte, 64*1024*blocks) } }
func Memory2(blocks int) { if Bytes2 == nil { Bytes2 = make([]byte, 64*1024*blocks) } }
func Memorysize() int32  { return int32(len(Bytes) >> 16) }
func Memorysize2() int32 { return int32(len(Bytes2) >> 16) }
func Memorygrow(blocks int32) (previous int32) {
	previous = int32(len(Bytes) >> 16)
	delta := 64 * 1024 * int(blocks)
	if len(Bytes)+delta > (1<<32 - 1) { return -1 }
	Bytes = append(Bytes, make([]byte, delta)...)
	return previous
}
func Memorygrow2(blocks int32) int32 {
	Bytes, Bytes2 = Bytes2, Bytes
	r := Memorygrow(blocks)
	Bytes, Bytes2 = Bytes2, Bytes
	return r
}
func Memorycopy(dst, src, n int32) { copy(Bytes[dst:], Bytes[src:src+n]) }
func Memoryfill(dst, val, n int32) { b := byte(val); for i := int32(0); i < n; i++ { Bytes[dst+i] = b } }
func Memorycopy2(dst, src, n int32) { copy(Bytes2[dst:], Bytes[src:src+n]) }
func Memorycopy3(dst, src, n int32) { copy(Bytes[dst:], Bytes2[src:src+n]) }
func Functions(i int, f interface{}) { for i >= len(Func) { Func = append(Func, nil) }; Func[i] = f }
func Data(o int, x string) { copy(Bytes, x) }
func I8(addr int32) int32           { return int32(int8(Bytes[addr])) }
func SetI8(addr int32, value int32) { Bytes[addr] = byte(value) }
func I32(addr int32) int32  { return int32(binary.LittleEndian.Uint32(Bytes[addr:])) }
func I64(addr int32) int64  { return int64(binary.LittleEndian.Uint64(Bytes[addr:])) }
func F64(addr int32) float64 { return math.Float64frombits(binary.LittleEndian.Uint64(Bytes[addr:])) }
func SetI32(addr int32, value int32) { binary.LittleEndian.PutUint32(Bytes[addr:], uint32(value)) }
func SetI64(addr int32, value int64) { binary.LittleEndian.PutUint64(Bytes[addr:], uint64(value)) }
func SetF64(addr int32, value float64) { binary.LittleEndian.PutUint64(Bytes[addr:], math.Float64bits(value)) }
func I32clz(x uint32) int32    { return int32(bits.LeadingZeros32(x)) }
func I64popcnt(x uint64) int64 { return int64(bits.OnesCount64(x)) }
func I32B(b bool) int32 { if b { return 1; }; return 0}
func Exit(x int32) { os.Exit(int(x)) }
func Args() int32 { return int32(len(os.Args)) }
func Arg(i, r int32) int32 {
	if i >= int32(len(os.Args)) { return 0 }
	if r == 0 { return int32(len(os.Args[i])) }
	copy(Bytes[r:], []byte(os.Args[i]))
	return 0
}
var filebuf []byte
func Read(file, nfile, dst int32) int32 {
	if dst != 0 { copy(Bytes[dst:], filebuf); return 0 }
	b, e := os.ReadFile(string(Bytes[file : file+nfile]))
	if e != nil { return -1 }
	filebuf = b
	return int32(len(filebuf))
}
func Write(file, nfile, src, n int32) int32 {
	b := Bytes[src : src+n]
	if nfile == 0 { os.Stdout.Write(b); return 0 }
	name := string(Bytes[file : file+nfile])
	e := os.WriteFile(name, b, 0644)
	if e != nil { return -1 }
	return 0
}
func ReadIn(dst, n int32) int32 {
	nr, e := os.Stdin.Read(Bytes[dst : dst+n])
	if e != nil { return 0 }
	if nr > 0 && Bytes[dst+int32(nr)-1] == 10 { nr -= 1 }
	return int32(nr)
}
func Native(x, y int64) int64 { return x }
func F64reinterpret_i64(x uint64) float64 { return math.Float64frombits(x) }
func I64reinterpret_f64(x float64) uint64 { return math.Float64bits(x) }
func F64abs(x float64) float64         { return math.Abs(x) }
func F64sqrt(x float64) float64        { return math.Sqrt(x) }
func F64ceil(x float64) float64        { return math.Ceil(x) }
func F64floor(x float64) float64       { return math.Floor(x) }
func F64min(x, y float64) float64      { return math.Min(x, y) }
func F64max(x, y float64) float64      { return math.Max(x, y) }
func F64copysign(x, y float64) float64 { return math.Copysign(x, y) }
"
