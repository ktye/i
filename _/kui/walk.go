package main

import (
	"fmt"
	"strings"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var tags *walk.LineEdit
var list *walk.ListBox
var edit *walk.TextEdit

var mwin *walk.MainWindow
var wfnt *walk.Font
var listpanel, editpanel *walk.Composite
var family = "Iosevka Term SS08"
var fontsize int = 18
var history []string
var histidx int
var commit func()

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
			Composite{
				AssignTo: &listpanel,
				Visible:  true,
				Layout:   VBox{MarginsZero: true, SpacingZero: true},
				Children: []Widget{
					LineEdit{Text: "a:`abc`b`c!(1 2;3;`x`y`z!(2 3;4 5.;(1 2;3 4)))", AssignTo: &tags, Background: bg, OnKeyDown: tagsKeyDownEval, OnMouseUp: tagsMouseUp},
					ListBox{Model: []string{"k.w(go) " + today[1:]}, AssignTo: &list, MultiSelection: true, Background: bg, OnItemActivated: listDblClick, OnKeyDown: listKeyDown}, // HScroll: true, VScroll: true},
				},
			},
			Composite{
				AssignTo: &editpanel,
				Visible:  false,
				Layout:   VBox{MarginsZero: true, SpacingZero: true},
				Children: []Widget{
					TextEdit{Text: "", AssignTo: &edit, Background: bg, HScroll: true, VScroll: true},
					Composite{
						Layout: HBox{MarginsZero: true, SpacingZero: true},
						Children: []Widget{
							PushButton{Text: "ok", OnClicked: func() { exEdit(true) }},
							PushButton{Text: "cancel", OnClicked: func() { exEdit(false) }},
						},
					},
				},
			},
		},
	}.Create())
	mwin.MouseWheel().Attach(cntrlWheelFontSize)
	mwin.Run()
}
func show(s []string)   { list.SetModel(s); showpath() }
func showstr(s string)  { show([]string{s}) }
func evaluate(s string) { show(E(s)) } // edit.AppendText(" " + s + "\r\n" + E(s) + "\r\n") }
func showpath()         { mwin.SetTitle(fmt.Sprintf("k %s %s\n", X(root), X(path))) }

func dropfiles(f []string) {
	show(f)
	for _, p := range f {
		// b, e := ioutil.ReadFile(p)
		_ = p
		fmt.Println("todo drop")
	}
} // nyi: drop`file1`file2!(0x1234;0x5678..)
func tagsKeyDownEval(k walk.Key) {
	fmt.Println("tags key", k)
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
func listDblClick() { // double-click/enter
	defer showpath()
	sel := list.SelectedIndexes()
	n := len(sel)
	if n != 1 {
		return
	}
	v := list.Model().([]string)

	s := v[sel[0]]
	row, col := 0, 0
	if nn, ee := fmt.Sscanf(s, "k.w:%d:%d", &row, &col); nn == 2 && ee == nil {
		unroot()
		show(strings.Split(source, "\n"))
		list.SetSelectedIndexes([]int{row - 1})
		return
	}
	if root == 0 {
		return
	}
	x := lookup()
	t := tp(x)
	if t < 6 {
		ed()
		return
	}
	if t > 5 {
		path = lcat(path, mki(I(sel[0])))
		dx(x)
		x = lookup()
		if tp(x) < 6 {
			fmt.Println("edit/call+")
			rx(x)
			return
		}
		show(kstr(x))
	}
}
func listKeyDown(k walk.Key) {
	if k != walk.KeyBack {
		return
	}
	if path != 0 && nn(path) > 0 {
		p1 := X(path)
		path = drop(path, 4294967295)
		p2 := X(path)
		fmt.Printf("drop path: %s => %s\n", p1, p2)
		show(kstr(lookup()))
	}
}

func ed() {
	x := lookup()
	// n := nn(x)
	switch tp(x) {
	case 0:
	case 1, 2, 3, 4, 5:
		editor(kstring(x))
	case 6:
	case 7:
	}
}
func editor(s string) {
	editpanel.SetVisible(true)
	listpanel.SetVisible(false)
	edit.SetText(s)
}
func exEdit(ok bool) {
	editpanel.SetVisible(false)
	listpanel.SetVisible(true)
	s := edit.Text()
	x, _ := ES(strings.Join(strings.Split(s, "\r\n"), "\n"))
	if x == 0 {
		edit.SetText(s + "\r\n\r\n?")
	}
	fmt.Println("apply?", x != 0)
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
