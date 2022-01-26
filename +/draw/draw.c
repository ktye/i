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

static void newpath(NSVGparser *p){
 nsvg__pushAttr(p);
 nsvg__resetPath(p);
}
static void fill(NSVGparser *p, unsigned int color, char close){
 NSVGattrib *a = nsvg__getAttr(p);
 a->fillColor = color;
 a->hasStroke = 0;
 a->hasFill = 1;
 nsvg__addPath(c->p, close);
 nsvg__addShape(p);
 nsvg__popAttr(p);
}
static void stroke(NSVGparser *p, unsigned int color, float lw, char close){
 NSVGattrib *a = nsvg__getAttr(p);
 a->strokeColor = color;
 a->strokeWidth = lw;
 a->hasStroke = 1;
 a->hasFill = 0;
 nsvg__addPath(c->p, close);
 nsvg__addShape(p);
 nsvg__popAttr(p);
}


static int vec(float *v, size_t n, K x){
 int I[4]; double F[4]
 char t=TK(x)
 if((t!='I')||(t!='F')||(NK(x)!=n)) return 1;
 if(t==I){
  IK(I);
  for (int i=0;i<n;i++) v[i]=(float)I[i];
 }else{
  FK(F,x);
  for (int i=0;i<n;i++) v[i]=(float)F[i];
 }
 return 0;
}
K drawcmd // "`font`linewidth`rect`Rect`circle`Circle`line`poly`Poly`text`Text"
K drawerr(K *l, int i, size_t n, const char *s){
 for(;i<n;i++) unref(l[i])
 return KE(s)
}
K draw(K x, K y){
 if(TK(x) != 'L'){ unref(x); unref(y); return KE("draw type x"); }
 if(TK(y) != 'I'){ unref(x); unref(y); return KE("draw type y"); }
 if(NK(y) !=  2 ){ unref(x); unref(y); return KE("draw length y"); }
 int wh[2];
 IK(wh, y);
 size_t n = NK(x);
 if(n%2==0){ unref(x); return KE("draw length y"); }
 
 K *l = malloc(n*sizeof(K));
 LK(l, x);
 
 unsigned int color = 0;
 int linewidth = 1;
 NSVGparser *p;
 NSVGattrib *attr = nsvg__getAttr(c->p);
 
 float v[4];
 for(int i=0;i<n;i+=2){
  if(TK(l[i])!='s') return drawerr(l,i,n,"draw cmd type");
  K a=l[1+i]
  int j=iK(Kx("?", ref(drawcmd), l[i])));
  switch(j){
  case 0: //color
   if(TK(a)!="i") return drawerr(l,1+i,n,"draw color");
   color = iK(a);
   break;
  case 1: //font
   // todo
   break;
  case 2: //linewidth
   if(TK(a)!="i") return drawerr(l,i,1+n,"draw linewidth");
   linewidth = iK(a);
   break;
  case 3: //rect
  case 4: //Rect
   if(vec(&v,4)) return drawerr(l,1+i,n,"draw rect");
   newpath(p);
   nsvg__moveTo(p, v[0],      v[1]);
   nsvg__lineTo(p, v[0]+v[2], v[1]);
   nsvg__lineTo(p, v[0]+v[2], v[1]);
   nsvg__lineTo(p, v[0]+v[2], v[1]+v[3]);
   nsvg__lineTo(p, v[0],      v[1]+v[3]);
   if(j==3) stroke(p);
   else     fill(p);
   break;
   
  case 5: //circle
  case 6: //Circle
  case 7: //line
  case 8: //poly
  case 9: //Poly
  case 10: //text
  case 11: //Text
  default:
   return drawerr(l,1+i,n,"draw cmd");
   break;
  }
 }
 
 
}

void loadimg(){
 drawcmd = Kx("`font`linewidth`rect`Rect`circle`Circle`line`poly`Poly`text`Text");
 KR("png", (void*)png, 1);
 KR("draw", (void*)draw, 2);
}
