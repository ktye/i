#include<stdlib.h>
#include"../k.h"

#include<stdio.h>

// http://nothings.org/stb/stb_image_write.h
#define STBI_WRITE_NO_STDIO
#define STBI_ONLY_PNG
//#define STB_IMAGE_WRITE_IMPLEMENTATION (included in raylib: ../ray)
#include "stb_image_write.h" 

// http://nothings.org/stb/stb_truetype.h
#define STB_TRUETYPE_IMPLEMENTATION 
#include "stb_truetype.h" 

// https://github.com/memononen/nanosvg/tree/master/src
#define NANOSVG_IMPLEMENTATION
#define NANOSVGRAST_IMPLEMENTATION
#include "nanosvg.h"
#include "nanosvgrast.h"

static void wpng(void *context, void *data, int size);
static uint32_t *U32I(int *m, size_t n);
static void imgk(K x, size_t *w, size_t *h, uint32_t **data);
static int vec(float *v, size_t n, K x);
static int vecn(K x);
static float *veca(K x, int n);
static K drawerr(K *l, int i, size_t n, NSVGparser *p, char *s);
static void drawOver(uint32_t *dst, uint32_t *src, size_t n);
static void drawChar(uint32_t *dst, int x, int y, size_t w, size_t h, unsigned int co, uint8_t *bm, int cw, int ch, int x0, int y0);
static void drawRect(NSVGparser *p, float x, float y, float w, float h, unsigned int co, int lw, char dofill);
static void drawCircle(NSVGparser *p, float cx, float cy, float r, unsigned int co, int lw, char dofill);	
static void drawPoly(NSVGparser *p, float *px, float *py, int n, unsigned int co, int lw, char dofill);
static void newpath(NSVGparser *p);
static void fill(NSVGparser *p, unsigned int color);
static void stroke(NSVGparser *p, unsigned int color, float lw, char close);


#define MAXFONTS 8
K fontnames;
stbtt_fontinfo ttfinfo[MAXFONTS];
stbtt_fontinfo *currentfont = NULL;
float fontscale;
int   fontascent;

K loadfont(K name, K ttfdata, float *scale){ //data persists
 static int nttf = 0;
 if((TK(name) != 's')||(TK(ttfdata) != 'C')){ unref(name); unref(ttfdata); return KE("loadfont args"); }
 if(nttf == MAXFONTS){ unref(ttfdata); return name; }
 fontnames = Kx(",", fontnames, ref(name));
 char *buf = malloc(NK(ttfdata));
 CK(buf, ttfdata);
 if(!stbtt_InitFont(&ttfinfo[nttf], buf, 0)){ unref(name); return KE("loadfont: load ttf"); }
 nttf++;
 return name;
}
static void setfont(K x){ // "20px monospace"
 size_t n = NK(x);
 char *p = dK(x);
 unref(x);
 int h;
 if((n<4)||(n>99)){ printf("setfont ignored (n)\n");  return;}
 if(p[1] == 'p'){       h = (int)(p[0]-'0');                 p+=4; n-=4; }
 else if(p[2] == 'p'){  h = (int)(10*(p[0]-'0')+(p[1]-'0')); p+=5; n-=5; }
 else {   printf("setfont ignored (px)\n"); return; }
 
 char   b[100];
 memcpy(b, p, n); b[n] = (char)0;
 K name = Ks(b);
 
 int i = iK(Kx("?", ref(fontnames), name));
 if((i<0)||(i>=NK(fontnames))) { printf("setfont ignored (find)\n"); return; }
 currentfont = &ttfinfo[i];
 
 // see: https://github.com/justinmeiners/stb-truetype-example/blob/master/main.c
 fontscale = stbtt_ScaleForPixelHeight(&ttfinfo[i], h);
 
 int descent, lineGap;
 stbtt_GetFontVMetrics(currentfont, &fontascent, &descent, &lineGap);
 fontascent = roundf(fontascent * fontscale);
 descent = roundf(descent * fontscale);
 
 //printf("setfont %s size %d found at i=%d scale=%f ascent=%d descent=%d linegap=%d\n", b, h, i, fontscale, fontascent, descent, lineGap);
}
void drawText(uint32_t *dst, size_t w, size_t h, unsigned int co, int x, int y, char *cc, size_t nc){
 if(currentfont==NULL){ printf("no current font\n"); return; }
 for(int i=0;i<nc;i++){ //todo: decode utf8 (currently ascii-only)
  int ax, lsb, cw, ch, x0, y0, kern;
  stbtt_GetCodepointHMetrics(currentfont, cc[i], &ax, &lsb);
  unsigned char *bm = stbtt_GetCodepointBitmap(currentfont, fontscale, fontscale, (int)cc[i], &cw, &ch, &x0, &y0);  
  drawChar(dst, x, y, w, h, co, (uint8_t*)bm, cw, ch, x0, y0);
  stbtt_FreeBitmap(bm, NULL);
  kern = (i<nc-1) ? stbtt_GetCodepointKernAdvance(currentfont, (int)cc[i], (int)cc[1+i]) : 0;
  x += (int)roundf(fontscale*(ax+kern));
 }
}


