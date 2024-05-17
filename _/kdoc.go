package main

// this file builds ktye.github.io/kdoc.htm
// $cd ..;go run _/kdoc.go

import (
	"fmt"
	"html"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	O := fmt.Print
	H := html.EscapeString
	R := func(x string) string { b, _ := os.ReadFile(x); return string(b) }
	h := func(x, y, z string) { O(`<hr/><span id="` + x + `"><h2>` + y + `</h2>` + z + `</span>`) }
	tick := func(x string) string { return strings.ReplaceAll(x, "°", "`") }

	O(head)

	//refcard
	x := R("readme")
	a := strings.Fields("flp add neg sub fst mul sqr div til key wer min rev max asc les dsc mor grp eql not mtc enl cat srt cut cnt tak flr drp str cst unq fnd typ atx val cal abs sin cos exp log find angle imag conj types intro")
	b := strings.Fields("Flp Add Neg Sub Fst Mul Sqr Div Til Key Wer Min Rev Max Asc Les Dsc Mor Grp Eql Not Mtc Enl Cat Srt Cut Cnt Tak Flr Drp Str Cst Unq Fnd Typ Atx Val Cal Abs Sin Cos Exp Log Find Angle Imag Conj types intro")
	for i := range a {
		x = strings.Replace(x, a[i], `<a href="#`+b[i]+`">`+a[i]+`</a>`, 1)
	}
	O(x)

	O(tty)

	//toc todo links
	x = toc
	a = strings.Fields("kinit z.k intro invoke types literals adverbs control-flow symbols heap memory/alloc kvm instructions ir syscalls abstract-machine extend/native")
	b = strings.Fields("kinit z.k intro invoke types literals adverbs control      symtab  heap allocator    kvm kvm          ir syscalls mach             extend")
	for i := range a {
		x = strings.Replace(x, a[i], `<a href="#`+b[i]+`">`+a[i]+`</a>`, 1)
	}
	O(x)

	//src
	O(src)
	m := make(map[string]string)
	k := []string{"a.go", "k.go"}
	d, _ := os.ReadDir(".")
	for _, f := range d {
		if strings.HasSuffix(f.Name(), ".go") && !strings.HasSuffix(f.Name(), "_test.go") {
			x := R(f.Name())
			m[f.Name()] = x
			if f.Name() != "a.go" && f.Name() != "k.go" {
				k = append(k, f.Name())
			}
		}
	}
	cid := 0
	ex := func(s string) int {
		i := strings.Index(s, ")")
		v := strings.Split(s[:i], ", ")
		for i := range v {
			O(`<a href="#` + v[i] + `" id='c:` + strconv.Itoa(cid) + `'>` + v[i] + `</a>`)
			cid++
			if i < len(v)-1 {
				O(", ")
			}
		}
		return i
	}
	fn := func(s string) int {
		i := strings.Index(s, ")")
		v := strings.Split(s[:i], ", ")
		O("(" + v[0] + ", ")
		v = v[1:]
		for i := range v {
			O(`<a href="#` + v[i] + `" id='c:` + strconv.Itoa(cid) + `'>` + v[i] + `</a>`)
			cid++
			if i < len(v)-1 {
				O(", ")
			}
		}
		return i
	}
	L := func(s string) {
		r := regexp.MustCompile(`\w+\(`)
		v := r.FindAllIndex([]byte(s), -1)
		i := 0
		for _, ab := range v {
			a, b := ab[0], ab[1]-1
			O(H(s[i:a]))
			O(`<a href="#` + s[a:b] + `" id='c:` + strconv.Itoa(cid) + `'>` + s[a:b] + `</a>`)
			cid++
			i = b
			if s[a:b] == "Export" {
				O("(")
				b++
				i = b + ex(s[b:])
			} else if s[a:b] == "Functions" {
				b++
				i = b + fn(s[b:])
			}
		}
		O(H(s[i:]))
	}
	F := func(s string) int {
		i := strings.Index(s, "(")
		O(`<span id="` + s[6:i] + `">`)
		O("\nfunc <span class='q'>" + s[6:i] + "</span>")
		j := 1 + strings.Index(s[1:], "\n")
		k := j
		if strings.Index(s, "\n\t") == j { //multi line function body
			j = 2 + strings.Index(s, "\n}")
		}
		L(s[i:k])
		O("<span id='c:" + s[6:i] + "' class='n'>\n</span>")
		if k != j {
			L(s[k:j])
		}
		O("</span>")
		return j
	}
	S := "\n"
	for _, f := range k {
		s := m[f]
		s = strings.ReplaceAll(s, "package main\n\nimport (\n\t. \"github.com/ktye/wg/module\"\n)\n\n", "")
		s = strings.ReplaceAll(s, "package main\n\n", "")
		s = strings.ReplaceAll(s, "import . \"github.com/ktye/wg/module\"\n\n", "")
		S += s
	}
	for i := 0; ; {
		j := strings.Index(S[i:], "\nfunc")
		if j < 0 {
			O(H(S[i:]))
			break
		}
		j += i
		O(H(S[i:j]))
		i = j + F(S[j:])
	}

	h("z.k", "z.k", zk+H(R("z.k")))
	h("intro", "intro", intro)
	h("invoke", "invoke", invoke)
	h("types", "k type system", types)
	h("literals", "literals", tick(literals))
	h("adverbs", "adverbs", tick(adverbs))
	h("control", "control flow", control)
	h("heap", "heap memory", heap)
	h("allocator", "memory allocator", allocator)
	h("symtab", "symbol table", tick(symtab))
	h("kvm", "kvm - virtual machine - instruction set", kvm)
	h("ir", "intermediate representation", tick(ir))
	h("mach", `abstract machine model`, mach)
	h("syscalls", `system interface`, syscalls)
	h("extend", "extend", tick(extend))
	h("embed", "embed", embed)
	h("f77", "fortran", f77)
	h("c", "c", cc)
	h("wasm", "WebAssembly", wasm)
	h("wasi", "WebAssembly (standalone/wasi)", wasi)

	O(tail)
}

