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

typedef struct {
	float xform[6];
	void *prev;
} draw_state;

typedef struct {
	NSVGparser *p;
	unsigned int color, strokeColor, fillColor;
	float lw;
	char fill, stroke, close, drawing;
	float px, py, cx, cy;
	draw_state *state;
} ctx_t;

typedef struct{
	const char *s;
	const char *a;
	void (*f)(ctx_t*, K*);
} draw_call;


K draw_func_names;
static NSVGattrib *startShape(NSVGparser *p, unsigned int color, float lw, char fill);
static void drawFlush(ctx_t *c, K *args);

static void float_args(float *f, K *x, size_t n){
	for(int i=0;i<n;i++){
		if(TK(x[i]) == 'f') f[i] = (float)fK(x[i]);
		else                f[i] = (float)iK(x[i]);
	}
}
static float float_arg(K *args){
	float f[1];
	float_args(f, args, 1);
	return f[0];
}
static void drawNewpath(ctx_t *c, float x, float y){
	drawFlush(c, NULL);
	nsvg__pushAttr(c->p);
	nsvg__resetPath(c->p);
	c->px=x; c->py=y;
	c->cx=x; c->cy=y;
	c->close=0; c->fill=0; c->stroke=0; c->drawing=1;
}
static void drawM(ctx_t *c, K *args){
	float f[2];
	float_args(f, args, 2);
	drawNewpath(c, f[0], f[1]);
	nsvg__moveTo(c->p, f[0], f[1]);
}
static void drawL(ctx_t *c, K *args){
	float f[2];
	float_args(f, args, 2);
	//nsvg__lineTo(c->p, f[0], f[1]);
	nsvg__pathLineTo(c->p, &c->px, &c->py, f, 0);
}
static void drawl(ctx_t *c, K *args){
	float f[2];
	float_args(f, args, 2);
	//nsvg__lineTo(c->p, f[0], f[1]);
	nsvg__pathLineTo(c->p, &c->px, &c->py, f, 1);
}
static void drawH(ctx_t *c, K *args){
	float f[2];
	float_args(f, args, 1);
	f[1] = 0.0;
	nsvg__pathLineTo(c->p, &c->px, &c->py, f, 0);
}
static void drawh(ctx_t *c, K *args){
	float f[2];
	float_args(f, args, 1);
	f[1] = 0.0;
	nsvg__pathLineTo(c->p, &c->px, &c->py, f, 1);
}
static void drawV(ctx_t *c, K *args){
	float f[2];
	float_args(f+1, args, 1);
	f[0] = 0.0;
	nsvg__pathLineTo(c->p, &c->px, &c->py, f, 0);
}
static void drawv(ctx_t *c, K *args){
	float f[2];
	float_args(f+1, args, 1);
	f[0] = 0.0;
	nsvg__pathLineTo(c->p, &c->px, &c->py, f, 1);
}
static void drawC(ctx_t *c, K *args){
	float f[6];
	float_args(f, args, 6);
	nsvg__pathCubicBezTo(c->p, &c->px, &c->py, &c->cx, &c->cy, f, 0);
}
static void drawc(ctx_t *c, K *args){
	float f[6];
	float_args(f, args, 6);
	nsvg__pathCubicBezTo(c->p, &c->px, &c->py, &c->cx, &c->cy, f, 1);
}
static void drawS(ctx_t *c, K *args){
	float f[4];
	float_args(f, args, 4);
	nsvg__pathCubicBezShortTo(c->p, &c->px, &c->py, &c->cx, &c->cy, f, 0);
}
static void draws(ctx_t *c, K *args){
	float f[4];
	float_args(f, args, 4);
	nsvg__pathCubicBezShortTo(c->p, &c->px, &c->py, &c->cx, &c->cy, f, 1);
}
static void drawQ(ctx_t *c, K *args){
	float f[4];
	float_args(f, args, 4);
	nsvg__pathQuadBezTo(c->p, &c->px, &c->py, &c->cx, &c->cy, f, 0);
}
static void drawq(ctx_t *c, K *args){
	float f[4];
	float_args(f, args, 4);
	nsvg__pathQuadBezTo(c->p, &c->px, &c->py, &c->cx, &c->cy, f, 1);
}
static void drawT(ctx_t *c, K *args){
	float f[2];
	float_args(f, args, 2);
	nsvg__pathQuadBezTo(c->p, &c->px, &c->py, &c->cx, &c->cy, f, 0);
}
static void drawt(ctx_t *c, K *args){
	float f[2];
	float_args(f, args, 2);
	nsvg__pathQuadBezTo(c->p, &c->px, &c->py, &c->cx, &c->cy, f, 1);
}
static void drawA(ctx_t *c, K *args){
	float f[7];
	float_args(f, args, 7);
	nsvg__pathArcTo(c->p, &c->px, &c->py, f, 0);
	c->cx=c->px; c->cy=c->py;
}
static void drawa(ctx_t *c, K *args){
	float f[7];
	float_args(f, args, 7);
	nsvg__pathArcTo(c->p, &c->px, &c->py, f, 1);
	c->cx=c->px; c->cy=c->py;
}

