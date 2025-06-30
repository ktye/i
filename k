#!/usr/bin/env -S deno --allow-read --allow-write

let argv=process.argv.slice(2),readfile=x=>Deno.readFileSync(x)
let su=x=>t_.decode(x),t_=new TextDecoder("utf-8"),us=x=>_t.encode(x),_t=new TextEncoder("utf-8")
let U=(o,n)=>new Uint8Array(K.memory.buffer,o,n),lo=x=>Number(BigInt.asUintN(32,x))

let nyi=_=>{throw new Error("nyi")}

let I={"env":{
 Exit:Deno.exit,
 Args:_=>argv.length,
 Arg:(x,y)=>y?U().set(us(argv[x]),r):argv[x].length,
 Read:(x,y,z)=>nyi(),
 Write:(f,nf,s,n)=>(nf?nyi():Deno.stdout.writeSync(U(s,n))),
 //ReadIn:(x,y)=>(x=U(x,y),y=os.stdin.readSync(x),y-=(x[n-1]==10)?1:0,y-=(x[n-1]==13)?1:0,y),
 ReadIn:(x,y)=>nyi(),
 Native:(x,y)=>nyi(),
}}

let kw=readfile((Deno.build.os=="windows"?"c:/":"/c/")+"k/ktye.github.io/k.wasm")
let K=new WebAssembly.Instance(new WebAssembly.Module(kw),I).exports
K.kinit()

console.log("ktye/k(deno)")
let _k=1,repl=_=>{
 let buf=new Uint8Array(256),p=new Uint8Array([32]);
 while(1){Deno.stdout.writeSync(p);let n=Deno.stdin.readSync(buf);if(!n)return;
  if(n&&buf[n-1]==10)--n;if(n&&buf[n-1]==13)--n;
  let x=buf.subarray(0,n);
  try{console.log("x",x); let r=K.mk(18,x.length);U().set(x,lo(r));K.repl(r) }
  catch(e){console.log(e)}}}

repl()

