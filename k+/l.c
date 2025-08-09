enum fn {
#ifdef BLAS
  nrm2, asum, imax,
#endif
};
void reg(int64_t id, int64_t name, int32_t arity){
 int64_t l=ti(14,(uint32_t)l2(id, Ku(name)));
 SetI32((int32_t)l-12,arity);
 dx(Asn(sc(Ku(name)),l));
}
void libs(void){ //encode strings with: https://play.golang.org/p/4ethx6OEVCR
#ifdef BLAS
 reg(nrm2, 846033518ull,1); //id k-name arity (name is same as id-enum)
 reg(asum,1836413793ull,1);
 reg(imax,2019650921ull,1);
#endif
}

int64_t fzvec(int64_t x){
 int32_t xt=tp(x);
 if(xt<16){x=Enl(x);xt=tp(x);}
 if(xt<Ft){x=Add(Kf(0.0),x);}
 if(xt>Zt)trap();
 return x;
}

#ifdef BLAS
double cblas_dnrm2(int,double*,int);
double cblas_dasum(int,double*,int);
int cblas_idamax(int,double*,int);
int cblas_izamax(int,double*,int);
#endif

#define fp(x) ((double*)(M_+(int32_t)x))

int64_t cnative(int64_t x, int64_t y){ //function-id, list of args(length arity)
 int32_t xt,xn,i;
 double r;
 switch(x){ //switch registered function id
#ifdef BLAS
 case nrm2:x=fzvec(Fst(y));xn=nn(x);if(xt==Zt)xn*=2;
  r=cblas_dnrm2(xn,fp(x),1);
  dx(x);return Kf(r);
 case asum:x=fzvec(Fst(y));xn=nn(x);if(xt==Zt)xn*=2;
  r=cblas_dasum(xn,fp(x),1);
  dx(x);return Kf(r);
 case imax:x=fzvec(Fst(y));xn=nn(x);
  i=tp(x)==Ft?cblas_idamax(xn,fp(x),1):cblas_izamax(xn,fp(x),1);
  dx(x);return Ki(i);
#endif
 default: return y;
}}

