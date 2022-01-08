#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include"raylib.h"
#include"../k.h"

K png(K); // ../img/img.c

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
	Texture2D t = LoadTextureFromImage(m);
	UnloadImage(m);
	while(!WindowShouldClose()){
		BeginDrawing();
            	DrawTexture(t, 0, 0, WHITE);
		EndDrawing();
	}
	UnloadTexture(t);
	CloseWindow();
	return Ki(0);
}
void loadray(){
	KR("show", (void*)show, 1);
}
