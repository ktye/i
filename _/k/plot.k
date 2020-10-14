
plot: {[x;y]drawplots`xy!,(x;y)}
polar:{[z]  drawplots`polar!,,z }
axis: {[x]`axis`nyi}

drawplots:{[d]dplot plots d}
replot:{[]dplot saveplots}
plots:{[d]saveplots::$[0~@f:(`xy`polar!(plotxy;plotpo))@!d;,f;f].'.d}

plotsize:   800 600
plotfont:   "20px monospace"
fontwidth:  10                         /or measure text dynamically in js
plotcolors: ("#8dd3c7";"#feffb3";"#bfbbd9";"#fa8174";"#81b1d2";"#fdb462";"#b3de69";"#bc82bd";"#ccebc4";"#ffed6f") /("#003FFF";"#03ED3A";"#E8000B";"#8A2BE2";"#FFC400";"#00D7FF")
plotfg:     "white"

dplot:{[p]plotsize 'd,/drawp./:&(p;rects[0 0,plotsize;#p])}
rects:{[r;n]r[2]%:n;((!n)*r 2),\:1_r}


plotxy:{[x;y]a:autoxy.d:xydat[x;y];x:d 0;y:d 1
 lines:{`x`y`type`size!(x 0;x 1;`line;2;x 2)}'&(x;y;(#x)#plotcolors)
 `type`title`xlabel`ylabel`axis`lines!(`xy;"- plot -";"x";"y";a;lines)}
plotpo:{[z];p:plotxy[!z;_z];p[`type]:`polar;p[`axis]:,/(-;+;-;+)@'|/+p`axis;p}

xydat: {x:3$'$[~6~@x;,x;x];((#y)#x;y:3$'$[~6~@y;,y;y])}
podat: {xydat[!x;_x]}
autoxy:{,/(&/&/x;|/|/x;&/&/y;|/|/y)}
autopo:{,/(-;+;-;+)@'|/+autoxy[x;y]}



/plt: 
/ `type    `xy                                 `xy|`polar
/ `title   "plot title"
/ `xlabel  "time [s]"
/ `ylabel  "value"
/ `axis    [0 10 -2 2] /xmin xmax ymin ymax
/ `lines   (line1;line2;...)
/   `x     0.0 0.1 0.2 0.3
/   `y     -0.1 2.1 -0.3 -0.4
/   `type  `line | `point | `bar
/   `size  2
/   `color "#ff0000"

pplot:{[x;y;ax]x:3$'$[~6~@x;,x;x];y:3$'$[~6~@y;,y;y];x:(#y)#x
 xlabel:"x";ylabel:"y";ax:$[0~ax;,/(&/&/x;|/|/x;&/&/y;|/|/y);ax]
 lines:{`x`y`type`size!(x 0;x 1;`line;2;x 2)}'&(x;y;(#x)#plotcolors)
 `type`title`xlabel`ylabel`axis`lines!(`xy;"- plot -";"x";"y";ax;lines)}
polarp:{[p];p[`type]:`polar;p[`axis]:,/(-;+;-;+)@'|/+p`axis;p}

drawp: {[p;r]$[`xy~p`type;drawxy;drawpo][p;r]}
drawxy:{[p;r]
 o:`font`textBaseline`fillStyle!(plotfont;"bottom";plotfg)
 r:clip[r;left:5;right:5;top:20;bottom:20] /left:5;right:5;top:20;bottom:20
 o,:rtext[rr:+/r 0 2;bottom++/r 1 3;((p`xlabel),"  "), 'k(p`axis)0 1]
 o,:ltext[r 0;r 1;((p`ylabel),"  ", 'k(p`axis)2 3)]
 o,:rtext[rr;r 1;p`title]
 o,:`beginPath`lineWidth`strokeStyle`rect`stroke!(();2;plotfg;r;())
 line:drawline[;p`axis;r]
 ,/o,line'p`lines}
 
drawpo:{[p;rect]o:`font`textBaseline`fillStyle`strokeStyle`lineWidth!(plotfont;"bottom";plotfg;"gray";2)
 left:5;right:5;top:20;bottom:20;rect:squareclip clip[rect;left;right;top;bottom]
 o,:ltext[rect 0;ty:rect 1;(p`title)] /," ",$|/p`axis]
 o,:ctext[(rect 0)+2\rect 3;ty;"0"],rtext[+/rect 0 2;ty;$|/p`axis]
 o,:strokeLine.(,/+\rect 0 2),\:(rect 1)+2\rect 3
 o,:strokeLine.((rect 0)+2\rect 2),/:+\rect 1 3
 o,:`strokeStyle`beginPath`arc`stroke!(plotfg;();arct rect;())
 line:drawline[;p`axis;rect]
 ,/o,line'p`lines}
 

drawline:{[l;axis;rect]o:`beginPath`lineWidth`strokeStyle!(();l`size;l`color)
 o,:`lineTo!(scale[l`x;axis 0 1;,/+\rect 0 2]),'scale[l`y;axis 2 3;|,/+\rect 1 3]
 o,:`stroke!,()
 o}
 
clip:{[rect;left;right;top;bottom]rect+:,/(left;top;-left+right;-top+bottom)}
squareclip:{[r]s:(r 2)&r 3;r+,/(2\(r 2 3)-s;s-r 2 3)}
arct:{[r](R+r 0 1),(R:0.+2\r 2),0 2p}                 /arc arguments (within square rect)
strokeLine:{[a;b]`beginPath`moveTo`lineTo`stroke!(();a;b;())}

scale:{[x;from;to](*to)+(x-*from)*(-/to)%-/from}      /scale[1 2 3.;-5 5.;0 100]
ctext:{[x;y;s]`fillText!,(s;x-2\w;y;w:fontwidth*#s)}
ltext:{[x;y;s]`fillText!,(s;x;y;fontwidth*#s)}
rtext:{[x;y;s]`fillText!,(s;x-w;y;w:fontwidth*#s)}

sin:129 'F;cos:130 'F
/x:!20
/y:(sin@0.3*x;cos@0.3*x)
/plot[x;*y]
/plot[x;y]
/drawplots`xy`xy!((x;*y);(x;y))

/z:1.1a20 1.2a30 1.4a40 1.5a50 1.6a60
/polar z
/polar (z;-z)
/drawplots`polar`polar!(z;(z;-z))

testPlot:{[]x:!20;y:(sin@0.3*x;cos@0.3*x);z:1.1a20 1.2a30 1.4a40 1.5a50 1.6a60
 drawplots`xy`polar`polar`xy!((x;*y);,z;,(z;-z);(x;y))}
