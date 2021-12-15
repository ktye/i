#include<stdlib.h>
#include<string.h>
#include"../k.h"

#include<stdio.h> //printf


// dgesv solve linear system (real)
//  x: L columns (input matrix)
//  y: F or L    (rhs or multi-rhs)
//  r: F or L    (result)
// zgesv                     (complex)
//  x: L columns (input matrix)
//  y: Z or L    (rhs or multi-rhs)
//  r: Z or L    (result)
//
// dgels solve overdetermined system (real)
//  x: L columns (input matrix)
//  y: F or L    (rhs or multi-rhs)
//  r: F or L    (result)
// zgels                             (complex)
//  x: L columns (input matrix)
//  y: Z or L    (rhs or multi-rhs)
//  r: Z or L    (result)
//
// dsyev eigenvalues                 (real symmetric)
//  x: L columns (square input matrix, using upper triag)
//  r: F eigenvalues
// dsyeV eigenvalues and vectors     (real symmetric)
//  x: L columns (square input matrix, using upper triag)
//  r: (F eigenvalues;L eigenvectors)
// zheev eigenvalues                 (complex hermitian)
//  x: L columns (square input matrix, using upper triag)
//  r: F eigenvalues
// zheeV eigenvalues and vectors     (complex hermitian)
//  x: L columns (square input matrix, using upper triag)
//  r: (F eigenvalues;LZ eigenvectors)
// dgeev eigenvalues                 (real unsymmetric)
//  x: L columns (square input matrix)
//  r: Z eigenvalues
// dgeeV eigenvalues and vectors     (real unsymmetric)
//  x: L columns (square input matrix)
//  r: (Z eigenvalues;L left-eigenvectors;L right-eigenvectors)
// zgeev eigenvalues                 (complex unsymmetric)
//  x: L columns (square input matrix)
//  r: Z eigenvalues
// zgeeV eigenvalues and vectors     (complex unsymmetric)
//  x: L columns (square input matrix)
//  r: (Z eigenvalues;L left-eigenvectors;L right-eigenvectors)
//
// todo: dgesvd dgesvD zgesvd zgesvD  
//
// complex numbers
//  backends that do not support complex types natively use F:
//  r0 i0 r1 i1 ..

// Example
//  A:3^?12    /4x3 matrix col-major
//  B:?4
//  dgels[A;B]
// lapack versions dgesv zgesv dgels zgels should return the same as solve[x;y] from ktye/z.k

/*
A:4^?16;b:?4;B:(?4;?4);dgesv[A;b];dgesv[A;B]
A:4^?32;b:?8;B:(?8;?8);zgesv[A;b];zgesv[A;B]
A:3^?12;b:?4;B:(?4;?4);dgels[A;b];dgels[A;B]
A:3^?24;b:?8;B:(?8;?8);zgels[A;b];zgels[A;B]
A:+(1.96 -6.49 -0.47 -7.20 -0.65;-6.49 3.80 -6.39 1.50 -6.34;-0.47 -6.39 4.17 -1.51 2.67;-7.20 1.50 -1.51 5.70 1.80;-0.65 -6.34 2.67 1.80 -7.10)
*/

int LAPACKE_dgesv(int, int, int, double*, int, int*, double*, int);
int LAPACKE_zgesv(int, int, int, double*, int, int*, double*, int);
int LAPACKE_dgels(int, char, int, int, int, double*, int, double*, int);
int LAPACKE_zgels(int, char, int, int, int, double*, int, double*, int);
int LAPACKE_dsyev(int, char, char, int, double*, int, double*);
int LAPACKE_zheev(int, char, char, int, double*, int, double*);
int LAPACKE_dgeev(int, char, char, int, double*, int, double*, double*, double*, int, double*, int);
int LAPACKE_zgeev(int, char, char, int, double*, int, double*, double*, double*, int, double*, int);


