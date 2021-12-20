#include<stdlib.h>
#include"../k.h"

#include<stdio.h>

// http://nothings.org/stb/stb_image_write.h
#define STBI_WRITE_NO_STDIO
#define STBI_ONLY_PNG
#define STB_IMAGE_WRITE_IMPLEMENTATION
#include "stb_image_write.h" 

// http://nothings.org/stb/stb_truetype.h
#define STB_TRUETYPE_IMPLEMENTATION 
#include "stb_truetype.h" 

// https://github.com/memononen/nanosvg/tree/master/src
#define NANOSVG_IMPLEMENTATION
#define NANOSVGRAST_IMPLEMENTATION
#include "nanosvg.h"
#include "nanosvgrast.h"


void wpng(void *context, void *data, int size){
	K *x = (K*)context;
	K r = KC((char*)data, (size_t)size);
	*x = r;
}

uint32_t *U32I(int *m, size_t n){
	if(4 == sizeof(int)) return (uint32_t*)m;
	uint32_t *u = (uint32_t *)malloc(4*n);
	for(int i=0;i<n;i++) u[i] = (uint32_t)m[i];
	free(m);
	return u;
}

// convert pixels to png data.
//  x: I pixel data, height: (#x)%y
//  y: i width
//  r: C png bytes
K png(K x, K y){
	uint8_t *data;
	size_t   width, height;
	
	if(TK(x) != 'I') { unref(x); unref(y); return KE("png type x"); }
	if(TK(y) != 'i') { unref(x); unref(y); return KE("png type y"); }

	width = iK(y);
	size_t n = NK(x);
	height = n / width;
	if(n != width * height){ unref(x); return KE("x length"); }

	int *I = malloc(n*sizeof(int));
	IK(I, x);
	uint32_t *u = U32I(I, n);
	
	for(int i=0;i<n;i++) u[i] |= 0xff000000; // always opaque
		
	printf("write png %d x %d\n", width, height); //rm
	
	K r;
	stbi_write_png_to_func(wpng, &r, (int)width, (int)height, 4, u, 4*width);
	free(u);
	return r;
}

K rgb(K x){
	if((TK(x) != 'I')||(NK(x) != 3)){ unref(x); return KE("rgb type x"); }
	int i[3];
	IK(i, x);
	uint32_t c = 0xff000000 | (i[0]&0xff) | ((i[1]&0xff)<<8) | ((i[2]&0xff)<<16);
	return Ki((int)c);
}

K Scolor, SlineWidth, SdrawRect, SfillRect, SdrawCirc, SfillCirc, SdrawElli, SfillElli, SdrawPoly, SfillPoly;
K Sbegin, SlineTo, SquadTo, ScubeTo, SarcTo, Sclose, Sstroke, Sfill;

static int eqs(K x, K y){ return x == y; } //todo: this works only if the k backend can compare symbols that way.
static int drawApi(K *e, size_t m, K s, const char *a){
	if(m != 1+strlen(a))                          return 0;
	if(!eqs(e[0], s))                             return 0;
	for(int i=0;i<m-1;i++) {
		char t = TK(e[1+i]);
		if(!((t == a[i])||((a[i]=='n') && ((t=='i') || (t == 'f')))))
			return 0;
	}
	unref(e[0]);
	return 1;
}
static float nK(K x){
	if(TK(x) == 'i') return (float)iK(x);
	return (float)fK(x);
}

