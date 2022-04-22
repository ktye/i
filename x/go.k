
t:+`n`t`p`i`s!(!#T;T;P;I;S)

cld:{&x=t`p}
emit:{[l;i] ,/ (l(t`t)i)@'i }


goroot:{[i] c:t{p=0}
 mem:"Memory(",(","/:$(c{t=`mem})`i),")"
 data:{ "\n"/:{"Data(",($x),",`",(`x C[x+!y]),"`)"}'[x@i:2*!2\#x;(x@1+i)] }
 dat:data(t{t=`dat})`i
 tab:";"/:{ "F[",($*x),"]=",$(x 1) }'+(t{t=`te})`i`s
 glo:"\n"/:govar'(c{(t=`var)|t=`con})`n
 "\n"/:("program main"
        glo
        "func init(){"
        mem;dat;tab;"} ")}


goglo:{[i]"(glo)\n"}
gotab:{[i]"(tab)\n"}
gotyp:`i`I`j`J`f!("int32";"uint32";"int64";"uint64";"float64")
govar:{[i]"var ",($(t i)`s)," ",(gotyp@(t 1+i)`s),"=",golit t 1+i}
golit:{$[`i~x`s;$x`i;($x`s),$x`i]} // todo: init litJ1234=..
gofn:{[i]"func ",($(t`s)i),"\n"}
go:``glo`tab`fn!(goroot;goglo;gotab;gofn)
gomit:emit go

 `<emit[go;0]
