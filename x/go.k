
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
nyi:{y;em x}

/goprg:{[i] c:t@&P=0
/ mem:"Memory(",(","/:$(c{t=`mem})`i),")"
/ data:{ "\n"/:{"Data(",($x),",`",(`x C[x+!y]),"`)"}'[x@i:2*!2\#x;(x@1+i)] }
/ dat:data I@&T=`dat
/ tab:";"/:{ "F[",($*x),"]=",$(x 1) }'+(t{t=`tab})`i`s
/ con:"\n"/:gocon'(c{t=`con})`n
/ glo:"\n"/:govar'(c{t=`var})`n
/ fun:"\n"/:gofun'(c{t=`fun})`n
/ "\n"/:("program main"
/        con;glo
/        "func init(){"
/        mem;dat;tab;" }"
/	fun;gort)}
/gotyp:`i`u`j`k`f`C`I`F!("int32";"uint32";"int64";"uint64";"float64";"[8]int8";"[32]int32";"[2]float64")
/gocon:{[i]"const ",($S i)," ",(gotyp S 1+i),"=",golit 1+i}
/govar:{[i]"var ",  ($S i)," ",(gotyp S 1+i),"=",golit 1+i}
/golit:{[j]i:I j;$[`i~tp:S j;$[0N~i;"-0x80000000";$i] /x
/ `I~tp;$[i<0;"0x",|`x@,i;$i]
/ `J~tp;$[i~0;,"0";hx8 i]
/ `j~tp;$[i~0;,"0";0<`i C 4+i+!4;hx8 i;"??"]
/ `f~tp;$[i~0;"0.0";"0x400921fb54442d18"~h:(hx8 i);"math.Pi";h~"0x7fefffffffffffff";"math.MaxFloat64";"math.Float64frombits(",h,")"]
/ ($s),$i]}
/gofun:{[i]
/ arg:{[i]n:&(P=i)&T=`arg; ","/:{x," ",gotyp y}'[$S 1+n; S n]}
/ res:{[i]","/:gotyp S@&(P=i)&T=`res}
/ var:{[i]"\n"/:{y;"var ",(gotyp x)," ",","/:$S 1+y`n}'[!g;.g:`s?t@&(P=i)&T=`loc]}
/ bdy:gox@&(P=i)&T=`ast
/ "func ",($S i),"(",(arg i),")(",(res i),"){\n",(var i),"\n",bdy,"\n}"}


gotype:`i`u`j`k`f`C`I`F!("int32";"uint32";"int64";"uint64";"float64";"[8]int8";"[32]int32";"[2]float64")

goprg:{
 fncs:"\n"/:y@&`fun=T c:1_&P=0
 cons:"\n"/:y@&`con=T c
 vars:"\n"/:y@&`var=T c
 runt:$[`nort':opts;"";gornt]
 ("package ";$S 0;"\n";cons;vars;runt),,fncs}
 
gofun:{
 args:", "/:y@&`arg=T c:&x=P
 rets:", "/:y@&`res=T c
 locs:"\n"/:"var ",/:(.g),'" ",/:gotype@!g:","/:'(y i)@=(S c)i:&`loc=T c
 defr:     *y@&`dfr=T c
 body:      y@*&`ast=T c
 ("func ";$S x;"(";args;")(";rets;"){";defr;locs,"\n";body;"}")}
 

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
goneg:go1"-";gonot:go1"!"
goadd:go2"+";gosub:go2"-"; gomul:go2"*";godiv:go2"/"; gomod:go2"%"
goeql:go2"=";gogte:go2">=";goles:go2"<";golte:go2"<=";gomor:go2">"; goneq:go2"!="
goand:go2"&";goant:go2"&^";goorr:go2"|";goxor:go2"^"; goshl:go2"<<";goshr:go2">>"

go:{
 F[`prg`mem`con`var`tab`dat`fun`arg`sym`res`loc`ast]:(goprg;gomem;gocon;govar;gotab;godat;gofun;goarg;gosym;gores;goloc;goast)
 F[`add`and`ant`asn`cal`cli`cnd`dfr`div`drp`eql`cst]:(goadd;goand;goant;goasn;gocal;gocli;gocnd;godfr;godiv;godrp;goeql;gocst)
 F[`for`get`Get`gte`jmp`les`lit`lte`mod`mor`mul`neg]:(gofor;goget;goGet;gogte;gojmp;goles;golit;golte;gomod;gomor;gomul;goneg)
 F[`neq`nop`not`orr`ret`shl`shr`stm`sub`swc`typ    ]:(goneq;gonop;gonot;goorr;goret;goshl;goshr;gostm;gosub;goswc;gotyp      )
 x}


gornt:"/* rt */\n"

/ t:go t
/ \emit[]