static K* Lk(K x, size_t *n){
	  *n = NK(x);
	K *r = (K*)malloc(sizeof(K)**n);
	LK(r, x);
	return r;
}
static void ul(K *l, size_t n){
	for(int i=0; i<n; i++) unref(l[i]);
	free(l);
}
static double *Fk(K x, size_t *n){
	if(TK(x) != 'F') return NULL;
	*n=NK(x);
	double *r = malloc(8**n);
	memcpy(r, (double*)dK(x), 8**n);
	unref(x);
	return r;
}
#ifdef KTYE
extern char *_M;
#endif
static K ZF(K x){
#ifdef KTYE
	size_t n = NK(x);
	printf("ZF n=%d\n", n);
	*(int32_t*)(_M + (int32_t)x - 12) = (int32_t)(n / 2);
	x = (((K)22)<<59)|(K)(int32_t)x;
	printf("ZF: xp=%d xt=%c x=%ld n=%d\n", (int32_t)x, TK(x), x, NK(x)); 
#endif
	return x;
}
static K KZ(double *x, size_t n, size_t z) {
	K r = KF(x, n);
#ifdef KTYE
	if(z > 1) r = ZF(r);
#endif
	return r;
}
static K KZ2(double *re, double *im, size_t n){
	K r = ZF(KF(NULL, 2*n));
	double *p = dK(r);
	for(int i=0;i<n;i++){
		p[0] = re[i];
		p[1] = im[i];
		p += 2;
	}
	return r;
}

