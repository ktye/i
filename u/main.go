package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"strconv"

	"gioui.org/ui/app"
	"gioui.org/ui/f32"
	gkey "gioui.org/ui/key"
	"gioui.org/ui/layout"
	"gioui.org/ui/paint"
	"gioui.org/ui/pointer"
)

var screen *image.RGBA
var netescape bool // pretend to be a webbrowser and connect to k at addr, otherwise use built-in k

func main() {
	args, addr := os.Args[1:], ""
	if len(args) > 0 {
		if args[0] == "-p" {
			args = args[1:]
			netescape = true
			addr = ":2019"
			if len(args) > 0 {
				addr = args[0]
				args = args[1:]
			}
		}
	}
	if netescape {
		connect(addr)
	} else {
		kinit()
	}
	// TODO: evl files given at remaining args
	go func() {
		w := app.NewWindow(app.WithTitle("u"))
		if err := loop(w); err != nil {
			panic(err)
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

func hkey(c c, shift, alt, cntrl int) {
	if netescape {
		get(c, shift, alt, cntrl)
	} else {
		// kxy(mks("key"), ...
	}
}

func loop(w *app.Window) error {
	c := &layout.Context{
		Queue: w.Queue(),
	}
	for {
		e := <-w.Events()
		switch e := e.(type) {
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

		// minor incompatibility: gio works with Cntrl and Shift only (no alt/meta)
		// It also cannot track Ctrl-[ etc, so don't rely on them.
		case gkey.EditEvent:
			// EditEvent is needed, because key.Event does not fire for +-*/..
			// But it does not know about modifiers, so we prefer key.Event for letters to track Cntrl-a etc
			if len(e.Text) == 1 && !craZ(e.Text[0]) {
				hkey(e.Text[0], 0, 0, 0)
			}
		case gkey.Event:
			shift, cntrl := int(e.Modifiers&2>>1), int(e.Modifiers&1) // uint32 (cntrl|shift<<1)
			c := specialKeys[e.Name]
			if c == 0 {
				if e.Name >= 'A' && e.Name <= 'Z' {
					c = byte(e.Name)
					if shift != 0 {
						c += 32
					}
				} else {
					break
				}
			}
			hkey(c, shift, 0, cntrl)
		case pointer.Event:
			fmt.Printf("mouse: %+v\n", e)

		}
	}
}

// keys map godoc.org/gioui.org/ui/key#pkg-constants to ../a.go:/^j,:"keycode/
var specialKeys = map[rune]byte{'⌫': 8, '⏎': 13, '⌤': 13, '⎋': 27, '⌦': 46, '⇞': 14, '⇟': 15, '⇱': 16, '⇲': 17, '←': 18, '↑': 19, '↓': 20, '→': 21}

func connect(addr string) {
}
func get(c c, a ...int) {
	b := make([]byte, 2, 32)
	b[0], b[1] = '/', c
	for i := range a {
		b = append(b, ',')
		b = append(b, []byte(strconv.Itoa(a[i]))...)
	}
	addr := string(b)         // e.g. "/k,49,1,0,0"
	fmt.Println("get ", addr) // TODO
}
