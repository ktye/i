package main

import (
	"image"
	"image/color"
	"image/draw"
	"syscall"
	"unsafe"
)

/*
func main() {
	m := image.NewRGBA(image.Rectangle{Max: image.Point{800, 600}})
	draw.Draw(m, m.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)
	bm := NewBitmapFromImage(m)
	bm.draw()
}
*/

type winHANDLE uintptr
type winHDC winHANDLE
type winHWND winHANDLE
type winHGDIOBJ winHANDLE
type winHGLOBAL winHANDLE
type winHBITMAP winHGDIOBJ
type winBITMAP struct {
	BmType       int32
	BmWidth      int32
	BmHeight     int32
	BmWidthBytes int32
	BmPlanes     uint16
	BmBitsPixel  uint16
	BmBits       unsafe.Pointer
}
type winDIBSECTION struct {
	DsBm        winBITMAP
	DsBmih      winBITMAPINFOHEADER
	DsBitfields [3]uint32
	DshSection  winHANDLE
	DsOffset    uint32
}
type winCIEXYZ struct{ CiexyzX, CiexyzY, CiexyzZ int32 }
type winCIEXYZTRIPLE struct{ CiexyzRed, CiexyzGreen, CiexyzBlue winCIEXYZ }
type winBITMAPINFOHEADER struct {
	BiSize                           uint32
	BiWidth, BiHeight                int32
	BiPlanes, BiBitCount             uint16
	BiCompression, BiSizeImage       uint32
	BiXPelsPerMeter, BiYPelsPerMeter int32
	BiClrUsed, BiClrImportant        uint32
}
type winBITMAPV4HEADER struct {
	winBITMAPINFOHEADER
	BV4RedMask    uint32
	BV4GreenMask  uint32
	BV4BlueMask   uint32
	BV4AlphaMask  uint32
	BV4CSType     uint32
	BV4Endpoints  winCIEXYZTRIPLE
	BV4GammaRed   uint32
	BV4GammaGreen uint32
	BV4GammaBlue  uint32
}
type winBITMAPV5HEADER struct {
	winBITMAPV4HEADER
	BV5Intent, BV5ProfileData, BV5ProfileSize, BV5Reserved uint32
}

var (
	libuser32          = syscall.NewLazyDLL("user32.dll")
	libgdi32           = syscall.NewLazyDLL("gdi32.dll")
	libkernel32        = syscall.NewLazyDLL("kernel32.dll")
	getDC              = libuser32.NewProc("GetDC")
	releaseDC          = libuser32.NewProc("ReleaseDC")
	createCompatibleDC = libgdi32.NewProc("CreateCompatibleDC")
	createDIBSection   = libgdi32.NewProc("CreateDIBSection")
	deleteDC           = libgdi32.NewProc("DeleteDC")
	selectObject       = libgdi32.NewProc("SelectObject")
	bitBlt             = libgdi32.NewProc("BitBlt")
	deleteObject       = libgdi32.NewProc("DeleteObject")
	globalAlloc        = libkernel32.NewProc("GlobalAlloc")
	globalFree         = libkernel32.NewProc("GlobalFree")
	globalLock         = libkernel32.NewProc("GlobalLock")
	globalUnlock       = libkernel32.NewProc("GlobalUnlock")
	getObject          = libgdi32.NewProc("GetObjectW")
	moveMemory         = libkernel32.NewProc("RtlMoveMemory")
)

func winGetDC(hWnd winHWND) winHDC {
	ret, _, _ := syscall.Syscall(getDC.Addr(), 1, uintptr(hWnd), 0, 0)
	return winHDC(ret)
}
func winReleaseDC(hWnd winHWND, hDC winHDC) bool {
	ret, _, _ := syscall.Syscall(releaseDC.Addr(), 2, uintptr(hWnd), uintptr(hDC), 0)
	return ret != 0
}
func winCreateCompatibleDC(hdc winHDC) winHDC {
	ret, _, _ := syscall.Syscall(createCompatibleDC.Addr(), 1, uintptr(hdc), 0, 0)
	return winHDC(ret)
}
func winCreateDIBSection(hdc winHDC, pbmih *winBITMAPINFOHEADER, iUsage uint32, ppvBits *unsafe.Pointer, hSection winHANDLE, dwOffset uint32) winHBITMAP {
	ret, _, _ := syscall.Syscall6(createDIBSection.Addr(), 6, uintptr(hdc), uintptr(unsafe.Pointer(pbmih)), uintptr(iUsage), uintptr(unsafe.Pointer(ppvBits)), uintptr(hSection), uintptr(dwOffset))
	return winHBITMAP(ret)
}

