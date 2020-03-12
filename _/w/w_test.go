package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

var broken = false // ../../k.w

func TestB(t *testing.T) {
	testCases := []struct {
		sig string
		b   string
		e   string
	}{
		{"I:IF", "x::y;x", "2000 2001 390300 2000"},
		{"I:II", "y+I?C x", "2001 2000 2d0000 6a"},
		{"I:II", "y+I?C x", "2001 2000 2d0000 6a"},
		{"I:I", "x?[;x:4;x:6];x", "024002400240024020000e020001020b0c010b410421000c010b410621000c000b2000"},
		{"I:I", "x?[x:4;;x:6];x", "024002400240024020000e020001020b410421000c020b0c000b410621000c000b2000"},
		{"0:I", "x::5", "2000 4105 360200"},
		{"0:I", "x::C 5", "2000 4105 2d0000 3a0000"},
		{"0:I", "0::1130366807310592j", "4100 42 8082 90c0 8082 8102 370300"},
		{"I:I", "x?!;x", "20000440 00 0b 2000"},
		{"I:I", "-1+x", "417f 2000 6a"},
		{"I:I", "x-1", "200041016b"},
		{"I:II", "x:I 4+x;x", "4104 2000 6a 280200 2100 2000"},
		{"I:I", "x?[x:4;x:5;x:6];x", "024002400240024020000e020001020b410421000c020b410521000c010b410621000c000b2000"},
		{"I:I", "I?255j&1130366807310592j>>J?8*x", "42ff0142808290c080828102410820006cad8883a7"},
		{"I:I", "(x<6)?/x+:1;x", "0240 0340 2000 4106 49 45 0d01 20004101 6a 2100 0c00 0b0b2000"},
		{"I:I", "(x<6)?/(x+:1;x+:1);x", "0240 0340 2000 4106 49 45 0d01 200041016a2100 200041016a2100 0c00 0b0b2000"},
		{"I:I", "1/(x+:1;?x>5);x", "0240 0340 2000 4101 6a 2100 2000 4105 4b  0d01 0c00 0b0b2000"},
		{"I:III", "$[x;y;z]", "2000 047f 2001 05 2002 0b"},
		{"I:I", "(x>3)?(:-x);x", "2000 4103 4b 0440 4100 2000 6b 0f 0b 2000"},
		{"I:I", "(x>3)?x+:1;x", "2000 4103 4b 0440 2000 4101 6a 2100 0b 2000"},
		{"I:II", "x::y;I x", "2000 2001 360200 2000 280200"},
		{"I:I", "x/r:r+i;r", "20000440410021020340200120026a2101200241016a22022000490d000b0b2001"},
		{"I:I", "x/r+:i;r", "20000440410021020340200120026a2101200241016a22022000490d000b0b2001"},
		{"I:II", "x+y", "20002001 6a"},
		{"I:II", "x\\y", "20002001 70"},
		{"I:II", "r:x;r+:y;r", "2000 2102 2002 2001 6a 2102 2002"},
		{"I:I", "x/r:i;r", "2000044041002102034020022101200241016a22022000490d000b0b2001"},
		{"I:II", "(3+x)*y", "4103 2000 6a 2001 6c"},
		{"I:I", "1+x", "410120006a"},
		{"F:FF", "(x*y)", "20002001 a2"},
		{"F:FF", "x-y", "20002001 a1"},
		{"F:FF", "3.*x+y", "44 0000000000000840 20002001 a0 a2"},
		{"I:I", "x:1+x;x*2", "4101 2000 6a 2100 2000 4102 6c"},
	}
	for n, tc := range testCases {
		f := newfn(tc.sig, tc.b)
		e := f.parse(nil, nil, nil)
		b := string(hex(e.bytes()))
		s := trim(tc.e)
		if b != s {
			t.Fatalf("%d: expected/got:\n%s\n%s", n+1, s, b)
		}
		// fmt.Println(b)
		ctest(t, tc.sig, tc.b)
	}
}
func TestRun(t *testing.T) {
	m, data := run(strings.NewReader("add:I:II{x+y}/cnt\n/\n/sum:I:I{x/r+:i;r}\n/"))
	g := s(hex(m.wasm(data)))
	e := "0061736d0100000001070160027f7f017f030201000503010001070d02036d656d02000361646400000a09010700200020016a0b"
	if e != g {
		t.Fatalf("expected/got\n%s\n%s\n", e, g)
	}
}
func ctest(t *testing.T, sig, b s) {
	b = jn("f:", sig, "{", b, "}")
	m, data := run(strings.NewReader(b))
	out := m.cout(data)
	if len(out) == 0 {
		t.Fatal("no output")
	}
	//fmt.Println(string(out))
}
func hex(a []c) []c {
	var r bytes.Buffer
	for _, b := range a {
		hi, lo := hxb(b)
		r.WriteByte(hi)
		r.WriteByte(lo)
	}
	return r.Bytes()
}
func newfn(sig string, body string) fn {
	var buf bytes.Buffer
	buf.WriteString(body)
	buf.WriteByte('}')
	v := strings.Split(sig, ":")
	if len(v) != 2 {
		panic("signature")
	}
	f := fn{src: [2]int{1, 0}, Buffer: buf}
	f.t = typs[v[0][0]]
	for _, c := range v[1] {
		f.locl = append(f.locl, typs[byte(c)])
	}
	f.args = len(v[1])
	return f
}
func trim(s string) string { return strings.Replace(s, " ", "", -1) }
func TestHtml(t *testing.T) { // write k.html from ../../k.w
	if broken {
		t.Skip()
	}
	m, data, src, err := KWasmModule()
	if err != nil {
		t.Fatal(err)
	}
	if idx := bytes.Index(src, []byte{'\n', '\\'}); idx != -1 {
		src = src[:idx+1]
	}
	tests, err := ioutil.ReadFile("t")
	if err != nil {
		t.Fatal(err)
	}
	wasm := m.wasm(data)
	_, exp := m.exports()
	var txt, fns bytes.Buffer
	fmt.Fprintf(&txt, "kwasm(%d b) %s tests src", len(wasm), time.Now().Format("2006.01.02"))
	for _, f := range exp {
		if f.t != 0 && (f.args == 1 || f.args == 2) {
			fmt.Fprintf(&txt, " %s", f.name)
			if fns.Len() != 0 {
				fns.WriteByte(',')
			}
			fmt.Fprintf(&fns, "\"%s\":%d", f.name, f.args)
		}
	}
	txt.WriteString(`\n `)
	var b bytes.Buffer
	s := hh
	s = strings.Replace(s, `{{wasm}}`, base64.StdEncoding.EncodeToString(wasm), 1)
	s = strings.Replace(s, `{{tests}}`, base64.StdEncoding.EncodeToString(tests), 1)
	s = strings.Replace(s, `{{src}}`, base64.StdEncoding.EncodeToString(src), 1)
	s = strings.Replace(s, `{{cons}}`, string(txt.Bytes()), 1)
	s = strings.Replace(s, `{{fncs}}`, string(fns.Bytes()), 1)
	b.WriteString(s)
	if e := ioutil.WriteFile("k.html", b.Bytes(), 0644); e != nil {
		t.Fatal(e)
	}
}

