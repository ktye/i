#include<string.h>
#include"../k.h"

K add(K x, K y) { return K2('+', x, y); }

void loadlib(const char *prefix) {
 char b[20];
 b[0] = 0;
 strncat(b, prefix, sizeof(b)-1);
 strncat(b, "add", 4);
 KR(b, (void*)add, 2);
}
