#include<stdio.h>
#include<string.h>
#include"../k.h"

K add(K x, K y) { 
 if((TK(x)!='i')||TK(y)!='i') return KE("type");
 return Ki(iK(x)+iK(y));
}
K Add (K x, K y){ return Kx("+", x, y); }
K Flip(K x)     { return Kx("+", x   ); }

void loadlib(const char *prefix) {
 char b[20];
 
 b[0] = 0;
 strncat(b, prefix, sizeof(b)-1);
 strncat(b, "add", 4);
 KR(b, (void*)add, 2);
 
 b[0] = 0;
 strncat(b, prefix, sizeof(b)-1);
 strncat(b, "Add", 4);
 KR(b, (void*)Add, 2);
 
 b[0] = 0;
 strncat(b, prefix, sizeof(b)-1);
 strncat(b, "Flip", 5);
 KR(b, (void*)Flip, 1);
}
