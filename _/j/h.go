// +build ignore

// generate j.html
package main

import (
	"bytes"
	"fmt"
	"html"
	"io/ioutil"
	"regexp"
	"strings"
)

func main() {
	var def []struct{ name, value string }
	var o bytes.Buffer
	o.WriteString(head)

	t, e := ioutil.ReadFile("readme")
	fatal(e)
	fmt.Fprintf(&o, "<pre id=\"ref\">\n%s</pre>\n", string(t))

	t, e = ioutil.ReadFile("j.w")
	fatal(e)
	fmt.Fprintf(&o, "<h1 id=\"src\">j.w source of the wasm module</h1> <span style='color:blue' onclick=\"toggle('jw')\">(show/hide)</span><br/>\n<pre id='jw' style='display:none'>\n%s</pre>\n", string(t))

	b, e := ioutil.ReadFile("t")
	fatal(e)

	words := regexp.MustCompile(`[a-z]+`)
	toc := []string{"ref", "src"}
	var te *bytes.Buffer
	nl := []byte("\n")
	v := bytes.Split(b, nl)
	var id string
	var localsyms []string
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
		} else if s == "" {
			if te != nil {
				b := te.Bytes()
				rows := bytes.Count(b, nl)
				s := string(b)

				fmt.Fprintf(&o, "</p>\n<img src='javascript:favi()' onclick='run(\"%s\")'/>\n", id)
				fmt.Fprintf(&o, "%s\n", strings.Join(uniq(localsyms), " "))
				fmt.Fprintf(&o, "\n<textarea rows=\"%d\" name='%s'>%s</textarea>\n<button onclick='run(\"%s\")'>run</button>", rows, id, s, id)
				localsyms = nil
				te = nil
			}
		} else if te != nil {
			fmt.Fprintln(te, s)
			if idx := strings.Index(s, "]:"); idx != -1 {
				x := s[:idx]
				y := strings.LastIndex(x, "[")
				n := strings.TrimSpace(x[y+1:])
				v := strings.TrimSpace(x[:y])
				v = v[1 : len(v)-2]
				def = append(def, struct{ name, value string }{name: n, value: v})
			}
			c := s
			if idx := strings.Index(s, "("); idx != -1 {
				c = s[:idx]
			}
			v := words.FindAllString(c, -1)
			for _, sym := range v {
				localsyms = append(localsyms, fmt.Sprintf("<a href='#%s'>%s</a>", sym, sym))
			}
		}
	}

	toc = append(toc, "definitions")
	fmt.Fprintf(&o, "<h1 id='definitions'>definitions</h1>\n")
	fmt.Fprintf(&o, "<table><tr><th>symbol</th><th>quotation</th></tr>\n")
	for _, d := range def {
		fmt.Fprintf(&o, "<tr><td id='%s'>%s</td><td>%s</tr>\n", d.name, hs(d.name), hs(d.value))
	}

	fmt.Fprintf(&o, "<ul>\n")
	for _, s := range toc {
		fmt.Fprintf(&o, "<li><a href=\"#%s\">%s</a></li>\n", s, s)
	}
	fmt.Fprintf(&o, "</ul>\n")

	o.WriteString(tail)

	fatal(ioutil.WriteFile("j.html", o.Bytes(), 0744))
	//io.Copy(os.Stdout, &o)
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

const head = `<head><meta charset="utf-8"><title>j</title></head>
<link rel=icon href='data:;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsQAAA7EAZUrDhsAAABSSURBVDhPY/wPBAwUACYozcDAyAjBhACaOoQBZAKCBjACbQNhXID2LgCFMb5wHgJeIARoYwAhfyMD2nqBGFdgNQA93slKB7AEhE8zCFCYnRkYAOS/HR/UGSYjAAAAAElFTkSuQmCC'>
<style>
 html{font-family:monospace}
 pre{background:#ffffea}
 textarea{background:black;color:#cccccc;width:100%;resize:none;overflow-y:hidden}textarea:focus{color:white;}
 ul{position:fixed;top:0;right:10}
 button{background:#ffffea;border:1px solid black;float:left}
 img{float:right}
 th{text-align:left;}
</style>
<script>
function ge(x){return document.getElementById(x)}
function toggle(id,e){e=ge(id);e.style.display=(e.style.display=='block')?'none':'block'}
function favi(){var f=undefined;var l=document.getElementsByTagName("link");for(var i=0;i<l.length;i++)if(l[i].getAttribute("rel")=="icon"){f=l[i].getAttribute("href");}return f;}
function run(x){console.log(run,x)}

</script>
<body>
`
const tail = `
<script>
var l=document.getElementsByTagName("img");for(var i=0;i<l.length;i++)l[i].src=favi();
</script>
</body></html>
`

/*
 html,body,textarea,input,select{margin:0;padding:0;overflow:hidden;font-family:monospace;overflow-x:hidden}
 table{position:absolute;width:100%;height:100%;border-collapse:collapse;}td{width:50%;}
 textarea{top:0;left:0;width:100%;height:100%;background:black;color:#cccccc;border:none;resize:none;overflow-y:auto;overflow-x:hidden;scrollbar-width:none;}
 ::-webkit-scrollbar{width:0;height:0;}
 textarea:focus{color:white;}.hold{background:#666666}
*/
