/k -fsrc k.w w.k
/   compile src

compile:{macro:funcs:extrn:(0#`)!();table:mem:();class./:mf'-1_"}"\:*(_10 92)\:x;p:parse[;macro];`macro`funcs`extrn`table`mem!(macro;p'funcs;extrn;table;mem)}
afun:{exported:1<+/":"=x; x:"."\:@[x;&x=":";"."]; name:`$*x; rtyp:`$x 1; argt:,/`$'""\:x 2;name!,`exported`rtyp`argt`body!(exported;rtyp;argt;y)}
class:{isnum:{x~^x:,/0+("0";x;"9")};$["!"~(:x);mem,:`off`data!(.-1_x;."0x",y) ;isnum(*x);table,:`off`funcs!(.-1_x;";"\:y); ":"~(:x);macro,:(`$-1_x)!,y; #y;funcs,:afun[x;y]; extrn,:afun[x;y]];}
trim:{(+/(_10 13 32)=*x)_x}
mf:{r:"{"\:x;(trim/:r 0;r 1)}
parse:{[x;macro]x}

a:compile src
macro:a`macro

tokenize: {r:0#`!0;{r::r,$[0~n:(l:noun x)1;0#r;(*n)!,1_n];*l}/:x;r}
noun:     {t:toks;r:0;({$[~#t;0;0~r::(*t)x;1;0]};{t::1_t;x})/:xpnd@ws/:x;((*r)_x; :r)}


first:{0+*x}
toks:(pSeq:{x;0} /todo
 pTyp:{$[~ctype first x;0;*(2>#x;r:(1;(`typ;`$*x)));r;alphanum 0+x 1;0;r]}
 pSym:{$[~alpha first x;0;r:(alphanum 0+x)?0;(r;(`sym;`$r#x));0]}
 pNlp:{$["/"~*x;(1;,`nlp);0]}
 pCon:{$[~num first x;0;$["j"~t:*(*n_x;i:&(n:(num 0+x)?0)#x);(n+1;(`con;`J;i));(n;(`con;`I;i))]]}    / todo 234f
 pOp: {$[(#ops)~n:(x{y~(#y)#x}/:ops)?,1;0;$[*(":"~*l_x;l:#s:ops n);(1+l;(`op;s,":"));(l;(`op;s))]]}
 pCnd:{$[3>#x;0;~"$["~x 0 1;0;*(2+*r;(,`cnd),1_ :r;r:pSeq 2_x)]}
 pSem:{$[";"~*x;(1;(`sem;0));0]})
 

ws:  {$[space first x;1_x;x]}
xpnd:{$[k: :pSym x;$[0~k;x;(#macro)>(!macro)?k 1;(macro k 1),(#$k 1)_x;x];x]}
ops: " "\:(">=' %' \' << >> <' <= >= ?/ ?' : + - * % ~ _ * | \ & ^ !") /long-first
set: {@[256#0;x;1]};num:set"0"+!10;alphanum:num+alpha:set"a"+!26;ctype:set"VCIJF"+0;space:set@!33;op:set"+-*%"+0

 \tokenize "-?/3+4"

code:(.a`funcs)@\:`body

\
expand:{k:!macro;$[(#k)>k?x;macro x;""]}
symbol:...
noun:{x:$[#s:symbol(x);(expand s),x;x]
c:~x;i:0+c;$["("~i;          seq[x;")"]    /(..)
             ctype i;     (`$c;1_x)        /type
	     #s:sym(x);   x:_#s
;]}
and:{[f;x]({$[~#f;0;(*f)x;1;0]};{f::1_f;x})/:x;~#f}
set:{@[256#0;x;1]};num:set"0"+!10;alphanum:alpha+num;ctype:set"VCJF";op:set"+-*%";space:set(!33)


todo convert float constants to hex. 0faabbccddee0011223344

package main
import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
)
func xx(u uint64) string {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, u)
	return hex.EncodeToString(b)
}
func main() {
	f := 3.141592653589793
	s := xx(math.Float64bits(f))
	fmt.Println(s)
	b, _ := hex.DecodeString(s)
	u := binary.LittleEndian.Uint64(b)
	fmt.Println(math.Float64frombits(u))
}
