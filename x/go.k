
/C  bytes(e.g.64bit literals)
/D  data (initial memory)
F:`!()                                         /function map: type to stringer
t:+`n`t`p`i`s!(!#T;T;P;I;S)


go:{[o]opts::o
 F[`prg`mem`con`var`tab`fun`arg`sym`res`loc`ast`lod`sto]:(goprg;gomem;gocon;govar;gotab;gofun;goarg;gosym;gores;goloc;goast;golod;gosto)
 F[`add`and`bnd`ant`asn`cal`cli`cnd`dfr`div`drp`eql`cst]:(goadd;goand;gobnd;goant;goasn;gocal;gocli;gocnd;godfr;godiv;godrp;goeql;gocst)
 F[`for`get`Get`gte`jmp`les`lit`lte`mod`mor`mul`neg    ]:(gofor;goget;goGet;gogte;gojmp;goles;golit;golte;gomod;gomor;gomul;goneg)
 F[`neq`nop`not`orr`bor`xor`ret`shl`shr`stm`sub`swc`typ]:(goneq;gonop;gonot;goorr;gobor;goxor;goret;goshl;goshr;gostm;gosub;goswc;gotyp)

 S[0]:$[0N~i:o?`lib;S 0;o 1+i]                 /emit`lib`name (instead of main + remove func main)
 n:(&#P){x[y]+:1}/P;n[0]-:1;                   /number of children per node
 s:{c:(0,(#x)-z)^x                             /cut args from stack
   (c 0),,,/(F T y)[y;|c 1]}/[();|!#P;|n]      /fold over the tree in reverse order
*s}


/set pkg name and delete main
pkg:{S[0]::x; rm@*&(`fun=T)&`main=S}

/cld:{&x=P}
hx8:{[i]"0x",`x@|C i+!8}
/dp:{+/0<(x@)\:!#x} /depth from parent
em:{"(",x,")"}

nl:{$[#x:"\n"/:x;x,"\n";x]}
nyi:{y;em x}

gotype:`i`u`j`k`f`C`I`F!("int32";"uint32";"int64";"uint64";"float64";"[16]int8";"[4]int32";"[2]float64")

goprg:{
 fncs:nl y@&`fun=T c:1_&P=0
 cons:nl y@&`con=T c
 vars:nl y@&`var=T c
 runt:$[`nort':opts;"";gornt]
 data:$[#D;"Data(0,\"",("\\x","\\x"/:(#D)^`x D),"\")";""]
 tabl:nl y@&`tab=T c
 mems:nl y@&`mem=T c
 init:$[|/(#mems;#data;#tabl);nl("func init(){";mems;data;tabl;"}\n");""]
 ("package ";$S 0;"\n";runt;cons;vars;init),,fncs}
 
gofun:{
 $[(`main~S x)&`lib':opts; :""; 0]
 args:", "/:y@&`arg=T c:&x=P
 rets:", "/:y@&`res=T c
 locs:"\n"/:"var ",/:(.g),'" ",/:gotype@!g:","/:'(y i)@=(S c)i:&`loc=T c
 defr:     *y@&`dfr=T c
 body:      y@*&`ast=T c
 ("func ";$S x;"(";args;")(";rets;"){";defr;locs,"\n";body;"}")}
 

gomem:{y;("Memory";$[`b~S x;"2";""];"( ";$I x;" )")}
gotab:{y;("Functions(";$I x;", ";$S x;") ")}

goflts:(,0x182d4454fb210940)!,"3.141592653589793" /predefined constants that don't roundtrip
goflt:{x:C x+!8;$[#r:goflts x;r;$*`f x]}
golit:{y;i:I x;$[`f~t:S x;goflt i
                 `i~t;$[0N~i;"-0x80000000";$i]
                 `u~t;$[i<0;"0x",`x@|`c@,i;$i]
		 `j~t;$[i~0;,"0";(-1~h:*`i C 4+i+!4)&0>j:*`i C i+!4; $j; (h~0)&j>0; $j; hx8 i]
		 `k~t;$[i~0;,"0";hx8 i];"???"]}

govar:{" "/:("var";  $S x;gotype S 1+x;,"=";*y)}
gocon:{" "/:("const";$S x;gotype S 1+x;,"=";*y)}

gocst:{(gotype S x;,"(";y 1;,")")}

gotyp:{y;gotype S x} 
gores:{y;gotype S x}
goloc:{y}
gosym:{y;$S x}
goarg:{((y 0);" ";gotype S x)}

gostm:{"\n"/:y}
goast:gostm
goget:{y;$S x}
goGet:goget
goret:{"return ",(", "/:y)}

gosiz:`b`i`j`f!`I8`I32`I64`F64
golod:{($gosiz S x;"( ";y 0;") ")}
gosto:{("Set";$gosiz S x;"( ";y 0;", ";y 1;") ")}

gocal:{($S x;,"("; ","/:y ;,")")}
godrp:{y}

gocli:{n:I x;("Func[";*y;"].(func(";  ","/:n#(1+n)_y; ")("; ","/:(1+2*n)_y ;"))(";  ","/:n#1_y; ") ")}
goasn:{(", "/:(I x)#y)," = ",*|y}
gocnd:{"if ",(y 0)," {",(y 1),"} ",$[3~#y;" else {",(y 2)," }";""]}
gonop:{y;""}
goswc:{("switch ";*y;" {"),(1_(@[("case ",/:$-1+!#y); (-1+#y)*I x; "default"] ,\:":\n"),'y,\:"\n")," }" }
 
gofor:{($[#l:$S x;l,":";""];"for ;";y 0;"; ";y 1;" {";y 2;"} ")}
gojmp:{y;$[I x;"break ";"continue "],$S x}
godfr:{("defer func(){if recover()!=nil{";*y;"}}()")}

go1:{x,*z}
go2:{z;("(";z 0;x;z 1;")")}
goneg:go1"-"; gonot:go1"!"
goadd:go2"+"; gosub:go2"-"; gomul:go2"*"; godiv:go2"/"; gomod:go2"%"
goeql:go2"==";gogte:go2">=";goles:go2"< ";golte:go2"<=";gomor:go2">"; goneq:go2"!="
goand:go2"&"; gobnd:go2"&&";goant:go2"&^";goorr:go2"|"; gobor:go2"||";goxor:go2"^"; goshl:go2"<<";goshr:go2">>"



/ run time
gornt:"import ( \"encoding/binary\"; \"math\"; \"math/bits\"; \"os\" )
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
func I8x16load(addr int32) (r [16]int8) { for i := range r { r[i] = int8(Bytes[addr+int32(i)]); }; return r }
func I32x4load(addr int32) [4]int32 { return [4]int32{I32(addr), I32(4 + addr), I32(8 + addr), I32(12 + addr)} }
func F64x2load(addr int32) [2]float64 { return [2]float64{F64(addr), F64(8 + addr)} }
func I8x16store(addr int32, x [16]int8) { for i := range x { Bytes[addr+int32(i)] = byte(x[i]) } }
func I32x4store(addr int32, x [4]int32) { for i := range x { SetI32(addr, x[i]); addr += 4 } }
func F64x2store(addr int32, x [2]float64) { SetF64(addr, x[0]); SetF64(8+addr, x[1]) }
func I8x16splat(x int32) (r [16]int8) { for i := range r { r[i] = int8(x); }; return r }
func I32x4splat(x int32) [4]int32 { return [4]int32{x, x, x, x} }
func F64x2splat(x float64) [2]float64 { return [2]float64{x, x} }
func I8x16Extract_lane_s0(x [16]int8) int32 { return int32(x[0]) }
func F64x2Replace_lane1(x [2]float64, f float64) [2]float64 { x[1] = f; return x }
func I8x16All_true(x [16]int8) int32 { for i := range x { if (x[i] != 0) == false { return 0; }; }; return 1 }
func I8x16Any_true(x [16]int8) int32 { for i := range x { if x[i] != 0 { return 1; }; }; return 0 }
func I32x4Any_true(x [4]int32) int32 { for i := range x { if x[i] != 0 { return 1; }; }; return 0 }
func I8x16Eq(x, y [16]int8) (r [16]int8) { for i := range r { if x[i] == y[i] { r[i] = -1; }; }; return r; }
func I32x4Eq(x, y [4]int32) (r [4]int32) { for i := range r { if x[i] == y[i] { r[i] = -1; }; }; return r }
func I8x16And(x, y [16]int8) (r [16]int8) { for i := range r { r[i] = x[i] & y[i] }; return r }
func I8x16Not(x [16]int8) (r [16]int8) { for i := range r { r[i] = ^x[i]; }; return r }
func I8x16Neg(x [16]int8) (r [16]int8) { for i := range r { r[i] = -x[i] }; return r }
func I32x4Neg(x [4]int32) (r [4]int32) { for i := range r { r[i] = -x[i] }; return r }
func F64x2Neg(x [2]float64) [2]float64 { return [2]float64{-x[0], -x[1]} }
func absi8(x int8) int8 { return int8(absi32(int32(x))) }
func absi32(x int32) int32 { if x < 0 { return -x }; return x }
func I8x16Abs(x [16]int8) (r [16]int8) { for i := range r { r[i] = absi8(x[i]) }; return r }
func I32x4Abs(x [4]int32) (r [4]int32) { for i := range r { r[i] = absi32(x[i]) }; return r }
func F64x2Abs(x [2]float64) [2]float64 { return [2]float64{math.Abs(x[0]), math.Abs(x[1])} }
func F64x2Sqrt(x [2]float64) [2]float64 { return [2]float64{math.Sqrt(x[0]), math.Sqrt(x[1])} }
func I8x16Add(x, y [16]int8) (r [16]int8) { for i := range r { r[i] = x[i] + y[i] }; return r }
func I32x4Add(x, y [4]int32) (r [4]int32) { for i := range r { r[i] = x[i] + y[i] }; return r }
func F64x2Add(x, y [2]float64) [2]float64 { return [2]float64{x[0]+y[0], x[1]+y[1]} }
func I8x16Sub(x, y [16]int8) (r [16]int8) { for i := range r { r[i] = x[i] - y[i] }; return r }
func I32x4Sub(x, y [4]int32) (r [4]int32) { for i := range r { r[i] = x[i] - y[i] }; return r }
func F64x2Sub(x, y [2]float64) [2]float64 { return [2]float64{x[0]-y[0], x[1]-y[1]} }
func I32x4Mul(x, y [4]int32) (r [4]int32) { for i := range r { r[i] = x[i] * y[i] }; return r }
func F64x2Mul(x, y [2]float64) [2]float64 { return [2]float64{x[0]*y[0], x[1]*y[1]} }
func F64x2Div(x, y [2]float64) [2]float64 { return [2]float64{x[0]/y[0], x[1]/y[1]} }
func I32x4Shr_s(x [4]int32, y int32) (r [4]int32) { for i := range r { r[i] = x[i] >> y }; return r }
func mini8(x, y int8)   int8  { return int8(mini32(int32(x), int32(y))) }
func mini32(x, y int32) int32 { if x < y { return x } else { return y } }
func maxi8(x, y int8)   int8  { return int8(maxi32(int32(x), int32(y))) }
func maxi32(x, y int32) int32 { if x > y { return x } else { return y } }
func pmin(x, y float64) float64 { if y < x { return y } else { return x } }
func pmax(x, y float64) float64 { if x < y { return y } else { return x } }
func I8x16Min_s(x, y [16]int8) (r [16]int8) { for i := range r { r[i] = mini8(x[i], y[i]) }; return r }
func I32x4Min_s(x, y [4]int32) (r [4]int32) { for i := range r { r[i] = mini32(x[i], y[i]) }; return r }
func F64x2Pmin(x, y [2]float64) [2]float64 { return [2]float64{ pmin(x[0], y[0]), pmin(x[1], y[1]) } }
func F64x2Pmax(x, y [2]float64) [2]float64 { return [2]float64{ pmax(x[0], y[0]), pmax(x[1], y[1]) } }	
func I8x16Max_s(x, y [16]int8) (r [16]int8) { for i := range r { r[i] = maxi8(x[i], y[i]) }; return r }
func I32x4Max_s(x, y [4]int32) (r [4]int32) { for i := range r { r[i] = maxi32(x[i], y[i]) }; return r }
func I8x16Lt_s(x, y [16]int8) (r [16]int8) { for i := range r { if x[i] < y[i] { r[i] = -1 } }; return r }
func I8x16Gt_s(x, y [16]int8) (r [16]int8) { for i := range r { if x[i] > y[i] { r[i] = -1 } }; return r }

"
