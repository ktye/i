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

	t, e = ioutil.ReadFile("j.wasm")
	fatal(e)
	cnt := len(t)

	t, e = ioutil.ReadFile("j.w")
	fatal(e)
	fmt.Fprintf(&o, "<h1 id=\"src\">j.w source of the wasm module</h1> <span style='color:blue' onclick=\"toggle('jw')\">(show/hide)</span><br/>\n<pre id='jw' style='display:none'>\n%s</pre>\n<code>w < j.w > j.wasm(%d byte)</code>\n", hs(string(t)), cnt)

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
	fmt.Fprintf(&o, "</table>\n")
	jj.Write([]byte{32})
	fatal(ioutil.WriteFile("j.j", jj.Bytes(), 0644))

	toc = append(toc, "examples")
	fmt.Fprintf(&o, "<h1 id='examples'>examples</h1>\n")
	examples := []struct{ value, url string }{
		{"ellipse", "eJyLNjEwUABhIyA2BGIQNDYziOWKTipILMlQUDJxMNE2qVNRSs3JySwoViguKcrPTo2NTitKzE2NtQIA6qoR2Q=="},
		{"roll", "eJx1j8EKwjAMhu99irCjdpBuVGGvEnKos8Nht0oX5s1ntzIUD90tf778fxJCaBBZkpuXAJeHkxshIE9x9RLhCajDOH/KRVK8e0XbTAsI7Qm5dwIu9TCMITBdo3CnUJECIAPf3KzMsSJ+6RTFic867zCItYZfks7mDZxt7peAxR3Q7DlyF0qgsboMLO6AfK0ugApNffj7jGlIbvLcvQGyg1bT"},
		{"koch", "eJxFjs1Ow0AMhO9+ihEnqlK6ragQHBCnXNo3sHxwWreJSHZhu+rPpc+OE0BoZMn6xmP7ftPGncVer5bn67RtUGH2hmpazVzTakLEvAgIUrLGY8dD27XRShK+yKu7WD4HyaloMWH9QXj5R7UjZ3LjB0kny63wvtMyDgYsELD8Le/FzcG5C/aFc2PZsIf2FncYQ3wYthGvAp5Wf09R/amlwfhb7ydKIgJfoKiFPH8Y5ORdTLdNfiQ6lpw+7Bvl1Uvo"},
		{"sierpinski", "eJxVj8FOwzAQRO/+ilFPoCbFDuHSQ9UP4MbR2oNDN47VxCmOBeTCt7NJ1SJkrT2at7JmHl5DPHEc3Mzp6S1wuoQ4nQPa0pe+kOcgattufdkW8OXB+0elrDUv0JSTi1Nvtcg+RM4j2W/aQy470162UNWa0phdZrLuapnqz2rEEo9+bEHjJ6dAtu1dXhc1nlGjgpapRWsSuBIDI9oveqP5A18dJ0YLN3A8YflAYWPuwP8H9R2YKyDbLTGUNRrmVko1F5c7rN0GiZZHJb0ww6EhhTWdWYbQ3Y7YR2L33qWdmnIaz/wLLcJqIw=="},
		{"arrowhead", "eJxVj0FOw0AMRfdziq+uQE1ghrYsuqg4ADcYeeFQzyQiycB0BGTD2XFSWoQsW9b7tvT/zXM3HmUceJJ8zzmnz1b4iFAh1IdYhzpWiPUhrOM63BrjvbOwVDKPp97Pa9+NUhL5L9pDh59or1fYWEs5FS5Cns8Ij3+oUaSMvn1F6UNyRz70XJbDLTaw2ltSuBCLBzhtqyTOZGXlHeo0CwJ4kFEd67vByl2F+F/YXgV3Fsi3swkzJ9rtLplM88alxRJtUGclGY2FCYyGdHW/5hyhvZTiJxJ+afOdOZWcXuUHPm1oAw=="},
		{"tree", "eJxNUMtOwzAQvPsrRj2BFMBueEg9ID6AP1jtwaHbJiKJW8cq8YVvZx0KrSzZo5nd2VnfvHfjVsbBZ4kPKYrAVsh3rzlXmPWBx4wG860xRM7Ccop+nHoqsO9GSYFp5g30oswbQ5M/CUCo3RPHkHwSJr8IKYZPQZxSiILm4FO7uAzhpC70eClvtFzH8TdVrGLsmHa9T8WEHNawqLVPyTPjFO8LXjk54qsV9d/DDzJuURoNVvZf2F0JTO0yqrYWz/ZvN9D6xZ7DmCUnroIaXRXlXxpWaBnt5SjxxuI/2nhvftf9AT4ab90="},
		{"square-spiral", "eJwlijEOgzAMAHe/wupERQdHdGm/YnlwJSMQgYBjeH9BDKcb7prYPbJh3Xb1U+vomp8AwB0Rnki4LjXDG/i3agzIhCRzOSwKpvaB9LqPO+RxuUINL5Mhf5J4CQ0T7l1nky/8ATTBJSA="},
		{"reuleux", "eJw9jDEOgEAIBPt7BaV2aGHhVwgFmjMaPc8gJj5fvBiLpdhZptJ4bVGuuw4hECEgtPjeDll0ZBLuvS8dm8p+btS0yJpNLDKlFw+H2AwC6c9pmlfHposPQvn3fAYvHICLfs+kkiL3DwDFLH8="},
		{"plant", "eJxNkM9uwjAMxu99ik89gUq3tMBhHKadctoDVLJ8CBBoRWmZGzF62bPPKeyPHFmff7ZlO7P3ptv77uxGL8+X1nUB1QJV/mozoopzfZZyW3Gm2Cq28ySh7cWFmgwMn/urDz2tVQZx3dBOtG06pRiC9CfPdOMNiGnkTUKE5XLN0gcXNOPuCOUf2kY0uKvKIUoZQi8aiAZay1+0YB0qDdOhdUHpDBYZSK0CI394q3H0McpQzbXXoMRKrcAaS5iHX00+shIFMx2nneIZOiHq1PgPfNZePA5wZ9/tEScnSIvfxPFfgqmedr1hhMMWA4STI+q7JfTG3u1q0QodINo3XYzH/5UvP3+B9EmYdm2zO8WKwuiuxvAEvgFz64MO"},
	}
	fmt.Fprintf(&o, "<p>\n")
	for _, e := range examples {
		fmt.Fprintf(&o, "<a href=\"/js.html#%s\">%s</a>\n", e.url, hs(e.value))
	}
	fmt.Fprintf(&o, "<p>\n")

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
			e := ""
			if idx == 1 {
				e = "[]"
				def = def[2:]
			}
			o := "[" + e + fmt.Sprint([]byte(long)) + "&][" + sym + "]:"
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
lineTo [x y]lineto
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
