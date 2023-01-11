data $data0 = { b "\x61\x62\x63\x64\x01\x02\x03" }
function init(){
@start

Memory(1);

%t =l call memcpy(%_M, $data0, 7)


}
