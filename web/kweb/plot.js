import { K } from '../k.js'

function ce(x){return document.createElement(x)}

export default function plot(dst,x){
 if(plotdict==0){
  K.KA(K.Ks("tics."),K.Kx(".",K.KC(tics_))) // `tics@
  plotdict=K.Kx(".",K.KC(plotdict_))
  plotwh  =K.Kx(".",K.KC(plotwh_))
 }
 K.unref(x)
 let cnv=ce("canvas")
 if("width" in dst.dataset)cnv.width =dst.dataset.width
 if("height"in dst.dataset)cnv.height=dst.dataset.height
 let dbl=ce("pre")                                        //dbl-click output
 dst.appendChild(cnv);dst.appendChild(dbl)
 window.requestAnimationFrame(function(){replot(cnv,[])}) //draw when mapped (and w/h is set)
}

function replot(cnv,ax){
 let x=K.ref(cnv.parentElement.k)
 let w=cnv.width,h=cnv.height              //default 300,150; show[`type`width`height!(`plot;400;300);..]
 let fs=('fs'in cnv)?cnv.fs:20             //show[`type`fs!(`plot;30);..]
 let ctx=cnv.getContext("2d")
 x=K.Kx("@",K.ref(plotdict),x)
 if(ax.length)x=K.Kx("@[;`a`t;]",x,K.KL([K.KF(ax),K.Ks("xy")])) //set axes after zoom
 let t=K.sK(K.Kx("@",K.ref(x),K.Ks("t")))      //"xy" or "po"
 let a=K.FK(K.Kx("@",K.ref(x),K.Ks("a")))      //[xmin xmax ymin ymax]
 let A=[fs,w-fs,h-fs,fs]                       //rect
 let C=[w/2,h/2],R=Math.min(w/2,h/2-fs)
 let f=(t=="xy")?xyclick(a,A):poclick(a[3],C,R)
 x=K.Kx(".",K.ref(plotwh),K.KL([x,K.Ki(fs),K.Ki(w),K.Ki(h)]))
 //x=K.Kx("plotwh",x,K.Ki(fs),K.Ki(w),K.Ki(h))
 draw(ctx,x,w,h)
 cnv.ondblclick=dblclick(cnv.nextSibling,ctx,f)
 cnv.onmousedown=zoom(cnv,ctx,a,A)
}



