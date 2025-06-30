#!/usr/bin/env -S deno --allow-read --allow-env
//--allow-write

//k [-m] [file.k] [datafile] [-e]
//   -m: fft.k lu.k qr.k stats.k svd.k

let argv=process.argv.slice(2),readfile=x=>Deno.readFileSync(x)
let su=x=>t_.decode(x),t_=new TextDecoder("utf-8"),us=x=>_t.encode(x),_t=new TextEncoder("utf-8")
let U=(o,n)=>new Uint8Array(K.memory.buffer,o,n),lo=x=>Number(BigInt.asUintN(32,x))
let KC=(x,r)=>(x=("string"==typeof x?us(x):x),r=K.mk(18,x.length),U().set(x,lo(r)),r)
let Ks=s=>K.sc(KC(s))
let asn=(s,y)=>K.dx(K.Asn(Ks(s),y))

let nyi=_=>{throw new Error("nyi")}
let fb

let I={"env":{Exit:Deno.exit,Args:nyi,Arg:nyi,ReadIn:nyi,Native:nyi,
 Read:(x,y,z)=>z?(U(z,fb.length).set(fb),0):(fb=readfile(su(U(x,y))),fb.length),
 Write:(f,nf,s,n)=>(nf?nyi():Deno.stdout.writeSync(U(s,n))),
}}

let kt=(Deno.build.os=="windows"?"c:/":"/c/")+"k/ktye.github.io/"
let kw=readfile(kt+"k.wasm")
let K=new WebAssembly.Instance(new WebAssembly.Module(kw),I).exports
K.kinit()

let cols=Deno.env.get("COLUMNS"),lines=Deno.env.get("LINES")
asn("l.",K.Atx(Ks("lxy"),K.Val(KC((cols?cols:80)+" "+(lines?lines:20)))))

let lk=x=>K.dx(K.Val(KC(us(su(readfile(x)).split("\n\\")[0]))))
let lf=x=>K.dx(K.Asn(KS(x),su(readfile(x))))
let lm=_=>"fft lu qr stats svd".split(" ").forEach(x=>lk(kt+"klib/"+x+".k"))
argv.forEach(x=>"-e"==x?Deno.exit(0):"-m"==x?lm():x.endsWith(".k")?lk(x):lf(x))

console.log("ktye/k(deno)",lines,cols)
let _k=1,repl=_=>{
 let buf=new Uint8Array(256),p=new Uint8Array([32]);
 while(1){Deno.stdout.writeSync(p);let n=Deno.stdin.readSync(buf);if(!n)return;
  if(n&&buf[n-1]==10)--n;if(n&&buf[n-1]==13)--n;
  let x=buf.subarray(0,n);
  try{let r=K.mk(18,x.length);U().set(x,lo(r));K.repl(r) }
  catch(e){console.log(e)}}}

repl()

