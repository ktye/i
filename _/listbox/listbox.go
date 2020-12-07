package main

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"os"
	"strconv"
	"strings"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var mainWindow *walk.MainWindow
var tagEdit *walk.LineEdit
var listBox *walk.ListBox
var textEdit *walk.TextEdit
var splitter *walk.Splitter
var canvas *walk.CustomWidget
var bitmap *walk.Bitmap
var B []byte

func main() {
	ini(16)
	fmt.Println(val(mkcs([]byte(listbox))))
	fmt.Println(val(mkcs([]byte(appinit))))
	for _, a := range os.Args[1:] {
		dropfile(a)
	}
	B = make([]byte, len(C))
	copy(B, C)

	bg := SolidColorBrush{walk.RGB(255, 255, 234)}
	mw := MainWindow{
		AssignTo:   &mainWindow,
		Title:      "listbox",
		MinSize:    Size{300, 200},
		Size:       Size{600, 400},
		Layout:     VBox{MarginsZero: true, SpacingZero: true},
		Font:       Font{Family: "Consolas", PointSize: 13},
		Background: bg,
		Children: []Widget{
			LineEdit{
				AssignTo:   &tagEdit,
				OnMouseUp:  TagClick,
				Text:       "edit list canvas listcanvas image flip",
				OnKeyPress: TagKey,
				Background: bg,
			},
			TextEdit{
				AssignTo: &textEdit,
				//Background:      bg, // walk#746
				Visible: false,
			},
			VSplitter{
				AssignTo: &splitter,
				Children: []Widget{
					CustomWidget{
						AssignTo:           &canvas,
						Paint:              Paint,
						PaintMode:          PaintBuffered,
						AlwaysConsumeSpace: false,
						Visible:            false,
					},
					ListBox{
						AssignTo:        &listBox,
						Model:           []string{"alpha", "beta"},
						MultiSelection:  true,
						OnItemActivated: push,
						OnKeyDown:       ListKey,
						Background:      bg,
					},
				},
			},
		},
	}
	fatal(mw.Create())
	setIcon()
	mainWindow.DropFiles().Attach(DropFiles)
	mainWindow.Run()
}
func setIcon() {
	if ico, err := walk.NewIconFromImage(kimg); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		mainWindow.SetIcon(ico)
	}
}
func Paint(canvas *walk.Canvas, updateBounds walk.Rectangle) error {
	if bitmap != nil {
		canvas.DrawImage(bitmap, walk.Point{0, 0})
	}
	return nil
}
func SetBitmap(m image.Image) {
	// bounds := canvas.ClientBounds()
	if bm, e := walk.NewBitmapFromImage(m); e != nil {
		return
	} else {
		old := bitmap
		bitmap = bm
		canvas.Invalidate()
		if old != nil {
			old.Dispose()
		}
	}
}
func ListKey(key walk.Key) {
	if key == walk.KeyBack {
		back()
	}
}
func TagKey(key walk.Key) {
	if key == walk.KeyReturn {
		fmt.Println("execute", tagEdit.Text())
	}
}
func TagClick(x, y int, button walk.MouseButton) {
	if button != 1 {
		return
	}
	t := tagEdit.Text()
	a, b := tagEdit.TextSelection()
	if b-a <= 0 || t == "" {
		return
	}
	s := strings.TrimSpace(t[a:b])
	if len(s) == 0 {
		return
	}
	switch s {
	/*
		case "edit":
			show(true, false, false)
		case "list":
			show(false, false, true)
		case "canvas":
			show(false, true, false)
		case "listcanvas":
			show(false, true, true)
		case "image":
			SetBitmap(kimg)
			show(false, true, false)
	*/
	case "ok":
		ok(textEdit.Text())
	case "flip":
		o := splitter.Orientation()
		if o == walk.Horizontal {
			o = walk.Vertical
		} else {
			o = walk.Horizontal
		}
		splitter.Layout().(flipper).SetOrientation(o)
	default:
		exec(s)
	}
}
func exec(x string) { disp(ktry("exec`" + x)); tag() }     // double-click word in tag bar
func push()         { disp(ktry("push " + sel())); tag() } // double click on list entry or press enter(multiple selections)
func back()         { disp(ktry("back[]")); tag() }        // go backwards(upwards) one level on ESC key
func disp(x uint32) { // display result in a listbox or the editor
	if x == 0 {
		return
	} else if x == 0xffffffff {
		EO("!")
		return
	}
	fmt.Println("disp ", x)
	switch tp(x) {
	case 1:
		EO(sk(x))
	case 6:
		LO(lk(x))
	}
}
func ok(x string) { // "ok" clicked when editing a variable
	r := ktry("ok@" + x)
	if r == 0xffffffff {
		back()
	} else {
		disp(r)
	}
}
func EO(s string) {
	textEdit.SetText(s)
	show(true, false, false)
}
func LO(m []string) {
	listBox.SetModel(m)
	show(false, false, true)
}
func tag() { tagEdit.SetText(sk(ktry("tag path"))) }
func sel() string {
	x := listBox.SelectedIndexes()
	v := make([]string, len(x))
	for i := range x {
		v[i] = strconv.Itoa(x[i])
	}
	return strings.Join(v, " ")
}
func sk(x uint32) string {
	if x == 0xffffffff {
		return ""
	}
	defer dx(x)
	if tp(x) == 1 {
		n := nn(x)
		return string(C[x+8 : x+8+n])
	}
	return ""
}
func lk(x uint32) []string {
	defer dx(x)
	if tp(x) != 6 {
		return nil
	}
	rl(x)
	n := nn(x)
	m := make([]string, int(n))
	for i := uint32(0); i < n; i++ {
		m[i] = sk(I[x>>2+2+i])
	}
	return m
}
func DropFiles(files []string) {
	for _, f := range files {
		dropfile(f)
	}
}
func dropfile(f string) {
	if strings.HasSuffix(f, ".k") {
		load(f)
	}
}
func show(t, c, l bool) {
	textEdit.SetVisible(t)
	canvas.SetVisible(c)
	listBox.SetVisible(l)
	//splitter.Children().Remove(canvas)
	//splitter.Children().Remove(listBox)
	//splitter.Children().Add(canvas)
	//splitter.Children().Add(listBox)
}
func ktry(s string) (r uint32) {
	fmt.Println("ktry", s)
	defer func() {
		if x := recover(); x != nil {
			fmt.Println("!")
			// todo backup memory
			r = 0xffffffff
			copy(C, B)
			C = C[:len(B)]
		} else {
			if B == nil || len(B) != len(C) {
				B = make([]byte, len(C))
			}
			copy(B, C)
			fmt.Println(r)
		}
	}()
	return val(mkcs([]byte(s)))
}

type flipper interface {
	SetOrientation(walk.Orientation) error
}

const kpng = `iVBORw0KGgoAAAANSUhEUgAAABAAAAAQAgMAAABinRfyAAAACVBMVEX/AAAAAAD////KksOZAAAAMElEQVR4nGJYtWrVKoYFq1ZxMSyYhkZMgxNRXAwLpmbBCDAXSRZEgAwAGQUIAAD//+QzHr+8V1EyAAAAAElFTkSuQmCC`

var kimg image.Image

func init() {
	kimg, _ = png.Decode(base64.NewDecoder(base64.StdEncoding, strings.NewReader(kpng)))
}

const appinit = "list:$`a`b`c;text:(_10)/:list;dict:`a`b`c!(1 2;list;`f`g);table:`a`b`c!(0.+!3;_97 98 99;`e`f`g);Tags:`List`Table"
