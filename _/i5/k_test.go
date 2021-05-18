package k

import (
	"bytes"
	"fmt"
	"html"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"
)

type testCase struct {
	file, name, i, e, comment string
	line                      int
}
type ref struct {
	name, mon, dya string
	tc             []testCase
}

//t 1 + 2 /3

func TestK(t *testing.T) {
	refs := getRef(t)
	tests := getTests(t)
	for _, tc := range tests {
		for i := range refs {
			if refs[i].name == tc.name {
				refs[i].tc = append(refs[i].tc, tc)
				break
			}
			if len(refs) == 1+i {
				t.Fatalf("wrong comment: %+v ", tc)
			}
		}
	}
	for _, tc := range tests {
		fmt.Println(tc)
		r := run([]byte(tc.i))
		if g := sstr(r); g != tc.e {
			t.Fatalf("in: %s\nexp %s\ngot %s", tc.i, tc.e, g)
		}
	}
	writeReadme(t, refs)
}
func run(src []byte) T {
	k := New()
	k.Trap = true
	return k.exec(k.fold(k.parse(src), src), src)
}
func writeReadme(t *testing.T, refs []ref) {
	h := html.EscapeString
	var buf bytes.Buffer
	for _, r := range refs {
		if r.name == "t" {
			continue
		}
		fmt.Fprintf(&buf, "<details><summary><code>%s</code>%s %s</summary>\n", h(r.name), h(r.mon), h(r.dya))
		for _, tc := range r.tc {
			fmt.Fprintf(&buf, "<a href=\"../master/%s#L%d\"><code>%s /%s</code>%s</a><br>\n", tc.file, 1+tc.line, h(tc.i), h(tc.e), h(tc.comment))
		}
		fmt.Fprintf(&buf, "</details>\n")
	}
	e := ioutil.WriteFile("readme.md", buf.Bytes(), 644)
	if e != nil {
		t.Fatal(e)
	}
}

func (tc testCase) String() string {
	return tc.file + ":" + strconv.Itoa(1+tc.line) + " " + tc.i + " /" + tc.e
}

func getRef(t *testing.T) (r []ref) {
	b, e := os.ReadFile("types.go")
	if e != nil {
		t.Fatal(e)
	}
	unil := func(s string) string {
		if s == "nil" {
			return ""
		}
		return s
	}
	v := bytes.Split(b, []byte("\n"))
	for i := range v {
		if s := string(v[i]); strings.HasSuffix(s, "//r") {
			var ri ref
			var e error
			k := strings.Index(s, `"`)
			s = s[1+k:]
			k = strings.Index(s, `"`)
			ri.name = s[:k]
			ri.name, e = strconv.Unquote(`"` + ri.name + `"`)
			if e != nil {
				t.Fatalf("%v: %s", e, s[:1+k])
			}
			k = 2 + k
			s = s[1+k:]
			k = strings.Index(s, ",")
			ri.mon = unil(strings.TrimPrefix(s[:k], "k."))
			s = s[2+k:]
			k = strings.Index(s, ",")
			ri.dya = unil(strings.TrimPrefix(s[:k], "k."))
			//fmt.Printf("%+v\n", ri)
			r = append(r, ri)
		}
	}
	def := []ref{
		{name: ":", dya: "assign"},
		{name: "λ", dya: "func"},
		{name: ";", dya: "statement"},
		{name: "t", dya: ""},
	}
	return append(r, def...)
}
func getTests(t *testing.T) (r []testCase) {
	d, e := os.ReadDir(".")
	if e != nil {
		t.Fatal(e)
	}
	linetest := func(f, s string, line int) (tc testCase) {
		k := strings.Index(s, " ")
		tc.name = s[:k]
		s = s[1+k:]
		k = strings.Index(s, " /")
		tc.file = f
		tc.i = s[:k]
		tc.e = s[2+k:]
		if k := strings.Index(tc.e, "  ("); k > 0 {
			tc.comment = tc.e[2+k:]
			tc.e = tc.e[:k]
		}
		tc.line = line
		return tc
	}
	testFile := func(n string) (r []testCase) {
		b, e := os.ReadFile(n)
		if e != nil {
			t.Fatal(e)
		}
		v := bytes.Split(b, []byte("\n"))
		for i := range v {
			s := string(v[i])
			if strings.HasPrefix(s, "//") {
				r = append(r, linetest(n, s[2:], i))
			}
		}
		return r
	}
	for _, f := range d {
		if n := f.Name(); f.IsDir() == false && strings.HasSuffix(n, ".go") {
			r = append(r, testFile(n)...)
		}
	}
	return r
}
