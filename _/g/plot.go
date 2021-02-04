package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ktye/plot"
)

// plot floats (xy)
// plot complex (polar)
// plot (..) (multiline)
// plot dict (KTablePlot `x`y!..)
// plot grouped dict (multiline)
func plot1(x uint32) (r uint32) {
	xt, xn := tp(x), nn(x)
	_ = xn
	switch xt {
	case 2:
		return plot1(add(kf(0), x))
	case 3:
		p := plot.Plot{Type: plot.XY}
		p.Lines = []plot.Line{pline(x, 0)}
		return K(p)
	case 4:
		p := plot.Plot{Type: plot.Polar}
		p.Lines = []plot.Line{pline(x, 0)}
		return K(p)
	case 6, 7:
		if xn == 0 {
			return 0
		}
		// plot `c=`x`y`c!(x,x;(sin x),cos x:!10;(10#`a),10#`b)
		if tx, ty, _, _, _ := iskeytab(x); tx != 0 {
			rx(ty)
			dx(x)
			k, v := kvd(ty)
			x = ecr(k, flp(v), 33) // k!/:&v
			xn = nn(x)
		}
		plts, e := plot.KTablePlot(x, MC, MI, MF)
		perr(e)
		if len(plts) == 1 {
			r = K(plts[0])
		} else {
			r = mk(6, 0)
			for _, p := range plts {
				r = lcat(r, K(p))
			}
		}
		return dxr(x, r)
	}
	return dxr(x, 0)
}
func plot2(x, y uint32) (r uint32) { // plot with caption(table)
	x = cal('p', enl(x))
	if tp(y) != 7 {
		panic("plot2: y type: not a table")
	}
	var c plot.Caption
	keys := Sk(MI[2+y>>2])
	v := MI[3+y>>2]
	n := nn(v)
	for i := uint32(0); i < n; i++ {
		var col plot.CaptionColumn
		col.Name = keys[i]
		var vi = MI[2+i+v>>2]
		switch tp(vi) {
		case 2:
			col.Data = Ik(vi)
		case 3:
			col.Data = Fk(vi)
		case 4:
			col.Data = Zk(vi)
		case 5:
			col.Data = Sk(vi)
		default:
			panic(fmt.Errorf("illegal caption column type: %d", i))
		}
		c.Columns = append(c.Columns, col)
	}
	dx(y)
	return asi(x, l2(mki(0), ks("Caption")), K(c))
}

//type Caption map[string]interface{}

func caption1(x uint32) (r uint32) { // caption plot`sig
	p := pk(x)
	if p == nil {
		return x
	}
	c, e := p.MergedCaption()
	perr(e)
	n := uint32(len(c.Columns))
	v := mk(6, n)
	keys := make([]string, n)
	for i := uint32(0); i < n; i++ {
		keys[i] = c.Columns[i].Name
		MI[2+i+v>>2] = K(c.Columns[i].Data)
	}
	return mkd(kS(keys), v)
}

func pk(x uint32) (r plot.Plots) { // plot from k (does not unref)
	xt, xn := tp(x), nn(x)
	if xt == 7 {
		if p, ok := toPlot(x); ok {
			return plot.Plots{p}
		}
	} else if xt == 6 && xn >= 1 {
		for i := uint32(0); i < xn; i++ {
			xi := MI[2+i+x>>2]
			if p, ok := toPlot(xi); ok {
				r = append(r, p)
			} else {
				return nil
			}
		}
		return r
	}
	return nil
}
func toPlot(x uint32) (plot.Plot, bool) { // does not unref
	if tp(x) == 7 && match(MI[2+x>>2], plotKeys) != 0 {
		var p plot.Plot
		G(x, &p)
		return p, true
	}
	return plot.Plot{}, false
}
func pline(y uint32, id int) (r plot.Line) {
	x := add(kf(0.0), til(mki(nn(y))))
	return pline2(x, y, id)
}
func pline2(x, y uint32, id int) (r plot.Line) {
	r.Id = id
	if t := tp(y); t == 3 {
		G(y, &r.Y)
	} else if t == 4 {
		G(y, &r.C)
	}
	G(x, &r.X)
	dx(x)
	dx(y)
	return r
}
func screensize() (int, int) {
	w := lupInt("WIDTH")
	h := lupInt("HEIGHT")
	return w, h
}

var lastCaption *plot.Caption

func showPlot(p plot.Plots) {
	c, e := p.MergedCaption()
	perr(e)
	lastCaption = &c
	if dir := lupString("PLTDIR"); dir == "" {
		w, h := screensize()
		ip, e := p.IPlots(w, h, 0)
		perr(e)
		m := plot.Image(ip, nil, w, h, 0).(*image.RGBA)
		drawTerm(pngData(m))
	} else {
		w, e := ioutil.TempFile(dir, "*.plt")
		name := w.Name()
		perr(e)
		defer w.Close()
		perr(p.Encode(w))
		fmt.Println(strings.Replace(name, `\`, "/", -1))
	}

	//if len(p) > 0 && p.Merg
}
func pngData(m image.Image) []byte {
	var buf bytes.Buffer
	perr(png.Encode(&buf, m))
	return buf.Bytes()
}
func drawTerm(b []c) {
	// iterm2.com/documentation-images.html
	// github.com/mintty/mintty/blob/master/src/termout.c ...1337
	// github.com/mintty/utils/blob/master/showimg
	os.Stdout.Write([]byte{27})
	os.Stdout.Write([]byte("]1337;File=:"))
	enc := base64.NewEncoder(base64.StdEncoding, os.Stdout)
	enc.Write(b)
	enc.Close()
	os.Stdout.Write([]byte{7})
}