const head = `<!KDOCTYPE html>
<head><meta charset="utf-8">
<link rel="icon" type="image/png" sizes="32x32" href="kelas32.png">
<link rel="icon" type="image/png" sizes="16x16" href="kelas16.png">
<title>kdoc</title>
<style>
*{font-family:monospace;margin:0}
.q{color:green}.q:hover{cursor:pointer}
.l{color:blue}.l:hover{cursor:pointer}
.n{display:none}
.h{font-weight:bold}
.hidden{display:none}
body{display:flex;flex-direction:column;height:100vh;overflow:hidden}
a{text-decoration:none}
#ref{display:inline-block;background:#ffe}
#tty{display:inline-block;background:#004687;color:white;width:100%;margin-left:0.5em;margin-right:0.5em;outline:none;overflow:hidden}
#top{display:flex}
#doc{overflow:auto;outline:none;margin:0}
#mono{position:absolute;left:100;top:100;visibility:hidden}
</style>
</head>
<body onload="init()">

<script>
function ge(x){return document.getElementById(x)}
function ce(x){return document.createElement(x)}
function ct(x){return document.createTextNode(x)}
function tc(x,y){x.textContent=y;return x}
function ac(x,y){x.appendChild(y);return y}
function id(x,y){x.id=y;return x}
function fe(u,f){fetch(u).then(r=>r.text()).then(f)}
function rm(p){while(p.firstChild)p.removeChild(p.firstChild);return p}
const td_=new TextDecoder("utf-8"),su=x=>td_.decode(x),te_=new TextEncoder("utf-8"),us=x=>te_.encode(x)
let cur=null
function callers(x,y){rm(y);
 let c=Array.from(document.querySelectorAll("a")).filter(a=>a.id.startsWith("c:")&&a.getAttribute("href")=="#"+x.textContent)
 tc(y,"\ncallers("+(c.length)+")\n")
 c.forEach(a=>{let d=ce("a");d.href="#"+a.id;d.textContent="│"+a.parentElement.id;ac(y,d);ac(y,ct(": "+linestr(a)+"\n"))})
}
function linestr(a){let r=a.textContent,x=a;while(1){x=x.nextSibling;let s=x.textContent,i=s.indexOf("\n");if(i>=0){r+=s.slice(0,1+i);break};r+=s};x=a;while(1){x=x.previousSibling;let s=x.textContent,i=s.lastIndexOf("\n");if(i>=0){r=s.slice(i)+r;break};r=s+r};return r.trim()}
function init(){
 Array.from(document.getElementsByClassName("q")).forEach(x=>x.onclick=function(){let y=ge("c:"+x.textContent);callers(x,y);y.classList.toggle("n")})
 kinit();ge("tty").onclick=end;ge("tty").onkeydown=key
}
//function goto(x){let s=ge("fn:"+x);s.scrollIntoView();s.classList.add("h")}
function ttysize(){let mono=ge("mono");let n=mono.textContent.length;return[Math.floor(ge("tty").clientHeight/mono.clientHeight),Math.floor(n*ge("tty").clientWidth/mono.clientWidth)]}
function end(){let b=ge("tty"),s=window.getSelection();s.removeAllRanges();let r=document.createRange();r.selectNodeContents(b);r.collapse(false);s.addRange(r);b.focus()}
function key(e){if(e.key!="Enter")return;e.preventDefault();let s=ge("tty").textContent.split("\n").slice(-1)+"\n";O("\n");try{K.repl(KC(s.startsWith(" ")?s.slice(1):s));O(" ");end()}catch(e){kinit()}}
window.onhashchange=e=>{if(cur)cur.classList.remove("h");cur=ge(window.location.hash.slice(1));console.log("cur",cur);if(cur==null)return;cur.classList.add("h");cur.scrollIntoView()}



let /*there be*/ K
let C=()=>new Int8Array(K.memory.buffer),I=()=>new Int32Array(K.memory.buffer),J=()=>new BigInt64Array(K.memory.buffer),F=()=>new Float64Array(K.memory.buffer),lo=x=>Number(BigInt.asUintN(32,x))
let kenv={env:{ 
 Exit:  function(x      ){},
 Args:  function(       ){return 0},
 Arg:   function(x,y    ){return 0},
 Read:  function(a,b,c  ){return -1},
 Write: function(a,b,c,d){O(su(new Uint8Array(K.memory.buffer,c,d)));return 0},
 ReadIn:function(x,y    ){return 0},
 Native:function(x,y    ){let i=lo(x);K.dx(x);return xcal[i](K.Atx(4n,y))}}}
function O(s){let o=ge("tty");console.log("kout", s);o.textContent=(o.textContent+s).split("\n").slice(-25).join("\n")}
function kinit(){let s,sz=x=>{s=x.byteLength;return x}
 fetch("k.wasm").then(r=>r.arrayBuffer()).then(r=>WebAssembly.instantiate(sz(r),kenv)).then(r=>{
 K=r.instance.exports
 K.kinit()
 let[rows,cols]=ttysize();K.dx(K.Asn(Ks("l."),K.Atx(Ks("lxy"),K.Val(KC((cols-2)+" 50")))))
 O("ktye/k "+s+"\n ");end()})
}
let
KC=x=>{let r=K.mk(18,x.length  );C().set(("string"===typeof x)?us(x):x,lo(r));return r},
Ks=x=>K.sc(KC(x))


let full=s=>{let a=[" ▲ "," ▼ "];ge("top").style.display=(s.textContent==a[0])?"none":"flex";s.textContent=a[1-a.indexOf(s.textContent)]}
</script>

<div id="top">
<pre id="ref">`

