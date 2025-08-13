/*
#include<stdlib.h>
#include<string.h>
#include"../k.h"

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
//  x: L columns (square symmetrix input matrix)
//  r: F eigenvalues
// dsyeV eigenvalues and vectors     (real symmetric)
//  x: L columns (square symmetrix matrix)
//  r: (F eigenvalues;L eigenvectors)
// zheev eigenvalues                 (complex hermitian)
//  x: L columns (square hermitian input matrix)
//  r: F eigenvalues
// zheeV eigenvalues and vectors     (complex hermitian)
//  x: L columns (square hermitian input matrix)
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
// dgesvd singular values             (real)
//  x: L columns (input matrix)
//  r: F singular values
// dgesvD singular values and vectors (real)
//  x: L columns (input matrix)
//  r: (F singular values;L left singular vectors U;L right singular vectors V^T)
// zgesvd singular values             (complex)
//  x: L columns (input matrix)
//  r: Z singular values
// zgesvD singular values and vectors (complex)
//  x: L columns (input matrix)
//  r: (F singular values;L left singular vectors U;L right singular vectors V^T)
//
// complex numbers
//  backends that do not support complex types natively use F:
//  r0 i0 r1 i1 ..

int LAPACKE_dgesv(int, int, int, double*, int, int*, double*, int);
int LAPACKE_zgesv(int, int, int, double*, int, int*, double*, int);
int LAPACKE_dgels(int, char, int, int, int, double*, int, double*, int);
int LAPACKE_zgels(int, char, int, int, int, double*, int, double*, int);
int LAPACKE_dsyev(int, char, char, int, double*, int, double*);
int LAPACKE_zheev(int, char, char, int, double*, int, double*);
int LAPACKE_dgeev(int, char, char, int, double*, int, double*, double*, double*, int, double*, int);
int LAPACKE_zgeev(int, char, char, int, double*, int, double*,          double*, int, double*, int);
int LAPACKE_dgesvd(int, char, char, int, int, double*, int, double*, double*, int, double*, int, double*);
int LAPACKE_zgesvd(int, char, char, int, int, double*, int, double*, double*, int, double*, int, double*);

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
	FK(r, x);
	return r;
}
#ifdef KTYE
extern char *_M;
#endif
static K ZF(K x){
#ifdef KTYE
	size_t n = NK(x);
	*(int32_t*)(_M + (int32_t)x - 12) = (int32_t)(n / 2);
	x = (((K)22)<<59)|(K)(int32_t)x;
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
static K Krect(double *x, size_t rows, size_t cols, size_t z){
	K *l = (K*)malloc(cols*sizeof(K));
	for(int i=0;i<cols;i++) l[i] = KZ(x+i*z*rows, z*rows, z); 
	K r = KL(l, cols);
	free(l);
	return r;
}
static K Ksq(double *x, size_t n, size_t z, double *im){
	K *l = (K*)malloc(n*sizeof(K));
	if(im){ //compact storage: conjugate complex if im[i]==0 (see dgeev)
		double *f = (double*)malloc(16*n);
		int i = 0;
		for(i = 0;i<2*n;i++) f[i] = 0.0;
		i = 0;
		while(i < n){
			if(im[i] == (double)0.0){
				for(int j=0;j<n;j++) {
					f[  2*j] = x[j];
					f[1+2*j] = 0.0;
				}
				l[  i] = KZ(f, 2*n, 2);
			}else{
				for(int j=0;j<n;j++){
					f[  2*j] =  x[  j];
					f[1+2*j] =  x[n+j];
				}
				l[  i] = KZ(f, 2*n,   2); 
				for(int j=0;j<n;j++){
					f[1+2*j] = -x[n+j];
				}
				l[1+i] = KZ(f, 2*n,   2);
				i++;
				x += n;
			}
			i++;
			x += n;
		}
		free(f);
	}
	else   for(int i=0;i<n;i++) l[i] = KZ(x+i*z*n, z*n, z); 
	K r = KL(l, n);
	free(l);
	return r;
}

K solve(K x, K y, size_t z, int sq){ //[dz]gesv [dz]gels
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

	int multirhs = 0;
	size_t nrhs = 1;
	char yt = TK(y);
	if(yt == 'L'){ //multi-rhs
		multirhs = 1;
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
		free(A);
		if(info!=0){ free(B); return KE("lapack-dgels"); }
		if(multirhs){
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
		
		free(B);
		return r;
	} else {
		unref(y); free(A); return KE("type dgels-B?");
	}
}
K eig(K x, size_t z, char job, int sym){ // dsyev zheev [dz]geev
	size_t cols = square(&x, x, z);
	if(cols == 0){ return KE("eig A square"); }
	
	size_t xn;
	double *A = Fk(x, &xn);
	if(A == NULL){ return KE("eig: type A"); }
	
	int info;
	double *wr = (double*)malloc(cols*z*sizeof(double));
	double *wi = (double*)NULL;
	double *vl = (double*)NULL;
	double *vr = (double*)NULL;
	int n = (int)cols;
	if(sym){
		if(z > 1) info = LAPACKE_zheev(102, job, 'L', n, A, n, wr);
		else      info = LAPACKE_dsyev(102, job, 'U', n, A, n, wr);
	}else{
		if(job == 'V'){
			vl = (double*)malloc(8*z*cols*cols);
			vr = (double*)malloc(8*z*cols*cols);
		}
		if(z > 1){ 
			info = LAPACKE_zgeev(102, job, job, n, A, n, wr,     vl, n, vr, n);
		}else{
			wi = (double*)malloc(8*cols);
			info = LAPACKE_dgeev(102, job, job, n, A, n, wr, wi, vl, n, vr, n);
		}
	}
	
	K r, rw;
	if(sym){
		rw = KF(wr, cols);
		free(wr);
	}else{
		if(z > 1) rw = KZ(wr, 2*cols, 2);
		else      rw = KZ2(wr, wi, cols);
		free(wr);
	}
	if(job == 'V'){
		if(sym){
			K rl[2] = {rw, Ksq(A, cols, z, NULL)};
			r = KL(rl, 2);
		}else{
			K rl[3] = {rw, Ksq(vl, cols, z, wi), Ksq(vr, cols, z, wi)};
			free(vl);
			free(vr);
			r = KL(rl, 3);
		}
	} else  r = rw;
	if(wi) free(wi);
	free(A);
	return r;
}
K svd(K x, size_t z, char job){ // [dz]gesvd
	size_t rows, cols;
	rows = rect(&x, x, &cols);
	if(cols == 0){ return KE("gesvd: rect A"); }
	rows /= z;
	size_t mn = rows; if(cols < mn) mn = cols; // mn: min(rows, cols);
	
	size_t xn;
	double *A = Fk(x, &xn);
	if(A == NULL){ return KE("gesvd: type A"); }
	
	double *S = (double*)malloc(8*mn);
	double *U = (double*)malloc(8*z*rows*mn);
	double *V = (double*)malloc(8*z*mn*cols);
	double *s = (double*)malloc(8*mn);
	int info;
	if(z > 1) info = LAPACKE_zgesvd(102, job, job, (int)rows, (int)cols, A, (int)rows, S, U, (int)rows, V, (int)mn, s);
	else      info = LAPACKE_dgesvd(102, job, job, (int)rows, (int)cols, A, (int)rows, S, U, (int)rows, V, (int)mn, s);
	
	
	K r;
	K rs = KF(S, mn);
	if(job == 'S'){
		K l[3]={rs, Krect(U, rows, mn, z), Krect(V, mn, mn, z)};
		r = KL(l, 3);
	}else r = rs;
	free(S);
	free(U);
	free(V);
	free(s);
	free(A);
	return r;
}
K dgesv(K x, K y){ return solve(x, y, 1, 1); }
K zgesv(K x, K y){ return solve(x, y, 2, 1); }
K dgels(K x, K y){ return solve(x, y, 1, 0); }
K zgels(K x, K y){ return solve(x, y, 2, 0); }
K dsyev(K x){ return eig(x, 1, 'N', 1); }
K dsyeV(K x){ return eig(x, 1, 'V', 1); }
K zheev(K x){ return eig(x, 2, 'N', 1); }
K zheeV(K x){ return eig(x, 2, 'V', 1); }
K dgeev(K x){ return eig(x, 1, 'N', 0); }
K dgeeV(K x){ return eig(x, 1, 'V', 0); }
K zgeev(K x){ return eig(x, 2, 'N', 0); }
K zgeeV(K x){ return eig(x, 2, 'V', 0); }
K dgesvd(K x){ return svd(x, 1, 'N'); }
K dgesvD(K x){ return svd(x, 1, 'S'); }
K zgesvd(K x){ return svd(x, 2, 'N'); }
K zgesvD(K x){ return svd(x, 2, 'S'); }

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
 KR("dgesvd", (void*)dgesvd, 1);
 KR("dgesvD", (void*)dgesvD, 1);
 KR("zgesvd", (void*)zgesvd, 1);
 KR("zgesvD", (void*)zgesvD, 1);
}
*/
