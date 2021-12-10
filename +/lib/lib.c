#include<stdio.h>
#include<string.h>
#include"../k.h"

K add(K x, K y) { 
 printf("add.. xt=%c yt=%c %d+%d\n", TK(x), TK(y), iK(x), iK(y));fflush(stdout);
 if((TK(x)!='i')||TK(y)!='i') return KE("type");
 return Ki(iK(x)+iK(y));
}
K Add(K x, K y) { return K2('+', x, y); }

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
}
