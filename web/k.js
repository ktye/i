// k.js is a native javascript module that loads k.wasm and provides the api similar to +/k.h

let K={}  // exported
let _     // k wasm export object

function lo(x){return Number(BigInt.asUintN(32,x))}         // 32-bit of BigInt serves as the wasm memory index (pointer).
function us(s){return new TextEncoder("utf-8").encode(s)}
function su(u){return (u.length)?new TextDecoder("utf-8").decode(u):""}
//function ku(u){_.memory.buffer
function dxr(x,r){_.dx(x);return r}

function C(){ return new     Int8Array(_.memory.buffer) }
function I(){ return new    Int32Array(_.memory.buffer) }
function J(){ return new BigInt64Array(_.memory.buffer) }
function F(){ return new  Float64Array(_.memory.buffer) }

// type/length
K.TK = function(x){ 
 const t="-icisfF----------ICISFFLDT"; //like ../+/api
 return t[_.tp(x)]
}
K.NK = function(x){ return _.nn(x) }

// create k atoms
K.Kc = function(x){ return _.Kc( ("string"===typeof(x)) ? x.charCodeAt(0) : x ) }
K.Ks = function(x){ return _.sc(K.Kc(x)) }
K.Ki = function(x){ return _.Ki(x) }
K.Kf = function(x){ return _.Kf(x) }

// create k vectors
K.KC = function(x){
 if("string"===typeof(x))x=us(x);
 let r=_.mk(18,x.length);C().set(x,lo(r));return r
}
K.KS = function(x){
 let n=x.length
 var r=_.mk(20,n)
 var p=I().slice(lo(r))
 for(let i=0;i<n;i++)p[i]=lo(K.s(x[i]))
 return r
}
K.KI = function(x){
 x = (x.constructor===Int32Array) ? x : new Int32Array(x)
 let r=_.mk(19,x.length)
 I().set(x,lo(r)>>>2)
 return r
}
K.KF = function(x){
 x = (x.constructor===Float64Array) ? x : new Int32Array(x)
 let r=_.mk(21,x.length)
 F().set(x,lo(r)>>>3)
 return r
}
K.KL = function(x){
 let n=x.length
 console.log("KL", n)
 let r=_.mk(23,n)
 let j=J()
 let p=lo(r)>>>3
 for(let i=0;i<n;i++)j[i+p]=x[i]
 return r
}

K.cK = function(x){ return lo(x) << 24 >> 24 } // signed int8
K.iK = function(x){ return lo(x) << 0 } // signed int32
K.fK = function(x){ 
 let p=lo(x)>>>3;return dxr(x,F()[p])
}

K.CK = function(x){ return dxr(x, su(C().slice(lo(x),lo(x)+_.nn(x)))) }
K.IK = function(x){ let p=lo(x)>>>2;return dxr(x,I().slice(p,p+_.nn(x))) }
K.FK = function(x){
 let t=_.tp(x); let n=(t==6) ? 2 : (t==22) ? 2*_.nn(x) : _.nn(x);
 let p=lo(x)>>>3;return dxr(x,F().slice(p,p+n))
}
K.LK = function(x){
 let n=(_.tp(x)==23) ? _.nn(x) : 2 // L vs D,T
 let r=new Array(n); let p=lo(x)>>>3; let j=J()
 for(let i=0;i<n;i++)r[i]=_.rx(J[p+i])
 return dxr(x,r)
}
K.dK = lo

K.Kx   = function(s,...args){ let f=_.Val(K.KC(s)); return (args.length>0) ? _.Cal(f,K.KL(args)) : f }
K.ref  = function(x){return _.rx(x)}
K.unref= function(x){       _.dx(x)}

function xx(x){console.log(x);return x}

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

K.kinit = function(f,r,w){ // f:onsuccess, r:read, w:write
 let n=0; function binsize(x){n=x.byteLength;return x}
 fetch('k.wasm').then(r=>r.arrayBuffer()).then(r=>WebAssembly.instantiate(binsize(r),env)).then(r=>{
  _=r.instance.exports
  _.kinit()
  _.reset()
  // todo
  // set WIDTH HEIGHT FH FW?
  // ksave()
  K._=_
  f()
 })
 
}

export { K }
window.K = K // for browser console
