// +build ui,!js

package main

// gui version
// build with
//	go build -tags ui

import (
	"github.com/eaburns/T/rope"
	"github.com/ktye/iv/cmd/lui/font"
	"github.com/ktye/ui"
)

func main() {
	var interp interp
	repl := &ui.Repl{Reply: true}
	repl.Nowrap = true
	repl.SetText(rope.New(" "))
	interp.repl = repl
	repl.Interp = &interp
	interp.a = kinit()

	w := ui.New(nil)
	w.SetFont(font.APL385(), 20)
	w.Top.W = repl
	w.Render()

	for {
		select {
		case e := <-w.Inputs:
			w.Input(e)

		case err, ok := <-w.Error:
			if !ok {
				return
			}
			println("ui:", err.Error())
		}
	}
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
			s = fmt(x).(string)
		}
		i.repl.Write([]byte(s))
	}
	i.repl.Write([]byte{'\n', ' '})
	i.repl.Edit.MarkAddr("$")
}

func (i *interp) Cancel() {}