static size_t rect(K *r, K x, size_t *cols){ // r:,/x  rows:#*x  cols:#x
	if(TK(x) != 'L'){ unref(x); return 0; }
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
static size_t square(K *r, K x, size_t z){ // r:,/x cols:#x
	size_t cols;
	size_t rows = rect(r, x, &cols);
	if(!rows) return rows;
	rows /= z;
	if(rows != cols){ unref(x); return 0; }
	return rows;
}
static K Ksq(double *x, size_t n, size_t z) {
	K *l = (K*)malloc(n*sizeof(K));
	for(int i=0;i<n;i++) l[i] = KZ(x+i*z*n, z*n, z); 
	K r = KL(l, n);
	free(l);
	return r;
}

K solve(K x, K y, size_t z, int sq){
	size_t rows, cols;
	rows = rect(&x, x, &cols);
	if(cols == 0){ unref(y); return KE("gels: rect A"); }
	
	if(z*(rows/z) != rows){ unref(x); unref(y); return KE("gels: zlength A"); }
	rows /= z;
	if(rows < cols){ unref(x); unref(y); return KE("gels: length A"); }
	if((rows != cols)&&sq){ unref(x); unref(y); return KE("gesv: A square"); }
	
	double *A, *B;
	size_t xn;
	A = Fk(x, &xn);
	if(A == NULL){ unref(y); return KE("gels: type A"); }

	int rl = 0;
	size_t nrhs = 1;
	char yt = TK(y);
	if(yt == 'L'){ //multi-rhs
		rl = 1;
		size_t yn = rect(&y, y, &nrhs);
		if(nrhs == 0) { return KE("gels: rect B"); }
		yn /= z;
		if(yn != rows){ unref(y); return KE("gels: rows B"); }
		yt = TK(y);
	}
	K r;
	if(yt == 'F'){
		size_t yn;
		B = Fk(y, &yn);
		if(B == NULL){ free(A); return KE("gels: type B"); }
		yn /= z;
		if(yn/nrhs != rows){ free(A); free(B); return KE("gels: length B"); }
		int info;
		if(sq){
			int *ipiv = malloc(rows*sizeof(int));
			if(z > 1) info = LAPACKE_zgesv(102, (int)rows, (int)nrhs, A, (int)rows, ipiv, B, (int)rows);
			else      info = LAPACKE_dgesv(102, (int)rows, (int)nrhs, A, (int)rows, ipiv, B, (int)rows);
			free(ipiv);
		}else{
			if(z > 1) info = LAPACKE_zgels(102, 'N', (int)rows, (int)cols, (int)nrhs, A, (int)rows, B, (int)rows);
			else      info = LAPACKE_dgels(102, 'N', (int)rows, (int)cols, (int)nrhs, A, (int)rows, B, (int)rows);
		}
		if(info!=0){ free(A); free(B); return KE("lapack-dgels"); }
		free(A);
		if(rl) {
			K *l = (K*)malloc(nrhs * sizeof(K));
			for(int i=0;i<nrhs;i++){
				l[i] = KZ(B+i*z*rows, z*cols, z);
			}
			r = KL(l, nrhs);
			free(l);
		} else {
			r = KZ(NULL, z*cols, z);
			memcpy(dK(r), B, 8*z*cols);
		}
		free(A); free(B);
		return r;
	} else {
		unref(y); free(A); return KE("type dgels-B?");
	}
}
K eig(K x, size_t z, char jobz, int sym){
	printf("eig x=%ld z=%d jobz=%c sym=%d\n", x, z, jobz, sym);
	size_t cols = square(&x, x, z);
	if(cols == 0){ return KE("eig A square"); }
	
	size_t xn;
	double *A = Fk(x, &xn);
	if(A == NULL){ return KE("eig: type A"); }
	
	int info;
	double *wr = (double*)malloc(cols*sizeof(double));
	double *wi = (double*)NULL;
	double *vl = (double*)NULL;
	double *vr = (double*)NULL;
	int n = (int)cols;
	if(sym){
		if(z > 1) info = LAPACKE_zheev(102, jobz, 'U', n, A, n, wr);
		else      info = LAPACKE_dsyev(102, jobz, 'U', n, A, n, wr);
	}else{
		if(jobz == 'V'){
			printf("alloc vl/vr %d\n", 8*z*cols);fflush(stdout);
			vl = (double*)malloc(8*z*cols*cols);
			vr = (double*)malloc(8*z*cols*cols);
		}
		printf("eig unsym z=%d n=%d\n", z, n);fflush(stdout);
		wi = (double*)malloc(8*cols);
		if(z > 1) info = LAPACKE_zgeev(102, jobz, jobz, n, A, n, wr, wi, vl, n, vr, n);
		else      info = LAPACKE_dgeev(102, jobz, jobz, n, A, n, wr, wi, vl, n, vr, n);
		printf("lapack info=%d\n", info); fflush(stdout);
	}
	
	K r, rw;
	if(sym){
		K rw = KF(wr, cols);
		free(wr);
	}else{
		printf("assign rw nonsym\n");
		rw = KZ2(wr, wi, cols);
		printf("rw = %ld\n", rw);
		free(wr);
		free(wi);
	}
	if(jobz == 'V'){
		printf("assign vectors\n");
		if(sym){
			K rl[2] = {rw, Ksq(A, cols, z)};
			r = KL(rl, 2);
		}else{
			K rl[3] = {rw, Ksq(vl, cols, z), Ksq(vr, cols, z)};
			free(vl);
			free(vr);
			r = KL(rl, 3);
		}
	} else  r = rw;
	free(A);
	printf("return r=%ld\n", r);
	return r;
}
K dgesv(K x, K y){ return solve(x, y, 1, 1); }
K zgesv(K x, K y){ return solve(x, y, 2, 1); }
K dgels(K x, K y){ return solve(x, y, 1, 0); }
K zgels(K x, K y){ return solve(x, y, 2, 0); }
K dsyev(K x){ return eig(x, 1, 'N', 1); }
K dsyeV(K x){ return eig(x, 1, 'V', 1); }
K zheev(K x){ return eig(x, 2, 'N', 1); }
K zheeV(K x){ return eig(x, 2, 'N', 1); }
K dgeev(K x){ return eig(x, 1, 'N', 0); }
K dgeeV(K x){ return eig(x, 1, 'V', 0); }
K zgeev(K x){ return eig(x, 2, 'N', 0); }
K zgeeV(K x){ return eig(x, 2, 'V', 0); }

void loadmat(){
 KR("dgesv", (void*)dgesv, 2);
 KR("zgesv", (void*)zgesv, 2);
 KR("dgels", (void*)dgels, 2);
 KR("zgels", (void*)zgels, 2);
 KR("dsyev", (void*)dsyev, 1);
 KR("dsyeV", (void*)dsyeV, 1);
 KR("zheev", (void*)zheev, 1);
 KR("zheeV", (void*)zheeV, 1);
 KR("dgeev", (void*)dgeev, 1);
 KR("dgeev", (void*)dgeev, 1);
 KR("dgeeV", (void*)dgeeV, 1);
 KR("zgeev", (void*)zgeev, 1);
 KR("zgeeV", (void*)zgeeV, 1);
}
