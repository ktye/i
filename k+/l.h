#define NATIVE 1
void libs(void);
int64_t cnative(int64_t,int64_t);
static uint64_t se(double);
void sincos(double,double*,double*);
static void sin1(int32_t,int32_t,int32_t);
static void cos1(int32_t,int32_t,int32_t);
#define cosin_(x,r,_) sincos(x,(double*)(M_+8+r),(double*)(M_+r));

//we collect names and arity of extention functions in a dict assigned to "help"
uint64_t hlp=0,help=0;

//initialization magick, thanks @ovs.
extern void(*__start_reg)(),(*__stop_reg)();
