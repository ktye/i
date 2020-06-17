package main

import (
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
	//bg, e := walk.NewSolidColorBrush(RGB(0xFF, 0xFF, 0xEA))
	//fatal(e)
	bg := SolidColorBrush{walk.RGB(0xFF, 0xFF, 0xEA)}
	fatal(MainWindow{
		AssignTo: &mwin,
		Title:    "k",
		Font:     Font{Family: family, PointSize: fontsize},
		MinSize:  Size{600, 400},
		Layout:   VBox{MarginsZero: true, SpacingZero: true},
		Children: []Widget{
			LineEdit{CueBanner: "+/0w", AssignTo: &tags, OnEditingFinished: tagsEnter, Background: bg},
			TextEdit{Text: "k.w(go) " + today[1:] + "\r\n", AssignTo: &edit, Background: bg, HScroll: true, VScroll: true},
		},
	}.Create())
	mwin.MouseWheel().Attach(cntrlWheelFontSize)
	mwin.Run()
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

var lastTag string

func tagsEnter() {
	s := tags.Text()
	if s == lastTag {
		return
	}
	lastTag = s
	edit.AppendText(" " + s + "\r\n")
	edit.AppendText(E(s) + "\r\n")
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
