#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include"raylib.h"
#include"../k.h"


K png(K); // ../img/img.c

//show(50;10000#255) /red window 50x100
K show(K x){
	if(TK(x) == 'L') x = png(x);
	if(TK(x) != 'C'){ return KE("show: type x"); }
	Image m = LoadImageFromMemory(".png", (const unsigned char*)dK(x), (int)NK(x));
	int w = m.width;
	int h = m.height;
	unref(x);
	if(w*h <= 0){
		UnloadImage(m);
		return Ki(1);
	}

	InitWindow(w, h, "k+");


	int    mx=0, my=0, lx, ly, rw, rh;
	double mt=0, lt;


	Texture2D t = LoadTextureFromImage(m);
	UnloadImage(m);
	while(!WindowShouldClose()){

		if(IsMouseButtonPressed(0)){
			lx=mx;          ly=my;          lt=mt;
			mx=GetMouseX(); my=GetMouseY(); mt=GetTime();
			//printf("pressed %f %f %f\n", mx, my, mt);
			if((0.3>mt-lt)&&(10>(mx-lx)*(mx-lx))&&(10>(my-ly)*(my-ly))){
				printf("double-click\n");
			}
		}
		if(IsMouseButtonDown(0)){
			int x=GetMouseX();
			int y=GetMouseY();
			rw=x-mx;
			rh=y-my;
			SetMouseCursor(MOUSE_CURSOR_CROSSHAIR);
		}
		if(IsMouseButtonReleased(0)){
			int x=GetMouseX();
			int y=GetMouseY();
			rw=0; rh=0;
			SetMouseCursor(MOUSE_CURSOR_DEFAULT);
			//printf("released %f %f %f\n", GetMouseX(), GetMouseY(), GetTime());
			if((10<(lx-x)*(lx-x))&&(10<(ly-y)*(ly-y))){
				printf("zoom\n");
			}
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
