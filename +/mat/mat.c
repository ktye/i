#include<stdlib.h>
#include<string.h>
#include"../k.h"

#include<stdio.h> //printf

// dgels solve overdetermined system (real)
//  x: L columns (input matrix)
//  y: F or L    (rhs or multi-rhs)
//  r: F or L    (result)

// Example
//  A:3^?12    /4x3 matrix col-major
//  B:?4
//  solve[A;B] /builtin z.k
//  dgels[A;B] /lapack version

// A:3^?12;b:?4;B:(?4;?4);dgels[A;b] /dgels[A;B]

int LAPACKE_dgels(int, char, int, int, int, double*, int, double*, int);

static K* Lk(K x, size_t *n){
	  *n = NK(x);
	K *r = (K*)malloc(sizeof(K)**n);
	LK(r, x);
	return r;
}
static void ul(K *l, size_t n){
	for(int i=0;i<n;i++) unref(l[i]);
	free(l);
}
static double *Fk(K x, size_t *n){
	if(TK(x)!='F') return NULL;
	*n=NK(x);
	double *r = malloc(8**n);
	memcpy(r, (double*)dK(x), 8**n);
	unref(x);
	return r;
}

static size_t rect(K *r, K x, size_t *cols) { // r:,/x  rows:#*x  cols:#x
	if(TK(x)!='L') { unref(x); return 0; }
	*cols = NK(x);
	x = Kx(",/", x);
	size_t n = NK(x);
	size_t rows = n / *cols;
	if (n != rows**cols) {
		unref(x);
		return 0;
	}
	*r = x;
	return rows;
}


K dgels(K x, K y){
	size_t rows, cols;
	rows = rect(&x, x, &cols);
	if(cols == 0) { unref(y); return KE("rect dgels-A");  }
	if(rows<cols) { unref(x); unref(y); return KE("length dgels-A"); }
	
	double *A, *B;
	size_t xn;
	A = Fk(x, &xn);
	if(A==NULL) { unref(y); return KE("type dgels-A"); }

	int rl = 0;
	size_t nrhs = 1;
	char yt=TK(y);
	if(yt=='L') {
		rl = 1;
		size_t yn;
		yn = rect(&y, y, &nrhs);
		if(nrhs==0) { return KE("rect dgels-B"); }
		if(yn!=rows){ unref(y); return KE("rows dgels-B"); }
		yt=TK(y);
	}
	K r;
	if(yt=='F') {
		size_t yn;
		B = Fk(y, &yn);
		if(B==NULL){free(A); return KE("type dgels-B");}
		int info = LAPACKE_dgels(102, 'N', (int)rows, (int)cols, (int)nrhs, A, (int)rows, B, (int)rows);
		if(info!=0){free(A); free(B); return KE("lapack-dgels");}
		free(A);
		if(rl) {
			K *l = (K*)malloc(nrhs * sizeof(K));
			for(int i=0;i<nrhs;i++){
				l[i] = KF(B+i*rows, cols);
			}
			r = KL(l, nrhs);
			free(l);
		} else {
			r = KF(NULL, cols);
			memcpy(dK(r), B, 8*cols);
		}
		free(A); free(B);
		return r;
	} else {
		unref(y); free(A); return KE("type dgels-B?");
	}
}

void loadmat(){
 KR("dgels", (void*)dgels, 2);
}