const tty = `</pre>
<pre id="tty" contenteditable spellcheck="false"></pre>
<pre id="toc">`

const toc = `kinit . . . . . . . . . . .z.k

intro                   invoke 

types literals    memory/alloc

adverb            control-flow

examples algebra stats queries

html/js/wasm         embedding





implementation
types        heap    allocator
symbols    variables     scope
exec    kvm       instructions
compilers boot   extend/native
            f77 c go wasm wasi
ir  abstract-machine  syscalls
`

const src = `</pre>
</div>
<pre id="doc"><span class="l" onclick="full(this)" style="position:sticky;float:right;top:0"> ▲ </span>`

const zk = `
z.k is placed in the initial memory section at 600 by <a href='#zk'>zk</a> called from <a href='#kinit'>kinit</a>
it contains the runtime part of the interpreter that is written in k. it also prepopulates the <a href="#symtab">symbol table</a>.

todo: s-dot@ k@ l@ space sin/cos..
`

//----------------------------------------------------------------------------------------------|

const intro = `
ktye/k is an implementation of the k programming language.
this page documents the implementation of the interpreter.
it is not an introduction to k or array programming in general.

the focus of the implementation is a simple system at small size, portability and the possibility
to embed and extend it into applications.

 size     is measured in bytes of the wasm binary, below 40k. the actual size in bytes is the
          number in the banner of the repl.
 port     the source code looks like go but it isn't. it's a structural assembly language that
          also happens to satisfy the go compiler.
          a <a href="https://github.com/ktye/wg">separate compiler</a> parses the source and writes a k table <a href="#ir">intermediate representation</a>.
	  <a href="#compilers">compilers</a> written in k transform the IR to other languages: <a href="#cc">c</a> <a href="#go">go</a> <a href="#wasm">wasm</a> <a href="#wasi">wasi</a>.
	  the fact that the source can also be directly compiled with a go compiler simplifies
	  bootstrapping testing and development by no small amount.
 embed    ..
 extend   ..

try/compile/run
try in the repl above, or see the any of the ports <a href="#cc">c</a> <a href="#go">go</a> <a href="#wasm">wasm</a> <a href="#wasi">wasi</a> <a href="#f77">fortran</a>.
`

