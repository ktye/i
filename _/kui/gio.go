package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"

	"gioui.org/font/opentype"

	"gioui.org/app"
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
	editor.SetText("xyz")
	go func() {
		w := app.NewWindow(app.Title("k"), app.Size(unit.Dp(800), unit.Dp(650)))
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
				lineEditor.SetText(e.Text)
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
	editor     = new(widget.Editor)
	lineEditor = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	list = &layout.List{
		Axis: layout.Vertical,
	}
	topLabel = "abc"
)

type (
	D   = layout.Dimensions
	Ctx = layout.Context
)

func gui(gtx layout.Context, th *material.Theme) layout.Dimensions {
	for _, e := range lineEditor.Events() {
		if e, ok := e.(widget.SubmitEvent); ok {
			topLabel = e.Text
			lineEditor.SetText("")
		}
	}
	widgets := []layout.Widget{
		material.H3(th, topLabel).Layout,
		func(gtx Ctx) D {
			gtx.Constraints.Max.Y = gtx.Px(unit.Dp(200))
			return material.Editor(th, editor, "Hint").Layout(gtx)
		},
		func(gtx Ctx) D {
			e := material.Editor(th, lineEditor, "Hint")
			e.Font.Style = text.Italic
			return e.Layout(gtx)
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
	if err != nil {
		panic(err)
	}
	c := new(text.Collection)
	face, err := opentype.Parse(ttf)
	if err != nil {
		panic(fmt.Sprintf("font.ttf: %v", err))
	}
	c.Register(text.Font{Typeface: "kui"}, face)
	collection = c
}
