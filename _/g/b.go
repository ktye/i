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
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"
)

var nyi = errors.New("nyi")
var plotKeys uint32
var symbols = make(map[string]uint32)
var xfile string
var xline int

func ginit() {
	MT['R'+128] = read1
	MT['D'+128] = dir1
	MT['C'+128] = csv1
	MT['C'] = csv2
	MT['r'+128] = rand1 // shuffle
	MT['r'] = rand2
	MT['p'+128] = plot1
	MT['q'+128] = qr1
	MT['q'] = solve2
	MT['d'+128] = diag1
	MT['m'] = mul2
	MT['c'+128] = caption1
	assign("read", 'R')
	assign("dir", 'D')
	assign("csv", 'C')
	assign("randi", prj('r', l2(mki(2), 0), mki(1)))          // randi: 'r[2;]
	assign("randf", prj('r', l2(mki(3), 0), mki(1)))          // randf: 'r[3;]
	assign("randn", prj('r', l2(mki(4294967293), 0), mki(1))) // randn: 'r[-3;]
	assign("randz", prj('r', l2(mki(4), 0), mki(1)))          // randz: 'r[4;]
	assign("rands", prj('r', l2(mki(5), 0), mki(1)))          // rands: 'r[5;]
	assign("shuffle", 'r')
	assign("caption", 'c')
	assign("plot", 'p')
	assign("qr", 'q')
	assign("cond", prj('q', l2(mkc('c'), 0), mki(1))) // cond: 'q["c";]
	assign("diag", 'd')
	assign("mul", 'm') // XT*y or XT*Y r|z
	assign("solve", 'q')
	assign("sin", prj('F', l2(mki(129), 0), mki(1)))
	assign("cos", prj('F', l2(mki(130), 0), mki(1)))
	assign("exp", prj('F', l2(mki(131), 0), mki(1))) // pow10:{131 'F2.302585092994046*x}
	assign("log", prj('F', l2(mki(132), 0), mki(1))) // log10:{0.4342944819032518* 'Fx}
	for _, s := range []string{"WIDTH", "HEIGHT", "COLUMNS", "LINES", "FFMT", "ZFMT"} {
		symbols[s] = ks(s)
	}
	assign("WIDTH", mki(800))
	assign("HEIGHT", mki(400))
	assign("COLUMNS", mki(80))
	assign("LINES", mki(20))
	assign("FFMT", kC([]byte("%.4g")))
	assign("ZFMT", kC([]byte("%.4ga%.0f")))
	plotKeys = kS([]string{"Type", "Style", "Limits", "Xlabel", "Ylabel", "Title", "Xunit", "Yunit", "Zunit", "Lines", "Foto", "Caption", "Data"})

}

var stdout io.Writer

