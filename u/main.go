package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"net/http"
	"os"
	"strconv"

	"gioui.org/app"
	"gioui.org/f32"
	gkey "gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/paint"
)

/*
func writePng() {
	println("w screen.png")
	w, err := os.Create("screen.png")
	if err != nil {
		panic(err)
	}
	defer w.Close()
	png.Encode(w, screen)
}
*/

var screen *image.RGBA
var addr string // pretend to be a webbrowser and connect to k at addr, if empty use built-in k ([ak].go)
var window *app.Window

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		if args[0] == "-i" {
			args = args[1:]
			addr = "http://localhost:2019"
			if len(args) > 0 {
				addr = args[0]
				args = args[1:]
			}
		}
	}
	if addr == "" {
		kinit()
	}
	// TODO: evl files given at remaining args
	go func() {
		screen = image.NewRGBA(image.Rectangle{})
		window = app.NewWindow(app.WithTitle("u"))
		if err := loop(window); err != nil {
			panic(err)
		}
	}()
	app.Main()
}

func hk(c c, shift, alt, cntrl int) {
	fmt.Println("key", int(c), shift, alt, cntrl)
	if addr == "" {
		// kxy(mks("key"), ...
	} else {
		get('k', int(c), shift, alt, cntrl)
	}
}
func hm(b, x0, x1, y0, y1, shift, alt, cntrl int) {
	fmt.Println("mouse", b, x0, x1, y0, y1, shift, alt, cntrl)
	if addr == "" {
	} else {
		get('m', b, x0, x1, y0, y1, shift, alt, cntrl)
	}
}
func hs(w, h int) {
	fmt.Println("screen", w, h)
	r := screen.Bounds()
	r.Max.X /= 2
	r.Max.Y /= 2
	draw.Draw(screen, r, &image.Uniform{color.RGBA{255, 0, 0, 255}}, image.ZP, draw.Src)
	if addr == "" {
	} else {
		get('s', w, h)
	}
}

func loop(w *app.Window) error {
	ctx := &layout.Context{
		Queue: w.Queue(),
	}
	x, y := 0, 0
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.UpdateEvent:
			if screen.Bounds().Max != e.Size {
				println("update size")
				screen = image.NewRGBA(image.Rectangle{Max: e.Size})
				hs(e.Size.X, e.Size.Y)
				break
			}
			fmt.Printf("update onscreen: e=%+v\n", e)
			// writePng()
			ctx.Reset(&e.Config, layout.RigidConstraints(e.Size))
			paint.ImageOp{screen, screen.Bounds()}.Add(ctx.Ops)
			paint.PaintOp{Rect: f32.Rectangle{Max: f32.Point{float32(e.Size.X), float32(e.Size.Y)}}}.Add(ctx.Ops)
			w.Update(ctx.Ops)

		// minor incompatibility: gio works with Cntrl and Shift only (no alt/meta)
		// It also cannot track Ctrl-[ etc, so don't rely on them.
		case gkey.EditEvent:
			// EditEvent is needed, because key.Event does not fire for +-*/..
			// But it does not know about modifiers, so we prefer key.Event for letters to track Cntrl-a etc
			if len(e.Text) == 1 && !craZ(e.Text[0]) {
				hk(e.Text[0], 0, 0, 0)
			}
		case gkey.Event:
			shift, cntrl := int(e.Modifiers&2>>1), int(e.Modifiers&1) // uint32 (cntrl|shift<<1)
			c := specialKeys[e.Name]
			if c == 0 {
				fmt.Println("key e.Name", e.Name, e.Modifiers, shift)
				if e.Name >= 'A' && e.Name <= 'Z' {
					c = byte(e.Name)
					if shift == 0 {
						c += 32
					}
				} else {
					break
				}
			}
			hk(c, shift, 0, cntrl)
		case pointer.Event:
			xxx := 0 // TODO: track modifiers
			if e.Scroll.Y != 0 {
				if e.Scroll.Y < 0 {
					hm(3, 0, 0, 0, 0, xxx, xxx, xxx)
				} else {
					hm(4, 0, 0, 0, 0, xxx, xxx, xxx)
				}
				break
			}
			if e.Type == pointer.Press {
				x, y = int(e.Position.X), int(e.Position.Y)
			} else if e.Type == pointer.Release {
				hm(0, x, int(e.Position.X), y, int(e.Position.Y), xxx, xxx, xxx)
			}
		}
	}
}

// specialKeys map godoc.org/gioui.org/ui/key#pkg-constants to ../a.go:/^j,:"keycode/
var specialKeys = map[rune]byte{'⌫': 8, '⏎': 13, '⌤': 13, '⎋': 27, '⌦': 46, '⇞': 14, '⇟': 15, '⇱': 16, '⇲': 17, '←': 18, '↑': 19, '↓': 20, '→': 21}

func get(c c, a ...int) {
	b := make([]byte, 2, 32)
	b[0], b[1] = '/', c
	for i := range a {
		b = append(b, ',')
		b = append(b, []byte(strconv.Itoa(a[i]))...)
	}
	u := addr + string(b) // e.g. "http://localhost/k,49,1,0,0"
	fmt.Println("get ", u)
	if r, err := http.Get(u); err != nil {
		fmt.Println(err)
	} else {
		defer r.Body.Close()
		if n, err := r.Body.Read(screen.Pix); n != len(screen.Pix) {
			fmt.Printf("got screen bytes: %d expected %d\n", n, len(screen.Pix))
		} else if err != nil {
			fmt.Println(err)
		} else {
			window.Invalidate()
		}
	}
}
