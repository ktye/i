//reinsert lost functions implemented with libm

/*
static uint64_t se(double x){
 int32_t ei;
 double f;
 int64_t e;
 f=x;
 e=(int64_t)(0ll);
 if(frexp1(x)){
  f=frexp2(x);
  e=frexp3(x);
 }
 x=(0.3010299956639812*(double)(e));
 ei=(int32_t)(F64floor(x));
 x=(x-(double)(ei));
 return ucat(cat1(sf((f*pow_(10.,x))),Kc(101)),si(ei));
}
static void sin1(int32_t xp, int32_t yp, int32_t rp){
 cosin_(F64(xp),rp,1);
}
static void cos1(int32_t xp, int32_t yp, int32_t rp){
 cosin_(F64(xp),rp,2);
}
*/
static uint64_t se(double x){
 int32_t ei;
 double f;
 int e;
 f=frexp(x,&e);
 x=(0.3010299956639812*(double)(e));
 ei=(int32_t)(F64floor(x));
 x=(x-(double)(ei));
 return ucat(cat1(sf((f*pow(10,x))),Kc(101)),si(ei));
}
static void sin1(int32_t xp, int32_t yp, int32_t rp){
 SetF64(rp,sin(F64(xp)));
}
static void cos1(int32_t xp, int32_t yp, int32_t rp){
 SetF64(rp,cos(F64(xp)));
}

