/copied
F:`!()
emit:{[o]opts::o
 S[0]:$[0N~i:o?`lib;S 0;o 1+i]
 n:(&#P){x[y]+:1}/P;n[0]-:1;
 s:{c:(0,(#x)-z)^x
   (c 0),,,/(F T y)[y;|c 1]}/[();|!#P;|n]
*s}
hx8:{[i]"0x",`x@|C i+!8}

waprg:{
 fncs:"\n"/:y@&`fun=T c:1_&P=0
 cons:"\n"/:y@&`con=T c
 mems:"(memory (export \"memory\") 1)"
 "\n"/:("(module";mems;cons;fncs;")\n")}
wamem:{y;""}
wacon:{" "/:("(global"; "$",$S x; watype S 1+x; "(",(*y),"))")} /export?
wavar:{y;0+`nyi}
watab:{y;0+`nyi}
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
walod:{(y 0),"\n",(watype S x),".load"}
wasto:{"\n"/:y,,(watype S x),".store"}
waand:{y;0+`nyi}
wabnd:{y;0+`nyi}
waant:{y;0+`nyi}
waasn:{                   /todo global, modified, multi
 n:I x; $[`~`S x;0+`modasn;~1~I x;0+`multiasn; (y 2),"\n","local.set ",y 0]}
wacal:{("\n"/:y),"\n",$[#c:wainst s:S x;c;"call $",$s]}
wacli:{y;0+`nyi}
wadfr:{y;0+`nyi}
wadiv:{y;0+`nyi}
wadrp:{"\n"/:y,,"drop"}
wacst:{t:!watype;(y 1),"\n",$[0~r:wacast i:(t?y 0)+8*t?S x;`castyp+ \i;r]}
wafor:{y;0+`nyi}
waget:{y;"local.get $",$S x}
waGet:{y;"global.get $",$S x}
wagte:{y;0+`nyi}
wajmp:{y;0+`nyi}

waflt:{(*&" "=c)#c:C 8+x+!24}                        /todo: wasm literals
walit:{y;i:I x;(watype t),".const ",$[`f~t:S x;waflt i
                 `i~t;$[0N~i;"-0x80000000";$i]
                 `u~t;$[i<0;"0x",`x@|`c@,i;$i]
		 `j~t;$[i~0;,"0";(-1~h:*`i C 4+i+!4)&0>j:*`i C i+!4; $j; (h~0)&j>0; $j; hx8 i]
		 `k~t;$[i~0;,"0";hx8 i];"???"]}

walte:{y;0+`nyi}
wamod:{y;0+`nyi}
waneg:{$[`f~t:S x;(*y),"\nf64.neg";"\n"/:(t,".const 0";(*y);(t:watype t),".sub")]}
wanop:{y;""} /or drop?
wanot:{y;0+`nyi}
waorr:{y;0+`nyi}
wabor:{y;0+`nyi}
waxor:{y;0+`nyi}
waret:{("\n"/:y),"\nreturn"}
washl:{y;0+`nyi}
washr:{y;0+`nyi}
wasub:{y;0+`nyi}
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
waadd:wa2".add"
wamul:wa2".mul"
waneq:wa2".ne"
wamor:wa2s".gt"
wales:wa2s".lt"
waeql:wa2".eq"

watype:`i`u`j`k`f`C`I`F!("i32";"i32";"i64";"i64";"f64";"v128";"v128";"v128")
wasign:`i`u`j`k`f`C`I`F!("_s";  "_u"; "_s"; "_u";   "";    "";    "";    "")
wacast:@[&64;0 1 2 3 4 8 16 19 24 32;("";"";"i32.wrap_i64";"i32.wrap_i64";"i32.trunc_f64_s";"";"i64.extend_i32_s";"i32.wrap_i64";"i64.extend_i32_u";"f64.convert_i32_s")]  /wg/wat.go:wasmcst
wainst:`Memorysize`Memorygrow`Memorycopy`Memoryfill!("memory.size";"memory.grow";"memory.copy";"memory.fill")


/todo:
/TEE: local.set $x local.get $x => local.tee $x
/RET: local.tee $r return => return, also remove (local $r typ)

wa:{
 
 RCND::i@&{`ret`ret~T@1_&P=x}@/:i:&T=`cnd  /these cnd(if-else) nodes always return
 
 F[`prg`mem`con`var`tab`fun`arg`sym`res`loc`ast`lod`sto]:(waprg;wamem;wacon;wavar;watab;wafun;waarg;wasym;wares;waloc;waast;walod;wasto)
 F[`add`and`bnd`ant`asn`cal`cli`cnd`dfr`div`drp`eql`cst]:(waadd;waand;wabnd;waant;waasn;wacal;wacli;wacnd;wadfr;wadiv;wadrp;waeql;wacst)
 F[`for`get`Get`gte`jmp`les`lit`lte`mod`mor`mul`neg    ]:(wafor;waget;waGet;wagte;wajmp;wales;walit;walte;wamod;wamor;wamul;waneg)
 F[`neq`nop`not`orr`bor`xor`ret`shl`shr`stm`sub`swc`typ]:(waneq;wanop;wanot;waorr;wabor;waxor;waret;washl;washr;wastm;wasub;waswc;watyp)
 x}


