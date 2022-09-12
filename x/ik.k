
F:`!()
hx8:{[i]"0x",`x@|C i+!8}
lp:{(~~)#@[&#x;x i;i:!#x]}


ik:{[o]opts::o
 T::@[T;i@&(1=I j)&`for=T j:P P i:i@&`jmp=T i:lp P;`nop] /delete continue in simple loops
 
 F[`prg`mem`con`var`tab`fun`arg`sym`res`loc`ast`lod`sto]:(ikprg;ikmem;ikcon;ikvar;iktab;ikfun;ikarg;iksym;ikres;ikloc;ikast;iklod;iksto)
 F[`add`and`bnd`ant`asn`cal`cli`cnd`dfr`div`drp`eql`cst]:(ikadd;ikand;ikbnd;ikant;ikasn;ikcal;ikcli;ikcnd;ikdfr;ikdiv;ikdrp;ikeql;ikcst)
 F[`for`get`Get`gte`jmp`les`lit`lte`mod`mor`mul`neg    ]:(ikfor;ikget;ikGet;ikgte;ikjmp;ikles;iklit;iklte;ikmod;ikmor;ikmul;ikneg)
 F[`neq`nop`not`orr`bor`xor`ret`shl`shr`stm`sub`swc`typ]:(ikneq;iknop;iknot;ikorr;ikbor;ikxor;ikret;ikshl;ikshr;ikstm;iksub;ikswc;iktyp)
 
 S[0]:$[0N~i:o?`lib;S 0;o 1+i]
 n:(&#P){x[y]+:1}/P;n[0]-:1;
 s:{c:(0,(#x)-z)^x
   (c 0),,,/(F T y)[y;|c 1]}/[();|!#P;|n]
*s}


ikprg:{
 c:1_&P=0
 excl:!0
 fncs:"\n"/:y@&`fun=T c
 cons:"\n"/:y@&`con=T c
 vars:"\n"/:y@&`var=T c
 data:$[#D;"(data (i32.const 0) \"",(ikdat D),"\")";""]
 tabd:$[0<i:1+|/I c t:&`tab=T c;"(table (export \"table\") ",($i)," funcref)";""]
 tabl:iktbl y t
 "\n"/:(""~/:)_(cons;vars;data;fncs;tabd;tabl)}
 
ikdat:{,/@["\\",/:256^`x@_!256;0+ikc;ikc:_("a"+!26),("A"+!26),("0"+!10),"~!@#$%^&*()_+{}[]|;:<>,./?"]x}
ikmem:{y;""}
ikcon:{"(global $",($S x)," ",(iktype S 1+x)," (",(*y),"))"}
ikvar:{"(global $",($S x)," (mut ",(iktype S 1+x),") (",(*y),"))"}
iktab:{y;(I x;S x)}
iktbl:{"\n"/:{"(elem (i32.const ",($x),") func ",y,")"}'[(*x)i;" "/:'"$",''$(i:&~1=-':*x)^(x:+x)1]}
ikfun:{s:$S x
 args:";"/:y@&`arg=T c:&x=P
 rets:   ,/y@&`res=T c
 locs:" "/:y@&`loc=T c
 body:    y@*&`ast=T c
 expo:$[I x;"(export \"",s,"\")";""]
 s,":",rets,":{[",args,"]",body,"}"}
 
 
ikarg:{$*y}
iksym:{y;$[`"_"~s:S x;"";$s]}
ikres:{y;,/$S x}
ikloc:{"(local ",(*y)," ",(iktype S x),")"}
ikstm:{";"/:y}
ikast:{";"/:$[":"~**|y;y[l]:1_y[l:-1+#y];y]}

iklod:{$[`b~s:S x;"B";_($s)-32]," ",*y}
iksto:{$[`b~s:S x;"B";_($s)-32],"[",(*y),"]::",y 1}

ikant:{(y,'"\n"),t,".const -1\n",t,".xor\n",(t:iktype S x),".and"}

ikasn:{$[1~i:(I x);*y;,/"`",'i#y],":",*|y}
ikcal:{n:#y;y:";"/:y;$[`panic~s:S x;"!",y;($s),$[1~n;" ",y;"[",y,"]"]]}

/cli: .[rt;at;e][args..]  e.g. .[i;ii;1+x]x+y
ikcli:{n:I x;".[",(,/(1+2*n)_y),";",(","/:n#(1+n)_y),";",(*y),"]",$[1~n;y 1;"[",(";"/:n#1_y),"]"]}

/(1+xi)[i;ii][1;2]
/.[i;ii;1+xi][1;2]


ikdfr:{y;""}
ikdrp:{y}
ikcst:{$[(S x)~y 0;y 1;"`",($S x)," ",y 1]}

ikfor:{"(",(";"/:2#y),$[S x;";",$s;""],")/[",(y 2),"]"}

ikget:{y;$S x}
ikGet:{y;$S x}
ikjmp:{y;(,$[I x;"<";">"]),$[`~S x;"";"`",($S x)]}

ikflt:{(*&" "=c)#c:C 8+x+!24}
iklit:{y;i:I x;$[`f~t:S x;ikflt i
                 `i~t;$[0N~i;"-2147483648";$i]
                 `u~t;$[i<0;"0x",`x@|`c@,i;$i]
		 $[i~0;,"0";(-1~h:*`i C 4+i+!4)&0>j:*`i C i+!4; $j; (h~0)&j>0; $j; hx8 i]]}

ikneg:{"-",*y}
iknot:{(*y),"\ni32.eqz"}
iknop:{y;""}
ikret:{$[1~#y; ":",*y; ":(",(";"/:y),")"]}
ikswc:{
 r:(,"block")@&n:#$[I x;y;y,:,""]
 r[0],:$[S x;" (result ",(iktype S x),")";""]
 r,:,*y
 r,:,"br_table ",(" "/:$!n-1),"\nend"
 r,:(1_y),'"\nbr ",/:($|!n-1),\:"\nend"
 ("\n"/:r)}
 
ikcnd:{"$[",(";"/:y),"]"}

iktyp:{y;S x} /keep as symbol

ik2: {$[(T 1+y)':`get`Get`lit;z 0;"(",(z 0),")"],x,z 1}
ikn: {"~",ik2[x;z;y]}

ikadd:ik2"+";iksub:ik2"-";ikmul:ik2"*";ikdiv:ik2s"%";ikmod:ik2s"\\"
ikand:ik2"&";ikbnd:ik2"&&"
ikneq:ik2"~"
ikorr:ik2"|";ikbor:ik2"||";ikxor:ik2"^"
ikeql:ik2"="
ikshl:ik2"<<";ikshr:ik2">>"
ikmor:ik2">";ikles:ik2"<";ikgte:ikn"<";iklte;ikn">"


iktype:`i`u`j`k`f!("i32";"i32";"i64";"i64";"f64")
iksign:`i`u`j`k`f!("_s";  "_u"; "_s"; "_u";   "")
ikinst:`I32B`Memorysize`Memorygrow`Memorycopy`Memoryfill!("nop";"memory.size";"memory.grow";"memory.copy";"memory.fill")
ikinst,:`Memorysize2`Memorygrow2`Memorycopy2`Memorycopy3!((,"unreachable\ni32.const 0")0 0),(,"drop\ndrop\ndrop\nunreachable")0 0
ikinst,:{x!{(_32+*x),(1_3#x),".",3_x}'$x}`I32clz`F64abs`F64sqrt`F64floor`F64copysign`F64min`F64max`"F64reinterpret_i64"`"I64reinterpret_f64"



x:`<ik``nort