func KWasmModule() (module, []byte, []byte, error) {
	var src io.Reader
	var srcb []byte
	if k, e := ioutil.ReadFile("../../k.w"); e != nil {
		return nil, nil, nil, e
	} else {
		src = bytes.NewReader(k)
		srcb = k
	}
	m, data := run(src)
	return m, data, srcb, nil
}
func TestCout(t *testing.T) { // write k_c from ../../k.w
	if broken {
		t.Skip()
	}
	m, data, _, err := KWasmModule()
	if err != nil {
		t.Fatal(err)
	}
	var dst bytes.Buffer
	io.Copy(&dst, strings.NewReader(kh))
	dst.Write(m.cout(data))
	io.Copy(&dst, strings.NewReader(kt1))
	for _, f := range m {
		if f.args == 1 && f.t == I && f.locl[0] == I && f.name != "ini" && f.name != "mki" {
			s := "\t\t} else if (Match(\"" + f.name + "\", a)) { n = f1(" + f.name + ", stack, n);\n"
			dst.WriteString(s)
		} else if f.args == 2 && f.t == I && f.locl[0] == I && f.locl[1] == I {
			s := "\t\t} else if (Match(\"" + f.name + "\", a)) { n = f2(" + f.name + ", stack, n);\n"
			dst.WriteString(s)
		}
	}
	io.Copy(&dst, strings.NewReader(kt2))
	if e := ioutil.WriteFile("k_c", dst.Bytes(), 0744); e != nil {
		t.Fatal(e)
	}
}

