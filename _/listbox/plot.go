package main

import (
	"github.com/ktye/plot"
)

// ktye/plot interface for k.

func plot1(x uint32) uint32 {
	plts, e := plot.KTablePlot(x, C, I, F)
	if e != nil {
		panic(e)
	}
	// todo setplot
	pltui.SetPlot(plts, nil)
	return 0
}
