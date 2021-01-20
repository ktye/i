package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"os"

	"github.com/ktye/plot"
)

// plot floats (xy)
// plot complex (polar)
// plot (..) (multiline)
// plot dict (KTablePlot `x`y!..)
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
	case 6:
		if xn == 0 {
			return 0
		}
		p := plot.Plot{}
		if xt := tp(MI[2+x>>2]); xt == 4 {
			p.Type = plot.Polar
		}
		p.Lines = make([]plot.Line, int(xn))
		rl(x)
		for i := uint32(0); i < xn; i++ {
			p.Lines[i] = pline(MI[2+i+x>>2], int(i))
		}
		dx(x)
		return K(p)
	case 7:
		plts, e := plot.KTablePlot(x, MC, MI, MF)
		// todo: unref x?
		perr(e)
		return K(plts[0])
	}
	return dxr(x, 0)
}

type Caption map[string]interface{}

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
	if e != nil {
		panic(e)
	}
	lastCaption = &c
	w, h := screensize()
	ip, e := p.IPlots(w, h, 0)
	perr(e)
	m := plot.Image(ip, nil, w, h, 0).(*image.RGBA)
	drawTerm(pngData(m))
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
