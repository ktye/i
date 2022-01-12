#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include"raylib.h"
#include"rgestures.h"
#include"../k.h"


K png(K); // ../img/img.c

//show(50;10000#255) /red window 50x100
K show(K x){
	if(TK(x) == 'L') x = png(x);
	if(TK(x) != 'C'){ return KE("show: type x"); }

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
	SetGesturesEnabled(GESTURE_DOUBLETAP|GESTURE_DRAG);


	int mx=0, my=0, rw=0, rh=0, dragging=0;


	Texture2D t = LoadTextureFromImage(m);
	UnloadImage(m);
	while(!WindowShouldClose()){
		if(IsMouseButtonPressed(0)){
			mx=GetMouseX(); my=GetMouseY();
		}

		switch(GetGestureDetected()){
		case GESTURE_DOUBLETAP: 
			printf("doubleclick %d %d\n", mx, my);
			break;
		case GESTURE_DRAG: 
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
				printf("zoom %d %d+%d %d\n", mx, my, rw, rh);
			}
			rw=0; rh=0;
		}

		BeginDrawing();
		DrawTexture(t, 0, 0, WHITE);
		if((rw!=0)&&(rh!=0)) DrawRectangleLines(mx, my, rw, rh, WHITE);
		EndDrawing();
	}
	UnloadTexture(t);
	CloseWindow();
	return Ki(0);
}
void loadray(){
	KR("show", (void*)show, 1);
}
