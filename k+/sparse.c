//umfpack solve sparse Ax=b
int  umfpack_di_symbolic(int,int,int*,int*,double*,        void**,double*,double*);
int  umfpack_zi_symbolic(int,int,int*,int*,double*,double*,void**,double*,double*);
int  umfpack_di_numeric(         int*,int*,double*,        void*,void**,double*,double*);
int  umfpack_zi_numeric(         int*,int*,double*,double*,void*,void**,double*,double*);
void umfpack_di_free_symbolic(void**);
void umfpack_zi_free_symbolic(void**);
int  umfpack_di_solve(int,int*,int*,double*,        double*,        double*,        void*,double*,double*);
int  umfpack_zi_solve(int,int*,int*,double*,double*,double*,double*,double*,double*,void*,double*,double*);
void umfpack_di_free_numeric(void**);
void umfpack_zi_free_numeric(void**);

//spqr sparse least squares (cholmod.h spqr.h)
typedef struct choldmod_dense_struct {
 size_t nrow,ncol,nzmax,d; //A is nrow x ncol, nmax(entries) d(lda>=nrow)
 double *x;void *z;
 int xtype,dtype;
} cholmod_dense;
typedef struct cholmod_sparse_struct{
 size_t nrow,ncol,nzmax;
 void *p,*i,*nz,*x,*z;
 int stype,itype,xtype,dtype,sorted,packed; //0:unsymmetric, >0:square sym upper, <0:square sym lower
} cholmod_sparse;
int cholmod_start(void*cc);
int cholmod_finish(void*cc);
int cholmod_free_dense(cholmod_dense**,void*);
cholmod_dense*SuiteSparseQR_C_backslash_default(cholmod_sparse*A,cholmod_dense*B,void*cc);

uint64_t spqr(uint64_t x){ //todo: b maybe m-by-k dense..
 uint64_t common[400];//enough to hold cholmod_common structure (sizeof 2680)
 uint64_t A=rx(x0(x)),B=rx(x1(x));dx(x);
 if(tp(A)!=Lt||nn(A)!=5)trap();
 uint64_t M=x0(A),N=x1(A),Ap=x2(A),Ai=x3(A),Ax=FZ(x4(A));dx(A);int z=(tp(Ax)==Zt);
 if(tp(M)!=it||tp(N)!=it||tp(Ap)!=It||tp(Ai)!=It)trap();
 int32_t m=(int32_t)M,n=(int32_t)N;
 if(nn(Ap)!=1+n||nn(Ai)!=nn(Ax))trap();
 int32_t nb=nn(B),k=nb/m;if(nb!=k*m)trap(); //k>1: multirhs flat colmajor vector
 if(z&&tp(B)!=Zt)B=uptype(B,zt);
 void *cc=&common;cholmod_start(cc);
 cholmod_sparse a;
 a.nrow=m;a.ncol=n;a.nzmax=nn(Ax);a.p=IK(Ap);a.i=IK(Ai);a.nz=NULL;a.x=FK(Ax);a.z=NULL;
 a.stype=0;a.itype=0;a.xtype=1+z;a.dtype=0;a.sorted=0;a.packed=1;
 cholmod_dense b;
 b.nrow=m;b.ncol=k;b.nzmax=m;b.x=FK(B);b.z=NULL;b.xtype=1+z;b.dtype=0;
 cholmod_dense*r=SuiteSparseQR_C_backslash_default(&a,&b,cc);
 dx(Ap);dx(Ai);dx(Ax);dx(B);
 x=mk(z?Zt:Ft,nb);
 memcpy(M_+(int32_t)x,r->x,nb<<(z?4:3));
 cholmod_free_dense(&r,cc);
 cholmod_finish(cc);
 return x;
}

//dense matrix representation
//(i;i;I; I; F or Z)
// m n Ap Ai A

