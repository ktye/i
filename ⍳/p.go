// +build !js

package main

import "github.com/ktye/plot"

func regplot(a map[v]v) {
	a["plot"] = plot.Plot{}
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
