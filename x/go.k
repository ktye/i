
t:+`n`t`p`i`s!(!#T;T;P;I;S)

emit:{[l;i] ,/ (l(t`t)i)@'i }
/cld:{&x=t`p}
hx8:{[i]"0x",`x@|C i+!8}



goroot:{[i] c:t{p=0}
 mem:"Memory(",(","/:$(c{t=`mem})`i),")"
 data:{ "\n"/:{"Data(",($x),",`",(`x C[x+!y]),"`)"}'[x@i:2*!2\#x;(x@1+i)] }
 dat:data(t{t=`dat})`i
 tab:";"/:{ "F[",($*x),"]=",$(x 1) }'+(t{t=`tab})`i`s
 con:"\n"/:gocon'(c{t=`con})`n
 glo:"\n"/:govar'(c{t=`var})`n
 fun:"\n"/:gofun'(c{t=`fun})`n
 "\n"/:("program main"
        con;glo
        "func init(){"
        mem;dat;tab;,"}"
	fun)}


gotab:{[i]"(tab)\n"}
gotyp:`i`I`j`J`f!("int32";"uint32";"int64";"uint64";"float64")
gocon:{[i]"const ",($(t i)`s)," ",(gotyp@(t 1+i)`s),"=",golit t 1+i}
govar:{[i]"var ",  ($(t i)`s)," ",(gotyp@(t 1+i)`s),"=",golit t 1+i}
golit:{i:x`i;s:x`s;$[`i~tp:s;$[0N~i;"-0x80000000";$i]
 `I~tp;$[i<0;"0x",|`x@,i;$i]
 `J~tp;$[i~0;,"0";hx8 i]
 `j~tp;$[i~0;,"0";0<`i C 4+i+!4;hx8 i;"??"]
 `f~tp;$[i~0;"0.0";"0x400921fb54442d18"~h:(hx8 i);"math.Pi";h~"0x7fefffffffffffff";"math.MaxFloat64";"math.Float64frombits(",h,")"]
 ($s),$i]}
gofun:{[i]"func ",($(t`s)i),"\n"}
go:``tab`fn!(goroot;gotab;gofn)
gomit:emit go


`<emit[go;0]
