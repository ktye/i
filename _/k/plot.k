
/ p:pplot[x;y]  -> (x;y)data to plt dict
/ d:drawplot[p] -> plt to 2d canvas api dict
/ 300 200 'd d  -> draw on js canvas

plot: {plotsize 'd drawplot[,/(0;0;plotsize);pplot[x;y]]}
polar:{plot[();x]}

plotsize:   300 200
plotfont:   "20px monospace"
fontwidth:  10                         /or measure text dynamically in js
plotcolors: ("#003FFF";"#03ED3A";"#E8000B";"#8A2BE2";"#FFC400";"#00D7FF")

/ p:plot[x;y] -> plt dict type
/ 

/plt: 
/ `type    `xy                                 `polar | `absang
/ `title   "plot title"
/ `xlabel  "time [s]"
/ `ylabel  "value"
/ `axis    [0 10 -2 2] /xmin xmax ymin ymax
/ `lines   (line1;line2;...)
/   `x     0.0 0.1 0.2 0.3
/   `y     -0.1 2.1 -0.3 -0.4                  1a20 2a30 4a40
/   `type  `line | `point | `bar
/   `size  2
/   `color "#ff0000"


pplot:{[x;y]xlabel:"x";ylabel:"y";axis:,/(&/x;|/x;&/y;|/y)
 line:`x`y`type`size`color!(x;y;`line;2;plotcolors 0)
 `type`title`xlabel`ylabel`axis`lines!(`xy;"";"x";"y";axis;,line)}
 

drawplot:{[rect;p];o:`rect`clip!(rect;())
 left:50;right:0;top:20;bottom:20
 o,:`font`textBaseline`fillText!(plotfont;"bottom";(p`xlabel;2\rect[2]-fontwidth*#p`xlabel;rect 2))
 rect[0]+:left;rect[1]+:top;rect[2]-:left+right;rect[3]-:top+bottom      /translate to plot axis
 l:p[`lines;0]
 o,:`rect`clip`lineWidth`strokeStyle!(rect;();l`size;l`color)
 o,:`lineTo!(l`x),'l`y
 o,:`stroke!,()
 o}


 \drawplot[0 0 300 200;pplot[!5;-1 1 0 2 -1]]

\
