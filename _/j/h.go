// +build ignore

// generate j.html
package main

import (
	"bytes"
	"fmt"
	"html"
	"io/ioutil"
	"strings"
)

func main() {
	var def []struct{ name, value string }
	var o bytes.Buffer
	o.WriteString(head)

	t, e := ioutil.ReadFile("readme")
	fatal(e)
	fmt.Fprintf(&o, "<pre id=\"ref\">\n%s</pre>\n", hs(string(t)))

	t, e = ioutil.ReadFile("j.w")
	fatal(e)
	fmt.Fprintf(&o, "<h1 id=\"src\">j.w source of the wasm module</h1> <span style='color:blue' onclick=\"toggle('jw')\">(show/hide)</span><br/>\n<pre id='jw' style='display:none'>\n%s</pre>\n", hs(string(t)))

	b, e := ioutil.ReadFile("t")
	fatal(e)
	if idx := bytes.Index(b, []byte("\n\\\n")); idx > 0 {
		b = b[:idx+1]
	}
	b = append(b, canvas()...)

	toc := []string{"ref", "src"}
	var te *bytes.Buffer
	nl := []byte("\n")
	v := bytes.Split(b, nl)
	var id string
	for _, b := range v {
		s := string(b)
		if strings.HasPrefix(s, "(-") {
			s = btrim(s)[1:]
			id = s
			fmt.Fprintf(&o, "<h1 id=\"%s\">%s</h1>\n ", s, s)
			fmt.Fprintf(&o, "<p>")
			toc = append(toc, s)
			te = bytes.NewBuffer(nil)
		} else if strings.HasPrefix(s, "(") {
			fmt.Fprintf(&o, "%s<br>\n", hs(btrim(s)))
			te = bytes.NewBuffer(nil)
		} else if s == "" {
			if te != nil {
				b := te.Bytes()
				rows := bytes.Count(b, nl)
				s := string(b)
				fmt.Fprintf(&o, "\n<textarea rows=\"%d\" name='%s'>%s</textarea>\n<br/>\n", rows, id, s)
				te = nil
			}
		} else if te != nil {
			fmt.Fprintln(te, s)
			if idx := strings.Index(s, "]:"); idx != -1 {
				x := s[:idx]
				y := strings.LastIndex(x, "[")
				n := strings.TrimSpace(x[y+1:])
				v := strings.TrimSpace(x[:y])
				v = v[1 : len(v)-1]
				def = append(def, struct{ name, value string }{name: n, value: v})
			}
		}
	}

	toc = append(toc, "definitions")
	var jj bytes.Buffer
	fmt.Fprintf(&o, "<h1 id='definitions'>definitions</h1>\n")
	fmt.Fprintf(&o, "<table><tr><th>symbol</th><th>quotation</th></tr>\n")
	for _, d := range def {
		fmt.Fprintf(&o, "<tr><td id='%s'>%s</td><td>%s</tr>\n", d.name, hs(d.name), hs(d.value))
		fmt.Fprintf(&jj, "[%s][%s]:", d.value, d.name)
	}
	jj.Write([]byte{32})
	fatal(ioutil.WriteFile("j.j", jj.Bytes(), 0644))

	fmt.Fprintf(&o, "<ul>\n")
	for _, s := range toc {
		fmt.Fprintf(&o, "<li><a href=\"#%s\">%s</a></li>\n", s, s)
	}
	fmt.Fprintf(&o, "</ul>\n")

	o.WriteString(tail)

	fatal(ioutil.WriteFile("j.html", o.Bytes(), 0744))
}
func hs(s string) string    { return html.EscapeString(s) }
func btrim(s string) string { return strings.TrimSuffix(strings.TrimPrefix(s, "("), ")") }
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}
func uniq(x []string) (r []string) {
	m := make(map[string]bool)
	for _, s := range x {
		if m[s] == false {
			r = append(r, s)
		}
		m[s] = true
	}
	return r
}
func canvas() []byte {
	var buf bytes.Buffer
	buf.WriteString("(-canvas api)\n")
	v := strings.Split(cnvapi, "\n")
	for _, p := range v {
		if len(p) > 0 {
			idx := strings.Index(p, " ")
			long := p[:idx]
			def := p[idx+1:]
			idx = strings.Index(def, "]")
			sym := def[idx+1:]
			o := "[" + fmt.Sprint([]byte(long)) + "&][" + sym + "]:"
			buf.WriteString(o)
			sp := strings.Repeat(" ", 87-len(o))
			buf.WriteString(sp)
			buf.WriteString("() (" + long + ": " + def + ")\n")
		}
	}
	return buf.Bytes()
}

const head = `<head><meta charset="utf-8"><title>j</title></head>
<link rel=icon href='data:;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsQAAA7EAZUrDhsAAABSSURBVDhPY/wPBAwUACYozcDAyAjBhACaOoQBZAKCBjACbQNhXID2LgCFMb5wHgJeIARoYwAhfyMD2nqBGFdgNQA93slKB7AEhE8zCFCYnRkYAOS/HR/UGSYjAAAAAElFTkSuQmCC'>
<style>
 html{font-family:monospace}
 pre{background:#ffffea}
 textarea{background:black;color:white;width:100%;resize:none;overflow-y:hidden}
 ul{position:fixed;top:0;right:10}
 img{float:right}
 th{text-align:left;}
</style>
<script>
function ge(x){return document.getElementById(x)}
function toggle(id,e){e=ge(id);e.style.display=(e.style.display=='block')?'none':'block'}

</script>
<body>
`
const tail = `
</body></html>
`

// https://developer.mozilla.org/en-US/docs/Web/API/CanvasRenderingContext2D
const cnvapi = `arc [x y r sa ea]arc
arcTo [xa ya xb yb r]arcto
beginPath []bpath
bezierCurveTo [cxa cya cxb cyb x y]bezito
clearRect [x y w h]crect
clip []clip
closePath []cpath
createLinearGradient [xa ya xb yb]lingrd
createRadialGradient [xa ya ra xb yb rb]radgrd
ellipse [x y rx ry rot sa ea]ellips
fill []fill
fillRect [x y w h]frect
fillStyle [color|gradient]fstyle
fillText [s x y]ftext
font [name]font
lineCap [butt|round|square]lcap
lineDashOffset [i]ldoff
lineJoin [bevel|round|miter]ljoin
lineto [x y]lineto
lineWidth [i]lwidth
moveTo [x y]moveto
quadraticCurveTo [cx cy x y]qcto
rect [x y w h]rect
resetTransform []rstra
restore []rstore
rotate [a]rotate
save []save
scale [x y]scale
setLineDash [segments]sldash
setTransform [a b c d e f]setra
shadowBlur [i]shblur
shadowColor [c]shcol
shadowOffsetX [i]shoffx
shadowOffsetY [i]shoffy
stroke []stroke
strokeRect [x y w h]srect
strokeStyle [color|gradient]sstyle
strokeText [s x y]stext
textAlign [left|right|center|start|end]talign
textBaseline [b]tbline
transform [a b c d e f]transf
translate [x y]transl
`