const invoke = `
$k              /interactive
$k a.k b.k      /run both files in order, drop to interactive mode
$k a.k -e       /run a.k then exit
$k a.k -e 'f x' /run a.k, eval f x then exit
$k k.t          /runs tests, e.g. <https://raw.githubusercontent.com/ktye/i/master/k.t>k.t</a>
-e stands for both: exit and eval.
`

const types = `
k values are represented by a 64 bit integer called K in the go implementation.
depending on the type, the value has a different meaning.
the type is stored in the upper 5 bits (and accessed by t := x&gt;&gt;59).
t&lt;16 are atoms and t&gt;16 are vectors/compound types
vector types t&lt;23 are flat
the base type of atoms or vectors is t&15

 atoms      t  t  vectors         t functions         t compound
 c char     2  18 C chars         0 v primitives      23 L list
 i int32    3  19 I ints         10 m compositions    24 D dict
 s symbol   4  20 S symbols      11 d derived         25 T table
 f float64  5  21 F floats       12 p projections
 z complex  6  22 Z complexs     13 l lambda
                                 14 x native/extern

values of type c i s carry their value withing the k value in the lower 32 bits.
all other values live in heap memory. the lower 32 bits are the index to the start of the data section.
flat vectors of type I S F Z store their vector data consecutive as 4 4 8 16 byte widths per element.
compound types m d p l x L D T store 64 bit k values in the data section.
the length of vector values is accessed by <a href="#nn">nn</a> and stored as int32 12 bytes before the data.
values the live in heap memory also have their <a href="#allocator">refcount</a> stored at 4 bytes before the data.

L is a general list that <a href="#uf">collapses/unifies</a> to flat vector if all values are atoms of the same type.
D and T are always a 2 element list (keys and values of the same length) and their length field stores the length of the keys.
D and T only differ by type, but a table is more restricted: it's keys must be symbols and values a general list (L)
each containing a vector of the same size.

the function types m d p l x store:
 m(function-list)                   forming the composition <a href="#calltrain">call train</a>
 d(base-function;adverb)            the adverb value is the function table index at 85 + <a href="#Ech">1</a> <a href="#Rdc">2</a> <a href="#Scn">3</a> 
 p(base-function;arglist;emptylist) <a href="#callprj">call projection</a>
 l(code;string;locals)              <a href="#lambda">call lambda</a>. code is the <a href="#kvm">instruction list</a> and locals a list of symbols including the arguments. the length field stores the arity.
 x(index/ptr;string)                <a href="#native">call native</a>, see <a href="#extend">extend</a>

primitive verbs store a value between 0 and 64 that corresponds to the index into the indirect function table defined in <a href="#init">init</a>.
primitives are ambivalent, their value corresponds to the monadic case and is <i>fixed</i> at call time, 
e.g. when called with two arguments the function tables is indexed at value x+64.

some k values also store the source location offset in the lower bits of the upper 32 bit value.
the offset is retrieved by 0xffffff & int32(x&gt;&gt;32) at execution time and stored in the global <i>srcp</i>.
it is the offset to the executed source code which is stored in a k value written to <a href="#src">16</a> in the <a href=#heap>heap</a>.
all k source is catenated to this value by the parser.
`

