
/ p:pplot[x;y]  -> (x;y)data to plt dict
/ d:drawplot[p] -> plt to 2d canvas api dict
/ 300 200 'd d  -> draw on js canvas

plot: {plotsize 'ddrawplot[,/(0;0;plotsize);pplot[x;y]]}
polar:{plot[();x]}

plotsize:   800 600
plotfont:   "20px monospace"
fontwidth:  10                         /or measure text dynamically in js
plotcolors: ("#8dd3c7";"#feffb3";"#bfbbd9";"#fa8174";"#81b1d2";"#fdb462";"#b3de69";"#bc82bd";"#ccebc4";"#ffed6f") /("#003FFF";"#03ED3A";"#E8000B";"#8A2BE2";"#FFC400";"#00D7FF")
plotfg:     "white"
/ p:plot[x;y] -> plt dict type
/ 

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


pplot:{[x;y]x:3$'$[~6~@x;,x;x];y:3$'$[~6~@y;,y;y];x:(#y)#x
 xlabel:"x";ylabel:"y";axis: \,/(&/&/x;|/|/x;&/&/y;|/|/y)
 lines:{`x`y`type`size!(x 0;x 1;`line;2;x 2)}'&(x;y;(#x)#plotcolors)
 `type`title`xlabel`ylabel`axis`lines!(`xy;"";"x";"y";axis;lines)}

drawplot:{[rect;p];o:`beginPath`rect`clip!(();rect;())
 left:50;right:0;top:20;bottom:20
 o,:`font`textBaseline`fillStyle`fillText!(plotfont;"bottom";plotfg;(p`xlabel;2\rect[2]-fontwidth*#p`xlabel;rect 3))
 rect[0]+:left;rect[1]+:top;rect[2]-:left+right;rect[3]-:top+bottom      /translate to plot axis
 line:drawline[;p`axis;rect]
 ,/o,line'p`lines}

drawline:{[l;axis;rect]o:`beginPath`lineWidth`strokeStyle`rect`stroke`clip!(();2;plotfg;rect;();())
 o,:`beginPath`lineWidth`strokeStyle!(();l`size;l`color)
 o,:`lineTo!(scale[l`x;axis 0 1;,/+\rect 0 2]),'scale[l`y;axis 2 3;|,/+\rect 1 3]
 o,:`stroke!,()
 o}

scale:{[x;from;to]to[0]+x*(to[1]-to[0])%(from[1]-from[0])}      /scale[1 2 3.;-5 5.;0 100]


sin:129 'F;cos:130 'F
x:!20
y:(sin@0.3*x;cos@0.3*x)
plot[x;y]
