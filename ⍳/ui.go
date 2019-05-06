// +build ui

package main

// gui version
// build with
//	go build -tags ui

import (
	"image"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/eaburns/T/rope"
	"github.com/ktye/iv/cmd/lui/font"
	"github.com/ktye/plot"
	"github.com/ktye/ui"
	"golang.org/x/mobile/event/key"
)

var win *ui.Window
var ipr interp
var kt map[v]v
var cnt func(v) v
var atx func(v, v) v
var lnx func(v) int
var til func(v) v
var cst func(v, v) v

func main() {
	rpl := &ui.Repl{Reply: true}
	rpl.Nowrap = true
	rpl.SetText(rope.New(" "))
	rpl.Execute = plumb
	ipr.repl = rpl
	rpl.Interp = &ipr
	kt = kinit()
	cnt = kt["#:"].(func(v) v)
	atx = kt["@@"].(func(v, v) v)
	lnx = kt["ln"].(func(v) int)
	til = kt["!:"].(func(v) v)
	cst = kt["$$"].(func(v, v) v)

	p := kt["plot"].(plot.Plot)
	p.Style.Dark = false
	kt["plot"] = p

	win = ui.New(nil)
	win.SetFont(font.APL385(), 20)
	win.Top.W = rpl
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
				return
			} else {
				s = fmt(x).(string)
			}
		}
		i.repl.Write([]byte(s + "\n"))
	}
	i.repl.Edit.MarkAddr("$")
}

func setTop(w ui.Widget) { // set the top widget
	win.Top.W = w
	win.Top.Layout = ui.Dirty
	win.Top.Draw = ui.Dirty
	win.Render()
}

func push(w ui.Widget) {
	t := top{Widget: w, save: win.Top.W}
	setTop(t)
}

type top struct {
	ui.Widget
	save ui.Widget
}

func (t top) Key(w *ui.Window, self *ui.Kid, k key.Event, m ui.Mouse, orig image.Point) (res ui.Result) {
	if k.Code == key.CodeEscape && k.Direction == key.DirRelease {
		setTop(t.save)
		res.Consumed = true
		return res
	}
	return t.Widget.Key(w, self, k, m, orig)
}

func (i *interp) plot(p plot.Plots) {
	w := &ui.Plot{}
	w.SetPlots(p)
	push(w)
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
			setTop(save)
		}
		sam := ui.NewSam(win)
		sam.Commands = cmd
		adr := strconv.Itoa(line)
		if line > 0 {
			adr += " 0"
		}
		sam.Cmd.SetText(rope.New(adr + " $ q\n"))
		sam.Edt.SetText(rope.New(string(b)))
		setTop(sam)
		if line > 0 {
			sam.Edt.MarkAddr(strconv.Itoa(line))
		}
		return
	}
	show(s)
}

func show(s string) {
	x := run(s, kt)
	if p, o := isplot(x); o {
		ipr.plot(p)
		return
	}
	tr := tree{x: x}
	if tr.Leaf() {
		ipr.repl.Write([]byte(fmt(x).(string) + "\n"))
		ipr.repl.MarkAddr("$")
		return
	}
	t := &ui.Tree{}
	t.Single = true
	t.SetRoot(&tr)
	push(t)
}

type tree struct {
	x v
	s string
	c []string
}

func (t *tree) String() string {
	if t.s != "" {
		return t.s
	}
	return fmt(t.x).(s)
}
func (t *tree) Count() int {
	r := int(real(cnt(t.x).(complex128)))
	if lnx(t.x) < 0 && r != 1 { // dict
		d := [2]l{l{"d", "q"}, l{complex(1, 0), complex(1, 0)}}
		f := cst(d, t.x).(string)
		t.c = strings.Split(f, "\n")
	}
	return r
}
func (t *tree) Leaf() bool { return t.Count() == 1 && lnx(t.x) < 0 }
func (t *tree) Child(i int) ui.Plant {
	var y v = complex(float64(i), 0)
	var s = ""
	if lnx(t.x) < 0 { // dict
		keys := til(t.x)
		y = atx(keys, y)
		if i < len(t.c) {
			s = t.c[i]
		}
	}
	v := atx(t.x, y)
	return &tree{x: v, s: s}
}
