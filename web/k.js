// k.js is a native javascript module that loads k.wasm and provides the api similar to +/k.h


function reset(){
 
}

var filename="" 
function path_open(fd,dirflags,path,pathlen,oflags,baserights,inheritrights,fdflags,res){
// msl();filename=su(C.slice(path,path+pathlen));if(oflags==0)return 1;U[res>>2]=3;return 0 //no read in js.
} 
function fd_write(fd,p,nio,nw){
// msl();var x=U[p>>>2]; var n=U[(4+p)>>>2];var u=(C.slice(x,x+n));if(fd==1){O(su(u))}else{download(filename,u)};return 0
}

var env={wasi_unstable:{ 
// l32:function(x){console.log(x);return x},
// l64:function(x){console.log(BigInt.asUintN(64,x));return x},
 args_get: function(a,b){return 0},
 args_sizes_get: function(a,b){return 0},
 proc_exit:function(x){console.log("exit", x)},
 fd_read:  function(a,b,c,d){console.log("read",a,b,c,d);return 0},
 fd_write: fd_write,
 fd_seek:  function(fd,o,w,r){return 0},
 fd_close: function(fd){return 0},
 path_open:path_open,
 clock_time_get:function(a,b,p){msl();J[p>>>3]=BigInt.asIntN(64, BigInt(Math.floor(1e6*performance.now())));return 0}
}}

// kinit loads k.wasm then executes f passing k.wasm binary size as an argument.
//  w(name_string, data_uint8array) is called when k writes to a file, w("", data) is stdout, e.g. terminal. no return value.
//  r(name_string) returns data_uint8array when k reads from a file. there is no stdin in k.js/wasm.
function kinit(f,r,w){ 
 let n=0; function binsize(x){n=x.byteLength;return x}
 fetch('k.wasm').then(r=>r.arrayBuffer()).then(r=>WebAssembly.instantiate(binsize(r),env)).then(r=>{
  let K=r.instance.exports
  K.kinit();

  K.reset()
  // set WIDTH HEIGHT FH FW?
  // ksave()
  

 
  f(n);
 })
 
}

export { kinit }