const literals = `
literals are parsed by the <a href="#tok">tokenizer</a>.
the functions <a href="#tchr">tchr</a> <a href="#tnms">tnms</a> <a href="#tvrb">tvrb</a> <a href="#tpct">tpct</a> <a href="#tvar">tvar</a> <a href="#tsym">tsym</a> parse characters, numbers, primitive verbs, punctuation, variables and symbols.
<span id="classes">
class  bytes
  1    :+-*%!&amp;|&lt;&gt;=~,^#_$?@.&#39;/\
  2    abcdefghijklmnopqrstuvxwyzABCDEFGHIJKLMNOPQRSTUVWXYZ
  4    0123456789
  8    &#39;/\
 16    ([{
 32    \n;)]}
 64    abcdefghijklmnopqrstuvxwyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789)]}
128    0123456789abcdef
</span>

<a href="#is">is</a>(b,c) tests if the char value b belongs to class c. classes can can also be combined:
e.g. <a href="#tvar">tvar</a> tests if the first character is in class 2 (alphabetic) and the following characters are alphanumeric 2+4.

type  literal
 c    "a" or 0x23
 i    -123
 f    .5  leading/trailing 0 may be omitted
 s    °ab or °"anything" shortform only for valid varnames
 z    1.5a30  abs and angle in degree, both maybe floats, 3a is 3a0

vector types are combined                  lists
 C     "ab"                                (1;(2;3 4);5)   (1;2;3) single element lists collapse to 1 2 3
 I     1 2      single space only          dicts have no literal form, use the functional form   °a°b!2 3
 F     3 4.2 5  any float promotes to F                                              or !/+-2^(°a;2;°b;3)
 S     °a°b°c                              tables are flipped dicts:                      +°a°b!(2 3;4 5)
 Z     1a 2a30                             also key tables S!T and T!T
`

const adverbs = `
`

const control = `
`

const heap = `
the <a href="#mach">abstract machine</a> has access to a linear range of heap memory. memory can only <a href="#Memorygrow">grow</a> but never shrinks.
in the initial state it's size is one block of 64kb. addresses are byte indexes into the memory section.
low addresses have special meaning, setup by <a href="#init">init</a> or created by <a href="#kinit">kinit</a>.
memory above 4k is managed by the <a href="#allocator">allocator</a> and used to store k values.

byte-addr
   0....7  keys see <a href=#symtab>symbol table</a>
   8...15  vals
  16...19  <a href="#src">src</a>(int32) catenated for <a href="#trap">error indication</a>
  20..127  free list see <a href="#allocator>memory allocator</a>
 128..131  currently allocated memory (log2)
 132..226  <a href="#class">character classes</a> for the parser (starts at 100) ⎫
 228..252  :+-*%!&|<>=~,^#_$?@.':/:\:                                     ⎬ text section
 253..279  vbcisfzldtcdpl000BCISFZLDT  type symbols                       ⎮
 280.....  <a href="#z.k">z.k</a> (embedded source)                                          ⎭
 2k....4k  <a href="#kvm">kvm</a> stack
 4k....4g  vector-space
`