// draw(ctx,L,w,h)          // draw[(`color;255*256;`rect;10 10 300 200);400 300]
// `color;rgb               // i
// `font;"20px monospace"   // C
// `linewidth;w             // i or f
// `rect;(x;y;w;h)          // I or F        stroke
// `Rect;(x;y;w;h)          // I or F        fill
// `circle;(x;y;r)          // I or F        stroke
// `Circle;(x;y;r)          // I or F        fill
// `clip;(x;y;r)or(x;y;w;h) // I or F
// `line;(x0;y0;x1;y1)      // I or F
// `poly;(X;Y)              // X, Y: I or F  stroke
// `Poly;(X;Y)              // I or F        close fill
// `text;(x;y;"text")       // i,i,C
// `Text;(x;y;"text")       // i,i,C         rotated
// draw returns an image object (h;I)
function draw(ctx,x,w,h){
 let l=K.LK(x)
 let n=K.NK(x)
 if(n%2 != 0) return K.KE("draw: #x")
 
 ctx.fillStyle = "white"
 ctx.fillRect(0,0,w,h)
 ctx.fillStyle = "black"
 ctx.strokeStyle = "black"
 ctx.font = "20px monospace"
 
 let state = {}
 let saveState      = function(){
  state.fillStyle   = ctx.fillStyle
  state.strokeStyle = ctx.strokeStyle
  state.font        = ctx.font
  state.lineWidth   = ctx.lineWidth
 }
 let restoreState = function(){
  ctx.fillStyle   = state.fillStyle
  ctx.strokeStyle = state.strokeStyle
  ctx.font        = state.font
  ctx.lineWidth   = state.lineWidth
 }
 saveState()
 
 // reset clip needs save/restore which also resets all other changes
 let resetClip = function(){
   ctx.restore()
   ctx.save()
   restoreState()
 }
 
 let ck = function(s,c){ if(c==false) K.KE(s+" arg") }
 
 let setcolor = function(i){ 
  let s="rgb(" + String(i&255) + ", " + String((i>>>8)&255) + ", " + String((i>>>16)&255) + ")"
  ctx.fillStyle=s; ctx.strokeStyle=s; saveState()
 }
 
 let num = function(x){ let t=K.TK(x); return (t=='f') ? K.fK(x) : (t=='i') ? K.iK(x) : K.KE("num type") }
 let vec = function(x){ let t=K.TK(x); return (t=='F') ? K.FK(x) : (t=='I') ? K.IK(x) : K.KE("vec type") }
 
 ctx.save() //for clipping
 let cmd=function(s,a){
  let n
  switch(s){
   case "color":
    setcolor(num(a))
    break
   case "font":
    ck(s, ((K.TK(a)=="L") && (K.NK(a) == 2)))
    a = K.LK(a)
    ck(s, K.TK(a[0]) == "C")
    let px = num(a[1])
    ctx.font = String(px) + "px " + K.CK(a[0]); saveState()
    break
   case "linewidth":
    ctx.lineWidth=num(a); saveState()
    break
   case "rect":
   case "Rect":
    if(s=="rect") ctx.strokeRect(...vec(a))
    else          ctx.fillRect(...vec(a))
    break
   case "circle":
   case "Circle":
    a = vec(a)
    ck(s, a.length==3)
    ctx.beginPath()
    ctx.arc(a[0], a[1], a[2], 0, 2 * Math.PI)
    if(s=="circle") ctx.stroke()
    else            ctx.fill()
    break
   case "clip":
    a = vec(a)
    resetClip()
    ctx.beginPath()
    if(a.length==3) ctx.arc(a[0], a[1], a[2], 0, 2 * Math.PI)
    else{
     ck(a.length==4)
                    ctx.rect(a[0], a[1], a[2], a[3])
    }
    ctx.clip()
    break
   case "line":
    a = vec(a)
    ck(s, a.length==4)
    ctx.beginPath()
    ctx.moveTo(a[0], a[1])
    ctx.lineTo(a[2], a[3])
    ctx.stroke()
    break
   case "poly":
   case "Poly":
    ck(s, (K.TK(a)=="L" && K.NK(a) == 2))
    a = K.LK(a)
    let x = vec(a[0]), y = vec(a[1])
    if((x.length>1)&&(x.length==y.length)){
     ctx.beginPath()
     ctx.moveTo(x[0], y[0])
     for(let i=1;i<x.length;i++) ctx.lineTo(x[i], y[i])
     if(s=="poly") ctx.stroke()
     else         {ctx.closePath(); ctx.fill()}
    }
    break
   case "text":
   case "Text":
    ck(s, ((K.TK(a)=="L") && (K.NK(a) == 3)))
    a = K.LK(a)
    let xy = vec(a[0])
    ck(s, xy.length==2)
    ck(s, K.TK(a[1]) == 'i')
    ck(s, K.TK(a[2]) == 'C')
    align(ctx, num(a[1]))
    if(s=="Text"){ ctx.save(); ctx.translate(xy[0], xy[1]); ctx.rotate(-Math.PI/2); xy[0]=0; xy[1]=0 }
    ctx.fillText(K.CK(a[2]), xy[0], xy[1])
    if(s=="Text"){ ctx.restore() }
    break
   default:
    K.KE("draw: command: "+s)
    break
   }
 }
 
 for(let i=0;i<n;i+=2){
  if(K.TK(l[i]) != 's') return K.KE("draw cmd type")
  cmd(K.CK(K.Kx("$",l[i])), l[1+i])
 }
}

