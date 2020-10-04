#include"a.h"

extern unsigned char __heap_base;
char *mem = &__heap_base;
char *mmap(void *addr, int len, int prot, int flags, int filedes, int off) {
	char *r = mem;
	mem += len;
	return r;
}
int   munmap(void *addr, int len) { return 0; }
A   mk(int n) { return aC(n); }
int nn(A x) { return xn; }
void  exit(){}
