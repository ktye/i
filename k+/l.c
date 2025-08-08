void reg(int64_t id, int64_t name, int32_t arity){
 int64_t l=ti(14,(uint32_t)l2(id, Ku(name)));
 SetI32((int32_t)l-12,arity);
 dx(Asn(sc(Ku(name)),l));
}
void libs(void){ //encode strings with: https://play.golang.org/p/4ethx6OEVCR
 reg(0,846033518ull,1); //nrm2
 //reg(1,...);
 //reg(2,...);
}
int64_t cnative(int64_t x, int64_t y){
 switch(x){ //switch registered function id
 case 0:  printf("todo call blas nrm2\n"); return y; break;
 default: return y;
}}

