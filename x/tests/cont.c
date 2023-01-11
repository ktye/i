void f(){
 int32_t i;
 i=(int32_t)(0);
 for(;;){
  i++;
  if(i<2){
   continue;
  }
  i=(i*2);
 }
}
