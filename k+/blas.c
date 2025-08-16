double cblas_dnrm2(int,double*,int);
double cblas_dasum(int,double*,int);
int cblas_idamax(int,double*,int);
int cblas_izamax(int,double*,int);
void cblas_drot(int,double*,int,double*,int,double,double);
void cblas_zdrot(int,void*,int,void*,int,double,double);
void cblas_dgemv(int,int,int,int,double,double*,int,double*,int,double,double*,int);
#define major(x) (101+(x<0))
#define trans(x) (111+x)

void dim(uint64_t A,int32_t rA,int32_t*m,int32_t*n){rA<0?(*m=-rA,*n=-rA/nn(A)):(*m=rA,*n=rA/nn(A));}

uint64_t nrm2(uint64_t x){x=FZ(Fst(x));int32_t n=nn(x);if(tp(x)==Zt)n*=2;return Kf(cblas_dnrm2(n,FK(x),1));}
uint64_t asum(uint64_t x){x=FZ(Fst(x));int32_t n=nn(x);if(tp(x)==Zt)n*=2;return Kf(cblas_dasum(n,FK(x),1));}
uint64_t imax(uint64_t x){x=FZ(Fst(x));int32_t n=nn(x);int32_t i=tp(x)==Ft?cblas_idamax(n,FK(x),1):cblas_izamax(n,FK(x),1);return Ki(i);}
uint64_t rot(uint64_t y){double c=f0(y),s=f1(y);uint64_t x=FZ(x2(y)),q=x3(y);dx(y);x=use(x);y=use(q);
 tp(x)==Zt?cblas_zdrot(nn(x),FK(x),1,FK(y),1,c,s):cblas_drot(nn(x),FK(x),1,FK(y),1,c,s);return l2(x,y);}
uint64_t gemv(uint64_t u){int32_t m,n;
//       0 1 2  3  4 5 6
//y:gemm[a;A;rA;op;b;y;x]
 uint64_t A=FZ(x1(u)),y=FZ(x5(u)),x=FZ(x6(u));
 if(tp(A)==Ft){
  double a=f0(y),b=f4(y);int32_t rA=i2(y),op=i3(y);dx(u);y=use(y);dim(A,rA,&m,&n);
  cblas_dgemv(major(rA),trans(op),m,n,a,FK(A),n,FK(x),1,b,FK(y),1);
  dx(A);dx(x);return b;
 }else{
  printf("todo mv complex\n");
 }
}

uint64_t band(uint64_t x){ //symmetric layout only (1 1;2 2 2;3 3 3 3;4 4 4;5 5)
 x=Fst(x);if(tp(x)!=Lt)trap();
 int32_t m=nn(x),u=m/2;if(!(m&1))trap();
 uint64_t d=FZ(ati(rx(x),u));if(tp(x)<16)trap();
 int32_t n=nn(d),z=(tp(d)==Zt);
 uint64_t r=mk(z?Zt:Ft,n*(m+u)); //additional space for lapack LU
 char*p=M_+(int32_t)r;
 const int32_t s=z?4:3;
 memset(FK(r),0,(z?16:8)*n*(m+u));
 int32_t o=0;for(int32_t i=0;i<m;i++){
  uint64_t xi=FZ(ati(rx(x),i));
  int32_t nx=nn(xi);
  if(nx!=(i<u?n-u+i:n+u-i)||(z&&tp(xi)!=Zt))trap();
  int32_t k=nn(ati(rx(x),u));
  if(i>u)o++;memcpy(p+(o+(n*(m-1-i))<<s),FK(xi),nx<<s);dx(xi);
 }
 dx(x);return r;
}

void blas(void){
 reg(band,"band",1); //band matrix from diagonals, e.g. band(3 3 3;1 1 1 1;2 2 2) netlib.org/lapack/lug/node124.html
 reg(nrm2,"nrm2",1);
 reg(asum,"asum",1);
 reg(imax,"imax",1);
 reg(rot,  "rot",4);
 reg(gemv,"gemv",7);
}
__attribute((section("rek")))void*rblas=blas;
