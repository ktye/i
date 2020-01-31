package main

import (
	"syscall"
	"unsafe"
)

func drw(x, y k) (r k) { // x 9:y  x(width `i) y(pixel-data `I)
	xt, yt, xn, yn := typs(x, y)
	if xt != I || yt != I || xn != atom || yn == atom {
		panic("type")
	}
	w := m.k[2+x]
	h := (m.k[y] & atom) / w
	if w > 8192 || h > 8192 || h < 1 {
		panic("size")
	}
	p := 8 + y<<2
	draw(w, h, m.c[p:p+4*w*h])
	return decr(x, y, 0)
}

type wBM struct { // bitmap
	a, b, c, d i
	e, f       uint16
	h          unsafe.Pointer
}
type wcxyz struct{ x, y, z i }   // ciexyz
type wc3 struct{ r, g, b wcxyz } // ciexyztriple
type wBI struct {                // bitmapinfoheader
	s    k
	w, h i
	p, c uint16
	m, i k
	x, y i
	u, t k
}
type wBM4 struct { // bitmapv4header
	wBI
	r, g, b, a, t k
	e             wc3
	rr, gg, bb    k
}
type wBM5 struct { // bitmapv5header
	wBM4
	a5, b5, c5, d5 k
}

func draw(w, h k, c []c) {
	hbm := toBM(w, h, c)
	d := winCreateCompatibleDC(0)
	xif(d == 0, "create compatible dc")
	defer winDeleteDC(d)
	o := winSelectObject(d, hbm)
	xif(o == 0, "select object")
	defer winSelectObject(d, o)
	xif(!winBitBlt(winGetDC(0), 0, 0, int32(w), int32(h), d, 0, 0, 0x00CC0020), "bitblt")
	winDeleteObject(hbm) // ?
}
func toBM(w, h k, c []c) (r uintptr) {
	var bi wBM5
	bi.s = k(unsafe.Sizeof(bi))
	bi.w, bi.h, bi.p, bi.c, bi.m, bi.x, bi.y = i(w), -i(h), 1, 32, 3, 3780, 3780
	bi.r, bi.g, bi.b, bi.a = 0x000000FF, 0x0000FF00, 0x00FF0000, 0xFF000000
	d := winGetDC(0)
	defer winReleaseDC(0, d)
	var lpBits unsafe.Pointer
	r = winCreateDIBSection(d, &bi.wBI, 0, &lpBits, 0, 0)
	switch r {
	case 0, 87:
		panic("CreateDIBSection failed")
	}
	a := (*[1 << 30]byte)(unsafe.Pointer(lpBits))
	copy(a[:], c)
	return r
}
func xif(c bool, e string) {
	if c {
		panic(e)
	}
}

var (
	libuser32          = syscall.NewLazyDLL("user32.dll")
	libgdi32           = syscall.NewLazyDLL("gdi32.dll")
	libkernel32        = syscall.NewLazyDLL("kernel32.dll")
	getDC              = libuser32.NewProc("GetDC")
	releaseDC          = libuser32.NewProc("ReleaseDC")
	deleteDC           = libgdi32.NewProc("DeleteDC")
	createCompatibleDC = libgdi32.NewProc("CreateCompatibleDC")
	createDIBSection   = libgdi32.NewProc("CreateDIBSection")
	selectObject       = libgdi32.NewProc("SelectObject")
	deleteObject       = libgdi32.NewProc("DeleteObject")
	bitBlt             = libgdi32.NewProc("BitBlt")
)

func winGetDC(h uintptr) (r uintptr) {
	r, _, _ = syscall.Syscall(getDC.Addr(), 1, h, 0, 0)
	return
}
func winReleaseDC(w uintptr, d uintptr) bool {
	r, _, _ := syscall.Syscall(releaseDC.Addr(), 2, w, d, 0)
	return r != 0
}
func winDeleteDC(h uintptr) bool {
	r, _, _ := syscall.Syscall(deleteDC.Addr(), 1, h, 0, 0)
	return r != 0
}
func winCreateCompatibleDC(h uintptr) (r uintptr) {
	r, _, _ = syscall.Syscall(createCompatibleDC.Addr(), 1, h, 0, 0)
	return
}
func winCreateDIBSection(h uintptr, p *wBI, u uint32, b *unsafe.Pointer, s uintptr, o uint32) (r uintptr) {
	r, _, _ = syscall.Syscall6(createDIBSection.Addr(), 6, h, uintptr(unsafe.Pointer(p)), uintptr(u), uintptr(unsafe.Pointer(b)), s, uintptr(o))
	return
}
func winSelectObject(h uintptr, o uintptr) (r uintptr) {
	r, _, _ = syscall.Syscall(selectObject.Addr(), 2, h, o, 0)
	return
}
func winDeleteObject(o uintptr) bool {
	r, _, _ := syscall.Syscall(deleteObject.Addr(), 1, o, 0, 0)
	return r != 0
}
func winBitBlt(d uintptr, x, y, w, h int32, s uintptr, xs, ys int32, o uint32) bool {
	r, _, _ := syscall.Syscall9(bitBlt.Addr(), 9, d, uintptr(x), uintptr(y), uintptr(w), uintptr(h), s, uintptr(xs), uintptr(ys), uintptr(o))
	return r != 0
}
