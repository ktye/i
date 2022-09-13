import {kweb} from './kweb.js'
import {K} from '../k.js'

function ce(x) { return document.createElement(x) }
function ge(x) { return document.getElementById(x) }


function treeclick(e){
 K.Kx("clicktree",K.Ki(e.target.parentNode.rowIndex-1))
 kweb.update()
 return false
}

kweb.register("treetable",function(dst,x){
 let l=K.LK(x)
 let D=K.LK(l[0]),T=K.SK(l[1]),P=K.IK(l[2]),I=K.IK(l[3]),S=K.SK(l[4])
 let a=D.map(x=>K.CK(x).replaceAll("|","│").replace("+","├").replace("L","└").replace("-","─").replace("T","┬"))
 let t=ce("table")
 let h=ce("thead");h.innerHTML="<tr><th></th><th>T</th><th>P</th><th>I</th><th>S</th></tr>"
 t.appendChild(h)
 for(let i=0;i<a.length;i++){
  let r=ce("tr")
  let d=ce("td");d.textContent=a[i];r.appendChild(d)
  let s=ce("td");s.textContent=T[i];r.appendChild(s)
  let p=ce("td");p.textContent=P[i];r.appendChild(p)
  let k=ce("td");k.textContent=(I[i]==-2147483648)?"0N":I[i];r.appendChild(k)
  let v=ce("td");v.textContent=S[i];r.appendChild(v)
  r.ondblclick=treeclick
  t.appendChild(r)
 }
 dst.appendChild(t)
})

function setinput(s){ 
 kweb.init({start:[s+".k","compile.k"],post:compile})
}

let sel=ge("sel")
let examples="asn`blank`cal`cast`cli`cnd`const`cont`drp`fun`heap`inc`label`lit`loop`loop2`mem`swtch`tab`k".split("`")
for(let i=0;i<examples.length;i++){
 let o=ce("option");o.textContent=examples[i]
 sel.appendChild(o)
}
let h=decodeURIComponent(window.location.hash.slice(1)).split(" ") //e.g. #go asn
if(h.length){
 ge("target").selectedIndex=Math.max(0,["help","go","wa","ik"].indexOf(h[0]))
 sel.selectedIndex=Math.max(0,examples.indexOf(h[1]))
}else sel.selectedIndex=0
sel.onchange=function(e){setinput(examples[e.target.selectedIndex])}


setinput(examples[sel.selectedIndex])

function compile(){
 let t=ge("target").value 
 if(t=="help"){
  fetch("compile.help").then(r=>r.text()).then(r=>ge("out").textContent=r)
  return
 }
 fetch(t+".k").then(r=>r.text()).then(r=>{
  kweb.ktry(r)
  ge("out").textContent=K.CK(kweb.ktry(t+"``nort"))
 })
}
ge("target").onchange=compile