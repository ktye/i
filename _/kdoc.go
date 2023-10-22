package main

// this file builds ktye.github.io/kdoc.html
// $cd ..;go run _/kdoc.go

import (
	"fmt"
	"html"
	"os"
	"regexp"
	"strings"
)

func main() {
	O := fmt.Print
	H := html.EscapeString
	R := func(x string) string { b, _ := os.ReadFile(x); return string(b) }
	h := func(x, y, z string) { O(`<hr/><span id="`+x+`">`+y+z+`</span>`) }

	O(head)

	//refcard
	x := R("readme")
	a := strings.Fields("flp add neg sub fst mul sqr div til key wer min rev max asc les dsc mor grp eql not mtc enl cat srt cut cnt tak flr drp str cst unq fnd typ atx val cal abs sin cos exp log find angle imag conj types")
	b := strings.Fields("Flp Add Neg Sub Fst Mul Sqr Div Til Key Wer Min Rev Max Asc Les Dsc Mor Grp Eql Not Mtc Enl Cat Srt Cut Cnt Tak Flr Drp Str Cst Unq Fnd Typ Atx Val Cal Abs Sin Cos Exp Log Find Angle Imag Conj K    ")
	for i := range a {
		x = strings.Replace(x, a[i], `<a href="#`+b[i]+`">`+a[i]+`</a>`, 1)
	}
	O(x)

	O(tty)

	//toc todo links
	x = toc
	a  = strings.Fields("kinit z.k tree/ir abstract-machine")
	b  = strings.Fields("kinit z.k ir      mach            ")
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
	L := func(s string) {
		r := regexp.MustCompile(`\w+\(`)
		v := r.FindAllIndex([]byte(s), -1)
		i := 0
		for _, ab := range v {
			a, b := ab[0], ab[1]-1
			O(H(s[i:a]))
			O(`<a href="#` + s[a:b] + `">` + s[a:b] + `</a>`)
			i = b
		}
		O(H(s[i:]))
	}
	F := func(s string) int {
		i := strings.Index(s, "(")
		O(`<span id="` + s[6:i] + `">`)
		O(H(s[:i]))
		j := 1 + strings.Index(s[1:], "\n")
		if strings.Index(s, "\n\t") == j { //multi line function body
			j = 2 + strings.Index(s, "\n}")
		}
		L(s[i:j])
		O("</span>")
		return j
	}
	S := ""
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

	h("z.k", "<b>z.k</b> is placed in the initial memory section at 600 by <a href='#zk'>zk</a> and executed by <a href='#kinit'>kinit</a>:\n", H(R("z.k")))
	h("ir", "intermediate representation", ir)
	h("mach", `abstract machine model`, mach)

	O(tail)
}

const head = `<!KDOCTYPE html>
<head><meta charset="utf-8">
<link rel="icon" type="image/png" sizes="32x32" href="db/kelas32.png">
<link rel="icon" type="image/png" sizes="16x16" href="db/kelas16.png">
<title>kdoc</title>
<style>
*{font-family:monospace}
.q{color:blue}.q:hover{cursor:pointer}
.h{font-weight:bold}
body{display:flex;flex-direction:column;height:100vh;overflow:hidden}
a{text-decoration:none}
#ref{display:inline-block;background:#ffe}
#tty{display:inline-block;background:#004687;color:white;width:100%;margin-left:1em;margin-right:0.5em;outline:none;overflow:hidden}
#top{display:flex}
#src{overflow:auto;outline:none;margin:0}
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
const td_=new TextDecoder("utf-8"),su=x=>td_.decode(x),te_=new TextEncoder("utf-8"),us=x=>te_.encode(x)
let cur=null
function init(){
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

</script>

<div id="top">
<pre id="ref">`

const tty = `</pre>
<pre id="tty" contenteditable spellcheck="false"></pre>
<pre id="toc">`

const toc = `kinit . . . . . . . . . . .z.k

intro                         

adverb

control flow

examples algebra stats queries

html/js/wasm         embedding







implementation
types        heap    allocator
symbols    variables     scope
exec    kvm    instruction set

compilers boot  f77 c go  wasm
tree/ir       abstract-machine
`

const src = `</pre>
</div>
<pre id="src">`

const ir = `
 ...
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
are represented as function calls and are implemented individually by the target language runtime.

 <span id="int32"     >int32(any)        : convert to i32, e.g. int32(K) takes the lower 32 bit of a k value</span>
 <span id="int64"     >int64(any)        : convert to i64, e.g. int64(K) is usually a non-op but required by Go's type system</span>
 <span id="float64"   >float64(any)      : convert to f64</span>
 <span id="I32"       >I32(addr)         : get i32 from heap</span>
 <span id="I64"       >I64(addr)         : get i64 from heap</span>
 <span id="F64"       >F64(addr)         : get f64 from heap</span>
 <span id="SetI32"    >SetI32(addr,value): set heap at addr to i32 value</span>
 <span id="SetI64"    >SetI64(addr,value): set heap at addr to i64 value</span>
 <span id="SetF64"    >SetF64(addr,value): set heap at addr to f64 value</span>
 <span id="Memorycopy">Memorycopy(d,s,n) : e.g. memcpy in c, in wasm it needs bulk memory instructions</span>
 <span id="Data"      >Data(off,string)  : initializes memory at offset</span>
 <span id="panic"     >panic(x)          : traps on k errors, e.g. unreachable in wasm. it is up to the embedder to recover</span>
 ...
`

const tail = `</pre>
<pre id="mono">M 0 1 2 3 4 5 6 7 8 9 0</pre>
</body>
`
