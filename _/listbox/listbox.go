package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/ktye/plot"
	"github.com/ktye/plot/plotui"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

const (
	EDIT = 1 << iota
	LIST
	PLOT
)

var Files fs
var Src []string
var O output

var mainWindow *walk.MainWindow
var tagEdit *walk.LineEdit
var listBox *walk.ListBox
var textEdit *walk.TextEdit
var splitter *walk.Splitter
var pltui plotui.Plot
var B []byte //k-backup memory

func main() {
	T[99] = csv2   // 'c[x;y]
	T[240] = plot1 // 'px
	restart(listbox)
	DropFiles(os.Args[1:])
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
				Text:       "Src",
				OnKeyPress: TagKey,
				Background: bg,
			},
			TextEdit{
				AssignTo:    &textEdit,
				Background:  bg,
				OnMouseDown: TextClick,
				ContextMenuItems: []MenuItem{
					Action{Text: "Back", OnTriggered: back},
					Action{Text: "Search", OnTriggered: Search},
					Action{Text: "Save/Reload (lb.k)", OnTriggered: Save},
				},
				HScroll: true,
				VScroll: true,
				Visible: true,
			},
			ListBox{
				AssignTo:        &listBox,
				Model:           []string{"alpha", "beta"},
				MultiSelection:  true,
				OnItemActivated: push,
				OnKeyDown:       ListKey,
				Background:      bg,
				Visible:         false,
				ContextMenuItems: []MenuItem{
					Action{Text: "Back", OnTriggered: back},
					Action{Text: "Select", OnTriggered: push},
				},
			},
			VSplitter{
				AssignTo:      &splitter,
				Visible:       false,
				StretchFactor: 2,
				Children: []Widget{
					pltui.BuildPlot(nil),
					Composite{
						Layout: VBox{MarginsZero: true, SpacingZero: true},
						Children: []Widget{
							pltui.BuildSlider(),
							pltui.BuildCaption(nil),
						},
					},
				},
			},
		},
	}
	fatal(mw.Create())
	setIcon()
	mainWindow.DropFiles().Attach(DropFiles)
	tag()
	mainWindow.Run()
}
func restart(s string) {
	C = make([]byte, 1<<16)
	msl()
	ini(16)
	initdat()
	B = make([]byte, len(C))
	copy(B, C)
	Src = strings.Split(s, "\n")
	ktry(s)
	for _, f := range Files.l {
		kdrop(f, Files.m[f])
	}
	tag()
}
func setIcon() {
	if ico, err := walk.NewIconFromImage(kimg); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		mainWindow.SetIcon(ico)
	}
}
func ListKey(key walk.Key) {
	if key == walk.KeyBack {
		back()
	}
}
func TagKey(key walk.Key) {
	if key == walk.KeyReturn {
		s := tagEdit.Text()
		if i := strings.Index(s, " k)"); i != -1 {
			s = s[i+3:]
		}
		kdo(func() { dx(out(val(mkcs([]byte(s))))) })
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
	case "Src":
		setMenu(textEdit, "back", "Search", "Save")
		textEdit.SetText(strings.Join(Src, "\r\n"))
		show(EDIT)
	case "ok":
		ok(textEdit.Text())
	default:
		if strings.HasPrefix(s, "k)") {
			setMenu(textEdit, "Search")
			kdo(func() { dx(out(val(mkcs([]byte(s[2:]))))) })
		} else {
			exec(s)
		}
	}
}
func setMenu(w walk.Widget, items ...string) {
	m := make(map[string]bool)
	for _, s := range items {
		m[s] = true
	}
	al := w.ContextMenu().Actions()
	for i := 0; i < al.Len(); i++ {
		a := al.At(i)
		t := a.Text()
		v := false
		for _, p := range items {
			if strings.HasPrefix(t, p) {
				v = true
			}
		}
		a.SetVisible(v)
	}
}
func TextClick(x, y int, button walk.MouseButton) {
	if button == 2 {
		m := textEdit.ContextMenu()
		al := m.Actions()
		for i := 0; i < al.Len(); i++ {
			fmt.Println(al.At(i).Text())
		}
	}
}
func Search() {
	t := textEdit.Text()
	a, b := textEdit.TextSelection()
	if a >= 0 && a < len(t) && b >= 0 && b <= len(t) && b > a {
		s := t[a:b]
		i := strings.Index(t[a+1:], s)
		if i == -1 {
			i = strings.Index(t, s)
		} else {
			i += 1 + a
		}
		if i > 0 {
			textEdit.SetTextSelection(i, i+len(s))
		}
	}
}
func Save() {
	t := strings.Replace(textEdit.Text(), "\r", "", -1)
	if e := ioutil.WriteFile("lb.k", []byte(t), 644); e != nil {
		fmt.Println(e)
		return
	}
	restart(t)
}
func exec(x string) { disp(ktry("exec`" + x)); tag() }     // double-click word in tag bar
func push()         { disp(ktry("push " + sel())); tag() } // double click on list entry or press enter(multiple selections)
func back()         { disp(ktry("back[]")); tag() }        // go backwards(upwards) one level on ESC key
func disp(x uint32) { // display result in a listbox or the editor
	if x == 0 {
		return
	} else if x == 0xffffffff {
		//EO("!")
		return
	}
	fmt.Println("disp ", x)
	t := tp(x)
	switch t {
	case 1:
		EO(sk(x))
	case 6:
		LO(lk(x))
	default:
		fmt.Println("disp: type", t)
	}
}
func printc(x, y uint32) {
	s := string(C[x : x+y])
	fmt.Println(s)
	fmt.Fprintf(&O, "%s\n", s)
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
	show(EDIT)
}
func LO(m []string) {
	listBox.SetCurrentIndex(-1)
	fmt.Println("model", m)
	listBox.SetModel(m)
	show(LIST)
}
func PO(p plot.Plots) { pltui.SetPlot(p, nil); show(PLOT) }
func tag() {
	if mainWindow != nil {
		mainWindow.SetTitle(sk(ktry(" 'kpath")))
	}
	if tagEdit != nil {
		tagEdit.SetText(sk(ktry("tag path")))
	}
}
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
		if strings.HasSuffix(f, ".k") {
			if b, e := ioutil.ReadFile(f); e == nil {
				if n := bytes.Index(b, []byte("\n\\")); n != -1 {
					b = b[:n+1]
				}
				Src = strings.Split(string(b), "\n")
			}
		} else {
			Files.add(f)
		}
	}
}
func show(flag int) {
	textEdit.SetVisible(flag&EDIT != 0)
	listBox.SetVisible(flag&LIST != 0)
	splitter.SetVisible(flag&PLOT != 0)
}
func ktry(s string) (r uint32) {
	O.Reset()
	fmt.Println("ktry", s)
	defer func() {
		if x := recover(); x != nil {
			debug.PrintStack()
			fmt.Fprintf(&O, "%s\r\n^\r\n%s\r\n", s, x)
			show(EDIT)
			r = 0xffffffff
			restore()
		} else {
			backup()
		}
		if O.Len() > 0 {
			EO(string(bytes.Replace(O.Bytes(), []byte{'\n'}, []byte("\r\n"), -1)))
			r = 0
		}
		if len(O.Plots) > 0 {
			PO(O.Plots)
			r = 0
		}
	}()
	return val(mkcs([]byte(s)))
}
func kdo(f func()) {
	O.Reset()
	defer func() {
		if x := recover(); x != nil {
			debug.PrintStack()
			fmt.Fprintf(&O, "%s\r\n^\r\n", x)
			textEdit.SetText(fmt.Sprintf("%s\r\n^\n", x))
			show(EDIT)
			restore()
		} else {
			backup()
		}
		if O.Len() > 0 {
			EO(string(bytes.Replace(O.Bytes(), []byte{'\n'}, []byte("\r\n"), -1)))
		}
		if len(O.Plots) > 0 {
			PO(O.Plots)
		}
	}()
	f()
}
func backup() {
	if B == nil || len(B) != len(C) {
		B = make([]byte, len(C))
	}
	copy(B, C)
}
func restore() {
	copy(C, B)
	C = C[:len(B)]
}