const allocator = `
memory is allocated by the k implementation by dividing the linear memory section to chunks on request.
if more memory is needed, the total memory section grows before dividing it.
it is a buddy allocator, that divides by chunck sizes of powers of two.
memory is requested by <a href=#mk>mk(t,n)</a> which <i>makes</i> a k value for the given type and number of elements.
e.g. r := mk(19,1000) creates an integer vector with 1000 elements.
the required size is 4000 to store 1000 int32 values + a 16 byte header, so 4016.
the next power-of-2 block size where it fits into is 4096 or 2^12, called a block of bucket type 12.
the smallest block is type 5 or 32 bytes with space for 16 chars, 4 ints or 2 floats (or other k values) or 1 complex number.

the allocator keeps a free list (linked-list) for each bucket type t in the heap memory at 4*t.
so I32(4*12) returns memory index to the start of the next free block for a type-12 memory chunck.
it also writes the next free block to the free list at the top location SetI32(4*12, ..).

when no free chuck is available for the given type, the next larger chunck is tried and splitted into two.
when nothing is found, total memory is doubled.

memory is freed when it's refcount drops to 0 by calls to <a href="#dx">dx</a> which calls <a href="#mfree">mfree</a>.
in this case it is prepended to the free list at it's bucket size location, and stores index of the current free block in it's memory.
memory is never released to the system.
divided blocks are also never merged to larger blocks. to prevent an accumulation of small chuncks, no primitive splits memory.
memory is only reused if the refcount is 1 and the new value has the same bucket type (e.g. when dropping the last value from a vector).

 <a href="#mk">mk(t,n)</a> make k value, also shortcuts <a href="#Kc">Kc</a> <a href="#Ki">Ki</a> <a href="#Kf">Kf</a> a href=#Kz">Kz</a> for atoms,
 <a href="#l2">l2(x,y)</a> make a list of two k values
 <a href="#rx">rx(x)</a> increase refcount, <a href="#dx">dx(x)</a> decrease refcount (from derive: d/dx)
 <a href="#rl">rl(x)</a> increase refcount of each value of a list, but not x itself
also low level <a href="#alloc">alloc</a> <a href="#mfree">mfree</a> <a href="#bucket">bucket</a>.
`

const symtab = `
symbols are stored as interned strings. their value is 32 bit integer.
the integer is an offset to the symbol table, a k list of chars stored at memory location 0.
e.g. we can use °i@ to convert a list of symbols to their integer value: °i °°x°y°z returns 0 8 16 24.

the symbol table needs to be searched only once for each symbol: when it is created by the tokenizer.
<a href="#tsym">tsym</a> uses <a href="#sc">sc</a> symbol from char which returns a symbol value with the correct offset.
if the symbol is new, the character vector for the new symbol is appended to the symbol table. that means the value at memory 0 point to a new location.

symbol values depend on the creation order. <a href="#z.k">z.k</a> is executed first and parses a set of symbols that should have known values, such as
°x°y°z the <i>while</i> keyword and some k implementations of primitives that use <a href="#kx">kx</a> and <a href="#kxy">kxy</a>.

the symbol offset values also serve a second purpose: they serve as the offset in the <b>lookup table for variables</b>.
variables are stored in a k list of type L that is stored at location 8.
both lists (keys at 0 and values at 8) are always in sync and append-only.
variable lookup never searches, it already knows it's offset. but it uses two indirections as the number of symbols or variables is not restricted.

as a consequence there is no <i>value error</i> for undefined variables:
when parsing the symbol for a variable that does not exist, both the symbol table and the value list is extended. the new value is 0 (the verb not the integer).

the value list is the only location where variables are stored. both local and global variables.
opposed to most other implementations of k, variables have dynamic scope and there is no k-tree or namespaces.

<span id="#scope">dynamic scope means, lambda function arguments and their local variables shadow variables with the same name.
when a lambda function is called, the values corresponding to the same symbol as the locals are saved and restored when the function returns.
this also means that a lambda function can modify a local variable of it's caller when it is treated as a global variable.
</span>`

const kvm = `
`

