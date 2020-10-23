/k -fsrc k.w w.k
/   compile src

compile:{macro:funcs:extrn:(0#`)!();table:mem:();class./:mf'-1_"}"\:*(_10 92)\:x;p:parse[;macro];`macro`funcs`extrn`table`mem!(macro;p'funcs;extrn;table;mem)}
afun:{exported:1<+/":"=x; x:"."\:@[x;&x=":";"."]; name:`$*x; rtyp:`$x 1; argt:,/`$'""\:x 2;name!,`exported`rtyp`argt`body!(exported;rtyp;argt;y)}
class:{isnum:{x~^x:,/0+("0";x;"9")};$["!"~(:x);mem,:`off`data!(.-1_x;."0x",y) ;isnum(*x);table,:`off`funcs!(.-1_x;";"\:y); ":"~(:x);macro,:(`$-1_x)!,y; #y;funcs,:afun[x;y]; extrn,:afun[x;y]];}
trim:{(+/(_10 13 32)=*x)_x}
mf:{r:"{"\:x;(trim/:r 0;r 1)}

parse:   {[x;macro]x,`ast!,ex tokenize"(",(x`body),")"}
tokenize:{r:0#`!0;{r::r,$[0~n:(l:tok x)2;0#r;(*n)!(,l 1),1_n];*l}/:x;r}
tok:     {t:toks;r:0;x:({$[~#t;0;0~r::(*t)x;1;0]};{t::1_t;x})/:xpnd/:ws/:x;((*r)_x;(*r)#x; :r)}

first:{0+*x}
toks:(pSeq:{x;0} /todo
 pTyp:{$[~ctype first x;0;*(2>#x;r:(1;(`typ;`$*x)));r;alphanum 0+x 1;0;r]}
 pSym:{$[~alpha first x;0;r:(alphanum 0+x)?0;(r;(`sym;`$r#x));0]}
 pNlp:{$["/"~*x;(1;,`nlp);0]}
 pCon:{$[~num first x;0;"0x"~2#x;(18;(`con;`F;-3$unhex 16#2_x));"j"~t:*(*n_x;i:&(n:(num 0+x)?0)#x);(n+1;(`con;`J;i));~"."~t;(n;(`con;`I;i));0~d:pCon (n+1)_x;(n+1;(`con;`F;0.+i));(m;(`con;`F;&(m:1+n+*d)#x))]}    / todo 234f
 pOpa:{$[(#ops)~n:(x{y~(#y)#x}/:ops)?,1;0;$[*(":"~*l_x;l:#s:ops n);(1+l;(`op;s,":"));(l;(`op;s))]]}
 pCnd:{$[3>#x;0;~"$["~x 0 1;0;*(2+*r;(,`cnd),1_ :r;r:pSeq 2_x)]}
 pSem:{$[";"~*x;(1;(`sem;0));0]}
 pBra:{$["("~*x;(1;(`bra;"("));0]}
 pCbr:{$[")"~*x;(1;(`cbr;")"));0]}
 pCpa:{$["]"~*x;(1;(`cpa;"]"));0]})
 
unhex:{,/{+/_16 1*"0123456789abcdef"?x}'(,/(2\#x;2))#x}
ws:  {$[space first x;1_x;x]}
xpnd:{$[0~k: :pSym x;x;(#macro)>(!macro)?s:k 1;(macro s),(#$s)_x;x]}
ops: " "\:(">=' %' \' << >> <' <= >= ?/ ?' $[ : + - * % ~ _ * | \ & ^ ! ? < > =") /long-first
set: {@[256#0;x;1]};num:set"0"+!10;alphanum:num+alpha:set"a"+!26;ctype:set"VCIJF"+0;space:set@!33

seq:{r:,`seq;e:0;(r;1_({~0~*e::ex x};{r::r,,*e;x:e 1;$[`sem~*!x;1_x;x]})/:x)}  /(`seq;e0;e1;..)
cnd:{r:seq x;((`cnd,1_*r);r 1)}                                                /(`cnd;e1;e2;..)
mon:{$["$["~o:**.x;cnd 1_x;((`mon;o;*a);(a:ex 1_x)1)]}                         /(`mon;"+";arg)
dya:{a:*x;m:*n:mon x 1;((`dya;m 1;a;m 2);n 1)}                                 /(`dya;"+";a;b)

exx:{p:*!x;$[`bra~p;seq 1_x;+/`sym`con=p;(p,*x;1_x);`op~p;mon x;(0;x)]}        /(e;tail)
ex: {h:exx x;$[`op~*!(h 1);dya h;h]}

/ \e:*ex tokenize"$[3;4;5]"
/ \e:*ex tokenize"+4.5"
/ \e:ex tokenize"1+2"

a:compile src
 \a[`funcs;`asi;`body]
 \a[`funcs;`asi;`ast]


wasm:{}
leb:{r:"";({c:x/128;s:x/64;x%:128;c+:$[((~s)+~x~-1)*(~~x)+~~s;128;0];r::r,_c;$[c<128;0;1]};{x%128})/:x;r}
whdr:0x0061736d01000000
/wsig:0x60,...

\
func leb(v int64) []c { // encode signed leb128
	var b []c
	for {
		c := uint8(v & 0x7f)
		s := uint8(v & 0x40)
		v >>= 7
		if (v != -1 || s == 0) && (v != 0 || s != 0) {
			c |= 0x80
		}
		b = append(b, c)
		if c&0x80 == 0 {
			break
		}
	}
	return b
}
