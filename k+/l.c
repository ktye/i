enum fn {
#ifdef BLAS
  nrm2, asum, imax, rot, mv,
#endif
#ifdef CURL
  curl,
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
 reg(rot,    7630706ull,4);
 reg(mv,       30317ull,7);
#endif
#ifdef CURL
 reg(curl,1819440483ull,1);
#endif
}

static uint64_t FZ(int64_t x){
 int32_t xt=tp(x);
 if(xt<16)trap();
 if(xt<Ft){x=Add(Kf(0.0),x);}
 if(xt>Zt)trap();
 return x;
}
static uint64_t x3(uint64_t x){return x0((x+24ull));}
static uint64_t x5(uint64_t x){return x0((x+40ull));}
static uint64_t x6(uint64_t x){return x0((x+48ull));}

#ifdef BLAS
double cblas_dnrm2(int,double*,int);
double cblas_dasum(int,double*,int);
int cblas_idamax(int,double*,int);
int cblas_izamax(int,double*,int);
void cblas_drot(int,double*,int,double*,int,double,double);
void cblas_zdrot(int,void*,int,void*,int,double,double);
void cblas_dgemv(int,int,int,int,double,double*,int,double*,int,double,double*,int);
#define major(x) (101+(x<0))
#define trans(x) (111+x)
#endif

#ifdef CURL
void*curl_easy_init(void);
int  curl_easy_perform(void*);
void curl_easy_setopt(void*,int,void*);
void curl_easy_cleanup(void*);
static size_t curlback(void *c,size_t n,size_t m,uint64_t*r){
 uint64_t x=*r;int32_t xn=nn(x);n*=m;x=uspc(x,Ct,n);*r=x;
 memcpy(M_+(int32_t)x+xn,c,n);return n;}
#endif

double*FK(uint64_t x){return (double*)(M_+(int32_t)x);       }
double f0(uint64_t x){return F64((int32_t)I64((int32_t)x  ));}
double f1(uint64_t x){return F64((int32_t)I64((int32_t)x+8));}
double f4(uint64_t x){return F64((int32_t)I64((int32_t)x+32));}
int32_t i2(uint64_t x){return I32((int32_t)I64((int32_t)x+16));}
int32_t i3(uint64_t x){return I32((int32_t)I64((int32_t)x+24));}
void dim(uint64_t A,int32_t rA,int32_t*m,int32_t*n){rA<0?(*m=-rA,*n=-rA/nn(A)):(*m=rA,*n=rA/nn(A));}


//x:function-id(not refcounted), y:list of args(length arity)
int64_t cnative(int64_t x, int64_t y){
 int32_t m,n,i;
 uint64_t q;
 double r;
 switch(x){ //switch registered function id
#ifdef BLAS
 case nrm2:x=FZ(Fst(y));n=nn(x);if(tp(x)==Zt)n*=2;
  r=cblas_dnrm2(n,FK(x),1);
  return Kf(r);
 case asum:x=FZ(Fst(y));n=nn(x);if(tp(x)==Zt)n*=2;
  r=cblas_dasum(n,FK(x),1);
  return Kf(r);
 case imax:x=FZ(Fst(y));n=nn(x);
  i=tp(x)==Ft?cblas_idamax(n,FK(x),1):cblas_izamax(n,FK(x),1);
  return Ki(i);
 case rot:{double c=f0(x),s=f1(x);x=FZ(x2(y));q=x3(y);dx(y);x=use(x);y=use(q);
  tp(x)==Zt?cblas_zdrot(nn(x),FK(x),1,FK(y),1,c,s):cblas_drot(nn(x),FK(x),1,FK(y),1,c,s);  
  return l2(x,y);}
 case mv:{ //todo..
  if(tp(x)==Zt){
   double a=f0(y),b=f4(y);uint64_t A=FZ(x1(y));int32_t rA=i2(y),op=i3(y);uint64_t r=FZ(x6(y));x=FZ(x5(y));dx(y);r=use(r);dim(A,rA,&m,&n);
   cblas_dgemv(major(rA),trans(op),m,n,a,FK(A),n,FK(x),1,b,FK(r),1);
   dx(A);dx(x);return b;
  }else{
   printf("todo mv complex\n");
  }
 };
#endif
#ifdef CURL
 case curl:{x=Fst(y);if(tp(x)!=Ct)trap();x=cat1(x,Kc(0));
  void*c=curl_easy_init();
  curl_easy_setopt(c,10002,M_+(int32_t)x);
  uint64_t r=mk(Ct,0);
  curl_easy_setopt(c,20011,curlback);
  curl_easy_setopt(c,10001,(void*)&r);
  curl_easy_setopt(c,10018,"libcurl-agent/1.0");
  if(curl_easy_perform(c))trap();
  curl_easy_cleanup(c);
  dx(x);return r;}
#endif
 default: return y;
}}


