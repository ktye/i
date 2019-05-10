// +build ui

package main

// gui version
// build with
//	go build -tags ui

import (
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/eaburns/T/rope"
	"github.com/ktye/plot"
	"github.com/ktye/ui"
	"github.com/ktye/ui/base"
	"github.com/ktye/ui/dpy"
	"github.com/ktye/ui/editor"
	"github.com/ktye/ui/fonts/apl385"
	"golang.org/x/exp/shiny/screen"
)

var win *ui.Window
var rpl *editor.Repl
var edt *editor.Edit
var cnv ui.Widget
var sp1, sp2 *base.Split

var ipr interp
var kt map[v]v
var cnt func(v) v
var atx func(v, v) v
var lnx func(v) int
var til func(v) v
var cst func(v, v) v

func main() {
	base.SetFont(apl385.TTF(), 20)
	rpl = &editor.Repl{Reply: true, Prompt: " "}
	rpl.Edit.SetText(rope.New(" "))
	rpl.Execute = func(_ *editor.Edit, s string) int { rpl.DefaultExec(nil, s); return -1 }
	rpl.Nowrap = true
	rpl.Interp = &ipr
	rpl.Menu = rplmenu(rpl)

	edt = editor.New("")
	edt.Nowrap = true
	edt.Menu = edt.StandardMenu()
	evlbutton(edt)
	dotbutton(edt)
	cnv = &base.Blank{}
	sp2 = base.NewSplit(edt, cnv)
	sp1 = base.NewSplit(rpl, sp2)
	sp1.Vertical = true
	sp1.Ratio = 1

	kt = kinit()
	cnt = kt["#:"].(func(v) v)
	atx = kt["@@"].(func(v, v) v)
	lnx = kt["ln"].(func(v) int)
	til = kt["!:"].(func(v) v)
	cst = kt["$$"].(func(v, v) v)

	p := kt["plot"].(plot.Plot)
	p.Style.Dark = false
	kt["plot"] = p

	win = ui.New(dpy.New(&screen.NewWindowOptions{Title: "i"})) // win7 confuses iota and quad.

	win.Top = &base.Scale{Widget: sp1, Funcs: []func(){plotfont}}
	done := win.Run()
	<-done
}

type interp struct{}

func isplot(x v) (plot.Plots, bool) {
	if p, o := x.(plot.Plots); o {
		return p, true
	} else if p, o := x.(plot.Plot); o {
		return plot.Plots{p}, true
	} else if p, o := x.([]plot.Plot); o {
		return plot.Plots(p), true
	} else if u, o := x.(l); o {
		for i := range u {
			if _, o := u[i].(plot.Plot); !o {
				return nil, false
			}
		}
		p = make(plot.Plots, len(u))
		for i := range u {
			p[i] = u[i].(plot.Plot)
		}
		return p, true
	}
	return nil, false
}

func (i *interp) Eval(s string) string {
	s = plumb(s)
	x := run(s, kt)
	if x != nil {
		s, o := x.(string)
		if !o {
			if p, o := isplot(x); o {
				setplot(p)
				return ""
			} else {
				s = fmt(x).(string)
			}
		}
		return s
	}
	return ""
}

func rplmenu(r *editor.Repl) *base.Menu {
	kval := func() v {
		s := r.Selection()
		if s == "" {
			return nil
		}
		x := run(s, kt)
		if _, o := isplot(x); o {
			return nil
		}
		return x
	}
	edit := base.NewButton("edit", "", func() int {
		var t string
		x := kval()
		if str, o := x.(s); o {
			t = str
		} else {
			t = fmt(x).(s)
		}
		edit(t)
		return -1
	})
	show := base.NewButton("show", "", func() int {
		println("TODO show")
		return 0
	})
	m := r.StandardMenu()
	m.Buttons = append(m.Buttons, edit, show)
	return m
}
func evlbutton(e *editor.Edit) { // add a run menu entry to the editor
	b := base.NewButton("eval", "", func() int {
		rpl.Execute(nil, e.Selection())
		return -1
	})
	e.Menu.Buttons = append([]*base.Button{b}, e.Menu.Buttons...)
}
func dotbutton(e *editor.Edit) { // assign editor content or selection to buf variable
	b := base.NewButton("dot:←", "", func() int {
		s := e.Selection()
		if len(s) == 0 {
			s = e.Text().String()
		}
		kt["dot"] = s
		return 0
	})
	e.Menu.Buttons = append(e.Menu.Buttons, b)
}

// plumb intercepts execute.
// pathname: dirname: list files in the repl, filename: show file in the editor.
// variable: show in repl, or as a tree in the canvas.
// otherwise: return the input.
func plumb(s string) string {
	s = strings.TrimSpace(s)
	if (len(s) > 0 && s[0] == '/') || (len(s) > 2 && s[1] == ':' && (s[2] == '/' || s[2] == '\\')) {
		if fi, err := ioutil.ReadDir(s); err == nil {
			dir := s
			for _, f := range fi {
				s = filepath.Join(dir, f.Name())
				if f.IsDir() {
					s += "/"
				}
				rpl.Write([]byte(s + "\n"))
			}
			return ""
		}
		file, line := s, 0
		if c := strings.LastIndexByte(s, ':'); c > 0 {
			if n, err := strconv.Atoi(s[c+1:]); err == nil {
				file, line = s[:c], n
			}
		}
		b, err := ioutil.ReadFile(file)
		if err == nil {
			edit(string(b))
			if line > 0 {
				edt.MarkAddr(strconv.Itoa(line))
			}
			return ""
		}
	}
	switch s {
	case `\c`: // clear terminal
		rpl.SetText(rope.New(""))
	case `\h`:
		return "doc"
	case `\v`:
		println("TODO list vars")
	default:
		return s
	}
	return ""
}
func edit(t s) {
	edt.SetText(rope.New(t))
	if sp1.Ratio > 0.95 {
		sp1.Ratio = 0
	}
}
func setplot(p plot.Plots) {
	cnv = plot.NewUI(p)
	if sp1.Ratio > 0.95 {
		sp1.Ratio = 0
	}
	if sp2.Ratio > 0.95 {
		sp2.Ratio = 0
	}
	sp2.Kids[1].Widget = cnv
}
func plotfont() {
	s1 := base.Font.Size()
	s2 := (s1 * 8) / 10
	if s1 < 6 {
		s1 = 6
	}
	if s2 < 6 {
		s2 = 6
	}
	f1 := base.LoadFace(base.Font.TTF, s1)
	f2 := base.LoadFace(base.Font.TTF, s2)
	plot.SetFonts(f1, f2)
}

/*
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
*/

/* TODO: port tree to v2
// plumb executes.
// pathname: dirname: list files in the repl, filename: show file in the editor.
// variable: show in repl, or as a tree in the canvas.
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
*/
