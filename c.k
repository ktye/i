/k:( `mk!{[r;x;y]j:i-4;(j>/4*t)?/(u:a+1^j_2;I[u]:I j;I[j]:u;j-:4);I[a]:y|x^29;I[a+4]:1;a}; 
/k:( `match!{[r;x;y]1+x>/y}
/ `mk!{[r;x;y](x>/y)?x+:4;x)

/k:60#k

funcs:{(,/!'x@i)!,/.'x@i:&,/7=@'x}
ftab: {x@*&,/&5=@'x}
data: {x@&,/6=@'k}
exprt:{*|x}
hex:  {,/{"0123456789abcdef"@,/(x%16;16/x)}'0+x}
obj:{ 'ox} 

/go a way
go:    {f:funcs x;o:gohead,(godata x),(gotab x),/gofunc'[!f;.'.f];o}
gofunc:{$[#y 1;""/:("func ";$x;"(";gosig y;")";goret y;"{";golocs y;gobody y;"}";_10);""]}
gosig: {a:1_x 2;n:-1+x 3;","/:(n#a){($x)," ",gotype y}'n#a}
goret: {$[`V~r:*x 2;"";"(",($r)," ",(gotype r),")"]}
gotype:{$["Z"<$x;"uint32";"float64"]}
golocs:{l:(x 3)_x 2;$[#l;(";"/:{"var ",($x)," ",gotype x}'l),";";""]}
gotab: {"func initab(){copy(T[0:],[]interface{}{", (","/:$(ftab x)) ,"})};"}
godata:{d:data x; "func initdat(){copy(C[",($*d),":],[]byte{",(","/:$0+d 1),"})};"}
gobody:{b:x 1;";"/:return[`V~*x 2;()gotree/b]}
gotree:{p:-1+#x;m:16/y;n:16\y
$[~y;
 0~m;  x,:,gocon@$%y
 1~m;  x[p]:$[164~n;gocond[|x p];128>n;goforc x p;"*"~n1:_n-128;gclz x p;"~"~n1;gonot x p;en gmon[go1;n1;x p]]
 2~m;  x:(-2_x),,$[63~n;goif;64~n;goat;46~n;gocal;gody(_n)][x p;x p-1]
 4~m;  x,:,$_n
 5~m;  x,:,$%y
 6~m;  x:((-n)_x),,(-n)#x
 7~m;  x:(-4_x),,goas -4#x
 9~m;  x[p]:"return ",x p
15~m;
 $y];  x }

goas:{$[x 1;gost x;(x 2),($[(#go2)~n:go2?,v:x 3;v;go2 1+n]),"=",x 0]}
gocond:{(5_,/{"else if ",(gobi x 0),"{",(";"/:rev x 1),"}"}'((2\#x),2)#x),"else{",(";"/:rev@*|x),"}"  }
goforc:{"for i:=uint32(0);i<",x,";i++"}
goforb:{x,"{",(";"/:rev y),"}"}
goat:{$[|/"CIF"~\:x;gold[x;y];|/go1~\:x;gmon[go1;x;y];"for"~3#x;goforb[x;y];x,en y]}
gocal:{$["?{"~x;gowhl y;|/(go2," ",goc)~\:x; gody[x;y 1;y 0]; "T"~x;goicl y; x,en@ ","/:rev y]}
goicl:{a:","/:(d:`I`F`V!("uint32";"float64";""))[,/`$'x 0];r:d@`$x 1;"T[",(x 2),"].(func(",a,")",r,")"}
gowhl:{"for ",(gobi x[1])," {",(";"/:rev@*x),"}"}
gobi:{$["ib("~3#x;2_x;"n32(ib("~7#x;"!((",7_x;"n32"~3#x;"1==",x;"bi(",x,")"]}
gonot:{"n32(",x,")"}
gold:{x:x,"[",(en y),(shift x),"]";$["C"~*x;"uint32(",x,")";x]}
gost:{(x 2),"[",(en x 1),(shift x 2),"]=",$["C"~x 2;"byte(",(x 0),")";x 0]}
goif:{"if ",(gobi x),"{",(";"/:rev y),"}"}
go1:" "\:"+[ math.Abs( -[ - - - %[ math.Sqrt( ! trap( ~ ! _\ math.Floor("
go2:" "\:"% / ! % !( %' ^ << | ^ _ >> $( $' +[ + + + -[ - - - *[ * * * %( /' %[ /"
goc:" "\:"~ == ~[ == < < <( <' <{ <= <{( <=' <[ < >( >' >{ >= >{( >=' >{[ >= >[ > > >"
gclz:{"uint32(bits.LeadingZeros32(",x,"))"}
gmon:{s:$[(#x)~n:x?,y; y; x 1+n];$["("~*|s;s,z,")";s,z]}
gody:{[v;x;y]$[v~":"; :x,"=",y;v~"$"; :gocvt[x;y;""];v~"$("; :gocvt[x;y;"'"]; 0]; ib:""; l:$[(#goc)~n:goc?,v;go2;*(goc;ib:"ib")]; s:$[(#l)~n:l?,v;v;l 1+n]; ib:gosnb[s;ib]; ib,"(",gosn[s;x],gosn[s;s],gosn[s;y],")"}
gocvt:{$["I"~x;"uint32";"float64"],"(",gosn[z;y],")"}
gocon:{$["0n"~x;"math.Float64frombits(18444492273895866368)";"0w"~x;"math.Inf(1)";x~"3.141592";"math.Pi";x~"57.29577";"57.29577951308232";x~"ffffffffffffef7f";"1.7976931348623157e+308";x~"399d52a246df913f";"0.017453292519943295";x]}
gosn:{o:$["'"~*|x;-1_x;x];$[x~y;o;$[o~x;y;"int32(",y,")"]]}
gosnb:{$[~""~y;y;"'"~*|x;"uint32";""]}
en:{"(",x,")"}
ne:{$["("~*x;1_-1_x;x]}
shift:{$["I"~x;">>2";"F"~x;">>3";""]}
rev:{$[6~@x;|x;x]}

return:{$[x;:y;#y;;:y];p:-1+#y;$["return"~6#y p;y;y[p]:"return ",y[p]];y}

gohead:(_34)/:"`"\:"package main
import (`unsafe`;`fmt`;`bufio`;`os`;`math`;`math/bits`;`io/ioutil`;`bytes`;`strings`)
var C []byte;var I []uint32;var F []float64;var T [256]interface{}
func init(){initab();C=make([]byte,1<<16);msl();initdat()}
func trap(x uint32){panic(`trap`)}
func n32(x uint32)uint32{if x==0{return 1}else{return 0}}
func bi(x uint32)bool{if 0==x{return false}else{return true}}
func ib(x bool)uint32{if x{return 1}else{return 0}}
func grow(x uint32)uint32{if x>31{panic(`oom`)};c:=make([]byte,1<<x);copy(c,C);C=c;msl();return x}
var sin,cos,exp,atan2,hypot,log=math.Sin,math.Cos,math.Exp,math.Atan2,math.Hypot,math.Log
func draw(x, y, z uint32){fmt.Printf(`draw %x %x %x\n`, x, y, z)}
func printc(x, y uint32){fmt.Printf(`%s\n`, string(C[x:x+y]))}
type sl struct {p uintptr;l,c int}
func msl() {
 cp := *(*sl)(unsafe.Pointer(&C));ip := *(*sl)(unsafe.Pointer(&I));fp := *(*sl)(unsafe.Pointer(&F))
 ip.l, ip.c, ip.p = cp.l/4, cp.c/4, cp.p; fp.l, fp.c, fp.p = ip.l/2, ip.c/2, ip.p
 C = *(*[]byte)(unsafe.Pointer(&cp));I = *(*[]uint32)(unsafe.Pointer(&ip));F = *(*[]float64)(unsafe.Pointer(&fp))}
func mkcs(b []byte)uint32{x:=mk(1,uint32(len(b)));copy(C[x+8:],b);return x}
func fatal(e error){if e!=nil{panic(e)}}
func load(f string){b,e:=ioutil.ReadFile(f);fatal(e);if n:=bytes.Index(b,[]byte(`\n\\`));n!=-1{b=b[:n+1]};dx(out(val(mkcs(b))))}
func main(){
 ini(16);
 for _, a := range os.Args[1:]{if strings.HasSuffix(a,`.k`){load(a)};if a==`-e`{os.Exit(0)}}
 s:=bufio.NewScanner(os.Stdin);fmt.Printf(`k.k(go)\n `);
 for s.Scan(){t:=s.Text();if t==`\\\\`{os.Exit(0)};
  dx(out(val(mkcs([]byte(t)))));os.Stdout.Write([]byte{32})
 }
}

func dump(a, n uint32) uint32 {
	p := a >> 2
	fmt.Printf(`%.8x `, a)
	for i := uint32(0); i < n; i++ {
		x := I[p+i]
		fmt.Printf(` %.8x`, x)
		if i > 0 && (i+1)%8 == 0 {
			fmt.Printf(`\n%.8x `, a+4*i+4)
		} else if i > 0 && (i+1)%4 == 0 {
			fmt.Printf(` `)
		}
	}
	fmt.Println()
	return 0
}
"