const hh = `<html>
<head><meta charset="utf-8"><title>w</title></head>
<link rel='icon' type'image/png' href="k32.png">
<style>
 html,body,textarea{height:100%;margin:0;padding:0;overflow:hidden;}
 #kons{top:0;left:0;width:100%;height:100%;position:absolute;border:0pt;resize:none;font-family:monospace;overflow:auto;}
 .term{background:black;color:white}
 .hold{background:white;color:black}
 .edit{background:#ffffea;color:black}
 #cnv{width:100;height:100;top:0;right:0;position:absolute;}
 #dl{display:none;}
</style>
<body>
<textarea id="kons" class="term" wrap="off" autofocus spellcheck="false"></textarea>
<canvas id="cnv"></canvas>
<script>
var r = "{{wasm}}"
var rt = "{{tests}}"
var rs = "{{src}}"
function sa(s){var r=new Uint8Array(new ArrayBuffer(s.length));for(var i=0;i<s.length;i++)r[i]=s.charCodeAt(i);return r}
function pd(e){if(e){e.preventDefault();e.stopPropagation()}};
function ae(x,y,z){x.addEventListener(y,z)};
var kwasm = sa(atob(r))
var tests = atob(rt)
var src = atob(rs)
var K
// kons (k console)
var hit = kons
var konstore = ""
var edname = ""
var ed = false
function initKons() {
 kons.value = "{{cons}}"
 var hold = false
 kons.onkeydown = function(e) {
  if(e.which === 27) { // quit edit / toggle hold / close image
   delay = 0
   pd(e)
   if (ed) { qed(); hold=true }
   hold = !hold
   kons.className = (hold) ? "hold" : "term"
   imgSize(0, 0)
   hit = kons
  } else if (e.which === 13 && !hold && !ed) { // execute
   pd(e)
   var a = kons.selectionStart
   var b = kons.selectionEnd
   var s = kons.value.substring(a, b)
   if (b == a) {
    if (kons.value[a] == "\n") a -= 1
    a = kons.value.lastIndexOf("\n", a)
    if (a == -1) a = 0
    b = kons.value.indexOf("\n", b)
    if (b == -1) b = kons.value.length
    s = kons.value.substring(a, b)
   }
   if (kons.selectionEnd != kons.value.length) O(s)
   O("\n")
   s = s.trim()
   if (s === "tests")           { O(tests);                        return }
   if (s === "src")             { O(src);                          return }
   if (s === "\\c")             { kons.value=" ";imgSize(0, 0);    return }
   if (s === "\\h")             { O(atob(h));P();                  return }
   if (s.substr(0,2) === "\\e") { P();edit(s.substr(2));           return }
   if (s.substr(0,2) === "\\w") { download(s.substr(2).trim());P();return }
   if (s.substr(0,2) === "\\L") { P();loop(s.substr(2).trim());    return }
   if (s === "\\lf")            { s = "\\m #:'.fs"                        }
   hash(s);E(s);P()
  }
 }
 kons.onmousedown = function(e) { hit=kons; if(e.button==2)pd(e); }
 kons.onblur  = function(e) { kons.style.filter = "brightness(70%)" }
 kons.onfocus = function(e) { kons.style.filter = "brightness(100%)" }
}
function O(s) { kons.value += s; kons.scrollTo(0, kons.scrollHeight) }
function P()  { kons.value += " " }

function us(s) { return new TextEncoder("utf-8").encode(s) } // uint8array from string
function su(u) { return (u.length) ? new TextDecoder("utf-8").decode(u) : "" }
function kst(x) {
 var h = K.I[x>>>2]
 var t = (h>>>29)>>>0
 var n = (h&536870911)>>>0
 var o = []
 switch(t){
 case 1:
  return '"'+su(K.C.slice(8+x, 8+x+n))+'"'
 case 2:
  x >>>= 2
  return K.I.slice(2+x, 2+x+n).join(" ").split("2147483648").join("0N")
 case 3:
  x >>>= 3
  var s = K.F.slice(1+x, 1+x+n).join(" ")
  if(s.indexOf(".")==-1) s+="f"
  return s
 case 4:
  x >>>= 2
  var r = ""
  var v = K.I.slice(2+x, 2+x+n)
  var tr = function(s) { return s.substr(1, s.length-2) }
  for (var i=0; i<n; i++) r += String.fromCharCode(96) + tr(kst(v[i]))
  return r
 case 5:
  x >>>= 2
  var r = []
  for (var i=0; i<n; i++) r.push(kst(K.I[2+x+i]))
  return "("+r.join(";")+")"
 default:
  return "kst nyi: t=" + String(t)
 }
}

var funcs = {{{fncs}}}
function parseNoun(s) {
 var t = 2
 if(s.startsWith('"'))                     return chrVector(s);
 if(s.startsWith(String.fromCharCode(96))) return symVector(s);
 if(s.startsWith('('))                     return lstVector(s.substr(1));
 if(s.indexOf(".") != -1) t = 3
 var v = s.split(",").map(x=>Number(x))
 var n = v.length
 if (n==0) return 0;
 var x = K.exports.mk(t, n)
 if (t==2) for (var i=0;i<n;i++) K.I[2+i+(x>>>2)] = v[i];
 else      for (var i=0;i<n;i++) K.F[1+i+(x>>>3)] = v[i];
 return x
}
function chrVector(s) {
 s = s.substr(1,s.length-2)
 var n = s.length
 var x = K.exports.mk(1, n)
 for (var i=0;i<n;i++) K.C[8+x+i] = s.charCodeAt(i);
 return x
}
function symVector(s) {
 var v = s.substr(1).split(String.fromCharCode(96))
 var x = K.exports.mk(4, v.length)
 for (var i=0; i<v.length; i++) K.I[2+i+(x>>>2)] = chrVector("_"+v[i]+"_")
 return x
}
function lstVector(s) {
 if (s.length == 0 || s.endsWith(')') == false) { throw new Error("parse list: "+s) }
 if (s.length == 1) return K.exports.mk(5, 0);
 s = s.substr(0, s.length-1)
 var a = 0
 var l = 0
 var r = []
 for (var i=0; i<s.length; i++) {
  var c = s[i]
  if (c == '(') l++
  else if (c == ')') { l--; if(l<0) throw new Error("list)") }
  else if ((l==0) && (c==';')) {
   r.push(parseNoun(s.substring(a,i)))
   a = i + 1
  }
 }
 r.push(parseNoun(s.substr(a)))
 var x = K.exports.mk(5, r.length)
 for (var i=0; i<r.length; i++) K.I[2+i+(x>>>2)] = r[i]
 return x
}
function xx(x) { return x.toString(16).padStart(8,'0') }
function dump(x, n) {
 var p = x >>> 2
 O("\n"+xx(p)+" ")
 for (var i=0; i<n; i++) {
  O(" "+xx(K.I[p+i]))
  if ((i>0)&&((i+1)%8==0)) O("\n"+xx(x+4*i+4)+" ")
  else if ((i>0)&&((i+1)%4==0)) O(" ")
 }
 O("\n")
}
function E(s) {
 try{ // todo save/restore
  var stack = []
  var v = s.split(" ").filter(x => x)
  for (var i=0; i<v.length; i++) {
   s = v[i]
   if (s=="dump") { dump(0,100); continue; }
   if(s in funcs) {
    var n = funcs[s]
    x = stack.pop()
    if(n==2) {
     y = x
     x = stack.pop()
     stack.push(K.exports[s](x, y))
    } else {
     stack.push(K.exports[s](x))
    }
    continue
   }
   var x = parseNoun(s)
   if (x!=0) { stack.push(x); continue; }
   else      { throw new Error(s); }
  }
  if(stack.length != 1) { throw new Error("stack length "+String(stack.length)) }
  else {
   x = stack.pop()
   O(kst(x)+"\n")
   K.exports.dx(x)
  }
 } catch(e) {
  console.log(e)
  O("error")
 }
}

function edit(name) { 
 edname = name; ed = true; konstore = kons.value; 
 var u = getfile(name.trim())
 kons.value = (u.length>0) ? su(u) : ""
 kons.className = "edit"
 kons.scrollTo(0,0);
}
function qed() { // quit edit
 putfile(edname, us(kons.value))
 kons.value = konstore; kons.scrollTo(0, kons.scrollHeight)
 ed = false
}
ae(kons,"contextmenu", function(e) { // button-3 search
 var l = kons.selectionEnd-kons.selectionStart; if(e.button==2&l>0) {
  pd(e); var t = kons.value.substring(kons.selectionStart,kons.selectionEnd)
  var f = function(a){ return kons.value.indexOf(t,a) }; var n = f(kons.selectionEnd)
  if (n<0){n=f(0)}; kons.setSelectionRange(n,n+l); }
})
function hash(s){window.location.hash=encodeURIComponent(s.trim())}

(async () => {
 initKons()
 const module = await WebAssembly.compile(kwasm.buffer);
 K = await WebAssembly.instantiate(module);
 K.C = new Uint8Array(K.exports.mem.buffer)
 K.I = new Uint32Array(K.exports.mem.buffer)
 K.F = new Float64Array(K.exports.mem.buffer)
 K.exports.ini(16);
 var h = decodeURIComponent(window.location.hash.substr(1))
 window.location.hash = h
 if (h.length > 0) {
  var p = kons.value.length
  kons.value += h
  kons.setSelectionRange(p, kons.value.length)
 }
 kons.focus()
})();
</script></body></html>
`

