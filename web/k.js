// k.js is a native javascript module that loads k.wasm and provides the api similar to +/k.h

let K={}  // exported
let _     // k wasm export object

function lo(x){return Number(BigInt.asUintN(32,x))}         // 32-bit of BigInt serves as the wasm memory index (pointer).
function us(s){return new TextEncoder("utf-8").encode(s)}
function su(u){return (u.length)?new TextDecoder("utf-8").decode(u):""}
//function ku(u){_.memory.buffer
function dxr(x,r){_.dx(x);return r}

// these should create a view into the underlying wasm linear memory.
function C(){ return new     Int8Array(_.memory.buffer) }
function U(){ return new   Uint32Array(_.memory.buffer) }
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
K.Ks = function(x){ return _.sc(K.KC(x)) }
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
 for(let i=0;i<n;i++)p[i]=lo(K.Ks(x[i]))
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

K.BK = function(x){ return dxr(x, C().slice(lo(x),lo(x)+_.nn(x))) }
K.CK = function(x){ return su(K.BK(x)) }
K.IK = function(x){ let p=lo(x)>>>2;return dxr(x,I().slice(p,p+_.nn(x))) }
K.FK = function(x){
 let t=_.tp(x); let n=(t==6) ? 2 : (t==22) ? 2*_.nn(x) : _.nn(x);
 let p=lo(x)>>>3;return dxr(x,F().slice(p,p+n))
}
K.LK = function(x){
 let n=(_.tp(x)==23) ? _.nn(x) : 2 // L vs D,T
 let r=new Array(n); let p=lo(x)>>>3; let j=J()
 for(let i=0;i<n;i++)r[i]=_.rx(j[p+i])
 return dxr(x,r)
}
K.dK = lo

K.KE   = function(e){ console.log(e); _.trap(0); return 0 }
K.Kx   = function(s,...args){ let f=_.Val(K.KC(s)); return (args.length>0) ? _.Cal(f,K.KL(args)) : f }
K.KA   = function(sym,val){ K._.dx(K._.Asn(sym,val)) }
K.ref  = function(x){return _.rx(x)}
K.unref= function(x){       _.dx(x)}


let bak   // save/restore back buffer
K.save    = function(){ bak = new Uint8Array(K._.memory.buffer).slice(0, 1<<U()[32]) }
K.restore = function(){ let u=new Uint8Array(K._.memory.buffer).set(bak) }


// low-level wasi io functions
let usr_write
let usr_read
let filename="" 
let filedata      // for fd_read

function k_write(file,nfile,src,n){
 let d = new Uint8Array(K._.memory.buffer, src,  n)
 if(nfile==0){ usr_write("",       d); return 0; }
 let filename = su(new Uint8Array(K._.memory.buffer, file, nfile))
 return usr_write(filename, d)
}
 
function k_read(file,nfile,dst){
 let filename = su(new Uint8Array(K._.memory.buffer, file, nfile))
 let u = usr_read(filename)
 if(dst == 0) return u.length;
 let m = new Uint8Array(K._.memory.buffer, dst, u.length)
 m.set(u)
 return 0
}

// js implementation for external k functions
let xcal=[]
function register(name, idx, arity){ // this is similar to the c-api's KR(..) implementation in ../+/api
 // k representation of a native function: length-2 list, the arity is stored as the vector-length.
 //   c uses: (pointer;string) where the pointer is stored in an 8-byte char vector; string is the function name (used by $f).
 //  js uses: (,index;string) with the index into xcal.
 
 let f = K._.l2(K.Kx(",", K.Ki(idx)), K.KC(name))
 I()[(lo(f)>>>2)-3] = arity                           // length-field at offset -3*32bit
 f = (BigInt(14)<<BigInt(59)) + BigInt(lo(BigInt(f))) // type tag xf(extern function) as upper bits
 K.KA(K.Ks(name), f)                                  // assign
}

 
// import object for k.wasm
let kenv={env:{ 
 Exit: function( x){ console.log("exit", x) },
 Args: function(  ){ console.log("args", x); return 0},
 Arg: function(x,y){ console.log("args", x, y); return 0},
 Read:  k_read,
 Write: k_write,
 ReadIn: function(x,y){ return 0 },
 Native: function(x,y){ let i=I()[lo(x)>>>2]
  K._.dx(x)
  return xcal[i](...K.LK(y)) },
}}

// kinit fetches, compiles and initializes l.wasm then k.wasm then calls ext.init asychronously as a callback.
K.kinit = function(ext,kw){
 let init  = ext.init;  delete ext.init  // callback when k is loaded
 usr_read  = ext.read;  delete ext.read  // file read  implementation for k: read(name)=>Uint8Array
 usr_write = ext.write; delete ext.write // file write implementation for k: write(name,data_uint8array)
 
 function binsize(x){K.n=x.byteLength;return x}
 fetch(kw?kw:'k.wasm').then(r=>r.arrayBuffer()).then(r=>WebAssembly.instantiate(binsize(r),kenv)).then(r=>{
  _=r.instance.exports
  _.kinit()
  K._=_
  
  // todo
  //  set WIDTH HEIGHT FH FW?
  
  let keys = Object.keys(ext)
  for(let i=0;i<keys.length;i++){ let jsfn=ext[keys[i]]; xcal[i]=jsfn; register(keys[i],i,jsfn.length) }
  
  if(init !== undefined)init()
  K.save()
 })
 
}

export { K }
window.K = K // for browser console
