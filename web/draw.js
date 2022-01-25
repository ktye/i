import { K } from './k.js'

let D = {} // draw show

// display an image on a canvas with id="_cnv"
// show(20;(20*50)#255)
D.show = function(x){
 unhandle()
 if(K.TK(x)=='i') return x  // hide image: show 0
 return setimg(img(x))
}

function setimg(x){
 let cnv = ge("_cnv"); cnv.width=x.w; cnv.height=x.h
 let ctx = cnv.getContext("2d")
 let d = ctx.createImageData(x.w, x.h)
 let u = new Uint8Array(d.data.buffer)
 u.set(new Uint8Array(x.I.buffer))
 ctx.putImageData(d,0,0)
 cnv.style.display = "block"
 return K.Ki(x.w*x.h)
}

let click, zoom // k-values: callback functions

function unhandle(){
 if(click!=undefined){K.unref(click); click=undefined}
 if(zoom !=undefined){K.unref(zoom);  zoom =undefined}
 let cnv = ge("_cnv")
 cnv.ondblclick  = undefined
 cnv.onmousedown = undefined
 cnv.onmouseup   = undefined
 cnv.style.display = "none"
}


//im:{(20;(20*50)#*1?255*256*256)};showev[im[];{`click \(x;y);im[]};{[x;y;w;h]`zoom \(x;y;w;h);im[]}]
D.showev = function(x,y,z){
 let r=D.show(x); click=y; zoom=z
 
 let cnv = ge("_cnv");
 let ctx = cnv.getContext("2d")
 console.log("set ondblclick")
 cnv.ondblclick = function(ev){
  // console.log("dblclick", ev.offsetX, ev.offsetY)
  let x=K.Kx(".", K.ref(click), K.KI([ev.offsetX, ev.offsetY]))
  K.unref( (K.TK(x)==="L") ? setimg(img(x)) : x )
 }
 
 let zd = false, zm, zs, ze, bg
 zs = function(ev){ //zoom-start
  bg = ctx.getImageData(0,0,cnv.width,cnv.height)
  cnv.onmousemove = zm
  cnv.style.cursor="crosshair"
  zd = [ev.offsetX,ev.offsetY,0,0]
 }
 zm = function(ev){ //zoom-move
  if(zd!==false){
   zd=[zd[0],zd[1],ev.offsetX-zd[0],ev.offsetY-zd[1]]
   ctx.putImageData(bg,0,0)
   ctx.beginPath()
   ctx.rect(...zd)
   ctx.strokeStyle='red'
   ctx.stroke()
  }
 }
 ze = function(ev){ //zoom-end
  ctx.putImageData(bg,0,0)
  cnv.style.cursor=""
  cnv.onmousemove = undefined
  if(Math.abs(zd[2])< 5||Math.abs(zd[3])<5){zd=false;return}
  if(zd[2]<0){ zd[0]+=zd[2]; zd[2]=-zd[2] }
  if(zd[3]<0){ zd[1]+=zd[3]; zd[3]=-zd[3] }
  
  console.log("zoom", zd)
  
  let x=K.Kx(".", K.ref(zoom), K.KI(zd))
  K.unref( (K.TK(x)==="L") ? setimg(img(x)) : x )
  
  zd = false
  return
 }
 cnv.onmousedown = zs
 cnv.onmousemove = zm
 cnv.onmouseup   = ze
 
 return r
}

function img(x){
 if(K.TK(x) != 'L') return K.KE("img: type")
 if(K.NK(x) != 2)   return K.KE("img: L2")
 let r = K.LK(x)
 if(K.TK(r[0]) != 'i') return K.KE("img: h-type")
 if(K.TK(r[1]) != 'I') return K.KE("img: I-type")
 let h = K.iK(r[0])
 let I = K.IK(r[1])
 let w = Math.floor(I.length / h)
 if(I.length != w*h) return K.KE("img:rect")
 let u = new Uint8Array(I.buffer)
 for(let i=3;i<u.length;i+=4)u[i]=255 //alpha
 return {w:w, h:h, I:I}
}

function ge(x){return document.getElementById(x)}
function ce(x){return document.createElement(x)}

// cnv=document.getElementById("_cnv");ctx=cnv.getContext("2d")



export { D }
