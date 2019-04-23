// +build ui,!js

package main

// gui version
// build with
//	go build -tags ui

import (
	"errors"
	"image"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"

	"github.com/eaburns/T/rope"
	"github.com/ktye/iv/cmd/lui/font"
	"github.com/ktye/plot"
	"github.com/ktye/ui"
	"golang.org/x/mobile/event/key"
)

var win *ui.Window
var kt map[v]v
var cnt func(v) v

func main() {
	var interp interp
	repl := &ui.Repl{Reply: true}
	repl.Nowrap = true
	repl.SetText(rope.New(" "))
	repl.Execute = plumb
	interp.repl = repl
	repl.Interp = &interp
	kt = kinit()
	cnt = kt["#:"].(func(v) v)

	p := kt["plot"].(plot.Plot)
	p.Style.Dark = false
	kt["plot"] = p

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
	repl *ui.Repl
}

func isplot(x v) (plot.Plots, bool) {
	if p, o := x.(plot.Plots); o {
		return p, true
	} else if p, o := x.(plot.Plot); o {
		return plot.Plots{p}, true
	} else if p, o := x.([]plot.Plot); o {
		return plot.Plots(p), true
	} // TODO: convert l{p…}
	return nil, false
}

func (i *interp) Eval(s string) {
	i.repl.Write([]byte{'\n'})
	x := run(s, kt)
	if x != nil {
		s, o := x.(string)
		if !o {
			if p, o := isplot(x); o {
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

func log(e *ui.Edit, err error) {
	e.Write([]byte("\n" + err.Error() + "\n"))
	e.MarkAddr("$")
}

// plumb executes a selection.
// pathname: edit file
// variable: show
func plumb(e *ui.Edit, s string) {
	if (len(s) > 0 && s[0] == '/') || (len(s) > 3 && s[1] == ':' && (s[2] == '/' || s[2] == '\\')) {
		file, line := s, 0
		if c := strings.LastIndexByte(s, ':'); c > 0 {
			if n, err := strconv.Atoi(s[c+1:]); err == nil {
				file, line = s[:c], n
			}
		}
		b, err := ioutil.ReadFile(file)
		if err != nil {
			log(e, err)
			return
		}
		save := win.Top.W
		cmd := make(map[string]func(*ui.Sam, string))
		cmd["q"] = func(sam *ui.Sam, c string) {
			top(save)
		}
		sam := ui.NewSam(win)
		sam.Commands = cmd
		adr := strconv.Itoa(line)
		if line > 0 {
			adr += " 0"
		}
		sam.Cmd.SetText(rope.New(adr + " $ q\n"))
		sam.Edt.SetText(rope.New(string(b)))
		top(sam)
		if line > 0 {
			sam.Edt.MarkAddr(strconv.Itoa(line))
		}
		return
	}
	show(e, s)
}

func show(e *ui.Edit, s string) {
	x := run(s, kt)
	if _, o := isplot(x); o {
		log(e, errors.New("plot"))
	} else if s, o := x.(string); o {
		log(e, errors.New(s))
	}
	z := cnt(x)
	n := int(real(z.(complex128)))
	if n == 0 {
		log(e, errors.New("empty "+reflect.TypeOf(x).String()))
	} else if n == 1 {
		log(e, errors.New(fmt(x).(string)))
	} else {
		println("show", fmt(x).(string), "⍴x", n)
	}
}
