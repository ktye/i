void*curl_easy_init(void);
int  curl_easy_perform(void*);
void curl_easy_setopt(void*,int,void*);
void curl_easy_cleanup(void*);

static size_t curlback(void *c,size_t n,size_t m,uint64_t*r){
 uint64_t x=*r;int32_t xn=nn(x);
 n*=m;x=uspc(x,Ct,n);*r=x;memcpy(M_+(int32_t)x+xn,c,n);return n;}

uint64_t kurl(uint64_t x){
 x=Fst(x);if(tp(x)!=Ct)trap();x=cat1(x,Kc(0));
 void*c=curl_easy_init();
 curl_easy_setopt(c,10002,M_+(int32_t)x);
 uint64_t r=mk(Ct,0);
 curl_easy_setopt(c,20011,curlback);
 curl_easy_setopt(c,10001,(void*)&r);
 curl_easy_setopt(c,10018,"libcurl-agent/1.0");
 if(curl_easy_perform(c))trap();
 curl_easy_cleanup(c);
 dx(x);return r;}

void curl(void){reg(kurl,"curl",1);}
__attribute((section("reg")))void*rcurl=curl;