// convert pixels to png data.
//  x: (i height;I pixels)
//  r: C png bytes
K png(K x){
 uint32_t *u;
 size_t    width, height;
 imgk(x, &width, &height, &u);
 //printf("png: %d x %d\n", width, height);
 if(u == NULL) return KE("png: type img");

 for(int i=0;i<width*height;i++) u[i] |= 0xff000000; // always opaque
 
 K r;
 stbi_write_png_to_func(wpng, &r, (int)width, (int)height, 4, u, 4*width);
 free(u);
 return r;
}


K drawcmds; // "`color`font`linewidth`rect`Rect`circle`Circle`line`poly`Poly`text`Text"

K draw(K x, K y){
 if(TK(x) != 'L'){ unref(x); unref(y); return KE("draw type x"); } //dst
 
 // y-arg: image(draw over) or wh(new all white)
 size_t w, h;
 uint32_t *bg = (uint32_t *)NULL;
 if(TK(y) == 'L'){
  imgk(y, &w, &h, &bg);
  if(bg == NULL){ unref(x); return KE("draw y img"); }
 } else {
  if((TK(y) != 'I')||(NK(y) != 2)) { unref(x); return KE("draw y wh"); }
  int wh[2]; IK(wh, y); w=(size_t)wh[0]; h=(size_t)wh[1];
 }
 
 size_t n = NK(x);
 if(n%2!=0){ unref(x); return KE("draw length y"); }
 
 K *l = malloc(n*sizeof(K));
 LK(l, x);
 
 unsigned int co = 0xff000000;  //color
 int          lw = 1;           //linewidth
 char         cp = 0;           //closepath
 NSVGparser *p = nsvg__createParser(); p->dpi=96; p->image->width=w; p->image->height=h;
 NSVGattrib *attr = nsvg__getAttr(p);
 
 
 if(bg == NULL) drawRect(p, 0.0, 0.0, (float)w, (float)h, 0xffffffff, 1, 1); // fill white bg
 
 float v[4];
 for(int i=0;i<n;i+=2){
  if(TK(l[i])!='s') return drawerr(l,i,n,p,"draw cmd type");
  K a=l[1+i];
  int j=iK(Kx("?", ref(drawcmds), l[i]));
  // printf("drawcmd %d\n", j);
  switch(j){
  case 0: //color
   if(TK(a)!='i') return drawerr(l,1+i,n,p,"draw color");
   co = 0xff000000 | (unsigned int)iK(a);
   break;
  case 1: //font
   if(TK(a)!='C') return drawerr(l,1+i,n,p,"draw font");
   setfont(a);
   break;
  case 2: //linewidth
   if(TK(a)!='i') return drawerr(l,i,1+n,p,"draw linewidth");
   lw = iK(a);
   break;
  case 3: //rect
  case 4: //Rect
   if(vec(v,4,a)) return drawerr(l,1+i,n,p,"draw rect");
   drawRect(p, v[0], v[1], v[2], v[3], co, lw, (char)(j==4));
   break;
  case 5: //circle
  case 6: //Circle
   if(vec(v,3,a)) return drawerr(l,1+i,n,p,"draw circle");
   drawCircle(p, v[0], v[1], v[2], co, lw, (char)(j==6));
   break;
  case 7: //line
   if(vec(v,4,a)) return drawerr(l,1+i,n,p,"draw line");
   float px[2], py[2];
   px[0]=v[0]; px[1]=v[2]; py[0]=v[1]; py[1]=v[3];
   drawPoly(p, px, py, 2, co, lw, 0);
   break;
  case 8: //poly
  case 9: //Poly
   if((TK(a)!='L')||(NK(a)!=2)) return drawerr(l,1+i,n,p,"draw poly");
   K xy[2]; LK(xy, a);
   int nx=vecn(xy[0]);
   if((nx<2)||nx!=vecn(xy[1]))  return drawerr(l,1+i,n,p,"draw poly-x-y");
   float *xf = veca(xy[0], nx);
   float *yf = veca(xy[1], nx);
   drawPoly(p, xf, yf, nx, co, lw, (char)(j==9));
   free(xf); free(yf);
   break;
  case 10: //text
  case 11: //Text     (todo)
   if(bg == NULL)               return drawerr(l,1+i,n,p,"draw text (no bg)");
   if((TK(a)!='L')||(NK(a)!=3)) return drawerr(l,1+i,n,p,"draw text (arg)");
   K  c = Kx("*|", ref(a));
   if(vec(v,2,Kx("2#", a))){unref(c); return drawerr(l,2+i,n,p,"draw text (xy)");}
   if(TK(c)!='C')          {unref(c); return drawerr(l,2+i,n,p,"draw text (text)");}
   size_t nc = NK(c);
   char *cc = (char*)malloc(nc); CK(cc,c);
   drawText(bg, w, h, co, (int)(roundf(v[0])), (int)(roundf(v[1])), cc, nc);
   break;
  default:
   return drawerr(l,1+i,n,p,"draw cmd");
   break;
  }
 }
 
 free(l);
 NSVGimage* im=p->image; p->image=NULL;
 nsvg__deleteParser(p);
 
 // rasterize
 uint32_t *dst = malloc(w*h*4);
 NSVGrasterizer *rst = nsvgCreateRasterizer();
 nsvgRasterize(rst, im, 0,0,1, (uint8_t*)dst, w, h, w*4);
 
 // blend over background image. nsvgRasterize resets dst
 if(bg != NULL){
  drawOver(bg, dst, w*h);
  free(dst); dst=bg;
 }
 
 // return image
 K ri;
 if(4 == sizeof(int)){
  ri = KI((int*)dst, w*h);
 }else{
  int *I = (int*)malloc(w*h*sizeof(int));
  for(int i=0;i<w*h;i++) I[i] = (int)dst[i];
  ri = KI(I, w*h);
  free(I);
 }
 nsvgDeleteRasterizer(rst);
 nsvgDelete(im);
 free(dst);
 K rl[2] = {Ki(h), ri};
 return KL(rl, 2);
}


