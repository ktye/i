#include<stdlib.h>
#include<cairo.h>
#include"../k.h"

#include<stdio.h>


static void    imgk(K x, size_t *w, size_t *h, uint32_t **data);
static int     vec(double *v, size_t n, K x);
static int     vecn(K x);
static double *veca(K x, int n);

/*
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
*/


// flip between rgba and cairos rgb24 (alpha ignored).
static void rgb24(uint32_t *u, size_t n){ for(int i=0;i<n;i++) u[i] = ((u[i]&0xff)<<16) | ((u[i]&0xff0000)>>16) | u[i]&0xff00; }

static cairo_status_t wpng(void *p, const unsigned char *d, unsigned int n){ // write png stream (catenate)
 *(K*)p = Kx(",", *(K*)p, KC((char*)d, (size_t)n));
 return CAIRO_STATUS_SUCCESS;
}



// convert pixels to png data.
//  x: (i height;I pixels)
//  r: C png bytes
K png(K x){
 uint32_t *u;
 size_t    width, height;
 imgk(x, &width, &height, &u);
 if(u == NULL) return KE("png: type img");

 rgb24(u, width*height);
 cairo_surface_t *s = cairo_image_surface_create_for_data((unsigned char *)u, CAIRO_FORMAT_RGB24, width, height, 4*width);
 K r = KC(NULL, 0);
 cairo_surface_write_to_png_stream(s, wpng, &r);
 cairo_surface_destroy(s);
 free(u);
 return r;
}

static void fillstroke(cairo_t *cr, int fill){ if(fill) cairo_fill(cr); else cairo_stroke(cr); }



K drawcmds; // "`color`font`linewidth`rect`Rect`circle`Circle`line`poly`Poly`text`Text"


K drawclose(K *l, int i, int n, cairo_t *cr, cairo_surface_t *surf, uint32_t *bg, const char *err);