function align(ctx, a){
 ctx.textBaseline  = ["bottom","bottom","bottom","middle","top","top","top","middle","middle"][a]
 ctx.textAlign     = ["left","center","right","right","right","center","left","left","center"][a]
}

function dblclick(dbl,ctx,f){
 return function(ev){
  let X=ev.offsetX, Y=ev.offsetY
  ctx.beginPath();ctx.arc(X,Y,3,0,2*Math.PI);ctx.fillStyle="red";ctx.fill()
  //let x=K.Kx(".", K.ref(click), K.KI([X, Y]))
  //K.unref( (K.TK(x)==="L") ? setimg(img(x)) : x )
  dbl.textContent=String(f(X,Y))
 }
}

function xyclick(a,A){return function(x,y){
 let xy=xyscale([x,y],a,A)
 return "x:"+String(xy[0])+" y:"+String(xy[1])
}}
function poclick(r,C,R){return function(x,y){
 let X=(x-C[0])*r/R,Y=-(y-C[1])*r/R
 let abs=Math.hypot(X,Y)
 let ang=Math.atan2(X,Y)/Math.PI*180
 return String(abs)+"a"+String((ang<0)?360+ang:ang)
}}
function scale(x,x0,x1,y0,y1){return y0+(y1-y0)*(x-x0)/(x1-x0)}
function xyscale(xy,a,A){return [scale(xy[0],A[0],A[1],a[0],a[1]),scale(xy[1],A[2],A[3],a[2],a[3])]}


function zoom(cnv,ctx,a,A){
 let zd = false, zs, bg
 let zm = function(e){ //zoom-move
  if(zd!==false){
   zd=[zd[0],zd[1],e.offsetX-zd[0],e.offsetY-zd[1]]
   ctx.putImageData(bg,0,0)
   ctx.beginPath()
   ctx.rect(...zd)
   ctx.strokeStyle='red'
   ctx.stroke()
  }
 }
 let ze = function(ev){ //zoom-end
  ctx.putImageData(bg,0,0)
  cnv.style.cursor=""
  cnv.onmousemove = undefined
  if(Math.abs(zd[2])< 5||Math.abs(zd[3])<5){zd=false;return}
  if(zd[2]<0){ zd[0]+=zd[2]; zd[2]=-zd[2] }
  if(zd[3]<0){ zd[1]+=zd[3]; zd[3]=-zd[3] }
  let xya=xyscale([zd[0],zd[1]],a,A)
  let xyb=xyscale([zd[0]+zd[2],zd[1]+zd[3]],a,A)
  zd = false
  replot(cnv,[xya[0],xyb[0],xyb[1],xya[1]])
  return
 }
 return function(e){
  bg = ctx.getImageData(0,0,cnv.width,cnv.height)
  cnv.onmousemove = zm
  cnv.onmouseup = ze
  cnv.style.cursor="crosshair"
  zd = [e.offsetX,e.offsetY,0,0]
}}

var plotdict=0
var plotwh=0
var tics=0
const tics_="{[minmax]nice:{[x;r]f:x%0.+10^ex:_log[10;x];(1 2 5 10.@*&(~f>1 2 5 0w;f<1.5 3 7 0w)[r])*10^ex};e:nice[-/|minmax;0];s:nice[e%4.;1];n:_1.5+e%s;$[~(minmax 1)>*-2#r:s*(_(*minmax)%s)+!n;-1_r;r]}"