func init() {
	exit = exitRepl
	leak = bleak
	stdout = os.Stdout
	Out = gOut

	// -s lines, cols (terminal size)
	s := func(a string, tail []string) ([]string, bool) {
		if a != "-s" || len(tail) < 2 {
			return tail, false
		}
		lines, cols := atoi(tail[0]), atoi(tail[1])
		assign("LINES", ki(lines))
		assign("COLUMNS", ki(cols))
		assign("WIDTH", ki(cols*11))
		assign("HEIGHT", ki(lines*20))
		return tail[2:], true
	}
	argvParsers = append(argvParsers, s)

	// \h(help)
	h := func(a string) bool {
		if a != `\h` && a != `\` {
			return false
		}
		fmt.Println(help)
		return true
	}
	// \c(caption)
	c := func(a string) bool {
		if a != `\c` {
			return false
		}
		if lastCaption != nil {
			w, _ := clipTerminal()
			lastCaption.WriteTable(w, 0)
		}
		return true
	}
	replParsers = append(replParsers, h, c)
	kiniRunners = append(kiniRunners, ginit)
}

func exitRepl(x int) {
	if interactive {
		os.Exit(x)
	}
}
func bleak() {
	b := make([]uint64, len(MJ))
	copy(b, MJ)
	dx(plotKeys)
	for _, x := range symbols {
		dx(x)
	}
	_leak()
	copy(MJ, b)
	msl()
}
func memstore() []uint32 {
	m := make([]uint32, len(MI))
	copy(m, MI)
	return m
}
func memcompare(m []uint32, s string) {
	if len(m) != len(MI) {
		panic(fmt.Sprintf("%s modified memory size: before %d now %d\n", s, len(m), len(MI)))
	}
	for i, u := range m {
		if u != MI[i] {
			panic(fmt.Sprintf("%s modified memory at %x(%d): 0x%x != 0x%x", s, i, i, m[i], MI[i]))
		}
	}
}
func gOut(x uint32) {
	if x == 0 {
		return
	}
	o, lines := clipTerminal()
	//m := memstore()
	p := pk(x)
	//memcompare(m, "pk")
	if p != nil {
		showPlot(p)
		//memcompare(m, "showplot")
		return
	}
	if rows, cols := istab(x); rows > 1 {
		writeTables(x, 0, o, rows, cols, 0, lines)
		return
	} else if tp(x) == 7 {
		writeDict(x, o, lines)
		return
	} else if ismatrix(x) {
		writeMatrix(x, o, lines)
		return
	} else if tp(x) == 6 {
		if tx, ty, rows, xcols, ycols := iskeytab(x); tx != 0 {
			writeTables(tx, ty, o, rows, xcols, ycols, lines)
			return
		}
	}
	rx(x)
	o.Write(append(CK(kst(x)), 10))
}
func Loadfile(file string) error {
	b := make([]uint64, len(MJ))
	copy(b, MJ)
	defer func() {
		if r := recover(); r != nil {
			ics(I(140), I(144), stdout)
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
func read1(x uint32) uint32 {
	if tp(x) == 5 && nn(x) == 1 {
		x = cs(x)
	}
	if tp(x) != 1 {
		panic("type")
	}
	if nn(x) == 0 {
		dx(x)
		x = kC([]byte("./"))
	}
	if MC[7+x+nn(x)] == '/' {
		return readdir(string(CK(x)))
	}
	b, e := ioutil.ReadFile(string(CK(x)))
	fatal(e)
	return kC(b)
}
func readdir(s string) uint32 {
	fi, e := ioutil.ReadDir(s)
	fatal(e)
	keys, vals := make([]string, len(fi)), make([]int, len(fi))
	for i, f := range fi {
		keys[i] = f.Name()
		if f.IsDir() {
			keys[i] += "/"
		}
		vals[i] = int(f.Size())
	}
	return mkd(kS(keys), kI(vals))
}
func dir1(x uint32) uint32 { // dir"*.k"
	s, e := filepath.Glob(string(CK(x)))
	if e != nil {
		panic(e)
	}
	return kS(s)
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
		x = val(kC([]byte(" " + s)))
	}
	xline = 0
	return x, nil
}

// clip to COLUMNS/LINES if interactive
func clipTerminal() (io.Writer, int) {
	if !interactive {
		return stdout, 0
	}
	c := lupInt("COLUMNS")
	l := lupInt("LINES")
	if c <= 0 || l <= 0 {
		return stdout, 0
	}
	return &clipWriter{Writer: stdout, c: c - 2, l: l - 2}, l
}

func atoi(s string) int {
	i, e := strconv.Atoi(s)
	fatal(e)
	return i
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

func lupInt(s string) int { // no modification
	r := lookup(s)
	if tp(r) != 2 || nn(r) != 1 {
		panic("var " + s + " is not int#1")
	}
	return int(MI[2+r>>2])
}
func lupString(s string) string { // no modification
	r := lookup(s)
	if tp(r) != 1 {
		panic("var " + s + " is not a char")
	}
	return string(Ck(r))
}
func lookup(s string) uint32 { // no modification
	x, o := symbols[s]
	if o == false {
		panic("var " + s + " is not a registered symbol")
	}
	return I(I(kval) + I(x+8))
}
func assign(s string, v uint32) { dx(asn(ks(s), v)) }
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

func istab(x uint32) (rows, cols int) {
	if tp(x) == 7 {
		v := MI[3+x>>2]
		n := nn(v)
		if n == 0 {
			return 0, 0
		}
		n0 := nn(MI[2+v>>2])
		for i := uint32(0); i < n; i++ {
			if nn(MI[2+i+v>>2]) != n0 {
				return 0, -1
			}
		}
		return int(n0), int(nn(v))
	}
	return 0, -1
}
func iskeytab(x uint32) (rx, ry uint32, rows, xcols, ycols int) {
	if tp(x) == 6 && nn(x) == 2 {
		tx, ty := MI[2+x>>2], MI[3+x>>2]
		xrows, xcols := istab(tx)
		yrows, ycols := istab(ty)
		if xrows == yrows && xrows > 1 {
			return tx, ty, xrows, xcols, ycols
		}
	}
	return 0, 0, 0, 0, 0
}
func ismatrix(x uint32) (r bool) { // rectangular real/complex column major
	if tp(x) != 6 || nn(x) < 1 {
		return false
	}
	t := tp(MI[2+x>>2])
	if t < 3 || t > 5 {
		return false
	}
	n := nn(MI[2+x>>2])
	for i := uint32(0); i < nn(x); i++ {
		p := MI[2+i+x>>2]
		if tp(p) != t || nn(p) != n {
			return false
		}
	}
	return true
}
func writeMatrix(x uint32, ww io.Writer, clip int) {
	cols := nn(x)
	rows := nn(MI[2+x>>2])
	tab := []byte{'\t'}
	nl := []byte{'\n'}
	ffmt := lupString("FFMT")
	zfmt := lupString("ZFMT")
	fmt.Fprintf(ww, "%dx%d\n", rows, cols)
	w := tabwriter.NewWriter(ww, 2, 8, 1, ' ', 0)
	for i := uint32(0); i < rows; i++ {
		if clip > 0 && int(i) == clip {
			w.Flush()
			fmt.Fprintf(ww, "..\n")
			return
		}
		for k := uint32(0); k < cols; k++ {
			w.Write([]byte(fmtVecAt(MI[2+k+x>>2], i, ffmt, zfmt)))
			if k == cols-1 {
				w.Write(nl)
			} else {
				w.Write(tab)
			}
		}
	}
	w.Flush()
}
func writeDict(x uint32, w io.Writer, clip int) {
	k := Sk(I(8 + x))
	m := 1
	for i := range k {
		if n := len(k[i]); n > m {
			m = n
		}
	}
	rx(x)
	x = val(x)
	x = ech(x, 'k'+128)
	for i, s := range k {
		fmt.Fprintf(w, "%s%s|%s\n", s, strings.Repeat(" ", m-len(s)), string(Ck(I(8+x+uint32(4*i)))))
		if clip > 0 && i > clip {
			fmt.Fprintf(w, "..\n")
			break
		}
	}
	dx(x)
}
func writeTables(x, y uint32, ww io.Writer, rows, xcols, ycols int, clip int) {
	tab := []byte{'\t'}
	nl := []byte{'\n'}
	ffmt := lupString("FFMT")
	zfmt := lupString("ZFMT")
	w := tabwriter.NewWriter(ww, 2, 8, 1, ' ', 0)
	xkeys, xvals := Sk(MI[2+x>>2]), MI[3+x>>2]
	ykeys, yvals := []string{}, uint32(0)
	if y != 0 {
		ykeys, yvals = Sk(MI[2+y>>2]), MI[3+y>>2]
	}
	for i := range xkeys {
		w.Write([]byte(xkeys[i]))
		if i != int(xcols-1) {
			w.Write(tab)
		}
	}
	for i := range ykeys {
		if i == 0 {
			w.Write([]byte("|\t"))
		}
		w.Write([]byte(ykeys[i]))
		if i != int(ycols-1) {
			w.Write(tab)
		}
	}
	w.Write(nl)
	for k := 0; k < rows; k++ {
		for i := 0; i < xcols; i++ {
			w.Write([]byte(fmtVecAt(MI[2+uint32(i)+xvals>>2], uint32(k), ffmt, zfmt)))
			if i != xcols-1 {
				w.Write(tab)
			}
		}
		for i := 0; i < ycols; i++ {
			if i == 0 {
				w.Write([]byte("|"))
			}
			w.Write(tab)
			w.Write([]byte(fmtVecAt(MI[2+uint32(i)+yvals>>2], uint32(k), ffmt, zfmt)))
		}
		w.Write(nl)
		if clip > 0 && int(k) > clip {
			w.Write([]byte("..\n"))
			break
		}
	}
	w.Flush()
}

func fmtVecAt(x uint32, i uint32, ffmt, zfmt string) string {
	switch tp(x) {
	case 2:
		return strconv.Itoa(int(int32(MI[2+i+x>>2])))
	case 3:
		return fmt.Sprintf(ffmt, MF[1+i+x>>3])
	case 4:
		z := complex(MF[1+2*i+x>>3], MF[2+2*i+x>>3])
		return absang(z, zfmt)
	case 5:
		return ski(MI[2+i+x>>2])
	case 6:
		xi := MI[2+i+x>>2]
		if nn(xi) != 1 {
			if tp(xi) == 1 {
				return string(Ck(xi))
			} else if tp(xi) < 6 {
				dots, n := "", nn(xi)
				if n > 4 {
					dots, n = "..", 4
				}
				r := make([]string, n)
				for i := uint32(0); i < n; i++ {
					r[i] = fmtVecAt(xi, i, ffmt, zfmt)
				}
				return strings.Join(r, " ") + dots
			}
			break
		}
		return fmtVecAt(xi, 0, ffmt, zfmt)
	}
	return "?"
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