const kh = `#include<stdlib.h>
#include<stdio.h>
#include<stddef.h>
#include<malloc.h>
#include<string.h>
#include<math.h>
#define R return
typedef void V;typedef char C;typedef int32_t I;typedef int64_t J;typedef double F;typedef uint32_t uI;typedef uint64_t uJ;
I __builtin_clz(I x){I r;__asm__("bsr %1, %0" : "=r" (r) : "rm" (x) : "cc");R r^31;}
V trap() { exit(1); }
C *MC;I* MI;J* MJ;F *MF;
//F NaN = &((unt64_t)9221120237041090561ull);
`
const kt1 = `// Postfix test interface: e.g. 5 mki til rev fst 0 500 dump
const trace = 0;
I pop1(I *s, I n, I *x) {
	*x = s[n-1];
	return n-1;
}
I pop2(I *s, I n, I *x, I *y) {
	*x = s[n-2];
	*y = s[n-1];
	return n-2;
}
I push(I *s, I n, I x) {
	s[n] = x;
	return n+1;
}
I f1(I (*f)(I), I *s, I n) {
	if(trace) printf("%d: ", s[n-1]);
	s[n-1] = f(s[n-1]);
	if(trace) printf("%d(%x)\n", s[n-1], s[n-1]);
	return n;
}
I f2(I (*f)(I,I), I *s, I n) {
	if(trace) printf("%d %d: ", s[n-2], s[n-1]);
	s[n-2] = f(s[n-2], s[n-1]);
	if(trace) printf("%d\n", s[n-2]);
	return n-1;
}
I Match(C *a, C *b) {
	for (I i=0; ;i++) {
		if (a[i] != b[i]) R 0;
		if (a[i] == 0)    R 1;
	}
}
V Dump(I x, I n) {
	I p = x>>2;
	printf("\n%08x  ", x);
	for (I i=0; i<n; i++) {
		printf(" %08x", MI[p+i]);
		if ((i > x) && ((i+1)%8 == 0))      printf("\n%08x  ", x+4*i+4);
		else if ((i > 0) && ((i+1)%4 == 0)) printf(" ");
	}
	printf("\n");
}
V kst(I x) {
	I i, j, y, m, tof;
	I t = ((uI)MI[x>>2])>>29;
	I n = ((uI)MI[x>>2])&536870911;
	switch(t){
	case 1:
		printf("\"");
		for(i=0;i<n;i++) printf("%c", MC[8+x+i]);
		printf("\"");
		break;
	case 2:
		x = 2 + (x>>2);
		for(i=0;i<n;i++) {
			if (i>0)  printf(" ");
			if (MI[x+i] == -2147483648) printf("0N");
			else                        printf("%d", MI[x+i]);
		}
		break;
	case 3:
		tof = 1;
		x = 1 + (x>>3);
		for(i=0;i<n;i++) {
			if (i>0)  printf(" ");
			if (MF[x+i] != MF[x+i]) { printf("0n"); tof = 0; } else printf("%g", MF[x+i]);
			if (MF[x+i] != (F)(I)MF[x+i]) tof = 0;
		}
		if(tof) printf("f");
		break;
	case 4:
		x = 2 + (x>>2);
		for(i=0;i<n;i++) {
			printf("%c", 96);
			y = MI[x+i];
			m = MI[y>>2]&536870911;
			for(j=0;j<m;j++) printf("%c", MC[8+y+j]);
		}
		break;
	case 5:
		x = 2 + (x>>2);
		printf("(");
		for(i=0;i<n;i++) {
			if(i>0) printf(";");
			kst(MI[x+i]);
		}
		printf(")");
		break;
	default:
		printf("nyi: kst %x t=%d\n", x, t);trap();
	}
}
V O(I x) { kst(x); printf("\n"); }
I parseNoun(C *s);
I findChr(C *s, I n, C p) {
	I i;
	for (i=0; i<n; i++) if (s[i] == p) return i;
	R n;
}
I chrVector(C *s) {
	I i, x;
	I n = strlen(s);
	s++; n-=2;
	if(n<0) trap();
	x = mk(1, n);
	for (i=0; i<n; i++) MC[8+x+i] = s[i];
	R x;
}
I lstVector(C *s) {
	I r[9];
	I l = 0;
	I a = 0;
	I rn = 0;
	I i;
	C p;
	I n = strlen(s);
	if ((n == 0) || s[n-1] != ')') { printf("parse list: %s\n", s); trap(); }
	if (n == 1) return mk(5, 0);
	s[n-1] = 0; n--;
	for (i=0; i<n; i++) {
		p = s[i];
		if (p=='(') l++;
		else if (p==')') {
			l--;
			if (l<0) { printf(")"); trap(); }
		} else if (l==0 && p == ';') {
			s[i] = 0;
			r[rn++] = parseNoun(s+a);
			if (r[rn-1] == 0) { printf("null element"); trap(); }
			a = i + 1;
			if (rn==8) { printf("list limit"); trap(); }
		}
	}
	r[rn++] = parseNoun(s+a);
	l = mk(5, rn);
	for (i=0;i<rn;i++) MI[2+(l>>2)+i] = r[i];
	R l;
}
I symVector(C *s) {
	I i, j, x, y, m;
	I n = strlen(s);
	I ns = 0;
	C *p;
	for (i=0;i<n;i++) if(s[i] == 96) ns++;
	x = mk(4, ns);
	s++;n--;
	for (i=0;i<ns;i++) {
		m = findChr(s, n, 96);
		y = mk(1, m);
		for (j=0;j<m;j++) MC[8+y+j] = s[j];
		MI[2+i+(x>>2)] = y;
		s += m+1; n -= m+1;
	}
	R x;
}
I numVector(C *s) {
	I n = strlen(s);
	I x, i;
	F fv[8];
	I iv[8];
	C *p;
	I isf = (strchr(s, '.') != NULL);
	if (s[0] == ',')      { s++;      n--; }
	if (n>0&&s[n-1]==',') { s[n-1]=0; n--; }
	p = strtok(s, ",");
	n = 0;
	while(p != NULL) {
		if (n==8) break;
		if (isf)  fv[n++] = atof(p);
		else      iv[n++] = atoi(p);
 		p = strtok(NULL, ",");
	}
	if (isf) {
		x = mk(3, n);
		for (i=0; i<n; i++) MF[i+1+(x>>3)] = fv[i];
	} else {
		x = mk(2, n);
		for (i=0; i<n; i++) MI[i+2+(x>>2)] = iv[i];
	}
	R x;
}
I parseNoun(C *s) {
	C c = s[0];
	if (c == '"')                return chrVector(s);
	if (c == 96)                 return symVector(s);
	if (c >= '0' && s[0] <= '9') return numVector(s);
	if (c == '(')                return lstVector(s+1);
	if (strchr(s, ',') != NULL)  return numVector(s);
	R 0;
}
#define M0 16
I main(int args, C **argv){
	MC=malloc(1<<M0);MI=(I*)MC;MJ=(J*)MC;MF=(F*)MC;
	memset(MC, 0, 1<<M0);
	I stack[32];
	I i, n = 0;
	I x, y, r;
	C *a;
	ini(M0);
	for (i=1; i<args; i++) {
		a = argv[i];
		x = parseNoun(a);
		if (x != 0) {
			n = push(stack, n, x);
			continue;
		}
		//printf("%s ", argv[i]);
		if (Match("dump", a)) {
			Dump(0, 100);
`
const kt2 = `
		} else {
			printf("arg!");
			trap();
		}
	}
	if (n != 1) { printf("stack (%d)", n);trap(); }
	O(stack[0]);
}
`
