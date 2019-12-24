#include<stdio.h>
#include<stdlib.h>

typedef long long I;
I* til(I* a, I n) {
	I i;
	for (i = 0; i<n; i++) a[i] = i;
	return a;
}
I sum(I* a, I n) {
	I i, r;
	r = 0;
	for (i = 0; i<n; i++) r += a[i];
	return r;
}
int main(int args, char **argv) {
	I  ns, np;
	I* a;
	ns = atoi(argv[1]);
	np = atoi(argv[2]);
	a = calloc(np*ns, sizeof(I));
	printf("%lld\n", sum(til(a, ns*np), ns*np));
}