K draw(K x, K y){ //dst
 if(TK(x) != 'L'){ unref(x); unref(y); return KE("draw type x"); } //dst
 
 cairo_surface_t *surf;
 
 // y-arg: image(draw over) or wh(new all white)
 size_t w, h;
 uint32_t *bg = (uint32_t *)NULL;
 if(TK(y) == 'L'){
  imgk(y, &w, &h, &bg);
  if(bg == NULL){ unref(x); return KE("draw y img"); }
  rgb24(bg, w*h);
  surf = cairo_image_surface_create_for_data((unsigned char *)bg, CAIRO_FORMAT_RGB24, w, h, 4*w);
 } else {
  if((TK(y) != 'I')||(NK(y) != 2)) { unref(x); return KE("draw y wh"); }
  int wh[2]; IK(wh, y); w=(size_t)wh[0]; h=(size_t)wh[1];
  surf = cairo_image_surface_create(CAIRO_FORMAT_RGB24, w, h);
 }
 
 size_t n = NK(x);
 if(n%2!=0){ unref(x); return KE("draw length y"); }
 
 K *l = malloc(n*sizeof(K));
 LK(l, x);
 
 
 cairo_t *cr = cairo_create(surf);
 cairo_set_source_rgb(cr, 1, 1, 1);
 
 if(bg == NULL){ cairo_rectangle(cr, 0, 0, (double)w, (double)h); cairo_fill(cr); } // fill white bg
 cairo_set_source_rgb(cr, 0, 0, 0); // default color black
 cairo_set_line_width(cr, 1);       // default line width
 
 char *err = NULL;
 double v[4];
 int i=0;
 for(;i<n;i+=2){
  if(TK(l[i])!='s') { i--; err="draw cmd type"; goto E; };
  K a=l[1+i];
  int j=iK(Kx("?", ref(drawcmds), l[i]));
  // printf("drawcmd %d\n", j);
  switch(j){
  case 0: //color
   if(TK(a)!='i') { err="draw color"; goto E; };
   unsigned int co = (unsigned int)iK(a);
   cairo_set_source_rgb(cr, (double)(co&0xff)/255.0, (double)((co&0xff00)>>8)/255.0, (double)((co&0xff0000)>>24)/255.0);
   break;
  case 1: //font
   if(TK(a)!='C') { err="draw font"; goto E; };
   //setfont(a);
   printf("font todo..\n");
   break;
  case 2: //linewidth
   if(TK(a)=='i')      cairo_set_line_width(cr, (double)iK(a));
   else if(TK(a)=='f') cairo_set_line_width(cr, fK(a));
   else { err="draw linewidth"; goto E; };
   break;
  case 3: //rect
  case 4: //Rect
   if(vec(v,4,a)) { err="draw rect"; goto E; };
   cairo_rectangle(cr, v[0], v[1], v[2], v[3]);
   fillstroke(cr, j==4);
   break;
  case 5: //circle
  case 6: //Circle
   if(vec(v,3,a)) { err="draw circle"; goto E; };
   cairo_arc(cr, v[0], v[1], v[3], 0, 6.283185307179586);
   fillstroke(cr, j==6);
   break;
  case 7: //line
   if(vec(v,4,a)) { err="draw line"; goto E; };
   cairo_move_to(cr, v[0], v[1]);
   cairo_line_to(cr, v[2], v[3]);
   cairo_stroke(cr);
   break;
  case 8: //poly
  case 9: //Poly
   if((TK(a)!='L')||(NK(a)!=2)) { err="draw poly"; goto E; };
   K xy[2]; LK(xy, a);
   int nx=vecn(xy[0]);
   if((nx<2)||nx!=vecn(xy[1]))  { err="draw poly xy"; goto E; };
   double *xf = veca(xy[0], nx);
   double *yf = veca(xy[1], nx);
   cairo_move_to(cr, xf[0], yf[0]);
   for(int i=1;i<nx;i++) cairo_line_to(cr, xf[i], yf[i]);
   fillstroke(cr, j==9);
   free(xf); free(yf);
   break;
  case 10: //text
  case 11: //Text
   if((TK(a)!='L')||(NK(a)!=3)) { err="draw text arg"; goto E; };
   K  c = Kx("*|", ref(a));
   if(vec(v,2,Kx("2#", a))){unref(c); err="draw text xy"; goto E; };
   if(TK(c)!='C')          {unref(c); err="draw text"; goto E; };
   size_t nc = NK(c);
   char *cc = (char*)malloc(1+nc); CK(cc,c);
   cc[nc] = 0;
   cairo_move_to(cr, v[0], v[1]);
   cairo_show_text(cr, cc);
   free(cc);
   break;
  default:
   i--; err="draw draw text xy"; goto E;
   break;
  }
 }
 
E:
 i++;
 for(;i<n;i++) unref(l[i]);
 free(l);
 
 K r;
 if(err==NULL){
  cairo_surface_flush(surf);
  uint32_t* dst = (uint32_t *)cairo_image_surface_get_data(surf);
  rgb24(dst, w*h);
 
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
  K rl[2] = {Ki(h), ri};
  r = KL(rl, 2);
 }
 
 cairo_destroy(cr);
 cairo_surface_destroy(surf);
 if(bg != NULL) free(bg);
 if(err != NULL) return KE(err);
 return r;
}




static int vec(double *v, size_t n, K x){
 int I[4];
 char t=TK(x);
 if(((t!='I')&&(t!='F'))||(NK(x)!=n)) return 1;
 if(t=='F'){ FK(v,x); }
 else{
  IK(I,x);
  for (int i=0;i<n;i++) v[i]=(float)I[i];
 }
 return 0;
}
static int vecn(K x){ char t=TK(x); if(((t!='I')&&(t!='F'))) return -1; return (int)NK(x); }
static double *veca(K x, int n){
 double *r = malloc(sizeof(double)*(size_t)n);
 if(TK(x)=='F'){ FK(r, x); return r; };
 
 int *p = malloc(sizeof(int)*(size_t)n);
 IK(p, x);
 for(int i=0;i<n;i++) r[i]=(float)p[i];
 free(p);
 return r;
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
 //fontnames = KS(NULL, 0);
 KR("png", (void*)png, 1);
 KR("draw", (void*)draw, 2);
 //KR("loadfont", (void*)loadfont, 2);
}
