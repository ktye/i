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

// A:(5;0 2 5 9 10 12;0 1 0 2 4 1 2 3 4 1 2 1 4;2. 3. 3. -1. 4. 4. -3. 1. 2. 2. 6. 1.)
// b:8. 45. -3. 3. 19.
uint64_t spsolve(uint64_t x){ // spsolve[A;b]
 uint64_t A=x0(A),B=x1(A);dx(A);
 if(tp(A)!=Lt)trap();if(nn(A)!=4)trap();
 uint64_t N=x0(A),Ap=x1(A),Ai=x2(A),Ax=FZ(x3(A));dx(A);
 if(tp(N)!=it||tp(Ap)!=It||tp(Ai)!=It)trap();
 int z=0,n=(int32_t)N;int*ap=IK(Ap),*ai=IK(Ai);double*ax=FK(Ax);void*sy,*nu;
 if(tp(Ax)==Ft){
  umfpack_di_symbolic(n,n,ap,ai,ax,&sy,NULL,NULL);
  umfpack_di_numeric(ap,ai,ax,sy,&nu, NULL,NULL);
  umfpack_di_free_symbolic(&sy);
 }else{z=1;
  umfpack_zi_symbolic(n,n,ap,ai,ax,NULL,&sy,NULL,NULL);
  umfpack_zi_numeric(ap,ai,ax,NULL,sy,&nu,NULL,NULL);
  umfpack_zi_free_symbolic(&sy);
 }
 if(tp(B)<Lt)B=l1(B);
 int nb=nn(B),s=0;
 uint64_t r=mk(Lt,nb),ri;
 for(int i=0;i<nb;i++){
  uint64_t bi=(uint64_t)I64(8*i+(int32_t)B);
  if(!z){
   if(tp(bi)!=Ft){umfpack_di_free_numeric(&nu);trap();}
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
 return(tp(B)<Lt)?Fst(r):r;
}

void umfpack(void){
 reg(spsolve,"spsolve",2);
}
__attribute((section("reg")))void*rumpfpack=umfpack;

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

