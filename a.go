package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"sync"
)

var rd func() []c
var dr sync.Mutex

func main() {
	ini()
	table[21] = red
	table[40] = exi
	table[21+dyad] = wrt
	args, addr := os.Args[1:], ""
	if len(args) == 1 && args[0] == "-kwac" {
		inikwac()
		return
	} else if len(args) > 1 && args[0] == "-p" {
		addr = args[1]
		if _, o := atoi([]c(args[1])); o {
			addr = ":" + args[1]
		}
		args = args[2:]
	} else if len(args) > 0 && args[0] == "-u" {
		addr = ":2019"
		dec(evl(prs(mkb([]c(ui)))))
		args = args[1:]
	}
	if len(args) > 0 {
		defer stk(false)
		rd = read
		zx := mk(L, k(len(args))) // .z.x: args
		for i, a := range args {
			m.k[2+k(i)+zx] = mkb([]c(a))
		}
		asn(mks(".z.x"), inc(zx), mk(N, atom))
		lod(inc(m.k[2+zx]))
		dec(zx)
	}
	if addr != "" {
		go http.ListenAndServe(addr, http.HandlerFunc(srv))
	}
	rd = readline(bufio.NewScanner(os.Stdin)) // 0:` or 1:` read a single line in interactive mode
	for {
		try()
	}
}
func try() {
	defer stk(true)
	evp(red(wrt(mku(0), enl(mkc(' '))))) // r: 1: ("" 1: ," ")
}
func red(x k) (r k) { // 1:x
	t, n := typ(x)
	if t == S {
		x = str(x)
		t, n = typ(x)
	}
	if t != C {
		panic("type")
	}
	var b []c
	if n == 0 {
		b = rd()
		if b == nil {
			exi(mki(0))
		}
	} else {
		xp := 8 + x<<2
		p, err := ioutil.ReadFile(string(m.c[xp : xp+n]))
		if err != nil {
			panic(err)
		}
		b = p
	}
	n = k(len(b))
	r = mk(C, n)
	rp := 8 + r<<2
	copy(m.c[rp:rp+n], b)
	return r
}
func read() []c { // read all from stdin (non-interactive)
	b, err := ioutil.ReadAll(os.Stdin)
	if err == nil {
		return b
	}
	return []c{}
}
func readline(sc *bufio.Scanner) func() []c { // read single line (interactive)
	return func() []c {
		if sc.Scan() == false {
			return nil
		}
		return sc.Bytes()
	}
}
func wrt(x, y k) k { // x 1:y
	t, n := typ(x)
	if t == S {
		x = str(x)
		t, n = typ(x)
	}
	if t != C {
		panic("type")
	}
	if n != 0 {
		panic("nyi") // write to a file
	}
	t, n = typ(y)
	if t != C || n == atom {
		panic("type")
	}
	yp := 8 + y<<2
	w := bufio.NewWriter(os.Stdout)
	w.Write(m.c[yp : yp+n])
	w.Flush()
	return decr(y, x)
}
func exi(x k) (r k) { // exit built-in
	t, n := typ(x)
	if t == I && n == atom {
		os.Exit(int(m.k[2+x]))
	}
	os.Exit(1)
	return mk(N, atom)
}
func stk(hide bool) {
	if r := recover(); r != nil {
		a, b := stack(r)
		if hide { // interactive
			dec(asn(mks(".stk"), mkb([]byte(a)), mk(N, atom))) // stack trace: \s
		} else {
			println(a + "\n")
		}
		dec(wrt(mku(0), ano(m.k[srcp], mkb([]byte(b)))))
	}
}
func stack(c interface{}) (stk, err string) {
	h := false
	for _, s := range strings.Split(string(debug.Stack()), "\n") {
		if h && strings.HasPrefix(s, "\t") {
			if i := strings.Index(s, "/ktye/i/"); i > 0 {
				s = strings.TrimSpace(s[i+7:])
			}
			if len(s) > 0 {
				stk += "\n" + s
			}
		}
		if strings.Index(s, "panic.go") > 0 { // skip first lines
			h = true
		}
	}
	err = "?"
	if s, o := c.(string); o {
		err = s
	} else if e, o := c.(error); o {
		err = e.Error()
	}
	return stk, err
}
func srv(w http.ResponseWriter, r *http.Request) {
	dr.Lock()
	defer dr.Unlock()
	buf := bytes.NewBuffer(nil)
	defer func() {
		w.Write(buf.Bytes())
		r.Body.Close()
	}()
	defer func() {
		if rec := recover(); rec != nil {
			a, b := stack(rec)
			println(a)
			buf = bytes.NewBuffer(nil)
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(500)
			dec(wrt(mku(0), mkb([]byte(b))))
		}
	}()
	z, f := lupo(mku(0)), k(0)
	if z == 0 {
		return
	}
	z = atx(z, mks("Z"))
	if r.Method == "GET" { // .Z.G
		f = atx(z, mks("G"))
	} else if r.Method == "POST" { // .Z.P
		f = atx(z, mks("P"))
	} else {
		dec(z)
		return
	}
	if m.k[f]>>28 != N+1 {
		dec(f)
		return
	}
	hk, hv := mk(S, k(len(r.Header))), mk(L, k(len(r.Header)))
	kp, j := 8+hk<<2, k(0)
	for key := range r.Header {
		kv := key
		if len(kv) > 8 {
			kv = kv[:8]
		}
		mys(kp, btou([]c(kv)))
		m.k[2+j+hv] = mkb([]c(r.Header.Get(key)))
		kp, j = kp+8, j+1
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	l := mk(L, 3)                   // TODO? Query(dict)?
	m.k[2+l] = mkb([]c(r.URL.Path)) // ?
	m.k[3+l] = key(hk, hv)          // ?
	m.k[4+l] = mkb(b)               // ?
	y := cal(f, enlist(l))          // ? assume ( hdr(`a); body(`c) )
	t, n := typ(y)
	if t == I && n == atom { // error status code
		http.Error(w, "", int(decr(y, m.k[2+y])))
		return
	}
	if t == L && n == 2 && m.k[m.k[2+y]]>>28 == A { // (hdr;"body") or "body"
		hdr := fst(inc(y))
		keys, vals := str(inc(m.k[2+hdr])), m.k[3+hdr]
		for i := k(0); i < atm1(m.k[keys]&atom); i++ {
			v := str(atx(inc(vals), mki(i)))
			kp, vp := ptr(m.k[2+i+keys], C), ptr(v, C)
			kn, vn := m.k[2+i+keys]&atom, m.k[v]&atom
			w.Header().Set(string(m.c[kp:kp+kn]), string(m.c[vp:vp+vn]))
		}
		y = fst(drop(1, y))
		t, n = typ(y)
	}
	if t != C || n == atom {
		panic("type")
	} else if n > 0 {
		p := ptr(y, C)
		buf.Write(m.c[p : p+n])
	}
	dec(y)
}
func inikwac() { // write initial memory as data section
	skip := 0
	fmt.Printf("(0;0x")
	for i, c := range m.c {
		if c == 0 {
			skip++
		} else {
			if skip < 8 {
				for i := 0; i < skip; i++ {
					fmt.Printf("00")
				}
			} else if skip != 0 {
				fmt.Printf(";%d;0x", i)
			}
			fmt.Printf("%02x", c)
			skip = 0
		}
	}
	fmt.Println(")")
}
func pr(x k, a ...interface{}) {
	fmt.Printf(":%x ", x)
	r := kst(inc(x))
	_, n := typ(r)
	s := s(m.c[8+r<<2 : 8+n+r<<2])
	dec(r)
	fmt.Println(a, s)
}
func fatal(s string) { println(s); os.Exit(1) }

// ui webserver
// j: (js client program) sends requests to k for each event(mouse, key, size) and draws result to canvas
// handlers:
//  "/"                                     p(main page)
//  "/k.png"                                ico(favicon)
//  "/m,button,x0,y0,x1,y1,shift,alt,ctrl"  m(mouse event, buttons 0..4)
//  "/k,key,shift,alt,ctrl"                 k(key event)
//  "/s,width,height"                       s(size event)
const ui = ` \"ui:localhost:2019"
ico:0x89504e470d0a1a0a0000000d49484452000000100000001008060000001ff3ff61000000017352474200aece1ce90000000467414d410000b18f0bfc6105000000097048597300000ec400000ec401952b0e1b0000006349444154384fad93e10a8020108377bdff3b9793a428dc16f681723fdcc71d6aed0d902ae02c2db7b35bdf177809aad9b952feefe02b91408d650523382eeb8914b830990a9230911db8308946504c05d70bd7926804259102e22456409464f13b03070b7f28230cf1c9ad0000000049454e44ae426082
j:"w=window;d=document;b=d.body;b.style.margin=0;b.style.padding=0;b.style.overflow='hidden';N=Number;function pd(e){if(e)e.preventDefault();e.stopPropagation()};"
j,:"c=d.createElement('canvas');b.appendChild(c);ctx=c.getContext('2d');"
j,:"function draw(s){ctx.putImageData(new ImageData(new Uint8ClampedArray(s), w.innerWidth),0,0)};" /TODO: send w in header
j,:"function debounce(f,w,i){var t;return function e(){var c=this;var a=arguments;var l=function(){t=null;if(!i)f.apply(c,a);};var n=i&&!t;clearTimeout(t);t=setTimeout(l,w);if(n)f.apply(c,a);}}"
j,:"function get(p,f){var r = new XMLHttpRequest();r.responseType='arraybuffer';r.onreadystatechange=function(){if(this.readyState==4&&this.status == 200){if(f)f(this.response,this);}};r.open('GET',p);r.send()};"
j,:"function mod(e){return ','+[N(e.shiftKey),N(e.altKey),N(e.ctrlKey)]};"
j,:"xd=0;yd=0;down=function(e){xd=e.clientX;yd=e.clientY;pd(e)};nomenu=function(e){pd(e)};"
j,:"up=function(e){pd(e);get('m,'+[e.button,xd,e.clientX,yd,e.clientY]+mod(e),draw)};" /bs tab ret esc   delete     page;end;home;arrows->14-21
j,:"keycode=function(e){var k=e.keyCode;return (e.key.length==1)?e.key.charCodeAt():(k==8)?8:(k==9)?9:(k==13)?13:(k==27)?27:(k==46)?127:(k>32&&k<41)?k-19:null};"
j,:"key=function(e){var k=keycode(e);if(!k)return;get('k,'+k+mod(e),draw);pd(e);};"
j,:"wheel=function(e){var x=e.clientX;var y=e.clientY;var m=(e.deltaY>0)?4:(e.deltaY<0)?5:null;if(m)get('m,'+m+','+[x,y,x,y]+mod(e),draw)};"
j,:"size=function(e){c.width=w.innerWidth;c.height=w.innerHeight;get('s,'+[c.width,c.height],draw);pd(e)};"
j,:"function ae(x,y,z){x.addEventListener(y,z)};ae(w,'contextmenu',nomenu);ae(w,'mousedown',down);ae(w,'mouseup',up);ae(w,'wheel',wheel);ae(w,'keydown',key);ae(w,'resize',debounce(size,100));size()"
p:"<!DOCTYPE html>\n<html><head><link rel='icon' type='image/png' href='k.png'></head><body><script>",j,"</script></body></html>"
(w;h;d):(0;0;!0) /d pixel buffer w*h
size:{d::(x*y) #0;w::x;h::y;d}
opaque:255*256*256*256
flush:{8_` + "`" + `@opaque+x} /send pixel buffer(uint32 array) as binary data
(fw;fh):16 32 /fontsize
(dx;dy):0 0   / cursor offset(pixels)
key:{y;d[ \dst[dx;dy;chars[x]]]::255*256;dx+::16;$[x=13;(dx;dy)::(0;dy+32);];d}
mouse:{(w*h)#x+y}
.Z.G:{ \x
 x:","\:*x;u:*x;a:` + "`" + `i$'1_x
 $[(,"/")  ~u;p
   "/k.png"~u;((,"Content-type")!,"image/png";ico)
   "/s"~u;flush size[a 0;a 1]
   "/k"~u;flush key[a 0;a 1 2 3]
   "/m"~u;flush mouse[a 1;a 2]
  404]}
. 1:"u/f/f3.k"             /load font(works only when running in this directory)
font:{&,/(8#2)\:'0+x}'font /unpack
x:"0123456789ABCDEF:+-*%&|<>=!~,^#_$?@.0123456789'/\\;` + "`" + `\"(){}[]abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
chars:256#,!0
chars[x+0]:font[!#x]
dst:{o:(x+y*w);o+(16\z)+w*16/z}
`
