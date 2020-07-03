package main

import (
	"fmt"
	"strings"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var tags *walk.LineEdit
var list *walk.ListBox

// var edit *walk.TextEdit
var mwin *walk.MainWindow
var wfnt *walk.Font
var family = "Iosevka Term SS04"
var fontsize int = 18
var history []string
var histidx int

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
			LineEdit{Text: "`abc`b`c!(1 2;3;\"hallo\")", CueBanner: "+/0w", AssignTo: &tags, Background: bg, OnKeyDown: tagsKeyDownEval, OnMouseUp: tagsMouseUp},
			ListBox{Model: []string{"k.w(go) " + today[1:]}, AssignTo: &list, MultiSelection: true, Background: bg}, // HScroll: true, VScroll: true},
		},
	}.Create())
	mwin.MouseWheel().Attach(cntrlWheelFontSize)
	mwin.Run()
}
func show(s []string)   { list.SetModel(s) }
func showstr(s string)  { show([]string{s}) }
func evaluate(s string) { show(E(s)) } // edit.AppendText(" " + s + "\r\n" + E(s) + "\r\n") }

func dropfiles(f []string) { showstr(fmt.Sprintf("drop %v\r\n", f)) } // nyi: drop`file1`file2!(0x1234;0x5678..)
func tagsKeyDownEval(k walk.Key) {
	fmt.Println(k, walk.KeyEscape)
	if k != walk.KeyReturn {
		if len(history) > 0 {
			h := func(i int) {
				histidx += i
				histidx = (histidx%len(history) + len(history)) % len(history)
				tags.SetText(history[histidx])
			}
			if k == walk.KeyUp {
				h(1)
			}
			if k == walk.KeyDown {
				h(-1)
			}
		}
		if k == walk.KeyEscape {
			tags.SetText("")
		}
		return
	}
	s := tags.Text()
	evaluate(s)
	for _, h := range history {
		if h == s {
			return
		}
	}
	history = append(history, s)
	histidx = len(history) - 1
	tags.SetText("")
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
	walk.SetWindowFont(list.Handle(), wfnt)
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