static void drawZ(ctx_t *c, K *args)     { c->close  = 1; }
static void drawFill(ctx_t *c, K *args)  { c->fill   = 1; c->fillColor=c->color;}
static void drawStroke(ctx_t *c, K *args){ c->stroke = 1; c->strokeColor=c->color; }
static void drawFlush(ctx_t *c, K *args) { 
	if(!c->drawing) return;
	NSVGattrib *a = nsvg__getAttr(c->p);
	a->hasFill = c->fill;
	a->hasStroke = c->stroke;
	memcpy(a->xform, c->state->xform, 6*sizeof(float));
	if(c->fill) {
		a->fillColor   = c->fillColor;
	}
	if(c->stroke){
		a->hasStroke = 1;
		a->strokeColor = c->strokeColor;
		a->strokeWidth = c->lw;
	}
	nsvg__addPath(c->p, c->close);
	nsvg__addShape(c->p);
	nsvg__popAttr(c->p);
	c->drawing = 0;
}
static void drawSave(ctx_t *c, K *args){
	draw_state *n = (draw_state*)malloc(sizeof(draw_state));
	memcpy(n->xform, c->state->xform, 6*sizeof(float));
	n->prev = c->state;
	c->state = n;
}
static void drawRest(ctx_t *c, K *args){
	draw_state *p = c->state->prev;
	if(p){
		free(c->state);
		c->state = p;
	}else	nsvg__xformIdentity(c->state->xform);
}
static void rectPath(ctx_t *c, float x, float y, float w, float h){
	drawNewpath(c, x, y);
	nsvg__moveTo(c->p, x, y);
	nsvg__lineTo(c->p, x+w, y);
	nsvg__lineTo(c->p, x+w, y+h);
	nsvg__lineTo(c->p, x, y+h);
	c->close = 1;
}
static void drawRect(ctx_t *c, K *args){
	float f[4];
	float_args(f, args, 4);
	rectPath(c, f[0], f[1], f[2], f[3]);	
}
static void drawCirc(ctx_t *c, K *args){
	float f[3];
	float_args(f, args, 3);
	float cx=f[0], cy=f[1], r=f[2];
	drawNewpath(c, cx+r, cy);
	nsvg__moveTo(c->p, cx+r, cy);
	nsvg__cubicBezTo(c->p, cx+r, cy+r*NSVG_KAPPA90, cx+r*NSVG_KAPPA90, cy+r, cx, cy+r);
	nsvg__cubicBezTo(c->p, cx-r*NSVG_KAPPA90, cy+r, cx-r, cy+r*NSVG_KAPPA90, cx-r, cy);
	nsvg__cubicBezTo(c->p, cx-r, cy-r*NSVG_KAPPA90, cx-r*NSVG_KAPPA90, cy-r, cx, cy-r);
	nsvg__cubicBezTo(c->p, cx+r*NSVG_KAPPA90, cy-r, cx+r, cy-r*NSVG_KAPPA90, cx+r, cy);
	c->close = 1;
}
static void drawElli(ctx_t *c, K *args){
	float f[4];
	float_args(f, args, 4);
	float cx=f[0], cy=f[1], rx=f[2], ry=f[3];
	drawNewpath(c, cx+rx, cy);
	nsvg__moveTo(c->p, cx+rx, cy);
	nsvg__cubicBezTo(c->p, cx+rx, cy+ry*NSVG_KAPPA90, cx+rx*NSVG_KAPPA90, cy+ry, cx, cy+ry);
	nsvg__cubicBezTo(c->p, cx-rx*NSVG_KAPPA90, cy+ry, cx-rx, cy+ry*NSVG_KAPPA90, cx-rx, cy);
	nsvg__cubicBezTo(c->p, cx-rx, cy-ry*NSVG_KAPPA90, cx-rx*NSVG_KAPPA90, cy-ry, cx, cy-ry);
	nsvg__cubicBezTo(c->p, cx+rx*NSVG_KAPPA90, cy-ry, cx+rx, cy-ry*NSVG_KAPPA90, cx+rx, cy);
	c->close = 1;
}
static void drawPoly(ctx_t *c, K *args){
	size_t n=NK(args[0]);
	double *x = (double*)malloc(8*n);
	FK(x, args[0]);
	if((NK(args[1]) != n)||(n < 2)){
		free(x);
		unref(args[1]);
		return;
	}
	double *y = (double*)malloc(8*n);
	FK(y, args[1]);
	drawNewpath(c, (float)x[0], (float)y[0]);
	nsvg__moveTo(c->p, (float)x[0], (float)y[0]);
	for(int i=1;i<n;i++){
		nsvg__lineTo(c->p, (float)x[i], (float)y[i]);
	}
	free(x);
	free(y);
}
static void drawCo(ctx_t *c, K *x){ c->color = (unsigned int)(0xff000000 | (uint32_t)iK(x[0])); }
static void drawLw(ctx_t *c, K *args){ c->lw = float_arg(args); }
static void drawTr(ctx_t *c, K *args){
	float f[2];
	float_args(f, args, 2);
	float t[6];
	nsvg__xformSetTranslation(t, f[0], f[1]);
	nsvg__xformPremultiply(c->state->xform, t);
}
static void drawRo(ctx_t *c, K *args){ 	
	float f = float_arg(args);
	float t[6];
	nsvg__xformSetRotation(t, f/180.0*NSVG_PI);
	nsvg__xformPremultiply(c->state->xform, t);
}
static void drawSc(ctx_t *c, K *args){ 	
	float f[2];
	float_args(f, args, 2);
	float t[6];
	nsvg__xformSetScale(t, f[0], f[1]);
	nsvg__xformPremultiply(c->state->xform, t);
}
	

