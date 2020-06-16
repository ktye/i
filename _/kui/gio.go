package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/font/opentype"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type scaledConfig struct {
	Scale float32
}

func gio() {
	flag.Parse()
	editor.SetText("k.w(go) " + today[1:] + " +/∞\n")
	go func() {
		w := app.NewWindow(app.Title("k"), app.Size(unit.Dp(640), unit.Dp(480)))
		if e := loop(w); e != nil {
			fatal(e)
		}
	}()
	app.Main()
}

func loop(w *app.Window) error {
	th := material.NewTheme(collection)
	var ops op.Ops
	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case system.ClipboardEvent:
				tags.SetText(e.Text)
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)

				gui(gtx, th)
				e.Frame(gtx.Ops)
			}
		}
	}
}

var (
	editor = new(widget.Editor)
	tags   = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	list = &layout.List{
		Axis: layout.Vertical,
	}
)

type (
	D   = layout.Dimensions
	Ctx = layout.Context
)

func gui(gtx layout.Context, th *material.Theme) layout.Dimensions {
	for _, e := range tags.Events() {
		if e, ok := e.(widget.SubmitEvent); ok {
			tags.SetText("")
			editor.Move(12345678) // end?
			editor.Insert(" " + e.Text + "\n")
			editor.Insert(E(e.Text))
			editor.Insert("\n ")
		}
	}
	widgets := []layout.Widget{
		func(gtx Ctx) D {
			e := material.Editor(th, tags, "+/0w")
			return e.Layout(gtx)
		},
		func(gtx Ctx) D {
			gtx.Constraints.Max.Y = gtx.Px(unit.Dp(200))
			return material.Editor(th, editor, "Hint").Layout(gtx)
		},
	}

	return list.Layout(gtx, len(widgets), func(gtx Ctx, i int) D {
		return layout.UniformInset(unit.Dp(16)).Layout(gtx, widgets[i])
	})
}

func (s *scaledConfig) Px(v unit.Value) int {
	scale := s.Scale
	if v.U == unit.UnitPx {
		scale = 1
	}
	return int(math.Round(float64(scale * v.V)))
}

var collection *text.Collection

func init() {
	ttf, err := ioutil.ReadFile("font.ttf")
	if err == nil {
		c := new(text.Collection)
		face, err := opentype.Parse(ttf)
		if err != nil {
			panic(fmt.Sprintf("font.ttf: %v", err))
		}
		c.Register(text.Font{Typeface: "kui"}, face)
		collection = c
	} else {
		collection = gofont.Collection()
	}
}
