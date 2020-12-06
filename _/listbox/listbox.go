package main

// rsrc -manifest listbox.manifest -ico k.ico -o listbox.syso

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"os"
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

func main() {
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
						AssignTo:       &listBox,
						Model:          []string{"alpha", "beta"},
						MultiSelection: true,
						Background:     bg,
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
func fatal(e error) {
	if e != nil {
		panic(e)
	}
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
	case "flip":
		o := splitter.Orientation()
		if o == walk.Horizontal {
			o = walk.Vertical
		} else {
			o = walk.Horizontal
		}
		splitter.Layout().(flipper).SetOrientation(o)
	}
	fmt.Println(x, y, button, tagEdit.Text(), a, b, s)
}
func DropFiles(files []string) {
	for _, f := range files {
		fmt.Println("drop", f)
	}
}

type flipper interface {
	SetOrientation(walk.Orientation) error
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

const kpng = `
iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAAXNSR0IArs4c6QAAAARnQU1BAACx
jwv8YQUAAAAJcEhZcwAADsQAAA7EAZUrDhsAAACoSURBVDhPvZMxDsMgDEVNJ0bYGFnZuFGuzS04
ArVdU6AphCpSn4TwJ18fkIMqCNzgITOAUq9xxYevBXQoNNDY4WvAL/wnYHWlrQBjjFRnLgOstTzP
ur0MCCGAcw5yzrJyZhmQUgKttagJ9CcyVIqk5Rgj13j/gqfgmul8xDRgzIbiva/i7SPaW6htQllb
1j51uvMRW21cMQYcB0+0U92dGLR4KjefM8ATxsGGMjsZFKMAAAAASUVORK5CYII=`

var kimg image.Image

func init() {
	kimg, _ = png.Decode(base64.NewDecoder(base64.StdEncoding, strings.NewReader(kpng)))
}
