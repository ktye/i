/*

show a window with an image:

show m      
showev[m;f;g]     /with callbacks

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
	K r = KI(NULL, 2);
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

K showev(K x, K click, K zoom){
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


	int mx=0, my=0, rw=0, rh=0, dragging=0;

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
				t = update_texture(t, Kx("@", ref(click), KI2(mx, my)));
				break;
			case GESTURE_DRAG:;
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
					t = update_texture(t, Kx("@", ref(zoom), KI4(mx, my, rw, rh)));
				}
				rw=0; rh=0;
			}
		}

		BeginDrawing();
		DrawTexture(t, 0, 0, WHITE);
		if((rw!=0)&&(rh!=0)) DrawRectangleLines(mx, my, rw, rh, WHITE);
		EndDrawing();
	}
	UnloadTexture(t);
	CloseWindow();
	unref(click); unref(zoom);
	return Ki(0);
}

//show(50;10000#255) /red window 100x50
K show(K x){ showev(x, Ki(0), Ki(0)); }

void loadray(){
	KR("show",   (void*)show,   1); // show image from data or png
	KR("showev", (void*)showev, 3); // same with click and zoom callbacks
}