// the following k code has been escaped with: sh ./escape.sh
const plotdict_="{[d]l:$!d;v:.d; t:$[2~#d;`xy;`polar];\n\
 y:$[t~`xy; $[`L~@y:v 1;y;,y];          $[`L~@y:_*v;y;,y]]\n\
 x:$[t~`xy; $[`L~@x:v 0;x;(,x)@(#y)#0]; $[`L~@x:imag@*v;x;,x]]\n\
 xt:`tics(&/&/x;|/|/x);yt:`tics(&/&/y;|/|/y)\n\
 a:$[t~`xy;(xt 0;*-1#xt;yt 0;*-1#yt);(-a;a;-a;a:*|`tics@0.,|/|/abs@*v)]\n\
 c:c@(#c:11826975 950271 2924588 2631638 12412820 4937356 12744675 8355711 2276796 13614615)/!#x\n\
 style:$[t~`polar;\"..\";`i~@**y;\"||\";\"--\"] / -.| line points bar\n\
 size: $[t~`polar;2;style~\"||\";(--/((**x),-1#*x))%-1+#*x ;2]\n\
 lines:{`style`size`color`x`y!(style;size;z;x;0.+y)}'[x;y;c]\n\
 `L`T`t`l`a!(lines;\"\";t;l;a)}"
const plotwh_="{[x;fs;W;H]; ; a:x`a;T:x`T;grey:13882323\n\
 C:(W%2;H%2);R:(W%2)&(H%2)-fs\n\
 dst:$[`xy~x`t;(fs;W-fs;H-fs;fs);((C-R),C+R)0 2 3 1];rdst:(fs;fs;W-2*fs;H-2*fs)\n\
 xs:(a 0 1)(dst 0 1)' /transform axis to canvas\n\
 ys:(a 2 3)(dst 2 3)'\n\
 bars:{[l]$[\"|\"':l`style;(`color;l`color),,/{(`Rect;((-dx%2)+xs x;ys y;dx:-/xs(l`size;0.);(ys a 2)-ys y))}'[l`x;l`y];()]}\n\
 line:{[l]$[\"-\"':l`style;(`linewidth;l`size;`color;l`color;`poly;(xs l`x;ys l`y));()]}\n\
 dots:{[l]$[\".\"':l`style;(`color;l`color),,/{(`Circle;(xs x;ys y;1.5*l`size))}'[l`x;l`y];()]}\n\
 c:(`clip;(0;0;W;H);`font;(\"monospace\";fs);`color;0;`text;((W%2;fs);1;T))\n\
 xy:{[]c,:(`text;((fs;H);0;$a 0);`text;((W%2;H);1;(x`l)0);`text;((W-fs;H);2;$a 1))\n\
       c,:(`Text;((fs;H-fs);0;$a 2);`Text;((fs;H%2);2;(x`l)1);`Text;((fs;fs);2;$a 3))\n\
       c,:(`color;0;`linewidth;2;`rect;rdst)      /todo: clip rdst\n\
       c,:(`linewidth;1;`color;grey)\n\
       c,:(`clip;rdst)\n\
       c,:,/{(`line;0.+(x;dst 2;x;dst 3))}'xs`tics x[`a;0 1]\n\
       c,:,/{(`line;0.+(dst 0;x;dst 1;x))}'ys`tics x[`a;2 3]}\n\
 po:{[]c,:(`text;((C 0;H);1;(x`l)0);`text;(C+.75*R;6;$(x`a)3))\n\
       c,:(`font;(\"monospace\";_fs*.8)),,/{(`text;(C+R*(_;imag)@'x;y;z))}'[1@270.+a;0 0 6 6 4 4 2 2;$a:30 60 120 150 210 240 300 330]\n\
       c,:(`color;0),/{(`line;,/+C+(R-fs%2;R)*/:(_;imag)@'x)}'1@30.*!12\n\
       /c,:(`clip;C,R) /bug in cairo?\n\
       c,:(`color;grey;`linewidth;1;`line;((-R)+*C;C 1;R+*C;C 1);`line;(*C;(-R)+C 1;*C;R+C 1))\n\
       c,:,/{(`circle;0.+C,x)}'r:(xs@`tics 0.,x[`a;3])-*C\n\
       c,:(`color;0;`linewidth;2;`circle;C,R)}\n\
 $[`xy~x`t;xy[];po[]]\n\
 c,:,/bars'x`L\n\
 c,:,/line'x`L\n\
 c,:,/dots'x`L}"
