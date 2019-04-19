// +build ui,!js

package main

// gui version
// build with
//	go build -tags ui

import (
	"image"

	"github.com/eaburns/T/rope"
	"github.com/ktye/iv/cmd/lui/font"
	"github.com/ktye/plot"
	"github.com/ktye/ui"
	"golang.org/x/mobile/event/key"
)

var win *ui.Window

func main() {
	var interp interp
	repl := &ui.Repl{Reply: true}
	repl.Nowrap = true
	repl.SetText(rope.New(" "))
	interp.repl = repl
	repl.Interp = &interp
	interp.a = kinit()

	p := interp.a["plot"].(plot.Plot)
	p.Style.Dark = false
	interp.a["plot"] = p

	win = ui.New(nil)
	win.SetFont(font.APL385(), 20)
	win.Top.W = repl
	win.Render()

	for {
		select {
		case e := <-win.Inputs:
			win.Input(e)

		case err, ok := <-win.Error:
			if !ok {
				return
			}
			println("ui:", err.Error())
		}
	}
}

func top(w ui.Widget) { // set the top widget
	win.Top.W = w
	win.Top.Layout = ui.Dirty
	win.Top.Draw = ui.Dirty
	win.Render()
}

type interp struct {
	a    map[v]v
	repl *ui.Repl
}

func (i *interp) Eval(s string) {
	i.repl.Write([]byte{'\n'})
	x := run(s, i.a)
	if x != nil {
		s, o := x.(string)
		if !o {
			if p, o := x.(plot.Plot); o {
				i.plot(plot.Plots{p})
				s = ""
			} else if p, o := x.(plot.Plots); o {
				i.plot(p)
				s = ""
			} else {
				s = fmt(x).(string)
			}
		}
		i.repl.Write([]byte(s))
	}
	i.repl.Write([]byte{'\n', ' '})
	i.repl.Edit.MarkAddr("$")
}

type plotui struct { // plotui overwrites the ESC button to return to the repl
	ui.Plot
	save ui.Widget
}

func (p *plotui) Key(w *ui.Window, self *ui.Kid, k key.Event, m ui.Mouse, orig image.Point) (r ui.Result) {
	if k.Code == key.CodeEscape {
		r.Consumed = true
		top(p.save)
	}
	return
}

func (i *interp) plot(p plot.Plots) {
	w := plotui{save: win.Top.W}
	w.SetPlots(p)
	top(&w)
}

func (i *interp) Cancel() {}