static NSVGattrib *startShape(NSVGparser *p, unsigned int color, float lw, char fill){
	nsvg__pushAttr(p);
	NSVGattrib *a = nsvg__getAttr(p);
	nsvg__resetPath(p);
	a->hasFill = fill;
	a->hasStroke = 0;
	if(fill) {
		a->fillColor   = color;
	} else {
		a->hasStroke = 1;
		a->strokeColor = color;
		a->strokeWidth = lw;
	}
	return a;
}
static void strokeFill(NSVGattrib *a, NSVGparser *p, unsigned int color, float lw, char fill, char close){
	a->hasFill = fill;
	a->hasStroke = 1-fill;
	if(!fill){
		a->strokeColor = color;
		a->strokeWidth = lw;
	}
	nsvg__addPath(p, close);
	nsvg__addShape(p);
	nsvg__popAttr(p);
}
static void drawRect(NSVGparser *p, float x, float y, float w, float h, unsigned int color, float lw, char fill){	
	NSVGattrib *a = startShape(p, color, lw, fill);
	nsvg__moveTo(p, x, y);
	nsvg__lineTo(p, x+w, y);
	nsvg__lineTo(p, x+w, y+h);
	nsvg__lineTo(p, x, y+h);
	nsvg__addPath(p, 1);
	nsvg__addShape(p);
	nsvg__popAttr(p);
}
static void drawCirc(NSVGparser *p, float cx, float cy, float r, unsigned int color, float lw, char fill){
	NSVGattrib *a = startShape(p, color, lw, fill);
	nsvg__moveTo(p, cx+r, cy);
	nsvg__cubicBezTo(p, cx+r, cy+r*NSVG_KAPPA90, cx+r*NSVG_KAPPA90, cy+r, cx, cy+r);
	nsvg__cubicBezTo(p, cx-r*NSVG_KAPPA90, cy+r, cx-r, cy+r*NSVG_KAPPA90, cx-r, cy);
	nsvg__cubicBezTo(p, cx-r, cy-r*NSVG_KAPPA90, cx-r*NSVG_KAPPA90, cy-r, cx, cy-r);
	nsvg__cubicBezTo(p, cx+r*NSVG_KAPPA90, cy-r, cx+r, cy-r*NSVG_KAPPA90, cx+r, cy);
	nsvg__addPath(p, 1);
	nsvg__addShape(p);
	nsvg__popAttr(p);
}
static void drawElli(NSVGparser *p, float cx, float cy, float rx, float ry, unsigned int color, float lw, char fill){
	NSVGattrib *a = startShape(p, color, lw, fill);
	nsvg__moveTo(p, cx+rx, cy);
	nsvg__cubicBezTo(p, cx+rx, cy+ry*NSVG_KAPPA90, cx+rx*NSVG_KAPPA90, cy+ry, cx, cy+ry);
	nsvg__cubicBezTo(p, cx-rx*NSVG_KAPPA90, cy+ry, cx-rx, cy+ry*NSVG_KAPPA90, cx-rx, cy);
	nsvg__cubicBezTo(p, cx-rx, cy-ry*NSVG_KAPPA90, cx-rx*NSVG_KAPPA90, cy-ry, cx, cy-ry);
	nsvg__cubicBezTo(p, cx+rx*NSVG_KAPPA90, cy-ry, cx+rx, cy-ry*NSVG_KAPPA90, cx+rx, cy);
	nsvg__addPath(p, 1);
	nsvg__addShape(p);
	nsvg__popAttr(p);
}
static void drawPoly(NSVGparser *p, K x, K y, unsigned int color, float lw, char fill){
	size_t xn = NK(x);
	size_t yn = NK(y);
	if(xn != yn || xn < 2) { unref(x); unref(y); return; }
	double *xp = (double*)malloc(8*xn);
	double *yp = (double*)malloc(8*xn);
	FK(xp, x);
	FK(yp, y);
	
	NSVGattrib *a = startShape(p, color, lw, fill);
	nsvg__moveTo(p, (float)xp[0], (float)yp[0]);
	for(int i=1;i<xn;i++) nsvg__lineTo(p, (float)xp[i], (float)yp[i]);
	nsvg__addPath(p, fill);
	nsvg__addShape(p);
	nsvg__popAttr(p);
	free(xp);
	free(yp);
}

static K svgRast(NSVGimage *im) { //(h;I)
	int w = im->width;
	int h = im->height;
	uint8_t *dst = malloc(w*h*4);
	NSVGrasterizer *rst = nsvgCreateRasterizer();
	nsvgRasterize(rst, im, 0,0,1, dst, w, h, w*4);
	
	K ri;
	if(4 == sizeof(int)){
		ri = KI((int*)dst, w*h);
	}else{
		int *I = (int*)malloc(w*h*sizeof(int));
		uint32_t *u = (uint32_t *)dst;
		for(int i=0;i<w*h;i++) I[i] = (int)u[i];
		ri = KI(I, w*h);
		free(I);
	}
	nsvgDeleteRasterizer(rst);
	nsvgDelete(im);
	free(dst);
	K l[2] = {Ki(h), ri};
	return KL(l, 2);
}

static K drawErr(K *l, int i, size_t n, NSVGparser *p){
	for(;i<n;i++) unref(l[i]);
	free(l);
	p->image = NULL;
	nsvg__deleteParser(p);
	return KE("draw data x");
}

