F:(0#`)!()
hx8:{[i]"0x",`x@|C i+!8}
delta:{x-(*x),-1_x}

prg:{
 c:1_&P=0
 excl:!0
 fncs:"\n"/y@&`fun=T c
 cons:"\n"/y@&`con=T c
 vars:"\n"/y@&`var=T c
 data:$[#D;"(data (i32.const 0) \"",(wadat D),"\")";""]
 tabd:$[0<i:1+|/I c t:&`tab=T c;"(table (export \"table\") ",($i)," funcref)";""]
 tabl:watbl y t
 mems:"(memory (export \"memory\") 1)"
 "\n"/(""~/)_("(module";mems;cons;vars;data;fncs;tabd;tabl;")\n")}
 
wadat:{,/@["\\",/256^`x@_!256;0+wac;wac:_("a"+!26),("A"+!26),("0"+!10),"~!@#$%^&*()_+{}[]|;:<>,./?"]x}
mem:{y;""}
con:{"(global $",($S x)," ",(type S 1+x)," (",(*y),"))"}
var:{"(global $",($S x)," (mut ",(type S 1+x),") (",(*y),"))"}
tab:{y;(I x;S x)}
watbl:{"\n"/{"(elem (i32.const ",($x),") func ",y,")"}'[(*x)i;" "/'"$",''$(i:&~1=delta@*x)^(x:+x)1]}
fun:{s:$S x
 args:" "/y@&`arg=T c:&x=P
 rets:" "/y@&`res=T c
 locs:" "/y@&`loc=T c
 body:    y@*&`ast=T c
 expo:$[I x;"(export \"",s,"\")";""]
 (" "/(""~/)_("(func";"$",s;expo;args;rets;locs)),"\n",body,")"}
arg:{" "/(""~/)_("(param";(y 0);(type S x),")")}
sym:{y;"$",$S x}
res:{y;"(result ",(type S x),")"}
loc:{"(local ",(*y)," ",(type S x),")"}
stm:{"\n"/y}
ast:stm
lod:{(y 0),"\n",$[`b~s:S x;"i32.load8_s";(type s),".load"]}
sto:{"\n"/y,,$[`b~s:S x;"i32.store8";(type s),".store"]}
ant:{(y,'"\n"),t,".const -1\n",t,".xor\n",(t:type S x),".and"}
asn:{ y,("\nlocal.set $";"\nglobal.set $")[I x],$S x }

cal:{$[`panic~s:S x;"unreachable";("\n"/y),"\n",$[#c:inst s;c;"call $",$s]]}

cli:{("\n"/((I x)#1_y),,*y),"\n","call_indirect "," "/((1+I x)_y),$[S x;,"(result ",(type S x),")";()]}
dfr:{y;""}
drp:{"\n"/y,,"drop"}
cst:{t:!type;(y 1),"\n",cast i:(t?y 0)+5*t?S x}
for:{$[I x;slp[x;y];lop[x;y]]}
slp:{"\n"/("loop";-5_"\n"/2_y;y 1;*y;"br_if 0\nend")}
lop:{bl:("block";"loop"),'$[`~S x;2^"";(l,"1";(l:" $",$S x),"0")]
 "\n"/bl,($[#*y;(*y),"\ni32.eqz\nbr_if 1";""]; "\n"/2_y;y 1;"br 0\nend\nend")}
get:{y;"local.get $",$S x}
Get:{y;"global.get $",$S x}
jmp:{y;"br ",$[`~S x;$1+I x;"$",($S x),$I x]}

flts:(,0x182d4454fb210940)!,"3.141592653589793" /predefined constants that don't roundtrip
flt:{x:C x+!8;$[#r:flts x;r;$*`f x]}
lit:{y;i:I x;(type t),".const ",$[`f~t:S x;flt i
                 `i~t;$i
                 `u~t;$[i<0;"0x",`x@|`c@,i;$i]
		 $[i~0;,"0";(-1~h:*`i C 4+i+!4)&0>j:*`i C i+!4; $j; (h~0)&j>0; $j; hx8 i]]}

neg:{$[`f~t:S x;(*y),"\nf64.neg";"\n"/(t,".const 0";(*y);(t:type t),".sub")]}
not:{(*y),"\ni32.eqz"}
nop:{y;""}
ret:{("\n"/y),"\nreturn"}
swc:{
 r:(,"block")@&n:#$[I x;y;y,:,""]
 r[0],:$[S x;" (result ",(type S x),")";""]
 r,:,*y
 r,:,"br_table ",(" "/$!n-1),"\nend"
 r,:(1_y),'"\nbr ",/($|!n-1),\"\nend"
 ("\n"/r)}
cnd:{"\n"/(*y;$[#t:type S x;"if (result ",t,")";"if"];y 1;$[3~#y;"else\n",(y 2);""],"\nend")}
typ:{y;S x} /keep as symbol

wa2: {"\n"/(z 0;z 1;(type S y),x)}
wa2s:{"\n"/(z 0;z 1;(type t  ),x,sign@t:S y)}
add:wa2".add";sub:wa2".sub"
mul:wa2".mul";div:wa2s".div";mod:wa2s".rem"
and:wa2".and";bnd:wa2".and"
neq:wa2".ne"
mor:wa2s".gt";gte:wa2s".ge";les:wa2s".lt";lte:wa2s".le";
orr:wa2".or";bor:orr;xor:wa2".xor"
eql:wa2".eq"
shl:wa2".shl";shr:wa2s".shr"

type:`i`u`j`k`f!("i32";"i32";"i64";"i64";"f64")
sign:`i`u`j`k`f!("_s";  "_u"; "_s"; "_u";   "")
inst:`I32B`Memorysize`Memorygrow`Memorycopy`Memoryfill!("nop";"memory.size";"memory.grow";"memory.copy";"memory.fill")
inst,:`Memorysize2`Memorygrow2`Memorycopy2`Memorycopy3!((,"unreachable\ni32.const 0")0 0),(,"drop\ndrop\ndrop\nunreachable")0 0
inst,:{x!{(_32+*x),(1_3#x),".",3_x}'$x}`I32clz`F64abs`F64sqrt`F64floor`F64copysign`F64min`F64max`"F64reinterpret_i64"`"I64reinterpret_f64"

/      i                   u                  j                   k                   f
cast:(""                 ;""                 ;"i32.wrap_i64"     ;"i32.wrap_i64"     ;"i32.trunc_f64_s"  /i
      ""                 ;""                 ;"i32.wrap_i64"     ;"i32.wrap_i64"     ;"i32.trunc_f64_u"  /u
      "i64.extend_i32_s" ;""                 ;""                 ;""                 ;"i64.trunc_f64_s"  /j
      "i64.extend_i32_u" ;"i64.extend_i32_u" ;""                 ;""                 ;"i64.trunc_f64_u"  /k
      "f64.convert_i32_s";"f64.convert_i32_u";"f64.convert_i64_s";"f64.convert_i64_u";"")                /f


F[`prg`mem`con`var`tab`fun`arg`sym`res`loc`ast`lod`sto]:(prg;mem;con;var;tab;fun;arg;sym;res;loc;ast;lod;sto)
F[`add`and`bnd`ant`asn`cal`cli`cnd`dfr`div`drp`eql`cst]:(add;and;bnd;ant;asn;cal;cli;cnd;dfr;div;drp;eql;cst)
F[`for`get`Get`gte`jmp`les`lit`lte`mod`mor`mul`neg    ]:(for;get;Get;gte;jmp;les;lit;lte;mod;mor;mul;neg)
F[`neq`nop`not`orr`bor`xor`ret`shl`shr`stm`sub`swc`typ]:(neq;nop;not;orr;bor;xor;ret;shl;shr;stm;sub;swc;typ)
 
n:{x[y]+:1}/(,&#P),P;n[0]-:1
s:{c:(0,(#x)-z)^x;(c 0),,,/(F T y)[y;|c 1]}/[();|!#P;|n]
`<*s
