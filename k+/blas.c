typedef struct{double r,i;}cmplex;
double cblas_dnrm2(int,double*,int);
double cblas_dasum(int,double*,int);
int cblas_idamax(int,double*,int);
int cblas_izamax(int,double*,int);
void cblas_drot(int,double*,int,double*,int,double,double);
void cblas_zdrot(int,void*,int,void*,int,double,double);
double cblas_ddot(int,double*,int,double*,int);
cmplex cblas_zdotu(int,double*,int,double*,int);
cmplex cblas_zdotc(int,double*,int,double*,int);
void cblas_dgemv(int,int,int,int,double, double*,int,double*,int,double, double*,int);
void cblas_zgemv(int,int,int,int,double*,double*,int,double*,int,double*,double*,int);
#define major(x) (101+(x<0))
#define trans(x) (111+x)

void dim(uint64_t A,int32_t rA,int32_t*m,int32_t*n){rA<0?(*m=-rA,*n=-rA/nn(A)):(*m=rA,*n=rA/nn(A));}

uint64_t nrm2(uint64_t x){x=FZ(Fst(x));int32_t n=nn(x);if(tp(x)==Zt)n*=2;return Kf(cblas_dnrm2(n,FK(x),1));}
uint64_t asum(uint64_t x){x=FZ(Fst(x));int32_t n=nn(x);if(tp(x)==Zt)n*=2;return Kf(cblas_dasum(n,FK(x),1));}
uint64_t imax(uint64_t x){x=FZ(Fst(x));int32_t n=nn(x);int32_t i=tp(x)==Ft?cblas_idamax(n,FK(x),1):cblas_izamax(n,FK(x),1);return Ki(i);}
uint64_t rot(uint64_t y){double c=f0(y),s=f1(y);uint64_t x=FZ(x2(y)),q=x3(y);dx(y);x=use(x);y=use(q);
 tp(x)==Zt?cblas_zdrot(nn(x),FK(x),1,FK(y),1,c,s):cblas_drot(nn(x),FK(x),1,FK(y),1,c,s);return l2(x,y);}
uint64_t dot(uint64_t L){int32_t n;
 uint64_t x=FZ(rx(x0(L))),y=FZ(ati(L,1));
 if(tp(x)!=tp(y)||(n=nn(x))!=nn(y))trap();dx(x);dx(y);
 if(tp(x)==Zt){cmplex c=cblas_zdotu(n,FK(x),1,FK(y),1);return Kz(c.r,c.i);}
 else         return Kf(cblas_ddot (n,FK(x),1,FK(y),1));
}
uint64_t dotc(uint64_t L){int32_t n;
 uint64_t x=rx(x0(L)),y=ati(L,1);
 if(tp(x)!=Zt||tp(y)!=Zt||(n=nn(x))!=nn(y))trap();dx(x);dx(y);
 cmplex c=cblas_zdotc(n,FK(x),1,FK(y),1);return Kz(c.r,c.i);
}
uint64_t gemv(uint64_t L){ //y:gemv[a;A;rA;op;b;y;x]
 uint64_t a=rx(x0(L)),A=FZ(rx(x1(L))),b=rx(x4(L)),y=rx(x5(L)),x=FZ(ati(rx(L),6));
 int32_t rA=iL(L,2),op=iL(L,3);dx(A);
 int32_t m=rA<0?-rA:rA;
 int z=tp(A)==Zt;int32_t t=ft+z,T=t+16;
 if(  tp(a)!=t||tp(A)!=T||tp(x)!=T)trap();
 if(b!=Ki(0)){if(tp(b)!=t||tp(y)!=T)trap();y=use(y);}
 else y=mk(T,m);
 int n=nn(A)/m;
 if(nn(A)!=n*m||m!=nn(y))trap();
 if(z)cblas_zgemv(major(rA),trans(op),m,n, FK(a),FK(A),rA<0?m:n,FK(x),1, FK(b),FK(y),1);
 else cblas_dgemv(major(rA),trans(op),m,n,*FK(a),FK(A),rA<0?m:n,FK(x),1,*FK(b),FK(y),1);
 dx(A);dx(b);dx(x);return y;
}

uint64_t band(uint64_t x){ //symmetric layout only (1 1;2 2 2;3 3 3 3;4 4 4;5 5)
 int32_t lu=iL(x,0);x=ati(x,1);if(tp(x)!=Lt)trap();
 int32_t m=nn(x),u=m/2;if(!(m&1))trap();lu*=u;
 uint64_t d=FZ(ati(rx(x),u));if(tp(x)<16)trap();
 int32_t n=nn(d),z=(tp(d)==Zt);
 uint64_t r=mk(z?Zt:Ft,n*(m+lu)); //lu: additional space for lapack LU
 char*p=M_+(int32_t)r;
 const int32_t s=z?4:3;
 memset(FK(r),0,(z?16:8)*n*(m+lu));
 int32_t o=lu;for(int32_t i=0;i<m;i++){
  uint64_t xi=FZ(ati(rx(x),i));
  int32_t nx=nn(xi);
  if(nx!=(i<u?n-u+i:n+u-i)||(z&&tp(xi)!=Zt))trap();
  int32_t k=nn(ati(rx(x),u));
  if(i>u)o++;memcpy(p+(o+(n*(m+lu-1-i))<<s),FK(xi),nx<<s);dx(xi);
 }
 dx(x);return r;
}

void blas(void){
 reg(band,"band",2); //band matrix from diagonals, e.g. band(3 3 3;1 1 1 1;2 2 2) netlib.org/lapack/lug/node124.html
 reg(nrm2,"nrm2",1);
 reg(asum,"asum",1);
 reg(imax,"imax",1);
 reg(rot,  "rot",4);
 reg(dot,  "dot",2);
 reg(dotc,"dotc",2);
 reg(gemv,"gemv",7);
}
__attribute((section("rek")))void*rblas=blas;
