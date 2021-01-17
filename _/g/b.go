package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/cmplx"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"text/tabwriter"
)

var nyi = errors.New("nyi")
var plotKeys uint32
var xfile string
var xline int
var interactive bool

func ginit() {
	MT['p'+128] = plot1
	MT['c'+128] = caption1
	assign("caption", 'c')
	assign("plot", 'p')
	assign("WIDTH", mki(800))
	assign("HEIGHT", mki(400))
	assign("COLUMNS", mki(80))
	assign("LINES", mki(50))
	assign("FFMT", mkchars("%g"))
	assign("ZFMT", mkchars("%ga%.0f"))
	plotKeys = mksymbols([]string{"Type", "Style", "Limits", "Xlabel", "Ylabel", "Title", "Xunit", "Yunit", "Zunit", "Lines", "Foto", "Caption", "Data"})
}
func init() {
	exit = exitRepl
	Out = gOut
	// \  \h  (help)
	h := func(a string) bool {
		if a != `\h` && a != `\` {
			return false
		}
		fmt.Println(help)
		return true
	}
	replParsers = append(replParsers, h)
	kiniRunners = append(kiniRunners, ginit)
}

func exitRepl(x int) {
	if interactive {
		os.Exit(x)
	}
}

/*
func ExecLine(w io.Writer, s string) {
	switch s {
	default:
		if strings.HasPrefix(s, `\l`) {
			s = strings.TrimSpace(s[2:])
			kerr(Loadfile(s))
			if xline != 0 {
				fmt.Fprintf(os.Stderr, "%s:%d\n", xfile, xline)
				xline = 0
			}
		} else {
			s := Exec(s, terminal(w))
			if len(s) > 0 {
				fmt.Fprintln(w, s)
			}
		}
	}
}
func Exec(s string, w io.Writer) string {
	b := make([]uint64, len(MJ))
	copy(b, MJ)
	defer func() {
		if r := recover(); r != nil {
			printStack(debug.Stack())
			fmt.Println(r)
			ics(I(140), I(144), os.Stdout)
			MJ = b
			msl()
		}
	}()
	//dx(out(val(mkchrs([]byte(s)))))

	if len(strings.TrimSpace(s)) == 0 {
		return ""
	}
	r := val(mkchars(s))
}
*/
func gOut(x uint32) {
	var o io.Writer = os.Stdout
	if interactive {
		o = clipTerminal()
	}

	p := pk(x)
	if p != nil {
		showPlot(p)
		return
	}
	if rows, cols := istab(x); rows+cols > 0 {
		writeTable(x, o, rows, cols)
		return
	}

	rx(x)
	o.Write(append(Ck(kst(x)), 10))
}
func SetInteractive() { interactive = true }
func Loadfile(file string) error {
	b := make([]uint64, len(MJ))
	copy(b, MJ)
	defer func() {
		if r := recover(); r != nil {
			ics(I(140), I(144), os.Stdout)
			MJ = b
			msl()
		}
	}()
	if strings.HasSuffix(file, ".k") {
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()
		xfile = file
		_, err = runscript(f)
		return err
	}
	return fmt.Errorf("loadfile: unknown file type: %s\n", file)
}
func LoadDataFile(file, sym string) error {
	b, e := ioutil.ReadFile(file)
	if e != nil {
		return e
	}
	x := mk(1, uint32(len(b)))
	copy(MC[8+x:], b)
	assign(sym, x)
	return nil
}
func runscript(r io.Reader) (uint32, error) {
	xline = 0
	scn := bufio.NewScanner(includes(r))
	var x uint32
	for scn.Scan() {
		xline++
		s := strings.TrimSpace(scn.Text())
		if s == "" || strings.HasPrefix(s, "/") {
			continue
		}
		if idx := strings.Index(s, " /"); idx != -1 {
			s = s[:idx]
		}
		dx(x)
		x = val(mkchars(" " + s))
	}
	xline = 0
	return x, nil
}

// clip to COLUMNS/LINES if interactive
func clipTerminal() io.Writer {
	c := lupInt("COLUMNS")
	l := lupInt("LINES")
	if c <= 0 || l <= 0 {
		return os.Stdout
	}
	return &clipWriter{Writer: os.Stdout, c: c - 2, l: l}
}

func printStack(stack []byte) {
	if !interactive {
		debug.PrintStack()
		return
	}
	v := bytes.Split(stack, []byte{10})
	var o []string
	for _, b := range v {
		s := string(b)
		if strings.HasPrefix(s, "\t") {
			s = strings.TrimPrefix(s, "\t")
			if i := strings.Index(s, " +"); i > 0 {
				s = s[:i]
			}
			w := strings.Split(s, "/")
			if len(w) > 2 {
				w = w[len(w)-2:]
			}
			s = strings.Join(w, "/")
			if strings.HasPrefix(s, "debug") || strings.HasPrefix(s, "runtime") {
				continue
			}
			o = append(o, " "+s)
		}
	}
	if len(o) > 10 {
		o = o[:10]
	}
	if len(o) > 1 {
		for i := len(o) - 1; i >= 0; i-- {
			fmt.Println(o[i])
		}
	}
}

type clipWriter struct {
	io.Writer
	c, l int
	x, y int
}

func (cw *clipWriter) Write(p []byte) (n int, err error) {
	size := len(p)
	for {
		if len(p) == 0 {
			break
		} else if cw.l > 0 && cw.y >= cw.l {
			if cw.y == cw.l {
				cw.Writer.Write([]byte("..\n"))
			}
			cw.y++
			break
		}
		idx := bytes.IndexByte(p, '\n')
		if idx == -1 {
			if xx := cw.x + len(p); xx > cw.c {
				p = p[:cw.c-cw.x]
			}
			cw.Writer.Write(p)
			cw.x += len(p)
			break
		} else {
			if xx := cw.x + idx - 1; xx > cw.c {
				cw.Writer.Write(p[:cw.c-cw.x])
				cw.Writer.Write([]byte("..\n"))
			} else {
				cw.Writer.Write(p[:idx+1])
			}
			p = p[idx+1:]
			cw.x = 0
			cw.y++
		}
	}
	return size, nil
}

func lupInt(s string) int {
	r := lookup(s)
	if r == 0 {
		panic("var " + s + " does not exist")
	}
	if tp(r) != 2 || nn(r) != 1 {
		panic("var " + s + " is not int#1")
	}
	dx(r)
	return int(MI[2+r>>2])
}
func lupString(s string) string {
	r := lookup(s)
	if r == 0 || tp(r) != 1 {
		panic("var " + s + " must exist as chars")
	}
	return kstr(r)
}
func assign(s string, v uint32)  { dx(asn(mksymbol(s), v)) }
func lookup(s string) (r uint32) { return lup(mksymbol(s)) }
func mksymbol(s string) (r uint32) {
	r = mk(1, uint32(len(s)))
	copy(MC[8+r:], s)
	return sc(r)
}
func mksymbols(v []string) (r uint32) {
	r = mk(5, 0)
	for _, s := range v {
		r = ucat(r, mksymbol(s))
	}
	return r
}
func mkchars(s string) (r uint32) {
	r = mk(1, uint32(len(s)))
	copy(MC[8+r:], s)
	return r
}
func mkfloat(f float64) (r uint32) {
	r = mk(3, 1)
	p := (r + 8) >> 3
	MF[p] = f
	return r
}
func mkfloats(f []float64) (r uint32) {
	n := uint32(len(f))
	r = mk(3, n)
	p := (r + 8) >> 3
	copy(MF[p:], f)
	return r
}
func mkcmplx(z []complex128) (r uint32) {
	n := uint32(len(z))
	r = mk(4, n)
	p := (r + 8) >> 3
	for i := uint32(0); i < n; i++ {
		MF[p+2*i] = real(z[i])
		MF[p+1+2*i] = imag(z[i])
	}
	return r
}
func kerr(e error) bool {
	if e == nil {
		return false
	}
	fmt.Fprintln(os.Stderr, e)
	if interactive == false {
		os.Exit(1)
	}
	return true
}
func perr(e error) {
	if e != nil {
		panic(e)
	}
}

func kstr(x uint32) (r string) { n := nn(x); return string(MC[x+8 : x+8+n]) }

func symstr(off uint32) string { return kstr(I(I(kkey) + off)) }
func istab(x uint32) (uint32, uint32) {
	if tp(x) == 7 {
		v := MI[3+x>>2]
		n := nn(v)
		if n == 0 {
			return 0, 0
		}
		n0 := nn(MI[2+v>>2])
		for i := uint32(0); i < n; i++ {
			if nn(MI[2+i+v>>2]) != n0 {
				return 0, 0
			}
		}
		return n0, nn(v)
	}
	return 0, 0
}
func writeTable(x uint32, ww io.Writer, rows, cols uint32) {
	tab := []byte{'\t'}
	nl := []byte{'\n'}
	ffmt := lupString("FFMT")
	zfmt := lupString("ZFMT")
	w := tabwriter.NewWriter(ww, 2, 8, 1, ' ', 0)
	keys, vals := MI[2+x>>2], MI[3+x>>2]
	for i := uint32(0); i < cols; i++ {
		w.Write([]byte(symstr(MI[2+keys>>2])))
		if i != cols-1 {
			w.Write(tab)
		}
	}
	w.Write(nl)
	for k := uint32(0); k < rows; k++ {
		for i := uint32(0); i < cols; i++ {
			w.Write([]byte(fmtVecAt(MI[2+i+vals>>2], k, ffmt, zfmt)))
			if i != cols-1 {
				w.Write(tab)
			}
		}
		w.Write(nl)
	}
	w.Flush()
}

func fmtVecAt(x uint32, i uint32, ffmt, zfmt string) string {
	switch tp(x) {
	case 2:
		return strconv.Itoa(int(MI[2+i+x>>2]))
	case 3:
		return fmt.Sprintf(ffmt, MF[1+i+x>>3])
	case 4:
		z := complex(MF[1+2*i+x>>3], MF[2+2*i+x>>3])
		return absang(z, zfmt)
	case 5:
		return symstr(MI[2+i+x>>2])
	default:
		return "?"
	}
}

// resolve backslash includes: \file.k
func includes(r io.Reader) io.Reader {
	b, e := ioutil.ReadAll(r)
	if e != nil {
		panic(e)
	}
	b = bytes.Replace(b, []byte("\r"), []byte{}, -1)
	v := bytes.Split(b, []byte("\n"))

	var buf bytes.Buffer
	for _, u := range v {
		if len(u) > 3 && u[0] == '\\' && bytes.HasSuffix(u, []byte(".k")) {
			if f, e := ioutil.ReadFile(string(u[1:])); e != nil {
				panic(e)
			} else {
				f = bytes.Replace(f, []byte("\r"), []byte{}, -1)
				buf.Write(f)
			}
		} else {
			buf.Write(u)
			buf.Write([]byte{'\n'})
		}
	}
	return &buf
}
func absang(x complex128, format string) string {
	if format == "" {
		format = "%va%v"
	}
	r, phi := cmplx.Polar(x)
	phi *= 180.0 / math.Pi
	if phi < 0 {
		phi += 360.0
	}
	if r == 0.0 {
		phi = 0.0
	}
	if phi == -0.0 || phi == 360.0 {
		phi = 0.0
	}
	return fmt.Sprintf(format, r, phi)
}
