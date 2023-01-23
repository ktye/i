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
 
 ifun::-1+\T=`fun                              /function index for each node
 locs::(!g)!(S 1+i)@.g:=ifun i:&(T=`arg)|T=`loc
 f:&T=`fun
 a:&T=`arg
 d:(S f)!(#f)#,0#`
 d[S@!g]:S a@.g:=f@f'a
 fsig::wbtype d,'@[(#f)#`;ifun r;S r:&T=`res]
 fsig::$[`nort':opts;fsig;fsig,wbsys]          /imports
 sigs::?.fsig
 glob::     S@&`con=T
 glob::glob,S@&`var=T
 
 S[0]:$[0N~i:o?`lib;S 0;o 1+i]
 n:(&#P){x[y]+:1}/P;n[0]-:1;
 s:{c:(0,(#x)-z)^x;(c 0),,,/(F T y)[y;|c 1]}/[();|!#P;|n]
*s}
 

wbprg:{
 s:{(_x),(leb@#y),y} /0(custom) 1(type) 2(import) 3(func) 4(table) 5(mem) 6(global) 7(export) 8(start) 9(element) 10(code) 11(data) 12(data count)
 c:1_&P=0;f:y fn:&`fun=T c;
 sysi:$[`nort':opts;"";s[2;(leb@#i),,/i:{0x03,"env",(leb@#x),x,0x00,(leb sigs?y)}'[$!wbsys;.wbsys]]]
 fncs:$[#f;s[10;(leb@#f),/f];""]
 fnsg:$[#f;s[3;(leb@#f),/leb'sigs?fsig S c fn];""]
 typs:s[1;(leb@#t),,/t:wbsig'sigs]
 cons:,/y@&`con=T c
 vars:,/y@&`var=T c
 mems:s[5;0x010001]
 glbs:s[6;(leb@#glob),,/cons,vars]
 tabl:$[#t:&`tab=T c;s[4;0x017000,(leb@|/I c t)];""]
 elem:$[#t;s[9;(leb@#t),,/t:wbtbl y t];""]
 data:$[#D;s[11;0x010041000b,(leb@#D),D];""]
 0x0061736d01000000,typs,sysi,fnsg,tabl,mems,glbs,elem,fncs,data}

wbsig:{0x60,((leb@#a),/a:-1_x),$[0x00~r:*|x;_0;(_1),r]} //wbsig args, ret
wbmem:{y;""}
wbcon:{(wbtype S 1+x),0x00,(*y),0x0b}
wbvar:{(wbtype S 1+x),0x01,(*y),0x0b}
wbtab:{y;(I x;(!fsig)?S x)}
wbtbl:{{0x0041,(leb x),0x0b,(leb@#y),,/leb'y}'[(*x)i;(i:&~1=-':*x)^(x:+x)1]}

wbfun:{s:$S x;c:&x=P
 locs:(leb@#l),,/l:{(leb@#x),wbtype@*x}'(&~`x~':l)^l:S c@&`loc=T c
 body:    y@*&`ast=T c
 /expo:$[I x;"(export \"",s,"\")";""]
 (leb@#r),r:locs,((- 0x0f~*|body)_body),0x0b}
 
wbarg:{y;wbtype S x}
wbsym:{y;S x}
wbres:{y;"(result ",(wbtype S x),")"}
wbloc:{"(local ",(*y)," ",(wbtype S x),")"}
wbstm:{,/y}
wbast:wbstm
wblod:{(y 0),(0x2d282829292b i),(0x000202030303 i:`b`i`u`j`k`f?S x),0x00}
wbsto:{(,/y),(0x3a3636373739 i),(0x000202030303 i:`b`i`u`j`k`f?S x),0x00}
wbant:{(,/y),(0x417f7371;0x427f8583)@(S x)':`j`k} /y0 y1 -1 xor and 
wbasn:{y,$[I x;0x24,leb glob?S x;0x21,leb(locs ifun x)?S x]}

/ $[#wbinst S x;0;$[0N~(!fsig)?S x;`null \S x;0]];
wbcal:{$[`panic~s:S x;,0x00;`I32B~s;*y;(,/y),$[#c:wbinst s;c;0x10,leb@(!fsig)?s]]}

wbcli:{(,/(I x)#1_y),(*y),0x11,(leb sigs?(,/(1+I x)_y),wbtype S x),0x00 } 
wbdfr:{y;""}
wbdrp:{y,0x1a}
wbcst:{(y 1),$[0x00~c:wbcast[t?S x;(t:!wbtype)?S 1+x];"";c]}
wbfor:{$[I x;wbslp[x;y];wblop[x;y]]}
wbslp:{0x0340,(-2_y 2),(y 1),(*y),0x0d000b }
wblop:{0x02400340,(*y),$[#*y;0x450d01;""],(y 2),(y 1),0x0c000b0b}
 
wbget:{y;0x20,leb(locs ifun x)?S x}
wbGet:{y;0x23,leb glob?S x}
wbjmp:{y;0x0c,(_1+I x)} /todo:label  /{"br ",$[`~S x;$1+I x;"$",($S x),$I x]}
wblit:{y;i:I x;((`i`u`k`j`f!0x4141424244)t),,$[`f~t:S x;C i+!8;`i~t;leb i;`u~t;leb i;leb8 C i+!8]}

wbneg:{wbsub[x;($[`f~t:S x;0x44,(8#0x00);(0x4141424244@`i`u`j`k?S x),0x00];*y)]}
wbnot:{(*y),"\ni32.eqz"}
wbnop:{y;""}
wbret:{(*y),0x0f}

wbswc:{n:#$[I x;y;y,:,""]
 r:@[,/n#(,0x0240);1;$[t:S x;wbtype t;0x40]],(*y),0x0e,(leb n-2),(,/leb'!n-1),0x0b
 r,,/(1_y),'0x0c,/:(leb'|!n-1),\:0x0b}
wbcnd:{(*y),0x04,$[`~t:S x;0x40;wbtype t],(y 1),$[3~#y;0x05,(y 2);""],0x0b}

wbtyp:{y;wbtype S x} /keep as symbol

wb2: {(,/z),x(`i`u`j`k`f?S y)}
wb2s:{"\n"/:(z 0;z 1;(wbtype t  ),x,wbsign@t:S y)}
wbadd:wb2 0x6a6a7c7ca0;wbsub:wb2 0x6b6b7d7da1
wbmul:wb2 0x6c6c7e7ea2;wbdiv:wb2 0x6d6e7f80a3;wbmod:wb2 0x6f70818200 
wband:wbbnd:wb2 0x7171838300;wborr:wbbor:wb2 0x7272848400;wbxor:wb2 0x7373858500
wbneq:wb2 0x4747525262;wbeql:wb2 0x4646515161
wbmor:wb2 0x4a4b555664;wbgte:wb2 0x4e4f595a66;wbles:wb2 0x4849535463; wblte:wb2 0x4c4d575856
wbshl:wb2 0x7474868600;wbshr:wb2 0x7576878800

wbtype:`i`u`j`k`f`!0x7f7f7e7e7c00
wbsign:`i`u`j`k`f!("_s";  "_u"; "_s"; "_u";   "")
wbinst:`Memorysize`Memorygrow`Memorycopy`Memoryfill!(0x3f00;0x4000;0xfc0a0000;0xfc0b00)
wbinst,:`Memorysize2`Memorygrow2`Memorycopy2`Memorycopy3!(0x3f01;0x4001;0xfc0a0001;0xfc0a0100)
wbinst,:`I32clz`F64abs`F64sqrt`F64floor`F64copysign`F64min`F64max`"F64reinterpret_i64"`"I64reinterpret_f64"!,'0x67999f9ca6a4a5bfbd

wbcast:(0x0000a7a7aa;0x0000a7a7ab;0xacad0000b0;0xadad0000b1;0xb7b8b9ba00) /iujkf


leb:{leb8(`c@,x),4#$[x<0;0xff;0x00]}
leb8:{r:!0;B:1v:,/+(8#2)\\256/256+|x;while[B&(8<#v)&+/v;b:-7#v;s:*b;v:-7_v;r,:$[((0~&/v)|0x00~s)&(+/v)|s;1;B:0],b];_2//'-8^r}

wbsys:`Exit`Args`Arg`Read`Write`ReadIn`Native!wbtype@/:(`i`;,`i;3#`i;4#`i;5#`i;3#`i;3#`k)
/(import \"env\" \"Exit\"  (func $Exit  (param i32)))
/(import \"env\" \"Args\"  (func $Args  (result i32)))
/(import \"env\" \"Arg\"   (func $Arg   (param i32) (param i32) (result i32)))
/(import \"env\" \"Read\"  (func $Read  (param i32) (param i32) (param i32) (result i32)))
/(import \"env\" \"Write\" (func $Write (param i32) (param i32) (param i32) (param i32) (result i32)))
/(import \"env\" \"ReadIn\" (func $ReadIn (param i32) (param i32) (result i32)))
/(import \"env\" \"Native\" (func $Native (param i64) (param i64) (result i64)))


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

