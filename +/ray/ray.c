/*

show a window with an image:

show m      
Show[m;f;g]     /with callbacks

image m is given as png bytes(C) or (height;dataI)

callback functions f and g (both monadic) return a new image

f x: double-click                x:xi yi
g x: zoom(click hold and move)   x:xi yi wi hi

*/
#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include"raylib.h"
#include"rgestures.h"
#include"../k.h"


K png(K); // ../img/img.c

K KI2(int x, int y) {
	K r = KI(NULL, 2);
	int *p = dK(r);
	p[0] = x;
	p[1] = y;
	return r;
}
K KI4(int x, int y, int u, int v) {
	K r = KI(NULL, 4);
	int *p = dK(r);
	p[0] = x;
	p[1] = y;
	p[2] = u;
	p[3] = v;
	return r;
}

Texture2D update_texture(Texture2D t, K x){
	if(TK(x) == 'L') x = png(x);
	if(TK(x) != 'C'){ unref(x); return t; }
	
	Image m = LoadImageFromMemory(".png", (const unsigned char*)dK(x), (int)NK(x));
	int w = m.width;
	int h = m.height;
	unref(x);
	if(w*h <= 0){ UnloadImage(m); return t; }
	UnloadTexture(t);
	return LoadTextureFromImage(m);
}

K Show(K x, K click, K zoom){
	if(TK(x) == 'L') x = png(x);
	if(TK(x) != 'C'){ return KE("show: type x"); }

	int interactive = 0;
	if((TK(click) != 'i')&&(TK(zoom) != 'i')) interactive = 1;

	SetTraceLogLevel(LOG_WARNING); // suppress info(raylib)

	Image m = LoadImageFromMemory(".png", (const unsigned char*)dK(x), (int)NK(x));
	int w = m.width;
	int h = m.height;
	unref(x);
	if(w*h <= 0){
		UnloadImage(m);
		return Ki(1);
	}

	InitWindow(w, h, "k+");
	if(interactive) SetGesturesEnabled(GESTURE_DOUBLETAP|GESTURE_DRAG);


	int px=0, py=0, mx=0, my=0, rw=0, rh=0, dragging=0;

	Texture2D t = LoadTextureFromImage(m);
	UnloadImage(m);
	while(!WindowShouldClose()){
		if(interactive){
			if(IsMouseButtonPressed(0)){
				mx=GetMouseX(); my=GetMouseY();
			}

			switch(GetGestureDetected()){
			case GESTURE_DOUBLETAP: 
				//printf("doubleclick %d %d\n", mx, my);
				px=mx; py=my;
				K im = Kx(".", ref(click), KI2(mx, my));
				if(TK(im)=='L') t = update_texture(t, im);
				else unref(im);
				break;
			case GESTURE_DRAG:;
			        px=0;py=0;
				int x=GetMouseX();
				int y=GetMouseY();
				rw=x-mx;
				rh=y-my;
				SetMouseCursor(MOUSE_CURSOR_CROSSHAIR);
				dragging = 1;
				break;
			default: break;
			}

			if(IsMouseButtonReleased(0)){
				if(dragging){
					SetMouseCursor(MOUSE_CURSOR_DEFAULT);
					dragging = 0;
					//printf("zoom %d %d+%d %d\n", mx, my, rw, rh);
					K im = Kx(".", ref(zoom), KI4(mx+((rw<0)?rw:0), my+((rh<0)?rh:0), (rw<0)?-rw:rw, (rh<0)?-rh:rh));
					if(TK(im)=='L') t = update_texture(t, im);
					else unref(im);
				}
				rw=0; rh=0;
			}
		}

		BeginDrawing();
		DrawTexture(t, 0, 0, WHITE);
		if((px!=0)&&(py!=0)) DrawCircle(mx, my, 3, RED);
		if((rw!=0)&&(rh!=0)) DrawRectangleLines(mx+((rw<0)?rw:0), my+((rh<0)?rh:0), (rw<0)?-rw:rw, (rh<0)?-rh:rh, RED);
		EndDrawing();
	}
	UnloadTexture(t);
	CloseWindow();
	unref(click); unref(zoom);
	return Ki(0);
}

//show(50;10000#255) /red window 100x50
K show(K x){ Show(x, Ki(0), Ki(0)); }

void loadray(){
	KR("show", (void*)show, 1); // show image from data or png
	KR("Show", (void*)Show, 3); // same with click and zoom callbacks
}
