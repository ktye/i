#include<stdlib.h>
#include<string.h>
#include<cairo.h>
#include"../k.h"

#include<stdio.h>


static void    imgk(K x, size_t *w, size_t *h, uint32_t **data);
static int     vec(double *v, size_t n, K x);
static int     vecn(K x);
static double *veca(K x, int n);
static void    rgb24(uint32_t *u, size_t n);
static void    align(cairo_t*, double *, int, int, cairo_text_extents_t*);



// write png stream (catenate)
static cairo_status_t wpng(void *p, const unsigned char *d, unsigned int n){ *(K*)p = Kx(",", *(K*)p, KC((char*)d, (size_t)n)); return CAIRO_STATUS_SUCCESS; }


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


static char *c0K(K x){
 size_t n = NK(x);
 char *r = malloc(1+n);
 CK(r, x);
 r[n] = (char)0;
 return r;
}


K drawcmds; // "`color`font`linewidth`rect`Rect`circle`Circle`clip`line`poly`Poly`text`Text"


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
  if(TK(l[i])!='s') { i--; err="draw cmd type"; goto E; }
  K a=l[1+i];
  int j=iK(Kx("?", ref(drawcmds), l[i]));
  //printf("drawcmd %d\n", j);
  cairo_new_path(cr);
  switch(j){
  case 0: //color
   if(TK(a)!='i') { err="draw color"; goto E; }
   unsigned int co = (unsigned int)iK(a);
   cairo_set_source_rgb(cr, (double)(co&0xff)/255.0, (double)((co&0xff00)>>8)/255.0, (double)((co&0xff0000)>>16)/255.0);
   break;
  case 1: //font
   if((TK(a)!='L')||(NK(a)!=2)) { err="draw font arg"; goto E; }
   K l2[2]; LK(l2, ref(a));
   if((TK(l2[0])!='C')||(TK(l2[1])!='i')){unref(l[0]);unref(l[1]);err="draw font args"; goto E; };
   unref(a);
   char *family = c0K(l2[0]);
   cairo_select_font_face(cr, family, CAIRO_FONT_SLANT_NORMAL, CAIRO_FONT_WEIGHT_NORMAL);
   cairo_set_font_size(cr, (double)iK(l2[1]));
   free(family);
   break;
  case 2: //linewidth
   if(TK(a)=='i')      cairo_set_line_width(cr, (double)iK(a));
   else if(TK(a)=='f') cairo_set_line_width(cr, fK(a));
   else { err="draw linewidth"; goto E; }
   break;
  case 3: //rect
  case 4: //Rect
   if(vec(v,4,a)) { err="draw rect"; goto E; }
   cairo_rectangle(cr, v[0], v[1], v[2], v[3]);
   fillstroke(cr, j==4);
   break;
  case 5: //circle
  case 6: //Circle
   if(vec(v,3,a)) { err="draw circle"; goto E; }
   cairo_arc(cr, v[0], v[1], v[2], 0, 6.283185307179586 );
   cairo_close_path(cr);
   fillstroke(cr, j==6);
   break;
  case 7: //clip
   cairo_reset_clip(cr);
   if(vec(v,3,a)){
    if(vec(v,4,a)){ err="clip"; goto E; }
    cairo_rectangle(cr, v[0], v[1], v[2], v[3]);
   }else{
    cairo_new_path(cr);
    cairo_arc(cr, v[0], v[1], v[2], 0.0, 6.283185307179586 );
   //cairo_close_path(cr);
   }
   cairo_clip(cr);
   break;
  case 8: //line
   if(vec(v,4,a)) { err="draw line"; goto E; }
   cairo_move_to(cr, v[0], v[1]);
   cairo_line_to(cr, v[2], v[3]);
   cairo_stroke(cr);
   break;
  case 9:  //poly
  case 10: //Poly
   if((TK(a)!='L')||(NK(a)!=2)) { err="draw poly"; goto E; }
   K xy[2]; LK(xy, a);
   int nx=vecn(xy[0]);
   if((nx<2)||nx!=vecn(xy[1]))  { err="draw poly xy"; goto E; }
   double *xf = veca(xy[0], nx);
   double *yf = veca(xy[1], nx);
   cairo_move_to(cr, xf[0], yf[0]);
   for(int i=1;i<nx;i++) cairo_line_to(cr, xf[i], yf[i]);
   fillstroke(cr, j==10);
   free(xf); free(yf);
   break;
  case 11: //text
  case 12: //Text
   if((TK(a)!='L')||(NK(a)!=3)) { err="draw text arg"; goto E; }
   K l3[3]; LK(l3, ref(a));
   if(vec(v,2,l3[0])||(TK(l3[1])!='i')||(TK(l3[2])!='C')){
    unref(l[1]);unref(l[2]);err="draw text args";goto E; 
   }
   unref(a);
   char *cc = c0K(l3[2]);
   cairo_text_extents_t ex;
   cairo_text_extents(cr, cc, &ex);
   align(cr, v, iK(l3[1]), j==12, &ex);
   cairo_move_to(cr, v[0], v[1]);
   if(j==12){ cairo_save(cr); cairo_rotate(cr,-1.5707963267948966); }
   cairo_show_text(cr, cc);
   if(j==12) cairo_restore(cr);
   free(cc);
   break;
  default:
   i--; err="draw cmd"; goto E;
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

// 6   5   4
// 7abc8def3
// 0   1   2
static void align(cairo_t *cr, double *v, int a, int rot, cairo_text_extents_t *ex){
 double x = ex->x_bearing;
 double y = ex->y_bearing;
 double w = ex->width;
 double h = ex->height; 
 double X[9] = {0, w/2.0, w, w, w, w/2.0, 0, 0, w/2.0};
 double Y[9] = {0, 0, 0, h/2.0, h, h, h, h/2.0, h/2.0};
 double dx = X[a] + x;
 double dy = (y+h) - Y[a];
 if(rot){ double t=dx;dx=dy;dy=-t; }
 v[0] -= dx;
 v[1] -= dy;
}


static int vec(double *v, size_t n, K x){
 int I[4];
 char t=TK(x);
 if(((t!='I')&&(t!='F'))||(NK(x)!=n)) return 1;
 if(t=='F'){ FK(v,x); }
 else{
  IK(I,x);
  for (int i=0;i<n;i++) v[i]=(double)I[i];
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

// flip between rgba and cairos rgb24 (alpha ignored).
static void rgb24(uint32_t *u, size_t n){ for(int i=0;i<n;i++) u[i] = ((u[i]&0xff)<<16) | ((u[i]&0xff0000)>>16) | u[i]&0xff00; }



void loaddrw(){
 drawcmds = Kx("`color`font`linewidth`rect`Rect`circle`Circle`clip`line`poly`Poly`text`Text");
 const char *v = cairo_version_string();
 KA(Ks("cairoversion"), KC((char*)v, strlen(v)));
 //fontnames = KS(NULL, 0);
 KR("png", (void*)png, 1);
 KR("draw", (void*)draw, 2);
 //KR("loadfont", (void*)loadfont, 2);
}
