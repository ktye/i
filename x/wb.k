/webassembly compiler (text format)
F:`!()
hx8:{[i]"0x",`x@|C i+!8}
lp:{(~~)#@[&#x;x i;i:!#x]}

wb:{[o]opts::o
 F[`prg`mem`con`var`tab`fun`arg`sym`res`loc`ast`lod`sto]:(wbprg;wbmem;wbcon;wbvar;wbtab;wbfun;wbarg;wbsym;wbres;wbloc;wbast;wblod;wbsto)
 F[`add`and`bnd`ant`asn`cal`cli`cnd`dfr`div`drp`eql`cst]:(wbadd;wband;wbbnd;wbant;wbasn;wbcal;wbcli;wbcnd;wbdfr;wbdiv;wbdrp;wbeql;wbcst)
 F[`for`get`Get`gte`jmp`les`lit`lte`mod`mor`mul`neg    ]:(wbfor;wbget;wbGet;wbgte;wbjmp;wbles;wblit;wblte;wbmod;wbmor;wbmul;wbneg)
 F[`neq`nop`not`orr`bor`xor`ret`shl`shr`stm`sub`swc`typ]:(wbneq;wbnop;wbnot;wborr;wbbor;wbxor;wbret;wbshl;wbshr;wbstm;wbsub;wbswc;wbtyp)
 
 S[0]:$[0N~i:o?`lib;S 0;o 1+i]
 n:(&#P){x[y]+:1}/P;n[0]-:1;
 s:{c:(0,(#x)-z)^x
   (c 0),,,/(F T y)[y;|c 1]}/[();|!#P;|n]
*s}


wbprg:{
 c:1_&P=0
 fncs:"\n"/:y@&`fun=T c
 cons:"\n"/:y@&`con=T c
 vars:"\n"/:y@&`var=T c
 sysi:$[`nort':opts;"";wbsys]
 data:$[#D;"(data (i32.const 0) \"",(wbdat D),"\")";""]
 tabd:$[0<i:1+|/I c t:&`tab=T c;"(table (export \"table\") ",($i)," funcref)";""]
 tabl:wbtbl y t
 mems:"(memory (export \"memory\") 1)"
 "\n"/:(""~/:)_("(module",sysi;mems;cons;vars;data;fncs;tabd;tabl;")\n")}
 
wbdat:{,/@["\\",/:256^`x@_!256;0+wbc;wbc:_("a"+!26),("A"+!26),("0"+!10),"~!@#$%^&*()_+{}[]|;:<>,./?"]x}
wbmem:{y;""}
wbcon:{"(global $",($S x)," ",(wbtype S 1+x)," (",(*y),"))"}
wbvar:{"(global $",($S x)," (mut ",(wbtype S 1+x),") (",(*y),"))"}
wbtab:{y;(I x;S x)}
wbtbl:{"\n"/:{"(elem (i32.const ",($x),") func ",y,")"}'[(*x)i;" "/:'"$",''$(i:&~1=-':*x)^(x:+x)1]}
wbfun:{s:$S x
 args:" "/:y@&`arg=T c:&x=P
 rets:" "/:y@&`res=T c
 locs:" "/:y@&`loc=T c
 body:    y@*&`ast=T c
 expo:$[I x;"(export \"",s,"\")";""]
 (" "/:(""~/:)_("(func";"$",s;expo;args;rets;locs)),"\n",wbfmt[body],")"}
wbarg:{" "/:(""~/:)_("(param";(y 0);(wbtype S x),")")}
wbsym:{y;"$",$S x}
wbres:{y;"(result ",(wbtype S x),")"}
wbloc:{"(local ",(*y)," ",(wbtype S x),")"}
wbstm:{"\n"/:y}
wbast:wbstm
wblod:{(y 0),"\n",$[`b~s:S x;"i32.load8_s";(wbtype s),".load"]}
wbsto:{"\n"/:y,,$[`b~s:S x;"i32.store8";(wbtype s),".store"]}
wbant:{(y,'"\n"),t,".const -1\n",t,".xor\n",(t:wbtype S x),".and"}
wbasn:{ y,("\nlocal.set $";"\nglobal.set $")[I x],$S x }

wbcal:{$[`panic~s:S x;"unreachable";("\n"/:y),"\n",$[#c:wbinst s;c;"call $",$s]]}

wbcli:{("\n"/:((I x)#1_y),,*y),"\n","call_indirect "," "/:((1+I x)_y),$[S x;,"(result ",(wbtype S x),")";()]}
wbdfr:{y;""}
wbdrp:{"\n"/:y,,"drop"}
wbcst:{t:!wbtype;(y 1),"\n",wbcast i:(t?y 0)+5*t?S x}
wbfor:{$[I x;wbslp[x;y];wblop[x;y]]}
wbslp:{"\n"/:("loop";-5_"\n"/:2_y;y 1;*y;"br_if 0\nend")}
wblop:{bl:("block";"loop"),'$[`~S x;2^"";(l,"1";(l:" $",$S x),"0")]
 "\n"/:bl,($[#*y;(*y),"\ni32.eqz\nbr_if 1";""]; "\n"/:2_y;y 1;"br 0\nend\nend")}
wbget:{y;"local.get $",$S x}
wbGet:{y;"global.get $",$S x}
wbjmp:{y;"br ",$[`~S x;$1+I x;"$",($S x),$I x]}

wbflts:(,0x182d4454fb210940)!,"3.141592653589793" /predefined constants that don't roundtrip
wbflt:{x:C x+!8;$[#r:wbflts x;r;$*`f x]}
wblit:{y;i:I x;(wbtype t),".const ",$[`f~t:S x;wbflt i
                 `i~t;$[0N~i;"-2147483648";$i]
                 `u~t;$[i<0;"0x",`x@|`c@,i;$i]
		 $[i~0;,"0";(-1~h:*`i C 4+i+!4)&0>j:*`i C i+!4; $j; (h~0)&j>0; $j; hx8 i]]}

wbneg:{$[`f~t:S x;(*y),"\nf64.neg";"\n"/:(t,".const 0";(*y);(t:wbtype t),".sub")]}
wbnot:{(*y),"\ni32.eqz"}
wbnop:{y;""}
wbret:{("\n"/:y),"\nreturn"}
wbswc:{
 r:(,"block")@&n:#$[I x;y;y,:,""]
 r[0],:$[S x;" (result ",(wbtype S x),")";""]
 r,:,*y
 r,:,"br_table ",(" "/:$!n-1),"\nend"
 r,:(1_y),'"\nbr ",/:($|!n-1),\:"\nend"
 ("\n"/:r)}
wbcnd:{"\n"/:(*y;$[#t:wbtype S x;"if (result ",t,")";"if"];y 1;$[3~#y;"else\n",(y 2);""],"\nend")}
wbtyp:{y;S x} /keep as symbol

wb2: { (,/z),x(`i`u`j`k`f?S y)}
wb2s:{"\n"/:(z 0;z 1;(wbtype t  ),x,wbsign@t:S y)}
wbadd:wb2 0x6a6a7c7ca0;wbsub:wb2 0x6b6b7d7da1
wbmul:wb2 0x6c6c7e7ea2;wbdiv:wb2 0x6d6e7f80a3;wbmod:wb2s".rem"
wband:wbbnd:wb2 0x7171838300;wborr:wbbor:wb2 0x7272848400;wbxor:wb2 0x7373858500
wbneq:wb2 0x4747525262;wbeql:wb2 0x4646515161
wbmor:wb2 0x4a4b555664;wbgte:wb2 0x4e4f595a66;wbles:wb2 0x4849535463; wblte:wb2 0x4c4d575856
wbshl:wb2 0x7474868600;wbshr:wb2 0x7576878800

wbtype:`i`u`j`k`f!("i32";"i32";"i64";"i64";"f64")
wbsign:`i`u`j`k`f!("_s";  "_u"; "_s"; "_u";   "")
wbinst:`I32B`Memorysize`Memorygrow`Memorycopy`Memoryfill!("nop";"memory.size";"memory.grow";"memory.copy";"memory.fill")
wbinst,:`Memorysize2`Memorygrow2`Memorycopy2`Memorycopy3!((,"unreachable\ni32.const 0")0 0),(,"drop\ndrop\ndrop\nunreachable")0 0
wbinst,:{x!{(_32+*x),(1_3#x),".",3_x}'$x}`I32clz`F64abs`F64sqrt`F64floor`F64copysign`F64min`F64max`"F64reinterpret_i64"`"I64reinterpret_f64"

/         i                   u                  j                   k                   f
wbcast:(""                 ;""                 ;"i32.wrap_i64"     ;"i32.wrap_i64"     ;"i32.trunc_f64_s"  /i
        ""                 ;""                 ;""                 ;"i32.wrap_i64"     ;""                 /u
        "i64.extend_i32_s" ;""                 ;""                 ;""                 ;"i64.trunc_f64_s"  /j
        "i64.extend_i32_u" ;"i64.extend_i32_u" ;""                 ;""                 ;"i64.trunc_f64_u"  /k
        "f64.convert_i32_s";"f64.convert_i32_u";"f64.convert_i64_s";"f64.convert_i64_u";"")                /f

wbsys:"
(import \"env\" \"Exit\"  (func $Exit  (param i32)))
(import \"env\" \"Args\"  (func $Args  (result i32)))
(import \"env\" \"Arg\"   (func $Arg   (param i32) (param i32) (result i32)))
(import \"env\" \"Read\"  (func $Read  (param i32) (param i32) (param i32) (result i32)))
(import \"env\" \"Write\" (func $Write (param i32) (param i32) (param i32) (param i32) (result i32)))
(import \"env\" \"ReadIn\" (func $ReadIn (param i32) (param i32) (result i32)))
(import \"env\" \"Native\" (func $Native (param i64) (param i64) (result i64)))"

wbfmt:{x:"\n"\:x;
 ("\n"/:wbind wbopt x@&~(x~\:"nop")|x~\:"")}
wbind:{ /indent: +(if block loop else) -(end else)
 s:(#x)#1,1+\(e:"els"~/:h)|("blo"~/:h)|("loo"~/:h)|"if "~/:h:x@\:0 1 2  
 s-:+\e|"end"~/:h;(s#\:" "),'x}
wbopt:{ /tee, !lt_s, last return
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

/const.0 eq eqz br_if -> br_if
/ge_s eqz -> lt_s 

