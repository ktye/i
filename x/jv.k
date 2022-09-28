
/jvm   https://medium.com/@davethomas_9528/writing-hello-world-in-java-byte-code-34f75428e0ad
F:`!()
hx8:{[i]"0x",`x@|C i+!8}
lp:{(~~)#@[&#x;x i;i:!#x]}

h:{_(256\x;256/x)} /todo _256\\x (issue26)



jvm:{
 con:(0x070002;0x010004,"main"                   /class
      0x070004;0x010010,"java/lang/Object")      /superclass
 /           min.maj.  #con      constant pool
 r:0xcafebabe0000003f,(h 1+#con),,/con
 /   acc.idx.sup.ifc.fld.met.atr.
 r,0x0021000100000000000000000000 }

`"k.class"<jvm[]  /c/local/jdk/bin/javap.exe k.class



jv:{[o]opts::o
 
 F[`prg`mem`con`var`tab`fun`arg`sym`res`loc`ast`lod`sto]:(jvprg;jvmem;jvcon;jvvar;jvtab;jvfun;jvarg;jvsym;jvres;jvloc;jvast;jvlod;jvsto)
 F[`add`and`bnd`ant`asn`cal`cli`cnd`dfr`div`drp`eql`cst]:(jvadd;jvand;jvbnd;jvant;jvasn;jvcal;jvcli;jvcnd;jvdfr;jvdiv;jvdrp;jveql;jvcst)
 F[`for`get`Get`gte`jmp`les`lit`lte`mod`mor`mul`neg    ]:(jvfor;jvget;jvGet;jvgte;jvjmp;jvles;jvlit;jvlte;jvmod;jvmor;jvmul;jvneg)
 F[`neq`nop`not`orr`bor`xor`ret`shl`shr`stm`sub`swc`typ]:(jvneq;jvnop;jvnot;jvorr;jvbor;jvxor;jvret;jvshl;jvshr;jvstm;jvsub;jvswc;jvtyp)
 
 S[0]:$[0N~i:o?`lib;S 0;o 1+i]
 n:(&#P){x[y]+:1}/P;n[0]-:1;
 s:{c:(0,(#x)-z)^x
   (c 0),,,/(F T y)[y;|c 1]}/[();|!#P;|n]
*s}

jvprg:{y;
 c:1_&P=0
 excl:!0
 fncs:"\n"/:y@&`fun=T c
 cons:"\n"/:y@&`con=T c
 vars:"\n"/:y@&`var=T c
 sysi:$[`nort':opts;"";jvsys]
 data:$[#D;"(data (i32.const 0) \"",(jvdat D),"\")";""]
 tabd:$[0<i:1+|/I c t:&`tab=T c;"(table (export \"table\") ",($i)," funcref)";""]
 tabl:jvtbl y t
 mems:"(memory (export \"memory\") 1)"
 "\n"/:(""~/:)_("(module",sysi;mems;cons;vars;data;fncs;tabd;tabl;")\n")}
 
jvdat:{,/@["\\",/:256^`x@_!256;0+jvc;jvc:_("a"+!26),("A"+!26),("0"+!10),"~!@#$%^&*()_+{}[]|;:<>,./?"]x}
jvmem:{y;""}
jvcon:{"(global $",($S x)," ",(jvtype S 1+x)," (",(*y),"))"}
jvvar:{"(global $",($S x)," (mut ",(jvtype S 1+x),") (",(*y),"))"}
jvtab:{y;(I x;S x)}
jvtbl:{"\n"/:{"(elem (i32.const ",($x),") func ",y,")"}'[(*x)i;" "/:'"$",''$(i:&~1=-':*x)^(x:+x)1]}
jvfun:{s:$S x
 args:" "/:y@&`arg=T c:&x=P
 rets:" "/:y@&`res=T c
 locs:" "/:y@&`loc=T c
 body:    y@*&`ast=T c
 expo:$[I x;"(export \"",s,"\")";""]
 (" "/:(""~/:)_("(func";"$",s;expo;args;rets;locs)),"\n",jvfmt[body],")"}
jvarg:{" "/:(""~/:)_("(param";(y 0);(jvtype S x),")")}
jvsym:{y;"$",$S x}
jvres:{y;"(result ",(jvtype S x),")"}
jvloc:{"(local ",(*y)," ",(jvtype S x),")"}
jvstm:{"\n"/:y}
jvast:jvstm
jvlod:{(y 0),"\n",$[`b~s:S x;"i32.load8_s";(jvtype s),".load"]}
jvsto:{"\n"/:y,,$[`b~s:S x;"i32.store8";(jvtype s),".store"]}
jvant:{(y,'"\n"),t,".const -1\n",t,".xor\n",(t:jvtype S x),".and"}
jvasn:{y[&b:"$_"~/:-1_y]:,"";(*|y),,/|(("\nlocal.set ";"\nglobal.set ";"\ndrop")(2*b)|I 1+x+!I x),'-1_y}

jvcal:{$[`panic~s:S x;"unreachable";("\n"/:y),"\n",$[#c:jvinst s;c;"call $",$s]]}

jvcli:{("\n"/:((I x)#1_y),,*y),"\n","call_indirect "," "/:(1+I x)_y}
jvdfr:{y;""}
jvdrp:{"\n"/:y,,"drop"}
jvcst:{t:!jvtype;(y 1),"\n",jvcast i:(t?y 0)+5*t?S x}
jvfor:{$[I x;jvslp[x;y];jvlop[x;y]]}
jvslp:{"\n"/:("loop";"\n"/:2_y;y 1;*y;"br_if 0\nend")}
jvlop:{bl:("block";"loop"),'$[`~S x;2^"";(l,"1";(l:" $",$S x),"0")]
 "\n"/:bl,($[#*y;(*y),"\ni32.eqz\nbr_if 1";""]; "\n"/:2_y;y 1;"br 0\nend\nend")}
jvget:{y;"local.get $",$S x}
jvGet:{y;"global.get $",$S x}
jvjmp:{y;"br ",$[`~S x;$1+I x;"$",($S x),$I x]}

jvflt:{(*&" "=c)#c:C 8+x+!24}
jvlit:{y;i:I x;(jvtype t),".const ",$[`f~t:S x;jvflt i
                 `i~t;$[0N~i;"-2147483648";$i]
                 `u~t;$[i<0;"0x",`x@|`c@,i;$i]
		 $[i~0;,"0";(-1~h:*`i C 4+i+!4)&0>j:*`i C i+!4; $j; (h~0)&j>0; $j; hx8 i]]}

jvneg:{$[`f~t:S x;(*y),"\nf64.neg";"\n"/:(t,".const 0";(*y);(t:jvtype t),".sub")]}
jvnot:{(*y),"\ni32.eqz"}
jvnop:{y;""}
jvret:{("\n"/:y),"\nreturn"}
jvswc:{
 r:(,"block")@&n:#$[I x;y;y,:,""]
 r[0],:$[S x;" (result ",(jvtype S x),")";""]
 r,:,*y
 r,:,"br_table ",(" "/:$!n-1),"\nend"
 r,:(1_y),'"\nbr ",/:($|!n-1),\:"\nend"
 ("\n"/:r)}
jvcnd:{"\n"/:(*y;$[#t:jvtype S x;"if (result ",t,")";"if"];y 1;$[3~#y;"else\n",(y 2);""],"\nend")}
jvtyp:{y;S x} /keep as symbol

jv2: {"\n"/:(z 0;z 1;(jvtype S y),x)}
jv2s:{"\n"/:(z 0;z 1;(jvtype t  ),x,jvsign@t:S y)}
jvadd:jv2".add";jvsub:jv2".sub"
jvmul:jv2".mul";jvdiv:jv2s".div";jvmod:jv2s".rem"
jvand:jv2".and";jvbnd:jv2".and"
jvneq:jv2".ne"
jvmor:jv2s".gt";jvgte:jv2s".ge";jvles:jv2s".lt";jvlte:jv2s".le";
jvorr:jv2".or";jvbor:jvorr;jvxor:jv2".xor"
jveql:jv2".eq"
jvshl:jv2".shl";jvshr:jv2s".shr"

jvtype:`i`u`j`k`f!("i32";"i32";"i64";"i64";"f64")
jvsign:`i`u`j`k`f!("_s";  "_u"; "_s"; "_u";   "")
jvinst:`I32B`Memorysize`Memorygrow`Memorycopy`Memoryfill!("nop";"memory.size";"memory.grow";"memory.copy";"memory.fill")
jvinst,:`Memorysize2`Memorygrow2`Memorycopy2`Memorycopy3!((,"unreachable\ni32.const 0")0 0),(,"drop\ndrop\ndrop\nunreachable")0 0
jvinst,:{x!{(_32+*x),(1_3#x),".",3_x}'$x}`I32clz`F64abs`F64sqrt`F64floor`F64copysign`F64min`F64max`"F64reinterpret_i64"`"I64reinterpret_f64"

/         i                   u                  j                   k                   f
jvcast:(""                 ;""                 ;"i32.wrap_i64"     ;"i32.wrap_i64"     ;"i32.trunc_f64_s"  /i
        ""                 ;""                 ;""                 ;"i32.wrap_i64"     ;""                 /u
        "i64.extend_i32_s" ;""                 ;""                 ;""                 ;"i64.trunc_f64_s"  /j
        "i64.extend_i32_u" ;"i64.extend_i32_u" ;""                 ;""                 ;"i64.trunc_f64_u"  /k
        "f64.convert_i32_s";"f64.convert_i32_u";"f64.convert_i64_s";"f64.convert_i64_u";"")                /f

jvsys:"
(import \"env\" \"Exit\"  (func $Exit  (param i32)))
(import \"env\" \"Args\"  (func $Args  (result i32)))
(import \"env\" \"Arg\"   (func $Arg   (param i32) (param i32) (result i32)))
(import \"env\" \"Read\"  (func $Read  (param i32) (param i32) (param i32) (result i32)))
(import \"env\" \"Write\" (func $Write (param i32) (param i32) (param i32) (param i32) (result i32)))
(import \"env\" \"ReadIn\" (func $ReadIn (param i32) (param i32) (result i32)))
(import \"env\" \"Native\" (func $Native (param i64) (param i64) (result i64)))"

jvfmt:{x:"\n"\:x;
 ("\n"/:jvind jvopt x@&~(x~\:"nop")|x~\:"")}
jvind:{ /indent: +(if block loop else) -(end else)
 s:(#x)#1,1+\(e:"els"~/:h)|("blo"~/:h)|("loo"~/:h)|"if "~/:h:3#'x   
 s-:+\e|"end"~/:h;(s#\:" "),'x}
jvopt:{ /tee, !lt_s, last return
 i:&"local.get"~/:9#'x
 i:-1+i@&(x i-1)~'("local.set"),/:t:9_'x i
 x:@[x;i;"local.tee",/:9_'x i]
 x@:(1+i)_!#x

 i:&("i32.eqz"~/:x 1+!#x)&"i32.lt_s"~/:x
 x:@[x;i;(,"i32.ge_s")@&#i]
 x@:(1+i)_!#x
 
 i:&("if"~/:x 2+!#x)&("i32.ne"~/:x 1+!#x)&"i32.const 0"~/:x
 x@:(i,1+i)_!#x

 (-"return"~*|x)_x}


