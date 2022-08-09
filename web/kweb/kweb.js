import { K } from '../k.js'


let nodes={}

function update(){
 for(let id in nodes){
  let dst=nodes[id]
  if(dst.offsetParent===null)continue   //skip invisible
  if(dst.s){ 
   let k=K.Kx(".",dst.s)
   if(k!=dst.k)gui(id,null,dst.s,null)  //global variable has changed
  }
  if(dst.e)    gui(id,null,null,dst.e)  //evaluate expr
}}


function show(id,x){return gui(K.sK(id),   x,   null,null)}
function link(id,x){return gui(K.sK(id),null,   x   ,null)}
function expr(id,x){return gui(K.sK(id),null,null,K.CK(x))}

function gui(id,x,s,e){id=(id=="")?"uid"+String(Object.keys(nodes).length):id
 let dst=ge(id);
 if(dst==null){dst=ce("div");dst.id=id;document.body.appendChild(dst)}
 dst.s=s;dst.e=e                                 //symbol,expr
 if(('k' in dst)&&dst.k)K.unref(dst.k); dst.k=x  //k-value
 removeAll(dst)
 nodes[id]=dst
 if(s){dst.k=K.Kx(".",     s )}                  //link
 if(e){dst.k=K.Kx(".",K.KC(e))}                  //expr
 let t=K.TK(dst.k)
 switch(t){
 case"B":case"I":case"S":case"F":case"Z":case"L":
         Lshow(dst,dst.k);break
 case"T":Tshow(dst,dst.k);break
 case"D":Dshow(dst,dst.k);break
 default:console.log("gui:type",t)
 }
 return K.Ks(id)
}



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

function Lshow(dst,x){ //create select element from vectors
 let n=K.NK(x)
 let S=Array(n)
 if(K.TK(x)=="L"){
  for(let i=0;i<n;i++)S[i]=K.Kx("`k@",K.Kx("@",K.ref(x),K.Ki(i)))
  K.unref(x)
 }else S=K.LK(K.Kx("$",x)).map(K.CK)
 let s=ce("select")
 for(let i=0;i<n;i++){
  let o=ce("option")
   o.textContent=S[i]
  s.appendChild(o)
 }
 dst.appendChild(s)
}
function Dshow(dst,x){ //create table from x(dict), only for symbol keys
 let [k,v]=K.LK(K.ref(x))
 if(K.TK(k)!="S"){K.unref(k);K.unref(v);return}
 let S=K.SK(k), n=S.length
 let ta=ce("table")
 ta.classList.add("kweb-dict")
 for(let i=0;i<n;i++){
  let tr=ce("tr")
   let th=ce("th")
    th.textContent=S[i]
   tr.appendChild(th)
   let td=ce("td")
    td.textContent=K.Kx("`k@",K.Kx("@",K.ref(v),K.Ki(i)))
   tr.appendChild(td)
  ta.appendChild(tr)
 }
 dst.appendChild(ta)
 K.unref(v)
 //todo: editable
}
function Tshow(dst,x){ //create table from x(table)
 let N=K.iK(K.Kx("#",K.ref(x)))
 let L=K.LK(K.ref(x))             //[keys,values]
 let S=K.SK(L[0])
 let ta=ce("table") 
 ta.classList.add("kweb-table")
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
 if(dst.s)editTable(dst)
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
    v=K.Kx(".",dst.k,K.Kx(",",K.Ks(t.cols[j]),K.Ki(i)),v)
    K.KA(dst.s,v)
    dst.k=K.ref(v)
    delete et.old
    update()
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

/*
function style(x){
 let s=ce("style");s.innerText=K.CK(x);
 document.head.appendChild(s);return BigInt(0)
}
*/

document.update=update //for custom updates from k

let kweb = {show,link,expr,update}
export { kweb }