const ir = `
the IR stores the parse tree of the source of the k interpreter as a k table.
this is about the source notation and the compilers that transpile it to other languages not about the execution of k code, see <a href="#kvm">kvm</a> for the latter.

node type                    i-value             s-value
 prg     root node           -                   prog/libname    first node only
 mem     memory segment      #64k blocks         °a|°b           a|b: memory1|memory2
 con     constant            -                   name            child: lit
 var     global variable     -                   name            child: lit
 lit     literal (con|var)   val(32bit)|C-index  type
 tab     func table entry    index               func name
 fun     function            exported            func name       children: args res locs ast dfr
 arg     func argument       -                   type            child: sym
 sym     symbol              1(global)|0N(func)  name
 res     return value        -                   type            unnamed
 loc     local var decl      -                   type            child: sym
 ast     func ast root       -                   -               one per func
 stm     statement list      -                   -
 ret     return              -                   °|type          children: return values, s-type only for single res
 cal     function call       -                   func name       children: args
 cli     indirect call       #args               res-type        children: func-expr args arg-types
 drp     drop return vals    -                   -               child: cal
 get     get local           -                   varname
 Get     get global          -                   varname
 lod     load                -                   type(bijf)      child:    addr
 sto     store               -                   type(bijf)      children: addr, value
 asn     assignment          1(global)           varname         children: expr
 cst     cast                -                   dst type        2 children: typ(src), arg
 typ     type                -                   type
 cnd     if condition        -                   °|result-type   2|3 children: if then [else]
 swc     switch              1(has default)      °|result-type   children: expr cases [default]
 jmp     break/continue      1(break)|0(cont)    label
 for     loop                1(simple)           label           children: (cond|nop) (post|nop) body
 dfr     defer stmt node     -                   -               child: cal
 nop     ignore              -                   -
unary operator nodes
 neg|not                     1                   type            1 child
binary operator nodes
 eql|les|mor|gte|lte|and|orr 2                   type            2 children
 add|sub|mul|div|mod|shr|shl	
 xor|neq|ant(andnot)|bnd|bor(&& ||)			     
 
types: °i°u°j°k°f!(i32;u32;i64;u64;f64)
`

const mach = `
 the <a href="#ir">IR</a> targets an abstract machine model that the target languages need to implement.
 it is roughly based on the characteristics of wasm:
 - basic data types are i32 i64 f64
 - infinite set of registers (local variables) 
 - one linear memory section (heap) addressable by a 32bit index (0..)
 - memory can grow by sections of 64k but never shrinks
 - code and data are separated, no jit
 - in addition to i32 i64 f64, heap can read/write a <a href="#I8">single byte as i32</a>
 - heap access must be aligned to the size of the data type
 - no stack, use locals or heap, no pointers to local variables, call by value
 - no multi return values
 - only wasm floating point ops are available, e.g. f64.sqrt everyting else is implemented in software, e.g. <a href="#cosin_">sin cos</a> <a href="#exp">exp</a>
 - undefined order of evaluating of function arguments
 - sequential operations, no threads
 - a function table and indirect calls by index is available (e.g. array of function pointers)
 - global variables
 - structured type safe control flow: <a href="#if">if</a> <a href="#if-else">if-else</a> <a href="#loop">loop</a> <a href="#jump-table">jump-table</a>

no-multi-return and undefined-order-of-evaluation was added to make the c target simpler.
an earlier compiler generated ssa temporaries for c to enforce the evaluation order which bloats the generated source and makes it hard to read.

intrinsic functions 
some builtins are represented as function calls instead of IR nodes and are implemented individually by the target language runtime:

 <span id="int32"     >int32(any)                   : convert to i32, e.g. int32(K) takes the lower 32 bit of a k value</span>
 <span id="int64"     >int64(any)                   : convert to i64, e.g. int64(K) is usually a non-op but required by Go's type system</span>
 <span id="float64"   >float64(any)                 : convert to f64</span>
 <span id="F64reinterpret_i64">F64reinterpret_i64(x): reinterpret/cast u64 value as f64</span>
 <span id="I64reinterpret_f64">I64reinterpret_f64(x): reinterpret/cast f64 value as i64</span>
 <span id="I8"        >I32(addr)                    : get byte as i32 from heap</span>
 <span id="I32"       >I32(addr)                    : get i32 from heap</span>
 <span id="I64"       >I64(addr)                    : get i64 from heap</span>
 <span id="F64"       >F64(addr)                    : get f64 from heap</span>
 <span id="SetI8"     >I32(addr)                    : set byte at addr to i32 low byte value</span>
 <span id="SetI32"    >SetI32(addr,value)           : set heap at addr to i32 value</span>
 <span id="SetI64"    >SetI64(addr,value)           : set heap at addr to i64 value</span>
 <span id="SetF64"    >SetF64(addr,value)           : set heap at addr to f64 value</span>
 <span id="I32clz"    >I32clz(x)                    : count leading zeros</span>
 <span id="Memory"    >Memory(blocks)               : declare the number of blocks of initial memory</span>
 <span id="Data"      >Data(offset,bytes)           : prepopulate initial memory at offset</span>
 <span id="Export"    >Export(functions..)          : declare exported functions, e.g. for wasm</span>
 <span id="Functions" >Functions(offset,functions..): add functions to function table at offset</span>
 <span id="Memorycopy">Memorycopy(d,s,n)            : e.g. memcpy in c, in wasm it needs bulk memory instructions</span>
 <span id="Memorysize">Memorysize()                 : returns the current number of blocks of the memory section</span>
 <span id="Memorygrow">Memorygrow(blocks)           : grows the memorysection by the given number of 64k blocks</span>
 <span id="Data"      >Data(off,string)             : initializes memory at offset</span>
 <span id="panic"     >panic(x)                     : traps on k errors, e.g. unreachable in wasm. it is up to the embedder to recover</span>

in the bootstrap interpreter k.go they are defined and imported from <a href="https://raw.githubusercontent.com/ktye/wg/master/module/module.go">wg/module/module.go</a>
`

