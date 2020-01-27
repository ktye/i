// draw device for win32 (no window, draw on DC 0)
// include<windows.h> tcc: -luser32 -lgdi32

HBITMAP toBM(I w, I h, char *p) {
	BITMAPV5HEADER   bi;
	HBITMAP          hbm;
	HDC              hdc;
	I *dst;
	ZeroMemory(&bi, sizeof(BITMAPV5HEADER));
	bi.bV5Size = sizeof(bi);
	bi.bV5Height = w; bi.bV5Width = h; bi.bV5Planes = 1; bi.bV5BitCount = 32; bi.bV5Compression = 3; 
	bi.bV5XPelsPerMeter = 3780; bi.bV5YPelsPerMeter = 3780;
	bi.bV5RedMask = 0x000000FF; bi.bV5GreenMask = 0x0000FF00; bi.bV5BlueMask = 0x00FF0000, bi.bV5AlphaMask = 0xFF000000;
	hdc = GetDC(0);
	hbm = CreateDIBSection(hdc, &bi, 0, &dst, 0, 0);
	ReleaseDC(0, hdc);
	if (hbm == NULL) {
        	// free(pbmi);
        	return NULL;
	}
	memcpy(dst, p, 4*w*h);
	return hbm;
}
void draw(I w, I h, I *p) {
	HDC     dc;
	HGDIOBJ  o;
	HBITMAP hbm = toBM(w, h, (char *)p);
	if (hbm == NULL) return;
	dc = CreateCompatibleDC(0);
	o = SelectObject(dc, hbm);
	if (o == 0) return;
	BitBlt(GetDC(0), 0, 0, w, h, dc, 0, 0, 0x00CC0020);
	SelectObject(dc, o);
	DeleteDC(dc);
	DeleteObject(hbm);
}