func winDeleteDC(hdc winHDC) bool {
	ret, _, _ := syscall.Syscall(deleteDC.Addr(), 1, uintptr(hdc), 0, 0)
	return ret != 0
}
func winSelectObject(hdc winHDC, hgdiobj winHGDIOBJ) winHGDIOBJ {
	ret, _, _ := syscall.Syscall(selectObject.Addr(), 2, uintptr(hdc), uintptr(hgdiobj), 0)
	return winHGDIOBJ(ret)
}
func winBitBlt(hdcDest winHDC, nXDest, nYDest, nWidth, nHeight int32, hdcSrc winHDC, nXSrc, nYSrc int32, dwRop uint32) bool {
	ret, _, _ := syscall.Syscall9(bitBlt.Addr(), 9, uintptr(hdcDest), uintptr(nXDest), uintptr(nYDest), uintptr(nWidth), uintptr(nHeight), uintptr(hdcSrc), uintptr(nXSrc), uintptr(nYSrc), uintptr(dwRop))
	return ret != 0
}
func winDeleteObject(hObject winHGDIOBJ) bool {
	ret, _, _ := syscall.Syscall(deleteObject.Addr(), 1, uintptr(hObject), 0, 0)
	return ret != 0
}
func winGlobalAlloc(uFlags uint32, dwBytes uintptr) winHGLOBAL {
	ret, _, _ := syscall.Syscall(globalAlloc.Addr(), 2, uintptr(uFlags), dwBytes, 0)
	return winHGLOBAL(ret)
}
func winGlobalFree(hMem winHGLOBAL) winHGLOBAL {
	ret, _, _ := syscall.Syscall(globalFree.Addr(), 1, uintptr(hMem), 0, 0)
	return winHGLOBAL(ret)
}
func winGlobalLock(hMem winHGLOBAL) unsafe.Pointer {
	ret, _, _ := syscall.Syscall(globalLock.Addr(), 1, uintptr(hMem), 0, 0)
	return unsafe.Pointer(ret)
}
func winGlobalUnlock(hMem winHGLOBAL) bool {
	ret, _, _ := syscall.Syscall(globalUnlock.Addr(), 1, uintptr(hMem), 0, 0)
	return ret != 0
}
func winGetObject(hgdiobj winHGDIOBJ, cbBuffer uintptr, lpvObject unsafe.Pointer) int32 {
	ret, _, _ := syscall.Syscall(getObject.Addr(), 3, uintptr(hgdiobj), uintptr(cbBuffer), uintptr(lpvObject))
	return int32(ret)
}
func winMoveMemory(destination, source unsafe.Pointer, length uintptr) {
	syscall.Syscall(moveMemory.Addr(), 3, uintptr(unsafe.Pointer(destination)), uintptr(source), uintptr(length))
}

type Bitmap struct {
	hBmp       winHBITMAP
	hPackedDIB winHGLOBAL
	w, h       int32
}

func NewBitmapFromImage(im image.Image) *Bitmap {
	return newBitmapFromHBITMAP(hBitmapFromImage(im))
}
func (bmp *Bitmap) draw() {
	hdc := winCreateCompatibleDC(0)
	xif(hdc == 0, "create compatible dc")
	defer winDeleteDC(hdc)
	hBmpOld := winSelectObject(hdc, winHGDIOBJ(bmp.hBmp))
	xif(hBmpOld == 0, "select object")
	defer winSelectObject(hdc, hBmpOld)
	xif(!winBitBlt(winGetDC(0), 0, 0, bmp.w, bmp.h, hdc, 0, 0, 0x00CC0020), "bitblt")
	winDeleteObject(winHGDIOBJ(bmp.hBmp))
	winGlobalUnlock(bmp.hPackedDIB)
	winGlobalFree(bmp.hPackedDIB)
	bmp.hPackedDIB = 0
	bmp.hBmp = 0
}
func hBitmapFromImage(im image.Image) winHBITMAP {
	var bi winBITMAPV5HEADER
	bi.BiSize = uint32(unsafe.Sizeof(bi))
	bi.BiWidth = int32(im.Bounds().Dx())
	bi.BiHeight = -int32(im.Bounds().Dy())
	bi.BiPlanes = 1
	bi.BiBitCount = 32
	bi.BiCompression = 3
	bi.BiXPelsPerMeter = 3780 // 96 dpi * 39.37 inches per meter
	bi.BiYPelsPerMeter = 3780
	bi.BV4RedMask = 0x00FF0000
	bi.BV4GreenMask = 0x0000FF00
	bi.BV4BlueMask = 0x000000FF
	bi.BV4AlphaMask = 0xFF000000
	hdc := winGetDC(0)
	defer winReleaseDC(0, hdc)
	var lpBits unsafe.Pointer
	hBitmap := winCreateDIBSection(hdc, &bi.winBITMAPINFOHEADER, 0, &lpBits, 0, 0)
	switch hBitmap {
	case 0, 87:
		panic("CreateDIBSection failed")
	}
	bitmap_array := (*[1 << 30]byte)(unsafe.Pointer(lpBits))
	i := 0
	for y := im.Bounds().Min.Y; y != im.Bounds().Max.Y; y++ {
		for x := im.Bounds().Min.X; x != im.Bounds().Max.X; x++ {
			r, g, b, a := im.At(x, y).RGBA()
			bitmap_array[i+3] = byte(a >> 8)
			bitmap_array[i+2] = byte(r >> 8)
			bitmap_array[i+1] = byte(g >> 8)
			bitmap_array[i+0] = byte(b >> 8)
			i += 4
		}
	}
	return hBitmap
}
func newBitmapFromHBITMAP(hBmp winHBITMAP) (bmp *Bitmap) {
	var dib winDIBSECTION
	xif(winGetObject(winHGDIOBJ(hBmp), unsafe.Sizeof(dib), unsafe.Pointer(&dib)) == 0, "getobject")
	bmih := &dib.DsBmih
	bmihSize := uintptr(unsafe.Sizeof(*bmih))
	pixelsSize := uintptr(int32(bmih.BiBitCount)*bmih.BiWidth*bmih.BiHeight) / 8
	totalSize := uintptr(bmihSize + pixelsSize)
	hPackedDIB := winGlobalAlloc(0x0042, totalSize)
	dest := winGlobalLock(hPackedDIB)
	defer winGlobalUnlock(hPackedDIB)
	src := unsafe.Pointer(&dib.DsBmih)
	winMoveMemory(dest, src, bmihSize)
	dest = unsafe.Pointer(uintptr(dest) + bmihSize)
	src = dib.DsBm.BmBits
	winMoveMemory(dest, src, pixelsSize)
	return &Bitmap{hBmp: hBmp, hPackedDIB: hPackedDIB, w: bmih.BiWidth, h: bmih.BiHeight}
}
func xif(c bool, e string) {
	if c {
		panic(e)
	}
}