// draw rasterizes 2d vector graphics
//  x: L api-calls
//  y: I (w;h) size in pixels
//  r: L (h;I) pixels as ints
K draw(K x, K y){
	if(TK(x) != 'L'){ unref(x); unref(y); return KE("draw type x"); }
	if(TK(y) != 'I'){ unref(x); unref(y); return KE("draw type y"); }
	if(NK(y) !=  2 ){ unref(x); unref(y); return KE("draw length y"); }
	int wh[2];
	IK(wh, y);
	size_t n = NK(x);
	printf("draw w=%d h=%d n=%d\n", wh[0], wh[1], n);fflush(stdout); //rm
	
	
	K *l = malloc(n*sizeof(K));
	LK(l, x);
	
	
	NSVGparser* p;
	p = nsvg__createParser();
	p->dpi = 96;
	p->image->width  = wh[0];
	p->image->height = wh[1];
	drawRect(p, 0, 0, (float)wh[0], (float)wh[1], NSVG_RGB(255,255,255), 0, 1); //white-bg
	
	//nsvg__parseXML(input, nsvg__startElement, nsvg__endElement, nsvg__content, p);
	
	NSVGattrib *attr = NULL;
	char close = 0;
	
	unsigned int color = 0xff000000;
	float lw=1, px, py, cx, cy;
	int a, b, c, d, ie;
	K e[8];
	for(int i=0;i<n;i++){
		if(TK(l[i] != 's') && ((TK(l[i]) != 'L')||(NK(l[i]) >= 8))){ return drawErr(l, i, n, p); }
		size_t m = 1;
		if(TK(l[i]) != 's') m = NK(l[i]);              
		if(m >= 8){ return drawErr(l, i, n, p); }
		if(m == 1) e[0] = l[i]; 
		else    LK(e,     l[i]);
		printf("draw element i=%d n=%d m=%d\n", i, n, m);
		if     (drawApi(e, m, Scolor, "i"))        color = 0xff000000|(unsigned int)Ki(e[1]);
		else if(drawApi(e, m, SlineWidth, "n"))    lw = nK(e[1]);
		else if(drawApi(e, m, SdrawRect, "nnnn"))  drawRect(p, nK(e[1]), nK(e[2]), nK(e[3]), nK(e[4]), color, lw, 0);
		else if(drawApi(e, m, SfillRect, "nnnn"))  drawRect(p, nK(e[1]), nK(e[2]), nK(e[3]), nK(e[4]), color, lw, 1);
		else if(drawApi(e, m, SdrawCirc, "nnn"))   drawCirc(p, nK(e[1]), nK(e[2]), nK(e[3]),           color, lw, 0);
		else if(drawApi(e, m, SfillCirc, "nnn"))   drawCirc(p, nK(e[1]), nK(e[2]), nK(e[3]),           color, lw, 1);
		else if(drawApi(e, m, SdrawElli, "nnnn"))  drawElli(p, nK(e[1]), nK(e[2]), nK(e[3]), nK(e[4]), color, lw, 0);
		else if(drawApi(e, m, SfillElli, "nnnn"))  drawElli(p, nK(e[1]), nK(e[2]), nK(e[3]), nK(e[4]), color, lw, 1);
		else if(drawApi(e, m, SdrawPoly, "FF"))    drawPoly(p, e[1], e[2],                             color, lw, 0);
		else if(drawApi(e, m, SfillPoly, "FF"))    drawPoly(p, e[1], e[2],                             color, lw, 1);
		else if(drawApi(e, m, Sbegin, "nn")){
			attr = startShape(p, color, lw, 0);
			close = 0;
			px = nK(e[1]); py=nK(e[2]); cx=px; cy=py;
			nsvg__moveTo(p, px, py);
		}
		else if(drawApi(e, m, SlineTo, "nn")) {
			nsvg__pathLineTo(p, &px, &py, (float[]){nK(e[1]), nK(e[2])}, 0);
			cx=px; cy=py;
		}
		else if(drawApi(e, m, SquadTo, "nnnn")) nsvg__pathQuadBezTo(p, &px, &py, &cx, &cy, (float[]){nK(e[1]), nK(e[2]), nK(e[3]), nK(e[4])}, 0);
		else if(drawApi(e, m, ScubeTo, "nnnnnn")) nsvg__pathCubicBezTo(p, &px, &py, &cx, &cy, (float[]){nK(e[1]), nK(e[2]), nK(e[3]), nK(e[4]), nK(e[5]), nK(e[6])}, 0);
		else if(drawApi(e, m, SarcTo, "nnnnnn")){
			nsvg__pathArcTo(p, &px, &py, (float[]){nK(e[1]), nK(e[2]), nK(e[3]), nK(e[4]), nK(e[5]), nK(e[6])}, 0);
			cx=px; cy=py;
		}
		else if(drawApi(e, m, Sclose, ""))   close = 1;
		else if(drawApi(e, m, Sstroke,""))   strokeFill(attr, p, color, lw, 0, close);
		else if(drawApi(e, m, Sfill, ""))    strokeFill(attr, p, color, lw, 1, close);
		
		else{
			printf("draw unknown case i=%d m=%d\n", i, m);fflush(stdout);
			for(int j=0; j<m; j++) unref(e[i]);
			return drawErr(l, 1+i, n, p);
		}
	}
	
	//nsvg__scaleToViewbox(p, "px");
	NSVGimage* im = p->image;
	p->image = NULL;
	nsvg__deleteParser(p);
	return svgRast(im);
}

