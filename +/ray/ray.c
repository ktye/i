#include<stdlib.h>
#include<string.h>
#include"raylib.h"
#include"../k.h"

K raywin(K width, K height, K frame){
	if('i' != TK(width)){  return KE("type raywin width");  }
	if('i' != TK(height)){ return KE("type raywin height"); }
	if('s' != TK(frame)){  return KE("type raywin frame");  }
	InitWindow(iK(width), iK(height), "k");
	unref(width);
	unref(height);
	while(!WindowShouldClose()){
		BeginDrawing();
		ClearBackground(RAYWHITE);
		DrawText("alpha", 190, 200, 20, LIGHTGRAY);
		EndDrawing();
	}
	CloseWindow();
	unref(frame);
	return Ki(0);
}
void loadray(){
	KR("raywin", raywin, 3);
}
