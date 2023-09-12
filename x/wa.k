/webassembly compiler (text format)      // \\ ': /: \:
F:(0#`)!()
hx8:{[i]"0x",`x@|C i+!8}
lp:{(~~)#@[&#x;x i;i:!#x]}


wa:{[o]opts::o
 /T::@[T;i@&(1=I j)&`for=T j:P P i:i@&`jmp=T i:lp P;`nop] /delete continue in simple loops
 
 F[`prg`mem`con`var`tab`fun`arg`sym`res`loc`ast`lod`sto]:(waprg;wamem;wacon;wavar;watab;wafun;waarg;wasym;wares;waloc;waast;walod;wasto)
 F[`add`and`bnd`ant`asn`cal`cli`cnd`dfr`div`drp`eql`cst]:(waadd;waand;wabnd;waant;waasn;wacal;wacli;wacnd;wadfr;wadiv;wadrp;waeql;wacst)
 F[`for`get`Get`gte`jmp`les`lit`lte`mod`mor`mul`neg    ]:(wafor;waget;waGet;wagte;wajmp;wales;walit;walte;wamod;wamor;wamul;waneg)
 F[`neq`nop`not`orr`bor`xor`ret`shl`shr`stm`sub`swc`typ]:(waneq;wanop;wanot;waorr;wabor;waxor;waret;washl;washr;wastm;wasub;waswc;watyp)
 
 S[0]:$[0N~i:o?`lib;S 0;o 1+i]
 n:{x[y]+:1}/(,&#P),P;n[0]-:1
 s:{c:(0,(#x)-z)^x
   (c 0),,,/(F T y)[y;|c 1]}/[();|!#P;|n]
*s}


waprg:{
 c:1_&P=0
 excl:!0
 fncs:"\n"/y@&`fun=T c
 cons:"\n"/y@&`con=T c
 vars:"\n"/y@&`var=T c
 sysi:$[`nort?opts;"";wasys]
 data:$[#D;"(data (i32.const 0) \"",(wadat D),"\")";""]
 tabd:$[0<i:1+|/I c t:&`tab=T c;"(table (export \"table\") ",($i)," funcref)";""]
 tabl:watbl y t
 mems:"(memory (export \"memory\") 1)"
 "\n"/(""~/)_("(module",sysi;mems;cons;vars;data;fncs;tabd;tabl;")\n")}
 
wadat:{,/@["\\",/256^`x@_!256;0+wac;wac:_("a"+!26),("A"+!26),("0"+!10),"~!@#$%^&*()_+{}[]|;:<>,./?"]x}
wamem:{y;""}
wacon:{"(global $",($S x)," ",(watype S 1+x)," (",(*y),"))"}
wavar:{"(global $",($S x)," (mut ",(watype S 1+x),") (",(*y),"))"}
watab:{y;(I x;S x)}
watbl:{"\n"/{"(elem (i32.const ",($x),") func ",y,")"}'[(*x)i;" "/'"$",''$(i:&~1=`d@*x)^(x:+x)1]}
wafun:{s:$S x
 args:" "/y@&`arg=T c:&x=P
 rets:" "/y@&`res=T c
 locs:" "/y@&`loc=T c
 body:    y@*&`ast=T c
 expo:$[I x;"(export \"",s,"\")";""]
 (" "/(""~/)_("(func";"$",s;expo;args;rets;locs)),"\n",wafmt[body],")"}
waarg:{" "/(""~/)_("(param";(y 0);(watype S x),")")}
wasym:{y;"$",$S x}
wares:{y;"(result ",(watype S x),")"}
waloc:{"(local ",(*y)," ",(watype S x),")"}
wastm:{"\n"/y}
waast:wastm
walod:{(y 0),"\n",$[`b~s:S x;"i32.load8_s";(watype s),".load"]}
wasto:{"\n"/y,,$[`b~s:S x;"i32.store8";(watype s),".store"]}
waant:{(y,'"\n"),t,".const -1\n",t,".xor\n",(t:watype S x),".and"}
waasn:{ y,("\nlocal.set $";"\nglobal.set $")[I x],$S x }

wacal:{$[`panic~s:S x;"unreachable";("\n"/y),"\n",$[#c:wainst s;c;"call $",$s]]}

wacli:{("\n"/((I x)#1_y),,*y),"\n","call_indirect "," "/((1+I x)_y),$[S x;,"(result ",(watype S x),")";()]}
wadfr:{y;""}
wadrp:{"\n"/y,,"drop"}
wacst:{t:!watype;(y 1),"\n",wacast i:(t?y 0)+5*t?S x}
wafor:{$[I x;waslp[x;y];walop[x;y]]}
waslp:{"\n"/("loop";-5_"\n"/2_y;y 1;*y;"br_if 0\nend")}
walop:{bl:("block";"loop"),'$[`~S x;2^"";(l,"1";(l:" $",$S x),"0")]
 "\n"/bl,($[#*y;(*y),"\ni32.eqz\nbr_if 1";""]; "\n"/2_y;y 1;"br 0\nend\nend")}
waget:{y;"local.get $",$S x}
waGet:{y;"global.get $",$S x}
wajmp:{y;"br ",$[`~S x;$1+I x;"$",($S x),$I x]}

waflts:(,0x182d4454fb210940)!,"3.141592653589793" /predefined constants that don't roundtrip
waflt:{x:C x+!8;$[#r:waflts x;r;$*`f x]}
walit:{y;i:I x;(watype t),".const ",$[`f~t:S x;waflt i
                 `i~t;$[0N~i;"-2147483648";$i]
                 `u~t;$[i<0;"0x",`x@|`c@,i;$i]
		 $[i~0;,"0";(-1~h:*`i C 4+i+!4)&0>j:*`i C i+!4; $j; (h~0)&j>0; $j; hx8 i]]}

waneg:{$[`f~t:S x;(*y),"\nf64.neg";"\n"/(t,".const 0";(*y);(t:watype t),".sub")]}
wanot:{(*y),"\ni32.eqz"}
wanop:{y;""}
waret:{("\n"/y),"\nreturn"}
waswc:{
 r:(,"block")@&n:#$[I x;y;y,:,""]
 r[0],:$[S x;" (result ",(watype S x),")";""]
 r,:,*y
 r,:,"br_table ",(" "/$!n-1),"\nend"
 r,:(1_y),'"\nbr ",/($|!n-1),\"\nend"
 ("\n"/r)}
wacnd:{"\n"/(*y;$[#t:watype S x;"if (result ",t,")";"if"];y 1;$[3~#y;"else\n",(y 2);""],"\nend")}
watyp:{y;S x} /keep as symbol

wa2: {"\n"/(z 0;z 1;(watype S y),x)}
wa2s:{"\n"/(z 0;z 1;(watype t  ),x,wasign@t:S y)}
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
wainst,:`Memorysize2`Memorygrow2`Memorycopy2`Memorycopy3!((,"unreachable\ni32.const 0")0 0),(,"drop\ndrop\ndrop\nunreachable")0 0
wainst,:{x!{(_32+*x),(1_3#x),".",3_x}'$x}`I32clz`F64abs`F64sqrt`F64floor`F64copysign`F64min`F64max`"F64reinterpret_i64"`"I64reinterpret_f64"

/         i                   u                  j                   k                   f
wacast:(""                 ;""                 ;"i32.wrap_i64"     ;"i32.wrap_i64"     ;"i32.trunc_f64_s"  /i
        ""                 ;""                 ;"i32.wrap_i64"     ;"i32.wrap_i64"     ;"i32.trunc_f64_u"  /u
        "i64.extend_i32_s" ;""                 ;""                 ;""                 ;"i64.trunc_f64_s"  /j
        "i64.extend_i32_u" ;"i64.extend_i32_u" ;""                 ;""                 ;"i64.trunc_f64_u"  /k
        "f64.convert_i32_s";"f64.convert_i32_u";"f64.convert_i64_s";"f64.convert_i64_u";"")                /f

wasys:"
(import \"env\" \"Exit\"  (func $Exit  (param i32)))
(import \"env\" \"Args\"  (func $Args  (result i32)))
(import \"env\" \"Arg\"   (func $Arg   (param i32) (param i32) (result i32)))
(import \"env\" \"Read\"  (func $Read  (param i32) (param i32) (param i32) (result i32)))
(import \"env\" \"Write\" (func $Write (param i32) (param i32) (param i32) (param i32) (result i32)))
(import \"env\" \"ReadIn\" (func $ReadIn (param i32) (param i32) (result i32)))
(import \"env\" \"Native\" (func $Native (param i64) (param i64) (result i64)))"

wafmt:{x:"\n"\x;
 ("\n"/waind waopt x@&~(x~\"nop")|x~\"")}
waind:{ /indent: +(if block loop else) -(end else)
 /s:(#x)#1,1+\(e:"els"~/h)|("blo"~/h)|("loo"~/h)|"if "~/h:x@\0 1 2  
 s:(#x)#1,+\(,1),(e:"els"~/h)|("blo"~/h)|("loo"~/h)|"if "~/h:x@\0 1 2  
 s-:+\e|"end"~/h;(s#\" "),'x}
waopt:{ /tee, !lt_s, last return
 i:&"local.get"~/9#'x
 i:-1+i@&(x i-1)~'("local.set"),/t:9_'x i
 x:@[x;i;"local.tee",/9_'x i]
 x@:(1+i)_!#x

 i:&("i32.eqz"~/x 1+!#x)&"i32.lt_s"~/x
 x:@[x;i;(,"i32.ge_s")@&#i]
 x@:(1+i)_!#x
 
 i:&("if"~/x 2+!#x)&("i32.ne"~/x 1+!#x)&"i32.const 0"~/x
 x@:(i,1+i)_!#x

 (-"return"~*|x)_x}

/const.0 eq eqz br_if -> br_if
/ge_s eqz -> lt_s 