static uint32_t blend(uint32_t dst, uint32_t src){ // draw src with alpha channel over dst.
 uint32_t a = (src >> 24 ) &0xff;
 uint32_t s = 255 - a;
 uint32_t r = ((( dst     &0xff) * s) + a * ( src     &0xff)) / 255;
 uint32_t g = ((((dst>> 8)&0xff) * s) + a * ((src>> 8)&0xff)) / 255;
 uint32_t b = ((((dst>>16)&0xff) * s) + a * ((src>>16)&0xff)) / 255;
 return 0xff000000|(b<<16)|(g<<8)|r;
}

// draw src over dst. src contains alpha, dst is opaque.
static void drawOver(uint32_t *dst, uint32_t *src, size_t n){ for (int i=0;i<n;i++) dst[i] = blend(dst[i], src[i]); }

// draw char (8-bit alpha bitmap) with color co over dst.
// x and y are the coordinates in dst.
// x0, y0  are the coordinates of the glyph origin within the char bitmap bm.
static void drawChar(uint32_t *dst, int x, int y, size_t w, size_t h, unsigned int co, uint8_t *bm, int cw, int ch, int x0, int y0){
 uint32_t a, sr, sg, sb, sa;
 int o, p, X, Y;
 int wh = w*h;
 co &= 0x00ffffff;
 int out = 0;
 o = 0;
 Y = y + y0;
 for(int j=0;j<ch;j++){
  if((Y>=0)&&(Y<h)){
   X = x - x0;
   p = X + w*Y;
   for(int i=0;i<cw;i++){
    if((X>=0)&&(X<w)){
     sa = (uint32_t)(bm[o++]);
     if((sa)&&((p>=0)&&(p<wh))) dst[p] = blend(dst[p], (sa<<24)|co);
    }
    p++;
   }
  }
  Y++;
  }
}



static K drawerr(K *l, int i, size_t n, NSVGparser *p, char *s){
 for(;i<n;i++) unref(l[i]);
 free(l);
 nsvg__deleteParser(p);
 return KE(s);
}


static int vec(float *v, size_t n, K x){
 int I[4]; double F[4];
 char t=TK(x);
 if(((t!='I')&&(t!='F'))||(NK(x)!=n)) return 1;
 if(t=='I'){
  IK(I,x);
  for (int i=0;i<n;i++) v[i]=(float)I[i];
 }else{
  FK(F,x);
  for (int i=0;i<n;i++) v[i]=(float)F[i];
 }
 return 0;
}
static int vecn(K x){ char t=TK(x); if(((t!='I')&&(t!='F'))) return -1; return (int)NK(x); }
static float *veca(K x, int n){
 float *r = malloc(sizeof(float)*(size_t)n);
 if(TK(x)=='F'){ 
  double *p = dK(x); 
  for(int i=0;i<n;i++) r[i] = (float)p[i];
  unref(x);
 } else { 
  int *p = malloc(sizeof(int)*(size_t)n);
  IK(p, x);
  for(int i=0;i<n;i++) r[i]=(float)p[i];
  free(p);
 }
 return r;
}

