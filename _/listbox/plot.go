package main

import (
	"fmt"

	"github.com/ktye/plot"
)

// ktye/plot interface for k.

func plot1(x uint32) uint32 {
	plts, e := plot.KTablePlot(x, C, I, F)
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
	O.Plots = plts
	return mki(1)
}
