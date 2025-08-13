#define NX 64 //large enough to hold all external functions

typedef uint64_t (*XF)(uint64_t);
uint64_t IX=0;XF xtab[NX];
void reg(XF f, const char*name, int32_t arity){
 uint64_t c=mk(Ct,(int32_t)strlen(name));
 if(!help){help=sc(Ku(1886152040ull));hlp=Key(mk(St,0),mk(It,0));}
 
 for(int i=0;i<strlen(name);i++)M_[i+(int32_t)c]=name[i];
 uint64_t s=sc(rx(c));

 hlp=Cat(hlp,Key(Enl(s),Enl(Ki(arity))));

 //a native function has type 14 and stores 2 values:
 // -any identifier, we use a primary function type 0 (no refcounts)
 // -a strings with the name used to display the function object
 //the identifier & the arglist is passed to native() when the function is called.
 int64_t l=ti(14,(uint32_t)l2(IX, c));
 SetI32((int32_t)l-12,arity);
 dx(Asn(s,l)); //assign the function object to a global symbol, for the user to access it.
 if(NX==IX){fprintf(stderr,"xtab is too small: increase NX in l.c\n");exit(1);}
 xtab[IX++]=f;
}
int64_t cnative(int64_t x,int64_t y){return(int64_t)xtab[(int32_t)x](y);}



//some helpers to unpack arguments
static uint64_t FZ(int64_t x){ //real or complex vector arg
 int32_t xt=tp(x);
 if(xt<16)trap();
 if(xt<Ft){x=Add(Kf(0.0),x);}
 if(xt>Zt)trap();
 return x;
}
static uint64_t fF(int64_t x){int32_t t=tp(x); //f or F
 if(t==ft||t==Ft)return x;
 if(t<ft||(t>16&&t<Ft))return Add(Kf(0.0),x);
 else trap();return 0;
}
static uint64_t x3(uint64_t x){return x0((x+24ull));}
static uint64_t x5(uint64_t x){return x0((x+40ull));}
static uint64_t x6(uint64_t x){return x0((x+48ull));}

int   *IK(uint64_t x){return (int   *)(M_+(int32_t)x);       }
double*FK(uint64_t x){return (double*)(M_+(int32_t)x);       }
double f0(uint64_t x){return F64((int32_t)I64((int32_t)x  ));}
double f1(uint64_t x){return F64((int32_t)I64((int32_t)x+8));}
double f4(uint64_t x){return F64((int32_t)I64((int32_t)x+32));}
int32_t i2(uint64_t x){return I32((int32_t)I64((int32_t)x+16));}
int32_t i3(uint64_t x){return I32((int32_t)I64((int32_t)x+24));}

