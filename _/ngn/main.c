#include<stdio.h>
#include<stdlib.h>
#include<stddef.h> 
#include "a.h"
void init();
A kst(A);
void  *mmap(void *addr, size_t len, int prot, int flags, int filedes, off_t off) {
	printf("mmap addr=%p len=%d prot=%d flags=%d filedes=%d off=%d\n", addr, len, prot, flags, filedes, off);	
	return malloc(len);
}
int    munmap(void *addr, size_t len) {
	printf("munmap addr=%p len=%d\n");
}

int main() {
	int n, i;
	char buf[128];
	char *c;
	A x;
	init();
	while (fgets(buf, 128, stdin) != NULL) {
		for (i=0;i<128;i++) {
			if (buf[i] == '\n')
				buf[i] = 0;
		}
		buf[127] = 0;
		n = strlen(buf);
		x = aC(n);
		for(i=0;i<n;i++) xci=buf[i];
		x = val(x);
		if (!x) {
			printf("!\n");
			continue;
		}
		x = kst(x);
		if(xn==2 && xc[0] == ':' && xc[1] == ':') continue;
		for(i=0;i<xn;i++) putchar(xc[i]); putchar('\n');
	}
	return 0;
}
