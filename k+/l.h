#define NATIVE 1
void libs(void);
int64_t cnative(int64_t,int64_t);
static uint64_t se(double);
void sincos(double,double*,double*);
static void sin1(int32_t,int32_t,int32_t);
static void cos1(int32_t,int32_t,int32_t);
#define cosin_(x,r,_) sincos(x,(double*)(M_+8+r),(double*)(M_+r));