type fs struct {
	l []string
	m map[string][]byte
}

func (f *fs) add(path string) {
	b, e := ioutil.ReadFile(path)
	if e != nil {
		EO(e.Error())
	} else {
		if f.m == nil {
			f.m = make(map[string][]byte)
		}
		base := filepath.Base(path)
		if _, ok := f.m[base]; ok == false {
			f.l = append(f.l, base)
		}
		f.m[base] = b
		kdrop(base, b)
	}
}
func kdrop(s string, b []byte) {
	kdo(func() {
		drop := ktry("drop")
		if drop != 0xffffffff && drop != 0 {
			rx(drop)
			key := sc(mkcs([]byte(s)))
			val := mkcs(b)
			dx(cal(drop, l2(key, val)))
		}
	})
}

type output struct {
	bytes.Buffer
	plot.Plots
}

func (o *output) Reset() {
	o.Buffer.Reset()
	o.Plots = nil
}

const kpng = `iVBORw0KGgoAAAANSUhEUgAAABAAAAAQAgMAAABinRfyAAAACVBMVEX/AAAAAAD////KksOZAAAAMElEQVR4nGJYtWrVKoYFq1ZxMSyYhkZMgxNRXAwLpmbBCDAXSRZEgAwAGQUIAAD//+QzHr+8V1EyAAAAAElFTkSuQmCC`

var kimg image.Image

func init() {
	kimg, _ = png.Decode(base64.NewDecoder(base64.StdEncoding, strings.NewReader(kpng)))
}
