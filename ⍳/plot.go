package main

import (
	"github.com/golang/freetype/truetype"
	"github.com/ktye/plot"
	"github.com/ktye/ui/fonts/apl385"
	"golang.org/x/image/font"
)

func regplot(a map[v]v) {
	p := plot.Plot{}
	p.Style.Dark = true
	a["plot"] = p
	a["line"] = plot.Line{}
	a["plots"] = plots
}

// plots converts a list of plots to plot.Plots which results in multiple plots next to each other.
func plots(x v) v {
	l, o := x.(l)
	if !o {
		panic("type")
	}
	plts := make(plot.Plots, len(l))
	for i := range l {
		p, o := l[i].(plot.Plot)
		if !o {
			panic("type")
		}
		plts[i] = p
	}
	return plts
}

func init() {
	ttf, err := truetype.Parse(apl385.TTF())
	if err != nil {
		panic(err)
	}
	face := func(size int) font.Face {
		opt := truetype.Options{
			Size: float64(size),
			DPI:  72,
		}
		return truetype.NewFace(ttf, &opt)
	}
	font1, font2 := face(20), face(16)
	plot.SetFonts(font1, font2)
}