const syscalls = `
the abstrace machine requires a set of simplified syscalls.
all arguments except for Native() are 32 bit integers:

 <span id="Args"      ><a href="#getargv">Args()</a>         :return number of command line arguments</span>
 <span id="Arg"       ><a href="#getargv">Arg(i,r)</a>       :call twice: first with r=0 returns the size of the argument, call again with allocated memory index in r</span>
 <span id="Read"      ><a href="#readfile">Read(f,n,d)</a>    :call twice: first with d=0 returns the size of the argument, call again with allocated memory index in r
                 n is the length of the filename and f it's memory index</span>
 <span id="Write"     ><a href="#writefile">Write(f,n,b,m)</a> :write content at b and length m to file with name at f with length n</span>
 <span id="ReadIn"    ><a href="#readfile">ReadIn(d,n)</a>    :read from stdin at most n bytes and write to d. return number of bytes written</span>
 <span id="Native"    ><a href="#cal">Native(x,y)</a>    :<a href="#extend">native</a> function call x,y are 64bit. the native function is registered at x (e.g. an index or a pointer) and
                 called with argument list y (a k value).</span>
 <span id="Exit"      ><a href="#repl">Exit(c)</a>        :exit with code c</span>
the reference implementation is imported from <a href="https://raw.githubusercontent.com/ktye/wg/master/module/system.go">wg/module/system.go</a>
the wasm version imports the system interface as an import module implemented in js which can be different for each application,
e.g. a Write call may write to an in-memory file system, trigger a download or use the file system api.
`

const embed = `
`

const extend = `
the k interpreter can be extended when <a href="#embed">embedded</a> to a host application assigning native functions to a k variable.
a native function is represented as a k value with type 14 that contains a list of two values: an identifier for
the host system and the string form.
the identifier may be an index to a table of registered native functions,
e.g. as a k-verb (type 0) to prevent refcounting, or a pointer disguised as a character vector.
native function have <b>fixed arity</b>. the arity is stored in the <a href="#types">length field</a>.
when called, the argument list has the length of the arity, otherwise a projection/error is triggered
before the native function call.

there is also a mild protection for k variables used in <a href="#z.k">z.k</a> which can be used in any k source:
monadic functions assigned to a symbol including the backtick, e.g. °f:{+/x} use the symbol with a dot attached internally.
these functions are called as overloads to @ e.g. with °f 2 3 4.
z.k defines °x(hex) °t(token) °p(parse) °c(as-chars) °i(as-ints) °s(as symbols) °f(as-floats) °z(as-complexs)
°(int) converts the numeric value of a basic verb or instruction to verb type(0). e.g. °@'1 2 3 is (:;+;-)
`

const f77 = `
compile with: gfortran <a href="https://github.com/ktye/i/releases/download/latest/k.f">k.f</a>
`

const cc = `
`

const wasm = `
`

const wasi = `
`

const tail = `</pre>
<pre id="mono">M 0 1 2 3 4 5 6 7 8 9 0</pre>
</body>
`
