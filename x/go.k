
/C  bytes(e.g.64bit literals)
/D  data (initial memory)
F:`!()                                         /function map: type to stringer
t:+`n`t`p`i`s!(!#T;T;P;I;S)

emit:{[o]opts::o
 n:(&#P){x[y]+:1}/P;n[0]-:1;                   /number of children per node
 s:{c:(0,(#x)-z)^x                             /cut args from stack
   (c 0),,,/(F T y)[y;|c 1]}/[();|!#P;|n]      /fold over the tree in reverse order
*s}

/cld:{&x=P}
hx8:{[i]"0x",`x@|C i+!8}
/dp:{+/0<(x@)\:!#x} /depth from parent
em:{"(",x,")"}
nl:{$[#x:"\n"/:x;x,"\n";x]}
nyi:{y;em x}


gotype:`i`u`j`k`f`C`I`F!("int32";"uint32";"int64";"uint64";"float64";"[8]int8";"[32]int32";"[2]float64")

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
 args:", "/:y@&`arg=T c:&x=P
 rets:", "/:y@&`res=T c
 locs:"\n"/:"var ",/:(.g),'" ",/:gotype@!g:","/:'(y i)@=(S c)i:&`loc=T c
 defr:     *y@&`dfr=T c
 body:      y@*&`ast=T c
 ("func ";$S x;"(";args;")(";rets;"){";defr;locs,"\n";body;"}")}
 

gomem:{y;("Memory";$[`b~S x;"2";""];"( ";$I x;" )")}
gotab:{y;("Functions(";$I x;", ";$S x;") ")}
goflt:{(*&" "=c)#c:C 8+x+!24}
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
goasn:{n:I x;(","/:n#y; ($S x),"= "; ","/:(2*n)_y)}
gocnd:{"if ",(y 0)," {",(y 1),"} ",$[3~#y;" else {",(y 2)," }";""]}
gonop:{y;""}
goswc:{("switch ";*y;" {"),(1_(@[("case ",/:$-1+!#y); (-1+#y)*I x; "default"] ,\:":\n"),'y,\:"\n")," }" }
 
gofor:{($[#l:$S x;l,":";""];"for ;";y 0;"; ";y 1;" {";y 2;"} ")}
gojmp:{y;$[I x;"break ";"continue "],$S x}
godfr:{("defer ";*y)}

go1:{x,*z}
go2:{z;("(";z 0;x;z 1;")")}
goneg:go1"-"; gonot:go1"!"
goadd:go2"+"; gosub:go2"-"; gomul:go2"*"; godiv:go2"/"; gomod:go2"%"
goeql:go2"==";gogte:go2">=";goles:go2"< ";golte:go2"<=";gomor:go2">"; goneq:go2"!="
goand:go2"&"; gobnd:go2"&&";goant:go2"&^";goorr:go2"|"; gobor:go2"||";goxor:go2"^"; goshl:go2"<<";goshr:go2">>"

go:{
 F[`prg`mem`con`var`tab`fun`arg`sym`res`loc`ast`lod`sto]:(goprg;gomem;gocon;govar;gotab;gofun;goarg;gosym;gores;goloc;goast;golod;gosto)
 F[`add`and`bnd`ant`asn`cal`cli`cnd`dfr`div`drp`eql`cst]:(goadd;goand;gobnd;goant;goasn;gocal;gocli;gocnd;godfr;godiv;godrp;goeql;gocst)
 F[`for`get`Get`gte`jmp`les`lit`lte`mod`mor`mul`neg    ]:(gofor;goget;goGet;gogte;gojmp;goles;golit;golte;gomod;gomor;gomul;goneg)
 F[`neq`nop`not`orr`bor`xor`ret`shl`shr`stm`sub`swc`typ]:(goneq;gonop;gonot;goorr;gobor;goxor;goret;goshl;goshr;gostm;gosub;goswc;gotyp)
 x}


gornt:"import ( \"encoding/binary\"; \"fmt\"; \"math\"; \"math/bits\"; \"os\" )
var Bytes, Bytes2 []byte
var Func []interface{}
func Memory(blocks int) {
	if Bytes == nil {
		Bytes = make([]byte, 64*1024*blocks)
	}
}
func Memory2(blocks int) {
	if Bytes2 == nil {
		Bytes2 = make([]byte, 64*1024*blocks)
	}
}
func Memorysize() int32  { return int32(len(Bytes) >> 16) }
func Memorysize2() int32 { return int32(len(Bytes2) >> 16) }
func Memorygrow(blocks int32) (previous int32) {
	previous = int32(len(Bytes) >> 16)
	delta := 64 * 1024 * int(blocks)
	if len(Bytes)+delta > (1<<32 - 1) {
		return -1
	}
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
func Memoryfill(dst, val, n int32) {
	b := byte(val)
	for i := int32(0); i < n; i++ {
		Bytes[dst+i] = b
	}
}
func Memorycopy2(dst, src, n int32) { copy(Bytes2[dst:], Bytes[src:src+n]) }
func Memorycopy3(dst, src, n int32) { copy(Bytes[dst:], Bytes2[src:src+n]) }
func Functions(i int, f interface{}) {
	for i >= len(Func) {
		Func = append(Func, nil)
	}
	Func[i] = f
}
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
	if i >= int32(len(os.Args)) {
		return 0
	}
	if r == 0 {
		return int32(len(os.Args[i]))
	}
	copy(Bytes[r:], []byte(os.Args[i]))
	return 0
}
var filebuf []byte
func Read(file, nfile, dst int32) int32 {
	if dst != 0 {
		copy(Bytes[dst:], filebuf)
		return 0
	}
	b, e := os.ReadFile(string(Bytes[file : file+nfile]))
	if e != nil {
		return -1
	}
	filebuf = b
	return int32(len(filebuf))
}
func Write(file, nfile, src, n int32) int32 {
	b := Bytes[src : src+n]
	if nfile == 0 {
		os.Stdout.Write(b)
		return 0
	}
	name := string(Bytes[file : file+nfile])
	e := os.WriteFile(name, b, 0644)
	if e != nil {
		return -1
	}
	return 0
}
func ReadIn(dst, n int32) int32 {
	nr, e := os.Stdin.Read(Bytes[dst : dst+n])
	if e != nil {
		return 0
	}
	if nr > 0 && Bytes[dst+int32(nr)-1] == 10 {
		nr -= 1
	}
	return int32(nr)
}
func Native(x, y int64) int64 { return x }
"

/ t:go t
/ \emit[]
