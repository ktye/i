/webassembly compiler (text format)
F:`!()
hx8:{[i]"0x",`x@|C i+!8}
lp:{(~~)#@[&#x;x i;i:!#x]}
ifun:&#T                                       /function index per node
locs:0                                         /ifun!names
glob:0#`
fsig:`!()
sigs:()

wb:{[o]opts::o
 F[`prg`mem`con`var`tab`fun`arg`sym`res`loc`ast`lod`sto]:(wbprg;wbmem;wbcon;wbvar;wbtab;wbfun;wbarg;wbsym;wbres;wbloc;wbast;wblod;wbsto)
 F[`add`and`bnd`ant`asn`cal`cli`cnd`dfr`div`drp`eql`cst]:(wbadd;wband;wbbnd;wbant;wbasn;wbcal;wbcli;wbcnd;wbdfr;wbdiv;wbdrp;wbeql;wbcst)
 F[`for`get`Get`gte`jmp`les`lit`lte`mod`mor`mul`neg    ]:(wbfor;wbget;wbGet;wbgte;wbjmp;wbles;wblit;wblte;wbmod;wbmor;wbmul;wbneg)
 F[`neq`nop`not`orr`bor`xor`ret`shl`shr`stm`sub`swc`typ]:(wbneq;wbnop;wbnot;wborr;wbbor;wbxor;wbret;wbshl;wbshr;wbstm;wbsub;wbswc;wbtyp)
 
 ifun::-1+\fi:T=`fun                              /function index for each node
 locs::(!g)!(S 1+i)@.g:=ifun i:&(T=`arg)|T=`loc
 g:=f:S P i@:&`fun=T P i:&(T=`arg)|T=`res
 fsig::f!wbtype(S i)g f:?f
 
 fsig::fsig,fx!(#fx:fx@&0N=f?fx:S@&fi)#,,0x00    /function with no arg/res
 sigs:: ?.fsig
 
 
 S[0]:$[0N~i:o?`lib;S 0;o 1+i]
 n:(&#P){x[y]+:1}/P;n[0]-:1;
 s:{c:(0,(#x)-z)^x
   (c 0),,,/(F T y)[y;|c 1]}/[();|!#P;|n]
*s}


wbprg:{
 s:{(_x),(leb@#y),y} /0(custom) 1(type) 2(import) 3(func) 4(table) 5(mem) 6(global) 7(export) 8(start) 9(element) 10(code) 11(data) 12(data count)
 c:1_&P=0;f:y fn:&`fun=T c;
 fncs:$[#f;s[10;(leb@#f),/f];""]
 fnsg:$[#f;s[3;(leb@#f),/leb'sigs?fsig S c fn];""]
 glob::     S c@&(`con=T c)
 glob::glob,S c@& `var=T c
 cons:,/y@&`con=T c
 vars:,/y@&`var=T c
 glbs:s[6;(leb@#glob),,/cons,vars]
 sysi:$[`nort':opts;"";wbsys]
 data:$[#D;"(data (i32.const 0) \"",(wbdat D),"\")";""]
 tabd:$[0<i:1+|/I c t:&`tab=T c;"(table (export \"table\") ",($i)," funcref)";""]
 tabl:wbtbl y t
 mems:"(memory (export \"memory\") 1)"
 /"\n"/:(""~/:)_("(module",sysi;mems;cons;vars;data;fncs;tabd;tabl;")\n")} 
 typs:s[1;(leb@#t),/t:wbsig'.fsig]
 
 0x0061736d01000000,typs,fnsg,glbs,fncs}

wbsig:{0x60,((leb@#a),/a:-1_x),$[0x00~r:*|x;_0;(_1),r]} //wbsig args, ret


wbdat:{,/@["\\",/:256^`x@_!256;0+wbc;wbc:_("a"+!26),("A"+!26),("0"+!10),"~!@#$%^&*()_+{}[]|;:<>,./?"]x}
wbmem:{y;""}
wbcon:{(wbtype S 1+x),0x00,(*y),0x0b}
wbvar:{(wbtype S 1+x),0x01,(*y),0x0b}
wbtab:{y;(I x;S x)}
wbtbl:{"\n"/:{"(elem (i32.const ",($x),") func ",y,")"}'[(*x)i;" "/:'"$",''$(i:&~1=-':*x)^(x:+x)1]}

wbfun:{s:$S x;c:&x=P
 locs:(leb@#l),,/l:{(leb@#x),wbtype@*x}'(&~`x~':l)^l:S c@&`loc=T c
 body:    y@*&`ast=T c
 expo:$[I x;"(export \"",s,"\")";""]
 /expo;args;rets;locs
 (leb@#r),r:locs,((- 0x0f~*|body)_body),0x0b}
 
wbarg:{y;wbtype S x}
wbsym:{y;S x}
wbres:{y;"(result ",(wbtype S x),")"}
wbloc:{"(local ",(*y)," ",(wbtype S x),")"}
wbstm:{,/y}
wbast:wbstm
wblod:{(y 0),"\n",$[`b~s:S x;"i32.load8_s";(wbtype s),".load"]}
wbsto:{"\n"/:y,,$[`b~s:S x;"i32.store8";(wbtype s),".store"]}
wbant:{(y,'"\n"),t,".const -1\n",t,".xor\n",(t:wbtype S x),".and"}
wbasn:{y,$[I x;0x24,leb(glob?S x);0x21,leb(locs ifun x)?S x]}

wbcal:{$[`panic~s:S x;_!0;(,/y),$[#c:wbinst s;c;0x10,leb@(!fsig)?s]]}

wbcli:{(,/(I x)#1_y),(*y),0x11,(leb sigs?(,/(1+I x)_y),wbtype S x),0x00 } 
wbdfr:{y;""}
wbdrp:{y,0x1a}
wbcst:{(y 1),$[0x00~c:wbcast[t?S x;(t:!wbtype)?S 1+x];"";c]}
wbfor:{$[I x;wbslp[x;y];wblop[x;y]]}
wbslp:{"\n"/:("loop";-5_"\n"/:2_y;y 1;*y;"br_if 0\nend")}

/wblop:{bl:("block";"loop"),'$[`~S x;2^"";(l,"1";(l:" $",$S x),"0")]
/ "\n"/:bl,($[#*y;(*y),"\ni32.eqz\nbr_if 1";""]; "\n"/:2_y;y 1;"br 0\nend\nend")}
wblop:{ 0x02400340,(*y),$[#*y;0x450d01;""],(y 2),(y 1),0x0c000b0b }
 
wbget:{y;0x20,leb(locs ifun x)?S x}
wbGet:{y;"global.get $",$S x}
wbjmp:{y;0x0c,(_1+I x)} /todo:label  /{"br ",$[`~S x;$1+I x;"$",($S x),$I x]}

wbflts:(,0x182d4454fb210940)!,"3.141592653589793" /predefined constants that don't roundtrip
wbflt:{x:C x+!8;$[#r:wbflts x;r;$*`f x]}
wblit:{y;i:I x;((`i`u`k`j`f!0x4141424244)t),,$[`f~t:S x;C i+!8;`i~t;leb i;`u~t;leb i;leb8 C i+!8]}

wbneg:{wbsub[x;($[`f~t:S x;0x44,(8#0x00);(0x4141424244@`i`u`j`k?S x),0x00];*y)]}
wbnot:{(*y),"\ni32.eqz"}
wbnop:{y;""}
wbret:{(*y),0x0f}
wbswc:{
 r:(,"block")@&n:#$[I x;y;y,:,""]
 r[0],:$[S x;" (result ",(wbtype S x),")";""]
 r,:,*y
 r,:,"br_table ",(" "/:$!n-1),"\nend"
 r,:(1_y),'"\nbr ",/:($|!n-1),\:"\nend"
 ("\n"/:r)}

/wbcnd:{"\n"/:(*y;$[#t:wbtype S x;"if (result ",t,")";"if"];y 1;$[3~#y;"else\n",(y 2);""],"\nend")}
wbcnd:{(*y),0x04,$[`~t:S x;0x40;wbtype t],(y 1),$[3~#y;0x05,(y 2);""],0x0b}


wbtyp:{y;wbtype S x} /keep as symbol

wb2: {(,/z),x(`i`u`j`k`f?S y)}
wb2s:{"\n"/:(z 0;z 1;(wbtype t  ),x,wbsign@t:S y)}
wbadd:wb2 0x6a6a7c7ca0;wbsub:wb2 0x6b6b7d7da1
wbmul:wb2 0x6c6c7e7ea2;wbdiv:wb2 0x6d6e7f80a3;wbmod:wb2s".rem"
wband:wbbnd:wb2 0x7171838300;wborr:wbbor:wb2 0x7272848400;wbxor:wb2 0x7373858500
wbneq:wb2 0x4747525262;wbeql:wb2 0x4646515161
wbmor:wb2 0x4a4b555664;wbgte:wb2 0x4e4f595a66;wbles:wb2 0x4849535463; wblte:wb2 0x4c4d575856
wbshl:wb2 0x7474868600;wbshr:wb2 0x7576878800

wbtype:`i`u`j`k`f`!0x7f7f7e7e7c00
wbsign:`i`u`j`k`f!("_s";  "_u"; "_s"; "_u";   "")
wbinst:`I32B`Memorysize`Memorygrow`Memorycopy`Memoryfill!("nop";"memory.size";"memory.grow";"memory.copy";"memory.fill")
wbinst,:`Memorysize2`Memorygrow2`Memorycopy2`Memorycopy3!((,"unreachable\ni32.const 0")0 0),(,"drop\ndrop\ndrop\nunreachable")0 0
wbinst,:{x!{(_32+*x),(1_3#x),".",3_x}'$x}`I32clz`F64abs`F64sqrt`F64floor`F64copysign`F64min`F64max`"F64reinterpret_i64"`"I64reinterpret_f64"

wbcast:(0x0000a7a7aa;0x0000a7a7ab;0xacad0000b0;0xadad0000b1;0xb7b8b9ba00) /iujkf


leb:{leb8(`c@,x),4#$[x<0;0xff;0x00]}
leb8:{r:!0;B:1v:,/+(8#2)\\256/256+|x;while[B&(8<#v)&+/v;b:-7#v;s:*b;v:-7_v;r,:$[((0~&/v)|0x00~s)&(+/v)|s;1;B:0],b];_2//'-8^r}

wbsys:"
(import \"env\" \"Exit\"  (func $Exit  (param i32)))
(import \"env\" \"Args\"  (func $Args  (result i32)))
(import \"env\" \"Arg\"   (func $Arg   (param i32) (param i32) (result i32)))
(import \"env\" \"Read\"  (func $Read  (param i32) (param i32) (param i32) (result i32)))
(import \"env\" \"Write\" (func $Write (param i32) (param i32) (param i32) (param i32) (result i32)))
(import \"env\" \"ReadIn\" (func $ReadIn (param i32) (param i32) (result i32)))
(import \"env\" \"Native\" (func $Native (param i64) (param i64) (result i64)))"


/wbopt:{ /tee, !lt_s, last return
/ i:&"local.get"~/:9#'x
/ i:-1+i@&(x i-1)~'("local.set"),/:t:9_'x i
/ x:@[x;i;"local.tee",/:9_'x i]
/ x@:(1+i)_!#x
/
/ i:&("i32.eqz"~/:x 1+!#x)&"i32.lt_s"~/:x
/ x:@[x;i;(,"i32.ge_s")@&#i]
/ x@:(1+i)_!#x
/ 
/ i:&("if"~/:x 2+!#x)&("i32.ne"~/:x 1+!#x)&"i32.const 0"~/:x
/ x@:(i,1+i)_!#x
/
/ (-"return"~*|x)_x}

/const.0 eq eqz br_if -> br_if
/ge_s eqz -> lt_s 