static void drawRect(NSVGparser *p, float x, float y, float w, float h, unsigned int co, int lw, char dofill){
 newpath(p);
 nsvg__moveTo(p, x, y);
 nsvg__lineTo(p, x+w, y);
 nsvg__lineTo(p, x+w, y+h);
 nsvg__lineTo(p, x, y+h);
 if(dofill) fill(p, co); 
 else       stroke(p, co, lw, 1);
}
static void drawCircle(NSVGparser *p, float cx, float cy, float r, unsigned int co, int lw, char dofill){	
 newpath(p);
 nsvg__moveTo(p, cx+r, cy);
 nsvg__cubicBezTo(p, cx+r, cy+r*NSVG_KAPPA90, cx+r*NSVG_KAPPA90, cy+r, cx, cy+r);
 nsvg__cubicBezTo(p, cx-r*NSVG_KAPPA90, cy+r, cx-r, cy+r*NSVG_KAPPA90, cx-r, cy);
 nsvg__cubicBezTo(p, cx-r, cy-r*NSVG_KAPPA90, cx-r*NSVG_KAPPA90, cy-r, cx, cy-r);
 nsvg__cubicBezTo(p, cx+r*NSVG_KAPPA90, cy-r, cx+r, cy-r*NSVG_KAPPA90, cx+r, cy);
 if(dofill) fill(p, co); 
 else       stroke(p, co, lw, 1);
}
static void drawPoly(NSVGparser *p, float *px, float *py, int n, unsigned int co, int lw, char dofill) {
 if(n<2) return;
 newpath(p);
 nsvg__moveTo(p, px[0], py[0]);
 for(int i=1;i<n;i++) nsvg__lineTo(p, px[i], py[i]);
 if(dofill) fill(p, co);
 else       stroke(p, co, lw, 0);
}

static void newpath(NSVGparser *p){
 nsvg__pushAttr(p);
 nsvg__resetPath(p);
}

static void fill(NSVGparser *p, unsigned int color){
 NSVGattrib *a = nsvg__getAttr(p);
 a->fillColor = color;
 a->hasStroke = 0;
 a->hasFill = 1;
 nsvg__addPath(p, 1);
 nsvg__addShape(p);
 nsvg__popAttr(p);
}

static void stroke(NSVGparser *p, unsigned int color, float lw, char close){
 NSVGattrib *a = nsvg__getAttr(p);
 a->strokeColor = color;
 a->strokeWidth = lw;
 a->hasStroke = 1;
 a->hasFill = 0;
 nsvg__addPath(p, close);
 nsvg__addShape(p);
 nsvg__popAttr(p);
}

static uint8_t *svgRast(NSVGimage *im) { //(h;I)
	
}

static void wpng(void *context, void *data, int size){
 K *x = (K*)context;
 K r = KC((char*)data, (size_t)size);
 *x = r;
}

static uint32_t *U32I(int *m, size_t n){
 if(4 == sizeof(int)) return (uint32_t*)m;
 uint32_t *u = (uint32_t *)malloc(4*n);
 for(int i=0;i<n;i++) u[i] = (uint32_t)m[i];
 free(m);
 return u;
}

static void imgk(K x, size_t *w, size_t *h, uint32_t **data){
 *data = NULL;
 if((TK(x) != 'L')||(NK(x) != 2)) { unref(x); return; }
 K l[2];
 LK(l, x);
   x = l[0];
 K y = l[1];
 if((TK(x) != 'i') || (TK(y) != 'I')){ unref(x); unref(y); return; }
 
 size_t height  = iK(x);
 size_t n       = NK(y);
 size_t width   = n / height;
 if(n != width * height){ unref(y); return; }
 
 int *I = malloc(n*sizeof(int));
 IK(I, y);
 *data = U32I(I, n);
 *w = width; *h = height;
}


void loadimg(){
 drawcmds = Kx("`color`font`linewidth`rect`Rect`circle`Circle`line`poly`Poly`text`Text");
 fontnames = KS(NULL, 0);
 KR("png", (void*)png, 1);
 KR("draw", (void*)draw, 2);
 KR("loadfont", (void*)loadfont, 2);
}
