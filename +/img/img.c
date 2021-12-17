#include<stdlib.h>
#include"../k.h"

#include<stdio.h>

/* http://nothings.org/stb/stb_image_write.h */
#define STBI_WRITE_NO_STDIO
#define STBI_ONLY_PNG
#define STB_IMAGE_WRITE_IMPLEMENTATION
#include "stb_image_write.h" 

/* http://nothings.org/stb/stb_truetype.h */
#define STB_TRUETYPE_IMPLEMENTATION 
#include "stb_truetype.h" 



int32_t *ctx = (int32_t*)NULL;
int32_t  ctxwh[2];

void wpng(void *context, void *data, int size){
	K *x = (K*)context;
	K r = KC((char*)data, (size_t)size);
	*x = r;
}

// create new drawing context
//  x: I (w;h)
//  r: i  w*h
K newctx(K x){
	if(TK(x) != 'I'){ unref(x); return KE("newctx type"); }
	if(NK(x) !=  2) { unref(x); return KE("newctx length"); }
	int wh[2];
	IK(wh, x);
	if(wh[0] <= 0 || wh[1] <= 0){ return KE("newctx value"); }
	ctxwh[0] = wh[0];
	ctxwh[1] = wh[1];
	if(ctx != 0) free(ctx);
	size_t n = 4*(size_t)wh[0]*(size_t)wh[1];
	ctx = (int32_t*)malloc(n);
	memset(ctx, 0, n);
	return Ki(wh[0]*wh[1]);
}

// convert pixels to png data.
//  x: I pixel data, height: (#x)%y
//  y: i width
//  r: C png bytes
K png(K x, K y){
	uint8_t *data;
	size_t   width, height;
	if((TK(x) == 'i') && (TK(y) == 'i')){ // png[0;0] uses ctx
		unref(x); unref(y);
		if(ctx == NULL) return KE("png: empty ctx");
		width  = ctxwh[0];
		height = ctxwh[1];
		data = (uint8_t*)malloc(3*width*height);
		for(int i=0;i<width*height;i++){
			data[0+3*i] = (uint8_t)(ctx[i]&0x0000ff);
			data[1+3*i] = (uint8_t)(ctx[i]&0x00ff00);
			data[2+3*i] = (uint8_t)(ctx[i]&0xff0000);
		}
	}else{
		if(TK(x) != 'I') { unref(x); unref(y); return KE("png type x"); }
		if(TK(y) != 'i') { unref(x); unref(y); return KE("png type y"); }
	
		width = iK(y);
		size_t n = NK(x);
		height = n / width;
		if(n != width * height){ unref(x); return KE("x length"); }
	
		int *I = malloc(n*sizeof(int));
		IK(I, x);
	
		data = (uint8_t*)malloc(3*width*height);
		for(int i=0;i<n;i++){
			data[0+3*i] = (uint8_t)(I[i]&0x0000ff);
			data[1+3*i] = (uint8_t)(I[i]&0x00ff00);
			data[2+3*i] = (uint8_t)(I[i]&0xff0000);
		}
		free(I);
	}
	
	K r;
	stbi_write_png_to_func(wpng, &r, (int)width, (int)height, 3, data, 3*width);
	free(data);
	return r;
}

//void stbi_write_func(void *context, void *data, int size);
//int stbi_write_png_to_func(stbi_write_func *func, void *context, int w, int h, int comp, const void  *data, int stride_in_bytes);


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
 KR("setfont", (void*)setfont, 1);
 KR("newctx", (void*)newctx, 1);
}
