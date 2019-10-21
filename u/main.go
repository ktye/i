package main

import (
	"fmt"
	"image"
	"net/http"
	"os"
	"strconv"

	"gioui.org/app"
	"gioui.org/f32"
	gkey "gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/op"
	"gioui.org/op/paint"
)

//go:generate go run gen.go

var addr string // pretend to be a webbrowser and connect to k at addr, if empty use built-in k ([ak].go)
var Img, Im0, Im1 *image.RGBA
var Win *app.Window

func main() {
	tk = readAttachment(os.Args[0])
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
		kinit(args)
	}
	go func() {
		Img = image.NewRGBA(image.Rectangle{})
		Win = app.NewWindow(app.WithTitle("u"))
		ops := new(op.Ops)
		x, y := 0, 0
		for e := range Win.Events() {
			switch e := e.(type) {
			case app.UpdateEvent:
				if Img.Bounds().Max != e.Size {
					Im0 = image.NewRGBA(image.Rectangle{Max: e.Size})
					Im1 = image.NewRGBA(image.Rectangle{Max: e.Size})
					flip()
					hs(e.Size.X, e.Size.Y)
					break
				}
				ops.Reset()
				paint.ImageOp{Img, Img.Bounds()}.Add(ops)
				paint.PaintOp{Rect: f32.Rectangle{Max: f32.Point{float32(e.Size.X), float32(e.Size.Y)}}}.Add(ops)
				Win.Update(ops)

				// minor incompatibility: gio works with Cntrl and Shift only (no alt/meta)
				// It also cannot track Ctrl-[ etc, so don't rely on them.
			case gkey.EditEvent:
				// EditEvent is needed, because key.Event does not fire for +-*/..
				// But it does not know about modifiers, so we prefer key.Event for letters to track Cntrl-a etc
				if len(e.Text) == 1 && !craZ(e.Text[0]) {
					flip()
					hk(int(e.Text[0]), 0, 0, 0)
				}
			case gkey.Event:
				shift, cntrl := int(e.Modifiers&2>>1), int(e.Modifiers&1) // uint32 (cntrl|shift<<1)
				c := specialKeys[e.Name]
				if c == 0 {
					// fmt.Println("key e.Name", e.Name, e.Modifiers, shift)
					if e.Name >= 'A' && e.Name <= 'Z' {
						c = byte(e.Name)
						if shift == 0 {
							c += 32
						}
					} else {
						break
					}
				}
				flip()
				hk(int(c), shift, 0, cntrl)
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
	}()
	app.Main()
}

// Flip backing image before k draws over it.
// That's not for double-buffering, but ImageOp does not draw, if it receives the same pointer.
func flip() {
	Im0, Im1 = Im1, Im0
	Img = Im0
}

func hk(c, shift, alt, cntrl int) {
	if addr == "" {
		set(call("k", mki(k(c)), kmod(shift, alt, cntrl)))
	} else {
		get('k', c, shift, alt, cntrl)
	}
}
func hm(b, x0, x1, y0, y1, shift, alt, cntrl int) {
	if addr == "" {
		set(call("m", l2(i2(x0, x1), i2(y0, y1)), kmod(shift, alt, cntrl)))
	} else {
		get('m', b, x0, x1, y0, y1, shift, alt, cntrl)
	}
}
func hs(w, h int) {
	if addr == "" {
		set(call("s", mki(k(w)), mki(k(h))))
	} else {
		get('s', w, h)
	}
}
func kmod(shift, alt, cntrl int) (r k) {
	r = mk(L, 3)
	m.k[2+r], m.k[3+r], m.k[4+r] = mki(k(shift)), mki(k(alt)), mki(k(cntrl))
	return r
}
func i2(x, y int) (r k) {
	r = mk(I, 2)
	m.k[2+r], m.k[3+r] = mki(k(x)), mki(k(y))
	return r
}

// specialKeys map godoc.org/gioui.org/ui/key#pkg-constants to ../a.go:/^j,:"keycode/
var specialKeys = map[rune]byte{'⌫': 8, '⏎': 13, '⌤': 13, '⎋': 27, '⌦': 46, '⇞': 14, '⇟': 15, '⇱': 16, '⇲': 17, '←': 18, '↑': 19, '↓': 21, '→': 20}

func get(c c, a ...int) {
	b := make([]byte, 2, 32)
	b[0], b[1] = '/', c
	for i := range a {
		b = append(b, ',')
		b = append(b, []byte(strconv.Itoa(a[i]))...)
	}
	u := addr + string(b) // e.g. "http://localhost/k,49,1,0,0"
	if r, err := http.Get(u); err != nil {
		fmt.Println(err)
	} else {
		defer r.Body.Close()
		if n, err := r.Body.Read(Img.Pix); n != len(Img.Pix) {
			fmt.Printf("got screen bytes: %d expected %d\n", n, len(Img.Pix))
		} else if err != nil {
			fmt.Println(err)
		} else {
			Win.Invalidate()
		}
	}
}
func set(d k) {
	t, n := typ(d)
	if t == I && int(4*n) == len(Img.Pix) { // k(Img.Bounds().Dx()*Img.Bounds().Dy()) {
		p := int(8 + d<<2)
		copy(Img.Pix, m.c[p:p+int(4*n)])
		for i := 3; i < len(Img.Pix); i += 4 {
			Img.Pix[i] = 255 // opaque
		}
		Win.Invalidate()
	} else {
		println("set: ignore frame", t, n)
	}
	dec(d)
}
