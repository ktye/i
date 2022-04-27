
t:+`n`t`p`i`s!(!#T;T;P;I;S)


emit:{[l;i] ,/ (l T i)@'i }                   / emit[target-lang;node-number]
emit:{[l;i] ,/ (l@`T \T@`i \i)@'i }
emit:{[l;i] ,/ (l T i)@'i }

/cld:{&x=P}
hx8:{[i]"0x",`x@|C i+!8}
em:{"(",x,")"}

OO:``
O:{OO::?OO,x}
nyi:{y;em x}

goprg:{[i] c:t@&P=0
 mem:"Memory(",(","/:$(c{t=`mem})`i),")"
 data:{ "\n"/:{"Data(",($x),",`",(`x C[x+!y]),"`)"}'[x@i:2*!2\#x;(x@1+i)] }
 dat:data I@&T=`dat
 tab:";"/:{ "F[",($*x),"]=",$(x 1) }'+(t{t=`tab})`i`s
 con:"\n"/:gocon'(c{t=`con})`n
 glo:"\n"/:govar'(c{t=`var})`n
 fun:"\n"/:gofun'(c{t=`fun})`n
 "\n"/:("program main"
        con;glo
        "func init(){"
        mem;dat;tab;" }"
	fun;gort)}
gotyp:`i`u`j`k`f`C`I`F!("int32";"uint32";"int64";"uint64";"float64";"[8]int8";"[32]int32";"[2]float64")
gocon:{[i]"const ",($S i)," ",(gotyp S 1+i),"=",golit 1+i}
govar:{[i]"var ",  ($S i)," ",(gotyp S 1+i),"=",golit 1+i}
golit:{[j]i:I j;$[`i~tp:S j;$[0N~i;"-0x80000000";$i] /x
 `I~tp;$[i<0;"0x",|`x@,i;$i]
 `J~tp;$[i~0;,"0";hx8 i]
 `j~tp;$[i~0;,"0";0<`i C 4+i+!4;hx8 i;"??"]
 `f~tp;$[i~0;"0.0";"0x400921fb54442d18"~h:(hx8 i);"math.Pi";h~"0x7fefffffffffffff";"math.MaxFloat64";"math.Float64frombits(",h,")"]
 ($s),$i]}
gofun:{[i]
 arg:{[i]n:&(P=i)&T=`arg; ","/:{x," ",gotyp y}'[$S 1+n; S n]}
 res:{[i]","/:gotyp S@&(P=i)&T=`res}
 var:{[i]"\n"/:{y;"var ",(gotyp x)," ",","/:$S 1+y`n}'[!g;.g:`s?t@&(P=i)&T=`loc]}
 bdy:gox@&(P=i)&T=`ast
 "func ",($S i),"(",(arg i),")(",(res i),"){\n",(var i),"\n",bdy,"\n}"}

goast:{i:&P=x;O T i;gox i}
gostm:goast
goget:{$S x}

gocal:nyi"cal"
gocli:nyi"cli"
goasn:nyi"asn"
gocnd:{[i]c:&P=i;"if",(gox c 0),"{",(gox c 1),"}"}
gonop:nyi"nop"
gojmp:nyi"jmp"
goswc:nyi"swc"
gofor:nyi"for"
godfr:nyi"dfr"
godrp:nyi"drp"

goret:{"return ",goast x}

go1:{x,gox@&P=y}
go2:{i:&P=y;em (gox i 0),x,gox i 1}
goneg:go1"-";gonot:go1"!"
goadd:go2"+";gosub:go2"-"; gomul:go2"*";godiv:go2"/"; gomod:go2"%"
goeql:go2"=";gogte:go2">=";goles:go2"<";golte:go2"<=";gomor:go2">"; goneq:go2"!="
goand:go2"&";goant:go2"&^";goorr:go2"|";goxor:go2"^"; goshl:go2"<<";goshr:go2">>"

go:`prg`add`and`ant`asn`ast`cal`cli`cnd`dfr`div`drp!(goprg;goadd;goand;goant;goasn;goast;gocal;gocli;gocnd;godfr;godiv;godrp)
go,:   `eql`for`get`Get`gte`jmp`les`lit`lte`mod`mor!(goeql;gofor;goget;goget;gogte;gojmp;goles;golit;golte;gomod;gomor)
go,:   `mul`neg`nop`not`orr`ret`shl`shr`stm`sub`swc!(gomul;goneg;gonop;gonot;goorr;goret;goshl;goshr;gostm;gosub;goswc)

gox:emit go
gort:"/* rt */"


x:gox 0