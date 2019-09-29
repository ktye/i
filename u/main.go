package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"

	"gioui.org/ui/app"
	"gioui.org/ui/f32"
	"gioui.org/ui/key"
	"gioui.org/ui/layout"
	"gioui.org/ui/pointer"
	"gioui.org/ui/paint"
)

var screen *image.RGBA

func main() {
	go func() {
		w := app.NewWindow(app.WithTitle("u"))
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
	}()
	app.Main()
}

func flush() {
	println("paint")
	r := screen.Bounds()
	r.Max.X /= 2
	r.Max.Y /= 2
	draw.Draw(screen, r, &image.Uniform{color.RGBA{255, 0, 0, 255}}, image.ZP, draw.Src)
}

func loop(w *app.Window) error {
	c := &layout.Context{
		Queue: w.Queue(),
	}
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case key.Event:
			fmt.Printf("key: %v\n", e)
		case pointer.Event:
			fmt.Printf("mouse: %v\n", e)
		case app.DestroyEvent:
			return e.Err
		case app.UpdateEvent:
			fmt.Printf("update %+v\n", e)
			if screen == nil || screen.Bounds().Max != e.Size {
				screen = image.NewRGBA(image.Rectangle{Max: e.Size})
				flush()
			}
			c.Reset(&e.Config, layout.RigidConstraints(e.Size))
			paint.ImageOp{screen, screen.Bounds()}.Add(c.Ops)
			paint.PaintOp{Rect: f32.Rectangle{Max: f32.Point{float32(e.Size.X), float32(e.Size.Y)}}}.Add(c.Ops)
			w.Update(c.Ops)
		}
	}
}
