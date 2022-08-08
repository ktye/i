import { K } from '../k.js'

//nodes stores all element ids that 
let nodes={}

function update(){
 for(id in nodes){
  let dst=nodes[i]
  if('s'in dst){
   let k=K.Kx(".",dst.s)
   if(k!=dst.k)show(id,dst.s) //global variable has changed
}}}

//show(dst,x) shows the k-value x in html element with id.
function show(id,x){
 let dst=ge(id)
 if(dst===null){
  dst=ce("div")
  document.body.appendChild(dst)
 }
 removeAll(dst)
 delete dst.s      //k-symbol(global variable)
 if('k'in dst){
  K.unref(dst.k)
  delete dst.k     //k-value
 }
 if(!x)return
 if(K.TK(x)=="s"){ //backed by a global
  nodes[id]=x
  dst.s=x;
  x=K.Kx(".",x)
 }else delete nodes[id]
 dst.k=x
 let t=K.TK(x)
 switch(K.TK(x)){
 case "T":Tshow(dst,x);break;
 default:console.log("type",t)
}}

function SV(x){ //strings from vector 
 let l
 switch(K.TK(x)){
 case 'B':case 'I':case 'S':case 'F':case 'Z':l=K.Kx("$",x);break
 default:{console.error("nyi");return[]}
 }
 let r=[]
 let n=K.NK(x)
 for(let i=0;i<n;i++)r.push(K.CK(K.Kx("@",K.ref(l),K.Ki(i))))
 K.unref(l);return r
}

function Tshow(dst,x){ //create table from x, append to dst
 let N=K.iK(K.Kx("#",K.ref(x)))
 let L=K.LK(K.ref(x))             //[keys,values]
 let S=K.SK(L[0])
 let ta=ce("table") 
 ta.cols=S                        //store column names/types
 ta.coltype=Array(S.length)
  let tr=ce("tr")
   for(let i=0;i<S.length;i++){let th=ce("th");th.textContent=S[i];tr.appendChild(th)}
  ta.appendChild(tr)
  let V=[]
  for(let i=0;i<S.length;i++){
   let v=K.Kx("@",K.ref(L[1]),K.Ki(i))
   ta.coltype[i]=K.TK(v)
   V[i]=SV(v)
  }
  K.unref(L[1])
  for(let i=0;i<N;i++){
   let tr=ce("tr")
   tr.i=i                         //store row index
   for(let j=0;j<S.length;j++){
    let td=ce("td")
     td.textContent=V[j][i]
    tr.appendChild(td)
   }
   ta.appendChild(tr)
  }
 dst.appendChild(ta)
 if('v'in dst)editTable(dst)
}

function editTable(dst){ //make table editable (link with k)
 dst.querySelectorAll("td").forEach(x=>{
  x.contentEditable=true
  x.onbeforeinput=function(e){let et=e.target
   if(!('old' in et))et.old=et.textContent
  }
  x.onkeydown=function(e){let et=e.target
   if(e.key=="Enter"){
    let i=et.parentElement.i,j=et.cellIndex
    let t=et.parentElement.parentElement
    let ct=t.coltype[j]
    let s=K.Ks(ct.toLowerCase())
    let v=K.Kx("$",s,K.KC((ct=="S")?'`"'+et.textContent+'"':et.textContent))
    if(v==0){x.classList.add("kweb-invalid");return false}
    et.textContent=K.CK(K.Kx("$",K.ref(v)))
    et.classList.remove("kweb-editing","kweb-invalid")
    et.blur()
    dst.k=K.Kx(".",dst.k,K.Kx(",",K.Ks(t.cols[j]),K.Ki(i)),v)
    if('v'in dst)dst.k=K.Kx(":",dst.v,dst.k)
    delete et.old
    return false
   }else if(e.key=="Escape"){
    if('old' in et)et.textContent=et.old
    delete et.old
    et.blur()
    et.classList.remove("kweb-editing","kweb-invalid")
    return false
   }
   et.classList.add("kweb-editing")
  }
})}

function ge(x){return document.getElementById(x)}
function ce(x){return document.createElement(x)}
function removeAll(p){while(p.firstChild)p.removeChild(p.firstChild)}

function body(x){
 let b=ge("main");
 ((b==null)?document.body:b).innerHTML=K.CK(x);return BigInt(0)
}
function style(x){
 let s=ce("style");s.innerText=K.CK(x);
 document.head.appendChild(s);return BigInt(0)
}

let kweb = {show,body,update,style}
export { kweb }
