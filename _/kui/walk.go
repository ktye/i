package main

import (
	"fmt"
	"strings"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var tags *walk.LineEdit
var edit *walk.TextEdit
var mwin *walk.MainWindow
var wfnt *walk.Font
var family = "Iosevka Term SS04"
var fontsize int = 18

func gui() {
	bg := SolidColorBrush{walk.RGB(0xFF, 0xFF, 0xEA)}
	fatal(MainWindow{
		AssignTo:    &mwin,
		Title:       "k",
		Font:        Font{Family: family, PointSize: fontsize},
		MinSize:     Size{600, 400},
		Layout:      VBox{MarginsZero: true, SpacingZero: true},
		OnDropFiles: dropfiles,
		Children: []Widget{
			LineEdit{Text: "File Select Plot Help 2 3⍴⍳5", CueBanner: "+/0w", AssignTo: &tags, Background: bg, OnKeyDown: tagsKeyDownEval, OnMouseUp: tagsMouseUp},
			TextEdit{Text: "k.w(go) " + today[1:] + "\r\n", AssignTo: &edit, Background: bg, HScroll: true, VScroll: true},
		},
	}.Create())
	mwin.MouseWheel().Attach(cntrlWheelFontSize)
	mwin.Run()
}

func evaluate(s string) { edit.AppendText(" " + s + "\r\n" + E(s) + "\r\n") }

func dropfiles(f []string) { edit.AppendText(fmt.Sprintf("drop %v\r\n", f)) } // nyi: drop`file1`file2!(0x1234;0x5678..)
func tagsKeyDownEval(k walk.Key) {
	if k != walk.KeyReturn {
		return
	}
	evaluate(tags.Text())
}
func tagsMouseUp(x, y int, b walk.MouseButton) {
	getSelectedText := func(a, b int) string { r := []rune(tags.Text()); return string(r[a:b]) }
	if b == walk.LeftButton {
		aa, bb := tags.TextSelection()
		if bb-aa > 0 {
			evaluate(getSelectedText(aa, bb))
		}
	}
	if b == walk.RightButton { //nyi: prevent contextmenu. Overwrite WndProc? WM_CONTEXTMENU
		aa, bb := tags.TextSelection() //
		if bb-aa > 0 {
			fmt.Println("selected text: ", getSelectedText(aa, bb))
		}
	}
}
func cntrlWheelFontSize(x, y int, button walk.MouseButton) {
	delta := walk.MouseWheelEventDelta(button)
	keyState := walk.MouseWheelEventKeyState(button)
	if keyState == 8 { // cntrl-mouse:
		fontsize++
		if delta < 0 {
			fontsize -= 2
			if fontsize < 10 {
				fontsize = 10
			}
		}
	}
	if wfnt != nil {
		wfnt.Dispose()
	}
	wfnt, _ := walk.NewFont(family, fontsize, 0)
	walk.SetWindowFont(tags.Handle(), wfnt)
	walk.SetWindowFont(edit.Handle(), wfnt)
	// layout? xxx.RequestLayout() ?
}

func todos(s string) string {
	return strings.Replace(strings.Replace(s, "\r", "", -1), "\n", "\r\n", -1)
}

/*
What's wrong with this one?
Nothing Turkish.
It's tip top.
I'm just not sure about the color.
*/