// A:(5;5;0 2 5 9 10 12;0 1 0 2 4 1 2 3 4 1 2 1 4;2. 3. 3. -1. 4. 4. -3. 1. 2. 2. 6. 1.)
// b:8. 45. -3. 3. 19.
uint64_t spsolve(uint64_t x){ // spsolve[A;b]
 uint64_t A=x0(x),B=x1(x);dx(x);
 if(tp(A)!=Lt)trap();if(nn(A)!=5)trap();
 uint64_t N=x1(A),Ap=x2(A),Ai=x3(A),Ax=FZ(x4(A));dx(A);
 if(tp(N)!=it||tp(Ap)!=It||tp(Ai)!=It)trap();
 int z=0,n=(int32_t)N;int*ap=IK(Ap),*ai=IK(Ai);double*ax=FK(Ax);void*sy,*nu;
 if(nn(Ap)!=1+n||nn(Ai)!=nn(Ax))trap();
 if(tp(Ax)==Ft){
  if(umfpack_di_symbolic(n,n,ap,ai,ax,&sy,NULL,NULL))trap();
  if(umfpack_di_numeric(ap,ai,ax,sy,&nu, NULL,NULL)){umfpack_di_free_symbolic(&sy);trap();}
  umfpack_di_free_symbolic(&sy);
 }else{z=1;
  if(umfpack_zi_symbolic(n,n,ap,ai,ax,NULL,&sy,NULL,NULL))trap();
  if(umfpack_zi_numeric(ap,ai,ax,NULL,sy,&nu,NULL,NULL)){umfpack_zi_free_symbolic(&sy);trap();}
  umfpack_zi_free_symbolic(&sy);
 }
 int mrhs=(tp(B)==Lt);if(tp(B)<Lt)B=l1(B);
 int nb=nn(B),s=0;
 uint64_t r=mk(Lt,nb),ri;
 for(int i=0;i<nb;i++){
  uint64_t bi=(uint64_t)I64(8*i+(int32_t)B);
  if(!z){
   if(tp(bi)!=Ft||nn(bi)!=n){umfpack_di_free_numeric(&nu);trap();}
   uint64_t ri=mk(Ft,n);
   s=umfpack_di_solve(0,ap,ai,ax,FK(ri),FK(bi),nu,NULL,NULL);
   if(s){umfpack_di_free_numeric(&nu);trap();}
   SetI64(8*i+(int32_t)r,ri);
  }else{
   if(tp(bi)!=Zt){umfpack_zi_free_numeric(&nu);trap();}
   uint64_t ri=mk(Zt,n);
   s=umfpack_zi_solve(0,ap,ai,ax,NULL,FK(ri),NULL,FK(bi),NULL,nu,NULL,NULL);
   if(s){umfpack_zi_free_numeric(&nu);trap();}
   SetI64(8*i+(int32_t)r,ri);
  }
 }
 if(!z)umfpack_di_free_numeric(&nu);
 else  umfpack_zi_free_numeric(&nu);
 return mrhs?r:Fst(r);
}
void umfpack(void){
 reg(spsolve,"spsolve",2);
 reg(spqr,   "spqr",   2);
}
__attribute((section("rek")))void*rumpfpack=umfpack;

/*
#include <stdio.h>
int n = 5 ;
int Ap [ ] = {0, 2, 5, 9, 10, 12} ;
int Ai [ ] = { 0, 1, 0, 2, 4, 1, 2, 3, 4, 2, 1, 4} ;
double Ax [ ] = {2., 3., 3., -1., 4., 4., -3., 1., 2., 2., 6., 1.} ;
double b [ ] = {8., 45., -3., 3., 19.} ;
double x [5] ;
int main (void)
{
double *null = (double *) NULL ;
int i ;
void *Symbolic, *Numeric ;
 (void) umfpack_di_symbolic (n, n, Ap, Ai, Ax, &Symbolic, null, null) ;
 (void) umfpack_di_numeric (Ap, Ai, Ax, Symbolic, &Numeric, null, null) ;
 umfpack_di_free_symbolic (&Symbolic) ;
 (void) umfpack_di_solve (0, Ap, Ai, Ax, x, b, Numeric, null, null) ;
 umfpack_di_free_numeric (&Numeric) ;
 for (i = 0 ; i < n ; i++) printf ("x [%d] = %g\n", i, x [i]) ;
 return (0) ;
}
*/

/*spqr:  
/c/dat/SuiteSparse-7.11.0/SPQR/Include/SuiteSparseQR_C.h
/c/dat/SuiteSparse-7.11.0/CHOLMOD/Include

#include<stdio.h>
#include<cholmod.h>
int main(){
 printf("%d\n", sizeof(cholmod_common)); //should be 2680
}
*/