// rasterize svg image
//  x: C svg string
//  r: L (i height;I pixels)
K svg(K x){
	if(TK(x) != 'C'){ unref(x); return KE("svg type x"); }
	
	size_t n = NK(x);
	char *s = malloc(1+n);
	CK(s, x);
	s[n] = (char)0;
	
	NSVGimage *im = nsvgParse(s, "px", 96);
	free(s);
	if(im == NULL){ return KE("svg parse"); }
	
	return svgRast(im);
}


#define NTTF 8 //max number of fonts
int nttf = 0, currentfont = -1;
char          *ttfdata[NTTF];
stbtt_fontinfo ttfinfo[NTTF];
float          ttfscal[NTTF];


// load ttf font from data.
//  x: C ttf bytes
//  y: i size
//  r: i font index
K ttf(K x, K y){
	if(TK(x) != 'C'){   unref(x); unref(y); return KE("ttf: x type"); }
	if(TK(y) != 'i'){   unref(x); unref(y); return KE("ttf: y type"); }
	if(nttf == NTTF-1){ unref(x); unref(y); return KE("ttf: max fonts"); }
	
	ttfdata[nttf] = malloc(NK(x));
	CK(ttfdata[nttf], x);
	int height = iK(y);
	if(!stbtt_InitFont(&ttfinfo[nttf], ttfdata[nttf], 0)) return KE("ttf load");
	ttfscal[nttf] = stbtt_ScaleForPixelHeight(&ttfinfo[nttf], height);
	nttf++;
	currentfont = nttf-1;
	return Ki(currentfont);
}

// change current font index
K setfont(K x){
	if(TK(x) != 'i'){ unref(x); return KE("setfont type"); }
	int i = iK(x);
	if((i<0)||(i>=nttf)){ return KE("setfont value"); }
	currentfont = i;
	return Ki(i);
}

/*


// draw text on bitmap
//  x: I (x;y) position
//  y: C text
//  r: i x+textwidth
K drawtext(K x, K y, K z){
 if(TK(x) != 'I'){ unref(x); unref(y);return KE("drawtext x type"); }
 if(TK(y) != 'C'){ unref(x); unref(y);return KE("drawtext y type"); }
 if(NK(x) !=  2) { unref(x); unref(y);return KE("drawtext x length"); }
 
 // example:
 // https://raw.githubusercontent.com/justinmeiners/stb-truetype-example/master/main.c
 
 int xy[2];
 IK(xy, x);
 
 char *text = dK(y);
 size_t n = NK(y);
 unref(y);
 
 stbtt_fontinfo *info = &ttfinfo[currentfont];
 float          scale =  ttfscal[currentfont];
 
 int asc, dsc, gap;
 stbtt_GetFontVMetrics(info, &asc, &dsc, &gap);
 asc = roundf(asc * scale);
 dsc = roundf(dsc * scale);
 
 // todo unicode..
 for(int i=0;i<n;i++){
 	int c = (int)text[i];
	int ax, lsb;        //advanceWidth, leftSideBearing
	int x0, y0, x1, y1; //bitmap bounding box
	stbtt_GetCodepointHMetrics( info, c, &ax, &lsb);
	stbtt_GetCodepointBitmapBox(info, c, scale, scale, &x0, &y0, &x1, &y1);
	
	int yc = asc + y1;
	//todo...
 }
 
 return Ki(0);
}

*/

void loadimg(){
	KR("png", (void*)png, 2);
	KR("ttf", (void*)ttf, 2);
	KR("svg", (void*)svg, 1);
	KR("setfont", (void*)setfont, 1);
	KR("draw", (void*)draw, 2);
	KR("rgb", (void*)rgb, 1);
	
	Scolor     = Ks("color");
	SlineWidth = Ks("lineWidth");
	SdrawRect  = Ks("drawRect");
	SfillRect  = Ks("fillRect");
	SdrawCirc  = Ks("drawCirc");
	SfillCirc  = Ks("fillCirc");
	SdrawElli  = Ks("drawElli");
	SfillElli  = Ks("fillElli");
	SdrawPoly  = Ks("drawPoly");
	SfillPoly  = Ks("fillPoly");
	Sbegin     = Ks("begin");
	SlineTo    = Ks("lineTo");
	SquadTo    = Ks("quadTo");
	ScubeTo    = Ks("cubeTo");
	SarcTo     = Ks("arcTo");
	Sstroke    = Ks("stroke");
	Sfill      = Ks("fill");
}
