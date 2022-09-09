/webassembly compiler (text format)
F:`!()
hx8:{[i]"0x",`x@|C i+!8}
lp:{(~~)#@[&#x;x i;i:!#x]}


wa:{[o]opts::o
 RCND::i@&{`ret`ret~T@1_&P=x}@/:i:&T=`cnd  /these cnd(if-else) nodes always return
 T::@[T;i@&(1=I j)&`for=T j:P P i:i@&`jmp=T i:lp P;`nop] /delete continue in simple loops
 
 F[`prg`mem`con`var`tab`fun`arg`sym`res`loc`ast`lod`sto]:(waprg;wamem;wacon;wavar;watab;wafun;waarg;wasym;wares;waloc;waast;walod;wasto)
 F[`add`and`bnd`ant`asn`cal`cli`cnd`dfr`div`drp`eql`cst]:(waadd;waand;wabnd;waant;waasn;wacal;wacli;wacnd;wadfr;wadiv;wadrp;waeql;wacst)
 F[`for`get`Get`gte`jmp`les`lit`lte`mod`mor`mul`neg    ]:(wafor;waget;waGet;wagte;wajmp;wales;walit;walte;wamod;wamor;wamul;waneg)
 F[`neq`nop`not`orr`bor`xor`ret`shl`shr`stm`sub`swc`typ]:(waneq;wanop;wanot;waorr;wabor;waxor;waret;washl;washr;wastm;wasub;waswc;watyp)
 
 S[0]:$[0N~i:o?`lib;S 0;o 1+i]
 n:(&#P){x[y]+:1}/P;n[0]-:1;
 s:{c:(0,(#x)-z)^x
   (c 0),,,/(F T y)[y;|c 1]}/[();|!#P;|n]
*s}


waprg:{
 c:1_&P=0
 excl:&(`fun=T c)&(`store=S c)|`catch=S c
 fncs:"\n"/:y@excl_&`fun=T c
 cons:"\n"/:y@&`con=T c
 vars:"\n"/:y@&`var=T c
 data:$[#D;"(data (i32.const 0) \"",(wadat D),"\")";""]
 tabl:"\n"/:y@&`tab=T c
 mems:"(memory (export \"memory\") 1)"
 "\n"/:("(module";mems;data;cons;vars;fncs;tabl;")\n")}
 
wadat:{,/@["\\",/:256^`x@_!256;0+wac;wac:_("a"+!26),("A"+!26),("0"+!10),"~!@#$%^&*()_+{}[]|;:<>,./?"]x}
wamem:{y;""}
wacon:{" "/:("(global"; "$",$S x; watype S 1+x; "(",(*y),"))")} /export?
wavar:{" "/:("(global"; "$",$S x; "(mut"; watype S 1+x; ")(",(*y),"))")}
watab:{y;"(elem (i32.const ",($I x),") func $",($S x),")"}
wafun:{s:$S x
 args:"",/y@&`arg=T c:&x=P
 rets:"",/y@&`res=T c
 locs:"",/y@&`loc=T c
 body:    y@*&`ast=T c
 "(func $",s," (export \"",s,"\")",args,rets,locs,"\n",body,")"}
waarg:{"(param ",(y 0)," ",(watype S x),")"}
wasym:{y;"$",$S x}
wares:{y;"(result ",(watype S x),")"}
waloc:{"(local ",(*y)," ",(watype S x),")"}
wastm:{"\n"/:y}
waast:wastm
walod:{(y 0),"\n",$[`b~s:S x;"i32.load8_s";(watype s),".load"]}
wasto:{"\n"/:y,,$[`b~s:S x;"i32.store8";(watype s),".store"]}
waant:{(y,'"\n"),t,".const -1\n",t,".xor\n",(t:watype S x),".and"}
waasn:{y[&b:"$_"~/:-1_y]:,"";(*|y),,/|(("\nlocal.set ";"\nglobal.set ";"\ndrop")(2*b)|I 1+x+!I x),'-1_y}
wacal:{("\n"/:y),"\n",$[#c:wainst s:S x;c;"call $",$s]}
wacli:{("\n"/:((I x)#1_y),,*y),"\n","call_indirect ",/(1+I x)_y}
wadfr:{y;""}
wadrp:{"\n"/:y,,"drop"}
wacst:{t:!watype;(y 1),"\n",wacast i:(t?y 0)+5*t?S x}
wafor:{$[I x;waslp[x;y];walop[x;y]]}
waslp:{"\n"/:("loop";"\n"/:2_y;y 1;*y;"br_if 0\nend")}
walop:{bl:("block";"loop"),'$[`~S x;2^"";(l,"1";(l:" $",$S x),"0")]
 "\n"/:bl,($[#*y;(*y),"\ni32.eqz\nbr_if 1";""]; "\n"/:2_y;y 1;"br 0\nend\nend")}
waget:{y;"local.get $",$S x}
waGet:{y;"global.get $",$S x}
wajmp:{y;"br ",$[`~S x;$1+I x;"$",($S x),$I x]}

waflt:{(*&" "=c)#c:C 8+x+!24}                        /todo: wasm literals
walit:{y;i:I x;(watype t),".const ",$[`f~t:S x;waflt i
                 `i~t;$[0N~i;"-0x80000000";$i]
                 `u~t;$[i<0;"0x",`x@|`c@,i;$i]
		 `j~t;$[i~0;,"0";(-1~h:*`i C 4+i+!4)&0>j:*`i C i+!4; $j; (h~0)&j>0; $j; hx8 i]
		 `k~t;$[i~0;,"0";hx8 i];"???"]}

waneg:{$[`f~t:S x;(*y),"\nf64.neg";"\n"/:(t,".const 0";(*y);(t:watype t),".sub")]}
wanot:{(*y),"\ni32.eqz"}
wanop:{y;""}

waret:{("\n"/:y),"\nreturn"}

waswc:{ /todo: transform to expression-block with return type in (S x)
 r:(,"block")@&n:#y
 r,:,*y
 r,:,"br_table ",(" "/:$!n),"\nend"
 r,:(1_y),'"\nbr ",/:($|!n-1),\:"\nend"
 "\n"/:r}
wacnd:{ /todo: result-type transformation (same as swc)
 
 "\n"/:(*y;$[#t:watype S x;"if (result ",t,")";"if"];y 1;$[3~#y;"else\n",(y 2);""],"\nend")}
watyp:{y;S x} /keep as symbol

wa2: {"\n"/:(z 0;z 1;(watype S y),x)}
wa2s:{"\n"/:(z 0;z 1;(watype t  ),x,wasign@t:S y)}
waadd:wa2".add";wasub:wa2".sub"
wamul:wa2".mul";wadiv:wa2s".div";wamod:wa2s".rem"
waand:wa2".and";wabnd:wa2".and"
waneq:wa2".ne"
wamor:wa2s".gt";wagte:wa2s".ge";wales:wa2s".lt";walte:wa2s".le";
waorr:wa2".or";wabor:waorr;waxor:wa2".xor"
waeql:wa2".eq"
washl:wa2".shl";washr:wa2s".shr"


watype:`i`u`j`k`f!("i32";"i32";"i64";"i64";"f64")
wasign:`i`u`j`k`f!("_s";  "_u"; "_s"; "_u";   "")
wainst:`I32B`Memorysize`Memorygrow`Memorycopy`Memoryfill!("nop";"memory.size";"memory.grow";"memory.copy";"memory.fill")
wainst,:{x!{(_32+*x),(1_3#x),".",3_x}'$x}`I32clz`F64abs`F64sqrt`F64floor`F64copysign`F64min`F64max`"F64reinterpret_i64"`"I64reinterpret_f64"

/         i                   u   j              k              f
wacast:("";                 "";"i32.wrap_i64";"i32.wrap_i64";"i32.trunc_f64_s"  /i
        "";                 "";""            ;""            ;""                 /u
        "i64.extend_i32_s"; "";""            ;"i32.wrap_i64";"i64.trunc_f64_s"  /j
        "i64.extend_i32_u"; "";""            ;""            ;"i64.trunc_f64_u"  /k
        "f64.convert_i32_s";"";""            ;""            ;"")                /f


/todo:
/TEE: local.set $x local.get $x => local.tee $x
/RET: local.tee $r return => return, also remove (local $r typ)



