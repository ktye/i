/copied
F:`!()
emit:{[o]opts::o
 S[0]:$[0N~i:o?`lib;S 0;o 1+i]
 n:(&#P){x[y]+:1}/P;n[0]-:1;
 s:{c:(0,(#x)-z)^x
   (c 0),,,/(F T y)[y;|c 1]}/[();|!#P;|n]
*s}


nl:{$[#x:"\n"/:x;x,"\n";x]} / \

waprg:{
 fncs:"\n"/:y@&`fun=T c:1_&P=0
 "\n"/:("(module";fncs;")\n")}
wamem:{y;0+`nyi}
wacon:{y;0+`nyi}
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
walod:{y;0+`nyi}
wasto:{y;0+`nyi}
waand:{y;0+`nyi}
wabnd:{y;0+`nyi}
waant:{y;0+`nyi}
waasn:{                   /todo global, modified, multi
 n:I x; $[`~`S x;0+`modasn;~1~I x;0+`multiasn; (y 2),"\n","local.set ",y 0]}
wacal:{("\n"/:y),"\n",$[#c:wainst s:S x;c;"call $",$s]}
wacli:{y;0+`nyi}
wacnd:{y;0+`nyi}
wadfr:{y;0+`nyi}
wadiv:{y;0+`nyi}
wadrp:{"\n"/:y,,"drop"}
waeql:{y;0+`nyi}
wacst:{t:!watype;(y 1),"\n",$[0~r:wacast i:(t?y 0)+8*t?S x;`castyp+ \i;r]}
wafor:{y;0+`nyi}
waget:{y;"local.get $",$S x}
waGet:{y;"global.get $",$S x}
wagte:{y;0+`nyi}
wajmp:{y;0+`nyi}
wales:{y;0+`nyi}

waflt:{(*&" "=c)#c:C 8+x+!24}                        /todo: wasm literals
walit:{y;i:I x;(watype t),".const ",$[`f~t:S x;waflt i
                 `i~t;$[0N~i;"-0x80000000";$i]
                 `u~t;$[i<0;"0x",`x@|`c@,i;$i]
		 `j~t;$[i~0;,"0";(-1~h:*`i C 4+i+!4)&0>j:*`i C i+!4; $j; (h~0)&j>0; $j; hx8 i]
		 `k~t;$[i~0;,"0";hx8 i];"???"]}

walte:{y;0+`nyi}
wamod:{y;0+`nyi}
wamor:{y;0+`nyi}
wamul:{y;0+`nyi}
waneg:{y;0+`nyi}
waneq:{y;0+`nyi}
wanop:{y;""} /or drop?
wanot:{y;0+`nyi}
waorr:{y;0+`nyi}
wabor:{y;0+`nyi}
waxor:{y;0+`nyi}
waret:{("\n"/:y),"\nreturn"}
washl:{y;0+`nyi}
washr:{y;0+`nyi}
wasub:{y;0+`nyi}
waswc:{y;0+`nyi}
watyp:{y;S x} /keep as symbol

wa2:{[z;x;y]"\n"/:(y 0;y 1;(watype S x),z)}
waadd:wa2".add"

watype:`i`u`j`k`f`C`I`F!("i32";"i32";"i64";"i64";"f64";"v128";"v128";"v128")
wacast:@[&64;1 2 3 4 8 16 24 32;("";"i32.wrap_i64";"i32.wrap_i64";"i32.trunc_f64_s";"";"i64.extend_i32_s";"i64.extend_i32_u";"f64.convert_i32_s")]  /wg/wat.go:wasmcst
wainst:`Memorysize`Memorygrow`Memorycopy`Memoryfill!("memory.size";"memory.grow";"memory.copy";"memory.fill")


/todo:
/TEE: local.set $x local.get $x => local.tee $x
/RET: local.tee $r return => return, also remove (local $r typ)

wa:{
 F[`prg`mem`con`var`tab`fun`arg`sym`res`loc`ast`lod`sto]:(waprg;wamem;wacon;wavar;watab;wafun;waarg;wasym;wares;waloc;waast;walod;wasto)
 F[`add`and`bnd`ant`asn`cal`cli`cnd`dfr`div`drp`eql`cst]:(waadd;waand;wabnd;waant;waasn;wacal;wacli;wacnd;wadfr;wadiv;wadrp;waeql;wacst)
 F[`for`get`Get`gte`jmp`les`lit`lte`mod`mor`mul`neg    ]:(wafor;waget;waGet;wagte;wajmp;wales;walit;walte;wamod;wamor;wamul;waneg)
 F[`neq`nop`not`orr`bor`xor`ret`shl`shr`stm`sub`swc`typ]:(waneq;wanop;wanot;waorr;wabor;waxor;waret;washl;washr;wastm;wasub;waswc;watyp)
 x}


