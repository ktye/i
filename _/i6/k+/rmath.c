/*
dnorm(x, mean = 0, sd = 1, log = FALSE)
pnorm(q, mean = 0, sd = 1, lower.tail = TRUE, log.p = FALSE)
qnorm(p, mean = 0, sd = 1, lower.tail = TRUE, log.p = FALSE)
rnorm(n, mean = 0, sd = 1)
*/

//c/dat/r/R-4.5.1/src/include/Rmath.h
double dnorm4(double,double,double,int);
double pnorm(double,double,double,int,int);
double qnorm(double,double,double,int,int);
double rnorm(double,double);
//unif
//gamma
//beta
//lnorm
//chisq(2double)
//nchisq
//f (F-distribution)
//t(2double) (student t dist)
//cauchy
//exp
//..

uint64_t dnrm(uint64_t x){x=fF(Fst(x));
 if(tp(x)<16){dx(x);return Kf(dnorm4(F64((int32_t)x),0,1,0));};
 x=use(x);int32_t e=ep(x);
 for(int32_t p=(int32_t)x;p<e;p+=8)SetF64(p,dnorm4(F64(p),0,1,0));
 return x;
}
void rmath(void){
 reg(dnrm,"dnorm",1);
}
__attribute((section("rek")))void*rrmath=rmath;