const char *draw_api_symbols = "`M`L`l`H`h`V`v`C`c`S`s`Q`q`T`t`A`a`Z`rect`circ`elli`poly`fill`stroke`flush`save`rest`co`lw`tr`ro`sc";
draw_call   draw_api[] = {
	{"M",    "nn",    drawM},      // create new path and move to x y
	{"L",    "nn",    drawL},      // line to x y
	{"l",    "nn",    drawl},      // rel line to +x +y
	{"H",    "n",     drawH},      // hor line to x
	{"h",    "n",     drawh},      // rel hor line to +x
	{"V",    "n",     drawV},      // ver line to y
	{"v",    "n",     drawv},      // rel ver line to +y
	{"C",   "nnnnnn", drawC},      // cube to x1 y1 x2 y2 x y
	{"c",   "nnnnnn", drawC},      // rel cube to +x1 +y1 +x2 +y2 +x +y
	{"S",     "nnnn", drawS},      // short cube to x2 y2 x y
	{"s",     "nnnn", draws},      // rel short cube to +x2 +y2 +x +y
	{"Q",     "nnnn", drawQ},      // quad to x1 y1 x y
	{"q",     "nnnn", drawq},      // rel quad to +x1 +y1 +x +y
	{"T",     "nn",   drawT},      // short quad to x y
	{"t",     "nn",   drawt},      // rel short quad to +x +y
	{"A",  "nnnnnnn", drawA},      // arc to rx ry xrot large sweep x y
	{"a",  "nnnnnnn", drawa},      // rel arc to rx ry xrot large sweep +x +y
	{"Z",     "",     drawZ},      // close path
	{"rect",  "nnnn", drawRect},   // add rectangle x y w h
	{"circ",  "nnn",  drawCirc},   // add circle cx cy r
	{"elli",  "nnnn", drawElli},   // add ellipse cx cy rx ry
	{"poly",  "FF",   drawPoly},   // add polygon path
	{"fill",  "",     drawFill},   // fill path
	{"stroke","",     drawStroke}, // stroke path
	{"flush", "",     drawFlush},  // flush pending stroke and fill
	{"save",  "",     drawSave},   // save transformation
	{"rest",  "",     drawRest},   // restore transformation
	{"co",    "i",    drawCo},     // set color
	{"lw",    "n",    drawLw},     // set line width
	{"tr",    "nn",   drawTr},     // translate
	{"ro",    "n",    drawRo},     // rotate deg
	{"sc",    "nn",   drawSc},     // scale sx sy
};

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
//  x: (i height;I pixels)
//  r: C png bytes
K png(K x){
	if((TK(x) != 'L')||(NK(x) != 2)) { unref(x); return KE("png type x"); }
	K l[2];
	LK(l, x);
	  x = l[0];
	K y = l[1];
	if(TK(x) != 'i'){ unref(x); unref(y); return KE("png type *x"); }
	if(TK(y) != 'I'){ unref(x); unref(y); return KE("png type x 1"); }
	
	uint8_t *data;
	size_t height  = iK(x);
	size_t n      = NK(y);
	size_t width = n / height;
	if(n != width * height){ unref(y); return KE("png length x 1"); }
	
	//printf("png %d x %d\n", width, height);

	int *I = malloc(n*sizeof(int));
	IK(I, y);
	uint32_t *u = U32I(I, n);
	
	for(int i=0;i<n;i++) u[i] |= 0xff000000; // always opaque
	
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

int find_draw_func(K x){
	if(TK(x) != 's') return -1;
	K r = Kx("?", ref(draw_func_names), x);
	if(TK(r) != 'i') return -1;
	int i = iK(r);
	if((i<0)||(i>=NK(draw_func_names))) return -1;
	return i;
}
int draw_args_check(draw_call f, K *x, size_t n){
	size_t m = strlen(f.a);
	if(m > n) return 1;
	for(int i=0;i<m;i++){
		char t = TK(x[i]);
		if(f.a[i] == 'n'){ if((t != 'i') &&  (t != 'f')) return 1; }
		else if(t != f.a[i]) return 1;
	}
	return 0;
}

K draw(K x, K y){
	if(TK(x) != 'L'){ unref(x); unref(y); return KE("draw type x"); }
	if(TK(y) != 'I'){ unref(x); unref(y); return KE("draw type y"); }
	if(NK(y) !=  2 ){ unref(x); unref(y); return KE("draw length y"); }
	int wh[2];
	IK(wh, y);
	size_t n = NK(x);
	
	K *l = malloc(n*sizeof(K));
	LK(l, x);
	
	
	ctx_t ctx;
	ctx.p = nsvg__createParser();
	ctx.p->dpi = 96;
	ctx.p->image->width  = wh[0];
	ctx.p->image->height = wh[1];
	ctx.lw = 1.0;
	ctx.drawing = 0;
	
	draw_state state;
	nsvg__xformIdentity(state.xform);
	ctx.state = &state;
	ctx.state->prev = NULL;
	
	// fill white bg
	ctx.color = 0xffffffff;
	rectPath(&ctx, 0, 0, wh[0], wh[1]);
	drawFill(&ctx, NULL);
	ctx.color = 0xff000000;
	
	int i = 0;
	while(i < n){
		int f;
		if(((f=find_draw_func(l[i])) < 0) || draw_args_check(draw_api[f], l+i+1, n-i-1)){
			printf("f=%d\n", f);
			for(;i<n;i++) unref(l[i]);
			free(l);
			nsvg__deleteParser(ctx.p);
			return KE("draw call"); 
		}
		//printf("[%d] draw i=%d\n", f, i);
		draw_api[f].f(&ctx, l+1+i);
		i += 1 + strlen(draw_api[f].a);
	}
	drawFlush(&ctx, NULL);
	free(l);
	
	NSVGimage* im = ctx.p->image;
	ctx.p->image = NULL;
	nsvg__deleteParser(ctx.p);
	return svgRast(im);
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
	KR("png", (void*)png, 1);
	KR("ttf", (void*)ttf, 2);
	KR("svg", (void*)svg, 1);
	KR("setfont", (void*)setfont, 1);
	KR("draw", (void*)draw, 2);
	KR("rgb", (void*)rgb, 1);
	draw_func_names = Kx(draw_api_symbols);
}
